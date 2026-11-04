package packagemanager

import (
	"context"
)

type Package struct {
	Name        string
	Version     string
	Description string
	// Architecture of the package, e.g. amd64 (apt), x86_64 (rpm, pacman) or noarch. A package can be installed once per
	// architecture, e.g. libc6 for amd64 and i386
	Architecture string `json:",omitempty"`
	// SourcePackage is the name of the source package this package was built from, e.g. curl for libcurl4t64
	SourcePackage string `json:",omitempty"`
	// SourceVersion is the version of the source package. It can differ from Version, e.g. bash 5.2.37-2+b10 is built from
	// the source 5.2.37-2. For rpm it has no epoch, because the source rpm file name does not contain one.
	SourceVersion string `json:",omitempty"`
	// SourceRPM is the file name of the source rpm, e.g. curl-8.18.0-8.fc44.src.rpm (only rpm based distributions)
	SourceRPM string `json:",omitempty"`
	// ModularityLabel is the module and stream of packages from an AppStream module on RHEL 8 and newer, e.g.
	// nginx:1.24:9080020260914161558:8cf767d6. Advisories of such packages are only valid for this stream.
	ModularityLabel string `json:",omitempty"`
}

type PackageUpdate struct {
	Name             string
	CurrentVersion   string
	AvailableVersion string
	IsSecurityUpdate bool
	IsPatch          bool // For distros that differentiate between security updates and patches (openSUSE)
}

type WindowsApp struct {
	Name      string
	Version   string
	Publisher string
}

type WindowsUpdate struct {
	Title            string
	Description      string
	KBArticleIDs     []string
	IsInstalled      bool
	IsSecurityUpdate bool
	IsOptional       bool
	UpdateID         string
	RevisionNumber   int
	RebootRequired   bool
}

type MacosUpdate struct {
	Name        string
	Description string
	Version     string
}

type PackageManager interface {
	IsAvailable() bool
	UpdateMetadata(ctx context.Context) error
	ListInstalledPackages(ctx context.Context) ([]Package, error)
	ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error)
	RebootRequired(ctx context.Context) (bool, error)
	CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error)
}

type WindowsManager interface {
	ListInstalledApps(ctx context.Context) ([]WindowsApp, error)
	ListAvailableUpdates(ctx context.Context) ([]WindowsUpdate, error)
	RebootRequired(ctx context.Context) (bool, error)
	CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error)
}

type MacOSManager interface {
	ListInstalledApps(ctx context.Context) ([]Package, error)
	ListAvailableUpdates(ctx context.Context) ([]MacosUpdate, error)
	CollectPackageInfo(ctx context.Context, limitDescriptionLength int64, enableUpdateCheck bool) (PackageInfo, error)
}

// truncateDescription truncates the given description to the specified limit.
func truncateDescription(desc string, limit int64) string {
	// -1 = No limit
	if limit <= -1 {
		return desc
	}

	// 0 = disable description
	if limit == 0 {
		return ""
	}

	if limit > 0 && int64(len(desc)) > limit {
		return desc[:limit]
	}
	return desc
}
