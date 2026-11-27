package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/platformpaths"
)

var agentVersion1ConfigBlank string = `[default]
interval = 30
# Port of the Agents build-in web server
port = 3333
address = 0.0.0.0
certfile =
keyfile =
try-autossl = true
autossl-csr-file =
autossl-crt-file =
autossl-key-file =
autossl-ca-file =
verbose = false
stacktrace = false
config-update-mode = false
auth =
customchecks = /etc/openitcockpit-agent/customchecks.ini
temperature-fahrenheit = false
dockerstats = false
qemustats = false
cpustats = true
sensorstats = true
processstats = true
processstats-including-child-ids = false
netstats = true
diskstats = true
netio = true
diskio = true
winservices = true
systemdservices = true
wineventlog = true
wineventlog-logtypes = System, Application, Security
alfrescostats = false
alfresco-jmxuser = monitorRole
alfresco-jmxpassword = change_asap
alfresco-jmxaddress = 0.0.0.0
alfresco-jmxport = 50500
alfresco-jmxpath = /alfresco/jmxrmi
alfresco-jmxquery =
alfresco-javapath = /usr/bin/java

[oitc]
enabled = false
hostuuid =
url =
apikey =
proxy =
interval = 60
`

var agentVersion1Config string = `[default]
interval = 60
# Port of the Agents build-in web server
port = 33333
address = 127.0.0.1
certfile = /foo/bar.cert
keyfile = /foo/bar.key
try-autossl = true
autossl-csr-file =  /etc/autossl/csr.csr
autossl-crt-file = /etc/autossl/crt.crt
autossl-key-file = /etc/autossl/key.key
autossl-ca-file = /etc/autossl/server_ca.ca
verbose = false
stacktrace = false
config-update-mode = false
auth = username:pass:word
customchecks = C:\Program Files\it-novum\openitcockpit-agent\customchecks.ini
temperature-fahrenheit = false
dockerstats = false
qemustats = false
cpustats = false
sensorstats = true
processstats = true
processstats-including-child-ids = false
netstats = true
diskstats = true
netio = true
diskio = true
winservices = true
systemdservices = true
wineventlog = true
wineventlog-logtypes = System, Application, Security
alfrescostats = false
alfresco-jmxuser = oitc-agent
alfresco-jmxpassword = change_asap
alfresco-jmxaddress = 0.0.0.0
alfresco-jmxport = 50500
alfresco-jmxpath = /alfresco/jmxrmi
alfresco-jmxquery =
alfresco-javapath = /usr/bin/java

[oitc]
enabled = true
hostuuid = 3a2d91e5-03f7-4d2c-b719-05bd69b312ee
url = https://demo.openitcockpit.io
apikey = aaaaabbbbbcccccdddddeeeeefffff
proxy = proxy.example.org
interval = 90
`

var agentVersion1ConfigEmpty = ""

var agentConfigWithCustomCheck string = `[default]
customchecks = "%s"
`

var agentConfigWithCustomCheck2 string = `[default]
interval = 5
customchecks = "%s"
`

var customChecksAgentVersion1Config string = `[default]
  # max_worker_threads should be increased with increasing number of custom checks
  # but consider: each thread needs (a bit) memory
  max_worker_threads = 10

[time_1]
  command = "C:\checks\check_time.exe"
  interval = 60
  timeout = 10
  enabled = true

[check_Windows_Services_Status_OSS]
  command = powershell.exe -nologo -noprofile -File "C:\checks\check_Windows_Services_Status_OSS.ps1"
  interval = 15
  timeout = 10
  enabled = false

[check_ping]
  command = /usr/lib/nagios/plugins/check_ping -H 127.0.0.1 -w 100.0,20%% -c 500.0,60%% -p 5
  interval = 15
  timeout = 10
  enabled = true

[check_users]
  command = /usr/lib/nagios/plugins/check_users -w 3 -c 7
  interval = 15
  timeout = 10
  enabled = true
`

var customChecksAgentEmptyConfig string = ``

var customChecksAgentVersion1ConfigEmptyCommand string = `
[time_1]
  command = "C:\checks\check_time.exe"
  interval = 60
  timeout = 10
  enabled = true

[empty_command_line]
  command =
  interval = 15
  timeout = 10
  enabled = false
`

var customChecksAgentVersion1ConfigMissingCommand string = `
[time_1]
  command = "C:\checks\check_time.exe"
  interval = 60
  timeout = 10
  enabled = true

[no_command_line_at_all]
timeout = 10
enabled = true
`

