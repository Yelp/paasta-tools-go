package volumes

import (
	"fmt"
	"log"
	"strings"

	"github.com/Yelp/paasta-tools-go/pkg/config"
	corev1 "k8s.io/api/core/v1"
)

func paastaVolumesToKubernetesVolumes(
	vols []Volume,
) ([]corev1.VolumeMount, []corev1.Volume) {
	volumeMounts := make([]corev1.VolumeMount, len(vols))
	volumes := make([]corev1.Volume, len(vols))
	for i, v := range vols {
		readOnly := true
		if v.Mode == "RW" {
			readOnly = false
		}
		name := formatMountName(v.HostPath)

		volumeMounts[i] = corev1.VolumeMount{
			Name:      name,
			ReadOnly:  readOnly,
			MountPath: v.ContainerPath,
		}

		volumes[i] = corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: v.HostPath,
				},
			},
		}
	}

	return volumeMounts, volumes
}

func formatMountName(hostPath string) string {
	var formatted string
	formatted = strings.TrimRight(hostPath, "/")
	formatted = strings.TrimLeft(formatted, "/")
	formatted = strings.Replace(formatted, "/", "-", -1)
	formatted = strings.Replace(formatted, "_", "--", -1)
	return formatted
}

// GetDefaultPaastaKubernetesVolumes ...
func GetDefaultPaastaKubernetesVolumes(configStore *config.Store) ([]corev1.VolumeMount, []corev1.Volume, error) {
	pvolumes, err := DefaultVolumesFromReader(configStore)
	if err != nil {
		err = fmt.Errorf("Error finding default volumes: %s", err)
		log.Print(err)
		return nil, nil, err
	}
	volumeMounts, volumes := paastaVolumesToKubernetesVolumes(pvolumes)
	return volumeMounts, volumes, err
}
