package deployments

import (
	"fmt"
	"testing"

	paastaconfig "github.com/Yelp/paasta-tools-go/pkg/config"
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

func TestMakeControlGroup(test *testing.T) {
	expected := "service:cluster.instance"
	actual := makeControlGroup("service", "instance", "cluster")
	if expected != actual {
		test.Errorf("Expected '%+v', got '%+v'", expected, actual)
	}
}

func TestDeploymentAnnotationsForControlGroup(test *testing.T) {
	fakeDeployments := &Deployments{
		V2: V2DeploymentsConfig{
			Controls: map[string]V2ControlGroup{
				"test-cg": V2ControlGroup{
					ForceBounce:  "test-bounce",
					DesiredState: "test-state",
				},
			},
		},
	}
	expectedAnns := map[string]string{
		"paasta.yelp.com/desired_state": "test-state",
		"paasta.yelp.com/force_bounce":  "test-bounce",
	}
	anns, err := deploymentAnnotationsForControlGroup(fakeDeployments, "test-cg")
	if err != nil {
		test.Errorf("Expected to not fail: %s", err)
		return
	}
	if len(anns) != len(expectedAnns) {
		test.Errorf("Expected '%+v', got '%+v'", expectedAnns, anns)
		return
	}
	for k, v := range anns {
		ev, _ := expectedAnns[k]
		if v != ev {
			test.Errorf("Expected %s to be '%+v', got '%+v'", k, ev, v)
			return
		}
	}
}
