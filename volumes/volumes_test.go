package volumes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type fakereader struct {
	r    io.Reader
	data map[string]interface{}
}

func (f fakereader) Read(b []byte) (n int, err error) {
	return f.r.Read(b)
}

func fakeDataReader() fakereader {
	data := map[string]interface{}{"volumes": []Volume{Volume{
		HostPath:      "/nail/etc/mrjob",
		ContainerPath: "/nail/etc/mrjob",
		Mode:          "RO",
	}}}
	content, _ := json.Marshal(data)
	return fakereader{
		data: data,
		r:    bytes.NewReader(content),
	}
}

func TestReadDefaultVolumes(t *testing.T) {
	fakeData := fakeDataReader()
	result, err := ReadDefaultVolumes(fakeData)
	if err != nil {
		t.Errorf("failed to decode environment")
	}
	if !reflect.DeepEqual(result, fakeData.data["volumes"]) {
		t.Errorf("environment was incorrect, got: %s, want: %s.", fakeData.data["volumes"], result)
	}
}
