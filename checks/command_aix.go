package checks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// aixCommandTimeout is the timeout of a single AIX command like lsvg or errpt
const aixCommandTimeout = 20 * time.Second

// runAixCommand runs an AIX command with the C locale (English messages and column names) and returns stdout.
// On errors the message contains the output of the command, e.g. the AIX message 0516-010.
func runAixCommand(ctx context.Context, name string, args ...string) (string, error) {
	return runAixCommandWithEnv(ctx, nil, name, args...)
}

// runAixCommandInUTC runs an AIX command like runAixCommand, but the command prints and expects all times in UTC.
//
// AIX sets the time zone as POSIX string (e.g. TZ=CET-1CEST,M3.5.0,M10.5.0 in /etc/environment). Go can not read
// such a TZ variable and silently uses UTC as local time, so local times of a command can not be converted correctly.
func runAixCommandInUTC(ctx context.Context, name string, args ...string) (string, error) {
	return runAixCommandWithEnv(ctx, []string{"TZ=UTC0"}, name, args...)
}

// runAixCommandWithEnv runs an AIX command with additional environment variables (later values win)
func runAixCommandWithEnv(ctx context.Context, env []string, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, aixCommandTimeout)
	defer cancel()

	command := exec.CommandContext(ctx, name, args...)
	command.Env = append(os.Environ(), "LANG=C", "LC_ALL=C")
	command.Env = append(command.Env, env...)
	output, err := command.Output()
	if err == nil {
		return string(output), nil
	}
	if ctx.Err() != nil {
		return "", fmt.Errorf("%s: timeout after %s", name, aixCommandTimeout)
	}
	var stderr string
	if exitError, ok := err.(*exec.ExitError); ok {
		stderr = string(exitError.Stderr)
	}
	return string(output), commandError(name, err, stderr)
}

// aixCommandAvailable returns true if the command is installed
func aixCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
