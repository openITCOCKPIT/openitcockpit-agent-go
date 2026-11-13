package checks

import (
	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckCpu gathers information about system CPU load
type CheckCpu struct {
	checkInterval int64
}

type cpuDetails struct {
	User   float64 // Linux, macOS, Windows, AIX - Seconds
	Nice   float64 // Linux, macOS           - Seconds
	System float64 // Linux, macOS, Windows, AIX - Seconds
	Idle   float64 // Linux, macOS, Windows, AIX - Seconds
	Iowait float64 // Linux, AIX
}

type resultCpu struct {
	PercentageTotal   float64      `json:"cpu_total_percentage"`
	PercentagePerCore []float64    `json:"cpu_percentage"`
	DetailsTotal      *cpuDetails  `json:"cpu_total_percentage_detailed"`
	DetailsPerCore    []cpuDetails `json:"cpu_percentage_detailed"`
	// Only on AIX shared processor partitions
	Partition *cpuPartition `json:"partition,omitempty"`
}

// cpuPartition is the processor usage of an AIX shared processor partition (LPAR)
type cpuPartition struct {
	PhysicalCpus          float64 `json:"physc"`                  // Consumed physical processors
	Entitlement           float64 `json:"entitlement"`            // Entitled processor units
	EntitlementPercentage float64 `json:"entitlement_percentage"` // physc in percent of the entitlement
}

// Name will be used in the response as check name
func (c *CheckCpu) Name() string {
	return "cpu"
}

// Configure the command or return false if the command was disabled
func (c *CheckCpu) Configure(config *config.Configuration) (bool, error) {
	c.checkInterval = config.CheckInterval // check interval in seconds (default: 30)
	return config.CPU, nil
}
