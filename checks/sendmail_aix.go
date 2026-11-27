package checks

import (
	"context"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckSendmail reports the mail queue of sendmail, the mail system of AIX (mailq, requires root)
type CheckSendmail struct {
}

// Name will be used in the response as check name
func (c *CheckSendmail) Name() string {
	return "sendmail"
}

// Configure the command or return false if the command was disabled
func (c *CheckSendmail) Configure(config *config.Configuration) (bool, error) {
	return config.Sendmail && aixCommandAvailable("mailq"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckSendmail) Run(ctx context.Context) (interface{}, error) {
	// mailq prints the queue times in local time, see runAixCommandInUTC
	output, err := runAixCommandInUTC(ctx, "mailq")
	if err != nil {
		return nil, err
	}
	return parseMailq(output, time.Now().UTC()), nil
}
