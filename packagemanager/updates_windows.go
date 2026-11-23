package packagemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	"github.com/shirou/gopsutil/v4/host"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

// WindowsUpdatesManager implements WindowsManager for Windows Updates
type WindowsUpdatesManager struct{}

type PowerShellUpdateJson struct {
	Title          string   `json:"Title"`
	Description    string   `json:"Description"`
	KBArticleIDs   []string `json:"KBArticleIDs"`
	Categories     []string `json:"Categories"`
	SupportURL     string   `json:"SupportUrl"`
	RebootRequired bool     `json:"RebootRequired"`
	EulaAccepted   bool     `json:"EulaAccepted"`
	UpdateID       string   `json:"UpdateID"`
	RevisionNumber int      `json:"RevisionNumber"`
	// IDs of the classification and the products, the same on all systems, see securityCategoryIDs
	CategoryIDs []string `json:"CategoryIDs"`
	// Optional updates which are only installed on request
	BrowseOnly bool `json:"BrowseOnly"`
}

// rebootRequiredKeys are the registry keys of Windows Update and the servicing stack (CBS). They only exist while a
// reboot is pending to finish the installation of updates and are removed by the reboot.
// PendingFileRenameOperations of the Session Manager is not checked on purpose: any application (installers, Office
// Click-to-Run, printer drivers, antivirus software, ...) uses it to replace files on the next boot and it is set
// again shortly after the reboot, so it reported a pending reboot on almost every system (#124).
// https://stackoverflow.com/a/47869761/11885414
var rebootRequiredKeys = []string{
	`SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\RebootRequired`,
	`SOFTWARE\Microsoft\Windows\CurrentVersion\Component Based Servicing\RebootPending`,
}

// IDs of the update classifications and products of Windows Update, WSUS and Configuration Manager (SCCM).
// The names of the categories depend on the language of the system and the update server, the IDs are the same
// on all systems. Definition updates (Microsoft Defender) are always security updates.
const (
	classificationSecurityUpdates     = "0fa1201d-4330-4fa8-8ae9-b877473b6441"
	classificationCriticalUpdates     = "e6cf1350-c01b-414d-a61f-263d14d133b4"
	classificationDefinitionUpdates   = "e0789628-ce08-4437-be74-2495b842f43b"
	productMicrosoftDefenderAntivirus = "8c3fcc84-7410-4a95-8b89-a166a0190486"
)

var securityCategoryIDs = []string{
	classificationSecurityUpdates,
	classificationCriticalUpdates,
	classificationDefinitionUpdates,
	productMicrosoftDefenderAntivirus,
}

// sccmClientServiceKey exists if the Configuration Manager (SCCM) client is installed
const sccmClientServiceKey = `SYSTEM\CurrentControlSet\Services\CcmExec`

