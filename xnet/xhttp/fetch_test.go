package xhttp_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestSafeFetcher_Fetch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		client           *http.Client
		header           http.Header
		trustedDomain    []string
		maxContentLength int64
		uri              string
		wantErr          error
	}{
		{
			name:   "Valid Request",
			client: http.DefaultClient,
			header: http.Header{
				"User-Agent": []string{"test"},
			},
			trustedDomain:    []string{"example.com"},
			maxContentLength: 1024,
			uri:              "https://example.com/",
		},
		{
			name:          "Nil HTTP Client",
			trustedDomain: []string{"example.com"},
			uri:           "https://example.com/nil-client",
			wantErr:       xhttp.ErrNilHTTPClient,
		},
		{
			name:          "Non-HTTPS URI",
			client:        http.DefaultClient,
			trustedDomain: []string{"example.com"},
			uri:           "http://example.com/non-https",
			wantErr:       xhttp.ErrNotHTTPS,
		},
		{
			name:          "Untrusted Domain",
			client:        http.DefaultClient,
			trustedDomain: []string{"example.com"},
			uri:           "https://untrusted.com",
			wantErr:       xhttp.ErrDomainNotTrusted,
		},
		{
			name:             "Response Status Not OK",
			client:           http.DefaultClient,
			trustedDomain:    []string{"example.com"},
			maxContentLength: 1024,
			uri:              "https://example.com/not-found",
			wantErr:          xhttp.ErrRequestFailed,
		},
		{
			name:             "Response Exceeds Max Content Length",
			client:           http.DefaultClient,
			trustedDomain:    []string{"i.cpimg.sh"},
			maxContentLength: 1,
			uri:              "https://i.cpimg.sh/r6BsGDijfFyl.png",
			wantErr:          xhttp.ErrExceededMaxContentLength,
		},
		{
			name:   "Response Successful With Content Length",
			client: http.DefaultClient,
			trustedDomain: []string{
				"i.cpimg.sh",
			},
			maxContentLength: 420000,
			uri:              "https://i.cpimg.sh/r6BsGDijfFyl.png",
			wantErr:          nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sf := &xhttp.SafeFetcher{
				Client:           tt.client,
				Header:           tt.header,
				TrustedDomain:    tt.trustedDomain,
				MaxContentLength: tt.maxContentLength,
			}

			_, err := sf.Fetch(context.Background(), tt.uri)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSafeFetcher_Fetch_ContextCancellation(t *testing.T) {
	t.Parallel()

	sf := &xhttp.SafeFetcher{
		Client:           &http.Client{},
		TrustedDomain:    []string{"example.com"},
		MaxContentLength: 1024,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := sf.Fetch(ctx, "https://example.com")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Fetch() error = %v, wantErr %v", err, context.Canceled)
	}
}
