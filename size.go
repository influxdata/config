package config

import (
	"fmt"
	"strconv"
)

const maxInt = int64(^uint(0) >> 1)

// Size is an integer representation of a capacity of data written in a
// human-readable format. The format consists of digits followed by a suffix
// indicating the magnitude of the capacity.
type Size int

// UnmarshalTOML converts the provided human-readable size and converts it into
// an integer representation
func (s *Size) UnmarshalTOML(data []byte) error {
	//parse numeric portion
	length := len(string(data))
	size, err := strconv.ParseInt(string(data[:length-1]), 10, 64)
	if err != nil {
		return err
	}

	//parse unit of measure
	switch suffix := data[len(data)-1]; suffix {
	case 'm':
		size *= 1 << 20 // MB
	case 'g':
		size *= 1 << 30 // GB
	default:
		return fmt.Errorf("unknown size suffix: %c", suffix)
	}

	if size > maxInt {
		return fmt.Errorf("size %d cannot be represented by an int", size)
	}

	*s = Size(size)
	return nil
}
