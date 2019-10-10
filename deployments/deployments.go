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

	paastaconfig "github.com/Yelp/paasta-tools-go/config"
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

func getDockerRegistry() (string, error) {
	dockerRegistry := &DockerRegistry{}
	configreader := paastaconfig.SystemPaaSTAConfigFileReader{
		Basedir:  fmt.Sprintf("/etc/paasta"),
		Filename: "docker_registry.json",
	}
	err := configreader.Read(dockerRegistry)
	return dockerRegistry.Registry, err

}

// DockerImageURLForService returns pullable docker image URL
// for service given a deployment.
func DockerImageURLForService(serviceName, deploymentGroup string) (string, error) {
	var image string
	configReader := paastaconfig.SystemPaaSTAConfigFileReader{
		Basedir:  fmt.Sprintf("/nail/etc/services/%s", serviceName),
		Filename: "deployments.json",
	}
	registry, err := getDockerRegistry()
	if err != nil {
		return "", fmt.Errorf("Failed to get docker registry")
	}
	image, err = getImageURL(configReader, deploymentGroup, registry)
	if err != nil {
		return "", fmt.Errorf(
			"Unable to read from deployments.json for %s",
			serviceName,
		)
	}
	return image, nil
}

func getImageURL(cReader paastaconfig.ConfigReader, deploymentGroup, dockerRepo string) (string, error) {
	deployments := &Deployments{}
	err := cReader.Read(deployments)
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

	imageurl := fmt.Sprintf("%s/%s", dockerRepo, deployment.DockerImage)
	return imageurl, nil
}

// DeploymentAnnotations returns a map of annotations for the relevant service
// deployment group
func DeploymentAnnotations(
	service, cluster, instance, deploymentGroup string,
) (map[string]string, error) {
	configReader := paastaconfig.SystemPaaSTAConfigFileReader{
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

func deploymentsFromConfig(cr paastaconfig.SystemPaaSTAConfigFileReader) (*Deployments, error) {
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
