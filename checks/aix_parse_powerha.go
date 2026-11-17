//go:build aix || linux
// +build aix linux

// Parsers for the commands of the AIX cluster PowerHA SystemMirror. The file is also built on Linux,
// so the unit tests run in the CI.

package checks

import (
	"regexp"
	"strings"
)

// resultPowerHA is the state of the PowerHA SystemMirror cluster of the node
type resultPowerHA struct {
	ClusterName     string `json:"cluster_name"`     // clmgr query cluster
	ClusterState    string `json:"cluster_state"`    // e.g.: STABLE
	Version         string `json:"version"`          // e.g.: 7.2.9.2
	UnsyncedChanges bool   `json:"unsynced_changes"` // The configuration was changed, but not synchronized to all nodes
	// State of the cluster manager clstrmgrES of the node: ST_STABLE, ST_INIT (cluster services are not running),
	// ST_RP_FAILED (an event script failed), ST_UNSTABLE, ST_JOINING, ST_BARRIER, ...
	ClusterManagerState string                         `json:"cluster_manager_state"`
	Nodes               []*resultPowerHANode           `json:"nodes"`
	ResourceGroups      []*resultPowerHAResourceGroup  `json:"resource_groups"`
	RepositoryDisks     []*resultPowerHARepositoryDisk `json:"repository_disks"`
	// Commands which failed, e.g. without root or if the cluster services are not running.
	// The values of a failed command are empty.
	Errors []string `json:"errors"`
}

// resultPowerHANode is a node of the cluster (clshowsrv -v)
type resultPowerHANode struct {
	Name                 string `json:"name"`
	Local                bool   `json:"local"`                 // The node of the agent
	ClusterServices      string `json:"cluster_services"`      // e.g.: NORMAL, INIT
	ClusterManagerState  string `json:"cluster_manager_state"` // e.g.: ST_STABLE, ST_INIT
	RemoteCommunications string `json:"remote_communications"` // UP or DOWN
	CaaState             string `json:"caa_state"`             // State of Cluster-Aware AIX: UP or DOWN
}

// resultPowerHAResourceGroup is a resource group with its state on each node (clRGinfo -s)
type resultPowerHAResourceGroup struct {
	Name           string                            `json:"name"`
	Type           string                            `json:"type"`            // non-concurrent or concurrent
	StartupPolicy  string                            `json:"startup_policy"`  // e.g.: OHN (online on home node only)
	FalloverPolicy string                            `json:"fallover_policy"` // e.g.: FNPN (fallover to next priority node)
	FallbackPolicy string                            `json:"fallback_policy"` // e.g.: NFB (never fallback)
	Nodes          []*resultPowerHAResourceGroupNode `json:"nodes"`           // In the order of the priority, the first node is the home node
	Applications   []*resultPowerHAApplication       `json:"applications"`
}

// resultPowerHAResourceGroupNode is the state of a resource group on a node
type resultPowerHAResourceGroupNode struct {
	Node  string `json:"node"`
	State string `json:"state"` // e.g.: ONLINE, OFFLINE, ERROR, ACQUIRING, RELEASING, UNMANAGED, ONLINE SECONDARY
}

// resultPowerHAApplication is an application controller of a resource group (clRGinfo -m)
type resultPowerHAApplication struct {
	Name  string `json:"name"`
	State string `json:"state"` // e.g.: ONLINE MONITORED, ONLINE NOT MONITORED, MONITOR FAILED, OFFLINE
}

// resultPowerHARepositoryDisk is a repository disk of Cluster-Aware AIX on a node (lscluster -d)
type resultPowerHARepositoryDisk struct {
	Node  string `json:"node"`
	Disk  string `json:"disk"`  // e.g.: hdisk100
	State string `json:"state"` // UP or DOWN
	Type  string `json:"type"`  // REPDISK (active) or BACKUP_DISK
}

// clRGinfoColumns is the minimum number of columns of clRGinfo -s
const clRGinfoColumns = 7

