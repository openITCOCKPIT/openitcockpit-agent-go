//go:build aix
// +build aix

// Live tests on a real AIX system. They compare the agent with the AIX tools and only read.
// Build on Linux and run on AIX (no root required):
//
//	GOOS=aix GOARCH=ppc64 CGO_ENABLED=0 go test -c -o oitc-aix-live.test ./checks
//	./oitc-aix-live.test -test.run 'TestAixLive' -test.v
//
// IP addresses and the host name are removed from the logs.

package checks

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

var (
	liveIPv4Pattern    = regexp.MustCompile(`\b(\d{1,3})\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	liveNetworkPattern = regexp.MustCompile(`\s\d{1,3}\.\d{1,3}\.\d{1,3}\s`)         // netstat -in, e.g. 10.0.1
	liveMacPattern     = regexp.MustCompile(`\b[0-9a-f]{1,2}(\.[0-9a-f]{1,2}){5}\b`) // netstat -in, e.g. 56.3f.9.17.d7.2
	livePvidPattern    = regexp.MustCompile(`\b[0-9a-f]{16}\b`)                      // lspv
	liveIdPattern      = regexp.MustCompile(`processor="[^"]*"`)                     // ntpq -c rv (machine id)
)

// liveRedact removes IPv4 addresses (except loopback and wildcard), networks, MAC addresses, PVIDs, the machine id
// and the host name
func liveRedact(text string) string {
	text = liveIPv4Pattern.ReplaceAllStringFunc(text, liveRedactIPv4)
	text = liveNetworkPattern.ReplaceAllString(text, " (net) ")
	text = liveMacPattern.ReplaceAllString(text, "(mac)")
	text = livePvidPattern.ReplaceAllString(text, "(pvid)")
	text = liveIdPattern.ReplaceAllString(text, `processor="(id)"`)
	if hostname, err := os.Hostname(); err == nil && len(hostname) > 2 {
		text = strings.ReplaceAll(text, hostname, "(host)")
		short, _, _ := strings.Cut(hostname, ".")
		if len(short) > 2 {
			text = strings.ReplaceAll(text, short, "(host)")
		}
	}
	return text
}

func liveRedactIPv4(address string) string {
	if strings.HasPrefix(address, "127.") || address == "0.0.0.0" || strings.HasPrefix(address, "255.") {
		return address
	}
	return "(ip)"
}

func liveCommand(t *testing.T, name string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	command.Env = append(os.Environ(), "LANG=C", "LC_ALL=C")
	output, err := command.CombinedOutput()
	t.Logf("$ %s\n%s", strings.Join(append([]string{name}, args...), " "), liveRedact(string(output)))
	if err != nil {
		t.Logf("(%s: %v)", name, err)
	}
	return string(output)
}

func liveJSON(t *testing.T, title string, value interface{}) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Errorf("%s: %v", title, err)
		return
	}
	t.Logf("%s:\n%s", title, liveRedact(string(data)))
}

func liveConfig() *config.Configuration {
	return &config.Configuration{
		CheckInterval: 5,
		Processes:     true,
		User:          true,
		DiskIo:        true,
		Netstats:      true,
		NetIo:         true,
		Ntp:           true,
		Mounts:        true,
		NetworkPorts:  true,
	}
}

