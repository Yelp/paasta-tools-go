package volumes

import (
	"fmt"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
)

type VolumeConfig struct {
	Volumes               []Volume `json:"volumes" mapstructure:"volumes"`
	HacheckSidecarVolumes []Volume `json:"hacheck_sidecar_volumes" mapstructure:"hacheck_sidecar_volumes"`
}

type Volume struct {
	HostPath      string `json:"hostPath" mapstructure:"hostPath"`
	ContainerPath string `json:"containerPath" mapstructure:"containerPath"`
	Mode          string `json:"mode" mapstructure:"mode"`
}

func DefaultVolumesFromReader(configStore *configstore.Store) ([]Volume, error) {
	volumeConfig := &VolumeConfig{}
	ok, err := configStore.Load("volumes", &volumeConfig.Volumes)
	if !ok {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig.Volumes, err
}

func DefaultHealthcheckVolumesFromReader(configStore *configstore.Store) ([]Volume, error) {
	volumeConfig := &VolumeConfig{}
	ok, err := configStore.Load("hacheck_sidecar_volumes", &volumeConfig.HacheckSidecarVolumes)
	if !ok {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig.HacheckSidecarVolumes, err
}
