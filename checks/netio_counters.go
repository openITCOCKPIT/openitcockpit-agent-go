//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/net"
)

// netIoCounters returns the counters of each network interface
func netIoCounters(ctx context.Context) ([]net.IOCountersStat, error) {
	return net.IOCountersWithContext(ctx, true)
}
