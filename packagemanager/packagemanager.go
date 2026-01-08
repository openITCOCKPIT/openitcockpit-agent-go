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

type PackageManager interface {
	IsAvailable() bool
	UpdateMetadata(ctx context.Context) error
	ListInstalledPackages(ctx context.Context) ([]Package, error)
	ListUpgradablePackages(ctx context.Context) ([]PackageUpdate, error)
	RebootRequired(ctx context.Context) (bool, error)
}
