//go:build !linux && !windows && !darwin
// +build !linux,!windows,!darwin

package packagemanager

import (
	"context"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// collectionSupported is false, the software inventory and the updates are not collected on this platform (e.g. FreeBSD and AIX)
const collectionSupported = false

func (s *SoftwareCollector) runCollection(parent context.Context, timeout time.Duration) {
	log.Errorln("Software collection not implemented for OS:", runtime.GOOS)
}
