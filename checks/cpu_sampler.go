//go:build linux || darwin || freebsd || aix
// +build linux darwin freebsd aix

package checks

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// The CPU times are sampled in the background with this interval
	cpuSampleInterval = 5 * time.Second
	// The background sampler stops if the check was not executed for this time (e.g. disabled by a reload)
	cpuSamplerIdleTimeout = 10 * time.Minute
)

type cpuSample struct {
	time  time.Time
	stats cpuTimes
}

// cpuSampler reads the CPU times in the background, so the check can report the average CPU usage
// of the whole check interval. A short measurement while the agent executes all checks in parallel
// would mainly measure the agent itself.
type cpuSampler struct {
	mutex      sync.Mutex
	samples    []*cpuSample // Oldest first
	window     time.Duration
	running    bool
	stop       chan struct{} // Closed by Close, every sampler goroutine has its own channel
	lastAccess time.Time
}

var cpuSamplerState = &cpuSampler{}

// Close stops the background sampler
func (c *CheckCpu) Close() {
	cpuSamplerState.close()
}

func (s *cpuSampler) close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		close(s.stop)
		s.running = false
	}
	s.samples = nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckCpu) Run(ctx context.Context) (interface{}, error) {
	base, current, err := cpuTimesOfWindow(ctx, c.checkInterval)
	if err != nil {
		return nil, err
	}
	return cpuResultFromTimes(base, current), nil
}

// cpuTimesOfWindow returns the CPU times at the beginning of the check interval (base) and the current CPU times,
// so the checks can report the average usage of the whole check interval
func cpuTimesOfWindow(ctx context.Context, checkInterval int64) (cpuTimes, cpuTimes, error) {
	window := time.Duration(checkInterval) * time.Second
	if window <= 0 {
		window = 30 * time.Second
	}

	current, err := readCpuTimes(ctx)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()

	base := cpuSamplerState.base(now, window, cpuTimesCount(current))
	if base == nil {
		// No sample of the background sampler yet, e.g. right after the start of the agent
		log.Debugln("CPU check: no background sample available, measuring for 1 second")
		if err := sleepWithContext(ctx, 1*time.Second); err != nil {
			return nil, nil, err
		}
		base = current
		if current, err = readCpuTimes(ctx); err != nil {
			return nil, nil, err
		}
	}

	if cpuTimesCount(base) != cpuTimesCount(current) {
		return nil, nil, fmt.Errorf("number of CPU cores has changed %v != %v", cpuTimesCount(base), cpuTimesCount(current))
	}
	return base, current, nil
}

// base returns the CPU times at the beginning of the window and starts the background sampler if required
func (s *cpuSampler) base(now time.Time, window time.Duration, cores int) cpuTimes {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.window = window
	s.lastAccess = now
	if !s.running {
		s.running = true
		s.stop = make(chan struct{})
		go s.run(s.stop)
	}
	return selectCpuBaseSample(s.samples, now, window, cores)
}

// selectCpuBaseSample returns the newest sample which is at least as old as the window.
// If the sampler runs shorter than the window, the oldest sample is used, if it is at least 1 second old.
func selectCpuBaseSample(samples []*cpuSample, now time.Time, window time.Duration, cores int) cpuTimes {
	var base *cpuSample
	for _, sample := range samples {
		if cpuTimesCount(sample.stats) != cores {
			// The number of CPU cores has changed (CPU hotplug)
			continue
		}
		if now.Sub(sample.time) >= window || base == nil {
			base = sample
		}
	}
	if base == nil || now.Sub(base.time) < time.Second {
		return nil
	}
	return base.stats
}

func (s *cpuSampler) run(stop chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			_ = recoveredPanic("CPU sampler", r)
			s.mutex.Lock()
			// The next check run starts the sampler again
			if s.stop == stop {
				s.running = false
			}
			s.mutex.Unlock()
		}
	}()
	for {
		stats, err := readCpuTimes(context.Background())
		now := time.Now()

		s.mutex.Lock()
		select {
		case <-stop:
			// Closed, a new sampler may already run
			s.mutex.Unlock()
			return
		default:
		}
		if now.Sub(s.lastAccess) > cpuSamplerIdleTimeout {
			// The check is not executed anymore
			s.running = false
			s.samples = nil
			s.mutex.Unlock()
			return
		}
		if err != nil {
			log.Debugln("CPU check: could not read the CPU times: ", err)
		} else {
			if len(s.samples) > 0 && cpuTimesCount(s.samples[len(s.samples)-1].stats) != cpuTimesCount(stats) {
				// The number of CPU cores has changed (CPU hotplug), the old samples can not be compared anymore
				s.samples = nil
			}
			s.samples = append(s.samples, &cpuSample{time: now, stats: stats})
			// Keep the samples of the window and one older sample as base
			keep := now.Add(-s.window - 2*cpuSampleInterval)
			for len(s.samples) > 2 && s.samples[1].time.Before(keep) {
				s.samples = s.samples[1:]
			}
		}
		s.mutex.Unlock()

		select {
		case <-stop:
			return
		case <-time.After(cpuSampleInterval):
		}
	}
}

// Like sleep but can be canceled by context
func sleepWithContext(ctx context.Context, interval time.Duration) error {
	var timer = time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
