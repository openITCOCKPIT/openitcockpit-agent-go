package packagemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	"github.com/shirou/gopsutil/v4/host"
)

// MacOSUpdatesManager implements MacOSManager for macOS Updates
type MacOSUpdatesManager struct{}

func (m MacOSUpdatesManager) ListInstalledApps(ctx context.Context) ([]Package, error) {
	output, err := m.getInstalledApps(ctx)
	if err != nil {
		return nil, err
	}
	return parseMacOSInstalledAppsOutput(output)
}

func (m MacOSUpdatesManager) getInstalledApps(ctx context.Context) (string, error) {
	timeout := time.Duration(120 * time.Second)
	// The app name is taken from the apps Info.plist CFBundleName key
	// This is localized and there is no easy way to get the non localized name
	// https://developer.apple.com/library/archive/documentation/General/Reference/InfoPlistKeyReference/Articles/AboutInformationPropertyListFiles.html#//apple_ref/doc/uid/TP40009254-SW1
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "system_profiler SPApplicationsDataType -json",
		Timeout: timeout,
	})

	// There is also a "pkgutil --pkgs" command avaiable but it does not provide versions
	// so I decided to ignore it for now.

	if err != nil {
		return "", err
	}

	if result.RC != 0 {
		return "", fmt.Errorf("error fetching installed apps, exit code %d: %s", result.RC, result.Stdout)
	}

	return result.Stdout, nil
}

func (m MacOSUpdatesManager) getAvailableUpdates(ctx context.Context) (string, error) {
	timeout := time.Duration(120 * time.Second)
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "softwareupdate --list --no-scan",
		Timeout: timeout,
	})

	if err != nil {
		return "", err
	}

	if result.RC != 0 {
		return "", fmt.Errorf("error fetching macOS Updates via softwareupdate, exit code %d: %s", result.RC, result.Stdout)
	}

	return result.Stdout, nil
}

func (m MacOSUpdatesManager) ListAvailableUpdates(ctx context.Context) ([]MacosUpdate, error) {
	output, err := m.getAvailableUpdates(ctx)
	if err != nil {
		return nil, err
	}

	return parseMacOSSoftwareUpdateOutput(output)
}

func (m MacOSUpdatesManager) CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return PackageInfo{
			Stats: PackageStats{
				LastError: err,
			},
		}, err
	}

	// Platform: "darwin"
	// PlatformFamily: "Standalone Workstation"
	// PlatformFamily: "Server"
	// PlatformVersion: "26.2"
	// KernelVersion: "25.2.0"
	result := PackageInfo{
		Enabled:    true,
		Pending:    false,
		LastUpdate: time.Now().Unix(),
		Stats: PackageStats{
			PackageManager:  "macos-updates",
			OperatingSystem: "macos",
			OsName:          "macos",
			OsVersion:       info.PlatformVersion,
			OsFamily:        info.Platform, // darwin
			Uptime:          int64(info.Uptime),
			AgentVersion:    config.AgentVersion,
		},
	}

	installedApps, err := m.ListInstalledApps(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.InstalledPackages = int64(len(installedApps))

	var availableUpdates []MacosUpdate
	if enableUpdateCheck {
		availableUpdates, err = m.ListAvailableUpdates(ctx)
		if err != nil {
			result.Stats.LastError = err
			return result, err
		}
		result.Stats.UpgradablePackages = int64(len(availableUpdates))
	}

	// Truncate descriptions if needed
	for i := range availableUpdates {
		availableUpdates[i].Description = truncateDescription(availableUpdates[i].Description, limitDescriptionLength)
	}

	result.MacOSApps = installedApps
	result.MacosUpdates = availableUpdates

	return result, nil
}
