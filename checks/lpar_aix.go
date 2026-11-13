package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckLpar reports the configuration and the processor usage of the AIX LPAR (logical partition):
// consumed physical processors, entitlement and the capacity the partition can use at most
type CheckLpar struct {
	checkInterval int64
}

// Name will be used in the response as check name
func (c *CheckLpar) Name() string {
	return "lpar"
}

// Configure the command or return false if the command was disabled
func (c *CheckLpar) Configure(config *config.Configuration) (bool, error) {
	c.checkInterval = config.CheckInterval
	return config.Lpar, nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckLpar) Run(ctx context.Context) (interface{}, error) {
	// The samples of the background sampler of the CPU check contain the partition values
	base, current, err := cpuTimesOfWindow(ctx, c.checkInterval)
	if err != nil {
		return nil, err
	}
	return aixLparResult(base, current)
}

// Close stops the background sampler, it is shared with the CPU check
func (c *CheckLpar) Close() {
	cpuSamplerState.close()
}
