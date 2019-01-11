package volumes

import (
	"fmt"
	"git.yelpcorp.com/paasta-tools-go/config"
)

type VolumeConfig struct {
	Volumes []Volume `json:"volumes"`
}

type Volume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	Mode          string `json:"mode"`
}

func DefaultVolumesFromReader(c config.ConfigReader) (volumes []Volume, err error) {
	volumeConfig := VolumeConfig{}
	err := c.Read(volumeConfig)
	return volumeConfig.Volumes, err
}
