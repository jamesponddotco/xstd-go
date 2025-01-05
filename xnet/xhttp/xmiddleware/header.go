package xmiddleware

import (
	"context"
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xcontext"
	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// contextKeyRequestID is the key used for the request ID in the context.
const contextKeyRequestID xcontext.ContextKey = "request_id"

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

// RequestID adds a request ID header to the response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(xhttp.RequestID)
		if requestID == "" {
			var uuid xrand.UUID

			requestID = uuid.GenerateV4()
		}

		ctx := context.WithValue(r.Context(), contextKeyRequestID, requestID)

		r = r.WithContext(ctx)

		w.Header().Set(xhttp.RequestID, requestID)

		next.ServeHTTP(w, r)
	})
}
