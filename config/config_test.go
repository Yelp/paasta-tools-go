package config

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
)

type FakeConfig struct {
	Foo string
}

type fakereader struct {
	r    io.Reader
	data *FakeConfig
}

func (f fakereader) Read(b []byte) (n int, err error) {
	return f.r.Read(b)
}

func fakeDataReader(c *FakeConfig) fakereader {
	content, _ := json.Marshal(*c)
	return fakereader{
		data: c,
		r:    bytes.NewReader(content),
	}
}

func TestParseContent(t *testing.T) {
	fakeData := &FakeConfig{}
	reader := fakeDataReader(fakeData)
	err := ParseContent(reader, fakeData)
	if err != nil {
		t.Errorf("failed to decode content")
	}
	if !reflect.DeepEqual(reader.data, fakeData) {
		t.Errorf("deserialized content was incorrect, got: %s, want: %s.", reader.data, fakeData)
	}
}

func TestFileNameForConfig(t *testing.T) {
	reader := SystemPaaSTAConfigFileReader{
		Basedir:  "/etc/paasta",
		Filename: "volumes.json",
	}
	expected := "/etc/paasta/volumes.json"
	actual := reader.FileNameForConfig()
	if actual != expected {
		t.Errorf("filename incorrect incorrect, got: %s, want: %s.", actual, expected)
	}
}
