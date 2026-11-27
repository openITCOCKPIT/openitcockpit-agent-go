package packagemanager

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
)

// A KB number has at least 6 digits, e.g. 890830 or 5079473
var kbNumberRegex = regexp.MustCompile(`\d{6,}`)

// parseInstalledKBsOutput parses the JSON output of the PowerShell script in getInstalledKBs.
// The script returns the KBs of multiple sources in different notations, e.g. "KB5079473" and "5079473".
// Because of that the result only contains the number of each KB once.
//
// The result is sorted to keep the order stable, the server calculates a hash of the inventory.
//
// This method is in here so we can test it easily on linux and macos systems
func parseInstalledKBsOutput(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return []string{}, nil
	}

	var entries []string
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		return nil, err
	}

	numbers := make([]string, 0, len(entries))
	for _, entry := range entries {
		number := kbNumberRegex.FindString(entry)
		if number != "" {
			numbers = append(numbers, number)
		}
	}

	slices.Sort(numbers)
	return slices.Compact(numbers), nil
}
