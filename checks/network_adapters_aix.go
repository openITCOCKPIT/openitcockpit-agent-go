package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckNetworkAdapters reports the Ethernet adapters of AIX with the state of EtherChannels,
// IEEE 802.3ad Link Aggregations and Network Interface Backups (lsdev, entstat -d)
type CheckNetworkAdapters struct {
}

// Name will be used in the response as check name
func (c *CheckNetworkAdapters) Name() string {
	return "network_adapters"
}

// Configure the command or return false if the command was disabled
func (c *CheckNetworkAdapters) Configure(config *config.Configuration) (bool, error) {
	return config.NetworkAdapters && aixCommandAvailable("lsdev") && aixCommandAvailable("entstat"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckNetworkAdapters) Run(ctx context.Context) (interface{}, error) {
	output, err := runAixCommand(ctx, "lsdev", "-Cc", "adapter")
	if err != nil {
		return nil, err
	}
	devices := parseLsdev(output)

	outputs := map[string]entstatOutput{}
	for _, device := range devices {
		if isEthernetAdapter(device) && device.State == "Available" {
			outputs[device.Name] = readEntstat(ctx, device.Name)
		}
	}
	return newNetworkAdapters(devices, outputs), nil
}

// readEntstat runs entstat -d. The members of an EtherChannel can not be opened (errno 19),
// their statistics are part of the output of the EtherChannel.
func readEntstat(ctx context.Context, name string) entstatOutput {
	output, err := runAixCommand(ctx, "entstat", "-d", name)
	return entstatOutput{Output: output, Err: err}
}
