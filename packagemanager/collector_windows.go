package packagemanager

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *SoftwareCollector) runCollection(parent context.Context, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	windows := WindowsUpdatesManager{}
	pkgInfo, err := windows.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
	if err != nil {
		log.Errorln("Error collecting software info for Windows:", err)
	}

	log.Debugln("Packagemanager: Software inventory collection for Windows completed")
	s.Result <- &pkgInfo
}