func saveTempConfig(config string, customchecks bool) string {
	filename := "config.ini"
	if customchecks {
		filename = "customchecks.ini"
	}
	tmpDir, err := os.MkdirTemp(os.TempDir(), "*-test")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(config), 0600); err != nil {
		panic(err)
	}
	return tmpDir
}

func saveTempConfigWithCC(config string, customchecks string) string {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "*-test")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "config.ini"), []byte(fmt.Sprintf(config, filepath.Join(tmpDir, "customchecks.ini"))), 0600); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "customchecks.ini"), []byte(customchecks), 0600); err != nil {
		panic(err)
	}
	return tmpDir
}

func TestAgentVersion1BlankConfig(t *testing.T) {
	cfgdir := saveTempConfig(agentVersion1ConfigBlank, false)
	defer os.RemoveAll(cfgdir)

	configPath := filepath.Join(cfgdir, "config.ini")
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}

	if c.CheckInterval != 30 {
		t.Error("Check interval expect to be 30")
	}

	if c.Port != 3333 {
		t.Error("WebServer port expect to be 3333")
	}

	if c.CustomchecksFilePath != "/etc/openitcockpit-agent/customchecks.ini" {
		t.Error("WebServer port expect to be /etc/openitcockpit-agent/customchecks.ini")
	}

	if c.CPU != true {
		t.Error("Checks CPU expect to be true")
	}

	if c.JmxUser != "monitorRole" {
		t.Error("Alfresco JmxUser expect to be monitorRole")
	}

	if c.OITC.Push != false {
		t.Error("Push Mode expect to be false")
	}

	js, _ := json.MarshalIndent(c, "", "    ")
	fmt.Println(string(js))
}

func TestAgentVersion1EmptyConfig(t *testing.T) {
	cfgdir := saveTempConfig(agentVersion1ConfigEmpty, false)
	defer os.RemoveAll(cfgdir)

	configPath := filepath.Join(cfgdir, "config.ini")
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}

	if c.CheckInterval != 30 {
		t.Error("Check interval expect to be 30")
	}

	if c.Port != 3333 {
		t.Error("WebServer port expect to be 3333")
	}

	if c.CustomchecksFilePath != filepath.Join(platformpaths.Get().ConfigPath(), "customchecks.ini") {
		t.Error("WebServer port expect to be: ", filepath.Join(platformpaths.Get().ConfigPath(), "customchecks.ini"))
	}

	if c.CPU != true {
		t.Error("Checks CPU expect to be true")
	}

	if c.JmxUser != "" {
		t.Error("Alfresco JmxUser expect to be monitorRole")
	}

	if c.OITC.Push != false {
		t.Error("Push Mode expect to be false")
	}

	js, _ := json.MarshalIndent(c, "", "    ")
	fmt.Println(string(js))

}

func TestAgentVersion1Config(t *testing.T) {
	cfgdir := saveTempConfig(agentVersion1Config, false)
	defer os.RemoveAll(cfgdir)

	configPath := filepath.Join(cfgdir, "config.ini")
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}

	if c.CheckInterval != 60 {
		t.Error("Check interval expect to be 60")
	}

	if c.Port != 33333 {
		t.Error("WebServer port expect to be 33333")
	}

	if c.CertificateFile != "/foo/bar.cert" {
		t.Error("TLS CertificateFile expect to be /foo/bar.cert")
	}

	if c.KeyFile != "/foo/bar.key" {
		t.Error("TLS KeyFile expect to be /foo/bar.key")
	}

	if c.AutoSslCsrFile != "/etc/autossl/csr.csr" {
		t.Error("TLS AutoSslCsrFile expect to be /etc/autossl/csr.csr")
	}

	if c.AutoSslCrtFile != "/etc/autossl/crt.crt" {
		t.Error("TLS AutoSslCrtFile expect to be /etc/autossl/crt.crt")
	}

	if c.AutoSslKeyFile != "/etc/autossl/key.key" {
		t.Error("TLS AutoSslKeyFile expect to be /etc/autossl/key.key")
	}

	if c.AutoSslCaFile != "/etc/autossl/server_ca.ca" {
		t.Error("TLS AutoSslCaFile expect to be /etc/autossl/server_ca.ca")
	}

	if c.CustomchecksFilePath != "C:\\Program Files\\it-novum\\openitcockpit-agent\\customchecks.ini" {
		t.Error("WebServer port expect to be C:\\Program Files\\it-novum\\openitcockpit-agent\\customchecks.ini")
	}

	if c.CPU != false {
		t.Error("Checks CPU expect to be false")
	}

	if c.JmxUser != "oitc-agent" {
		t.Error("Alfresco JmxUser expect to be oitc-agent")
	}

	if c.OITC.HostUUID != "3a2d91e5-03f7-4d2c-b719-05bd69b312ee" {
		t.Error("Push HostUUID expect to be 3a2d91e5-03f7-4d2c-b719-05bd69b312ee")
	}

	if c.OITC.URL != "https://demo.openitcockpit.io" {
		t.Error("Push url expect to be https://demo.openitcockpit.io")
	}

	if c.OITC.Apikey != "aaaaabbbbbcccccdddddeeeeefffff" {
		t.Error("Push Apikey expect to be aaaaabbbbbcccccdddddeeeeefffff")
	}

	if c.OITC.Proxy != "proxy.example.org" {
		t.Error("Push HostUUID expect to be proxy.example.org")
	}

	if c.OITC.Push != true {
		t.Error("Push Mode expect to be true")
	}

	if c.BasicAuth != "username:pass:word" {
		t.Error("BasicAuth username expect to be 'username'")
	}

	if !strings.Contains(c.OITC.AuthFile, "auth.json") {
		t.Error("auth.json file not set")
	}

	js, _ := json.MarshalIndent(c, "", "    ")
	fmt.Println(string(js))

}

