//go:build aix
// +build aix

// Live tests of the custom checks on a real AIX system. They only run harmless commands (echo, sleep, ps).
// Build on Linux and run on AIX (no root required):
//
//	GOOS=aix GOARCH=ppc64 CGO_ENABLED=0 go test -c -o oitc-aix-customcheck.test ./checkrunner
//	./oitc-aix-customcheck.test -test.run 'TestAixLive|TestRun' -test.v

package checkrunner

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

func liveRunCustomCheck(t *testing.T, check *config.CustomCheck) *utils.CommandResult {
	t.Helper()
	executor := &CustomCheckExecutor{
		Configuration: check,
		ResultOutput:  make(chan *CustomCheckResult, 1),
		shutdown:      make(chan struct{}),
	}
	start := time.Now()
	executor.runCheck(context.Background(), time.Duration(check.Timeout)*time.Second)
	result := (<-executor.ResultOutput).Result
	t.Logf("%s: rc=%d after %s, stdout=%q", check.Name, result.RC, time.Since(start).Round(time.Millisecond), result.Stdout)
	return result
}

// TestAixLiveCustomChecks runs custom checks like a monitoring plugin: without shell, with ksh and with a timeout
func TestAixLiveCustomChecks(t *testing.T) {
	plugin := liveRunCustomCheck(t, &config.CustomCheck{
		Name:    "plugin_without_shell",
		Command: `/usr/bin/printf "OK - hello from AIX | value=1;2;3\n"`,
		Timeout: 10,
	})
	if plugin.RC != 0 || !strings.Contains(plugin.Stdout, "OK - hello from AIX | value=1;2;3") {
		t.Errorf("unexpected result %+v", plugin)
	}

	ksh := liveRunCustomCheck(t, &config.CustomCheck{
		Name:    "ksh_with_pipe_and_exit_code",
		Command: `echo "CRITICAL - $(uname -s) $(oslevel) $(echo pipe | tr a-z A-Z)"; exit 2`,
		Shell:   "/usr/bin/ksh",
		Timeout: 10,
	})
	if ksh.RC != 2 || !strings.Contains(ksh.Stdout, "CRITICAL - AIX") || !strings.Contains(ksh.Stdout, "PIPE") {
		t.Errorf("unexpected result %+v", ksh)
	}

	sh := liveRunCustomCheck(t, &config.CustomCheck{
		Name:    "sh_warning",
		Command: `echo "WARNING - sh"; exit 1`,
		Shell:   "/bin/sh",
		Timeout: 10,
	})
	if sh.RC != 1 || !strings.Contains(sh.Stdout, "WARNING - sh") {
		t.Errorf("unexpected result %+v", sh)
	}

	notFound := liveRunCustomCheck(t, &config.CustomCheck{
		Name:    "not_found",
		Command: "/usr/local/nagios/libexec/check_does_not_exist",
		Timeout: 10,
	})
	if notFound.RC != utils.NotFound {
		t.Errorf("unexpected result %+v", notFound)
	}

	// The timeout kills the whole process group, also the child processes of the shell
	timeout := liveRunCustomCheck(t, &config.CustomCheck{
		Name:    "timeout",
		Command: "sleep 47 & sleep 47; echo never",
		Shell:   "/usr/bin/ksh",
		Timeout: 2,
	})
	if timeout.RC != 124 {
		t.Errorf("unexpected result %+v", timeout)
	}
	time.Sleep(time.Second)
	output, _ := exec.Command("ps", "-e", "-o", "pid,ppid,args").Output()
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "sleep 47") {
			t.Errorf("child process was not killed: %s", line)
		}
	}
}
