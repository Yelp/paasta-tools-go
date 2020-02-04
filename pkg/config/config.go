package config

import (
	"io"
	"os"
	"path"

	yaml "k8s.io/apimachinery/pkg/util/yaml"
)

type ConfigReader interface {
	Read(interface{}) error
}

type ConfigFileReader struct {
	Basedir  string
	Filename string
}

func ParseContent(reader io.Reader, content interface{}) error {
	decoder := yaml.NewYAMLToJSONDecoder(reader)
	return decoder.Decode(content)
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
