package webserver

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

// certificateUpdateMessage describes a certificate update of AutoSSL for the log: the sender, the new CA and whether
// it replaces the CA of another openITCOCKPIT server, e.g. if another server took over the agent. The second return
// value is true if the CA of another server is replaced.
func certificateUpdateMessage(remoteAddr string, newCaPem string, existingCaFile string) (string, bool) {
	newCa, err := describeCaCertificate([]byte(newCaPem))
	if err != nil {
		return fmt.Sprintf("AutoSSL certificate update from %s, the CA can not be read: %s", remoteAddr, err), false
	}

	existingData, err := os.ReadFile(existingCaFile)
	if err != nil {
		return fmt.Sprintf("AutoSSL certificate update from %s, new CA %s", remoteAddr, newCa), false
	}
	existingCa, err := describeCaCertificate(existingData)
	if err != nil {
		return fmt.Sprintf("AutoSSL certificate update from %s, new CA %s (the previous CA can not be read: %s)", remoteAddr, newCa, err), false
	}
	if existingCa == newCa {
		return fmt.Sprintf("AutoSSL certificate update from %s, same CA %s", remoteAddr, newCa), false
	}
	return fmt.Sprintf("AutoSSL certificate update from %s replaces the CA %s with %s, only the openITCOCKPIT server of the new CA can query this agent now",
		remoteAddr, existingCa, newCa), true
}

// describeCaCertificate returns the subject and the SHA-256 fingerprint of the first certificate of a PEM file, the
// fingerprint has the format of "openssl x509 -noout -fingerprint -sha256"
func describeCaCertificate(data []byte) (string, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("no PEM encoded certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%q (SHA-256 %s)", certificate.Subject.String(), formatFingerprint(sha256.Sum256(certificate.Raw))), nil
}

func formatFingerprint(sum [sha256.Size]byte) string {
	parts := make([]string, len(sum))
	for i, b := range sum {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}
