package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NonexistantFile(t *testing.T) {
	r := require.New(t)
	_, err := NewConfig("/tmp/does/not/exist.conf", nil)
	r.Error(err)
}

type testConfig struct {
	Hostname string
	Port     uint16
}

func Test_DecodeFile(t *testing.T) {
	r := require.New(t)
	cfgFile, err := ioutil.TempFile("", "config")
	r.NoError(err)
	defer func() {
		err := cfgFile.Close()
		r.NoError(err)
		os.Remove(cfgFile.Name())
	}()

	expect := `
hostname="localhost"
port=8080
`
	cfgFile.WriteString(expect)

	cfg, err := NewConfig(cfgFile.Name(), nil)
	r.NoError(err)

	target := testConfig{}
	cfg.Decode(&target)
	r.Equal("localhost", target.Hostname)
	r.Equal(uint16(8080), target.Port)
}

func Test_ExampleConfig(t *testing.T) {
	r := require.New(t)
	cfgDefault := testConfig{"localhost", uint16(8080)}

	expect := `hostname="localhost"
port=8080
`
	cfgFile, err := ioutil.TempFile("", "config")
	r.NoError(err)
	defer func() {
		err := cfgFile.Close()
		r.NoError(err)
		os.Remove(cfgFile.Name())
	}()

	cfg, err := NewConfig(cfgFile.Name(), cfgDefault)
	r.NoError(err)
	r.Equal(expect, cfg.EncodeDefault())
}

func Test_ExampleConfigNil(t *testing.T) {
	r := require.New(t)
	cfgFile, err := ioutil.TempFile("", "config")
	r.NoError(err)
	defer func() {
		err := cfgFile.Close()
		r.NoError(err)
		os.Remove(cfgFile.Name())
	}()

	cfg, err := NewConfig(cfgFile.Name(), struct{}{})
	r.NoError(err)
	r.Equal("", cfg.EncodeDefault())
}

// Ensure that megabyte sizes can be parsed.
func TestSize_UnmarshalTOML_MB(t *testing.T) {
	var s Size
	if err := s.UnmarshalTOML([]byte("200m")); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if s != 200*(1<<20) {
		t.Fatalf("unexpected size: %d", s)
	}
}

// Ensure that gigabyte sizes can be parsed.
func TestSize_UnmarshalTOML_GB(t *testing.T) {
	var s Size
	if err := s.UnmarshalTOML([]byte("1g")); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if s != 1073741824 {
		t.Fatalf("unexpected size: %d", s)
	}
}

func TestConfig_Encode(t *testing.T) {
	d := struct {
		WriteTimeout Duration `toml:"write-timeout"`
	}{Duration(time.Minute)}
	buf := new(bytes.Buffer)
	if err := NewEncoder(buf).Encode(d); err != nil {
		t.Fatal("Failed to encode: ", err)
	}
	got, search := buf.String(), `write-timeout="1m0s"`
	if !strings.Contains(got, search) {
		t.Fatalf("Encoding config failed.\nfailed to find %s in:\n%s\n", search, got)
	}
}
