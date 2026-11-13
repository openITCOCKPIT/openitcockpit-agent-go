package checks

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// blockingDiskReader blocks for the key "hanging" until release is closed
type blockingDiskReader struct {
	release chan struct{}
	calls   atomic.Int32
}

func (r *blockingDiskReader) read(key string) (string, error) {
	r.calls.Add(1)
	switch key {
	case "hanging":
		<-r.release
		return "late", nil
	case "error":
		return "", errors.New("permission denied")
	}
	return "usage of " + key, nil
}

func TestDiskProbe(t *testing.T) {
	reader := &blockingDiskReader{release: make(chan struct{})}
	probe := &diskProbe[string]{}

	start := time.Now()
	results := probe.readAll([]string{"/", "hanging", "error"}, reader.read, 200*time.Millisecond)
	if elapsed := time.Since(start); elapsed < 200*time.Millisecond || elapsed > 2*time.Second {
		t.Errorf("unexpected duration %s", elapsed)
	}
	if results[0].Value != "usage of /" || results[0].Hanging || results[0].Err != nil {
		t.Errorf("unexpected result %+v", results[0])
	}
	if !results[1].Hanging || results[1].Value != "" {
		t.Errorf("expected a hanging call %+v", results[1])
	}
	if results[2].Err == nil || results[2].Hanging {
		t.Errorf("expected an error %+v", results[2])
	}

	// The next run reuses the hanging call instead of starting a new one
	results = probe.readAll([]string{"hanging"}, reader.read, 50*time.Millisecond)
	if !results[0].Hanging || reader.calls.Load() != 3 {
		t.Errorf("hanging call was not reused: %+v, %d calls", results[0], reader.calls.Load())
	}

	// After the file system responds again, the late result is used and the call is removed
	close(reader.release)
	results = probe.readAll([]string{"hanging"}, reader.read, time.Second)
	if results[0].Hanging || results[0].Value != "late" {
		t.Errorf("unexpected result %+v", results[0])
	}
	time.Sleep(10 * time.Millisecond)
	results = probe.readAll([]string{"hanging"}, reader.read, time.Second)
	if results[0].Value != "late" || reader.calls.Load() != 4 {
		t.Errorf("expected a new call: %+v, %d calls", results[0], reader.calls.Load())
	}
}
