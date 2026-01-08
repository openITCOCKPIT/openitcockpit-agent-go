package packagemanager

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

// RpmManager implements PackageManager for rpm
// While rpm itself is not a package manager, the "rpm" command can be used
// to query installed packages on RPM-based systems.
type RpmManager struct{}

// IsAvailable checks if rpm is available on the system
func (r RpmManager) IsAvailable() bool {
	_, err := exec.LookPath("rpm")
	return err == nil
}

// UpdateMetadata updates the package metadata using yum or dnf
func (r RpmManager) UpdateMetadata(ctx context.Context) error {
	// RPM itself does not manage metadata, so we do nothing here.
	return nil
}

// ListInstalledPackages lists installed packages using rpm
func (r RpmManager) ListInstalledPackages(ctx context.Context) ([]Package, error) {
	output, err := r.getInstalledPackagesWithCancel(ctx)
	if err != nil {
		return nil, err
	}

	return r.parseRpmOutput(output)
}

func (r RpmManager) getInstalledPackagesWithCancel(ctx context.Context) (string, error) {
	// http://ftp.rpm.org/api/4.4.2.2/queryformat.html
	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "rpm -qa --qf '%{NAME} %{EPOCHNUM}:%{VERSION}-%{RELEASE} %{DESCRIPTION}\nEND\n'",
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

// parseRpmOutput parses the output of 'rpm -qa' into a slice of Package structs
func (r RpmManager) parseRpmOutput(output string) ([]Package, error) {
	lines := strings.Split(output, "\nEND\n")
	var pkgs []Package
	for _, line := range lines {
		// package version description
		// tar 2:1.34-7.el9 The GNU tar program saves many files together in one archive and can <until END>
		parts := strings.SplitN(line, " ", 3)

		if len(parts) < 3 {
			continue
		}

		version := r.removeEpochFromVersionIfZero(parts[1])

		pkgs = append(pkgs, Package{
			Name:        parts[0],
			Version:     version,
			Description: parts[2],
		})
	}
	return pkgs, nil
}

// ListUpgradablePackages lists upgradable packages using yum or dnf
func (r RpmManager) ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error) {
	// RPM itself does not handle upgrades, so we return an empty list.
	return []PackageUpdate{}, nil
}

// removeEpochFromVersionIfZero removes the epoch prefix from a version string if it is zero
func (r RpmManager) removeEpochFromVersionIfZero(version string) string {
	// To have the same version format as dnf/yum, we have to remove the epoch prefix if it is zero
	// but keep it if it is non-zero
	// Examples:
	// glibc 0:2.34-231.el9_7.2 -> 2.34-231.el9_7.2
	// bash 1:5.1.8-9.el9 -> 1:5.1.8-9.el9
	if after, ok := strings.CutPrefix(version, "0:"); ok {
		version = after
	}

	return version
}
