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

	OsName       string // e.g. "ubuntu", "Windows", "macOS"
	OsVersion    string
	AgentVersion string
	Uptime       int64

	// Exclude from json
	LinuxPackages  []Package
	LinuxUpdates   []PackageUpdate
	WindowsApps    []WindowsApp
	WindowsUpdates []WindowsUpdate
	MacOSApps      []Package
	MacosUpdates   []MacosUpdate
}

// PackageInfo holds summary information about the package manager state
// Will be avaiable in the default JSON response of the agent
type PackageStats struct {
	InstalledPackages  int64
	UpgradablePackages int64
	SecurityUpdates    int64
	RebootRequired     bool
	LastError          error
	OperatingSystem    string
	PackageManager     string
}

type SoftwareCollector struct {
	Configuration *config.Configuration
	Result        chan *PackageInfo

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func (s *SoftwareCollector) Shutdown() {
	close(s.shutdown)
	s.wg.Wait()
}

// Start the check runner and returns immediatly (SHOULD NOT RUN IN GOROUTINE)
func (s *SoftwareCollector) Start(ctx context.Context) error {
	s.shutdown = make(chan struct{})

	log.Infoln("Packagemanager: Software inventory is enabled")

	// convert check intervall from minutes (config) to seconds
	checkInterval := s.Configuration.Packagemanager.CheckInterval * 60

	// Start a first run delayed after startup
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Normaly the check interval is a large value such as 1 hour or more
		// To get fresh data after the agent was started, we want to run the first check
		// after the agent is running for 90 seconds.
		// this ensures that the agent is running long enough (certificate exchange etc)
		// to not do the heavy collection work too early.
		firstRunDelay := 9 * time.Second //todo set to 90
		firstRunTrigger := time.NewTimer(firstRunDelay)
		checkTimeout := time.Duration(checkInterval-1) * time.Second
		defer firstRunTrigger.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-s.shutdown:
				if !ok {
					return
				}
			case <-firstRunTrigger.C:
				go s.runCollection(ctx, checkTimeout)
			}
		}
	}()

	// Start a periodic ticker to run the collection
	checkTimeout := time.Duration(checkInterval-1) * time.Second
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
		defer ticker.Stop()

		// Tell the webserver that we have pending data collection
		s.Result <- &PackageInfo{
			Enabled: true,
			Pending: true,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-s.shutdown:
				if !ok {
					return
				}
			case <-ticker.C:
				go s.runCollection(ctx, checkTimeout)
			}
		}
	}()

	return nil
}
