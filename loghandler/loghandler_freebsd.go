//go:build freebsd

package loghandler

import (
	"log/syslog"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func (h *LogHandler) SetupNativeLogging() {
	// Connect to local Syslog daemon
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "openitcockpit-agent")
	if err == nil {
		log.AddHook(hook)
	}

	// Set formatter to text (Syslog usually doesn't need timestamps in the body)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
}
