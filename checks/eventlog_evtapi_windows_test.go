package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
)

func newEvtApiCheck(t *testing.T, age int64, maxRecords int64, logs ...string) *CheckWindowsEventLog {
	t.Helper()
	c := &CheckWindowsEventLog{}
	enabled, err := c.Configure(&config.Configuration{
		WindowsEventLog:           true,
		WindowsEventLogTypes:      logs,
		WindowsEventLogAge:        age,
		WindowsEventLogMethod:     "EvtApi",
		WindowsEventLogMaxRecords: maxRecords,
	})
	if err != nil || !enabled {
		t.Fatal("check should be enabled: ", err)
	}
	t.Cleanup(func() { closeEvtState(nil) })
	return c
}

func runEvtApiCheck(t *testing.T, c *CheckWindowsEventLog) (map[string][]*resultEvent, time.Duration) {
	t.Helper()
	start := time.Now()
	r, err := c.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	duration := time.Since(start)
	if _, err := json.Marshal(r); err != nil {
		t.Fatal(err)
	}
	return r.(map[string][]*resultEvent), duration
}

func TestEventLogEvtApiWindow(t *testing.T) {
	c := newEvtApiCheck(t, 86400, 1000, "Application", "System", "Does-Not-Exist")
	result, duration := runEvtApiCheck(t, c)
	t.Logf("first run: Application=%d System=%d records in %s", len(result["Application"]), len(result["System"]), duration)

	if _, ok := result["Does-Not-Exist"]; ok {
		t.Error("a missing event log must not be part of the result")
	}

	minTime := time.Now().Add(-86400 * time.Second).Unix()
	for _, channel := range []string{"Application", "System"} {
		records, ok := result[channel]
		if !ok {
			t.Fatalf("missing result for %s", channel)
		}
		for i, record := range records {
			if record.Index <= 0 || record.Source == "" || record.Message == "" || record.MachineName == "" {
				t.Errorf("%s: incomplete record: %+v", channel, record)
			}
			if record.TimeGenerated < minTime {
				t.Errorf("%s: record outside of the time window: %+v", channel, record)
			}
			if i > 0 && record.Index <= records[i-1].Index {
				t.Errorf("%s: records are not sorted by record id", channel)
			}
			switch record.EntryType {
			case entryTypeError, entryTypeWarning, entryTypeInformation, entryTypeAuditSuccess, entryTypeAuditFailure:
			default:
				t.Errorf("%s: invalid entry type: %+v", channel, record)
			}
		}
		if len(records) > 0 {
			js, _ := json.Marshal(records[len(records)-1])
			t.Logf("%s newest record: %s", channel, js)
		}
	}

	// Second run only queries new records
	_, duration = runEvtApiCheck(t, c)
	t.Logf("second run in %s", duration)
}

