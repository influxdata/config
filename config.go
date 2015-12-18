package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	ntoml "github.com/naoina/toml"
	nast "github.com/naoina/toml/ast"
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

func Decode(tomlBlob string, target interface{}) (toml.MetaData, error) {
	sreader := strings.NewReader(tomlBlob)
	decoder := ntoml.NewDecoder(sreader)
	err := decoder.Decode(target)
	meta := toml.MetaData{}

	if err != nil {
		return meta, err
	}

	return meta, nil
}

func DecodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	file, err := os.Open(fpath)
	meta := toml.MetaData{}

	if err != nil {
		return meta, err
	}
	decoder := ntoml.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func ParseFile(fpath string) (*nast.Table, error) {
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return ntoml.Parse(contents)
}

// Kapacitor
func NewEncoder(w io.Writer) *ntoml.Encoder {
	return ntoml.NewEncoder(w)
}

func UnmarshalTable(t *nast.Table, v interface{}) error {
	return ntoml.UnmarshalTable(t, v)
}
