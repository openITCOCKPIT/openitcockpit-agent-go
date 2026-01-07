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
}

type PackageManager interface {
	ListInstalledPackages() ([]Package, error)
	ListUpgradablePackages() ([]PackageUpdate, error)
}
