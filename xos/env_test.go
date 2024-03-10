package xos_test

import (
	"os"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xos"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envVal   string
		want     string
	}{
		{
			name:     "non-existent variable with a fallback",
			key:      "SOMETHING_THAT_DOES_NOT_EXIST",
			fallback: "fallback",
			envVal:   "",
			want:     "fallback",
		},
		{
			name:     "existent variable with a value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: "fallback",
			envVal:   "test-value",
			want:     "test-value",
		},
		{
			name:     "existent variable with an empty value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: "fallback",
			envVal:   "",
			want:     "fallback",
		},
	}

	for _, tt := range tests { //nolint:paralleltest // Test is not safe to run in parallel.
		t.Run(tt.name, func(t *testing.T) {
			// Set up the environment for the test case
			t.Setenv(tt.key, tt.envVal)
			defer os.Unsetenv(tt.key)

			got := xos.GetEnv(tt.key, tt.fallback)

			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback int
		envVal   string
		want     int
	}{
		{
			name:     "non-existent variable with a fallback",
			key:      "SOMETHING_THAT_DOES_NOT_EXIST",
			fallback: 123,
			envVal:   "",
			want:     123,
		},
		{
			name:     "existent variable with a valid integer value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: 123,
			envVal:   "456",
			want:     456,
		},
		{
			name:     "existent variable with an invalid integer value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: 123,
			envVal:   "not-an-integer",
			want:     123,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // Test is not safe to run in parallel.
		t.Run(tt.name, func(t *testing.T) {
			// Set up the environment for the test case
			t.Setenv(tt.key, tt.envVal)
			defer os.Unsetenv(tt.key)

			got := xos.GetIntEnv(tt.key, tt.fallback)

			if got != tt.want {
				t.Errorf("GetIntEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback bool
		envVal   string
		want     bool
	}{
		{
			name:     "non-existent variable with a fallback",
			key:      "SOMETHING_THAT_DOES_NOT_EXIST",
			fallback: true,
			envVal:   "",
			want:     true,
		},
		{
			name:     "existent variable with a valid boolean value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: true,
			envVal:   "false",
			want:     false,
		},
		{
			name:     "existent variable with an invalid boolean value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: true,
			envVal:   "not-a-boolean",
			want:     true,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // Test is not safe to run in parallel.
		t.Run(tt.name, func(t *testing.T) {
			// Set up the environment for the test case
			t.Setenv(tt.key, tt.envVal)
			defer os.Unsetenv(tt.key)

			got := xos.GetBoolEnv(tt.key, tt.fallback)

			if got != tt.want {
				t.Errorf("GetBoolEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback time.Duration
		envVal   string
		want     time.Duration
	}{
		{
			name:     "non-existent variable with a fallback",
			key:      "SOMETHING_THAT_DOES_NOT_EXIST",
			fallback: 5 * time.Minute,
			envVal:   "",
			want:     5 * time.Minute,
		},
		{
			name:     "existent variable with a valid duration value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: 5 * time.Minute,
			envVal:   "1h30m",
			want:     1*time.Hour + 30*time.Minute,
		},
		{
			name:     "existent variable with an invalid duration value",
			key:      "SOME_EXISTENT_VARIABLE",
			fallback: 5 * time.Minute,
			envVal:   "not-a-duration",
			want:     5 * time.Minute,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // Test is not safe to run in parallel.
		t.Run(tt.name, func(t *testing.T) {
			// Set up the environment for the test case
			t.Setenv(tt.key, tt.envVal)
			defer os.Unsetenv(tt.key)

			got := xos.GetDurationEnv(tt.key, tt.fallback)

			if got != tt.want {
				t.Errorf("GetDurationEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
