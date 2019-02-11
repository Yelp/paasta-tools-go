package volumes

import (
	"github.com/Yelp/paasta-tools-go/config"
)

type VolumeConfig struct {
	Volumes []Volume `json:"volumes"`
}

type Volume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	Mode          string `json:"mode"`
}

func DefaultVolumesFromReader(configReader config.ConfigReader) (volumes []Volume, err error) {
	volumeConfig := &VolumeConfig{}
	err = configReader.Read(volumeConfig)
	return volumeConfig.Volumes, err
}
