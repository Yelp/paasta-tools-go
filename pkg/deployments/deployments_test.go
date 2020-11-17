package deployments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"github.com/stretchr/testify/assert"
)

const (
	dockerRepo = "docker-paasta.yelpcorp.com:443"
)

func TestDefaultProviderGetDeployment(test *testing.T) {
	paastaConfigData := &sync.Map{}
	paastaConfigData.Store("docker_registry", "fakeregistry.yelp.com")

	serviceConfigData := &sync.Map{}
	serviceConfigData.Store("v2", map[string]interface{}{
		"deployments": map[string]interface{}{
			"dev.every": map[string]interface{}{
				"docker_image": "busybox:latest",
				"git_sha":      "abc123",
			},
			"test.every": map[string]interface{}{
				"docker_image": "ubuntu:latest",
				"git_sha":      "abc123",
			},
		},
	})

	imageProvider := DefaultImageProvider{
		PaastaConfig:  &configstore.Store{Data: paastaConfigData},
		ServiceConfig: &configstore.Store{Data: serviceConfigData},
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

func TestDeploymentsFromConfig(test *testing.T) {
	// create fake deployments file
	data, err := json.MarshalIndent(
		map[string]interface{}{
			"v2": map[string]interface{}{
				"deployments": map[string]interface{}{
					"deploy.group": map[string]interface{}{
						"docker_image": "image-abc123",
						"git_sha":      "abc123",
					},
				},
				"controls": map[string]interface{}{
					"a_service:a_cluster.an_instance": map[string]interface{}{
						"desired_state": "resurrection",
						"force_bounce":  nil,
					},
				},
			},
		},
		"",
		"  ",
	)
	if err != nil {
		test.Error("Failed to serialize test deployment JSON")
	}

	tempDir, err := ioutil.TempDir(os.TempDir(), "paasta-tools-go-test-service-*")
	if err != nil {
		test.Error("Failed to create temp services dir")
	}
	defer os.Remove(tempDir)
	err = ioutil.WriteFile(path.Join(tempDir, "deployments.json"), data, 0644)
	if err != nil {
		test.Error("Failed to write to deployments.json")
	}

	// now actually try to load the deployments
	configStore := configstore.NewStore(tempDir, map[string]string{"v2": "deployments"})

	deployments, err := deploymentsFromConfig(configStore)

	assert.NoError(test, err)
	assert.Equal(test, "abc123", deployments.V2.Deployments["deploy.group"].GitSHA)
	assert.Equal(
		test,
		"resurrection",
		deployments.V2.Controls["a_service:a_cluster.an_instance"].DesiredState,
	)
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
	dir, err := ioutil.TempDir("", "deployments-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = ioutil.WriteFile(
		fmt.Sprintf("%v/docker_registry.json", dir),
		[]byte(`{"docker_registry": "fakeregistry.yelp.com"}`),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	serviceConfigData := &sync.Map{}
	serviceConfigData.Store("v2", map[string]interface{}{
		"deployments": map[string]interface{}{},
	})

	imageProvider := DefaultImageProvider{
		PaastaConfig:  configstore.NewStore(dir, nil),
		ServiceConfig: &configstore.Store{Data: serviceConfigData},
	}
	url, err := imageProvider.getDockerRegistry()
	if err != nil {
		t.Errorf("expected %s actual: error %+v", "fakeregistry.yelp.com", err)
	} else if url != "fakeregistry.yelp.com" {
		t.Errorf("expected %s actual %+v", "fakeregistry.yelp.com", url)
	}
}
