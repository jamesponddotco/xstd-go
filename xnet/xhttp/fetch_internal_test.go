package xhttp

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestSafeFetcher_IsDomainTrusted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		giveTrustedDomains []string
		giveURL            string
		want               bool
	}{
		{
			name:               "Empty trusted domains",
			giveTrustedDomains: []string{},
			giveURL:            "https://example.com",
			want:               false,
		},
		{
			name: "URL not in trusted domains",
			giveTrustedDomains: []string{
				"trusted1.com",
				"trusted2.com",
			},
			giveURL: "https://example.com",
			want:    false,
		},
		{
			name: "Empty URL",
			giveTrustedDomains: []string{
				"example.com",
			},
			giveURL: "",
			want:    false,
		},
		{
			name: "URL in trusted domains",
			giveTrustedDomains: []string{
				"trusted.com",
				"example.com",
			},
			giveURL: "https://example.com",
			want:    true,
		},
		{
			name: "Subdomain not in trusted domains",
			giveTrustedDomains: []string{
				"example.com",
			},
			giveURL: "https://subdomain.example.com",
			want:    false,
		},
		{
			name: "Invalid URL",
			giveTrustedDomains: []string{
				"example.com",
			},
			giveURL: "://invalid-url",
			want:    false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sf := &SafeFetcher{
				TrustedDomain: tt.giveTrustedDomains,
				Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
			}

			got := sf.isDomainTrusted(context.Background(), tt.giveURL)
			if got != tt.want {
				t.Errorf("IsDomainTrusted(%q) got %t, want %t", tt.giveURL, got, tt.want)
			}
		})
	}
}
