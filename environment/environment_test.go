package environment

import (
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"testing"
)

func TestGetDefaultPaastaKubernetesEnvironment(test *testing.T) {
	actual := GetDefaultPaastaKubernetesEnvironment()
	fakeEnvironment := []corev1.EnvVar{
		{
			Name: "PAASTA_POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name: "PAASTA_SERVICE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.labels['yelp.com/paasta_service']",
				},
			},
		},
		{
			Name: "PAASTA_INSTANCE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.labels['yelp.com/paasta_instance']",
				},
			},
		},
	}

	if !reflect.DeepEqual(actual, fakeEnvironment) {
		test.Errorf("Expected:\n%+v\nGot:\n%+v", actual, fakeEnvironment)
	}
}
