package packagemanager

import (
	"os"
	"path/filepath"
	"testing"

	"howett.net/plist"
)

// createApp creates a fake .app bundle with the given Info.plist content and returns its path
func createApp(t *testing.T, infoPlist []byte) string {
	t.Helper()
	appPath := filepath.Join(t.TempDir(), "Test.app")
	if err := os.MkdirAll(filepath.Join(appPath, "Contents"), 0755); err != nil {
		t.Fatal(err)
	}
	if infoPlist != nil {
		if err := os.WriteFile(filepath.Join(appPath, "Contents", "Info.plist"), infoPlist, 0644); err != nil {
			t.Fatal(err)
		}
	}
	return appPath
}

func TestReadBundleIdentifier_XmlPlist(t *testing.T) {
	appPath := createApp(t, []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>com.apple.iCal</string>
	<key>CFBundleName</key>
	<string>Kalender</string>
</dict>
</plist>`))

	if result := readBundleIdentifier(appPath); result != "com.apple.iCal" {
		t.Errorf("expected com.apple.iCal, got %q", result)
	}
}

func TestReadBundleIdentifier_BinaryPlist(t *testing.T) {
	content, err := plist.Marshal(map[string]string{"CFBundleIdentifier": "com.google.Chrome"}, plist.BinaryFormat)
	if err != nil {
		t.Fatal(err)
	}
	appPath := createApp(t, content)

	if result := readBundleIdentifier(appPath); result != "com.google.Chrome" {
		t.Errorf("expected com.google.Chrome, got %q", result)
	}
}

func TestReadBundleIdentifier_MissingOrInvalid(t *testing.T) {
	tests := map[string]string{
		"empty path":         "",
		"missing app":        filepath.Join(t.TempDir(), "Missing.app"),
		"missing Info.plist": createApp(t, nil),
		"invalid Info.plist": createApp(t, []byte("not a plist")),
		"no identifier":      createApp(t, []byte(`<plist version="1.0"><dict><key>CFBundleName</key><string>Test</string></dict></plist>`)),
	}

	for name, appPath := range tests {
		t.Run(name, func(t *testing.T) {
			if result := readBundleIdentifier(appPath); result != "" {
				t.Errorf("expected an empty identifier, got %q", result)
			}
		})
	}
}

func TestParseMacOSInstalledAppsOutput_BundleID(t *testing.T) {
	appPath := createApp(t, []byte(`<plist version="1.0"><dict><key>CFBundleIdentifier</key><string>com.apple.Safari</string></dict></plist>`))
	output := `{"SPApplicationsDataType": [
		{"_name": "Safari", "version": "26.1", "path": "` + filepath.ToSlash(appPath) + `"},
		{"_name": "Ohne Pfad", "version": "1.0"}
	]}`

	apps, err := parseMacOSInstalledAppsOutput(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 || apps[0].BundleID != "com.apple.Safari" || apps[1].BundleID != "" {
		t.Errorf("unexpected bundle identifiers: %+v", apps)
	}
}
