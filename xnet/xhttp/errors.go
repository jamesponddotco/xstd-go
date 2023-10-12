package xhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// ResponseError represents the response returned by the HTTP server when an
// error occurs.
type ResponseError struct {
	// Message is a human-readable message describing the error.
	Message string `json:"message"`

	// Documentation is a link to documentation describing the error.
	Documentation string `json:"documentation,omitempty"`

	// Code is a machine-readable code describing the error.
	Code int `json:"code"`
}

// Error implements the Error interface.
func (e ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Write serializes the ResponseError as a JSON object and writes it to the
// given HTTP response writer.
func (e ResponseError) Write(ctx context.Context, logger *slog.Logger, w http.ResponseWriter) {
	js, _ := json.MarshalIndent(e, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	if _, err := w.Write(js); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to write error response",
			slog.Any("error", err),
		)

		return
	}

	logger.LogAttrs(
		ctx,
		slog.LevelError,
		"error",
		slog.Int("code", e.Code),
		slog.String("message", e.Message),
	)
}
