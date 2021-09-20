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

func DefaultVolumeConfigFromReader(configStore *configstore.Store) (*VolumeConfig, error) {
	volumeConfig := &VolumeConfig{}
	ok, err := configStore.Load("volumes", &volumeConfig)
	if !ok {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig, err
}

func DefaultVolumesFromReader(configStore *configstore.Store) ([]Volume, error) {
	volumeConfig, err := DefaultVolumeConfigFromReader(configStore)
	if err != nil {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig.Volumes, err
}

func DefaultHealthcheckVolumesFromReader(configStore *configstore.Store) ([]Volume, error) {
	volumeConfig, err := DefaultVolumeConfigFromReader(configStore)
	if err != nil {
		return nil, fmt.Errorf("volumes not found")
	}
	return volumeConfig.HacheckSidecarVolumes, err
}
