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

// dpkgQueryFormat is the format of dpkg-query. The description has to be the last field, as it has multiple lines.
// ${source:Package} and ${source:Version} are virtual fields, unlike ${Source} they are also set if the source package has
// the name and the version of the binary package, and they never contain a version in parentheses.
const dpkgQueryFormat = "${Status}\t${Package}\t${Version}\t${Architecture}\t${source:Package}\t${source:Version}\t${Description}\nEND\n"

// dpkgQueryFieldCount is the number of tab separated fields of dpkgQueryFormat
const dpkgQueryFieldCount = 7

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
		// dpkg-query -W also lists packages which got removed, but not purged (state "rc", config-files)
		Command: "dpkg-query -W -f=\"" + dpkgQueryFormat + "\"",
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
	records := strings.Split(output, "\nEND\n")
	var pkgs []Package
	for _, record := range records {
		// status<TAB>package<TAB>version<TAB>architecture<TAB>source package<TAB>source version<TAB>description
		// install ok installed	libcurl4t64	8.5.0-2ubuntu10.16	arm64	curl	8.5.0-2ubuntu10.16	easy-to-use client-side <until END>
		parts := strings.SplitN(record, "\t", dpkgQueryFieldCount)

		if len(parts) < dpkgQueryFieldCount {
			continue
		}

		// want flag status
		status := strings.Fields(parts[0])
		if len(status) != 3 {
			continue
		}

		switch status[2] {
		case "not-installed", "config-files":
			// Removed packages, "rc" only has its configuration files left
			continue
		}

		pkgs = append(pkgs, Package{
			Name:          parts[1],
			Version:       parts[2],
			Architecture:  parts[3],
			SourcePackage: parts[4],
			SourceVersion: parts[5],
			Description:   parts[6],
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
			OsFamily:        info.PlatformFamily,  // "debian"
			Uptime:          int64(info.Uptime),
			AgentVersion:    config.AgentVersion,
			KernelRelease:   info.KernelVersion,
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
			// The installed packages are reported, even if the updates can not be determined
			// (e.g. a repository is not reachable or a required tool like checkupdates is not installed)
			result.Stats.LastError = err
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
