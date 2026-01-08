package packagemanager

import (
	"testing"
)

var softwareUpdateOutputSample = `Software Update Tool

Finding available software
Software Update found the following new or updated software:
* Label: Command Line Tools for Xcode 26.2-26.2
	Title: Command Line Tools for Xcode 26.2, Version: 26.2, Size: 858715KiB, Recommended: YES,
* Label: macOS Tahoe 26.2-25C56
	Title: macOS Tahoe 26.2, Version: 26.2, Size: 3693075KiB, Recommended: YES, Action: restart,

`

func TestParseSoftwareUpdateOutput(t *testing.T) {
	updates, err := parseSoftwareUpdateOutput(softwareUpdateOutputSample)
	if err != nil {
		t.Fatalf("Error parsing software update output: %v", err)
	}

	if len(updates) != 2 {
		t.Fatalf("Expected 2 updates, got %d", len(updates))
	}

	expectedFirst := MacosUpdate{
		Name:        "Command Line Tools for Xcode 26.2-26.2",
		Description: "Command Line Tools for Xcode 26.2",
		Version:     "26.2",
	}

	if updates[0] != expectedFirst {
		t.Errorf("First update does not match expected.\nGot: %+v\nExpected: %+v", updates[0], expectedFirst)
	}

	expectedSecond := MacosUpdate{
		Name:        "macOS Tahoe 26.2-25C56",
		Description: "macOS Tahoe 26.2",
		Version:     "26.2",
	}

	if updates[1] != expectedSecond {
		t.Errorf("Second update does not match expected.\nGot: %+v\nExpected: %+v", updates[1], expectedSecond)
	}
}
