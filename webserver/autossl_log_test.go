package webserver

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

// newTestCaPem returns a self signed CA certificate like the one of openITCOCKPIT and its SHA-256 fingerprint
func newTestCaPem(t *testing.T, commonName string) (string, string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})), formatFingerprint(sha256.Sum256(der))
}

func TestCertificateUpdateMessage(t *testing.T) {
	ca, caFingerprint := newTestCaPem(t, "72e53fa3-09c9-40f4-bfa2-66191a63ce5e.agentserver.oitc")
	otherCa, otherFingerprint := newTestCaPem(t, "11111111-2222-3333-4444-555555555555.agentserver.oitc")
	dir := t.TempDir()
	caFile := filepath.Join(dir, "server_ca.crt")
	invalidFile := filepath.Join(dir, "invalid.crt")
	if err := os.WriteFile(caFile, []byte(ca), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(invalidFile, []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		newCa        string
		existingFile string
		replaces     bool
		contains     []string
	}{
		{"first certificate", ca, filepath.Join(dir, "missing.crt"), false,
			[]string{"from 10.0.8.10:51234, new CA", "72e53fa3-09c9-40f4-bfa2-66191a63ce5e.agentserver.oitc", caFingerprint}},
		{"renewal", ca, caFile, false, []string{"same CA", caFingerprint}},
		{"other server", otherCa, caFile, true,
			[]string{"replaces the CA", caFingerprint, otherFingerprint, "11111111-2222-3333-4444-555555555555.agentserver.oitc"}},
		{"previous CA invalid", ca, invalidFile, false, []string{"new CA", "the previous CA can not be read"}},
		{"new CA invalid", "invalid", caFile, false, []string{"the CA can not be read"}},
	}
	for _, test := range tests {
		message, replaces := certificateUpdateMessage("10.0.8.10:51234", test.newCa, test.existingFile)
		if replaces != test.replaces {
			t.Errorf("%s: expected replaces=%v: %s", test.name, test.replaces, message)
		}
		for _, text := range test.contains {
			if !strings.Contains(message, text) {
				t.Errorf("%s: %q not in %q", test.name, text, message)
			}
		}
	}
}

func TestWebserverHandlerAutoTLSLogsReplacedCa(t *testing.T) {
	hook := test.NewGlobal()
	defer hook.Reset()

	ca, _ := newTestCaPem(t, "first.agentserver.oitc")
	otherCa, _ := newTestCaPem(t, "second.agentserver.oitc")
	dir := t.TempDir()
	cfg := &config.Configuration{
		AutoSslCrtFile: filepath.Join(dir, "agent.crt"),
		AutoSslCaFile:  filepath.Join(dir, "server_ca.crt"),
	}
	if err := os.WriteFile(cfg.AutoSslCaFile, []byte(ca), 0600); err != nil {
		t.Fatal(err)
	}
	w := &handler{StateInput: make(chan []byte), Configuration: cfg}
	ts := httptest.NewServer(w.Handler())
	defer ts.Close()

	body, _ := json.Marshal(updateCrtRequest{Signed: "certificate", CA: otherCa})
	res, err := http.Post(ts.URL+"/autotls", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()

	found := false
	for _, entry := range hook.AllEntries() {
		if entry.Level == log.WarnLevel && strings.Contains(entry.Message, "replaces the CA") && strings.Contains(entry.Message, "127.0.0.1:") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a warning about the replaced CA: %v", hook.AllEntries())
	}
	if data, _ := os.ReadFile(cfg.AutoSslCaFile); string(data) != otherCa {
		t.Error("the new CA was not stored")
	}
}
