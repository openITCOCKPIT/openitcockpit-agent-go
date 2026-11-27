//go:build aix || linux
// +build aix linux

// Parsers for the network adapters of AIX (lsdev, entstat -d). The file is also built on Linux,
// so the unit tests run in the CI.

package checks

import (
	"regexp"
	"strconv"
	"strings"
)

// resultNetworkAdapter is an Ethernet adapter of AIX (lsdev -Cc adapter, entstat -d)
type resultNetworkAdapter struct {
	Name        string   `json:"name"`         // e.g.: ent4
	Description string   `json:"description"`  // lsdev, e.g.: EtherChannel / IEEE 802.3ad Link Aggregation
	State       string   `json:"state"`        // lsdev: Available, Defined
	DeviceType  string   `json:"device_type"`  // entstat, e.g.: EtherChannel, Virtual I/O Ethernet Adapter (l-lan)
	DriverFlags []string `json:"driver_flags"` // e.g.: ["Up", "Broadcast", "Running", ...]
	Running     bool     `json:"running"`      // Driver flag Running: the adapter is operational
	DataRate    uint64   `json:"data_rate"`    // Adapter Data Rate in Mbit/s

	LinkStatus        string `json:"link_status,omitempty"`        // Physical adapters, e.g.: Up, Down
	LanState          string `json:"lan_state,omitempty"`          // Virtual I/O Ethernet adapters, e.g.: LAN_OPERATIONAL
	AggregationStatus string `json:"aggregation_status,omitempty"` // IEEE 802.3ad, e.g.: Aggregated
	LacpActorSync     string `json:"lacp_actor_sync,omitempty"`    // IEEE 802.3ad member: Synchronization of the actor, e.g.: IN_SYNC
	LacpPartnerSync   string `json:"lacp_partner_sync,omitempty"`  // IEEE 802.3ad member: Synchronization of the partner (switch)

	ResetCount     uint64 `json:"reset_count"`     // Adapter Reset Count, e.g. failovers of an EtherChannel
	TransmitErrors uint64 `json:"transmit_errors"` // Counters since the adapter was configured
	ReceiveErrors  uint64 `json:"receive_errors"`

	EtherChannel *resultEtherChannel `json:"etherchannel,omitempty"` // Only for EtherChannels and Network Interface Backup
	Error        string              `json:"error,omitempty"`        // entstat failed
}

// resultEtherChannel is an EtherChannel, IEEE 802.3ad Link Aggregation or Network Interface Backup (NIB)
type resultEtherChannel struct {
	Mode           string                  `json:"mode"`            // Operating mode, e.g.: Standard mode, Round-robin mode, IEEE 802.3ad
	ActiveChannel  string                  `json:"active_channel"`  // primary or backup, empty if not reported
	Adapters       int64                   `json:"adapters"`        // Number of primary adapters
	BackupAdapters int64                   `json:"backup_adapters"` // Number of backup adapters
	Primary        []*resultNetworkAdapter `json:"primary"`
	Backup         []*resultNetworkAdapter `json:"backup"`
}

// aixDevice is a device of lsdev
type aixDevice struct {
	Name        string
	State       string
	Description string
}

// parseLsdev parses the output of lsdev -C (name, state, description):
//
//	ent0    Available  Virtual I/O Ethernet Adapter (l-lan)
//	ent4    Available  EtherChannel / IEEE 802.3ad Link Aggregation
func parseLsdev(output string) []aixDevice {
	devices := []aixDevice{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		devices = append(devices, aixDevice{
			Name:        fields[0],
			State:       fields[1],
			Description: strings.Join(fields[2:], " "),
		})
	}
	return devices
}

// isEthernetAdapter returns true for the Ethernet adapters (entX), entstat only supports them
func isEthernetAdapter(device aixDevice) bool {
	return strings.HasPrefix(device.Name, "ent")
}

// entstatOutput is the output of entstat -d of an adapter
type entstatOutput struct {
	Output string
	Err    error
}

// newNetworkAdapters combines the adapters of lsdev with the outputs of entstat -d. The members of an EtherChannel
// are reported within the EtherChannel only (entstat can not open them, errno 19).
func newNetworkAdapters(devices []aixDevice, outputs map[string]entstatOutput) []*resultNetworkAdapter {
	devicesByName := map[string]aixDevice{}
	for _, device := range devices {
		devicesByName[device.Name] = device
	}

	parsed := map[string]*resultNetworkAdapter{}
	members := map[string]bool{}
	for name, output := range outputs {
		if output.Err != nil {
			continue
		}
		adapter := parseEntstat(output.Output)
		if adapter == nil {
			continue
		}
		parsed[name] = adapter
		addEtherChannelMembers(adapter, devicesByName, members)
	}

	adapters := []*resultNetworkAdapter{}
	for _, device := range devices {
		if !isEthernetAdapter(device) || members[device.Name] {
			continue
		}
		adapter := newNetworkAdapter(device, parsed[device.Name], outputs)
		adapters = append(adapters, adapter)
	}
	return adapters
}

