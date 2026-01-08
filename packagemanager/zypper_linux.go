package packagemanager

import (
	"context"
	"encoding/xml"
	"os/exec"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

// ZypperManager implements PackageManager for zypper
type ZypperManager struct{}

// ZYPPER_EXIT_INF_REBOOT_NEEDED indicates that a reboot is needed
const ZYPPER_EXIT_INF_REBOOT_NEEDED = 102

// IsAvailable checks if zypper is available on the system
func (z ZypperManager) IsAvailable() bool {
	_, err := exec.LookPath("zypper")
	return err == nil
}

// UpdateMetadata updates the package metadata using zypper refresh
func (z ZypperManager) UpdateMetadata(ctx context.Context) error {
	timeout := 300 * time.Second
	_, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "zypper refresh",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})
	return err
}

// ListInstalledPackages lists installed packages using dpkg
func (z ZypperManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	// For simplicity, we reuse the rpm implementation here
	rpm := RpmManager{}
	return rpm.ListInstalledPackages(ctx)
}

// ListUpgradablePackages lists upgradable packages using apt
func (z ZypperManager) ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error) {
	output, err := z.getUpgradablePackages(ctx)
	if err != nil {
		return nil, err
	}

	return z.parseZypperUpgradeOutput(output)
}

func (z ZypperManager) getUpgradablePackages(ctx context.Context) (string, error) {
	timeout := 120 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "zypper --no-refresh --quiet --xmlout list-updates",
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

func (z ZypperManager) parseZypperUpgradeOutput(output string) ([]PackageUpdate, error) {
	var pkgs []PackageUpdate

	type Update struct {
		Name        string `xml:"name,attr"`
		EditionOld  string `xml:"edition-old,attr"`
		Edition     string `xml:"edition,attr"`
		Description string `xml:"description"`
	}

	type UpdateList struct {
		Updates []Update `xml:"update"`
	}

	type UpdateStatus struct {
		UpdateList UpdateList `xml:"update-list"`
	}

	type Stream struct {
		UpdateStatus UpdateStatus `xml:"update-status"`
	}

	var stream Stream

	err := xml.Unmarshal([]byte(output), &stream)
	if err != nil {
		return nil, err
	}

	for _, update := range stream.UpdateStatus.UpdateList.Updates {
		pkg := PackageUpdate{
			Name:             update.Name,
			CurrentVersion:   update.EditionOld,
			AvailableVersion: update.Edition,

			// Zypper / openSUSE does not have "security updates" for packages in the same way as other distros
			// openSUSE has "patches" but they are not listed in the zypper list-updates output
			// Patches can not (or at least I'm not aware of how to) be mapped to individual packages
			IsSecurityUpdate: false,
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func (z ZypperManager) getSecurityPatches(ctx context.Context) (string, error) {
	timeout := 120 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "zypper --no-refresh --quiet --xmlout list-patches --category security",
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

func (z ZypperManager) parseZypperSecurityPatchesOutput(output string) ([]PackageUpdate, error) {
	var pkgs []PackageUpdate

	type SecurityPatch struct {
		Name        string `xml:"name,attr"`
		Summary     string `xml:"summary"`
		Edition     string `xml:"edition,attr"`
		Description string `xml:"description"`
	}

	type UpdateList struct {
		Updates []SecurityPatch `xml:"update"`
	}

	type UpdateStatus struct {
		UpdateList UpdateList `xml:"update-list"`
	}

	type Stream struct {
		UpdateStatus UpdateStatus `xml:"update-status"`
	}

	var stream Stream

	err := xml.Unmarshal([]byte(output), &stream)
	if err != nil {
		return nil, err
	}

	for _, update := range stream.UpdateStatus.UpdateList.Updates {
		pkg := PackageUpdate{
			Name:             update.Name,
			CurrentVersion:   "",
			AvailableVersion: update.Edition,
			IsSecurityUpdate: true,
			IsPatch:          true,
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func (z ZypperManager) RebootRequired(ctx context.Context) (bool, error) {
	// Checks if a reboot is required based on zypper output
	// From the zypper man page:
	//needs-rebooting
	//   Checks if the reboot-needed flag was set by a previous update or install of a core library or service. +
	//   The reboot-needed flag is set when a package from a predefined list (/etc/zypp/needreboot) is updated or
	//   installed. Exit code ZYPPER_EXIT_INF_REBOOT_NEEDED(102) indicates that a reboot is needed, otherwise the exit
	//   code is set to ZYPPER_EXIT_OK(0).

	timeout := 60 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "zypper --no-refresh needs-rebooting",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if err != nil {
		return false, err
	}

	return result.RC == ZYPPER_EXIT_INF_REBOOT_NEEDED, nil
}
