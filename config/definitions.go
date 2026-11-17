package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/basiclog"
)

// Supported source types of a LogfileDefinition
const (
	LogfileTypeFile       = "file"
	LogfileTypeJournald   = "journald"
	LogfileTypeUnifiedLog = "unifiedlog"
)

// LogfileDefinition describes a log source which gets monitored for lines matching a regular expression
type LogfileDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	// file (all platforms), journald (Linux) or unifiedlog (macOS)
	Type string `mapstructure:"type"`
	// file: Path or glob pattern of the log file(s)
	Path string `mapstructure:"path"`
	// journald: Comma separated list of systemd units
	Unit string `mapstructure:"unit"`
	// journald: Priority or priority range like "err" or "0..4" (see journalctl -p)
	Priority string `mapstructure:"priority"`
	// unifiedlog: Predicate passed to "log show --predicate"
	Predicate string `mapstructure:"predicate"`
	// Regular expression a line has to match. If empty, every line matches
	Match string `mapstructure:"match"`
	// Regular expression to ignore lines which would match otherwise
	Exclude string `mapstructure:"exclude"`
	// Matches older than Window seconds will be discarded
	Window int64 `mapstructure:"window"`
	// Number of the latest matching lines to include in the check result
	MaxLines int64 `mapstructure:"max_lines"`
	// Maximum amount of bytes to read per run and source
	MaxReadBytes int64 `mapstructure:"max_read_bytes"`

	MatchRegexp   *regexp.Regexp `json:"-" mapstructure:"-"`
	ExcludeRegexp *regexp.Regexp `json:"-" mapstructure:"-"`
	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

// Supported types of a CertificateDefinition
const (
	CertificateTypeFile         = "file"
	CertificateTypeWindowsStore = "windows_store"
	CertificateTypeTLS          = "tls"
)

// Supported types of a FileDefinition
const (
	FileTypeFile      = "file"
	FileTypeDirectory = "directory"
)

// CertificateDefinition describes X.509 certificates which should be monitored
type CertificateDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	// file (all platforms), tls (all platforms) or windows_store (Windows)
	Type string `mapstructure:"type"`
	// file: Path or glob pattern of PEM or DER encoded certificate files
	Path string `mapstructure:"path"`
	// windows_store: Certificate store like LocalMachine\My or CurrentUser\My
	Store string `mapstructure:"store"`
	// tls: Address of a TLS endpoint like mail.example.com:993
	Address string `mapstructure:"address"`
	// tls: Server name for SNI. Default: host of the address
	ServerName string `mapstructure:"servername"`
	// tls: Connection timeout in seconds
	Timeout int64 `mapstructure:"timeout"`
	// Regular expression to filter certificates by subject or DNS names
	Subject string `mapstructure:"subject"`
	// Verify that the certificate is valid for this hostname. tls: Default is the server name
	Hostname string `mapstructure:"hostname"`
	// PEM file with additional trusted root certificates (e.g. an internal CA)
	CAFile string `mapstructure:"ca_file"`
	// Do not verify the certificate chain against the trusted root certificates
	SkipVerify bool `mapstructure:"skip_verify"`
	// file: Every certificate of a file is monitored on its own (e.g. CA bundles).
	// By default the first certificate is the monitored certificate and all other certificates are its chain
	SplitBundle bool `mapstructure:"split_bundle"`
	// windows_store: Monitor all certificates. By default only the newest certificate per subject is monitored,
	// so certificates which got replaced by a renewed certificate do not raise an alert
	IncludeAll bool `mapstructure:"include_all"`

	SubjectRegexp *regexp.Regexp `json:"-" mapstructure:"-"`
	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

// FileDefinition describes a file or directory which should be monitored
type FileDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	// file or directory
	Type string `mapstructure:"type"`
	// file: Path or glob pattern of the file(s). directory: Path of the directory
	Path string `mapstructure:"path"`
	// directory: Only count files matching this glob pattern, e.g. *.tar.gz
	Pattern string `mapstructure:"pattern"`
	// directory: Include files of sub directories
	Recursive bool `mapstructure:"recursive"`
	// directory: Stop after this number of files to protect the agent from huge directories
	MaxFiles int64 `mapstructure:"max_files"`

	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

