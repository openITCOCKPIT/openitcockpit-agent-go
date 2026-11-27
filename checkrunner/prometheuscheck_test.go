package checkrunner

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
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

// method = https has to be used for the request and the certificate has to be verified
func TestPrometheusCheckExporterHTTPS(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	var gotTLS, gotRequest atomic.Bool
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequest.Store(true)
	}))
	_ = ts.Listener.Close()
	ts.Listener = l
	ts.TLS = &tls.Config{
		GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
			gotTLS.Store(true)
			return nil, nil
		},
	}
	ts.StartTLS()
	defer ts.Close()

	executor := newTestPrometheusExecutor(l.Addr().(*net.TCPAddr).Port)
	executor.Configuration.Method = "https"
	executor.runCheck(context.Background(), 5*time.Second)

	if !gotTLS.Load() {
		t.Fatal("exporter with method = https was not queried via TLS")
	}

	// The test server uses a self-signed certificate
	if gotRequest.Load() {
		t.Fatal("request was sent to a server with an untrusted certificate")
	}
	select {
	case res := <-executor.ResultOutput:
		t.Fatalf("unexpected result from exporter with an untrusted certificate: %+v", res)
	default:
	}
}

func TestPrometheusCheckTimeoutGreaterInterval(t *testing.T) {
	executor := newTestPrometheusExecutor(9100)
	executor.Configuration.Interval = 5
	executor.Configuration.Timeout = 10

	if err := executor.Start(context.Background()); err == nil {
		t.Fatal("expected an error if timeout is greater than interval")
	}
}
