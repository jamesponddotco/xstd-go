package xnet

import (
	"net"
	"time"
)

const (
	// DefaultDialerTimeout is the default maximum amount of time a dial will wait for
	// a connection to complete.
	DefaultDialerTimeout = 30 * time.Second

	// DefaultDialerKeepAlive is the default amount of time a connection will be
	// kept alive.
	DefaultDialerKeepAlive = 30 * time.Second
)

// NewDialer returns a new Dialer with default timeout values.
//
// These values are not one-size-fits-all and should be adapted based on
// specific application needs and network characteristics.
func NewDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   DefaultDialerTimeout,
		KeepAlive: DefaultDialerKeepAlive,
	}
}
