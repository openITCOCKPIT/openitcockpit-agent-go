package pushclient

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// testCertificate is a certificate with its key
type testCertificate struct {
	certificate *x509.Certificate
	key         *ecdsa.PrivateKey
	pem         []byte
}

// newTestCertificate creates a CA or a server certificate for 127.0.0.1 and localhost, signed by parent (nil = self
// signed). Server certificates get the given DNS names instead, if set.
func newTestCertificate(t *testing.T, name string, isCA bool, parent *testCertificate, dnsNames ...string) *testCertificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: name},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	if isCA {
		template.IsCA = true
		template.BasicConstraintsValid = true
		template.KeyUsage |= x509.KeyUsageCertSign
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		if len(dnsNames) > 0 {
			template.DNSNames = dnsNames
		} else {
			template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
			template.DNSNames = []string{"localhost"}
		}
	}

	signer, signerKey := template, key
	if parent != nil {
		signer, signerKey = parent.certificate, parent.key
	}
	der, err := x509.CreateCertificate(rand.Reader, template, signer, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &testCertificate{
		certificate: certificate,
		key:         key,
		pem:         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
	}
}

// newTestTLSServer starts a HTTPS server with the given certificate on 127.0.0.1
func newTestTLSServer(t *testing.T, server *testCertificate) *httptest.Server {
	t.Helper()
	ts := httptest.NewUnstartedServer(&fakePushServer{})
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{server.certificate.Raw},
			PrivateKey:  server.key,
		}},
	}
	ts.StartTLS()
	t.Cleanup(ts.Close)
	return ts
}

func writeTestFile(t *testing.T, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ca.crt")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

// pushTLSRequest sends a request with the TLS configuration of the push client
func pushTLSRequest(configuration *config.PushConfiguration, serverURL string) error {
	p := &PushClient{}
	p.client.Transport = &http.Transport{TLSClientConfig: pushTLSConfig(configuration)}
	target, _ := url.Parse(serverURL + "/agentconnector/submit_checkdata.json")
	res := submitCheckDataResponse{}
	_, err := p.httpRequest(context.Background(), target, map[string]string{"agentuuid": "test", "checkdata": "data"}, &res)
	return err
}

func TestPushTLSServerCaFile(t *testing.T) {
	ca := newTestCertificate(t, "Customer Root CA", true, nil)
	server := newTestTLSServer(t, newTestCertificate(t, "openitcockpit", false, ca))
	otherCa := newTestCertificate(t, "Other CA", true, nil)

	tests := []struct {
		name          string
		configuration *config.PushConfiguration
		expectedError string
	}{
		{"verification disabled", &config.PushConfiguration{}, ""},
		{"verification disabled, ca file is ignored", &config.PushConfiguration{ServerCaFile: writeTestFile(t, otherCa.pem)}, ""},
		{"own CA not installed", &config.PushConfiguration{VerifyServerCertificate: true}, "certificate"},
		{"own CA", &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, ca.pem)}, ""},
		{"own CA in a file with several certificates", &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, append(otherCa.pem, ca.pem...))}, ""},
		{"other CA", &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, otherCa.pem)}, "certificate"},
		// The push fails instead of trusting any server
		{"missing ca file", &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: filepath.Join(t.TempDir(), "missing.crt")}, "certificate"},
		{"invalid ca file", &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, []byte("no certificate"))}, "certificate"},
	}
	for _, test := range tests {
		err := pushTLSRequest(test.configuration, server.URL)
		if test.expectedError == "" && err != nil {
			t.Errorf("%s: unexpected error: %s", test.name, err)
		}
		if test.expectedError != "" && (err == nil || !strings.Contains(err.Error(), test.expectedError)) {
			t.Errorf("%s: expected an error with %q, got %v", test.name, test.expectedError, err)
		}
	}
}

func TestPushTLSSelfSignedServerCertificate(t *testing.T) {
	server := newTestCertificate(t, "openitcockpit", false, nil)
	ts := newTestTLSServer(t, server)

	configuration := &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, server.pem)}
	if err := pushTLSRequest(configuration, ts.URL); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// The address of the url must be in the certificate
	otherServer := newTestCertificate(t, "openitcockpit", false, nil, "openitcockpit.example.org")
	otherTs := newTestTLSServer(t, otherServer)
	configuration = &config.PushConfiguration{VerifyServerCertificate: true, ServerCaFile: writeTestFile(t, otherServer.pem)}
	err := pushTLSRequest(configuration, otherTs.URL)
	if err == nil || !strings.Contains(err.Error(), "cannot validate certificate for 127.0.0.1") {
		t.Errorf("expected an error for the address, got %v", err)
	}
}
