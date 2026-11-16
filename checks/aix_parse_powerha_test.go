//go:build aix || linux
// +build aix linux

package checks

import (
	"reflect"
	"testing"
)

// The outputs are synthetic (node1, node2, example.com) with the format of PowerHA SystemMirror 7.2.9 on AIX 7.3

// Output of "lssrc -ls clstrmgrES" (shortened, works without root)
func TestParseClstrmgrState(t *testing.T) {
	output := `Current state: ST_STABLE
sccsid = "@(#)  1234567 43haes/usr/sbin/cluster/hacmprd/main.C, 61aha_r729, 2609D_aha729, Aug 24 2022 10:16 PM"
CLversion: 25
local node vrmf is: 7292
Current DNP values:
NodeId 1  NodeName node1
`
	if state := parseClstrmgrState(output); state != "ST_STABLE" {
		t.Errorf("unexpected state %q", state)
	}
	if state := parseClstrmgrState("Current state: ST_INIT\n"); state != "ST_INIT" {
		t.Errorf("unexpected state %q", state)
	}
	if state := parseClstrmgrState(""); state != "" {
		t.Errorf("unexpected state %q", state)
	}
}

// Output of "/usr/es/sbin/cluster/utilities/clRGinfo -s" (requires root)
func TestParseClRGinfo(t *testing.T) {
	// RG2 is concurrent, RG3 ran a fallover to its second node
	output := `RG1:ONLINE:node1:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
RG1:OFFLINE:node2:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
RG2:ONLINE:node1:concurrent:OAAN:BO:NFB:ignore::: : ::::
RG2:ONLINE:node2:concurrent:OAAN:BO:NFB:ignore::: : ::::
RG3:OFFLINE:node1:non-concurrent:OHN:FNPN:FBHPN:ignore::: : ::::
RG3:ONLINE:node2:non-concurrent:OHN:FNPN:FBHPN:ignore::: : ::::
`
	groups := parseClRGinfo(output)
	if len(groups) != 3 {
		t.Fatalf("unexpected resource groups %+v", groups)
	}
	want := &resultPowerHAResourceGroup{
		Name: "RG1", Type: "non-concurrent", StartupPolicy: "OHN", FalloverPolicy: "FNPN", FallbackPolicy: "NFB",
		Nodes: []*resultPowerHAResourceGroupNode{
			{Node: "node1", State: "ONLINE"},
			{Node: "node2", State: "OFFLINE"},
		},
		Applications: []*resultPowerHAApplication{},
	}
	if !reflect.DeepEqual(groups[0], want) {
		t.Errorf("got %+v, want %+v", groups[0], want)
	}
	if groups[1].Type != "concurrent" || groups[1].Nodes[1].State != "ONLINE" {
		t.Errorf("unexpected concurrent resource group %+v", groups[1])
	}
	if groups[2].Nodes[0].Node != "node1" || groups[2].Nodes[1].State != "ONLINE" {
		t.Errorf("the order of the nodes must be kept %+v", groups[2])
	}
	if len(parseClRGinfo("")) != 0 {
		t.Errorf("expected no resource groups")
	}
}

// Output of "/usr/es/sbin/cluster/utilities/clRGinfo -s" and "/usr/es/sbin/cluster/utilities/clRGinfo -m"
// (requires root)
func TestParseClRGinfoApplications(t *testing.T) {
	groups := parseClRGinfo(`RG1:ONLINE:node1:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
RG1:OFFLINE:node2:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
RG2:ONLINE:node2:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
`)
	output := `-----------------------------------------------------------------------------
Group Name     State            Application state            Node
-----------------------------------------------------------------------------
RG1            ONLINE                                        node1
 APP1                           ONLINE NOT MONITORED
 APP2                           ONLINE MONITORED
                                OFFLINE                      node2
RG2            ONLINE                                        node2
 APP3                           MONITOR FAILED
`
	parseClRGinfoApplications(output, groups)
	want := []*resultPowerHAApplication{
		{Name: "APP1", State: "ONLINE NOT MONITORED"},
		{Name: "APP2", State: "ONLINE MONITORED"},
	}
	if !reflect.DeepEqual(groups[0].Applications, want) {
		t.Errorf("got %+v, want %+v", groups[0].Applications, want)
	}
	if len(groups[1].Applications) != 1 || groups[1].Applications[0].State != "MONITOR FAILED" {
		t.Errorf("unexpected applications %+v", groups[1].Applications)
	}
}

