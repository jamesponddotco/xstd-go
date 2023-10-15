package xmiddleware

import (
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// PanicRecovery tries to recover from panics and returns a 500 error if there
// was one.
func PanicRecovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.LogAttrs(
					r.Context(),
					slog.LevelError,
					"panic recovered",
					slog.Any("error", err),
				)

				response := xhttp.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error. Please try again later.",
				}

				response.Write(r.Context(), logger, w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