// addEtherChannelMembers sets the lsdev values of the members and remembers their names
func addEtherChannelMembers(adapter *resultNetworkAdapter, devicesByName map[string]aixDevice, members map[string]bool) {
	if adapter.EtherChannel == nil {
		return
	}
	all := append(append([]*resultNetworkAdapter{}, adapter.EtherChannel.Primary...), adapter.EtherChannel.Backup...)
	for _, member := range all {
		members[member.Name] = true
		device, ok := devicesByName[member.Name]
		if ok {
			member.State = device.State
			member.Description = device.Description
		}
	}
}

// newNetworkAdapter returns the parsed adapter or an adapter with the error of entstat
func newNetworkAdapter(device aixDevice, parsed *resultNetworkAdapter, outputs map[string]entstatOutput) *resultNetworkAdapter {
	adapter := parsed
	if adapter == nil {
		adapter = &resultNetworkAdapter{Name: device.Name, DriverFlags: []string{}}
		output, ok := outputs[device.Name]
		if ok && output.Err != nil {
			adapter.Error = output.Err.Error()
		}
	}
	adapter.State = device.State
	adapter.Description = device.Description
	return adapter
}

var (
	entstatAdapterRegex        = regexp.MustCompile(`^ETHERNET STATISTICS \((\S+)\)\s*:`)
	entstatLanStateRegex       = regexp.MustCompile(`^LAN state\(\d+\)\s*:\s*(\S+)`)
	entstatTransmitErrorsRegex = regexp.MustCompile(`Transmit Errors:\s*(\d+)`)
	entstatReceiveErrorsRegex  = regexp.MustCompile(`Receive Errors:\s*(\d+)`)
)

// entstatParser reads the output of entstat -d. The output starts with the adapter itself, an EtherChannel
// is followed by the statistics of its members ("PRIMARY ADAPTERS", "BACKUP ADAPTERS").
type entstatParser struct {
	adapter       *resultNetworkAdapter // The adapter of the entstat command
	current       *resultNetworkAdapter // The adapter (or member) of the current line
	backupSection bool                  // The following members are backup adapters
	driverFlags   bool                  // The current line can continue the driver flags
	lacpState     string                // actor or partner: the current line belongs to the LACP state of the actor or partner
}

// parseEntstat parses the output of entstat -d. It returns nil if the output does not contain an adapter.
//
//	ETHERNET STATISTICS (ent4) :
//	Device Type: EtherChannel
//	...
//	Driver Flags: Up Broadcast Running
//		Simplex 64BitSupport ChecksumOffload
//	Number of adapters: 1
//	Number of backup adapters: 1
//	Active channel: primary channel
//	Operating mode: Standard mode
//	                     PRIMARY ADAPTERS
//	ETHERNET STATISTICS (ent0) :
//	...
//	LAN state(14): LAN_OPERATIONAL
//	                      BACKUP ADAPTERS
//	ETHERNET STATISTICS (ent1) :
//
// The format of IEEE 802.3ad (Aggregation status, LACP synchronization) and the link status of physical adapters
// is taken from the AIX documentation.
func parseEntstat(output string) *resultNetworkAdapter {
	parser := &entstatParser{}
	for _, line := range strings.Split(output, "\n") {
		parser.parseLine(line)
	}
	return parser.adapter
}

