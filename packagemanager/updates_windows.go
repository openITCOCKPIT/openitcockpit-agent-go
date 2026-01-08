package packagemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
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
}

func (w WindowsUpdatesManager) IsAvailable() bool {
	return true
}
func (w WindowsUpdatesManager) UpdateMetadata(ctx context.Context) error {
	return nil
}

func (w WindowsUpdatesManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	return nil, nil
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
	// DO NOT use WMI and the Win32_Product class!
	// It will may repair all installed MSI packages on the system!
	// DO NOT EVEN THINK ABOUT IT!
	// https://support.microsoft.com/kb/974524
	// https://sdmsoftware.com/wmi/why-win32_product-is-bad-news/
	//
	// The file 'updates_windows_ole.txt' contains a implemtation using OLE/COM
	// how ever, this was way to complicated for what we try to achieve here.

	// Alternaitve PowerShell script with JSON output
	psScript := `
		[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
		$Session = New-Object -ComObject Microsoft.Update.Session
		$Searcher = $Session.CreateUpdateSearcher()
		$Results = $Searcher.Search("IsInstalled=0 and IsHidden=0")
		$Updates = @()

		foreach ($Update in $Results.Updates) {
		    $categories = @()
		    foreach ($cat in $Update.Categories) {
		        $categories += $cat.Name
		    }
		    $Updates += [PSCustomObject]@{
		        Title         = $Update.Title
		        Description   = $Update.Description
		        KBArticleIDs  = $Update.KBArticleIDs
		        Categories    = $categories
		        SupportUrl    = $Update.SupportUrl
		        RebootRequired= $Update.RebootRequired
		        EulaAccepted  = $Update.EulaAccepted
				UpdateID       = $Update.Identity.UpdateID
        		RevisionNumber = $Update.Identity.RevisionNumber
		    }
		}

		$Updates | ConvertTo-Json -Compress -Depth 5
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

	for _, update := range dst {

		var isSecurity, isDefender, isCritical bool
		isSecurity = containsIgnoreCase(update.Categories, "Security")
		isDefender = containsIgnoreCase(update.Categories, "Defender")
		isCritical = containsIgnoreCase(update.Categories, "Critical")

		updates = append(updates, WindowsUpdate{
			Title:            update.Title,
			Description:      update.Description,
			KBArticleIDs:     update.KBArticleIDs,
			IsInstalled:      false,
			IsSecurityUpdate: isSecurity || isDefender || isCritical,
			IsOptional:       containsIgnoreCase(update.Categories, "Optional"),
			UpdateID:         update.UpdateID,
			RevisionNumber:   update.RevisionNumber,
			RebootRequired:   update.RebootRequired,
		})
	}
	return updates, nil
}

func (w WindowsUpdatesManager) RebootRequired(ctx context.Context) (bool, error) {
	return false, nil
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
