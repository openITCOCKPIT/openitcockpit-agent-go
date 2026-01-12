package packagemanager

import (
	"context"
)

type Package struct {
	Name        string
	Version     string
	Description string
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
