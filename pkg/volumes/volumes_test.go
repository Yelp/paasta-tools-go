package volumes

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
)

func TestDefaultVolumesFromReader(test *testing.T) {
	fakeVolumeConfig := &sync.Map{}
	fakeVolumeConfig.Store("volumes", []map[string]string{
		{
			"hostPath":      "/foo",
			"containerPath": "/bar",
			"mode":          "RO",
		},
	},
	)
	reader := &configstore.Store{Data: fakeVolumeConfig}
	actual, err := DefaultVolumesFromReader(reader)
	if err != nil {
		test.Errorf("failed to read config")
	}
	expectedVolume := []Volume{
		{HostPath: "/foo", ContainerPath: "/bar", Mode: "RO"},
	}
	if !reflect.DeepEqual(actual, expectedVolume) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", actual, expectedVolume)
	}
}
func TestDefaultHealthcheckVolumesFromReader(test *testing.T) {
	fakeVolumeConfig := &sync.Map{}
	fakeVolumeConfig.Store("hacheck_sidecar_volumes", []map[string]string{
		{
			"hostPath":      "/foo1",
			"containerPath": "/bar1",
			"mode":          "RW",
		},
	},
	)
	reader := &configstore.Store{Data: fakeVolumeConfig}
	actual, err := DefaultHealthcheckVolumesFromReader(reader)
	if err != nil {
		test.Errorf("failed to read config")
	}
	expectedVolume := []Volume{
		{HostPath: "/foo1", ContainerPath: "/bar1", Mode: "RW"},
	}
	if !reflect.DeepEqual(actual, expectedVolume) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", actual, expectedVolume)
	}
}
