package packagemanager

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	"github.com/shirou/gopsutil/v4/host"
)

// rpmQueryFormat is the query format of rpm. The description has to be the last field, as it has multiple lines.
// http://ftp.rpm.org/api/4.4.2.2/queryformat.html
const rpmQueryFormat = "%{NAME}\t%{EPOCHNUM}\t%{VERSION}\t%{RELEASE}\t%{ARCH}\t%{SOURCERPM}\t%{DESCRIPTION}\nEND\n"

// rpmQueryFieldCount is the number of tab separated fields of rpmQueryFormat
const rpmQueryFieldCount = 7

// rpmNoValue is printed by rpm for tags without a value
const rpmNoValue = "(none)"

// rpmPublicKeyPackage is the name of the pseudo package of each imported GPG key, it is no installed software and neither
// dnf nor yum list it
const rpmPublicKeyPackage = "gpg-pubkey"

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
	timeout := 10 * time.Second
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "rpm -qa --qf '" + rpmQueryFormat + "'",
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
	records := strings.Split(output, "\nEND\n")
	var pkgs []Package
	for _, record := range records {
		// name<TAB>epoch<TAB>version<TAB>release<TAB>architecture<TAB>source rpm<TAB>description
		// tar	2	1.34	7.el9	x86_64	tar-1.34-7.el9.src.rpm	The GNU tar program saves many files together <until END>
		parts := strings.SplitN(record, "\t", rpmQueryFieldCount)

		if len(parts) < rpmQueryFieldCount {
			continue
		}

		name := parts[0]
		if name == rpmPublicKeyPackage {
			continue
		}

		version := r.removeEpochFromVersionIfZero(parts[1] + ":" + parts[2] + "-" + parts[3])
		sourceRpm := r.removeNoValue(parts[5])
		sourcePackage, sourceVersion := r.parseSourceRpm(sourceRpm)

		pkgs = append(pkgs, Package{
			Name:          name,
			Version:       version,
			Architecture:  r.removeNoValue(parts[4]),
			SourcePackage: sourcePackage,
			SourceVersion: sourceVersion,
			SourceRPM:     sourceRpm,
			Description:   parts[6],
		})
	}
	return pkgs, nil
}

// removeNoValue returns an empty string for tags without a value
func (r RpmManager) removeNoValue(value string) string {
	if value == rpmNoValue {
		return ""
	}
	return value
}

// parseSourceRpm returns the name and the version-release of a source rpm file name.
// The name can contain dashes, the version and the release can not:
// openssl-3-3.5.0-160000.11.1.src.rpm -> openssl-3, 3.5.0-160000.11.1
func (r RpmManager) parseSourceRpm(sourceRpm string) (string, string) {
	nameVersionRelease, isSource := strings.CutSuffix(sourceRpm, ".src.rpm")
	if !isSource {
		// Source packages without the sources, e.g. for software which must not be redistributed
		nameVersionRelease, isSource = strings.CutSuffix(sourceRpm, ".nosrc.rpm")
	}
	if !isSource {
		return "", ""
	}

	releaseSeparator := strings.LastIndex(nameVersionRelease, "-")
	if releaseSeparator == -1 {
		return "", ""
	}

	versionSeparator := strings.LastIndex(nameVersionRelease[:releaseSeparator], "-")
	if versionSeparator == -1 {
		return "", ""
	}

	return nameVersionRelease[:versionSeparator], nameVersionRelease[versionSeparator+1:]
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

func (r RpmManager) RebootRequired(ctx context.Context) (bool, error) {
	// RpmManager is to generic to check for reboot requirements
	return false, nil
}

func (r RpmManager) CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return PackageInfo{
			Stats: PackageStats{
				LastError: err,
			},
		}, err
	}

	// Platform: "almalinux"
	// PlatformFamily: "rhel"
	// PlatformVersion: "9.7"
	// KernelVersion: "6.6.87.2-microsoft-standard-WSL2"
	result := PackageInfo{
		Enabled:    true,
		Pending:    false,
		LastUpdate: time.Now().Unix(),
		Stats: PackageStats{
			PackageManager:  "rpm",
			OperatingSystem: "linux",
			OsName:          info.Platform,
			OsVersion:       info.PlatformVersion,
			OsFamily:        info.PlatformFamily,
			Uptime:          int64(info.Uptime),
			AgentVersion:    config.AgentVersion,
		},
	}

	installedPackages, err := r.ListInstalledPackages(ctx)
	if err != nil {
		result.Stats.LastError = err
		return result, err
	}
	result.Stats.InstalledPackages = int64(len(installedPackages))

	// Truncate descriptions if needed
	for i := range installedPackages {
		installedPackages[i].Description = truncateDescription(installedPackages[i].Description, limitDescriptionLength)
	}

	result.LinuxPackages = installedPackages
	return result, nil
}
