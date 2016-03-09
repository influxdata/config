package config

import (
	"reflect"
	"testing"
	"time"
)

func Test_DurationMarshalTOML(t *testing.T) {
	d := Duration(time.Second)
	data, err := d.MarshalTOML()
	if err != nil {
		t.Fatal(err)
	}
	expData := []byte(`"1s"`)
	if !reflect.DeepEqual(expData, data) {
		t.Errorf("unexpected data: exp %s got %s", string(expData), string(data))
	}
}

func Test_DurationUnmarshalTOML(t *testing.T) {
	data := []byte(`"1s"`)
	d := new(Duration)
	err := d.UnmarshalTOML(data)
	if err != nil {
		t.Fatal(err)
	}
	expDur := Duration(time.Second)

	if expDur != *d {
		t.Errorf("unexpected duration: exp %s got %s", expDur, d)
	}
}
