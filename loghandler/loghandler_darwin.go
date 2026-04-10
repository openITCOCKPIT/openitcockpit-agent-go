//go:build darwin
// +build darwin

package loghandler

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func (h *LogHandler) SetupNativeLogging() {
	// With macOS 15: Just stay on Stderr. Launchd will automatically redirect it to the Unified Log.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
