package checks

import (
	"context"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckErrpt reports the entries of the AIX error log (errpt) within a time window
type CheckErrpt struct {
	window time.Duration
}

// Name will be used in the response as check name
func (c *CheckErrpt) Name() string {
	return "errpt"
}

// Configure the command or return false if the command was disabled
func (c *CheckErrpt) Configure(config *config.Configuration) (bool, error) {
	c.window = time.Duration(config.ErrptWindow) * time.Second
	if c.window <= 0 {
		c.window = 24 * time.Hour
	}
	return config.Errpt && aixCommandAvailable("errpt"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckErrpt) Run(ctx context.Context) (interface{}, error) {
	// errpt prints and expects the time in the time zone of the process, it runs in UTC (see runAixCommandInUTC)
	start := time.Now().Add(-c.window).UTC()
	output, err := runAixCommandInUTC(ctx, "errpt", "-s", start.Format(errptTimeFormat))
	if err != nil {
		return nil, err
	}
	return errptResult(parseErrpt(output, time.UTC), c.window), nil
}
