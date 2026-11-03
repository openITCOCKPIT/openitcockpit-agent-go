package webserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	log "github.com/sirupsen/logrus"
)

const (
	testBasicAuthUser     = "user"
	testBasicAuthPassword = "test"
	testBasicAuth         = "user:test"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestWebserverHandler(t *testing.T) {
	stateInput := make(chan []byte)
	w := &handler{
		StateInput: stateInput,
	}
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	stateInput <- []byte(`{"test": "tata"}`)
	cancel()
	w.wg.Wait()
}

func TestWebserverHandlerState(t *testing.T) {
	stateInput := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := &handler{
		StateInput: stateInput,
		Configuration: &config.Configuration{
			BasicAuth: "",
		},
	}
	ts := httptest.NewServer(w.Handler())
	defer ts.Close()

	w.Start(ctx)

	testState := []byte(`{"test": "tata"}`)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.SetBasicAuth(testBasicAuthUser, testBasicAuthPassword)

	r, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if string(body) != "{}" {
		t.Fatal("body does not match")
	}

	stateInput <- testState

	r, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err = io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if string(body) != string(testState) {
		t.Fatal("body does not match")
	}

	w.Shutdown()
}

func TestWebserverHandlerAuthFailed(t *testing.T) {
	stateInput := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := &handler{
		StateInput: stateInput,
		Configuration: &config.Configuration{
			BasicAuth: testBasicAuth,
		},
	}
	w.Start(ctx)

	ts := httptest.NewServer(w.Handler())
	defer ts.Close()

	stateInput <- []byte(`{"test": "tata"}`)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.SetBasicAuth("some", "fake")

	r, err := client.Do(req)
	if err != nil {
		t.Fatal("Request failed")
	}
	if r.StatusCode != http.StatusForbidden {
		t.Error("Unexpected status code: ", http.StatusForbidden)
	}
	_ = r.Body.Close()

	req, _ = http.NewRequest("GET", ts.URL, nil)
	r, err = client.Do(req)
	if err != nil {
		t.Fatal("Request failed")
	}
	if r.StatusCode != http.StatusForbidden {
		t.Error("Unexpected status code: ", http.StatusForbidden)
	}
	_ = r.Body.Close()

	w.Shutdown()
}

func TestWebserverHandlerConfig(t *testing.T) {
	state := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmpdir, err := os.MkdirTemp(os.TempDir(), "*-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpdir)
	}()
	cfgPath := filepath.Join(tmpdir, "config.ini")
	exporterPath := filepath.Join(tmpdir, "prometheus_exporters.ini")

	w := &handler{
		StateInput: state,
		Configuration: &config.Configuration{
			ConfigurationPath: cfgPath,
			ConfigUpdate:      true,
			Prometheus: &config.PrometheusConfiguration{
				ExportersFilePath: exporterPath,
				Enable:            false,
			},
		},
	}

	ts := httptest.NewServer(w.Handler())
	defer ts.Close()

	result := ""
	w.Start(ctx)

	data, err := json.Marshal(&configurationPush{
		Configuration:                   base64.StdEncoding.EncodeToString([]byte(`[default]`)),
		CustomCheckConfiguration:        "",
		PrometheusExporterConfiguration: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(ts.URL+"/config", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	d, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Error(err)
	} else {
		result = string(d)
		if result != `[default]` {
			t.Error("unexpected result ([default]): ", result)
		}
	}

	resp, err = http.Get(ts.URL + "/config")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Status code is not 200")
	}

	cp := &configurationPush{}
	if err := json.Unmarshal(body, cp); err != nil {
		t.Fatal(err)
	}
	cfg, err := base64.StdEncoding.DecodeString(cp.Configuration)
	if err != nil {
		t.Fatal(err)
	}
	if string(cfg) != "[default]" {
		t.Fatal("unexpected response for configuration get")
	}

	w.Shutdown()
}

func TestWebserverHandlerAutoTLSCreatesDirectories(t *testing.T) {
	// Every file is stored in its own directory, which does not exist yet
	dir := t.TempDir()
	cfg := &config.Configuration{
		AutoSslKeyFile: filepath.Join(dir, "private", "agent.key"),
		AutoSslCsrFile: filepath.Join(dir, "csr", "agent.csr"),
		AutoSslCrtFile: filepath.Join(dir, "certs", "agent.crt"),
		AutoSslCaFile:  filepath.Join(dir, "ca", "server_ca.crt"),
	}
	w := &handler{
		StateInput:    make(chan []byte),
		Configuration: cfg,
	}
	ts := httptest.NewServer(w.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/autotls?domain=agent.example.org")
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code of the csr request: %d", res.StatusCode)
	}
	for _, file := range []string{cfg.AutoSslKeyFile, cfg.AutoSslCsrFile} {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("file was not created: %s", err)
		}
	}

	body, _ := json.Marshal(updateCrtRequest{Signed: "certificate", CA: "ca"})
	res, err = http.Post(ts.URL+"/autotls", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code of the certificate update: %d", res.StatusCode)
	}
	for file, content := range map[string]string{cfg.AutoSslCrtFile: "certificate", cfg.AutoSslCaFile: "ca"} {
		if data, err := os.ReadFile(file); err != nil || string(data) != content {
			t.Errorf("unexpected content of %s: %q %v", file, data, err)
		}
	}
}

func TestWebserverHandlerGzip(t *testing.T) {
	stateInput := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := &handler{
		StateInput:    stateInput,
		Configuration: &config.Configuration{},
	}
	ts := httptest.NewServer(w.Handler())
	defer ts.Close()
	w.Start(ctx)

	state := bytes.Repeat([]byte(`{"name": "value"},`), 200)
	stateInput <- state

	// The transport does not add Accept-Encoding and does not decompress, so the raw response is visible
	client := &http.Client{Transport: &http.Transport{DisableCompression: true}}
	get := func(acceptEncoding string) (*http.Response, []byte) {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		if acceptEncoding != "" {
			req.Header.Set("Accept-Encoding", acceptEncoding)
		}
		r, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		return r, body
	}

	// Clients without Accept-Encoding (e.g. older receivers) get the uncompressed state
	r, body := get("")
	if r.Header.Get("Content-Encoding") != "" || !bytes.Equal(body, state) {
		t.Errorf("expected the uncompressed state, got %q with %d bytes", r.Header.Get("Content-Encoding"), len(body))
	}

	// Accept-Encoding: gzip
	r, body = get("gzip")
	if r.Header.Get("Content-Encoding") != "gzip" || r.Header.Get("Vary") != "Accept-Encoding" || len(body) >= len(state) {
		t.Fatalf("expected a compressed response, got %q with %d bytes", r.Header.Get("Content-Encoding"), len(body))
	}
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	decompressed, err := io.ReadAll(reader)
	if err != nil || !bytes.Equal(decompressed, state) {
		t.Errorf("the decompressed state does not match (%v)", err)
	}

	// gzip is not acceptable
	r, _ = get("gzip;q=0")
	if r.Header.Get("Content-Encoding") != "" {
		t.Error("expected an uncompressed response for gzip;q=0")
	}

	w.Shutdown()
}

func TestBasicAuthCredentialsMatch(t *testing.T) {
	b := &basicAuthMiddleware{Username: testBasicAuthUser, Password: testBasicAuthPassword}
	tests := []struct {
		user     string
		password string
		expected bool
	}{
		{testBasicAuthUser, testBasicAuthPassword, true},
		{testBasicAuthUser, "tes", false},
		{testBasicAuthUser, "test2", false},
		{testBasicAuthUser, "", false},
		{"use", testBasicAuthPassword, false},
		{"", "", false},
		{testBasicAuthPassword, testBasicAuthUser, false},
	}
	for _, test := range tests {
		if got := b.credentialsMatch(test.user, test.password); got != test.expected {
			t.Errorf("%q:%q: expected %v, got %v", test.user, test.password, test.expected, got)
		}
	}
}
