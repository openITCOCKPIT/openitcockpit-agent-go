//go:build !darwin && !linux && !freebsd && !windows

package loghandler

import (
	log "github.com/sirupsen/logrus"
)

// Fallback for all other operating systems
func (h *LogHandler) SetupNativeLogging() {
	// If nothing else is defined, just stick with the DefaultWriter (Stderr)
	log.Infoln("Native logging not supported for this OS, falling back to default writer.")
}
