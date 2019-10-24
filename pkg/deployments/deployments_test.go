package deployments

import (
	"fmt"
	"testing"
)

const (
	dockerRepo = "docker-paasta.yelpcorp.com:443"
)

type FakeDeploymentsReader struct {
	data Deployments
}

func (fakereader FakeDeploymentsReader) Read(content interface{}) error {
	*content.(*Deployments) = fakereader.data
	return nil
}

type FakeRegistryReader struct {
	registry DockerRegistry
}

func (fakereader FakeRegistryReader) Read(content interface{}) error {
	*content.(*DockerRegistry) = fakereader.registry
	return nil
}

type StaticImageProvider struct {
	DockerRegistry string
	Image          string
}

func NewStaticImageProvider(dockerRegistry, image string) *StaticImageProvider {
	return &StaticImageProvider{
		DockerRegistry: dockerRegistry,
		Image:          image,
	}
}

func (provider StaticImageProvider) DockerImageURLForService(serviceName, deploymentGroup string) (string, error) {
	return fmt.Sprintf("%s/%s", provider.DockerRegistry, provider.Image), nil
}

func TestDefaultProviderGetDeployment(test *testing.T) {
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
	registry := DockerRegistry{
		Registry: "fakeregistry.yelp.com",
	}
	imageReader := &FakeDeploymentsReader{data: fakeDeployments}
	registryReader := &FakeRegistryReader{registry: registry}
	imageProvider := DefaultImageProvider{
		RegistryURLReader: registryReader,
		ImageReader:       imageReader,
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
	fakeDeployments := Deployments{
		V2: V2DeploymentsConfig{
			Deployments: map[string]V2DeploymentGroup{},
		},
	}
	registry := DockerRegistry{
		Registry: "fakeregistry.yelp.com",
	}
	imageReader := &FakeDeploymentsReader{data: fakeDeployments}
	registryReader := &FakeRegistryReader{registry: registry}
	imageProvider := DefaultImageProvider{
		RegistryURLReader: registryReader,
		ImageReader:       imageReader,
	}
	url, _ := imageProvider.getDockerRegistry()
	if url != registry.Registry {
		t.Errorf("expected correct docker registry url")
	}
}
