package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/toml"
	"github.com/naoina/toml/ast"
)

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

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

// Config is a configuration management object, providing services for
// unmarshalling, remotely reading, and updating an underlying configuration
// object
type Config struct {
	path           string      // the path of the underlying configuration file
	exampleDefault interface{} // an instance of the target config object which is used to generate example configs
}

// HTTP returns handlers necessary to facilitate remotely reading and updating
// the underlying configuration object
func (c *Config) HTTP() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			cfgFile, err := os.Open(c.path)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer cfgFile.Close()
			io.Copy(rw, cfgFile)
		case "POST":
			cfgFile, err := os.OpenFile(c.path, os.O_WRONLY, 0666)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer cfgFile.Close()
			io.Copy(cfgFile, r.Body)
		}
	})
}

// Decode unmarshalls the underlying configuration file into the target object.
func (c *Config) Decode(target interface{}) error {
	return DecodeFile(c.path, target)
}

// EncodeDefault produces the default configuration from the example object.
func (c *Config) EncodeDefault() string {
	out, err := c.Encode(c.exampleDefault)
	if err != nil {
		panic(err)
	}
	return out
}

func (c *Config) Encode(target interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	enc := NewEncoder(buf)
	err := enc.Encode(target)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// NewConfig initializes a configuration management object for the config file
// located at the provided path
func NewConfig(path string, exampleDefault interface{}) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &Config{path, exampleDefault}, nil
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
