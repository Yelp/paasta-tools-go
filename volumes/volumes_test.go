package volumes

import (
	"reflect"
	"testing"
)

type FakeConfigReader struct {
	data VolumeConfig
}

func (f FakeConfigReader) Read(t interface{}) error {
	t = f.data
	return nil
}

func TestDefaultVolumesFromReader(t *testing.T) {
	fakeVolumeConfig := VolumeConfig{Volumes: []Volume{Volume{HostPath: "/foo", ContainerPath: "/bar", Mode: "RO"}}}
	reader := &FakeConfigReader{data: fakeVolumeConfig}
	actual, err := DefaultVolumesFromReader(reader)
	if err != nil {
		t.Errorf("failed to read config")
	}
	if reflect.DeepEqual(actual, fakeVolumeConfig.Volumes) {
		t.Errorf("volumes incorrect, got: %s, want: %s.", actual, fakeVolumeConfig.Volumes)
	}
}
