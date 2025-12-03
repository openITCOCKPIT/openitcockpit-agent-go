//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package cmd

import "os"

func PlatformMain() {
	if err := New().Execute(); err != nil {
		os.Exit(1)
	}
}
