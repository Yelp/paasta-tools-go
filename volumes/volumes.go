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

func DefaultVolumesFromReader(c config.ConfigReader) (volumes []Volume, e error) {
	volumeConfig := VolumeConfig{}
	err := c.Read(volumeConfig)
	if err != nil {
		fmt.Println("couldn't load default volumes")
	}
	return volumeConfig.Volumes, err
}