func liveRunCheck(t *testing.T, check Check) interface{} {
	t.Helper()
	enabled, err := check.Configure(liveConfig())
	if err != nil || !enabled {
		t.Fatalf("%s: enabled=%v err=%v", check.Name(), enabled, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	start := time.Now()
	result, err := check.Run(ctx)
	t.Logf("check %s took %s", check.Name(), time.Since(start).Round(time.Millisecond))
	if err != nil {
		t.Logf("check %s: ERROR %v", check.Name(), err)
		return nil
	}
	return result
}

// netstatInPackets returns Ipkts and Opkts of the first line of each interface of netstat -in
func netstatInPackets(output string) map[string][2]uint64 {
	packets := map[string][2]uint64{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 9 || fields[0] == "Name" {
			continue
		}
		if _, ok := packets[fields[0]]; ok {
			continue
		}
		// Name Mtu Network Address Ipkts Ierrs Opkts Oerrs Coll, the address is missing for some networks
		offset := len(fields) - 9
		inPackets, errIn := strconv.ParseUint(fields[4+offset], 10, 64)
		outPackets, errOut := strconv.ParseUint(fields[6+offset], 10, 64)
		if errIn == nil && errOut == nil {
			packets[fields[0]] = [2]uint64{inPackets, outPackets}
		}
	}
	return packets
}

// TestAixLiveNetInterfaces compares perfstat_netinterface with netstat -in (read before and after netstat)
func TestAixLiveNetInterfaces(t *testing.T) {
	before, err := perfstatNetInterfaces()
	if err != nil {
		t.Fatal(err)
	}
	output := liveCommand(t, "netstat", "-in")
	after, err := perfstatNetInterfaces()
	if err != nil {
		t.Fatal(err)
	}
	liveJSON(t, "perfstat_netinterface", after)

	netstat := netstatInPackets(output)
	if len(after) == 0 {
		t.Fatal("no network interfaces")
	}
	for i, netInterface := range after {
		packets, ok := netstat[netInterface.Name]
		if !ok || i >= len(before) || before[i].Name != netInterface.Name {
			t.Logf("%s: not compared", netInterface.Name)
			continue
		}
		if packets[0] < before[i].PacketsRecv || packets[0] > netInterface.PacketsRecv {
			t.Errorf("%s: netstat Ipkts %d not within perfstat %d..%d", netInterface.Name, packets[0], before[i].PacketsRecv, netInterface.PacketsRecv)
		}
		if packets[1] < before[i].PacketsSent || packets[1] > netInterface.PacketsSent {
			t.Errorf("%s: netstat Opkts %d not within perfstat %d..%d", netInterface.Name, packets[1], before[i].PacketsSent, netInterface.PacketsSent)
		}
		t.Logf("%s: packets match netstat -in (in %d, out %d)", netInterface.Name, packets[0], packets[1])
	}
	liveCommand(t, "entstat", "-d", "ent0")
}

// TestAixLiveDisks compares the names and volume groups of perfstat_disk with lspv
func TestAixLiveDisks(t *testing.T) {
	disks, err := perfstatDisks()
	if err != nil {
		t.Fatal(err)
	}
	liveJSON(t, "perfstat_disk", disks)
	liveCommand(t, "lsattr", "-E", "-l", "sys0", "-a", "iostat")

	output := liveCommand(t, "lspv")
	volumeGroups := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			volumeGroups[fields[0]] = fields[2]
		}
	}
	for _, disk := range disks {
		vg, ok := volumeGroups[disk.Name]
		if !ok {
			t.Errorf("%s: not reported by lspv", disk.Name)
			continue
		}
		if vg != "None" && vg != disk.VolumeGroup {
			t.Errorf("%s: volume group %q, lspv %q", disk.Name, disk.VolumeGroup, vg)
		}
	}
}

