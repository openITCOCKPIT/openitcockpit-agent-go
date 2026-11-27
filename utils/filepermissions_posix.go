//go:build linux || darwin || freebsd || aix
// +build linux darwin freebsd aix

package utils

import (
	"io/fs"
	"os"
)

// On Linux and macOS this is func just calls os.Chmod. This is only important for Windows systems
func Chmod(name string, mode fs.FileMode) error {
	return os.Chmod(name, mode)
}
