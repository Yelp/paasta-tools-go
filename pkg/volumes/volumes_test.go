package volumes

import (
	"reflect"
	"testing"
)

type FakeConfigReader struct {
	data VolumeConfig
}

func (fakereader FakeConfigReader) Read(content interface{}) error {
	*content.(*VolumeConfig) = fakereader.data
	return nil
}

func TestDefaultVolumesFromReader(test *testing.T) {
	fakeVolumeConfig := VolumeConfig{Volumes: []Volume{Volume{HostPath: "/foo", ContainerPath: "/bar", Mode: "RO"}}}
	reader := &FakeConfigReader{data: fakeVolumeConfig}
	actual, err := DefaultVolumesFromReader(reader)
	if err != nil {
		test.Errorf("failed to read config")
	}
	if !reflect.DeepEqual(actual, fakeVolumeConfig.Volumes) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", actual, fakeVolumeConfig.Volumes)
	}
}
