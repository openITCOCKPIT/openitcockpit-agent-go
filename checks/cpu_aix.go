package checks

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// The background sampler of cpu_sampler.go stores samples of the CPU times of the LPAR
type cpuTimes = *aixCpuTimes

func readCpuTimes(_ context.Context) (cpuTimes, error) {
	times := &aixCpuTimes{Time: time.Now()}

	cpuTotal, err := perfstatCpuTotal()
	if err != nil {
		return nil, err
	}
	if times.Logical, err = decodeCpuTotalTicks(cpuTotal); err != nil {
		return nil, err
	}

	if times.Cpus, err = perfstatCpus(); err != nil {
		return nil, err
	}

	// Without the partition values the usage is calculated with the ticks of the logical CPUs,
	// which are not correct for shared processor partitions.
	partition, err := perfstatPartitionTotal()
	if err != nil {
		log.Debugln("CPU check: ", err)
		return times, nil
	}
	if err := decodePartition(partition, times); err != nil {
		log.Debugln("CPU check: ", err)
	}
	return times, nil
}

func cpuTimesCount(times cpuTimes) int {
	return len(times.Cpus)
}

func cpuResultFromTimes(prev, current cpuTimes) *resultCpu {
	return aixCpuResult(prev, current)
}
