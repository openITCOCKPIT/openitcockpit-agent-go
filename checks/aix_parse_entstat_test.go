//go:build aix || linux
// +build aix linux

package checks

import (
	"errors"
	"reflect"
	"testing"
)

// Network Interface Backup: ent4 with the virtual adapters ent0 (primary) and ent1 (backup), running on the backup channel
// Output of "entstat -d ent4" (shortened, works without root)
const entstatNibBackupChannel = `-------------------------------------------------------------
ETHERNET STATISTICS (ent4) :
Device Type: EtherChannel
Hardware Address: 02:00:00:00:00:04
Elapsed Time: 10 days 3 hours 2 minutes 1 seconds

Transmit Statistics:                          Receive Statistics:
--------------------                          -------------------
Packets: 1000                                 Packets: 2000
Bytes: 100000                                 Bytes: 200000
Interrupts: 0                                 Interrupts: 1900
Transmit Errors: 3                            Receive Errors: 4
Packets Dropped: 0                            Packets Dropped: 0

General Statistics:
-------------------
No mbuf Errors: 0
Adapter Reset Count: 17
Adapter Data Rate: 20000
Driver Flags: Up Broadcast Running
	Simplex 64BitSupport ChecksumOffload
	LargeSend DataRateSet ETHERCHANNEL

=============================================================
=============================================================

Statistics for every adapter in the EtherChannel:
-------------------------------------------------

Number of adapters: 1
Number of backup adapters: 1
Active channel: backup channel
Operating mode: Standard mode
MAC swap: disabled
LLDP mode: none

=============================================================

                     PRIMARY ADAPTERS

=============================================================

-------------------------------------------------------------

ETHERNET STATISTICS (ent0) :
Device Type: Virtual I/O Ethernet Adapter (l-lan)
Hardware Address: 02:00:00:00:00:00

Transmit Errors: 0                            Receive Errors: 0

General Statistics:
-------------------
Adapter Reset Count: 0
Adapter Data Rate: 20000
Driver Flags: Up Broadcast Running
	Simplex 64BitSupport ChecksumOffload
	ECHAN_ELEM

Virtual I/O Ethernet Adapter (l-lan) Specific Statistics:
---------------------------------------------------------
RQ Length: 4481
Trunk Adapter: False
LAN state(14): LAN_OPERATIONAL

Port VLAN ID:     1
Switch ID: ETHERNET0

=============================================================

                      BACKUP ADAPTERS

=============================================================

-------------------------------------------------------------

ETHERNET STATISTICS (ent1) :
Device Type: Virtual I/O Ethernet Adapter (l-lan)
Hardware Address: 02:00:00:00:00:01

Transmit Errors: 0                            Receive Errors: 1

General Statistics:
-------------------
Adapter Reset Count: 0
Adapter Data Rate: 20000
Driver Flags: Up Broadcast Running
	Simplex 64BitSupport ChecksumOffload
	ECHAN_ELEM

Virtual I/O Ethernet Adapter (l-lan) Specific Statistics:
---------------------------------------------------------
LAN state(14): LAN_OPERATIONAL
`

// IEEE 802.3ad with two physical adapters, the link of ent3 is down (format of the AIX documentation)
// Output of "entstat -d ent6" (shortened)
const entstat8023ad = `ETHERNET STATISTICS (ent6) :
Device Type: IEEE 802.3ad Link Aggregation
Adapter Reset Count: 0
Adapter Data Rate: 2000
Driver Flags: Up Broadcast Running
	Simplex 64BitSupport

Statistics for every adapter in the IEEE 802.3ad Link Aggregation:
------------------------------------------------------------------

Number of adapters: 2
Operating mode: Standard mode
IEEE 802.3ad Link Aggregation Statistics:
Aggregation status: Aggregated
LACPDU Interval: Long

ETHERNET STATISTICS (ent2) :
Device Type: 2-Port 10/100/1000 Base-TX PCI-Express Adapter (14104003)
Driver Flags: Up Broadcast Running
	Simplex 64BitSupport
Link Status : Up
Media Speed Running: 1000 Mbps Full Duplex
IEEE 802.3ad Port Statistics:
-----------------------------
	Actor State:
		LACP activity: Active
		Synchronization: IN_SYNC
		Collecting: Enabled
	Partner State:
		LACP activity: Active
		Synchronization: IN_SYNC

ETHERNET STATISTICS (ent3) :
Device Type: 2-Port 10/100/1000 Base-TX PCI-Express Adapter (14104003)
Driver Flags: Up Broadcast
	Simplex 64BitSupport
Link Status : Down
IEEE 802.3ad Port Statistics:
-----------------------------
	Actor State:
		Synchronization: OUT_OF_SYNC
	Partner State:
		Synchronization: OUT_OF_SYNC
`

