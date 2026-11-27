package utils

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

var (
	commandSysproc = &syscall.SysProcAttr{
		// Do not open any cmd windows
		// Run all processes in an own process group to be able to kill the process group and all child processes
		HideWindow:    true, //This will also hide the powershell windows if you run go test.
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
)

func handleCommandError(arg0 string, err error) int {
	if strings.HasSuffix(err.Error(), "file does not exist") || strings.HasSuffix(err.Error(), "executable file not found in %PATH%") {
		if _, err := os.Stat(arg0); os.IsNotExist(err) {
			return NotFound
		}
		return NotExecutable
	}
	return Unknown
}

// processGroup contains a started command and all processes it starts. Windows has no process groups which can be
// killed like on POSIX systems, the processes are put into a job object instead. All processes started by a process
// of the job belong to the job too, so the job can be terminated without searching the child processes.
//
// Searching the child processes with gopsutil took a snapshot of all processes for every process of the system and
// used a lot of CPU after every custom check (issue #118). It could also kill processes of other programs, because
// Windows reuses process ids and keeps the parent process id of a process whose parent has exited.
type processGroup struct {
	process *os.Process
	job     windows.Handle
}

// newProcessGroup puts the started process into a new job object. Processes which the command starts before it is
// added to the job (within microseconds after the start) are not part of the job.
// If the job object can not be created, only the process itself is killed.
func newProcessGroup(p *os.Process) *processGroup {
	group := &processGroup{
		process: p,
	}

	job, err := createKillOnCloseJob()
	if err != nil {
		log.Debugln("Command: could not create job object: ", err)
		return group
	}

	if err := assignProcessToJob(job, p.Pid); err != nil {
		log.Debugln("Command: could not add process ", p.Pid, " to the job object: ", err)
		_ = windows.CloseHandle(job)
		return group
	}

	group.job = job
	return group
}

// createKillOnCloseJob creates a job object which terminates all of its processes when its last handle is closed,
// e.g. if the agent exits while a command is running
func createKillOnCloseJob() (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	_, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		_ = windows.CloseHandle(job)
		return 0, err
	}
	return job, nil
}

func assignProcessToJob(job windows.Handle, pid int) error {
	// os.Process does not provide its handle. The process id can not be reused, as long as os.Process holds its
	// handle (until Wait returns).
	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(process)

	return windows.AssignProcessToJobObject(job, process)
}

// kill terminates the process and all processes it started
func (g *processGroup) kill() {
	if g.job == 0 {
		// nolint:errcheck
		g.process.Kill()
		return
	}
	// nolint:errcheck
	windows.TerminateJobObject(g.job, 1)
}

// close releases the job object, processes which are still running are terminated
func (g *processGroup) close() {
	if g.job == 0 {
		return
	}
	// nolint:errcheck
	windows.CloseHandle(g.job)
	g.job = 0
}
