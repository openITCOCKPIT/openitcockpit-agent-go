package config

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/basiclog"
)

const (
	redfishDefaultInterval       = 300
	redfishMinInterval           = 60
	redfishDefaultTimeout        = 30
	redfishDefaultEventlogWindow = 86400
	// Event logs which contain hardware events: SEL (Dell, Supermicro, generic), IML (HPE), Log1 (DMTF mockups),
	// PlatformLog/ActiveLog (Lenovo). Huge logs like the Dell Lifecycle log or audit logs are skipped by default.
	RedfishDefaultEventlogServices = `(?i)^(sel|iml|log[0-9]*|eventlog|platformlog|activelog)$`
)

// RedfishDefinition describes a BMC (e.g. Dell iDRAC, HPE iLO, Lenovo XClarity, Supermicro) which is monitored with Redfish
type RedfishDefinition struct {
	Name    string `mapstructure:"-"`
	Enabled bool   `mapstructure:"enabled"`
	// URL of the BMC, e.g. https://10.0.0.10
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	// Do not verify the TLS certificate (most BMCs use self signed certificates)
	SkipVerify bool `mapstructure:"skip_verify"`
	// SHA-256 fingerprint of the certificate of the BMC. The certificate is accepted if the fingerprint matches.
	Fingerprint string `mapstructure:"fingerprint"`
	// http:// URLs are only allowed if this is true, as the password would be sent in plain text
	AllowHTTP bool `mapstructure:"allow_http"`
	// Seconds between two collections. Default: 300, minimum: 60
	Interval int64 `mapstructure:"interval"`
	// Timeout of a single request in seconds. Default: 30
	Timeout int64 `mapstructure:"timeout"`
	// Read the event logs of the BMC. Default: true
	Eventlog *bool `mapstructure:"eventlog"`
	// Report event log entries of the last N seconds. Default: 86400
	EventlogWindow int64 `mapstructure:"eventlog_window"`
	// Regular expression for the Id of the log services which are read. Default: RedfishDefaultEventlogServices
	EventlogServices string `mapstructure:"eventlog_services"`
	// HTTP proxy, e.g. http://proxy.example.com:3128. By default no proxy is used
	Proxy string `mapstructure:"proxy"`

	FingerprintBytes       []byte         `json:"-" mapstructure:"-"`
	EventlogServicesRegexp *regexp.Regexp `json:"-" mapstructure:"-"`
	// ConfigError is set if the definition is invalid. The check will report it as error
	ConfigError string `mapstructure:"-"`
}

func unmarshalRedfish(configPath string) ([]*RedfishDefinition, error) {
	cfg, err := readIniSections[RedfishDefinition](configPath)
	if err != nil {
		return nil, err
	}
	for _, file := range definitionFiles(configPath) {
		warnIfReadableByOthers(file)
	}

	definitions := make([]*RedfishDefinition, 0)
	for name, definition := range cfg {
		if !definition.Enabled {
			continue
		}
		definition.Name = name

		// Invalid definitions are kept, so the check can report the error for the affected definition only
		if err := definition.validate(); err != nil {
			definition.ConfigError = err.Error()
			logInvalidDefinition("redfish definition", name, err)
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

// warnIfReadableByOthers logs a warning if a file with credentials can be read by other users
func warnIfReadableByOthers(configPath string) {
	if runtime.GOOS == "windows" {
		return
	}
	if info, err := os.Stat(configPath); err == nil && info.Mode().Perm()&0077 != 0 {
		logger, _ := basiclog.New()
		logger.Errorln("Configuration: ", configPath, " contains passwords and should only be readable by the agent (chmod 600)")
	}
}

func (d *RedfishDefinition) validate() error {
	d.URL = strings.TrimRight(strings.TrimSpace(d.URL), "/")
	parsed, err := url.Parse(d.URL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return fmt.Errorf("invalid url '%s' in redfish definition %s, expected https://", d.URL, d.Name)
	}
	if parsed.Scheme == "http" && !d.AllowHTTP {
		return fmt.Errorf("http:// would send the password in plain text, use https:// or set allow_http = true in redfish definition %s", d.Name)
	}
	// Only scheme and host are used, the service root is always /redfish/v1
	d.URL = parsed.Scheme + "://" + parsed.Host

	if d.Interval <= 0 {
		d.Interval = redfishDefaultInterval
	} else if d.Interval < redfishMinInterval {
		d.Interval = redfishMinInterval
	}
	if d.Timeout <= 0 {
		d.Timeout = redfishDefaultTimeout
	}
	if d.Eventlog == nil {
		eventlog := true
		d.Eventlog = &eventlog
	}
	if d.EventlogWindow <= 0 {
		d.EventlogWindow = redfishDefaultEventlogWindow
	}
	if strings.TrimSpace(d.EventlogServices) == "" {
		d.EventlogServices = RedfishDefaultEventlogServices
	}
	if d.EventlogServicesRegexp, err = regexp.Compile(d.EventlogServices); err != nil {
		return fmt.Errorf("invalid eventlog_services regex in redfish definition %s: %s", d.Name, err)
	}

	if d.Fingerprint != "" {
		// Accept the formats AA:BB:..., aa bb ... and aabb...
		fingerprint := strings.NewReplacer(":", "", " ", "", "-", "").Replace(strings.TrimSpace(d.Fingerprint))
		if d.FingerprintBytes, err = hex.DecodeString(fingerprint); err != nil || len(d.FingerprintBytes) != 32 {
			return fmt.Errorf("invalid fingerprint in redfish definition %s, expected the SHA-256 fingerprint of the certificate", d.Name)
		}
	}

	if d.Proxy != "" {
		if _, err := url.Parse(d.Proxy); err != nil {
			return fmt.Errorf("invalid proxy in redfish definition %s: %s", d.Name, err)
		}
	}
	return nil
}

func (c *Configuration) ReadRedfishConfiguration() []byte {
	return readOptionalFile(c.RedfishFilePath)
}

func (c *Configuration) SaveRedfishConfiguration(config []byte) error {
	return saveOptionalFile(c.RedfishFilePath, config)
}
