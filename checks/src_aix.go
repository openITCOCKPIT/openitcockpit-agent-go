package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckSrcServices reports the subsystems of the AIX System Resource Controller (lssrc -a)
type CheckSrcServices struct {
}

// Name will be used in the response as check name
func (c *CheckSrcServices) Name() string {
	return "src_services"
}

// Configure the command or return false if the command was disabled
func (c *CheckSrcServices) Configure(config *config.Configuration) (bool, error) {
	return config.SrcServices && aixCommandAvailable("lssrc"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckSrcServices) Run(ctx context.Context) (interface{}, error) {
	output, err := runAixCommand(ctx, "lssrc", "-a")
	if err != nil {
		return nil, err
	}
	return parseLssrc(output), nil
}