// parseClstrmgrState returns the state of the cluster manager (lssrc -ls clstrmgrES)
//
//	Current state: ST_STABLE
func parseClstrmgrState(output string) string {
	for _, line := range strings.Split(output, "\n") {
		value, found := strings.CutPrefix(strings.TrimSpace(line), "Current state:")
		if found {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// parseClRGinfo parses the output of clRGinfo -s: a line for each resource group and node,
// the nodes of a resource group are in the order of their priority
//
//	RG1:ONLINE:node1:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
//	RG1:OFFLINE:node2:non-concurrent:OHN:FNPN:NFB:ignore::: : ::::
func parseClRGinfo(output string) []*resultPowerHAResourceGroup {
	groups := []*resultPowerHAResourceGroup{}
	byName := map[string]*resultPowerHAResourceGroup{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(strings.TrimSpace(line), ":")
		if len(fields) < clRGinfoColumns || fields[0] == "" {
			continue
		}
		group, ok := byName[fields[0]]
		if !ok {
			group = newPowerHAResourceGroup(fields)
			byName[group.Name] = group
			groups = append(groups, group)
		}
		group.Nodes = append(group.Nodes, &resultPowerHAResourceGroupNode{
			Node:  fields[2],
			State: fields[1],
		})
	}
	return groups
}

func newPowerHAResourceGroup(fields []string) *resultPowerHAResourceGroup {
	return &resultPowerHAResourceGroup{
		Name:           fields[0],
		Type:           fields[3],
		StartupPolicy:  fields[4],
		FalloverPolicy: fields[5],
		FallbackPolicy: fields[6],
		Nodes:          []*resultPowerHAResourceGroupNode{},
		Applications:   []*resultPowerHAApplication{},
	}
}

// parseClRGinfoApplications adds the application controllers of clRGinfo -m to the resource groups.
// The line of a resource group starts in the first column, its applications are indented by one space.
// Other nodes of a resource group are indented further and are ignored (clRGinfo -s contains them).
//
//	Group Name     State            Application state            Node
//	-----------------------------------------------------------------------------
//	RG1            ONLINE                                        node1
//	 APP1                           ONLINE NOT MONITORED
func parseClRGinfoApplications(output string, groups []*resultPowerHAResourceGroup) {
	byName := map[string]*resultPowerHAResourceGroup{}
	for _, group := range groups {
		byName[group.Name] = group
	}

	var current *resultPowerHAResourceGroup
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "---") || fields[0] == "Group" {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			current = byName[fields[0]]
			continue
		}
		if current == nil || !isPowerHAApplicationLine(line) || len(fields) < 2 {
			continue
		}
		current.Applications = append(current.Applications, &resultPowerHAApplication{
			Name:  fields[0],
			State: strings.Join(fields[1:], " "),
		})
	}
}

// isPowerHAApplicationLine returns true for the lines of clRGinfo -m which are indented by one or two spaces
func isPowerHAApplicationLine(line string) bool {
	indent := len(line) - len(strings.TrimLeft(line, " "))
	return indent > 0 && indent <= 2
}

// quotedValue matches the values of clshowsrv -v, e.g. "NORMAL" ("ST_STABLE")
var quotedValue = regexp.MustCompile(`"([^"]*)"`)

// parseClshowsrv parses the state of the nodes of clshowsrv -v. The subsystems after the nodes are ignored.
//
//	Local node: "node1" ("node1.example.com", "node1.example.com")
//	         Cluster services status:   "NORMAL" ("ST_STABLE")
//	         Remote communications:     "UP"
//	         Cluster-Aware AIX status:  "UP"
func parseClshowsrv(output string) []*resultPowerHANode {
	nodes := []*resultPowerHANode{}
	var current *resultPowerHANode
	for _, line := range strings.Split(output, "\n") {
		label, value, found := strings.Cut(strings.TrimSpace(line), ":")
		if !found {
			continue
		}
		values := quotedValues(value)
		if label == "Local node" || label == "Remote node" {
			current = newPowerHANode(label, values)
			nodes = append(nodes, current)
			continue
		}
		if current == nil || len(values) == 0 {
			continue
		}
		setPowerHANodeValue(current, label, values)
	}
	return nodes
}

func newPowerHANode(label string, values []string) *resultPowerHANode {
	node := &resultPowerHANode{
		Local: label == "Local node",
	}
	if len(values) > 0 {
		node.Name = values[0]
	}
	return node
}

func setPowerHANodeValue(node *resultPowerHANode, label string, values []string) {
	switch label {
	case "Cluster services status":
		node.ClusterServices = values[0]
		if len(values) > 1 {
			node.ClusterManagerState = values[1]
		}
	case "Remote communications":
		node.RemoteCommunications = values[0]
	case "Cluster-Aware AIX status":
		node.CaaState = values[0]
	}
}

func quotedValues(value string) []string {
	values := []string{}
	for _, match := range quotedValue.FindAllStringSubmatch(value, -1) {
		values = append(values, match[1])
	}
	return values
}

// parseLsclusterDisks parses the repository disks of each node of lscluster -d
//
//	Node node1.example.com
//	Node UUID = 6d1f...
//	Number of disks discovered = 2
//		hdisk101:
//		      State : UP
//		       Type : BACKUP_DISK
func parseLsclusterDisks(output string) []*resultPowerHARepositoryDisk {
	disks := []*resultPowerHARepositoryDisk{}
	var node string
	var current *resultPowerHARepositoryDisk
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if isLsclusterNodeLine(trimmed) {
			node = strings.TrimSpace(strings.TrimPrefix(trimmed, "Node "))
			current = nil
			continue
		}
		if strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, " ") {
			current = &resultPowerHARepositoryDisk{Node: node, Disk: strings.TrimSuffix(trimmed, ":")}
			disks = append(disks, current)
			continue
		}
		if current == nil {
			continue
		}
		setLsclusterDiskValue(current, trimmed)
	}
	return disks
}

// isLsclusterNodeLine returns true for the line with the name of a node, but not for "Node UUID = ..."
func isLsclusterNodeLine(line string) bool {
	return strings.HasPrefix(line, "Node ") && !strings.Contains(line, "=")
}

func setLsclusterDiskValue(disk *resultPowerHARepositoryDisk, line string) {
	label, value, found := strings.Cut(line, ":")
	if !found {
		return
	}
	switch strings.TrimSpace(label) {
	case "State":
		disk.State = strings.TrimSpace(value)
	case "Type":
		disk.Type = strings.TrimSpace(value)
	}
}

// parseClmgrQuery parses the attributes of clmgr query (KEY="VALUE")
//
//	CLUSTER_NAME="cluster1"
//	UNSYNCED_CHANGES="false"
func parseClmgrQuery(output string) map[string]string {
	attributes := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if !found || key == "" || strings.Contains(key, " ") {
			continue
		}
		attributes[key] = strings.Trim(value, `"`)
	}
	return attributes
}

// setPowerHAClusterAttributes sets the values of clmgr query cluster
func setPowerHAClusterAttributes(result *resultPowerHA, attributes map[string]string) {
	result.ClusterName = attributes["CLUSTER_NAME"]
	result.ClusterState = attributes["STATE"]
	result.Version = attributes["VERSION"]
	result.UnsyncedChanges = attributes["UNSYNCED_CHANGES"] == "true"
}
