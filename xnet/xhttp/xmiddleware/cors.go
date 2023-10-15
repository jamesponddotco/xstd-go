package xmiddleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
)

// DefaultCORSMaxAge is the default amount of time (in seconds) that a browser
// can cache the preflight response.
const DefaultCORSMaxAge uint = 600

// CORSConfig represents the basic configuration for the CORS middleware.
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to make requests to
	// the service.
	AllowedOrigins []string

	// AllowedMethods is a list of methods that are allowed to make requests to
	// the service.
	AllowedMethods []string

	// AllowedHeaders is a list of headers that can be used when making requests
	// to the service.
	AllowedHeaders []string

	// ExposedHeaders is a list of headers that are exposed to the client.
	ExposedHeaders []string

	// MaxAge is the maximum amount of time (in seconds) that a browser can
	// cache the preflight response.
	MaxAge uint

	// AllowCredentials indicates whether or not the request can include user
	// credentials.
	AllowCredentials bool
}

// DefaultCORSConfig returns the default configuration for the CORS middleware.
// The default configuration is fairly opinionated, read-only, set of options.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{xhttp.Wildcard},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			xhttp.Accept,
			xhttp.ContentType,
			xhttp.ContentLength,
			xhttp.AcceptEncoding,
		},
		ExposedHeaders:   []string{},
		MaxAge:           DefaultCORSMaxAge,
		AllowCredentials: false,
	}
}

// CORS adds CORS headers to the response given the provided configuration
// options. If no options are provided, DefaultCORSConfig() is used.
func CORS(config *CORSConfig, logger *slog.Logger, next http.Handler) http.Handler {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			origin         = r.Header.Get(xhttp.Origin)
			allowedOrigins string
		)

		// Special case: if any allowed origin is "*", allow any origin.
		for _, o := range config.AllowedOrigins {
			if o == "*" {
				allowedOrigins = xhttp.Wildcard

				break
			}
		}

		if allowedOrigins == "" {
			for _, o := range config.AllowedOrigins {
				if o == origin {
					allowedOrigins = origin

					break
				}
			}
		}

		if allowedOrigins == "" {
			logger.LogAttrs(
				r.Context(),
				slog.LevelWarn,
				"blocked cors request",
				slog.String("origin", origin),
			)

			response := xhttp.ResponseError{
				Code:    http.StatusForbidden,
				Message: fmt.Sprintf("Origin %s is not allowed.", origin),
			}

			response.Write(r.Context(), logger, w)

			return
		}

		w.Header().Set(xhttp.AccessControlAllowOrigin, allowedOrigins)

		w.Header().Set(xhttp.AccessControlAllowMethods, xstrings.JoinWithSeparator(", ", config.AllowedMethods...))
		w.Header().Set(xhttp.AccessControlAllowHeaders, xstrings.JoinWithSeparator(", ", config.AllowedHeaders...))
		w.Header().Set(xhttp.AccessControlExposeHeaders, xstrings.JoinWithSeparator(", ", config.ExposedHeaders...))

		if config.AllowCredentials {
			if allowedOrigins != xhttp.Wildcard {
				w.Header().Set(xhttp.AccessControlAllowCredentials, xhttp.True)
			} else {
				logger.LogAttrs(
					r.Context(),
					slog.LevelWarn,
					"AllowCredentials is true, but origin is a wildcard.",
					slog.String("origin", origin),
				)
			}
		}

		if config.MaxAge > 0 {
			w.Header().Set(xhttp.AccessControlMaxAge, strconv.Itoa(int(config.MaxAge)))
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
