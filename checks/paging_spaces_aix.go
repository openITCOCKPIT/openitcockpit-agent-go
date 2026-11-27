package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckPagingSpaces reports the paging spaces of AIX (lsps -a)
type CheckPagingSpaces struct {
}

// Name will be used in the response as check name
func (c *CheckPagingSpaces) Name() string {
	return "paging_spaces"
}

// Configure the command or return false if the command was disabled
func (c *CheckPagingSpaces) Configure(config *config.Configuration) (bool, error) {
	return config.PagingSpaces && aixCommandAvailable("lsps"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckPagingSpaces) Run(ctx context.Context) (interface{}, error) {
	output, err := runAixCommand(ctx, "lsps", "-a")
	if err != nil {
		return nil, err
	}
	return parseLsps(output), nil
}
