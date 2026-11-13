//go:build aix || linux
// +build aix linux

package checks

import (
	"reflect"
	"testing"
	"time"
)

func TestParseAixPs(t *testing.T) {
	now := time.Date(2026, 11, 13, 15, 0, 0, 0, time.UTC)
	output := `     PID     PPID     USER NI  %CPU  %MEM   VSZ   RSS S     ELAPSED COMMAND
       1        0     root 20   0.0   0.0   800   852 A 197-03:22:24 /etc/init
 5898584        1     root --   0.1   0.0  3156  4016 A    01:02:03 /usr/bin/topas_nmon  -f -s 900
 8061234  4587958     oitc 25  12.5   1.4 16360 37952 T       05:09 sshd: oitc@pts/0
  700000        1   daemon 20   0.0   0.0     0     0 Z       00:01
`
	processes := parseAixPs(output, now)
	if len(processes) != 4 {
		t.Fatalf("unexpected number of processes %d", len(processes))
	}

	want := &resultProcess{
		Pid:      1,
		Ppid:     0,
		Username: "root",
		Name:     "init",
		Cmdline:  "/etc/init",
		Status:   []string{"Active"},
		Exe:      "/etc/init",
		Nice:     0,
		Memory: &resultMemoryPosix{
			RSS: 852 * 1024,
			VMS: 800 * 1024,
		},
		CreateTime: now.Add(-(197*24*time.Hour + 3*time.Hour + 22*time.Minute + 24*time.Second)).Unix(),
	}
	if !reflect.DeepEqual(processes[0], want) {
		t.Errorf("got %+v, want %+v", processes[0], want)
	}

	if processes[1].Cmdline != "/usr/bin/topas_nmon  -f -s 900" || processes[1].Name != "topas_nmon" || processes[1].Nice != 0 {
		t.Errorf("unexpected process %+v", processes[1])
	}
	if processes[1].CreateTime != now.Add(-(time.Hour + 2*time.Minute + 3*time.Second)).Unix() {
		t.Errorf("unexpected create time %d", processes[1].CreateTime)
	}

	sshd := processes[2]
	if sshd.Username != "oitc" || sshd.Nice != 5 || sshd.CPUPercent != 12.5 || sshd.MemoryPercent != 1.4 {
		t.Errorf("unexpected process %+v", sshd)
	}
	if sshd.Name != "sshd:" || sshd.Exe != "sshd:" || sshd.Cmdline != "sshd: oitc@pts/0" || sshd.Status[0] != "Stop" {
		t.Errorf("unexpected process %+v", sshd)
	}

	zombie := processes[3]
	if zombie.Status[0] != "Zombie" || zombie.Cmdline != "" || zombie.Name != "" || zombie.Exe != "" {
		t.Errorf("unexpected process %+v", zombie)
	}
}

func TestParseAixElapsedTime(t *testing.T) {
	tests := map[string]time.Duration{
		"00:05":        5 * time.Second,
		"12:34":        12*time.Minute + 34*time.Second,
		"01:02:03":     time.Hour + 2*time.Minute + 3*time.Second,
		"2-00:00:01":   48*time.Hour + time.Second,
		"197-03:22:24": 197*24*time.Hour + 3*time.Hour + 22*time.Minute + 24*time.Second,
	}
	for value, want := range tests {
		got, ok := parseAixElapsedTime(value)
		if !ok || got != want {
			t.Errorf("%s: got %v %v, want %v", value, got, ok, want)
		}
	}
	for _, value := range []string{"", "-", "12", "a:b", "1:2:3:4"} {
		if _, ok := parseAixElapsedTime(value); ok {
			t.Errorf("%q should be invalid", value)
		}
	}
}

func TestSplitColumns(t *testing.T) {
	fields, rest := splitColumns("  a  b\tc   d e  ", 3)
	if !reflect.DeepEqual(fields, []string{"a", "b", "c"}) || rest != "d e" {
		t.Errorf("got %q %q", fields, rest)
	}
	fields, rest = splitColumns("a b", 3)
	if !reflect.DeepEqual(fields, []string{"a", "b"}) || rest != "" {
		t.Errorf("got %q %q", fields, rest)
	}
}

