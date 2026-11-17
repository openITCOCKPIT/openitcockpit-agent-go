package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const dropInMainCustomChecks = `[default]

[check_disk]
  command = /usr/lib/nagios/plugins/check_disk -w 10% -c 5%
  interval = 60
  enabled = true

[check_load]
  command = /usr/lib/nagios/plugins/check_load -w 5 -c 10
  enabled = true

[check_users]
  command = /usr/lib/nagios/plugins/check_users -w 5 -c 10
  enabled = true
`

// Deployed to selected hosts, e.g. by Ansible
const dropInAppCustomChecks = `[check_app]
  command = /opt/app/bin/check_app
  timeout = 30
  enabled = true
`

// Changes and disables checks of the main file
const dropInOverrideCustomChecks = `[check_disk]
  command = /usr/lib/nagios/plugins/check_disk -w 20% -c 10%
  enabled = true

[check_users]
  enabled = false
`

// writeDropInTestFiles writes the given files (relative path => content) to a new temporary directory
func writeDropInTestFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func customCheckCommands(checks []*CustomCheck) map[string]string {
	commands := map[string]string{}
	for _, check := range checks {
		commands[check.Name] = check.Command
	}
	return commands
}

func TestDropInDirectory(t *testing.T) {
	tests := map[string]string{
		filepath.Join("etc", "openitcockpit-agent", "customchecks.ini"):         filepath.Join("etc", "openitcockpit-agent", "customchecks.d"),
		filepath.Join("etc", "openitcockpit-agent", "prometheus_exporters.ini"): filepath.Join("etc", "openitcockpit-agent", "prometheus_exporters.d"),
		filepath.Join("etc", "openitcockpit-agent", "checks"):                   filepath.Join("etc", "openitcockpit-agent", "checks.d"),
	}
	for configPath, expected := range tests {
		if directory := dropInDirectory(configPath); directory != expected {
			t.Errorf("%s: expected %s, got %s", configPath, expected, directory)
		}
	}
}

func TestDefinitionFiles(t *testing.T) {
	dir := writeDropInTestFiles(t, map[string]string{
		"customchecks.ini":                 dropInMainCustomChecks,
		"customchecks.d/20-override.ini":   dropInOverrideCustomChecks,
		"customchecks.d/10-app.INI":        dropInAppCustomChecks,
		"customchecks.d/.10-app.ini.swp":   "editor swap file",
		"customchecks.d/.hidden.ini":       dropInAppCustomChecks,
		"customchecks.d/README.txt":        "not an ini file",
		"customchecks.d/30-old.ini.bak":    dropInAppCustomChecks,
		"customchecks.d/sub.ini/check.ini": dropInAppCustomChecks,
	})

	files := definitionFiles(filepath.Join(dir, "customchecks.ini"))
	expected := []string{
		filepath.Join(dir, "customchecks.ini"),
		filepath.Join(dir, "customchecks.d", "10-app.INI"),
		filepath.Join(dir, "customchecks.d", "20-override.ini"),
	}
	if fmt.Sprint(files) != fmt.Sprint(expected) {
		t.Errorf("expected %v, got %v", expected, files)
	}

	if len(definitionFiles(filepath.Join(dir, "http.ini"))) != 0 {
		t.Error("expected no files without http.ini and http.d")
	}
	if len(definitionFiles("")) != 0 {
		t.Error("expected no files for an empty path")
	}
}

func TestReadCustomChecksDropIn(t *testing.T) {
	dir := writeDropInTestFiles(t, map[string]string{
		"customchecks.ini":               dropInMainCustomChecks,
		"customchecks.d/10-app.ini":      dropInAppCustomChecks,
		"customchecks.d/20-override.ini": dropInOverrideCustomChecks,
	})

	checks, err := unmarshalCustomChecks(filepath.Join(dir, "customchecks.ini"))
	if err != nil {
		t.Fatal(err)
	}

	commands := customCheckCommands(checks)
	expected := map[string]string{
		"check_disk": "/usr/lib/nagios/plugins/check_disk -w 20% -c 10%", // replaced by 20-override.ini
		"check_load": "/usr/lib/nagios/plugins/check_load -w 5 -c 10",
		"check_app":  "/opt/app/bin/check_app",
		// check_users is disabled by 20-override.ini
	}
	if fmt.Sprint(commands) != fmt.Sprint(expected) {
		t.Errorf("expected %v, got %v", expected, commands)
	}

	for _, check := range checks {
		// The replaced section does not keep the interval of the main file
		if check.Name == "check_disk" && check.Interval != 60 {
			t.Errorf("check_disk: expected the default interval 60, got %d", check.Interval)
		}
		if check.Name == "check_app" && check.Timeout != 30 {
			t.Errorf("check_app: expected timeout 30, got %d", check.Timeout)
		}
	}
}