// HTTPDefinition describes an HTTP(S) endpoint which should be monitored
type HTTPDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
	// GET, HEAD or POST. Default: GET
	Method string `mapstructure:"method"`
	// Timeout of the whole request in seconds. Default: 10
	Timeout int64 `mapstructure:"timeout"`
	// Comma separated list of status codes and ranges, e.g. "200-299,301". Default: 200-399
	ExpectedStatus string `mapstructure:"expected_status"`
	// Regular expression the response body has to match (optional)
	Match string `mapstructure:"match"`
	// Follow redirects. Default: true
	FollowRedirects *bool `mapstructure:"follow_redirects"`
	// Do not verify the TLS certificate
	SkipVerify bool `mapstructure:"skip_verify"`
	// Basic authentication (optional)
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	// Value of the Host header, e.g. to check a virtual host by IP address (optional)
	HostHeader string `mapstructure:"host_header"`
	// HTTP proxy, e.g. http://proxy.example.com:3128. By default no proxy is used
	Proxy string `mapstructure:"proxy"`

	MatchRegexp          *regexp.Regexp `json:"-" mapstructure:"-"`
	ExpectedStatusRanges [][2]int       `json:"-" mapstructure:"-"`
	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

// PerfCounterDefinition describes a Windows performance counter which should be monitored
type PerfCounterDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	// English counter path, e.g. \Processor(_Total)\% Processor Time or \Processor(*)\% Processor Time
	Counter string `mapstructure:"counter"`

	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

// loadOptionalDefinitions loads an optional ini file and its drop-in directory. Missing files are not an error
func loadOptionalDefinitions[T any](configPath, description string, unmarshal func(string) ([]*T, error)) []*T {
	if hasDefinitionFiles(configPath) {
		definitions, err := unmarshal(configPath)
		if err == nil {
			return definitions
		}
		logger, _ := basiclog.New()
		logger.Errorln("Configuration: could not load ", description, ": ", err)
	}
	return []*T{}
}

func logInvalidDefinition(description, name string, err error) {
	logger, _ := basiclog.New()
	logger.Errorln("Configuration: invalid ", description, " ", name, ": ", err)
}

