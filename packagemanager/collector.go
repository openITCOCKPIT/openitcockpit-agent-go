package packagemanager

import (
	"context"
	"sync"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	log "github.com/sirupsen/logrus"
)

type PackageInfo struct {
	Enabled    bool
	Pending    bool
	LastUpdate int64
	Stats      PackageStats

	// Exclude from json
	LinuxPackages  []Package
	LinuxUpdates   []PackageUpdate
	WindowsApps    []WindowsApp
	WindowsUpdates []WindowsUpdate
	MacosApps      []Package
	MacosUpdates   []MacosUpdate
}

// PackageInfo holds summary information about the package manager state
// Will be avaiable in the default JSON response of the agent
type PackageStats struct {
	InstalledPackages  int64
	UpgradablePackages int64
	SecurityUpdates    int64
	RebootRequired     bool
	LastError          error  `json:"-"`         // exclude from json serialization
	LastErrorString    string `json:"LastError"` // used to serialize LastError
	OperatingSystem    string
	PackageManager     string
	OsName             string // e.g. "ubuntu", "Microsoft Windows 11 Enterprise", "macOS"
	OsVersion          string // e.g. "24.04", "24H2 (10.0.26100.7462 Build 26100.7462)", "12.5.1"
	OsFamily           string // e.g. "debian", "windows", "macos"
	AgentVersion       string // e.g. "3.4.0"
	Uptime             int64  // system uptime in seconds
}

type SoftwareCollector struct {
	Configuration *config.Configuration
	Result        chan *PackageInfo

	wg sync.WaitGroup
	//shutdown chan struct{} // deprecated, kept for compatibility if needed
	cancel context.CancelFunc
}

// Shutdown gracefully stops the collector by cancelling the context and waiting for goroutines to finish
func (s *SoftwareCollector) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

// Start the check runner and returns immediately (SHOULD NOT RUN IN GOROUTINE)
func (s *SoftwareCollector) Start(parentCtx context.Context) error {
	// Create a cancellable context for the collector
	ctx, cancel := context.WithCancel(parentCtx)
	s.cancel = cancel

	log.Infoln("Packagemanager: Software inventory is enabled")

	// convert check interval from minutes (config) to seconds
	checkInterval := s.Configuration.Packagemanager.CheckInterval * 60

	// Start a first run delayed after startup
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// Normally the check interval is a large value such as 1 hour or more
		// To get fresh data after the agent was started, we want to run the first check
		// after the agent is running for 90 seconds.
		// this ensures that the agent is running long enough (certificate exchange etc)
		// to not do the heavy collection work too early.
		firstRunDelay := 90 * time.Second
		firstRunTrigger := time.NewTimer(firstRunDelay)
		checkTimeout := time.Duration(checkInterval-1) * time.Second
		defer firstRunTrigger.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-firstRunTrigger.C:
				go s.runCollection(ctx, checkTimeout)
				return // exit after first run
			}
		}
	}()

	// Start a periodic ticker to run the collection
	checkTimeout := time.Duration(checkInterval-1) * time.Second
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
		defer ticker.Stop()

		// Tell the webserver that we have pending data collection
		select {
		case s.Result <- &PackageInfo{
			Enabled: true,
			Pending: true,
		}:
			// sent successfully
		default:
			// channel not ready, skip or log
			log.Warnln("Packagemanager: Unable to send pending status, channel not ready")
		}

		for {
			select {
			case <-ctx.Done():
				// Drain ticker channel to avoid goroutine leak
				for {
					select {
					case <-ticker.C:
						// drain
					default:
						return
					}
				}
			case <-ticker.C:
				go s.runCollection(ctx, checkTimeout)
			}
		}
	}()

	return nil
}
