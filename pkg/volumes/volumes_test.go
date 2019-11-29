package volumes

import (
	"reflect"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/config_store"
)

func TestDefaultVolumesFromReader(test *testing.T) {
	fakeVolumeConfig := map[string]interface{}{
		"volumes": []map[string]interface{}{
			map[string]interface{}{
				"hostPath":      "/foo",
				"containerPath": "/bar",
				"mode":          "RO",
			},
		},
	}
	reader := &config_store.Store{Data: fakeVolumeConfig}
	actual, err := DefaultVolumesFromReader(reader)
	if err != nil {
		test.Errorf("failed to read config")
	}
	expectedVolume := []Volume{
		Volume{HostPath: "/foo", ContainerPath: "/bar", Mode: "RO"},
	}
	if !reflect.DeepEqual(actual, expectedVolume) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", actual, expectedVolume)
	}
}
