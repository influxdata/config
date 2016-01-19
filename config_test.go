package config

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RemoteRead(t *testing.T) {
	cfgFile, err := ioutil.TempFile("", "config")
	r := require.New(t)
	defer func() {
		err := cfgFile.Close()
		r.NoError(err)
		os.Remove(cfgFile.Name())
	}()

	expect := `
[test]
hostname="localhost"
port=8080
`
	cfgFile.WriteString(expect)

	cfg, err := NewConfig(cfgFile.Name(), nil)
	r.NoError(err)
	defer os.Remove(cfgFile.Name())

	ts := httptest.NewServer(cfg.HTTP())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	r.NoError(err)

	actual, err := ioutil.ReadAll(res.Body)
	r.NoError(err)

	r.Equal(expect, string(actual))
}

func Test_NonexistantFile(t *testing.T) {
	r := require.New(t)
	_, err := NewConfig("/tmp/does/not/exist.conf", nil)
	r.Error(err)
}

func Test_RemoteUpdate(t *testing.T) {
	r := require.New(t)
	cfgFile, err := ioutil.TempFile("", "config")
	r.NoError(err)
	defer os.Remove(cfgFile.Name())

	cfg, err := NewConfig(cfgFile.Name(), nil)
	r.NoError(err)

	expect := `
[test]
hostname="localhost"
port=8080
`

	ts := httptest.NewServer(cfg.HTTP())
	defer ts.Close()

	resp, err := http.Post(ts.URL, "text/plain", strings.NewReader(expect))
	r.NoError(err)
	defer resp.Body.Close()

	r.Equal(200, resp.StatusCode)

	actual, err := ioutil.ReadAll(cfgFile)
	r.NoError(err)

	r.Equal(expect, string(actual))
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