func TestReadCustomChecksDropInOnly(t *testing.T) {
	dir := writeDropInTestFiles(t, map[string]string{
		"customchecks.d/10-app.ini": dropInAppCustomChecks,
	})

	configPath := filepath.Join(dir, "customchecks.ini")
	if !hasDefinitionFiles(configPath) {
		t.Fatal("expected definition files in customchecks.d")
	}
	checks, err := unmarshalCustomChecks(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || checks[0].Name != "check_app" {
		t.Errorf("expected only check_app, got %v", customCheckCommands(checks))
	}
}

func TestReadIniSectionsSkipsInvalidFiles(t *testing.T) {
	dir := writeDropInTestFiles(t, map[string]string{
		"customchecks.ini":             dropInMainCustomChecks,
		"customchecks.d/10-app.ini":    dropInAppCustomChecks,
		"customchecks.d/20-broken.ini": "[unclosed section\ncommand = /bin/false\n",
	})

	checks, err := unmarshalCustomChecks(filepath.Join(dir, "customchecks.ini"))
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0)
	for _, check := range checks {
		names = append(names, check.Name)
	}
	sort.Strings(names)
	if fmt.Sprint(names) != "[check_app check_disk check_load check_users]" {
		t.Errorf("expected the checks of the valid files, got %v", names)
	}

	brokenOnly := writeDropInTestFiles(t, map[string]string{
		"customchecks.d/20-broken.ini": "[unclosed section\ncommand = /bin/false\n",
	})
	if _, err := unmarshalCustomChecks(filepath.Join(brokenOnly, "customchecks.ini")); err == nil {
		t.Error("expected an error if no file can be read")
	}
}

func TestLoadConfigurationWithDropInDirectories(t *testing.T) {
	dir := writeDropInTestFiles(t, map[string]string{
		"customchecks.d/10-app.ini": dropInAppCustomChecks,
		"http.ini": `[intranet]
  url = https://intranet.example.com/
  enabled = true
`,
		"http.d/10-shop.ini": `[shop]
  url = https://shop.example.com/health
  enabled = true
`,
		"certificates.d/10-mail.ini": `[mail]
  type = tls
  address = mail.example.com:993
  enabled = true
`,
	})
	agentConfig := fmt.Sprintf(`[default]
  customchecks = %s
  http = %s
  certificates = %s
  logfiles = %s
`, filepath.Join(dir, "customchecks.ini"), filepath.Join(dir, "http.ini"), filepath.Join(dir, "certificates.ini"), filepath.Join(dir, "logfiles.ini"))
	configPath := filepath.Join(dir, "config.ini")
	if err := os.WriteFile(configPath, []byte(agentConfig), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.CustomCheckConfiguration) != 1 || cfg.CustomCheckConfiguration[0].Name != "check_app" {
		t.Errorf("expected the custom check of customchecks.d, got %v", customCheckCommands(cfg.CustomCheckConfiguration))
	}
	httpNames := make([]string, 0)
	for _, definition := range cfg.HTTPConfiguration {
		httpNames = append(httpNames, definition.Name)
	}
	sort.Strings(httpNames)
	if fmt.Sprint(httpNames) != "[intranet shop]" {
		t.Errorf("expected the http definitions of http.ini and http.d, got %v", httpNames)
	}
	if len(cfg.CertificateConfiguration) != 1 || cfg.CertificateConfiguration[0].Name != "mail" {
		t.Errorf("expected the certificate definition of certificates.d, got %d definitions", len(cfg.CertificateConfiguration))
	}
	if len(cfg.LogfileConfiguration) != 0 {
		t.Errorf("expected no logfile definitions, got %d", len(cfg.LogfileConfiguration))
	}
}
