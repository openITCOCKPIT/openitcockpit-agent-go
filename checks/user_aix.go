package checks

import (
	"context"
	"time"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckUser) Run(ctx context.Context) (interface{}, error) {
	// who prints the login time in the time zone of the process, see runAixCommandInUTC
	output, err := runAixCommandInUTC(ctx, "who")
	if err != nil {
		return nil, err
	}
	return parseAixWho(output, time.Now()), nil
}