func (p *entstatParser) parseLine(line string) {
	match := entstatAdapterRegex.FindStringSubmatch(line)
	if match != nil {
		p.startAdapter(match[1])
		return
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "PRIMARY ADAPTERS" {
		p.backupSection = false
		return
	}
	if trimmed == "BACKUP ADAPTERS" {
		p.backupSection = true
		return
	}
	if p.current == nil {
		return
	}

	if p.driverFlags && strings.HasPrefix(line, "\t") {
		p.addDriverFlags(trimmed)
		return
	}
	p.driverFlags = false
	p.parseValue(line, trimmed)
}

// startAdapter starts the adapter of the command or a member of the EtherChannel
func (p *entstatParser) startAdapter(name string) {
	adapter := &resultNetworkAdapter{Name: name, DriverFlags: []string{}}
	p.current = adapter
	p.driverFlags = false
	p.lacpState = ""
	if p.adapter == nil {
		p.adapter = adapter
		return
	}

	etherChannel := p.etherChannel()
	if p.backupSection {
		etherChannel.Backup = append(etherChannel.Backup, adapter)
	} else {
		etherChannel.Primary = append(etherChannel.Primary, adapter)
	}
}

func (p *entstatParser) etherChannel() *resultEtherChannel {
	if p.adapter.EtherChannel == nil {
		p.adapter.EtherChannel = &resultEtherChannel{
			Primary: []*resultNetworkAdapter{},
			Backup:  []*resultNetworkAdapter{},
		}
	}
	return p.adapter.EtherChannel
}

func (p *entstatParser) addDriverFlags(flags string) {
	p.current.DriverFlags = append(p.current.DriverFlags, strings.Fields(flags)...)
	for _, flag := range p.current.DriverFlags {
		if flag == "Running" {
			p.current.Running = true
		}
	}
}

// parseValue reads the values of a line, e.g. "Adapter Reset Count: 17"
func (p *entstatParser) parseValue(line string, trimmed string) {
	adapter := p.current

	if value, ok := entstatValue(trimmed, "Device Type"); ok {
		adapter.DeviceType = value
	} else if value, ok := entstatValue(trimmed, "Driver Flags"); ok {
		p.addDriverFlags(value)
		p.driverFlags = true
	} else if value, ok := entstatValue(trimmed, "Adapter Reset Count"); ok {
		adapter.ResetCount = parseEntstatUint(value)
	} else if value, ok := entstatValue(trimmed, "Adapter Data Rate"); ok {
		adapter.DataRate = parseEntstatUint(value)
	} else if strings.HasPrefix(trimmed, "Transmit Errors:") {
		adapter.TransmitErrors = parseEntstatCounter(entstatTransmitErrorsRegex, trimmed)
		adapter.ReceiveErrors = parseEntstatCounter(entstatReceiveErrorsRegex, trimmed)
	} else if match := entstatLanStateRegex.FindStringSubmatch(trimmed); match != nil {
		adapter.LanState = match[1]
	} else {
		p.parseLinkValue(trimmed)
	}
}

// parseLinkValue reads the link status and the values of an EtherChannel and IEEE 802.3ad
func (p *entstatParser) parseLinkValue(trimmed string) {
	adapter := p.current

	if value, ok := entstatValue(trimmed, "Physical Port Link Status"); ok {
		adapter.LinkStatus = value
	} else if value, ok := entstatValue(trimmed, "Link Status"); ok {
		if adapter.LinkStatus == "" {
			adapter.LinkStatus = value
		}
	} else if value, ok := entstatValue(trimmed, "Aggregation status"); ok {
		adapter.AggregationStatus = value
	} else if trimmed == "Actor State:" {
		p.lacpState = "actor"
	} else if trimmed == "Partner State:" {
		p.lacpState = "partner"
	} else if value, ok := entstatValue(trimmed, "Synchronization"); ok {
		p.setLacpSync(value)
	} else if adapter == p.adapter {
		p.parseEtherChannelValue(trimmed)
	}
}

func (p *entstatParser) setLacpSync(value string) {
	if p.lacpState == "actor" {
		p.current.LacpActorSync = value
	} else if p.lacpState == "partner" {
		p.current.LacpPartnerSync = value
	}
}

// parseEtherChannelValue reads the summary of an EtherChannel ("Statistics for every adapter in the EtherChannel")
func (p *entstatParser) parseEtherChannelValue(trimmed string) {
	if value, ok := entstatValue(trimmed, "Number of adapters"); ok {
		p.etherChannel().Adapters = parseEntstatInt(value)
	} else if value, ok := entstatValue(trimmed, "Number of backup adapters"); ok {
		p.etherChannel().BackupAdapters = parseEntstatInt(value)
	} else if value, ok := entstatValue(trimmed, "Active channel"); ok {
		p.etherChannel().ActiveChannel = strings.TrimSuffix(value, " channel")
	} else if value, ok := entstatValue(trimmed, "Operating mode"); ok {
		p.etherChannel().Mode = value
	}
}

// entstatValue returns the value of a line "key: value" or "key : value"
func entstatValue(line string, key string) (string, bool) {
	if !strings.HasPrefix(line, key) {
		return "", false
	}
	rest := strings.TrimLeft(line[len(key):], " ")
	if !strings.HasPrefix(rest, ":") {
		return "", false
	}
	return strings.TrimSpace(rest[1:]), true
}

func parseEntstatCounter(regex *regexp.Regexp, line string) uint64 {
	match := regex.FindStringSubmatch(line)
	if match == nil {
		return 0
	}
	return parseEntstatUint(match[1])
}

func parseEntstatUint(value string) uint64 {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return 0
	}
	number, _ := strconv.ParseUint(fields[0], 10, 64)
	return number
}

func parseEntstatInt(value string) int64 {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return 0
	}
	number, _ := strconv.ParseInt(fields[0], 10, 64)
	return number
}
