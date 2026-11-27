//go:build darwin || freebsd || aix
// +build darwin freebsd aix

package checks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// The result is cached, so public NTP servers are not queried with every check run
	sntpCacheDuration = 5 * time.Minute
	sntpQueryTimeout  = 3 * time.Second
)

var (
	sntpCacheMutex sync.Mutex
	sntpCache      *sntpResult
	sntpCacheTime  time.Time
)

// runSNTPCheck measures the offset of the local clock to the first reachable time server.
// If no time server is reachable, the openITCOCKPIT server compares the timestamps of the agent
// with its own clock (SyncStatus = false).
func runSNTPCheck(ctx context.Context, servers []string) *resultNtp {
	result := &resultNtp{
		Timestamp:      time.Now().Unix(),
		TimestampMicro: time.Now().UnixMicro(),
	}

	measurement, measuredAt, err := measureNtpOffset(ctx, servers)
	if err != nil {
		result.QueryError = err.Error()
		return result
	}

	result.SyncStatus = true
	result.Offset = measurement.Offset.Seconds()
	result.NtpServer = measurement.Server
	result.NtpAddress = measurement.Address
	result.Stratum = measurement.Stratum
	result.RoundTripDelay = measurement.Delay.Seconds()
	result.MeasuredAt = measuredAt.Unix()
	return result
}

// measureNtpOffset returns the cached measurement or queries the time servers
func measureNtpOffset(ctx context.Context, servers []string) (*sntpResult, time.Time, error) {
	sntpCacheMutex.Lock()
	defer sntpCacheMutex.Unlock()

	if sntpCache != nil && time.Since(sntpCacheTime) < sntpCacheDuration {
		return sntpCache, sntpCacheTime, nil
	}

	var errs []string
	for _, server := range servers {
		measurement, err := querySNTP(ctx, server, sntpQueryTimeout)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		sntpCache, sntpCacheTime = measurement, time.Now()
		return sntpCache, sntpCacheTime, nil
	}
	return nil, time.Time{}, fmt.Errorf("could not query time server: %s", strings.Join(errs, "; "))
}

// ntpServersFromConfig returns the servers of the first config file which defines servers
func ntpServersFromConfig(configFiles []string, defaultServer string) []string {
	for _, configFile := range configFiles {
		if data, err := os.ReadFile(configFile); err == nil {
			if servers := parseNtpConfServers(string(data)); len(servers) > 0 {
				return servers
			}
		}
	}
	return []string{defaultServer}
}