func TestParseEntstatNetworkInterfaceBackup(t *testing.T) {
	adapter := parseEntstat(entstatNibBackupChannel)
	if adapter == nil {
		t.Fatal("no adapter")
	}
	if adapter.Name != "ent4" || adapter.DeviceType != "EtherChannel" || !adapter.Running || adapter.ResetCount != 17 ||
		adapter.DataRate != 20000 || adapter.TransmitErrors != 3 || adapter.ReceiveErrors != 4 {
		t.Errorf("unexpected adapter %+v", adapter)
	}
	wantFlags := []string{"Up", "Broadcast", "Running", "Simplex", "64BitSupport", "ChecksumOffload", "LargeSend", "DataRateSet", "ETHERCHANNEL"}
	if !reflect.DeepEqual(adapter.DriverFlags, wantFlags) {
		t.Errorf("unexpected driver flags %v", adapter.DriverFlags)
	}

	channel := adapter.EtherChannel
	if channel == nil {
		t.Fatal("no EtherChannel")
	}
	if channel.Mode != "Standard mode" || channel.ActiveChannel != "backup" || channel.Adapters != 1 || channel.BackupAdapters != 1 {
		t.Errorf("unexpected EtherChannel %+v", channel)
	}
	if len(channel.Primary) != 1 || len(channel.Backup) != 1 {
		t.Fatalf("unexpected members %+v %+v", channel.Primary, channel.Backup)
	}
	primary := channel.Primary[0]
	if primary.Name != "ent0" || primary.DeviceType != "Virtual I/O Ethernet Adapter (l-lan)" || !primary.Running ||
		primary.LanState != "LAN_OPERATIONAL" || primary.EtherChannel != nil || primary.LinkStatus != "" {
		t.Errorf("unexpected primary adapter %+v", primary)
	}
	backup := channel.Backup[0]
	if backup.Name != "ent1" || backup.ReceiveErrors != 1 || backup.LanState != "LAN_OPERATIONAL" {
		t.Errorf("unexpected backup adapter %+v", backup)
	}
}

func TestParseEntstatLinkAggregation(t *testing.T) {
	adapter := parseEntstat(entstat8023ad)
	if adapter == nil || adapter.EtherChannel == nil {
		t.Fatalf("unexpected adapter %+v", adapter)
	}
	if adapter.AggregationStatus != "Aggregated" || adapter.EtherChannel.Adapters != 2 || adapter.EtherChannel.ActiveChannel != "" {
		t.Errorf("unexpected adapter %+v %+v", adapter, adapter.EtherChannel)
	}
	members := adapter.EtherChannel.Primary
	if len(members) != 2 || len(adapter.EtherChannel.Backup) != 0 {
		t.Fatalf("unexpected members %+v", members)
	}
	if members[0].LinkStatus != "Up" || !members[0].Running || members[0].LacpActorSync != "IN_SYNC" || members[0].LacpPartnerSync != "IN_SYNC" {
		t.Errorf("unexpected member %+v", members[0])
	}
	if members[1].LinkStatus != "Down" || members[1].Running || members[1].LacpActorSync != "OUT_OF_SYNC" || members[1].LacpPartnerSync != "OUT_OF_SYNC" {
		t.Errorf("unexpected member %+v", members[1])
	}
}

// Output of "entstat -d ent0" (shortened) and the error of entstat for an adapter without a link to a device
func TestParseEntstatPhysicalPort(t *testing.T) {
	adapter := parseEntstat(`ETHERNET STATISTICS (ent0) :
Device Type: PCIe2 4-Port Adapter (1GbE RJ45) (e4148a1614109304)
Driver Flags: Up Broadcast Running
Physical Port Link Status : Up
Logical Port Link Status : Down
Link Status : Down
`)
	if adapter == nil || adapter.LinkStatus != "Up" || adapter.EtherChannel != nil {
		t.Errorf("unexpected adapter %+v", adapter)
	}
	if parseEntstat("entstat: 0909-003 Unable to connect to device ent2, errno = 19") != nil {
		t.Errorf("expected nil without adapter")
	}
}

// Output of "lsdev -Cc adapter" and "entstat -d" of each Ethernet adapter (works without root)
func TestNewNetworkAdapters(t *testing.T) {
	devices := parseLsdev(`ent0    Available  Virtual I/O Ethernet Adapter (l-lan)
ent1    Available  Virtual I/O Ethernet Adapter (l-lan)
ent2    Available  Virtual I/O Ethernet Adapter (l-lan)
ent3    Defined    Virtual I/O Ethernet Adapter (l-lan)
ent4    Available  EtherChannel / IEEE 802.3ad Link Aggregation
vscsi0  Available  Virtual SCSI Client Adapter
`)
	if len(devices) != 6 || devices[4].Description != "EtherChannel / IEEE 802.3ad Link Aggregation" {
		t.Fatalf("unexpected devices %+v", devices)
	}

	outputs := map[string]entstatOutput{
		"ent0": {Err: errors.New("entstat: 0909-003 Unable to connect to device ent0, errno = 19")},
		"ent1": {Err: errors.New("entstat: 0909-003 Unable to connect to device ent1, errno = 19")},
		"ent2": {Err: errors.New("entstat: timeout")},
		"ent4": {Output: entstatNibBackupChannel},
	}
	adapters := newNetworkAdapters(devices, outputs)
	names := []string{}
	for _, adapter := range adapters {
		names = append(names, adapter.Name)
	}
	if !reflect.DeepEqual(names, []string{"ent2", "ent3", "ent4"}) {
		t.Fatalf("unexpected adapters %v", names)
	}
	if adapters[0].Error != "entstat: timeout" || adapters[0].State != "Available" {
		t.Errorf("unexpected adapter %+v", adapters[0])
	}
	if adapters[1].State != "Defined" || adapters[1].Error != "" {
		t.Errorf("unexpected adapter %+v", adapters[1])
	}
	etherChannel := adapters[2]
	if etherChannel.State != "Available" || etherChannel.Description != "EtherChannel / IEEE 802.3ad Link Aggregation" {
		t.Errorf("unexpected adapter %+v", etherChannel)
	}
	primary := etherChannel.EtherChannel.Primary[0]
	if primary.State != "Available" || primary.Description != "Virtual I/O Ethernet Adapter (l-lan)" {
		t.Errorf("unexpected member %+v", primary)
	}
}
