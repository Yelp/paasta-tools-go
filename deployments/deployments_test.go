package deployments

import (
	"fmt"
	"testing"

	paastaconfig "github.com/Yelp/paasta-tools-go/config"
)

const (
	dockerRepo = "docker-paasta.yelpcorp.com:443"
)

type FakeConfigReader struct {
	data Deployments
}

func (fakereader FakeConfigReader) Read(content interface{}) error {
	*content.(*Deployments) = fakereader.data
	return nil
}

func TestGetImageURL(test *testing.T) {
	fakeDeployments := Deployments{
		V2: V2DeploymentsConfig{
			Deployments: map[string]V2DeploymentGroup{
				"dev.every": V2DeploymentGroup{
					DockerImage: "busybox:latest",
					GitSHA:      "03d6f783c99695af0e716588abb9ba83ac957be2",
				},
				"test.every": V2DeploymentGroup{
					DockerImage: "ubuntu:latest",
					GitSHA:      "f3d6f783c99695af0e716588abb9ba83ac957be3",
				},
			},
		},
	}
	reader := &FakeConfigReader{data: fakeDeployments}
	testcases := map[string]string{
		"dev.every":  "busybox:latest",
		"test.every": "ubuntu:latest",
		"absent":     "",
	}
	var expected string
	for dment, imageurl := range testcases {
		actual, _ := getImageURL(reader, dment, dockerRepo)
		if imageurl != "" {
			expected = fmt.Sprintf("%s/%s", dockerRepo, imageurl)
		} else {
			expected = ""
		}
		if actual != expected {
			test.Errorf("Failed for %s %s, expected %s", dment, actual, expected)
		}
	}
}

func TestGetImageURLEmptyReader(test *testing.T) {
	actual, err := getImageURL(paastaconfig.SystemPaaSTAConfigFileReader{}, "", dockerRepo)
	if err == nil || actual != "" {
		test.Errorf("Expected to fail for nil interface")
	}
}
