package packagemanager

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

// DnfManager implements PackageManager for dnf
type DnfManager struct{}

// IsAvailable checks if dnf is available on the system
func (d DnfManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

// UpdateMetadata updates the package metadata using dnf
func (d DnfManager) UpdateMetadata(ctx context.Context) error {
	// DNF automatically updates metadata during operations, so no action is needed here.
	return nil
}

// ListInstalledPackages lists installed packages using dnf
func (d DnfManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	output, err := d.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		return nil, err
	}

	return d.parseDnfListInstalledOutput(output)
}

func (d DnfManager) getInstalledPackagesWithCancel(ctx context.Context) (string, error) {
	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "dnf list installed",
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

// parseDnfListInstalledOutput parses the output of 'dnf list installed' into a slice of Package structs
func (d DnfManager) parseDnfListInstalledOutput(output string) ([]Package, error) {
	lines := strings.Split(output, "\n")
	var pkgs []Package
	foundHeader := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip header lines until we find the actual package list
		if !foundHeader {
			if strings.HasPrefix(line, "Installed Packages") {
				foundHeader = true
			}
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// Example line: bash.x86_64           5.0.17-2.fc32           @fedora
		name := d.removeArchFromPkgName(fields[0])

		version := fields[1]
		// To get the description for a specific package with DNF or RPM, you must query each package individually:
		// dnf info <package> or rpm -qi <package>
		// Sounds like a horrible idea to do this for all installed packages.
		pkgs = append(pkgs, Package{
			Name:        name,
			Version:     version,
			Description: "",
		})
	}
	return pkgs, nil
}

// ListUpgradablePackages lists upgradable packages using dnf
func (d DnfManager) ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error) {
	// Get list of upgradable packages
	checkUpdateOutput, updatesAvailable, err := d.getUpgradablePackages(ctx)
	if err != nil {
		return nil, err
	}

	if !updatesAvailable {
		return []PackageUpdate{}, nil
	}

	// Get list of all installed packages to determine current versions
	installedPackages, err := d.ListInstalledPackages(ctx)
	if err != nil {
		return nil, err
	}

	// Get list of security relevant updates
	securityOutput, err := d.getSecurityRelevantUpdates(ctx)
	if err != nil {
		return nil, err
	}

	securityPackages, err := d.parseDnfSecurityUpdateOutput(securityOutput)
	if err != nil {
		return nil, err
	}

	updates, err := d.parseDnfCheckUpdateOutput(checkUpdateOutput, installedPackages, securityPackages)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

func (d DnfManager) getUpgradablePackages(ctx context.Context) (string, bool, error) {
	// check-update can run long, in case it refreshes metadata, so we set a higher timeout here
	// Also, a return code of 100 means there are updates available, so we should not treat it as an error
	timeout := 300 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "dnf --quiet check-update",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if err != nil {
		return "", false, err
	}

	updatesAvailable := false
	if result.RC == 100 {
		updatesAvailable = true
	}

	return result.Stdout, updatesAvailable, nil
}

func (d DnfManager) parseDnfCheckUpdateOutput(packagesToUpdateOutput string, installedPackages []Package, securityPackages map[string]bool) ([]PackageUpdate, error) {
	// Create a map of installed packages for quick lookup
	installedPkgMap := make(map[string]string)
	for _, pkg := range installedPackages {
		installedPkgMap[pkg.Name] = pkg.Version
	}

	lines := strings.Split(packagesToUpdateOutput, "\n")
	var pkgs []PackageUpdate
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Last metadata expiration check") || strings.HasPrefix(line, "Failed to set locale") {
			continue
		}

		if line == "Obsoleting Packages" {
			// Stop processing at the Obsoleting Packages section
			// These packages are replaced by others and should not be listed as upgradable
			break
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// Example line: bash.x86_64           5.0.17-2.fc32           @fedora
		name := d.removeArchFromPkgName(fields[0])

		// Check if the package is installed
		currentVersion, installed := installedPkgMap[name]
		if !installed {
			currentVersion = ""
		}

		// Determine if this is a security update
		_, isSecurityUpdate := securityPackages[name]

		availableVersion := fields[1]
		pkgs = append(pkgs, PackageUpdate{
			Name:             name,
			CurrentVersion:   currentVersion,
			AvailableVersion: availableVersion,
			IsSecurityUpdate: isSecurityUpdate,
		})
	}
	return pkgs, nil
}

