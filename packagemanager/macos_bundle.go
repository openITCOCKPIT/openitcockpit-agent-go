package packagemanager

import (
	"os"
	"path/filepath"

	"howett.net/plist"
)

type macosInfoPlist struct {
	CFBundleIdentifier string `plist:"CFBundleIdentifier"`
}

// readBundleIdentifier returns the CFBundleIdentifier of the app at appPath, e.g. com.apple.Safari.
// The name of an app is localized (e.g. "Kalender" or "Calendar"), the bundle identifier is not.
// Vulnerability scanners use it to identify Apple apps.
//
// An empty string is returned if the identifier can not be read. This is not an error, as some apps
// have no Info.plist or are not accessible.
//
// This method is in here so we can test it easily on linux and windows systems
func readBundleIdentifier(appPath string) string {
	if appPath == "" {
		return ""
	}

	content, err := os.ReadFile(filepath.Join(appPath, "Contents", "Info.plist"))
	if err != nil {
		return ""
	}

	// Info.plist files can be XML or binary, plist.Unmarshal handles both formats
	var info macosInfoPlist
	if _, err := plist.Unmarshal(content, &info); err != nil {
		return ""
	}

	return info.CFBundleIdentifier
}
