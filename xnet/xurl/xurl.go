// Package xurl provides helpers and utilities for working with URLs.
package xurl

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Normalize takes a raw URL string as input and returns a normalized URL
// string.
func Normalize(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	if (parsedURL.Scheme == "http" && parsedURL.Port() == "80") ||
		(parsedURL.Scheme == "https" && parsedURL.Port() == "443") {
		parsedURL.Host = strings.Split(parsedURL.Host, ":")[0]
	}

	if parsedURL.RawQuery == "" {
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	}

	query := parsedURL.Query()
	for key := range query {
		sort.Strings(query[key])
	}

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
