package packagemanager

import (
	"reflect"
	"testing"
)

func TestParseInstalledKBsOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected []string
	}{
		{
			name:     "both notations of the same KB are merged",
			output:   `["KB5079473","5079473"]`,
			expected: []string{"5079473"},
		},
		{
			name:     "result is sorted",
			output:   `["KB5085516","890830","KB5079473"]`,
			expected: []string{"5079473", "5085516", "890830"},
		},
		{
			name:     "lower case prefix",
			output:   `["kb5079473"]`,
			expected: []string{"5079473"},
		},
		{
			name:     "entries without a KB number are ignored",
			output:   `["File 1","","KB123","KB5079473"]`,
			expected: []string{"5079473"},
		},
		{
			name:     "no KBs installed",
			output:   `[]`,
			expected: []string{},
		},
		{
			name:     "empty output",
			output:   "",
			expected: []string{},
		},
		{
			name:     "trailing line break of PowerShell",
			output:   "[\"KB5079473\"]\r\n",
			expected: []string{"5079473"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := parseInstalledKBsOutput(test.output)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestParseInstalledKBsOutput_InvalidJson(t *testing.T) {
	if _, err := parseInstalledKBsOutput("Not JSON"); err == nil {
		t.Error("expected an error for invalid JSON")
	}
}

func TestParseInstalledKBsOutput_StableOrder(t *testing.T) {
	first, _ := parseInstalledKBsOutput(`["KB5079473","KB5085516"]`)
	second, _ := parseInstalledKBsOutput(`["KB5085516","KB5079473"]`)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("order must not depend on the input: %v, %v", first, second)
	}
}