func TestParseAixWho(t *testing.T) {
	now := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	output := `oitc        pts/0       Jan  2 09:05     (192.168.1.5)
root        lft0        Dec 30 08:00
`
	want := []*resultUser{
		{Name: "oitc", Terminal: "pts/0", Host: "192.168.1.5", Started: time.Date(2026, 1, 2, 9, 5, 0, 0, time.UTC).Unix()},
		// A login time in the future is from the previous year
		{Name: "root", Terminal: "lft0", Started: time.Date(2025, 12, 30, 8, 0, 0, 0, time.UTC).Unix()},
	}
	if got := parseAixWho(output, now); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if got := parseAixWho("", now); got == nil || len(got) != 0 {
		t.Errorf("expected an empty list, got %+v", got)
	}
}

func TestParseAixFilesystems(t *testing.T) {
	data := []byte(`* comment
/:
	dev		= /dev/hd4
	vfs		= jfs2
	log		= /dev/hd8
	mount		= automatic
	type		= bootfs

/home:
	dev		= /dev/hd1
	vfs		= jfs2
	mount		= true

/proc:
	dev	  = /proc
	vol	  = "/proc"
	mount	  = true
	vfs	  = procfs

/data:
	dev		= /dev/datalv
	vfs		= jfs2
	mount		= false
	mountguard	= yes

/mnt/backup:
	dev		= /export/backup
	vfs		= nfs
	nodename	= nas
	mount		= true
	options		= bg,hard,intr
`)
	filesystems := parseAixFilesystems(data)
	if len(filesystems) != 5 || filesystems[2].Values["vol"] != "/proc" {
		t.Fatalf("unexpected filesystems %+v", filesystems)
	}

	want := []*resultFstabEntry{
		{Device: "/dev/hd4", Mountpoint: "/", Fstype: "jfs2"},
		{Device: "/dev/hd1", Mountpoint: "/home", Fstype: "jfs2"},
		{Device: "/proc", Mountpoint: "/proc", Fstype: "procfs"},
		{Device: "nas:/export/backup", Mountpoint: "/mnt/backup", Fstype: "nfs", Options: "bg,hard,intr", Network: true},
	}
	if got := aixFstabEntries(filesystems); !reflect.DeepEqual(got, want) {
		for _, entry := range got {
			t.Logf("%+v", entry)
		}
		t.Errorf("unexpected entries")
	}
}

func TestParseAixNetworkOption(t *testing.T) {
	if value, ok := parseAixNetworkOption("tcp_ephemeral_low = 32768\n"); !ok || value != 32768 {
		t.Errorf("got %d %v", value, ok)
	}
	if _, ok := parseAixNetworkOption("0516-123 error"); ok {
		t.Errorf("expected an error")
	}
}

func TestParseAixNetstat(t *testing.T) {
	// netstat -an -f inet of AIX has the format of macOS, "tcp" and "udp" are dual stack sockets
	output := []byte(`Active Internet connections (including servers)
Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
tcp        0      0  *.*                    *.*                    CLOSED
tcp6       0      0  *.22                   *.*                    LISTEN
tcp4       0      0  *.25                   *.*                    LISTEN
tcp        0      0  *.111                  *.*                    LISTEN
tcp4       0      0  192.0.2.10.22          198.51.100.7.57422     ESTABLISHED
udp        0      0  *.*                    *.*                   
udp4       0      0  *.123                  *.*                   
Active UNIX domain sockets
SADR/PCB  Type   Recv-Q Send-Q   Inode     Conn     Refs  Nextref     Addr
f1000f0000b2c008 dgram       0      0 f1000a00002e9020                0        0        0 /dev/.SRC-unix/SRCJHcda
`)
	var bound []*networkSocket
	for _, socket := range parseDarwinNetstat(output) {
		if socket.LocalPort != 0 {
			bound = append(bound, socket)
		}
	}
	if len(bound) != 5 {
		t.Fatalf("unexpected sockets %d", len(bound))
	}
	if bound[0].Protocol != "tcp" || bound[0].Family != 6 || bound[0].LocalPort != 22 || bound[0].State != "listen" {
		t.Errorf("unexpected socket %+v", bound[0])
	}
	if bound[2].Family != 6 || bound[2].LocalPort != 111 {
		t.Errorf("dual stack socket %+v", bound[2])
	}
	established := bound[3]
	if established.LocalIP != "192.0.2.10" || established.RemoteIP != "198.51.100.7" || established.RemotePort != 57422 || established.State != "established" {
		t.Errorf("unexpected socket %+v", established)
	}
	if bound[4].Protocol != "udp" || bound[4].LocalPort != 123 || bound[4].State != "listen" {
		t.Errorf("unexpected socket %+v", bound[4])
	}
}