// TestAixLiveDiskCounters checks the units of perfstat_disk while iostat -D runs (iostat counts the disk I/O
// while it is running, even if the I/O history is disabled). Only with OITC_AIX_DISKIO=1: writes 64 MB to the
// temporary directory of the test (TMPDIR) and removes them.
func TestAixLiveDiskCounters(t *testing.T) {
	if os.Getenv("OITC_AIX_DISKIO") != "1" {
		t.Skip("set OITC_AIX_DISKIO=1 to write 64 MB to TMPDIR while iostat -D runs")
	}
	iostat := exec.Command("iostat", "-D", "12", "1")
	iostat.Env = append(os.Environ(), "LANG=C", "LC_ALL=C")
	var iostatOutput strings.Builder
	iostat.Stdout = &iostatOutput
	iostat.Stderr = &iostatOutput
	if err := iostat.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	before, err := perfstatDisks()
	if err != nil {
		t.Fatal(err)
	}
	startTime := time.Now()

	file := t.TempDir() + "/oitc-diskio-test"
	liveCommand(t, "dd", "if=/dev/zero", "of="+file, "bs=1048576", "count=64")
	liveCommand(t, "sync")
	time.Sleep(8 * time.Second)

	after, err := perfstatDisks()
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(startTime)
	_ = os.Remove(file)
	_ = iostat.Wait()
	t.Logf("$ iostat -D 12 1\n%s", liveRedact(iostatOutput.String()))

	t.Logf("perfstat over %s:", elapsed.Round(time.Millisecond))
	for i, disk := range after {
		if i >= len(before) || before[i].Name != disk.Name {
			continue
		}
		prev := before[i]
		writeBlocks := disk.WriteBlocks - prev.WriteBlocks
		transfers := disk.Transfers - prev.Transfers
		busyTicks := disk.BusyTicks - prev.BusyTicks
		writeService := disk.WriteService - prev.WriteService
		t.Logf("%s: transfers %d, read blocks %d, write blocks %d (%.1f MB), busy ticks %d (%.1f%% at 100 Hz), "+
			"read service %d, write service %d timebase ticks (avg write %.2f ms at 512 MHz)",
			disk.Name, transfers, disk.ReadBlocks-prev.ReadBlocks, writeBlocks, float64(writeBlocks*512)/1048576,
			busyTicks, float64(busyTicks)/elapsed.Seconds(), disk.ReadService-prev.ReadService, writeService,
			float64(writeService)/512000/float64(max(transfers, 1)))
	}
}

// TestAixLiveTimebase measures the frequency of the timebase register (the service times of perfstat_disk)
func TestAixLiveTimebase(t *testing.T) {
	first, err := liveTimebase()
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	time.Sleep(2 * time.Second)
	second, err := liveTimebase()
	if err != nil {
		t.Fatal(err)
	}
	frequency := float64(second-first) / time.Since(start).Seconds()
	t.Logf("timebase frequency %.0f Hz", frequency)
	if frequency < aixTimebaseFrequency*0.98 || frequency > aixTimebaseFrequency*1.02 {
		t.Errorf("timebase frequency %.0f Hz, expected %d Hz", frequency, aixTimebaseFrequency)
	}
}

func liveTimebase() (uint64, error) {
	buf, err := perfstatPartitionTotal()
	if err != nil {
		return 0, err
	}
	times := &aixCpuTimes{}
	if err := decodePartition(buf, times); err != nil {
		return 0, err
	}
	return times.Timebase, nil
}

// TestAixLiveChecks runs the checks of batch 2 and prints the AIX tools for comparison
func TestAixLiveChecks(t *testing.T) {
	processes, _ := liveRunCheck(t, &CheckProcess{}).([]*resultProcess)
	t.Logf("processes: %d", len(processes))
	for _, process := range processes {
		if process.Pid == 1 || process.Name == "sshd" || process.Name == "xntpd" || process.Status[0] != "Active" {
			liveJSON(t, "process", process)
		}
	}
	liveCommand(t, "sh", "-c", "ps -e | wc -l; ps -e -o pid,ppid,user,nice,pcpu,pmem,vsz,rssize,stat,etime,args | head -5")

	liveJSON(t, "users", liveRunCheck(t, &CheckUser{}))
	liveCommand(t, "who")

	liveJSON(t, "net_stats", liveRunCheck(t, &CheckNet{}))
	liveCommand(t, "ifconfig", "-a")

	netIo := &CheckNetIo{}
	liveRunCheck(t, netIo)
	time.Sleep(5 * time.Second)
	liveJSON(t, "net_io (5 seconds)", liveRunCheck(t, netIo))

	liveJSON(t, "disk_io", liveRunCheck(t, &CheckDiskIo{}))

	liveJSON(t, "ntp", liveRunCheck(t, &CheckNtp{}))
	liveCommand(t, "ntpq", "-c", "rv")

	liveJSON(t, "mounts", liveRunCheck(t, &CheckMounts{}))

	liveJSON(t, "network_ports", liveRunCheck(t, &CheckNetworkPorts{}))
	liveCommand(t, "sh", "-c", "netstat -an -f inet | grep -c LISTEN; no -o tcp_ephemeral_low -o tcp_ephemeral_high")
}
