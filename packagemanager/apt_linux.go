package packagemanager

import (
	"os/exec"
	"strings"
)

// AptManager implements PackageManager for apt
type AptManager struct{}

func (a AptManager) ListInstalledPackages() ([]Package, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package} ${Version} ${Description}\n")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var pkgs []Package
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pkgs = append(pkgs, Package{Name: fields[0], Version: fields[1]})
		}
	}
	return pkgs, nil
}

func (a AptManager) ListUpgradablePackages() ([]PackageUpdate, error) {
	cmd := exec.Command("apt", "list", "--upgradable")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var updates []PackageUpdate
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			nameVer := strings.Split(fields[0], "/")
			if len(nameVer) == 2 {
				updates = append(updates, PackageUpdate{
					Name:             nameVer[0],
					AvailableVersion: fields[1],
					CurrentVersion:   "", // apt does not show current version here
				})
			}
		}
	}
	return updates, nil
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
