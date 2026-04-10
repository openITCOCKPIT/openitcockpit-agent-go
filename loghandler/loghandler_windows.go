//go:build windows

package loghandler

import (
	"github.com/freman/eventloghook"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc/eventlog"
)

func (h *LogHandler) SetupNativeLogging() {

	elog, err := eventlog.Open("openITCOCKPIT-Agent")
	if err == nil {
		log.AddHook(eventloghook.NewHook(elog))
	}
}
