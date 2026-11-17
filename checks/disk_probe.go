package checks

import (
	"sync"
	"time"
)

// diskUsageTimeout is the time the disk check waits for the file systems. statfs of a network file system blocks
// if the server is not reachable (e.g. a hard NFS mount), such a file system is reported as hanging.
const diskUsageTimeout = 5 * time.Second

// diskProbe runs blocking file system calls like statfs in goroutines, as they can not be canceled.
// A call which did not return yet is reused by the next check run, so a hanging file system blocks
// only one goroutine and not the check.
type diskProbe[T any] struct {
	mutex   sync.Mutex
	pending map[string]*diskProbeCall[T]
}

type diskProbeCall[T any] struct {
	done   chan struct{}
	result T
	err    error
}

// diskProbeFunc reads the value of a file system, e.g. statfs of a mountpoint
type diskProbeFunc[T any] func(key string) (T, error)

// start returns the running call of the key or starts a new one
func (p *diskProbe[T]) start(key string, read diskProbeFunc[T]) *diskProbeCall[T] {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.pending == nil {
		p.pending = map[string]*diskProbeCall[T]{}
	}
	if call, ok := p.pending[key]; ok {
		return call
	}
	call := &diskProbeCall[T]{done: make(chan struct{})}
	p.pending[key] = call
	go p.run(key, read, call)
	return call
}

func (p *diskProbe[T]) run(key string, read diskProbeFunc[T], call *diskProbeCall[T]) {
	call.result, call.err = read(key)
	p.mutex.Lock()
	delete(p.pending, key)
	p.mutex.Unlock()
	close(call.done)
}

// wait returns the result of the call or hanging = true, if the call did not return before the deadline
func (call *diskProbeCall[T]) wait(deadline time.Time) (result T, err error, hanging bool) {
	// A call which already returned is used, even if the deadline passed while waiting for another call
	select {
	case <-call.done:
		return call.result, call.err, false
	default:
	}

	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()
	select {
	case <-call.done:
		return call.result, call.err, false
	case <-timer.C:
		var empty T
		return empty, nil, true
	}
}

// readAll starts the calls of all keys at once and waits until the common deadline
func (p *diskProbe[T]) readAll(keys []string, read diskProbeFunc[T], timeout time.Duration) []diskProbeResult[T] {
	calls := make([]*diskProbeCall[T], len(keys))
	for i, key := range keys {
		calls[i] = p.start(key, read)
	}
	deadline := time.Now().Add(timeout)
	results := make([]diskProbeResult[T], len(keys))
	for i, call := range calls {
		results[i].Value, results[i].Err, results[i].Hanging = call.wait(deadline)
	}
	return results
}

type diskProbeResult[T any] struct {
	Value   T
	Err     error
	Hanging bool
}
