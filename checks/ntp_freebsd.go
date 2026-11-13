package checks

import (
	"context"
)

// Configuration files of the NTP implementations available on FreeBSD.
// /etc/ntp.conf is part of the base system and always exists, so the configuration files of the packages come first
var freebsdNtpConfigFiles = []string{
	"/usr/local/etc/chrony.conf", // net/chrony
	"/usr/local/etc/ntpd.conf",   // net/openntpd
	"/etc/ntp.conf",              // ntpd of the base system
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckNtp) Run(ctx context.Context) (interface{}, error) {
	// If ntpd is running and synchronized, its offset is used
	if result, ok := ntpdOffset(ctx); ok {
		return result, nil
	}

	// Otherwise the agent measures the offset to the configured time server with an SNTP query.
	// This also works if no NTP service is running or ntpd is not synchronized yet.
	servers := ntpServersFromConfig(freebsdNtpConfigFiles, "0.freebsd.pool.ntp.org")
	return runSNTPCheck(ctx, servers), nil
}
