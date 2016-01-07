package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

type Duration time.Duration

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

func (d Duration) MarshalTOML() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

const maxInt = int64(^uint(0) >> 1)

type Size int

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

func (d Duration) String() string {
	return time.Duration(d).String()
}

func Decode(tomlBlob string, target interface{}) error {
	sreader := strings.NewReader(tomlBlob)
	decoder := toml.NewDecoder(sreader)
	err := decoder.Decode(target)

	if err != nil {
		return err
	}

	return nil
}

func DecodeFile(fpath string, v interface{}) error {
	file, err := os.Open(fpath)

	if err != nil {
		return err
	}
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func ParseFile(fpath string) (*ast.Table, error) {
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return toml.Parse(contents)
}

// Kapacitor
func NewEncoder(w io.Writer) *toml.Encoder {
	return toml.NewEncoder(w)
}

func UnmarshalTable(t *ast.Table, v interface{}) error {
	return toml.UnmarshalTable(t, v)
}
