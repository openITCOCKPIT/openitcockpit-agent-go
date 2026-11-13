//go:build linux || darwin || freebsd || aix
// +build linux darwin freebsd aix

package cmd

import "os"

func PlatformMain() {
	if err := New().Execute(); err != nil {
		os.Exit(1)
	}
}
