package checks

import (
	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckNtp gathers information about  time offset between the system clock and the chosen time source (NTP)
type CheckNtp struct {
}

// Name will be used in the response as check name
func (c *CheckNtp) Name() string {
	return "ntp"
}

type resultNtp struct {
	//Meta data
	Timestamp      int64 // Timestamp of the last check evaluation
	TimestampMicro int64 // Timestamp of the last check evaluation in micro seconds (can be used if no NTP server is configured or on macOS)

	SyncStatus bool    // Is synchronized with an NTP server
	Offset     float64 // Time offset between local system and NTP server in seconds

}

// Configure the command or return false if the command was disabled
func (c *CheckNtp) Configure(config *config.Configuration) (bool, error) {
	return config.Ntp, nil
}
