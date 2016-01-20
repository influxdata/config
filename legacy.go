package config

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

// Decode unmarshals a string of TOML into a target configuration struct.
func Decode(tomlBlob string, target interface{}) error {
	sreader := strings.NewReader(tomlBlob)
	decoder := toml.NewDecoder(sreader)
	err := decoder.Decode(target)

	if err != nil {
		return err
	}

	return nil
}

// DecodeFile loads a TOML configuration from a provided path and unmarshals it
// into a target configuration struct.
func DecodeFile(path string, target interface{}) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}
	return nil
}

// ParseFile loads a TOML configuration from a provided path and returns the
// AST produced from the TOML parser. This function was originally provided by naoina/toml
func ParseFile(fpath string) (*ast.Table, error) {
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	return toml.Parse(contents)
}

// NewEncoder produces a struct capable of writing TOML-encoded data to the
// provided io.Writer. This function was originally provided by BurntSushi/toml
func NewEncoder(w io.Writer) *toml.Encoder {
	return toml.NewEncoder(w)
}

// UnmarshalTable provides a mechanism to incrementally unmarshal an AST
// produced by ParseFile. This method was originally provided by naoina/toml
func UnmarshalTable(t *ast.Table, v interface{}) error {
	return toml.UnmarshalTable(t, v)
}
