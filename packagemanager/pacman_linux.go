package packagemanager

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

// PacmanManager implements PackageManager for pacman
type PacmanManager struct{}

// IsAvailable checks if pacman is available on the system
func (p PacmanManager) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

// UpdateMetadata updates the package metadata using pacman -Sy
func (p PacmanManager) UpdateMetadata(ctx context.Context) error {
	// While we could implement this, updating package metadata
	// on Arch Linux systems is is not recommended to update only the package database.
	// Therefore, this function is currently is disabled.
	// Please open an Issue if you think this is a bug. I'm not a Arch Linux user.

	//timeout := 300 * time.Second
	//_, err := utils.RunCommand(ctx, utils.CommandArgs{
	//	Command: "pacman -Sy",
	//	Timeout: timeout,
	//	Env: map[string]string{
	//		"LANG": "C",
	//	},
	//})
	//return err
	return nil
}

// ListInstalledPackages lists installed packages using pacman
func (p PacmanManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	output, err := p.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		return nil, err
	}

	return p.parsePacmanQiOutput(output)
}

func (p PacmanManager) getInstalledPackagesWithCancel(ctx context.Context) (string, error) {
	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "pacman -Qi",
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

func (p PacmanManager) parsePacmanQiOutput(output string) ([]Package, error) {
	var packages []Package
	var currentPackage Package

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			// End of current package entry
			if currentPackage.Name != "" {
				packages = append(packages, currentPackage)
				currentPackage = Package{}
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Name":
			currentPackage.Name = value
		case "Version":
			currentPackage.Version = value
		case "Description":
			currentPackage.Description = value
		}
	}

	// Append the last package if exists
	if currentPackage.Name != "" {
		packages = append(packages, currentPackage)
	}

	return packages, nil
}

// ListUpgradablePackages lists upgradable packages using checkupdates
func (p PacmanManager) ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error) {
	output, err := p.getUpgradablePackages(ctx)
	if err != nil {
		return nil, err
	}

	return p.parsePacmanUpgradeOutput(output)
}

func (p PacmanManager) getUpgradablePackages(ctx context.Context) (string, error) {
	// This command is taken from good old check_apt
	// https://github.com/monitoring-plugins/monitoring-plugins/blob/bfc6492562f6cef4badda192142a0d10a3ed870b/plugins/check_apt.c#L45-L47

	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "checkupdates",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if result.RC == 127 {
		return "", fmt.Errorf("checkupdates command not found; please install the 'pacman-contrib' package")
	}

	if err != nil {
		return "", err
	}

	if result.RC == 1 {
		// Unknown error
		return "", fmt.Errorf("error while executing checkupdates command: %s", result.Stdout)
	}

	if result.RC == 2 {
		// No updates available
		return "", nil
	}

	return result.Stdout, nil
}

func (p PacmanManager) parsePacmanUpgradeOutput(output string) ([]PackageUpdate, error) {
	var pkgs []PackageUpdate

	if output == "" {
		return pkgs, nil
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Example line:
		// package_name current_version -> available_version
		// libseccomp 2.5.6-1 -> 2.6.0-1
		// libtasn1 4.20.0-1 -> 4.21.0-1

		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue
		}

		pkgName := parts[0]
		currentVersion := parts[1]

		availableParts := strings.SplitN(parts[2], "->", 2)
		if len(availableParts) < 2 {
			continue
		}
		availableVersion := strings.TrimSpace(availableParts[1])

		pkgs = append(pkgs, PackageUpdate{
			Name:             pkgName,
			CurrentVersion:   currentVersion,
			AvailableVersion: availableVersion,
			IsSecurityUpdate: false, // Pacman does not provide security update info
		})
	}
	return pkgs, nil
}

// RebootRequired checks if a reboot is required on the system
func (p PacmanManager) RebootRequired(ctx context.Context) (bool, error) {
	return false, nil
}
