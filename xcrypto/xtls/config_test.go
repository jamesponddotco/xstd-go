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

func TestModernServerConfig(t *testing.T) {
	t.Parallel()

	cfg := xtls.ModernServerConfig()

	if cfg.MinVersion != tls.VersionTLS13 {
		t.Errorf("Unexpected MinVersion: got %v want %v", cfg.MinVersion, tls.VersionTLS13)
	}

	if len(cfg.CurvePreferences) != 3 || cfg.CurvePreferences[0] != tls.X25519 || cfg.CurvePreferences[1] != tls.CurveP256 || cfg.CurvePreferences[2] != tls.CurveP384 {
		t.Errorf("Unexpected CurvePreferences: got %v want [tls.X25519, tls.CurveP256, tls.CurveP384]", cfg.CurvePreferences)
	}

	if cfg.SessionTicketsDisabled {
		t.Errorf("Unexpected SessionTicketsDisabled: got %v want false", cfg.SessionTicketsDisabled)
	}

	if cfg.ClientSessionCache == nil {
		t.Errorf("ClientSessionCache was not initialized")
	}
}

func TestIntermediateServerConfig(t *testing.T) {
	t.Parallel()

	cfg := xtls.IntermediateServerConfig()

	if cfg.MinVersion != tls.VersionTLS12 {
		t.Errorf("Unexpected MinVersion: got %v want %v", cfg.MinVersion, tls.VersionTLS12)
	}

	if len(cfg.CurvePreferences) != 3 || cfg.CurvePreferences[0] != tls.X25519 || cfg.CurvePreferences[1] != tls.CurveP256 || cfg.CurvePreferences[2] != tls.CurveP384 {
		t.Errorf("Unexpected CurvePreferences: got %v want [tls.X25519, tls.CurveP256, tls.CurveP384]", cfg.CurvePreferences)
	}

	if cfg.SessionTicketsDisabled {
		t.Errorf("Unexpected SessionTicketsDisabled: got %v want false", cfg.SessionTicketsDisabled)
	}

	if cfg.ClientSessionCache == nil {
		t.Errorf("ClientSessionCache was not initialized")
	}

	if len(cfg.CipherSuites) == 0 {
		t.Errorf("CipherSuites was not set")
	}
}

func TestDefaultCipherSuites(t *testing.T) {
	t.Parallel()

	cipherSuites := xtls.DefaultCipherSuites()

	if cipherSuites == nil {
		t.Error("cipherSuites is nil")
	}
}
