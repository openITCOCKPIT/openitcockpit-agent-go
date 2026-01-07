package packagemanager

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

// AptManager implements PackageManager for apt
type AptManager struct{}

// UpdateMetadata updates the package metadata using apt-get update
func (a AptManager) UpdateMetadata() error {
	timeout := 300 * time.Second
	_, err := utils.RunCommand(context.Background(), utils.CommandArgs{
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
		parts := strings.SplitN(line, " ", 3)

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

// YumManager implements PackageManager for yum
type YumManager struct{}

func (y YumManager) ListInstalledPackages() ([]Package, error) {
	cmd := exec.Command("yum", "list", "installed")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var pkgs []Package
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pkgs = append(pkgs, Package{Name: fields[0], Version: fields[1]})
		}
	}
	return pkgs, nil
}

func (y YumManager) ListUpgradablePackages() ([]PackageUpdate, error) {
	cmd := exec.Command("yum", "check-update")
	out, err := cmd.Output()
	if err != nil {
		// yum returns exit code 100 if updates are available, so ignore error if output exists
		if len(out) == 0 {
			return nil, err
		}
	}
	lines := strings.Split(string(out), "\n")
	var updates []PackageUpdate
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			updates = append(updates, PackageUpdate{
				Name:             fields[0],
				CurrentVersion:   "", // yum does not show current version here
				AvailableVersion: fields[1],
			})
		}
	}
	return updates, nil
}
