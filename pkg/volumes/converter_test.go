package volumes

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	corev1 "k8s.io/api/core/v1"
)

func TestPaastaVolumesToKubernetesVolumes(t *testing.T) {
	fakeVols := []Volume{
		{
			HostPath:      "/tmp/rw",
			ContainerPath: "/tmp/bar",
			Mode:          "RW",
		},
		{
			HostPath:      "/tmp/ro",
			ContainerPath: "/tmp/bar",
			Mode:          "RO",
		},
	}
	volumeMounts, volumes := paastaVolumesToKubernetesVolumes(fakeVols)
	expectedMounts := []corev1.VolumeMount{
		{
			Name:      "tmp-rw",
			MountPath: "/tmp/bar",
			ReadOnly:  false,
		},
		{
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
	fakeVolumeConfig := &sync.Map{}
	fakeVolumeConfig.Store("volumes", []map[string]string{
		{
			"hostPath":      "/foo",
			"containerPath": "/bar",
			"mode":          "RO",
		},
	},
	)
	reader := &configstore.Store{Data: fakeVolumeConfig}
	volumeMounts, volumes, err := GetDefaultPaastaKubernetesVolumes(reader)
	if err != nil {
		t.Errorf("Error %s", err)
	}
	expectedMounts := []corev1.VolumeMount{
		{
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

func TestGetDefaultPaastaKubernetesHealthcheckVolumes(t *testing.T) {
	fakeVolumeConfig := &sync.Map{}
	fakeVolumeConfig.Store("hacheck_sidecar_volumes", []map[string]string{{
		"hostPath":      "/foo1",
		"containerPath": "/bar1",
		"mode":          "RW",
	},
	},
	)
	reader := &configstore.Store{Data: fakeVolumeConfig}
	volumeMounts, volumes, err := GetDefaultPaastaKubernetesHealthcheckVolumes(reader)
	if err != nil {
		t.Errorf("Error %s", err)
	}
	expectedMounts := []corev1.VolumeMount{
		{
			Name:      "foo1",
			MountPath: "/bar1",
			ReadOnly:  false,
		},
	}

	expectedVolumes := []corev1.Volume{
		{
			Name: "foo1",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/foo1",
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
