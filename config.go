package config

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

// Config is a configuration management object, providing services for
// unmarshalling, remotely reading, and updating an underlying configuration
// object
type Config struct {
	path           string      // the path of the underlying configuration file
	exampleDefault interface{} // an instance of the target config object which is used to generate example configs
}

// NewConfig initializes a configuration management object for the config file
// located at the provided path.
func NewConfig(path string, exampleDefault interface{}) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &Config{path, exampleDefault}, nil
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

// Encode marshals the provided struct and returns the resultant TOML
func (c *Config) Encode(target interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	enc := NewEncoder(buf)
	err := enc.Encode(target)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
