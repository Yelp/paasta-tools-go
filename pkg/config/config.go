package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type ConfigReader interface {
	Read(interface{}) error
}

type ConfigFileReader struct {
	Basedir  string
	Filename string
}

func ParseContent(reader io.Reader, content interface{}) error {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, content)
	return err
}

func (configReader ConfigFileReader) FileNameForConfig() string {
	return path.Join(configReader.Basedir, configReader.Filename)
}

func (configReader ConfigFileReader) Read(content interface{}) error {
	reader, err := os.Open(configReader.FileNameForConfig())
	defer reader.Close()
	if err != nil {
		return err
	}
	return ParseContent(reader, content)
}
