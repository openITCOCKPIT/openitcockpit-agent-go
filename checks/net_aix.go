package checks

import (
	"context"
	"net"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckNet) Run(ctx context.Context) (interface{}, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// The speed is the bit rate of perfstat_netinterface, the duplex mode is not available
	speeds := map[string]int64{}
	if netInterfaces, err := perfstatNetInterfaces(); err == nil {
		for _, netInterface := range netInterfaces {
			speeds[netInterface.Name] = int64(netInterface.Bitrate / 1000000)
		}
	}

	results := make(map[string]*resultNet)
	for _, iface := range interfaces {
		results[iface.Name] = &resultNet{
			Isup:   iface.Flags&net.FlagUp != 0,
			Duplex: DUPLEX_UNKNOWN,
			Speed:  speeds[iface.Name],
			MTU:    int64(iface.MTU),
		}
	}
	return results, nil
}
