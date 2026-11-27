package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

type testWindowsShell struct {
	command  string
	suffix   string
	args     string
	shell    string
	expected string
	exitCode int
}

var (
	testWS = []testWindowsShell{
		{
			command: `Write-Host $args[0]
			exit 1
			`,
			suffix:   "ps1",
			shell:    "powershell",
			args:     "Test",
			expected: "Test",
			exitCode: 1,
		},
		{
			command:  `Write-Host Test`,
			shell:    "powershell_command",
			expected: "Test",
			exitCode: 0,
		},
		{
			command: `
			WScript.Echo WScript.Arguments.item(0)
			`,
			suffix:   "vbs",
			shell:    "vbs",
			args:     "Test",
			expected: "Test",
			exitCode: 0,
		},
		{
			command: `
			echo %*
			exit 1
			`,
			suffix:   "bat",
			shell:    "bat",
			args:     "Test",
			expected: "Test",
			exitCode: 1,
		},
	}
)

func TestWindowsShells(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "*-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	writeCommand := func(command, suffix string) string {
		res := filepath.Join(tempDir, fmt.Sprintf("testcommand.%s", suffix))
		if err := os.WriteFile(res, []byte(command), 0600); err != nil {
			t.Fatal(err)
		}
		return res
	}

	for i, test := range testWS {
		cmd := test.command
		if test.suffix != "" {
			cmd = writeCommand(cmd, test.suffix)
		}
		res, err := RunCommand(context.Background(), CommandArgs{
			Command: cmd + " " + test.args,
			Timeout: time.Second * 5,
			Shell:   test.shell,
		})
		if strings.TrimSpace(res.Stdout) != test.expected {
			t.Error("test ", i, " expected output: ", test.expected, " got: ", res.Stdout)
		}
		if res.RC != test.exitCode {
			t.Error("test ", i, " unexpected exit code: ", res.RC)
		}
		if err != nil && res.RC == Unknown {
			t.Error(err)
		}
	}
}

// A PowerShell custom check which runs a program: the program inherits stdout of PowerShell
func TestCommandTimeoutWithPowershellChildHoldingStdout(t *testing.T) {
	assertTimeoutDoesNotHang(t, CommandArgs{
		Command: "ping.exe -n 30 127.0.0.1",
		Timeout: 2 * time.Second,
		Shell:   "powershell_command",
	})
}

// The script starts ping.exe as child process and writes its process id to the file of the first argument
const testStartChildScript = `$child = Start-Process -FilePath ping.exe -ArgumentList '-n','60','127.0.0.1' -WindowStyle Hidden -PassThru
Set-Content -Path $args[0] -Value $child.Id
`

func TestCommandTimeoutKillsChildProcesses(t *testing.T) {
	childPid := runCommandWithChild(t, testStartChildScript+"Start-Sleep -Seconds 60\n", time.Second*5)
	if !processExited(t, childPid) {
		t.Errorf("the child process %d of the command is still running after the timeout", childPid)
	}
}

func TestCommandExitKillsLeftChildProcesses(t *testing.T) {
	childPid := runCommandWithChild(t, testStartChildScript+"exit 0\n", time.Second*20)
	if !processExited(t, childPid) {
		t.Errorf("the child process %d which the command left behind is still running", childPid)
	}
}

// runCommandWithChild runs the PowerShell script and returns the process id of its child process
func runCommandWithChild(t *testing.T, script string, timeout time.Duration) int {
	t.Helper()
	tempDir := t.TempDir()
	scriptPath := filepath.Join(tempDir, "child.ps1")
	pidPath := filepath.Join(tempDir, "child.pid")
	if err := os.WriteFile(scriptPath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := RunCommand(context.Background(), CommandArgs{
		Command: scriptPath + " " + pidPath,
		Timeout: timeout,
		Shell:   "powershell",
	})
	if err != nil {
		t.Logf("command error: %s", err)
	}

	content, err := os.ReadFile(pidPath)
	if err != nil {
		t.Fatalf("the script did not write the process id of its child: %s", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		t.Fatalf("invalid process id %q", content)
	}
	return pid
}

// processExited returns true if the process does not exist anymore or exits within 5 seconds
func processExited(t *testing.T, pid int) bool {
	t.Helper()
	handle, err := windows.OpenProcess(windows.SYNCHRONIZE|windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return true
	}
	defer windows.CloseHandle(handle)

	event, err := windows.WaitForSingleObject(handle, 5000)
	if err != nil {
		t.Fatal(err)
	}
	if event == windows.WAIT_OBJECT_0 {
		return true
	}
	// Do not leave the test process behind
	// nolint:errcheck
	windows.TerminateProcess(handle, 1)
	return false
}
