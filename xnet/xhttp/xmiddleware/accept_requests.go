package xmiddleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

// AcceptRequests ensures that the request method is one of the allowed methods.
func AcceptRequests(methods []string, logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)

				return
			}
		}

		w.Header().Set("Allow", strings.Join(methods, ", "))

		response := xhttp.ResponseError{
			Code:    http.StatusMethodNotAllowed,
			Message: fmt.Sprintf("Method %s not allowed.", r.Method),
		}

		response.Write(r.Context(), logger, w)
	})
}
