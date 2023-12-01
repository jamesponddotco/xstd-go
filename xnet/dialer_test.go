package xnet_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet"
)

func TestNewDialer(t *testing.T) {
	t.Parallel()

	dialer := xnet.NewDialer()

	if dialer.Timeout != xnet.DefaultDialerTimeout {
		t.Errorf("expected Timeout to be %v, got %v", xnet.DefaultDialerTimeout, dialer.Timeout)
	}

	if dialer.KeepAlive != xnet.DefaultDialerKeepAlive {
		t.Errorf("expected KeepAlive to be %v, got %v", xnet.DefaultDialerKeepAlive, dialer.KeepAlive)
	}
}
