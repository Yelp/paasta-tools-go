package deployments

import (
	"fmt"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/config_store"
)

const (
	dockerRepo = "docker-paasta.yelpcorp.com:443"
)

func TestDefaultProviderGetDeployment(test *testing.T) {
	imageProvider := DefaultImageProvider{
		PaastaConfig: &config_store.Store{
			Data: map[string]interface{}{
				"docker_registry": map[string]interface{}{
					"registry": "fakeregistry.yelp.com",
				},
			},
		},
		ServiceConfig: &config_store.Store{
			Data: map[string]interface{}{
				"v2": map[string]interface{}{
					"deployments": map[string]interface{}{
						"dev.every": map[string]interface{}{
							"docker_image": "busybox:latest",
							"git_sha":      "03d6f783c99695af0e716588abb9ba83ac957be2",
						},
						"test.every": map[string]interface{}{
							"docker_image": "ubuntu:latest",
							"git_sha":      "f3d6f783c99695af0e716588abb9ba83ac957be3",
						},
					},
				},
			},
		},
	}
	testcases := map[string]string{
		"dev.every":  "busybox:latest",
		"test.every": "ubuntu:latest",
		"absent":     "",
	}
	var expected string
	for dment, image := range testcases {
		actual, _ := imageProvider.getImageForDeployGroup(dment)
		if image != "" {
			expected = fmt.Sprintf("%s", image)
		} else {
			expected = ""
		}
		if actual != expected {
			test.Errorf("Failed for %s %s, expected %s", dment, actual, expected)
		}
	}
}

func TestGetImageForDeployGroupEmptyReader(test *testing.T) {
	imageProvider := NewDefaultImageProviderForService("myservice")
	actual, err := imageProvider.getImageForDeployGroup("")
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

func TestDefaultGetRegistry(t *testing.T) {
	imageProvider := DefaultImageProvider{
		PaastaConfig: &config_store.Store{
			Data: map[string]interface{}{
				"registry": "fakeregistry.yelp.com",
			},
		},
		ServiceConfig: &config_store.Store{
			Data: map[string]interface{}{
				"v2": map[string]interface{}{
					"deployments": map[string]interface{}{},
				},
			},
		},
	}
	url, _ := imageProvider.getDockerRegistry()
	if url != "fakeregistry.yelp.com" {
		t.Errorf("expected %s actual %+v", "fakeregistry.yelp.com", url)
	}
}
