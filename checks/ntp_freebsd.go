package checks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	log "github.com/sirupsen/logrus"
)

func (c *CheckNtp) Run(ctx context.Context) (interface{}, error) {

	//ntpq -c rv
	//associd=0 status=c016 leap_alarm, sync_unspec, 1 event, restart,
	//version="ntpd 4.2.8p18-a (1)", processor="amd64",
	//system="FreeBSD/15.0-RELEASE", leap=11, stratum=16, precision=-25,
	//rootdelay=0.000, rootdisp=0.015, refid=INIT, reftime=(no time),
	//clock=ecdb11d2.88100972  Wed, Dec  3 2025 20:45:22.531, peer=0, tc=3,
	//mintc=3, offset=+0.000000, frequency=+0.000, sys_jitter=0.000000,
	//clk_jitter=0.000, clk_wander=0.000, leapsec=201701010000,
	//expire=202606280000
	//
	// At least on my test system this command and the ntp services was crashing
	// so for now we use the same workaround as on macos

	// For now, we send a unixtimestamp in microseconds to the openITCOCKPIT server and compare those two clocks

	//result := &resultNtp{
	//	Timestamp:      time.Now().Unix(),
	//	TimestampMicro: time.Now().UnixMicro(),
	//	SyncStatus:     false,
	//	Offset:         0,
	//}
	//
	//return result, nil

	var err error
	var checkResult *resultNtp

	timeout := 10 * time.Second
	command := "ntpq -c rv"

	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: command,
		Timeout: timeout,
	})

	if err != nil || result.RC > 0 {
		checkResult = &resultNtp{
			Timestamp:      time.Now().Unix(),
			TimestampMicro: time.Now().UnixMicro(),
			SyncStatus:     false,
			Offset:         0,
		}

		log.Debugf("Error while executing '%v': %v\n", command, err)
		return checkResult, nil
	}

	syncStatus := false
	offset := 0.0

	for _, line := range strings.Split(string(result.Stdout), ",") {
		if strings.Contains(line, "sync_ntp") {
			syncStatus = true
		}
		if strings.Contains(line, "offset=") {
			var val float64
			_, err := fmt.Sscanf(line, "offset=%f", &val)
			if err == nil {
				offset = val / 1000.0 // ntpq offset is in ms, convert to seconds
			}
		}
	}

	checkResult = &resultNtp{
		Timestamp:      time.Now().Unix(),
		TimestampMicro: time.Now().UnixMicro(),
		SyncStatus:     syncStatus,
		Offset:         offset,
	}

	return checkResult, nil

}
