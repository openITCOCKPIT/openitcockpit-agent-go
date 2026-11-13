package checks

import (
	"context"
	"time"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckProcess) Run(ctx context.Context) (interface{}, error) {
	// ps -e lists all processes except kernel processes. %CPU is the average since the start of the process.
	output, err := runAixCommand(ctx, "ps", aixPsArguments...)
	if err != nil {
		return nil, err
	}
	return parseAixPs(output, time.Now()), nil
}
