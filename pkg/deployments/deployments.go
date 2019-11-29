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

	"github.com/Yelp/paasta-tools-go/pkg/config_store"
)

// V2DeploymentGroup ...
type V2DeploymentGroup struct {
	DockerImage string `json:"docker_image" mapstructure:"docker_image"`
	GitSHA      string `json:"git_sha" mapstructure:"git_sha"`
}

// V2ControlGroup ...
type V2ControlGroup struct {
	DesiredState string `json:"desired_state" mapstructure:"desired_state"`
	ForceBounce  string `json:"force_bounce" mapstructure:"force_bounce"`
}

// V2DeploymentsConfig ...
type V2DeploymentsConfig struct {
	Deployments map[string]V2DeploymentGroup `json:"deployments" mapstructure:"deployments"`
	Controls    map[string]V2ControlGroup    `json:"controls" mapstructure:"controls"`
}

// Deployments ...
type Deployments struct {
	V2 V2DeploymentsConfig `json:"v2" mapstructure:"v2"`
}

// DockerRegistry ...
type DockerRegistry struct {
	Registry string `json:"docker_registry" mapstructure:"docker_registry"`
}

type ImageProvider interface {
	DockerImageURLForDeployGroup(deploymentGroup string) (string, error)
}

type DefaultImageProvider struct {
	Service       string
	ServiceConfig *config_store.Store
	PaastaConfig  *config_store.Store
}

// NewDefaultImageProviderForService ...
func NewDefaultImageProviderForService(service string) *DefaultImageProvider {
	serviceConfig := config_store.NewStore(
		path.Join("/nail/etc/services", service),
		map[string]string{"v2": "deployments"},
	)
	paastaConfig := config_store.NewStore(
		"/etc/paasta",
		map[string]string{"registry": "docker_registry"},
	)
	return &DefaultImageProvider{
		Service:       service,
		ServiceConfig: serviceConfig,
		PaastaConfig:  paastaConfig,
	}
}

// DockerImageURLForDeployGroup returns pullable docker image URL
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
	dockerRegistry := &DockerRegistry{Registry: ""}
	err := provider.PaastaConfig.Load("registry", &dockerRegistry.Registry)
	return dockerRegistry.Registry, err
}

func (provider *DefaultImageProvider) getImageForDeployGroup(deploymentGroup string) (string, error) {
	deployments := &Deployments{V2: V2DeploymentsConfig{}}
	err := provider.ServiceConfig.Load("v2", &deployments.V2)
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

// DeploymentAnnotations returns a map of annotations for the relevant service
// deployment group
func DeploymentAnnotations(
	service, cluster, instance, deploymentGroup string,
) (map[string]string, error) {
	configStore := config_store.NewStore(
		fmt.Sprintf("/nail/etc/services/%s", service),
		map[string]string{"v2": "deployments"},
	)
	deployments, err := deploymentsFromConfig(configStore)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading deployments for service %s: %s", service, err,
		)
	}
	controlGroup := makeControlGroup(service, instance, cluster)
	return deploymentAnnotationsForControlGroup(deployments, controlGroup)
}

func deploymentsFromConfig(cr *config_store.Store) (*Deployments, error) {
	deployments := &Deployments{}
	err := cr.Load("v2", deployments)
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
