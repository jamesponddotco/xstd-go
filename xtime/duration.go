package xtime

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a wrapper for time.Duration which supports JSON unmarshalling
// from a string.
type Duration time.Duration

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("could not unmarshal duration string: %w", err)
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("could not parse duration string: %w", err)
	}

	*d = Duration(duration)

	return nil
}