func (w WindowsUpdatesManager) ListInstalledApps(ctx context.Context) ([]WindowsApp, error) {
	apps, err := w.getInstalledApps()
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (w WindowsUpdatesManager) getInstalledApps() ([]WindowsApp, error) {
	// DO NOT use WMI and the Win32_Product class!
	// It will may repair all installed MSI packages on the system!
	// DO NOT EVEN THINK ABOUT IT!
	// https://support.microsoft.com/kb/974524
	// https://sdmsoftware.com/wmi/why-win32_product-is-bad-news/

	var apps []WindowsApp
	keys := []registry.Key{
		registry.LOCAL_MACHINE,
		registry.CURRENT_USER,
	}
	paths := []string{
		`SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall`,
		`SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall`,
	}

	// Keep a list of all found applications to avoid duplicates
	seenApps := make(map[string]bool)

	i := 0
	for _, key := range keys {
		for _, path := range paths {
			k, err := registry.OpenKey(key, path, registry.READ)
			if err != nil {
				continue
			}
			defer k.Close()
			names, _ := k.ReadSubKeyNames(-1)
			for _, name := range names {
				i++
				subKey, err := registry.OpenKey(key, path+"\\"+name, registry.READ)
				if err == nil {

					// Get name, version and install location
					displayName, _, err := subKey.GetStringValue("DisplayName")
					version, _, errVer := subKey.GetStringValue("DisplayVersion")
					// Publisher is optional. Many installers do not set it, and the app is still relevant for the inventory
					publisher, _, _ := subKey.GetStringValue("Publisher")

					// Avoid duplicates
					uniqueID := fmt.Sprintf("%s|%s|%s", displayName, version, publisher)
					if seenApps[uniqueID] {
						subKey.Close()
						continue
					}

					// Only add if we have a name
					if err == nil && displayName != "" && errVer == nil {
						seenApps[uniqueID] = true

						apps = append(apps, WindowsApp{
							Name:      displayName,
							Version:   version,
							Publisher: publisher,
						})
					}
					subKey.Close()
				}
			}
		}
	}

	return apps, nil
}

func (w WindowsUpdatesManager) ListAvailableUpdates(ctx context.Context) ([]WindowsUpdate, error) {
	output, err := w.getAvailableUpdates(ctx)
	if err != nil {
		return nil, err
	}

	updates, err := w.parsePowerShellUpdateSessionOutput(output)
	if err != nil {
		return nil, fmt.Errorf("error parsing Windows Updates output: %s", err)
	}

	return updates, nil
}

func (w WindowsUpdatesManager) getAvailableUpdates(ctx context.Context) (string, error) {
	// The file 'updates_windows_ole.txt' contains a implemtation using OLE/COM
	// how ever, this was way to complicated for what we try to achieve here.

	// PowerShell script with JSON output. It reads the updates of two sources:
	// 1. The Windows Update API (Windows Update or WSUS)
	// 2. The updates which Configuration Manager (SCCM) deployed to this computer. With SCCM the Windows Update API
	//    finds no updates, as SCCM does not approve the updates on its WSUS server (#116). The classification is read
	//    from the scan result of the SCCM client. The scan result also lists the updates which are not deployed
	//    (and some updates twice), so only the deployed updates are reported, like the approved updates of WSUS.
	// Only IDs and numbers are evaluated, the names and messages depend on the language of the system.
	psScript := `
		[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
		$Updates = @()
		$WuaError = $null

		try {
		    $Session = New-Object -ComObject Microsoft.Update.Session
		    $Searcher = $Session.CreateUpdateSearcher()
		    $Results = $Searcher.Search("IsInstalled=0 and IsHidden=0")
		    foreach ($Update in $Results.Updates) {
		        $categories = @()
		        $categoryIDs = @()
		        foreach ($cat in $Update.Categories) {
		            $categories += $cat.Name
		            $categoryIDs += $cat.CategoryID
		        }
		        $Updates += [PSCustomObject]@{
		            Title          = $Update.Title
		            Description    = $Update.Description
		            KBArticleIDs   = @($Update.KBArticleIDs)
		            Categories     = $categories
		            CategoryIDs    = $categoryIDs
		            SupportUrl     = $Update.SupportUrl
		            RebootRequired = $Update.RebootRequired
		            EulaAccepted   = $Update.EulaAccepted
		            BrowseOnly     = $Update.BrowseOnly
		            UpdateID       = $Update.Identity.UpdateID
		            RevisionNumber = $Update.Identity.RevisionNumber
		        }
		    }
		} catch {
		    $WuaError = $_.Exception.Message
		}

		# Configuration Manager (SCCM), the namespaces only exist if the SCCM client is installed
		$Deployed = @()
		try {
		    $Deployed = @(Get-CimInstance -Namespace 'root\ccm\ClientSDK' -ClassName CCM_SoftwareUpdate -ErrorAction Stop |
		        Where-Object { $_.ComplianceState -eq 0 })
		} catch {}

		if ($Deployed.Count -gt 0) {
		    $ScanResults = @{}
		    try {
		        Get-CimInstance -Namespace 'root\ccm\SoftwareUpdates\UpdatesStore' -ClassName CCM_UpdateStatus -ErrorAction Stop |
		            ForEach-Object { $ScanResults[$_.UniqueId] = $_ }
		    } catch {}

		    foreach ($Update in $Deployed) {
		        # Site_<site id>/SUM_<update id>
		        $updateID = (($Update.UpdateID -split '/')[-1]) -replace '^SUM_', ''
		        $categoryIDs = @()
		        $revisionNumber = 0
		        $scan = $ScanResults[$updateID]
		        if ($scan) {
		            $categoryIDs = @($scan.UpdateClassification, $scan.ProductID) | Where-Object { $_ }
		            $revisionNumber = $scan.RevisionNumber
		        }
		        $kbArticleIDs = @()
		        if ($Update.ArticleID) {
		            $kbArticleIDs = @($Update.ArticleID)
		        }
		        $Updates += [PSCustomObject]@{
		            Title          = $Update.Name
		            Description    = $Update.Description
		            KBArticleIDs   = $kbArticleIDs
		            Categories     = @()
		            CategoryIDs    = @($categoryIDs)
		            SupportUrl     = $Update.URL
		            # EvaluationState 8, 9 and 10: installed, waiting for the reboot
		            RebootRequired = @(8, 9, 10) -contains [int]$Update.EvaluationState
		            EulaAccepted   = $true
		            BrowseOnly     = $false
		            UpdateID       = $updateID
		            RevisionNumber = [int]$revisionNumber
		        }
		    }
		}

		if ($WuaError -and $Deployed.Count -eq 0) {
		    [Console]::Out.Write("Windows Update search failed: " + $WuaError)
		    exit 1
		}

		ConvertTo-Json -Compress -Depth 5 -InputObject @($Updates)
	`

	timeout := time.Duration(120 * time.Second)
	commandResult, err := utils.RunCommand(ctx, utils.CommandArgs{
		Timeout: timeout,
		Command: psScript,
		Shell:   "powershell_command",
	})

	if err != nil {
		return "", fmt.Errorf("error fetching Windows Updates via PowerShell: %s", commandResult.Stdout)
	}

	if commandResult.RC != 0 {
		return "", fmt.Errorf("error fetching Windows Updates via PowerShell, exit code %d: %s", commandResult.RC, commandResult.Stdout)
	}

	return commandResult.Stdout, nil
}

// ListInstalledKBs returns the numbers of all installed KB articles, e.g. "5079473"
func (w WindowsUpdatesManager) ListInstalledKBs(ctx context.Context) ([]string, error) {
	output, err := w.getInstalledKBs(ctx)
	if err != nil {
		return nil, err
	}

	kbs, err := parseInstalledKBsOutput(output)
	if err != nil {
		return nil, fmt.Errorf("error parsing installed KBs output: %s", err)
	}

	return kbs, nil
}

func (w WindowsUpdatesManager) getInstalledKBs(ctx context.Context) (string, error) {
	// Both sources describe the current state of the system. The update history is not used on purpose,
	// because it still lists updates that were removed later on.
	// Win32_QuickFixEngineering is safe. DO NOT use Win32_Product, see getInstalledApps.
	// Online = $false: The installed updates are known locally, no need to contact a server.
	psScript := `
		[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
		$KBs = @()

		try {
		    Get-CimInstance -ClassName Win32_QuickFixEngineering -ErrorAction Stop | ForEach-Object {
		        $KBs += $_.HotFixID
		    }
		} catch {}

		try {
		    $Session = New-Object -ComObject Microsoft.Update.Session
		    $Searcher = $Session.CreateUpdateSearcher()
		    $Searcher.Online = $false
		    $Results = $Searcher.Search("IsInstalled=1")
		    foreach ($Update in $Results.Updates) {
		        foreach ($KB in $Update.KBArticleIDs) {
		            $KBs += $KB
		        }
		    }
		} catch {}

		ConvertTo-Json -Compress -InputObject @($KBs)
	`

	timeout := time.Duration(120 * time.Second)
	commandResult, err := utils.RunCommand(ctx, utils.CommandArgs{
		Timeout: timeout,
		Command: psScript,
		Shell:   "powershell_command",
	})

	if err != nil {
		return "", fmt.Errorf("error fetching installed KBs via PowerShell: %s", commandResult.Stdout)
	}

	if commandResult.RC != 0 {
		return "", fmt.Errorf("error fetching installed KBs via PowerShell, exit code %d: %s", commandResult.RC, commandResult.Stdout)
	}

	return commandResult.Stdout, nil
}

func (w WindowsUpdatesManager) parsePowerShellUpdateSessionOutput(output string) ([]WindowsUpdate, error) {
	var updates []WindowsUpdate

	if output == "" {
		return updates, nil
	}

	// Handle both single object and array of objects
	var dst []PowerShellUpdateJson
	if strings.HasPrefix(output, "[") {
		err := json.Unmarshal([]byte(output), &dst)
		if err != nil {
			return nil, err
		}
	} else {
		var singleRecord PowerShellUpdateJson
		err := json.Unmarshal([]byte(output), &singleRecord)
		if err != nil {
			return nil, err
		}
		dst = []PowerShellUpdateJson{singleRecord}
	}

	reported := map[string]bool{}
	for _, update := range dst {
		// An update deployed by SCCM can also be found by the Windows Update API
		if isReportedUpdate(reported, update.UpdateID) {
			continue
		}

		updates = append(updates, WindowsUpdate{
			Title:            update.Title,
			Description:      update.Description,
			KBArticleIDs:     update.KBArticleIDs,
			IsInstalled:      false,
			IsSecurityUpdate: isSecurityUpdate(update),
			IsOptional:       update.BrowseOnly || containsIgnoreCase(update.Categories, "Optional"),
			UpdateID:         update.UpdateID,
			RevisionNumber:   update.RevisionNumber,
			RebootRequired:   update.RebootRequired,
		})
	}
	return updates, nil
}

// isReportedUpdate returns true if the update was already reported and remembers the update otherwise.
// Updates without ID are always reported.
func isReportedUpdate(reported map[string]bool, updateID string) bool {
	if updateID == "" {
		return false
	}
	id := strings.ToLower(updateID)
	if reported[id] {
		return true
	}
	reported[id] = true
	return false
}

// isSecurityUpdate checks the IDs of the classification and the products, which are the same in all languages.
// The names are checked as well for updates without IDs.
func isSecurityUpdate(update PowerShellUpdateJson) bool {
	for _, id := range update.CategoryIDs {
		if containsEqualFold(securityCategoryIDs, id) {
			return true
		}
	}
	if containsIgnoreCase(update.Categories, "Security") {
		return true
	}
	if containsIgnoreCase(update.Categories, "Defender") {
		return true
	}
	return containsIgnoreCase(update.Categories, "Critical")
}

func containsEqualFold(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func (w WindowsUpdatesManager) RebootRequired(ctx context.Context) (bool, error) {
	if anyRegistryKeyExists(rebootRequiredKeys, registryKeyExists) {
		return true, nil
	}
	if !registryKeyExists(sccmClientServiceKey) {
		return false, nil
	}
	return w.sccmRebootPending(ctx), nil
}

// sccmRebootPending asks the Configuration Manager (SCCM) client, if a reboot is pending to finish the installation
// of deployed updates. Errors are only logged, so the inventory is still reported.
func (w WindowsUpdatesManager) sccmRebootPending(ctx context.Context) bool {
	psScript := `
		$Result = Invoke-CimMethod -Namespace 'root\ccm\ClientSDK' -ClassName CCM_ClientUtilities -MethodName DetermineIfRebootPending -ErrorAction Stop
		ConvertTo-Json -Compress -InputObject ([bool]($Result.RebootPending -or $Result.IsHardRebootPending))
	`
	commandResult, err := utils.RunCommand(ctx, utils.CommandArgs{
		Timeout: 60 * time.Second,
		Command: psScript,
		Shell:   "powershell_command",
	})
	if err != nil || commandResult.RC != 0 {
		log.Errorln("Packagemanager: Error asking the SCCM client for a pending reboot:", err, commandResult.Stdout)
		return false
	}

	rebootPending, err := parseSccmRebootPendingOutput(commandResult.Stdout)
	if err != nil {
		log.Errorln("Packagemanager: Error parsing the pending reboot of the SCCM client:", err)
		return false
	}
	return rebootPending
}

func parseSccmRebootPendingOutput(output string) (bool, error) {
	var rebootPending bool
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &rebootPending); err != nil {
		return false, err
	}
	return rebootPending, nil
}

// anyRegistryKeyExists returns true if at least one of the keys exists. keyExists is replaced by the tests.
func anyRegistryKeyExists(paths []string, keyExists func(path string) bool) bool {
	for _, path := range paths {
		if keyExists(path) {
			return true
		}
	}
	return false
}

// registryKeyExists returns true if the key exists below HKEY_LOCAL_MACHINE
func registryKeyExists(path string) bool {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	key.Close()
	return true
}

func containsIgnoreCase(slice []string, item string) bool {
	for _, s := range slice {
		needle := strings.ToLower(item)
		haystack := strings.ToLower(s)
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func (w WindowsUpdatesManager) CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return PackageInfo{
			Stats: PackageStats{
				LastError: err,
			},
		}, err
	}

	// Platform: "Microsoft Windows 11 Enterprise"
	// PlatformFamily: "Standalone Workstation"
	// PlatformFamily: "Server (Domain Controller)"
	// PlatformFamily: "Server"
	// PlatformVersion: "24H2"
	// KernelVersion: "10.0.26100.7462 Build 26100.7462"
	result := PackageInfo{
		Enabled:    true,
		Pending:    false,
		LastUpdate: time.Now().Unix(),
		Stats: PackageStats{
			PackageManager:  "windows-updates",
			OperatingSystem: "windows",
			OsName:          info.Platform,
			OsVersion:       fmt.Sprintf("%s (%s)", info.PlatformVersion, info.KernelVersion),
			OsFamily:        "windows", // keep it simple
			Uptime:          int64(info.Uptime),
			AgentVersion:    config.AgentVersion,
		},
	}

	installedApps, err := w.ListInstalledApps(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.InstalledPackages = int64(len(installedApps))

	// The installed KBs are not collected for now. The server does not use them yet, and the
	// vulnerability scan delivers better results with the build number from OsVersion alone.
	// To enable it, remove the comments here and at the assignment to result.WindowsInstalledKBs
	// below.
	//
	// A failure here must not cost us the rest of the inventory.
	// installedKBs, err := w.ListInstalledKBs(ctx)
	// if err != nil {
	// 	log.Errorln("Packagemanager: Error collecting installed KBs:", err)
	// }

	var availableUpdates []WindowsUpdate
	if enableUpdateCheck {
		availableUpdates, err = w.ListAvailableUpdates(ctx)
		if err != nil {
			result.Stats.LastError = err
			return result, err
		}
		result.Stats.UpgradablePackages = int64(len(availableUpdates))

		// Count security updates
		var securityUpdates int64
		for _, update := range availableUpdates {
			if update.IsSecurityUpdate {
				securityUpdates++
			}
		}
		result.Stats.SecurityUpdates = securityUpdates
	}

	rebootRequired, err := w.RebootRequired(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.RebootRequired = rebootRequired

	// Truncate descriptions if needed
	for i := range availableUpdates {
		availableUpdates[i].Description = truncateDescription(availableUpdates[i].Description, limitDescriptionLength)
	}

	result.WindowsApps = installedApps
	result.WindowsUpdates = availableUpdates
	// result.WindowsInstalledKBs = installedKBs

	return result, nil
}
