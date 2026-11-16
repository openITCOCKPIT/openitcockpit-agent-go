//go:build aix || linux
// +build aix linux

// Parsers for the output of AIX commands (ps, who, no) and /etc/filesystems. The file is also built on Linux,
// so the unit tests run in the CI.

package checks

import (
	"bufio"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// aixPsArguments are the columns of ps on AIX. The arguments are the last column, as they can contain spaces.
// AIX does not support the columns rss and pagein of the BSD ps, rssize is the resident set size in KB.
// vsz is the size of the data in the virtual memory in KB, it does not include the text (it can be smaller than rssize).
// comm is the name of the program. The arguments can not be used for the name, as many processes overwrite them,
// e.g. "sshd-session: user@pts/0", "-ksh" (login shell) or "[IBM.Soft]".
var aixPsArguments = []string{"-e", "-o", "pid,ppid,user,nice,pcpu,pmem,vsz,rssize,stat,etime,comm,args"}

// aixPsColumns is the number of columns before the arguments
const aixPsColumns = 11

// aixPsCommColumn is the index of the column comm
const aixPsCommColumn = 10

// aixDefunct is the name ps reports for zombie processes
const aixDefunct = "<defunct>"

// AIX process states of ps -o stat
var aixProcessStates = map[string]string{
	"A": "Active", // Running or waiting for an event (AIX does not distinguish running and sleeping processes)
	"W": "Swapped",
	"I": "Idle", // Starting
	"Z": "Zombie",
	"T": "Stop",
	"O": "Nonexistent",
}

// parseAixPs parses the output of ps with aixPsArguments
//
//	PID     PPID     USER NI  %CPU  %MEM   VSZ   RSS S     ELAPSED COMMAND      COMMAND
//	  1        0     root 20   0.0   0.0   800   852 A 197-03:22:24 init         /etc/init
//	5898584    1     root --   0.0   0.0  3156  4016 A 197-03:21:47 topas_nmon   /usr/bin/topas_nmon -f -s 900
//	18809104 21496188 oitc 20  0.0   0.0  1780  7464 A       34:26 sshd-session sshd-session: oitc@pts/0
func parseAixPs(output string, now time.Time) []*resultProcess {
	processes := []*resultProcess{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		process, ok := parseAixPsLine(scanner.Text(), now)
		if ok {
			processes = append(processes, process)
		}
	}
	return processes
}

func parseAixPsLine(line string, now time.Time) (*resultProcess, bool) {
	fields, args := splitColumns(line, aixPsColumns)
	if len(fields) < aixPsCommColumn {
		return nil, false
	}
	comm := ""
	if len(fields) == aixPsColumns {
		comm = fields[aixPsCommColumn]
	}
	// Zombie processes have no program name
	if comm == aixDefunct {
		comm = ""
	}
	pid, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		// Header
		return nil, false
	}
	ppid, _ := strconv.ParseUint(fields[1], 10, 64)
	cpuPercent, _ := strconv.ParseFloat(fields[4], 64)
	memoryPercent, _ := strconv.ParseFloat(fields[5], 64)
	vsz, _ := strconv.ParseUint(fields[6], 10, 64)
	rss, _ := strconv.ParseUint(fields[7], 10, 64)

	process := &resultProcess{
		Pid:           pid,
		Ppid:          ppid,
		Username:      fields[2],
		Name:          aixProcessName(comm, args),
		CPUPercent:    cpuPercent,
		MemoryPercent: memoryPercent,
		Cmdline:       args,
		Status:        []string{aixProcessStatus(fields[8])},
		Exe:           aixProcessExe(comm, args),
		Nice:          aixNiceValue(fields[3]),
		Memory: &resultMemoryPosix{
			RSS: rss * 1024,
			VMS: vsz * 1024,
		},
	}
	if elapsed, ok := parseAixElapsedTime(fields[9]); ok {
		process.CreateTime = now.Add(-elapsed).Unix()
	}
	return process, true
}

// splitColumns returns the first count fields of a line separated by spaces and the rest of the line
func splitColumns(line string, count int) ([]string, string) {
	fields := make([]string, 0, count)
	rest := strings.TrimLeft(line, " \t")
	for len(fields) < count && rest != "" {
		end := strings.IndexAny(rest, " \t")
		if end < 0 {
			fields = append(fields, rest)
			rest = ""
			break
		}
		fields = append(fields, rest[:end])
		rest = strings.TrimLeft(rest[end:], " \t")
	}
	return fields, strings.TrimRight(rest, " \t")
}

// aixProcessExe returns the program, e.g. /usr/sbin/sshd for "/usr/sbin/sshd -D". If the first argument is not
// the program (e.g. "sshd-session: oitc@pts/0" or "-ksh"), the name of the program is returned.
func aixProcessExe(comm string, args string) string {
	fields := strings.Fields(args)
	if len(fields) > 0 && isAixProgramPath(fields[0], comm) {
		return fields[0]
	}
	return aixProcessName(comm, args)
}

// isAixProgramPath returns true if the first argument of a process is the path of its program comm
func isAixProgramPath(argument string, comm string) bool {
	if strings.HasPrefix(argument, "-") || strings.HasPrefix(argument, "[") || strings.Contains(argument, ":") {
		// The process changed its arguments, e.g. "-ksh" (login shell), "[IBM.Soft]" or "sshd: user@pts/0"
		return false
	}
	if comm == "" || strings.HasPrefix(argument, "/") {
		// An absolute path is the program, also if it is a link with another name, e.g. /usr/sbin/xntpd -> ntpd4
		return true
	}
	// comm can be shorter than the name of the program
	return strings.HasPrefix(path.Base(argument), comm)
}

