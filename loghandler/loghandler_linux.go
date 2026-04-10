//go:build linux

package loghandler

import (
	"fmt"
	"io"
	"os"

	"github.com/coreos/go-systemd/v22/journal"
	log "github.com/sirupsen/logrus"
)

// JournalHook ist ein einfacher Custom-Hook für Logrus
type JournalHook struct{}

func (h *JournalHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *JournalHook) Fire(entry *log.Entry) error {
	// Map Logrus Levels zu Systemd Priority Levels
	priority := journal.PriInfo
	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		priority = journal.PriDebug
	case log.WarnLevel:
		priority = journal.PriWarning
	case log.ErrorLevel:
		priority = journal.PriErr
	case log.FatalLevel, log.PanicLevel:
		priority = journal.PriCrit
	}

	// Create a message string that includes the entry message and fields if they exist
	msg := entry.Message
	if len(entry.Data) > 0 {
		// Optional: Append fields to the message if we want to see them in the journal text
		msg = fmt.Sprintf("%s %v", entry.Message, entry.Data)
	}

	// For real metadata in the journal, use Send instead of Print
	vars := make(map[string]string)
	vars["SYSLOG_IDENTIFIER"] = "openitcockpit-agent"
	for k, v := range entry.Data {
		vars[k] = fmt.Sprintf("%v", v)
	}

	// Send to Journald
	return journal.Send(msg, priority, vars)
}

func (h *LogHandler) SetupNativeLogging() {
	if h.LogPath == "" {
		if journal.Enabled() {
			log.AddHook(&JournalHook{})
			// Disable standard output to avoid duplicate logging
			log.SetOutput(io.Discard)
		} else {
			// Fallback in case no journal is available (e.g. Docker or old init)
			log.SetOutput(os.Stdout)
		}
	}
}
