package volumes

import (
	"encoding/json"
	"fmt"
	"git.yelpcorp.com/paasta-tools-go/config"
	"io"
	"io/ioutil"
)

type Volume struct {
	HostPath      string `json:"hostPath"`
	ContainerPath string `json:"containerPath"`
	Mode          string `json:"mode"`
}

func ReadDefaultVolumes(r io.Reader) (volumes []Volume, e error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Errorf("failed reading environment: %g", err)
		return make([]Volume, 0), err
	}
	var data map[string][]Volume
	e = json.Unmarshal(buf, &data)
	if err != nil {
		fmt.Errorf("failed to decode volumes: %g", err)
		return make([]Volume, 0), e
	}
	return data["volumes"], e
}

func DefaultVolumesFromFile() (volumes []Volume, e error) {
	r := config.ReadSystemPaaSTAConfig("/etc/paasta/volumes.json")
	defer r.Close()
	return ReadDefaultVolumes(r)
}
