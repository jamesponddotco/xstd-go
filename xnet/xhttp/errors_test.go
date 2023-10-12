package xhttp_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestResponseError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    xhttp.ResponseError
		expected string
	}{
		{
			name: "basic error",
			input: xhttp.ResponseError{
				Message: "Not Found",
				Code:    404,
			},
			expected: "404: Not Found",
		},
		{
			name: "empty message",
			input: xhttp.ResponseError{
				Message: "",
				Code:    500,
			},
			expected: "500: ",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := tt.input.Error()
			if output != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, output)
			}
		})
	}
}

func TestResponseError_Write(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          xhttp.ResponseError
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "basic error",
			input: xhttp.ResponseError{
				Message: "Not Found",
				Code:    404,
			},
			expectedBody: `{
  "message": "Not Found",
  "code": 404
}`,
			expectedStatus: 404,
		},
		{
			name: "empty message",
			input: xhttp.ResponseError{
				Message: "",
				Code:    500,
			},
			expectedBody: `{
  "message": "",
  "code": 500
}`,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				buf     bytes.Buffer
				w       = httptest.NewRecorder()
				handler = slog.NewJSONHandler(&buf, nil)
				logger  = slog.New(handler)
			)

			tt.input.Write(context.Background(), logger, w)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}