func TestEventLogEvtApiIncremental(t *testing.T) {
	const source = "oitc-agent-evtapi-test"
	c := newEvtApiCheck(t, 3600, 1000, "Application")
	result, _ := runEvtApiCheck(t, c)
	before := len(result["Application"])

	// Write events without registering the source (no admin rights required).
	// As the source has no message file, the message can not be formatted and the fallback message is used.
	writer, err := eventlog.Open(source)
	if err != nil {
		t.Fatal(err)
	}
	marker := fmt.Sprintf("evtapi test %d", time.Now().UnixNano())
	if err := writer.Error(1001, marker+" error"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Warning(1002, marker+" warning"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Info(1003, marker+" info"); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	result, _ = runEvtApiCheck(t, c)
	records := result["Application"]
	found := map[int64]*resultEvent{}
	for _, record := range records {
		if record.Source == source && strings.Contains(record.Message, marker) {
			found[record.EventID] = record
		}
	}
	if len(records) < before+3 {
		t.Errorf("expected at least %d records, got %d", before+3, len(records))
	}
	for eventID, entryType := range map[int64]int64{1001: entryTypeError, 1002: entryTypeWarning, 1003: entryTypeInformation} {
		record, ok := found[eventID]
		if !ok {
			t.Errorf("event %d not found", eventID)
			continue
		}
		if record.EntryType != entryType {
			t.Errorf("event %d: expected entry type %d: %+v", eventID, entryType, record)
		}
	}
	if record, ok := found[1001]; ok {
		t.Logf("written record: %+v", record)
	}

	// Each record is only reported once
	ids := map[int64]bool{}
	for _, record := range records {
		if ids[record.Index] {
			t.Errorf("duplicate record %d", record.Index)
		}
		ids[record.Index] = true
	}
}

func TestEventLogEvtApiMaxRecordsAndCleared(t *testing.T) {
	c := newEvtApiCheck(t, 86400*7, 5, "System")
	result, _ := runEvtApiCheck(t, c)
	records := result["System"]
	if len(records) != 5 {
		t.Fatalf("expected the 5 newest records, got %d", len(records))
	}

	// Simulate a cleared event log: the last record id is higher than the newest record
	evtStatesMutex.Lock()
	newest := evtStates["System"].lastRecordID
	evtStates["System"].lastRecordID = newest + 1000000
	evtStatesMutex.Unlock()

	result, _ = runEvtApiCheck(t, c)
	if len(result["System"]) != 5 {
		t.Errorf("expected the records of the time window after the log got cleared, got %d", len(result["System"]))
	}
	evtStatesMutex.Lock()
	if evtStates["System"].lastRecordID != newest {
		t.Errorf("expected last record id %d, got %d", newest, evtStates["System"].lastRecordID)
	}
	evtStatesMutex.Unlock()
}

// The records of the EvtApi method must match the records of the WMI method, so the check receiver works the same way
func TestEventLogEvtApiCompareWithWmi(t *testing.T) {
	wmiCheck := &CheckWindowsEventLog{}
	if _, err := wmiCheck.Configure(&config.Configuration{
		WindowsEventLog: true, WindowsEventLogTypes: []string{"Application", "System"}, WindowsEventLogAge: 3600, WindowsEventLogMethod: "WMI",
	}); err != nil {
		t.Fatal(err)
	}
	evtCheck := newEvtApiCheck(t, 3600, 10000, "Application", "System")

	wmiResult, err := wmiCheck.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	evtResult, _ := runEvtApiCheck(t, evtCheck)

	for _, channel := range []string{"Application", "System"} {
		wmiRecords := map[int64]*resultEvent{}
		for _, record := range wmiResult.(map[string][]*resultEvent)[channel] {
			wmiRecords[record.Index] = record
		}
		compared, sameMessage := 0, 0
		for _, record := range evtResult[channel] {
			wmiRecord, ok := wmiRecords[record.Index]
			if !ok {
				// Records which got written between the two queries
				continue
			}
			compared++
			// WMI reports records of the level LogAlways (e.g. Microsoft-Windows-Security-SPP 902/903) as type 0,
			// which the check receiver handles as information like the EvtApi method
			wmiEntryType := wmiRecord.EntryType
			if wmiEntryType == 0 {
				wmiEntryType = entryTypeInformation
			}
			if record.EventID != wmiRecord.EventID || record.EntryType != wmiEntryType || record.Source != wmiRecord.Source || record.TimeGenerated != wmiRecord.TimeGenerated {
				t.Errorf("%s record %d differs:\nevtapi: %+v\nwmi:    %+v", channel, record.Index, record, wmiRecord)
			}
			if strings.TrimSpace(record.Message) == strings.TrimSpace(wmiRecord.Message) {
				sameMessage++
			} else {
				t.Logf("%s record %d (%s %d) message differs:\nevtapi: %q\nwmi:    %q", channel, record.Index, record.Source, record.EventID, record.Message, wmiRecord.Message)
			}
		}
		t.Logf("%s: wmi=%d evtapi=%d compared=%d same message=%d", channel, len(wmiRecords), len(evtResult[channel]), compared, sameMessage)
		if compared == 0 && len(wmiRecords) > 0 {
			t.Errorf("%s: no records to compare", channel)
		}
	}
}

// Many check runs must not leak handles or memory
func TestEventLogEvtApiStress(t *testing.T) {
	c := newEvtApiCheck(t, 86400, 1000, "Application", "System")
	getHandleCount := func() uint32 {
		var count uint32
		proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetProcessHandleCount")
		_, _, _ = proc.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&count)))
		return count
	}
	// PROCESS_MEMORY_COUNTERS, PagefileUsage is the private memory of the process (including native allocations)
	type processMemoryCounters struct {
		cb                         uint32
		PageFaultCount             uint32
		PeakWorkingSetSize         uintptr
		WorkingSetSize             uintptr
		QuotaPeakPagedPoolUsage    uintptr
		QuotaPagedPoolUsage        uintptr
		QuotaPeakNonPagedPoolUsage uintptr
		QuotaNonPagedPoolUsage     uintptr
		PagefileUsage              uintptr
		PeakPagefileUsage          uintptr
	}
	getMemory := func() uint64 {
		var counters processMemoryCounters
		counters.cb = uint32(unsafe.Sizeof(counters))
		proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("K32GetProcessMemoryInfo")
		_, _, _ = proc.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&counters)), uintptr(counters.cb))
		return uint64(counters.PagefileUsage)
	}

	// Warm up: load the DLLs and publisher metadata
	for i := 0; i < 20; i++ {
		runEvtApiCheck(t, c)
		closeEvtState(nil)
	}
	runtime.GC()
	handlesBefore, memoryBefore := getHandleCount(), getMemory()

	// The duration depends on the size of the event log of the system, so the test stops after 30 seconds
	// (at least 300 runs), to not exceed the timeout of go test on systems with large event logs
	start := time.Now()
	records := 0
	runs := 0
	for i := 0; i < 1500; i++ {
		if i >= 300 && time.Since(start) > 30*time.Second {
			break
		}
		runs++
		if i%300 == 0 {
			runtime.GC()
			t.Logf("run %4d: handles %d, private memory %d MB", i, getHandleCount(), getMemory()/1024/1024)
		}
		if i%3 == 0 {
			// Read the whole time window again
			closeEvtState(nil)
		}
		result, _ := runEvtApiCheck(t, c)
		records += len(result["Application"]) + len(result["System"])
	}
	runtime.GC()
	handlesAfter, memoryAfter := getHandleCount(), getMemory()

	t.Logf("%d runs (%d records) in %s, handles: %d -> %d, private memory: %d MB -> %d MB",
		runs, records, time.Since(start), handlesBefore, handlesAfter, memoryBefore/1024/1024, memoryAfter/1024/1024)
	if handlesAfter > handlesBefore+50 {
		t.Errorf("handle leak: %d -> %d", handlesBefore, handlesAfter)
	}
	if memoryAfter > memoryBefore+50*1024*1024 {
		t.Errorf("memory leak: %d MB -> %d MB", memoryBefore/1024/1024, memoryAfter/1024/1024)
	}
}
