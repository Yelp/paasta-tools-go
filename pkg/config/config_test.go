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
	reader io.Reader
	data   *FakeConfig
}

func (fakereader fakereader) Read(bytes []byte) (n int, err error) {
	return fakereader.reader.Read(bytes)
}

func fakeDataReader(config *FakeConfig) fakereader {
	content, _ := json.Marshal(*config)
	return fakereader{
		data:   config,
		reader: bytes.NewReader(content),
	}
}

func TestParseContent(test *testing.T) {
	fakeData := &FakeConfig{}
	reader := fakeDataReader(fakeData)
	err := ParseContent(reader, fakeData)
	if err != nil {
		test.Errorf("failed to decode content")
	}
	if !reflect.DeepEqual(reader.data, fakeData) {
		test.Errorf("deserialized content was incorrect, got: %s, want: %s.", reader.data, fakeData)
	}
}

func TestFileNameForConfig(test *testing.T) {
	reader := ConfigFileReader{
		Basedir:  "/etc/paasta",
		Filename: "volumes.json",
	}
	expected := "/etc/paasta/volumes.json"
	actual := reader.FileNameForConfig()
	if actual != expected {
		test.Errorf("filename incorrect incorrect, got: %s, want: %s.", actual, expected)
	}
}
