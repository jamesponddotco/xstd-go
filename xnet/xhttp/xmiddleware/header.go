package xmiddleware

import (
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// UserAgent ensures that the request has the User-Agent header set.
func UserAgent(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() == "" {
			response := xhttp.ResponseError{
				Code:    http.StatusBadRequest,
				Message: "User agent is missing. Please provide a valid user agent.",
			}

			response.Write(r.Context(), logger, w)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// PrivacyPolicy adds a privacy policy header to the response.
func PrivacyPolicy(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Privacy-Policy", uri)

		next.ServeHTTP(w, r)
	})
}

// TermsOfService adds a terms of service header to the response.
func TermsOfService(uri string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Terms-Of-Service", uri)

		next.ServeHTTP(w, r)
	})
}
