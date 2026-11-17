//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package checks

import (
	"context"
	"math"

	"github.com/shirou/gopsutil/v4/cpu"
)

// The background sampler of cpu_sampler.go stores samples of the CPU times per core
type cpuTimes = []cpu.TimesStat

func readCpuTimes(ctx context.Context) (cpuTimes, error) {
	return cpu.TimesWithContext(ctx, true)
}

func cpuTimesCount(times cpuTimes) int {
	return len(times)
}

func saturatingSub(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return 0
}

func calcuateBusy(t cpu.TimesStat) (float64, float64, float64) {
	busy := t.User + t.Nice + t.System + t.Irq + t.Softirq + t.Steal
	idle := t.Idle + t.Iowait
	total := busy + idle
	return busy, idle, total
}

func calculateUsagePercentage(prev, current cpu.TimesStat) float64 {
	// This is highly inspired by
	// https://github.com/shirou/gopsutil/blob/aa9796d3d754cc700939236e2fa8c3a95374c61b/cpu/cpu.go#L93-L144
	// and htop
	// https://github.com/htop-dev/htop/blob/main/linux/LinuxProcessList.c#L1948-L2006
	// https://github.com/htop-dev/htop/blob/4102862d12695cdf003e2d51ef6ce5984b7136d7/linux/LinuxMachine.c#L402-L512 (code got moved 31.10.2024)
	// and yes - htop is reading /proc/stat as we do
	prevBusy, _, prevTotal := calcuateBusy(prev)
	currentBusy, _, currentTotal := calcuateBusy(current)

	if currentBusy <= prevBusy {
		return 0
	}

	if currentTotal <= prevTotal {
		return 100
	}
	return math.Min(100, math.Max(0, (currentBusy-prevBusy)/(currentTotal-prevTotal)*100))
}

// cpuResultFromTimes calculates the CPU usage between two samples of the CPU times
func cpuResultFromTimes(prevTimeStats, timeStats []cpu.TimesStat) *resultCpu {
	result := &resultCpu{}

	// Get CPU usage per core
	var cpuPercentages []float64
	for i := range prevTimeStats {
		cpuPercentage := calculateUsagePercentage(prevTimeStats[i], timeStats[i])
		cpuPercentages = append(cpuPercentages, cpuPercentage)
	}
	result.PercentagePerCore = cpuPercentages

	// Get total CPU usage as percentage
	var totalCpuPercentage float64
	for i := range cpuPercentages {
		totalCpuPercentage = totalCpuPercentage + cpuPercentages[i]
	}
	result.PercentageTotal = totalCpuPercentage / float64(len(cpuPercentages))

	// Get CPU timing details per core
	var timings []cpuDetails
	var totalUser, totalNice, totalSystem, totalIoWait, totalIdle, totalTotal float64
	for i := range prevTimeStats {
		// Ignore t.Irq + t.Softirq + t.Steal to get 100% over all values we list in the json
		prevTotal := prevTimeStats[i].User + prevTimeStats[i].Nice + prevTimeStats[i].System + prevTimeStats[i].Idle + prevTimeStats[i].Iowait
		currentTotal := timeStats[i].User + timeStats[i].Nice + timeStats[i].System + timeStats[i].Idle + timeStats[i].Iowait

		userDelta := saturatingSub(timeStats[i].User, prevTimeStats[i].User)
		niceDelta := saturatingSub(timeStats[i].Nice, prevTimeStats[i].Nice)
		systemDelta := saturatingSub(timeStats[i].System, prevTimeStats[i].System)
		ioWaitDelta := saturatingSub(timeStats[i].Iowait, prevTimeStats[i].Iowait)
		idleDelta := saturatingSub(timeStats[i].Idle, prevTimeStats[i].Idle)
		totalDelta := saturatingSub(currentTotal, prevTotal)

		if totalDelta == 0 {
			timings = append(timings, cpuDetails{})
		} else {
			// Do not Division by zero (not sure if this could ever happen to be honest but better safe than sorry)
			timings = append(timings, cpuDetails{
				User:   userDelta / totalDelta * 100,
				Nice:   niceDelta / totalDelta * 100,
				System: systemDelta / totalDelta * 100,
				Idle:   idleDelta / totalDelta * 100,
				Iowait: ioWaitDelta / totalDelta * 100,
			})
		}

		totalUser = totalUser + userDelta
		totalNice = totalNice + niceDelta
		totalSystem = totalSystem + systemDelta
		totalIoWait = totalIoWait + ioWaitDelta
		totalIdle = totalIdle + idleDelta
		totalTotal = totalTotal + totalDelta
	}
	result.DetailsPerCore = timings

	// Get total CPU timing details
	result.DetailsTotal = &cpuDetails{}
	if totalTotal > 0 {
		result.DetailsTotal = &cpuDetails{
			User:   totalUser / totalTotal * 100,
			Nice:   totalNice / totalTotal * 100,
			System: totalSystem / totalTotal * 100,
			Idle:   totalIdle / totalTotal * 100,
			Iowait: totalIoWait / totalTotal * 100,
		}
	}
	return result
}
