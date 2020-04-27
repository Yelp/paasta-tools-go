package volumes

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
)

func TestDefaultVolumesFromReader(test *testing.T) {
	fakeVolumeConfig := &sync.Map{}
	fakeVolumeConfig.Store("volumes", []map[string]interface{}{
		map[string]interface{}{
			"hostPath":      "/foo",
			"containerPath": "/bar",
			"mode":          "RO",
		},
	})
	reader := &configstore.Store{Data: fakeVolumeConfig}
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
