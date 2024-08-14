package environment

import (
	corev1 "k8s.io/api/core/v1"
)

func GetDefaultPaastaKubernetesEnvironment() []corev1.EnvVar {
	defaultEnvironment := []corev1.EnvVar{
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
					FieldPath: "metadata.labels['paasta.yelp.com/service']",
				},
			},
		},
		{
			Name: "PAASTA_INSTANCE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.labels['paasta.yelp.com/instance']",
				},
			},
		},
		{
			Name: "PAASTA_INSTANCE_TYPE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.labels['paasta.yelp.com/service']",
				},
			},
		},
		{
			Name: "PAASTA_CLUSTER",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.labels['paasta.yelp.com/cluster']",
				},
			},
		},
		{
			Name: "PAASTA_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		},
	}
	return defaultEnvironment
}