func TestReadConfigFromFile(t *testing.T) {
	dir, _ := os.Getwd()
	configPath := fmt.Sprintf("%s%s../example/config_example.ini", dir, string(os.PathSeparator))
	_, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadCustomChecksConfigAgentVersion1(t *testing.T) {
	// Read custom checks config
	cfgdir := saveTempConfig(customChecksAgentVersion1Config, true)
	defer os.RemoveAll(cfgdir)

	ccc, err := unmarshalCustomChecks(filepath.Join(cfgdir, "customchecks.ini"))
	if err != nil {
		t.Fatal(err)
	}

	if len(ccc) != 3 {
		t.Error("This config is expected to have 3 enabled custom checks")
	}

	for _, customcheck := range ccc {
		if customcheck.Name == "time_1" {
			if customcheck.Enabled != true {
				t.Error("Custom check time_1 is expected to be enabled")
			}
		}

		if customcheck.Name == "check_Windows_Services_Status_OSS" {
			if customcheck.Enabled != false {
				t.Error("Custom check time_1 is expected to be disabled")
			}
		}
	}
}

func TestReadCustomChecksConfigAgentVersion1EmptyCommandline(t *testing.T) {
	cfgdir := saveTempConfig(customChecksAgentVersion1ConfigEmptyCommand, true)
	defer os.RemoveAll(cfgdir)

	// The custom check without command is skipped, the other custom checks are loaded
	ccc, err := unmarshalCustomChecks(filepath.Join(cfgdir, "customchecks.ini"))
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if len(ccc) != 1 || ccc[0].Name != "time_1" {
		t.Fatalf("expected only the custom check time_1, got %+v", ccc)
	}
}

func TestReadCustomChecksConfigAgentVersion1MissingCommandline(t *testing.T) {
	cfgdir := saveTempConfig(customChecksAgentVersion1ConfigMissingCommand, true)
	defer os.RemoveAll(cfgdir)

	// The custom check without command is skipped, the other custom checks are loaded
	ccc, err := unmarshalCustomChecks(filepath.Join(cfgdir, "customchecks.ini"))
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if len(ccc) != 1 || ccc[0].Name != "time_1" {
		t.Fatalf("expected only the custom check time_1, got %+v", ccc)
	}
}

func TestReadCustomChecksConfigEmpty(t *testing.T) {
	cfgdir := saveTempConfig(customChecksAgentEmptyConfig, true)
	defer os.RemoveAll(cfgdir)

	ccc, err := unmarshalCustomChecks(filepath.Join(cfgdir, "customchecks.ini"))
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if len(ccc) != 0 {
		t.Fatal("unexpected number of custom checks (0): ", len(ccc))
	}
}

func TestReadAgentConfigWithCC(t *testing.T) {
	cfgdir := saveTempConfigWithCC(agentConfigWithCustomCheck, customChecksAgentVersion1Config)
	defer os.RemoveAll(cfgdir)

	configPath := filepath.Join(cfgdir, "config.ini")
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
	ccc := c.CustomCheckConfiguration

	if len(ccc) < 1 {
		t.Fatal("unexpected number of custom checks (>0): ", len(ccc))
	}
}

func TestReadAgentConfigWithCCAndNewConfig(t *testing.T) {
	cfgdir := saveTempConfigWithCC(agentConfigWithCustomCheck, customChecksAgentVersion1Config)
	defer os.RemoveAll(cfgdir)

	configPath := filepath.Join(cfgdir, "config.ini")
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
	ccc := c.CustomCheckConfiguration

	if len(ccc) < 1 {
		t.Fatal("unexpected number of custom checks (>0): ", len(ccc))
	}

	if err := c.SaveCustomCheckConfiguration([]byte(customChecksAgentEmptyConfig)); err != nil {
		t.Fatal(err)
	}
	if err := c.SaveConfiguration([]byte(agentConfigWithCustomCheck2)); err != nil {
		t.Fatal(err)
	}
	c, err = Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
	ccc = c.CustomCheckConfiguration

	if c.CheckInterval != 5 {
		t.Error("reload did not work, unexpected check interval (5): ", c.CheckInterval)
	}
	if len(ccc) != 0 {
		t.Error("reload did not work, unexpected number of custom checks (0): ", len(ccc))
	}
}

var logfilesConfig string = `[default]

[Nginx_Errors]
  path = /var/log/nginx/*error.log
  match = \[(crit|alert|emerg)\]
  exclude = favicon\.ico
  enabled = true

[ssh]
  type = journald
  unit = ssh.service
  priority = 0..4
  match = ` + "`Failed password #\\d+`" + `
  window = 600
  max_lines = 5
  max_read_bytes = 1024
  enabled = true

[disabled]
  path = /var/log/syslog
  enabled = false

[invalid_regex]
  path = /var/log/syslog
  match = (unclosed
  enabled = true

[missing_path]
  type = file
  enabled = true

[unknown_type]
  type = eventlog
  enabled = true
`

func TestReadLogfilesConfig(t *testing.T) {
	cfgdir := t.TempDir()
	path := filepath.Join(cfgdir, "logfiles.ini")
	if err := os.WriteFile(path, []byte(logfilesConfig), 0600); err != nil {
		t.Fatal(err)
	}

	definitions, err := unmarshalLogfiles(path)
	if err != nil {
		t.Fatal(err)
	}

	byName := map[string]*LogfileDefinition{}
	for _, definition := range definitions {
		byName[definition.Name] = definition
	}
	if len(byName) != 5 {
		t.Fatal("expected 5 enabled definitions, got: ", len(byName))
	}
	if _, ok := byName["disabled"]; ok {
		t.Error("disabled definition must be ignored")
	}

	// Section names are lower case
	nginx := byName["nginx_errors"]
	if nginx == nil {
		t.Fatal("missing nginx_errors")
	}
	if nginx.ConfigError != "" || nginx.Type != LogfileTypeFile || nginx.Window != 300 || nginx.MaxLines != 10 || nginx.MaxReadBytes != 16*1024*1024 {
		t.Errorf("unexpected defaults: %+v", nginx)
	}
	if !nginx.MatchRegexp.MatchString("2025/01/01 [crit] 123#0: foo") || !nginx.ExcludeRegexp.MatchString("/favicon.ico") {
		t.Errorf("unexpected regular expressions: %s %s", nginx.Match, nginx.Exclude)
	}

	ssh := byName["ssh"]
	if ssh.ConfigError != "" || ssh.Type != LogfileTypeJournald || ssh.Unit != "ssh.service" || ssh.Priority != "0..4" ||
		ssh.Window != 600 || ssh.MaxLines != 5 || ssh.MaxReadBytes != 1024 {
		t.Errorf("unexpected definition: %+v", ssh)
	}
	if ssh.Match != `Failed password #\d+` {
		t.Errorf("unexpected match: %s", ssh.Match)
	}

	for name, expected := range map[string]string{
		"invalid_regex": "invalid match regex",
		"missing_path":  "missing path",
		"unknown_type":  "unknown type",
	} {
		if !strings.Contains(byName[name].ConfigError, expected) {
			t.Errorf("expected error '%s' for %s, got: '%s'", expected, name, byName[name].ConfigError)
		}
	}
}

func TestReadAgentConfigWithLogfiles(t *testing.T) {
	cfgdir := t.TempDir()
	configPath := filepath.Join(cfgdir, "config.ini")
	logfilesPath := filepath.Join(cfgdir, "logfiles.ini")
	if err := os.WriteFile(configPath, []byte("[default]\nlogfiles = "+logfilesPath+"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// A missing logfiles.ini is not an error
	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
	if c.LogfileConfiguration == nil || len(c.LogfileConfiguration) != 0 {
		t.Fatal("expected empty logfile configuration")
	}

	if err := c.SaveLogfileConfiguration([]byte(logfilesConfig)); err != nil {
		t.Fatal(err)
	}
	if string(c.ReadLogfileConfiguration()) != logfilesConfig {
		t.Fatal("unexpected logfile configuration")
	}

	c, err = Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.LogfileConfiguration) != 5 {
		t.Fatal("unexpected number of logfile definitions (5): ", len(c.LogfileConfiguration))
	}
}

func TestReadCertificatesAndFilesConfig(t *testing.T) {
	cfgdir := t.TempDir()
	configPath := filepath.Join(cfgdir, "config.ini")
	certificatesPath := filepath.Join(cfgdir, "certificates.ini")
	filesPath := filepath.Join(cfgdir, "files.ini")

	config := fmt.Sprintf("[default]\ncertificates = %s\nfiles = %s\n", certificatesPath, filesPath)
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatal(err)
	}

	certificates := `[default]
[letsencrypt]
  path = /etc/letsencrypt/live/*/fullchain.pem
  hostname = www.example.com
  enabled = true
[iis]
  type = windows_store
  subject = (?i)example\.com
  enabled = true
[disabled]
  path = /tmp/foo.pem
  enabled = false
[missing_path]
  enabled = true
[invalid_subject]
  type = windows_store
  subject = (unclosed
  enabled = true
`
	files := `[backup]
  path = /backup/db-*.sql.gz
  enabled = true
[mailqueue]
  type = directory
  path = /var/spool/postfix/deferred
  recursive = true
  pattern = *
  enabled = true
[unknown]
  type = socket
  path = /tmp
  enabled = true
`
	if err := os.WriteFile(certificatesPath, []byte(certificates), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filesPath, []byte(files), 0600); err != nil {
		t.Fatal(err)
	}

	c, err := Load(context.Background(), configPath)
	if err != nil {
		t.Fatal(err)
	}

	certificateDefinitions := map[string]*CertificateDefinition{}
	for _, d := range c.CertificateConfiguration {
		certificateDefinitions[d.Name] = d
	}
	if len(certificateDefinitions) != 4 {
		t.Fatal("expected 4 enabled certificate definitions, got: ", len(certificateDefinitions))
	}
	if d := certificateDefinitions["letsencrypt"]; d.Type != CertificateTypeFile || d.Hostname != "www.example.com" || d.ConfigError != "" {
		t.Errorf("unexpected definition: %+v", d)
	}
	if d := certificateDefinitions["iis"]; d.Store != `LocalMachine\My` || d.SubjectRegexp == nil || !d.SubjectRegexp.MatchString("CN=WWW.EXAMPLE.COM") {
		t.Errorf("unexpected definition: %+v", d)
	}
	if !strings.Contains(certificateDefinitions["missing_path"].ConfigError, "missing path") {
		t.Error("expected missing path error")
	}
	if !strings.Contains(certificateDefinitions["invalid_subject"].ConfigError, "invalid subject regex") {
		t.Error("expected invalid subject regex error")
	}

	fileDefinitions := map[string]*FileDefinition{}
	for _, d := range c.FileConfiguration {
		fileDefinitions[d.Name] = d
	}
	if len(fileDefinitions) != 3 {
		t.Fatal("expected 3 file definitions, got: ", len(fileDefinitions))
	}
	if d := fileDefinitions["backup"]; d.Type != FileTypeFile || d.MaxFiles != 100000 || d.ConfigError != "" {
		t.Errorf("unexpected definition: %+v", d)
	}
	if d := fileDefinitions["mailqueue"]; d.Type != FileTypeDirectory || !d.Recursive || d.Pattern != "*" {
		t.Errorf("unexpected definition: %+v", d)
	}
	if !strings.Contains(fileDefinitions["unknown"].ConfigError, "unknown type") {
		t.Error("expected unknown type error")
	}

	if string(c.ReadCertificateConfiguration()) != certificates || string(c.ReadFileConfiguration()) != files {
		t.Error("unexpected content of the configuration files")
	}
}

func TestReadHTTPConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "http.ini")
	content := `[intranet]
  url = https://intranet.example.com/health
  match = (?i)healthy
  enabled = true
[api]
  url = http://127.0.0.1:8080/api/status
  method = head
  expected_status = 200-299, 301
  follow_redirects = false
  timeout = 3
  enabled = true
[invalid_url]
  url = ftp://example.com
  enabled = true
[invalid_status]
  url = http://example.com
  expected_status = 500-200
  enabled = true
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	definitions, err := unmarshalHTTP(path)
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]*HTTPDefinition{}
	for _, d := range definitions {
		byName[d.Name] = d
	}

	if d := byName["intranet"]; d.ConfigError != "" || d.Method != "GET" || d.Timeout != 10 || !*d.FollowRedirects || d.ExpectedStatus != "200-399" || d.MatchRegexp == nil {
		t.Errorf("unexpected defaults: %+v", d)
	}
	if d := byName["api"]; d.ConfigError != "" || d.Method != "HEAD" || d.Timeout != 3 || *d.FollowRedirects || fmt.Sprint(d.ExpectedStatusRanges) != "[[200 299] [301 301]]" {
		t.Errorf("unexpected definition: %+v", d)
	}
	if !strings.Contains(byName["invalid_url"].ConfigError, "invalid url") {
		t.Error("expected invalid url error")
	}
	if !strings.Contains(byName["invalid_status"].ConfigError, "invalid expected_status") {
		t.Error("expected invalid status error")
	}
}

func TestReadRedfishConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "redfish.ini")
	content := "[server01]\n" +
		"  url = https://10.0.0.10/redfish/v1/\n" +
		"  username = monitoring\n" +
		"  password = `pass#word;1`\n" +
		"  skip_verify = true\n" +
		"  enabled = true\n" +
		"[server02]\n" +
		"  url = https://bmc02.example.com:8443\n" +
		"  fingerprint = AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89:AB:CD:EF:01:23:45:67:89\n" +
		"  interval = 10\n" +
		"  timeout = 5\n" +
		"  eventlog = false\n" +
		"  eventlog_window = 3600\n" +
		"  eventlog_services = ^Sel$\n" +
		"  enabled = true\n" +
		"[plain_http]\n" +
		"  url = http://10.0.0.11\n" +
		"  enabled = true\n" +
		"[plain_http_allowed]\n" +
		"  url = http://10.0.0.12\n" +
		"  allow_http = true\n" +
		"  enabled = true\n" +
		"[invalid_fingerprint]\n" +
		"  url = https://10.0.0.13\n" +
		"  fingerprint = abc\n" +
		"  enabled = true\n" +
		"[disabled]\n" +
		"  url = https://10.0.0.14\n" +
		"  enabled = false\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	definitions, err := unmarshalRedfish(path)
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]*RedfishDefinition{}
	for _, d := range definitions {
		byName[d.Name] = d
	}
	if len(definitions) != 5 {
		t.Fatalf("expected 5 enabled definitions, got %d", len(definitions))
	}

	d := byName["server01"]
	if d.ConfigError != "" || d.URL != "https://10.0.0.10" || d.Password != "pass#word;1" || !d.SkipVerify {
		t.Errorf("unexpected definition: %+v", d)
	}
	if d.Interval != 300 || d.Timeout != 30 || !*d.Eventlog || d.EventlogWindow != 86400 || !d.EventlogServicesRegexp.MatchString("Sel") || d.EventlogServicesRegexp.MatchString("Lclog") {
		t.Errorf("unexpected defaults: %+v", d)
	}

	d = byName["server02"]
	if d.ConfigError != "" || d.URL != "https://bmc02.example.com:8443" || len(d.FingerprintBytes) != 32 || d.Interval != 60 || d.Timeout != 5 || *d.Eventlog || d.EventlogWindow != 3600 || d.EventlogServicesRegexp.MatchString("IML") {
		t.Errorf("unexpected definition: %+v", d)
	}

	if !strings.Contains(byName["plain_http"].ConfigError, "plain text") {
		t.Errorf("expected an error for http://: %s", byName["plain_http"].ConfigError)
	}
	if byName["plain_http_allowed"].ConfigError != "" {
		t.Errorf("unexpected error: %s", byName["plain_http_allowed"].ConfigError)
	}
	if !strings.Contains(byName["invalid_fingerprint"].ConfigError, "invalid fingerprint") {
		t.Errorf("expected an error for the fingerprint: %s", byName["invalid_fingerprint"].ConfigError)
	}
}
