// Package deployments provides functions for decoding V2 paasta deployments
// of the form:
//   "v2": {
//    "deployments": {
//		"everything": {
//		  "docker_image": "services-fluffy:paasta-abcdfff",
//		  "git_sha": "abcdfff"
//		},
//		"stagef": {
//		  "docker_image": "services-fluffy:paasta-abcdfff",
//		  "git_sha": "abcdfff"
//		}
//	  },
//  }
package deployments

import (
	"fmt"
	"path"
	"strings"

	"github.com/Yelp/paasta-tools-go/pkg/config"
)

// V2DeploymentGroup ...
type V2DeploymentGroup struct {
	DockerImage string `json:"docker_image"`
	GitSHA      string `json:"git_sha"`
}

// V2ControlGroup ...
type V2ControlGroup struct {
	DesiredState string `json:"desired_state"`
	ForceBounce  string `json:"force_bounce"`
}

// V2DeploymentsConfig ...
type V2DeploymentsConfig struct {
	Deployments map[string]V2DeploymentGroup `json:"deployments"`
	Controls    map[string]V2ControlGroup    `json:"controls"`
}

// Deployments ...
type Deployments struct {
	V2 V2DeploymentsConfig `json:"v2"`
}

// DockerRegistry ...
type DockerRegistry struct {
	Registry string `json:"docker_registry"`
}

type ImageProvider interface {
	DockerImageURLForDeployGroup(deploymentGroup string) (string, error)
}

type DefaultImageProvider struct {
	Service           string
	RegistryURLReader config.ConfigReader
	ImageReader       config.ConfigReader
}

func NewDefaultImageProviderForService(service string) *DefaultImageProvider {
	imageReader := config.ConfigFileReader{
		Basedir:  path.Join("/nail/etc/services", service),
		Filename: "deployments.json",
	}
	registryURLReader := config.ConfigFileReader{
		Basedir:  "/etc/paasta",
		Filename: "docker_registry.json",
	}
	return &DefaultImageProvider{
		Service:           service,
		RegistryURLReader: registryURLReader,
		ImageReader:       imageReader,
	}
}

// DockerImageURLForService returns pullable docker image URL
func (provider *DefaultImageProvider) DockerImageURLForDeployGroup(deploymentGroup string) (string, error) {
	var image string
	registry, err := provider.getDockerRegistry()
	if err != nil {
		return "", fmt.Errorf("Failed to get docker registry: %s", err)
	}
	image, err = provider.getImageForDeployGroup(deploymentGroup)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to read from deployments.json for %s: %s",
			provider.Service, err,
		)
	}

	return fmt.Sprintf("%s/%s", registry, image), nil
}

func (provider *DefaultImageProvider) getDockerRegistry() (string, error) {
	dockerRegistry := &DockerRegistry{}
	err := provider.RegistryURLReader.Read(dockerRegistry)
	return dockerRegistry.Registry, err

}

func (provider *DefaultImageProvider) getImageForDeployGroup(deploymentGroup string) (string, error) {
	deployments := &Deployments{}
	err := provider.ImageReader.Read(deployments)
	if err != nil {
		return "", err
	}
	deployment, ok := deployments.V2.Deployments[deploymentGroup]

	if !ok {
		return "", fmt.Errorf(
			"Deployment group %s not found in v2 deployments of %+v",
			deploymentGroup, deployments,
		)
	}
	return deployment.DockerImage, nil
}

func GetPaastaGitShaFromDockerURL(dockerUrl string) (string, error) {
	image := strings.Split(dockerUrl, "/")
	if len(image) != 2 {
		return "", fmt.Errorf(
			"Failed to extract paasta git sha from url: %s",
			dockerUrl,
		)
	}
	gitShaStrings := strings.Split(image[1], "-")
	if len(gitShaStrings) != 2 {
		return "", fmt.Errorf(
			"Failed to extract paasta git sha from url: %s",
			dockerUrl,
		)
	}
	if len(gitShaStrings[1]) < 8 {
		return "", fmt.Errorf(
			"%s doesn't look like a git sha, not long enough",
			gitShaStrings[1],
		)
	}
	return "git" + gitShaStrings[1][:8], nil
}

// DeploymentAnnotations returns a map of annotations for the relevant service
// deployment group
func DeploymentAnnotations(
	service, cluster, instance, deploymentGroup string,
) (map[string]string, error) {
	configReader := config.ConfigFileReader{
		Basedir:  fmt.Sprintf("/nail/etc/services/%s", service),
		Filename: "deployments.json",
	}
	deployments, err := deploymentsFromConfig(configReader)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading deployments for service %s: %s", service, err,
		)
	}
	controlGroup := makeControlGroup(service, instance, cluster)
	return deploymentAnnotationsForControlGroup(deployments, controlGroup)
}

func deploymentsFromConfig(cr config.ConfigFileReader) (*Deployments, error) {
	deployments := &Deployments{}
	err := cr.Read(deployments)
	return deployments, err
}

func makeControlGroup(service, instance, cluster string) string {
	return fmt.Sprintf("%s:%s.%s", service, cluster, instance)
}

func deploymentAnnotationsForControlGroup(ds *Deployments, cg string) (map[string]string, error) {
	annotations := map[string]string{}
	control, ok := ds.V2.Controls[cg]
	if !ok {
		return nil, fmt.Errorf("Control group %s does not exist", cg)
	}
	annotations["paasta.yelp.com/desired_state"] = control.DesiredState
	annotations["paasta.yelp.com/force_bounce"] = control.ForceBounce
	return annotations, nil
}
