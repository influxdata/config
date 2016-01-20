package config

import (
	"strconv"
	"time"
)

// Duration is a time.Duration that can marshal and unmarshal itself to a valid TOML representation
type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

// UnmarshalTOML converts a given human-readable duration into a time.Duration
func (d *Duration) UnmarshalTOML(data []byte) error {
	// Ignore if there is no value set
	if len(data) == 0 {
		return nil
	}

	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(unquoted)
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

// MarshalTOML converts a time.Duration into a human-readable duration.
func (d Duration) MarshalTOML() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}
