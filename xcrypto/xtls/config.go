package xtls

import (
	"crypto/tls"
)

// DefaultConfig returns a [*tls.Config] with optimized security and
// performance settings for common use cases.
//
// [*tls.Config]: https://godocs.io/crypto/tls#Config
func DefaultConfig() *tls.Config {
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		CipherSuites: DefaultCipherSuites(),
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		SessionTicketsDisabled: false,
		ClientSessionCache:     tls.NewLRUClientSessionCache(64),
	}
}

// DefaultCipherSuites returns a sensible default list of cipher suites based
// on [Mozilla's recommendations].
//
// [Mozilla's recommendations]: https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
func DefaultCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}
