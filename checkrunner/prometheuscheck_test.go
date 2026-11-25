package checkrunner

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

func newTestPrometheusExecutor(port int) *PrometheusCheckExecutor {
	return &PrometheusCheckExecutor{
		Configuration: &config.PrometheusExporter{
			Name: "test_exporter",
			Port: int64(port),
			Path: "/metrics",
		},
		ResultOutput: make(chan *PrometheusExporterResult, 1),
		shutdown:     make(chan struct{}),
	}
}

// An exporter that is not reachable must not crash the agent
func TestPrometheusCheckExporterNotReachable(t *testing.T) {
	// Reserve a free port and close it again, so nothing is listening on it
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}

	executor := newTestPrometheusExecutor(port)
	executor.runCheck(context.Background(), 5*time.Second)

	select {
	case res := <-executor.ResultOutput:
		t.Fatalf("unexpected result from unreachable exporter: %+v", res)
	default:
	}
}

func TestPrometheusCheckExporterReachable(t *testing.T) {
	const metrics = "test_metric 1\n"

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(metrics))
	}))
	_ = ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	executor := newTestPrometheusExecutor(l.Addr().(*net.TCPAddr).Port)
	executor.runCheck(context.Background(), 5*time.Second)

	select {
	case res := <-executor.ResultOutput:
		if res.Name != "test_exporter" || res.Result != metrics {
			t.Fatalf("unexpected result: %+v", res)
		}
	default:
		t.Fatal("no result from reachable exporter")
	}
}