// Output of "/usr/es/sbin/cluster/utilities/clshowsrv -v" (requires root)
func TestParseClshowsrv(t *testing.T) {
	output := `Local node: "node1" ("node1.example.com", "node1.example.com")
         Cluster services status:   "NORMAL" ("ST_STABLE")
         Remote communications:     "UP"
         Cluster-Aware AIX status:  "UP"

Remote node: "node2" ("node2.example.com", "node2.example.com")
         Cluster services status:   "INIT" ("ST_INIT")
         Remote communications:     "UP"
         Cluster-Aware AIX status:  "DOWN"

Status of the RSCT subsystems used by PowerHA SystemMirror:
Subsystem         Group            PID          Status
 cthags           cthags           1000001      active
 ctrmc            rsct             1000002      active

Status of the PowerHA SystemMirror subsystems:
Subsystem         Group            PID          Status
 clstrmgrES       cluster          1000003      active
`
	want := []*resultPowerHANode{
		{Name: "node1", Local: true, ClusterServices: "NORMAL", ClusterManagerState: "ST_STABLE", RemoteCommunications: "UP", CaaState: "UP"},
		{Name: "node2", Local: false, ClusterServices: "INIT", ClusterManagerState: "ST_INIT", RemoteCommunications: "UP", CaaState: "DOWN"},
	}
	nodes := parseClshowsrv(output)
	if !reflect.DeepEqual(nodes, want) {
		for _, node := range nodes {
			t.Logf("%+v", node)
		}
		t.Errorf("unexpected nodes")
	}
}

// Output of "lscluster -d" (works without root, the second disk and node are shortened)
func TestParseLsclusterDisks(t *testing.T) {
	output := `Storage Interface Query

Cluster Name: cluster1
Cluster UUID: 00000000-0000-0000-0000-000000000001
Number of nodes reporting = 2
Number of nodes expected = 2

Node node1.example.com
Node UUID = 00000000-0000-0000-0000-000000000002
Number of disks discovered = 2
	hdisk11:
	      State : UP
	       uDid : 00000000000000000000000000000000000000000000000000000000000000
	       PVID : 0000000000000001
	       uUid : 00000000-0000-0000-0000-000000000003
	  Site uUid : 00000000-0000-0000-0000-000000000004
	       Type : BACKUP_DISK
	hdisk10:
	      State : UP
	       Type : REPDISK

Node node2.example.com
Node UUID = 00000000-0000-0000-0000-000000000005
Number of disks discovered = 1
	hdisk10:
	      State : DOWN
	       Type : REPDISK
`
	want := []*resultPowerHARepositoryDisk{
		{Node: "node1.example.com", Disk: "hdisk11", State: "UP", Type: "BACKUP_DISK"},
		{Node: "node1.example.com", Disk: "hdisk10", State: "UP", Type: "REPDISK"},
		{Node: "node2.example.com", Disk: "hdisk10", State: "DOWN", Type: "REPDISK"},
	}
	disks := parseLsclusterDisks(output)
	if !reflect.DeepEqual(disks, want) {
		for _, disk := range disks {
			t.Logf("%+v", disk)
		}
		t.Errorf("unexpected disks")
	}
}

// Output of "/usr/es/sbin/cluster/utilities/clmgr query cluster" (requires root, shortened)
func TestParseClmgrQuery(t *testing.T) {
	output := `CLUSTER_NAME="cluster1"
CLUSTER_ID="1000000001"
STATE="STABLE"
TYPE="NSC"
REPOSITORIES="hdisk10 (0000000000000001), hdisk11 (0000000000000002)"
VERSION="7.2.9.2"
UNSYNCED_CHANGES="true"
CERTIFICATE=""
`
	result := &resultPowerHA{}
	setPowerHAClusterAttributes(result, parseClmgrQuery(output))
	want := &resultPowerHA{ClusterName: "cluster1", ClusterState: "STABLE", Version: "7.2.9.2", UnsyncedChanges: true}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("got %+v, want %+v", result, want)
	}
	if attributes := parseClmgrQuery(output); attributes["REPOSITORIES"] != "hdisk10 (0000000000000001), hdisk11 (0000000000000002)" || attributes["CERTIFICATE"] != "" {
		t.Errorf("unexpected attributes %+v", attributes)
	}
}
