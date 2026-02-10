//go:build !linux && !windows && !darwin
// +build !linux,!windows,!darwin

package packagemanager

import (
	"context"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *SoftwareCollector) runCollection(parent context.Context, timeout time.Duration) {
	log.Errorln("Software collection not implemented for OS:", runtime.GOOS)
}