func (d DnfManager) getSecurityRelevantUpdates(ctx context.Context) (string, error) {
	// alternatively, we could use
	// dnf updateinfo info --available
	timeout := 300 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "dnf updateinfo list --available security",
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

func (d DnfManager) parseDnfSecurityUpdateOutput(output string) (map[string]bool, error) {
	lines := strings.Split(output, "\n")
	securityUpdates := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Last metadata expiration check") || strings.HasPrefix(line, "Failed to set locale") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Example line: ALSA-2025:20532 Moderate/Sec.  grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64
		// Available Security Advisory Types:
		// ALSA = AlmaLinux Security Advisory
		// CESA = CentOS Errata and Security Advisory
		// RHSA = Red Hat Security Advisory
		// RLSA = Rocky Linux Security Advisory
		pkgWithVersionAndArch := fields[2]

		// Some packages may appear multiple times with different advisories.
		// We use a map to avoid duplicates.
		// ALSA-2025:19409 Moderate/Sec.  kernel-5.14.0-570.60.1.el9_6.x86_64
		// ALSA-2025:19930 Moderate/Sec.  kernel-5.14.0-570.62.1.el9_6.x86_64
		// ALSA-2025:22405 Moderate/Sec.  kernel-5.14.0-611.11.1.el9_7.x86_64
		// ALSA-2025:19409 Moderate/Sec.  kernel-core-5.14.0-570.60.1.el9_6.x86_64
		// ALSA-2025:19930 Moderate/Sec.  kernel-core-5.14.0-570.62.1.el9_6.x86_64
		// ALSA-2025:19409 Moderate/Sec.  kernel-headers-5.14.0-570.60.1.el9_6.x86_64
		// ALSA-2025:19930 Moderate/Sec.  kernel-headers-5.14.0-570.62.1.el9_6.x86_64

		// Remove the version from the package name
		// Find the last hyphen before a digit (start of version)
		pkgWithArch := d.removeVersionFromPkgName(pkgWithVersionAndArch)
		pkgName := d.removeArchFromPkgName(pkgWithArch)
		securityUpdates[pkgName] = true
	}
	return securityUpdates, nil
}

func (d DnfManager) removeArchFromPkgName(pkgName string) string {
	// Known architectures suffixes: .x86_64 .noarch .i386 .i686 .arm64 .ppc64le .aarch64
	// probably others as well - just remove everything after the last dot
	if idx := strings.LastIndex(pkgName, "."); idx != -1 {
		return pkgName[:idx]
	}
	return pkgName
}

// removeVersionFromPkgName removes the version from a package name with version and architecture
// Example: grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64 -> grub2-tools-efi
func (d DnfManager) removeVersionFromPkgName(pkgNameWithVersion string) string {
	// Remove the version from the package name
	// Find the last hyphen before a digit (start of version)
	// Example: grub2-tools-efi-1:2.06-114.el9_7.alma.1.x86_64
	lastHyphen := -1
	for i := 0; i < len(pkgNameWithVersion); i++ {
		if pkgNameWithVersion[i] == '-' && i+1 < len(pkgNameWithVersion) && pkgNameWithVersion[i+1] >= '0' && pkgNameWithVersion[i+1] <= '9' {
			//fmt.Println(pkgNameWithVersion[:i])
			lastHyphen = i
			break
		}
	}
	if lastHyphen != -1 {
		return pkgNameWithVersion[:lastHyphen]
	}
	return pkgNameWithVersion
}

// RebootRequired checks if a reboot is required on the system
func (d DnfManager) RebootRequired(ctx context.Context) (bool, error) {
	// Check for /run/reboot-required or /var/run/reboot-required
	paths := []string{"/run/reboot-required", "/var/run/reboot-required"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true, nil
		}
	}

	// yum-utils and dnf-utils provide needs-restarting command
	timeout := 60 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "needs-restarting -r",
		Timeout: timeout,
		Env: map[string]string{
			"LANG": "C",
		},
	})

	if err != nil {
		// If the command is not found, we cannot determine if a reboot is required
		return false, err
	}

	// If the command exits with code 1, a reboot is required
	if result.RC == 1 {
		return true, nil
	}

	return false, nil
}