// aixProcessName returns the name of the program, e.g. sshd-session for "sshd-session: oitc@pts/0".
// If ps reports no name, the name of the first argument is used.
func aixProcessName(comm string, args string) string {
	if comm != "" {
		return comm
	}
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return ""
	}
	return path.Base(fields[0])
}

func aixProcessStatus(state string) string {
	name, ok := aixProcessStates[state]
	if !ok {
		return state
	}
	return name
}

// aixNiceValue converts the nice value of AIX to the range of the other platforms, so nice_level means the same
// on all platforms: ps on AIX reports 0 to 39 with 20 as default, Linux, macOS and FreeBSD use -20 to 19 with
// 0 as default. The value is therefore the AIX value minus 20, e.g. a normal process (AIX 20) has 0.
// Processes with a fixed priority (scheduling policy FIFO or RR) have no nice value, ps reports "--", which is
// returned as 0.
func aixNiceValue(value string) int64 {
	nice, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return nice - 20
}

// parseAixElapsedTime parses the elapsed time of ps ([[dd-]hh:]mm:ss), e.g. 197-03:22:24
func parseAixElapsedTime(value string) (time.Duration, bool) {
	var days int64
	if i := strings.Index(value, "-"); i >= 0 {
		parsed, err := strconv.ParseInt(value[:i], 10, 64)
		if err != nil {
			return 0, false
		}
		days = parsed
		value = value[i+1:]
	}

	parts := strings.Split(value, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}
	var seconds int64
	for _, part := range parts {
		parsed, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return 0, false
		}
		seconds = seconds*60 + parsed
	}
	return time.Duration(days*86400+seconds) * time.Second, true
}

// parseAixWho parses the output of who. who runs in UTC (see runAixCommandInUTC), the year is not reported.
//
//	oitc        pts/0       Nov 13 15:05     (192.168.1.5)
//	root        lft0        Nov 10 08:00
func parseAixWho(output string, now time.Time) []*resultUser {
	users := []*resultUser{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		user := &resultUser{
			Name:     fields[0],
			Terminal: fields[1],
		}
		if started, ok := parseAixWhoTime(fields[2]+" "+fields[3]+" "+fields[4], now); ok {
			user.Started = started.Unix()
		}
		if len(fields) > 5 {
			user.Host = strings.Trim(fields[len(fields)-1], "()")
		}
		users = append(users, user)
	}
	return users
}

// parseAixWhoTime parses the login time of who (e.g. "Nov 13 15:05") in UTC. A login time in the future
// is from the previous year.
func parseAixWhoTime(value string, now time.Time) (time.Time, bool) {
	parsed, err := time.ParseInLocation("Jan 2 15:04", value, time.UTC)
	if err != nil {
		return time.Time{}, false
	}
	now = now.UTC()
	started := time.Date(now.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
	if started.After(now.Add(24 * time.Hour)) {
		started = started.AddDate(-1, 0, 0)
	}
	return started, true
}

// aixFilesystem is a stanza of /etc/filesystems
type aixFilesystem struct {
	Mountpoint string
	Values     map[string]string
}

// parseAixFilesystems parses /etc/filesystems
//
//	/home:
//	        dev             = /dev/hd1
//	        vfs             = jfs2
//	        mount           = true
func parseAixFilesystems(data []byte) []*aixFilesystem {
	var filesystems []*aixFilesystem
	var current *aixFilesystem
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && strings.HasSuffix(trimmed, ":") {
			current = &aixFilesystem{
				Mountpoint: strings.TrimSuffix(trimmed, ":"),
				Values:     map[string]string{},
			}
			filesystems = append(filesystems, current)
			continue
		}
		if current == nil {
			continue
		}
		key, value, found := strings.Cut(trimmed, "=")
		if found {
			current.Values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"`)
		}
	}
	return filesystems
}

// aixNetworkFilesystems are the vfs types of /etc/filesystems which require the network
var aixNetworkFilesystems = map[string]bool{"nfs": true, "nfs3": true, "nfs4": true, "cifs": true, "smbfs": true, "stnfs": true}

// aixFstabEntries returns the filesystems of /etc/filesystems which are mounted at boot (mount = automatic or true).
// Filesystems with mount = false are mounted by an administrator or a cluster (e.g. PowerHA), they are ignored.
func aixFstabEntries(filesystems []*aixFilesystem) []*resultFstabEntry {
	var entries []*resultFstabEntry
	for _, filesystem := range filesystems {
		mount := strings.ToLower(filesystem.Values["mount"])
		if mount != "automatic" && mount != "true" && mount != "readonly" {
			continue
		}
		if !strings.HasPrefix(filesystem.Mountpoint, "/") {
			continue
		}
		vfs := strings.ToLower(filesystem.Values["vfs"])
		entries = append(entries, &resultFstabEntry{
			Device:     aixFilesystemDevice(filesystem),
			Mountpoint: filepath.Clean(filesystem.Mountpoint),
			Fstype:     vfs,
			Options:    filesystem.Values["options"],
			Network:    aixNetworkFilesystems[vfs] || filesystem.Values["nodename"] != "",
		})
	}
	return entries
}

// aixFilesystemDevice returns the device, e.g. /dev/hd1 or server:/export for NFS
func aixFilesystemDevice(filesystem *aixFilesystem) string {
	device := filesystem.Values["dev"]
	if nodename := filesystem.Values["nodename"]; nodename != "" && nodename != "-" {
		return nodename + ":" + device
	}
	return device
}

// parseAixNetworkOption parses the value of a network option of "no -o <option>", e.g. "tcp_ephemeral_low = 32768"
func parseAixNetworkOption(output string) (int, bool) {
	_, value, found := strings.Cut(output, "=")
	if !found {
		return 0, false
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, false
	}
	return parsed, true
}
