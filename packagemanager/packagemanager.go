package packagemanager

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
}

type PackageManager interface {
	UpdateMetadata() error
	ListInstalledPackages() ([]Package, error)
	ListUpgradablePackages() ([]PackageUpdate, error)
}
