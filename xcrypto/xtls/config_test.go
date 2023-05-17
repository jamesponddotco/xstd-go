package xtls_test

import (
	"crypto/tls"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	config := xtls.DefaultConfig()

	if config.SessionTicketsDisabled {
		t.Error("SessionTicketsDisabled is true")
	}

	if config.MinVersion != tls.VersionTLS12 {
		t.Error("MinVersion is not TLS 1.2")
	}

	if config.CipherSuites == nil {
		t.Error("CipherSuites is nil")
	}

	if config.CurvePreferences == nil {
		t.Error("CurvePreferences is nil")
	}

	if config.ClientSessionCache == nil {
		t.Error("ClientSessionCache is nil")
	}
}

func TestDefaultServerConfig(t *testing.T) {
	t.Parallel()

	config := xtls.DefaultServerConfig()

	if config.SessionTicketsDisabled {
		t.Error("SessionTicketsDisabled is true")
	}

	if config.MinVersion != tls.VersionTLS13 {
		t.Error("MinVersion is not TLS 1.2")
	}

	if config.ClientSessionCache == nil {
		t.Error("ClientSessionCache is nil")
	}
}

func TestDefaultCipherSuites(t *testing.T) {
	t.Parallel()

	cipherSuites := xtls.DefaultCipherSuites()

	if cipherSuites == nil {
		t.Error("cipherSuites is nil")
	}
}
