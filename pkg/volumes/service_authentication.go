package volumes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

const (
	authenticatingServicesConfigPath    = "/nail/etc/services/authenticating.yaml"
	jwtServiceAuthConfigName            = "jwt_service_auth"
	authenticatingServicesCacheDuration = 10 * time.Minute
	defaultTokenExpirationSeconds       = 3600
)

type authenticatingServicesConfig struct {
	Services []string `yaml:"services"`
}

type jwtServiceAuthTokenSettings struct {
	Audience          string `json:"audience"`
	ContainerPath     string `json:"container_path"`
	ExpirationSeconds int64  `json:"expiration_seconds"`
}

type jwtServiceAuthConfig struct {
	TokenSettings jwtServiceAuthTokenSettings `json:"service_auth_token_settings"`
}

var authenticatingServicesCache map[string]bool
var authenticatingServicesLastLoaded time.Time

func getAuthenticatingServices(configPath string) (map[string]bool, error) {
	now := time.Now()
	if !now.After(authenticatingServicesLastLoaded.Add(authenticatingServicesCacheDuration)) {
		return authenticatingServicesCache, nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var authenticatingServices authenticatingServicesConfig
	if err = yaml.Unmarshal(data, &authenticatingServices); err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, service := range authenticatingServices.Services {
		result[service] = true
	}
	authenticatingServicesCache = result
	authenticatingServicesLastLoaded = now
	return result, nil
}

func ServiceRequiresAuthenticationToken(serviceName string) (bool, error) {
	authenticatingServices, err := getAuthenticatingServices(authenticatingServicesConfigPath)
	if err != nil {
		return false, err
	}
	_, found := authenticatingServices[serviceName]
	return found, nil
}

func formatServiceAccountVolumeName(audience string) string {
	formatted := strings.Replace(audience, ".", "dot-", -1)
	formatted = formatMountName(formatted)
	formatted = fmt.Sprintf("projected-sa--%s", formatted)
	if len(formatted) > 63 {
		return formatted[0:63]
	}
	return formatted
}

func GetProjectedServiceAccountVolume(audience string, path string, expirationSeconds int64) (corev1.VolumeMount, corev1.Volume) {
	volumeName := formatServiceAccountVolumeName(audience)
	if expirationSeconds == 0 {
		expirationSeconds = defaultTokenExpirationSeconds
	}
	volume := corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				Sources: []corev1.VolumeProjection{
					{
						ServiceAccountToken: &corev1.ServiceAccountTokenProjection{
							Audience:          audience,
							ExpirationSeconds: &expirationSeconds,
							Path:              "token",
						},
					},
				},
			},
		},
	}
	volumeMount := corev1.VolumeMount{
		Name:      volumeName,
		ReadOnly:  true,
		MountPath: path,
	}
	return volumeMount, volume
}

func GetServiceAuthenticationTokenVolume(configStore *configstore.Store) (corev1.VolumeMount, corev1.Volume, error) {
	var jwtServiceAuth jwtServiceAuthConfig
	if ok, err := configStore.Load(jwtServiceAuthConfigName, &jwtServiceAuth); !ok || err != nil {
		if err == nil {
			err = fmt.Errorf("%s configuration not found", jwtServiceAuthConfigName)
		}
		return corev1.VolumeMount{}, corev1.Volume{}, err
	}
	tokenSettings := jwtServiceAuth.TokenSettings
	if tokenSettings.Audience == "" || tokenSettings.ContainerPath == "" {
		return corev1.VolumeMount{}, corev1.Volume{}, fmt.Errorf("Missing token settings in %s configuration", jwtServiceAuthConfigName)
	}
	volume, volumeMount := GetProjectedServiceAccountVolume(tokenSettings.Audience, tokenSettings.ContainerPath, tokenSettings.ExpirationSeconds)
	return volume, volumeMount, nil
}
