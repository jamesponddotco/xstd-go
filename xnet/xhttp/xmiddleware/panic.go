package xmiddleware

import (
	"context"
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// PanicRecovery tries to recover from panics and returns a 500 error if there
// was one.
func PanicRecovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer func(ctx context.Context) {
			if err := recover(); err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"panic recovered",
					slog.Any("error", err),
				)

				response := xhttp.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error. Please try again later.",
				}

				response.Write(ctx, logger, w)
			}
		}(ctx)

		next.ServeHTTP(w, r)
	})
}
