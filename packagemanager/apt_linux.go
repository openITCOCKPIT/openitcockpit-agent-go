package packagemanager

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	"github.com/shirou/gopsutil/v4/host"
)

// AptManager implements PackageManager for apt
type AptManager struct{}

// IsAvailable checks if apt-get is available on the system
func (a AptManager) IsAvailable() bool {
	_, err := exec.LookPath("apt-get")
	return err == nil
}

// UpdateMetadata updates the package metadata using apt-get update
func (a AptManager) UpdateMetadata(ctx context.Context) error {
	timeout := 300 * time.Second
	_, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "apt-get update -q",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})
	return err
}

// ListInstalledPackages lists installed packages using dpkg
func (a AptManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	output, err := a.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		return nil, err
	}

	return a.parseDpkgOutput(output)
}

func (a AptManager) getInstalledPackagesWithCancel(ctx context.Context) (string, error) {
	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "dpkg-query -W -f=\"${Package} ${Version} ${Description}\nEND\n\"",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if err != nil {
		return "", err
	}

	return result.Stdout, nil
}

func (a AptManager) parseDpkgOutput(output string) ([]Package, error) {
	lines := strings.Split(output, "\nEND\n")
	var pkgs []Package
	for _, line := range lines {
		// package version description
		// zsh 5.9-6ubuntu2 shell with lots of features <until END>
		parts := strings.SplitN(line, " ", 3)

		if len(parts) < 3 {
			continue
		}

		pkgs = append(pkgs, Package{
			Name:        parts[0],
			Version:     parts[1],
			Description: parts[2],
		})
	}
	return pkgs, nil
}

// ListUpgradablePackages lists upgradable packages using apt
func (a AptManager) ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error) {
	output, err := a.getUpgradablePackages(ctx)
	if err != nil {
		return nil, err
	}

	return a.parseAptUpgradeOutput(output)
}

func (a AptManager) getUpgradablePackages(ctx context.Context) (string, error) {
	// This command is taken from good old check_apt
	// https://github.com/monitoring-plugins/monitoring-plugins/blob/bfc6492562f6cef4badda192142a0d10a3ed870b/plugins/check_apt.c#L45-L47

	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "apt-get -s -o Debug::NoLocking=1 upgrade",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if err != nil {
		return "", err
	}

	return result.Stdout, nil
}

func (a AptManager) parseAptUpgradeOutput(output string) ([]PackageUpdate, error) {
	lines := strings.Split(output, "\n")
	var pkgs []PackageUpdate
	for _, line := range lines {
		if strings.HasPrefix(line, "Inst ") {
			// Example line:
			// Inst package [current_version] (available_version repository)

			parts := strings.SplitN(line, " ", 4)

			if len(parts) < 4 {
				continue
			}

			pkgName := parts[1]
			currentVersion := strings.Trim(parts[2], "[]")
			availablePart := parts[3]
			// "(2.80.0-6ubuntu3.6 Ubuntu:24.04/noble-updates, Ubuntu:24.04/noble-security [all])"
			availableVersion := strings.SplitN(availablePart, " ", 2)[0]
			// Now we have "(2.80.0-6ubuntu3.6" - trim the leading "("
			availableVersion = strings.TrimPrefix(availableVersion, "(")

			// Check if it's a security update
			isSecurityUpdate := strings.Contains(availablePart, "Debian-Security:")
			if !isSecurityUpdate {
				// For Ubuntu, security updates are often indicated by "-security" in the version string
				// The complete string is "Ubuntu:24.04/noble-security" and we would need a regex to parse this
				// but i guess this will do for now
				isSecurityUpdate = strings.Contains(availablePart, "-security")
			}

			pkgs = append(pkgs, PackageUpdate{
				Name:             pkgName,
				CurrentVersion:   currentVersion,
				AvailableVersion: availableVersion,
				IsSecurityUpdate: isSecurityUpdate,
			})
		}
	}
	return pkgs, nil
}

// RebootRequired checks if a reboot is required on the system
func (a AptManager) RebootRequired(ctx context.Context) (bool, error) {
	// Official way as Debian documents it here:
	// 9.12. Signaling that a reboot is required
	// https://www.debian.org/doc/debian-policy/ch-opersys.html#signaling-that-a-reboot-is-required
	paths := []string{"/run/reboot-required", "/var/run/reboot-required", "/run/reboot-required.pkgs", "/var/run/reboot-required.pkgs"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true, nil
		}
	}
	return false, nil
}

func (a AptManager) CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error) {

	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return PackageInfo{
			Stats: PackageStats{
				LastError: err,
			},
		}, err
	}

	// Platform: "ubuntu"
	// PlatformFamily: "debian"
	// PlatformVersion: "24.04"
	// KernelVersion: "6.6.87.2-microsoft-standard-WSL2"
	result := PackageInfo{
		Enabled:    true,
		Pending:    false,
		LastUpdate: time.Now().Unix(),
		Stats: PackageStats{
			PackageManager:  "apt",
			OperatingSystem: "linux",
			OsName:          info.Platform,        // "ubuntu"
			OsVersion:       info.PlatformVersion, // "24.04"
			Uptime:          int64(info.Uptime),
			AgentVersion:    config.AgentVersion,
		},
	}

	if enableUpdateCheck {
		err := a.UpdateMetadata(ctx)
		if err != nil {
			result.Stats.LastError = err
		}
	}

	installedPackages, err := a.ListInstalledPackages(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.InstalledPackages = int64(len(installedPackages))

	var upgradablePackages []PackageUpdate
	if enableUpdateCheck {
		upgradablePackages, err = a.ListUpgradablePackages(ctx)
		if err != nil {
			result.Stats.LastError = err
			return result, err
		}
		result.Stats.UpgradablePackages = int64(len(upgradablePackages))

		// Count security updates
		var securityUpdates int64
		for _, pkg := range upgradablePackages {
			if pkg.IsSecurityUpdate {
				securityUpdates++
			}
		}
		result.Stats.SecurityUpdates = securityUpdates
	}

	rebootRequired, err := a.RebootRequired(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.RebootRequired = rebootRequired

	// Truncate descriptions if needed
	for i := range installedPackages {
		installedPackages[i].Description = truncateDescription(installedPackages[i].Description, limitDescriptionLength)
	}

	result.LinuxPackages = installedPackages
	result.LinuxUpdates = upgradablePackages

	return result, nil
}
