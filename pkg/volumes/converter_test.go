package volumes

import (
	"reflect"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/config"
	corev1 "k8s.io/api/core/v1"
)

func TestPaastaVolumesToKubernetesVolumes(t *testing.T) {
	fakeVols := []Volume{
		Volume{
			HostPath:      "/tmp/rw",
			ContainerPath: "/tmp/bar",
			Mode:          "RW",
		},
		Volume{
			HostPath:      "/tmp/ro",
			ContainerPath: "/tmp/bar",
			Mode:          "RO",
		},
	}
	volumeMounts, volumes := paastaVolumesToKubernetesVolumes(fakeVols)
	expectedMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      "tmp-rw",
			MountPath: "/tmp/bar",
			ReadOnly:  false,
		},
		corev1.VolumeMount{
			Name:      "tmp-ro",
			MountPath: "/tmp/bar",
			ReadOnly:  true,
		},
	}

	expectedVolumes := []corev1.Volume{
		{
			Name: "tmp-rw",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/tmp/rw",
				},
			},
		},
		{
			Name: "tmp-ro",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/tmp/ro",
				},
			},
		},
	}

	if !reflect.DeepEqual(volumeMounts, expectedMounts) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v", expectedMounts, volumeMounts)
	}
	if !reflect.DeepEqual(volumes, expectedVolumes) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v", expectedVolumes, volumes)
	}
}

func TestFormatMountName(t *testing.T) {
	in := "/var/my_mount/"
	out := "var-my--mount"
	actual := formatMountName(in)
	if actual != out {
		t.Errorf("Expected:\n%+v\nGot:\n%+v", out, actual)
	}
}

func TestGetDefaultPaastaKubernetesVolumes(t *testing.T) {
	fakeVolumeConfig := map[string]interface{}{
		"volumes": []map[string]interface{}{
			map[string]interface{}{
				"hostPath":      "/foo",
				"containerPath": "/bar",
				"mode":          "RO",
			},
		},
	}
	reader := &config.Store{Data: fakeVolumeConfig}
	volumeMounts, volumes, err := GetDefaultPaastaKubernetesVolumes(reader)
	if err != nil {
		t.Errorf("Error %s", err)
	}
	expectedMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      "foo",
			MountPath: "/bar",
			ReadOnly:  true,
		},
	}

	expectedVolumes := []corev1.Volume{
		{
			Name: "foo",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/foo",
				},
			},
		},
	}
	if !reflect.DeepEqual(volumeMounts, expectedMounts) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v", expectedMounts, volumeMounts)
	}
	if !reflect.DeepEqual(volumes, expectedVolumes) {
		t.Errorf("Expected:\n%+v\nGot:\n%+v", expectedVolumes, volumes)
	}
}
