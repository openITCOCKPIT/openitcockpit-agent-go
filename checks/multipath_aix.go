package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckMultipath reports the paths of the disks of AIX MPIO (lspath)
type CheckMultipath struct {
}

// Name will be used in the response as check name
func (c *CheckMultipath) Name() string {
	return "multipath"
}

// Configure the command or return false if the command was disabled
func (c *CheckMultipath) Configure(config *config.Configuration) (bool, error) {
	return config.Multipath && aixCommandAvailable("lspath"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckMultipath) Run(ctx context.Context) (interface{}, error) {
	output, err := runAixCommand(ctx, "lspath")
	if err != nil {
		return nil, err
	}
	return parseLspath(output), nil
}
