package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/net"
)

// netIoCounters returns the counters of each network interface (perfstat_netinterface). netstat -i of AIX
// does not report the bytes.
func netIoCounters(_ context.Context) ([]net.IOCountersStat, error) {
	interfaces, err := perfstatNetInterfaces()
	if err != nil {
		return nil, err
	}
	return aixNetIOCounters(interfaces), nil
}
