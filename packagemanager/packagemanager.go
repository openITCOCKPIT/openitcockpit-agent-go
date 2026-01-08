package packagemanager

import "context"

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

type PackageManager interface {
	IsAvailable() bool
	UpdateMetadata(ctx context.Context) error
	ListInstalledPackages(ctx context.Context) ([]Package, error)
	ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error)
	RebootRequired(ctx context.Context) (bool, error)
}

type WindowsManager interface {
	ListInstalledApps(ctx context.Context) ([]WindowsApp, error)
	ListAvailableUpdates(ctx context.Context) ([]WindowsUpdate, error)
}
