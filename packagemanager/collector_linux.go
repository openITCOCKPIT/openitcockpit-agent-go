package packagemanager

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *SoftwareCollector) runCollection(parent context.Context, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	apt := AptManager{}
	if apt.IsAvailable() {
		pkgInfo, err := apt.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
		if err != nil {
			log.Errorln("Error collecting package info with apt:", err)
		}

		log.Debugln("Packagemanager: Software inventory collection with apt completed")
		s.Result <- &pkgInfo
		return
	}

	dnf := DnfManager{}
	if dnf.IsAvailable() {
		pkgInfo, err := dnf.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
		if err != nil {
			log.Errorln("Error collecting package info with yum/dnf:", err)
		}

		log.Debugln("Packagemanager: Software inventory collection with dnf completed")
		s.Result <- &pkgInfo
		return
	}

	pacman := PacmanManager{}
	if pacman.IsAvailable() {
		pkgInfo, err := pacman.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
		if err != nil {
			log.Errorln("Error collecting package info with pacman:", err)
		}

		log.Debugln("Packagemanager: Software inventory collection with pacman completed")
		s.Result <- &pkgInfo
		return
	}

	zypper := ZypperManager{}
	if zypper.IsAvailable() {
		pkgInfo, err := zypper.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
		if err != nil {
			log.Errorln("Error collecting package info with zypper:", err)
		}

		log.Debugln("Packagemanager: Software inventory collection with zypper completed")
		s.Result <- &pkgInfo
		return
	}

	// Fallback to rpm if no other package manager was found
	// This is the last resort keep it the last option
	rpm := RpmManager{}
	if rpm.IsAvailable() {
		pkgInfo, err := rpm.CollectPackageInfo(ctx, s.Configuration.Packagemanager.LimitDescriptionLength, s.Configuration.Packagemanager.EnableUpdateCheck)
		if err != nil {
			log.Errorln("Error collecting package info with rpm:", err)
		}

		log.Debugln("Packagemanager: Software inventory collection with rpm completed")
		s.Result <- &pkgInfo
		return
	}

	log.Errorln("No known package manager found on Linux")
}
