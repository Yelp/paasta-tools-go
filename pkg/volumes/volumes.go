package volumes

import (
	"fmt"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
)

type VolumeConfig struct {
	Volumes []Volume `json:"volumes" mapstructure:"volumes"`
}

type Volume struct {
	HostPath      string `json:"hostPath" mapstructure:"hostPath"`
	ContainerPath string `json:"containerPath" mapstructure:"containerPath"`
	Mode          string `json:"mode" mapstructure:"mode"`
}

func DefaultVolumesFromReader(configStore *configstore.Store) ([]Volume, error) {
	volumeConfig := &VolumeConfig{Volumes: []Volume{}}
	ok, err := configStore.Load("volumes", &volumeConfig.Volumes)
	if !ok {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig.Volumes, err
}
