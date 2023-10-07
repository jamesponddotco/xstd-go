package xhttp_test

import (
	"errors"
	"net/http"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		req         *http.Request
		expectedIP  string
		expectedErr error
	}{
		{
			name:        "Nil Request",
			req:         nil,
			expectedIP:  "",
			expectedErr: xhttp.ErrMissingRequest,
		},
		{
			name:        "Missing X-Forwarded-For and RemoteAddr",
			req:         &http.Request{},
			expectedIP:  "",
			expectedErr: xhttp.ErrInvalidIP,
		},
		{
			name: "Valid RemoteAddr, Missing X-Forwarded-For",
			req: &http.Request{
				RemoteAddr: "192.168.1.1:8080",
			},
			expectedIP:  "192.168.1.1",
			expectedErr: nil,
		},
		{
			name: "Valid X-Forwarded-For",
			req: &http.Request{
				Header: http.Header{
					"X-Forwarded-For": []string{"192.168.1.1, 192.168.1.2"},
				},
			},
			expectedIP:  "192.168.1.2",
			expectedErr: nil,
		},
		{
			name: "Invalid IP in X-Forwarded-For",
			req: &http.Request{
				Header: http.Header{
					"X-Forwarded-For": []string{"192.168.1.2, invalid"},
				},
			},
			expectedIP:  "",
			expectedErr: xhttp.ErrInvalidIP,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ip, err := xhttp.ClientIP(tt.req)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if ip != tt.expectedIP {
				t.Fatalf("expected IP: %s, got: %s", tt.expectedIP, ip)
			}
		})
	}
}
