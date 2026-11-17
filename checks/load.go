package checks

import (
	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckLoad gathers information about system CPU load
type CheckLoad struct {
}

// Name will be used in the response as check name
func (c *CheckLoad) Name() string {
	return "system_load"
}

type resultLoad struct {
	Load1  float64 `json:"0"`
	Load5  float64 `json:"1"`
	Load15 float64 `json:"2"`
}

// Configure the command or return false if the command was disabled
func (c *CheckLoad) Configure(config *config.Configuration) (bool, error) {
	return config.Load, nil
}
