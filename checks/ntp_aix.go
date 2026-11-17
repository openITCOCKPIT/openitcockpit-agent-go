package checks

import (
	"context"
)

// Configuration file of xntpd, the NTP daemon of AIX (subsystem xntpd of the SRC)
var aixNtpConfigFiles = []string{"/etc/ntp.conf"}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckNtp) Run(ctx context.Context) (interface{}, error) {
	// If xntpd is running and synchronized, its offset is used (ntpq -c rv, NTP 4.2.8 on AIX 7.2 and newer)
	if result, ok := ntpdOffset(ctx); ok {
		return result, nil
	}

	// Otherwise the agent measures the offset to the configured time server with an SNTP query
	servers := ntpServersFromConfig(aixNtpConfigFiles, "pool.ntp.org")
	return runSNTPCheck(ctx, servers), nil
}
