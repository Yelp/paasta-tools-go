package volumes

import (
	"github.com/Yelp/paasta-tools-go/pkg/config"
)

type VolumeConfig struct {
	Volumes []Volume `json:"volumes" mapstructure:"volumes"`
}

type Volume struct {
	HostPath      string `json:"hostPath" mapstructure:"hostPath"`
	ContainerPath string `json:"containerPath" mapstructure:"containerPath"`
	Mode          string `json:"mode" mapstructure:"mode"`
}

func DefaultVolumesFromReader(configStore *config.Store) (volumes []Volume, err error) {
	volumeConfig := &VolumeConfig{Volumes: []Volume{}}
	err = configStore.Load("volumes", &volumeConfig.Volumes)
	return volumeConfig.Volumes, err
}