func unmarshalLogfiles(configPath string) ([]*LogfileDefinition, error) {
	cfg, err := readIniSections[LogfileDefinition](configPath)
	if err != nil {
		return nil, err
	}

	definitions := make([]*LogfileDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name
		if definition.Window <= 0 {
			definition.Window = 300
		}
		if definition.MaxLines < 0 {
			definition.MaxLines = 0
		} else if definition.MaxLines == 0 {
			definition.MaxLines = 10
		}
		if definition.MaxReadBytes <= 0 {
			definition.MaxReadBytes = 16 * 1024 * 1024
		}

		// Invalid definitions are kept, so the check can report the error for the affected log only
		if err := definition.validate(); err != nil {
			definition.ConfigError = err.Error()
			logInvalidDefinition("logfile definition", name, err)
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func (d *LogfileDefinition) validate() error {
	d.Type = strings.ToLower(strings.TrimSpace(d.Type))
	if d.Type == "" {
		d.Type = LogfileTypeFile
	}

	switch d.Type {
	case LogfileTypeFile:
		if strings.TrimSpace(d.Path) == "" {
			return fmt.Errorf("missing path in logfile definition: %s", d.Name)
		}
	case LogfileTypeJournald, LogfileTypeUnifiedLog:
	default:
		return fmt.Errorf("unknown type '%s' in logfile definition: %s", d.Type, d.Name)
	}

	var err error
	if d.Match != "" {
		if d.MatchRegexp, err = regexp.Compile(d.Match); err != nil {
			return fmt.Errorf("invalid match regex in logfile definition %s: %s", d.Name, err)
		}
	}
	if d.Exclude != "" {
		if d.ExcludeRegexp, err = regexp.Compile(d.Exclude); err != nil {
			return fmt.Errorf("invalid exclude regex in logfile definition %s: %s", d.Name, err)
		}
	}
	return nil
}

func unmarshalCertificates(configPath string) ([]*CertificateDefinition, error) {
	cfg, err := readIniSections[CertificateDefinition](configPath)
	if err != nil {
		return nil, err
	}

	definitions := make([]*CertificateDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name

		// Invalid definitions are kept, so the check can report the error for the affected definition only
		if err := definition.validate(); err != nil {
			definition.ConfigError = err.Error()
			logInvalidDefinition("certificate definition", name, err)
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (d *CertificateDefinition) validate() error {
	d.Type = strings.ToLower(strings.TrimSpace(d.Type))
	if d.Type == "" {
		d.Type = CertificateTypeFile
	}

	switch d.Type {
	case CertificateTypeFile:
		if strings.TrimSpace(d.Path) == "" {
			return fmt.Errorf("missing path in certificate definition: %s", d.Name)
		}
	case CertificateTypeWindowsStore:
		if strings.TrimSpace(d.Store) == "" {
			d.Store = `LocalMachine\My`
		}
	case CertificateTypeTLS:
		if strings.TrimSpace(d.Address) == "" {
			return fmt.Errorf("missing address in certificate definition: %s", d.Name)
		}
		if _, _, err := net.SplitHostPort(d.Address); err != nil {
			return fmt.Errorf("invalid address '%s' in certificate definition %s, expected host:port: %s", d.Address, d.Name, err)
		}
		if d.Timeout <= 0 {
			d.Timeout = 10
		}
	default:
		return fmt.Errorf("unknown type '%s' in certificate definition: %s", d.Type, d.Name)
	}

	if d.Subject != "" {
		var err error
		if d.SubjectRegexp, err = regexp.Compile(d.Subject); err != nil {
			return fmt.Errorf("invalid subject regex in certificate definition %s: %s", d.Name, err)
		}
	}
	return nil
}

func unmarshalFiles(configPath string) ([]*FileDefinition, error) {
	cfg, err := readIniSections[FileDefinition](configPath)
	if err != nil {
		return nil, err
	}

	definitions := make([]*FileDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name
		if definition.MaxFiles <= 0 {
			definition.MaxFiles = 100000
		}

		// Invalid definitions are kept, so the check can report the error for the affected definition only
		if err := definition.validate(); err != nil {
			definition.ConfigError = err.Error()
			logInvalidDefinition("file definition", name, err)
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (d *FileDefinition) validate() error {
	d.Type = strings.ToLower(strings.TrimSpace(d.Type))
	if d.Type == "" {
		d.Type = FileTypeFile
	}

	if d.Type != FileTypeFile && d.Type != FileTypeDirectory {
		return fmt.Errorf("unknown type '%s' in file definition: %s", d.Type, d.Name)
	}
	if strings.TrimSpace(d.Path) == "" {
		return fmt.Errorf("missing path in file definition: %s", d.Name)
	}
	return nil
}

func unmarshalHTTP(configPath string) ([]*HTTPDefinition, error) {
	cfg, err := readIniSections[HTTPDefinition](configPath)
	if err != nil {
		return nil, err
	}

	definitions := make([]*HTTPDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name

		// Invalid definitions are kept, so the check can report the error for the affected definition only
		if err := definition.validate(); err != nil {
			definition.ConfigError = err.Error()
			logInvalidDefinition("http definition", name, err)
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (d *HTTPDefinition) validate() error {
	d.URL = strings.TrimSpace(d.URL)
	parsed, err := url.Parse(d.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return fmt.Errorf("invalid url '%s' in http definition %s, expected http:// or https://", d.URL, d.Name)
	}

	d.Method = strings.ToUpper(strings.TrimSpace(d.Method))
	if d.Method == "" {
		d.Method = "GET"
	}
	if d.Method != "GET" && d.Method != "HEAD" && d.Method != "POST" {
		return fmt.Errorf("unsupported method '%s' in http definition %s", d.Method, d.Name)
	}

	if d.Timeout <= 0 {
		d.Timeout = 10
	}
	if d.FollowRedirects == nil {
		followRedirects := true
		d.FollowRedirects = &followRedirects
	}

	if strings.TrimSpace(d.ExpectedStatus) == "" {
		d.ExpectedStatus = "200-399"
	}
	if d.ExpectedStatusRanges, err = parseStatusRanges(d.ExpectedStatus); err != nil {
		return fmt.Errorf("invalid expected_status in http definition %s: %s", d.Name, err)
	}

	if d.Match != "" {
		if d.MatchRegexp, err = regexp.Compile(d.Match); err != nil {
			return fmt.Errorf("invalid match regex in http definition %s: %s", d.Name, err)
		}
	}

	if d.Proxy != "" {
		if _, err := url.Parse(d.Proxy); err != nil {
			return fmt.Errorf("invalid proxy in http definition %s: %s", d.Name, err)
		}
	}
	return nil
}

// parseStatusRanges parses a list like "200-299,301" into ranges
func parseStatusRanges(value string) ([][2]int, error) {
	var ranges [][2]int
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		bounds := strings.SplitN(part, "-", 2)
		from, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid status code '%s'", part)
		}
		to := from
		if len(bounds) == 2 {
			if to, err = strconv.Atoi(strings.TrimSpace(bounds[1])); err != nil || to < from {
				return nil, fmt.Errorf("invalid status code range '%s'", part)
			}
		}
		ranges = append(ranges, [2]int{from, to})
	}
	if len(ranges) == 0 {
		return nil, fmt.Errorf("no status code defined")
	}
	return ranges, nil
}

func unmarshalPerfCounters(configPath string) ([]*PerfCounterDefinition, error) {
	cfg, err := readIniSections[PerfCounterDefinition](configPath)
	if err != nil {
		return nil, err
	}

	definitions := make([]*PerfCounterDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name
		definition.Counter = strings.TrimSpace(definition.Counter)
		if !strings.HasPrefix(definition.Counter, `\`) {
			definition.ConfigError = fmt.Sprintf("invalid counter '%s' in perfcounter definition %s, expected e.g. \\Processor(_Total)\\%% Processor Time", definition.Counter, name)
			logInvalidDefinition("perfcounter definition", name, errors.New(definition.ConfigError))
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

// readOptionalFile returns the content of a file or an empty slice if the file does not exist
func readOptionalFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		data = []byte{}
	}
	return data
}

func saveOptionalFile(path string, data []byte) error {
	if path == "" {
		return fmt.Errorf("no path for the configuration file defined")
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Configuration) ReadLogfileConfiguration() []byte {
	return readOptionalFile(c.LogfilesFilePath)
}

func (c *Configuration) SaveLogfileConfiguration(config []byte) error {
	return saveOptionalFile(c.LogfilesFilePath, config)
}

func (c *Configuration) ReadCertificateConfiguration() []byte {
	return readOptionalFile(c.CertificatesFilePath)
}

func (c *Configuration) SaveCertificateConfiguration(config []byte) error {
	return saveOptionalFile(c.CertificatesFilePath, config)
}

func (c *Configuration) ReadHTTPConfiguration() []byte {
	return readOptionalFile(c.HTTPFilePath)
}

func (c *Configuration) SaveHTTPConfiguration(config []byte) error {
	return saveOptionalFile(c.HTTPFilePath, config)
}

func (c *Configuration) ReadPerfCounterConfiguration() []byte {
	return readOptionalFile(c.PerfCountersFilePath)
}

func (c *Configuration) SavePerfCounterConfiguration(config []byte) error {
	return saveOptionalFile(c.PerfCountersFilePath, config)
}

func (c *Configuration) ReadFileConfiguration() []byte {
	return readOptionalFile(c.FilesFilePath)
}

func (c *Configuration) SaveFileConfiguration(config []byte) error {
	return saveOptionalFile(c.FilesFilePath, config)
}
