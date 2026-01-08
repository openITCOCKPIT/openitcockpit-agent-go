package packagemanager

import "strings"

// parseSoftwareUpdateOutput parses the output of the softwareupdate --list command
// This method is in here so we can test it easily on windows and linux systems
func parseSoftwareUpdateOutput(output string) ([]MacosUpdate, error) {
	// This is the format we need to parse
	//* Label: Command Line Tools for Xcode 26.2-26.2
	//	Title: Command Line Tools for Xcode 26.2, Version: 26.2, Size: 858715KiB, Recommended: YES,
	//* Label: macOS Tahoe 26.2-25C56
	//	Title: macOS Tahoe 26.2, Version: 26.2, Size: 3693075KiB, Recommended: YES, Action: restart,

	var updates []MacosUpdate
	lines := strings.Split(output, "\n")
	var currentUpdate MacosUpdate
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "* Label:") {
			// If we have a current update, append it to the list
			if currentUpdate.Name != "" {
				updates = append(updates, currentUpdate)
			}
			// Start a new update
			currentUpdate = MacosUpdate{}
			labelParts := strings.SplitN(line, ": ", 2)
			if len(labelParts) == 2 {
				currentUpdate.Name = strings.TrimSpace(labelParts[1])
			}
		} else if strings.HasPrefix(line, "Title:") {
			titleParts := strings.SplitN(line, ",", 2)
			if len(titleParts) > 0 {
				titleSubParts := strings.SplitN(titleParts[0], ": ", 2)
				if len(titleSubParts) == 2 {
					currentUpdate.Description = strings.TrimSpace(titleSubParts[1])
				}
			}
			// Extract version
			versionParts := strings.SplitN(line, "Version: ", 2)
			if len(versionParts) == 2 {
				versionSubParts := strings.SplitN(versionParts[1], ",", 2)
				currentUpdate.Version = strings.TrimSpace(versionSubParts[0])
			}
		}
	}

	// Append the last update if exists
	if currentUpdate.Name != "" {
		updates = append(updates, currentUpdate)
	}

	return updates, nil
}
