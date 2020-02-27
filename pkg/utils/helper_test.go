package utils

import (
	assert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"encoding/json"
)

func TestGetHashObjectOfKubernetes(t *testing.T) {
	labels := map[string]string{
		"yelp.com/rick": "andmortyadventures",
		"yelp.com/operator_config_hash": "somerandomhash",
	}
	labelsWithoutHash := map[string]string{
		"yelp.com/rick": "andmortyadventures",
	}
	replicas := int32(2)
	someStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster",
			Namespace: "paasta-cassandra",
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &replicas,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	expectedHashObject := map[string]interface{}{
		"kind": someStatefulSet.TypeMeta.Kind,
		"apiVersion": someStatefulSet.TypeMeta.APIVersion,
		"spec": someStatefulSet.Spec,
		"metadata": map[string]interface{}{
			"name": someStatefulSet.ObjectMeta.Name,
			"namespace": someStatefulSet.ObjectMeta.Namespace,
			"labels": labelsWithoutHash,
		},
	}
	expectedOutString, err := json.Marshal(expectedHashObject)
	_ = json.Unmarshal(expectedOutString, &expectedHashObject)
	hashObject, err := GetHashObjectOfKubernetes(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash object")
	}
	assert.Equal(t, expectedHashObject, hashObject)
}

func TestGetYamlOfObject(t *testing.T) {
	labels := map[string]string{
		"yelp.com/rick": "andmortyadventures",
		"yelp.com/operator_config_hash": "somerandomhash",
	}
	replicas := int32(2)
	someStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster",
			Namespace: "paasta-cassandra",
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &replicas,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	expectedOutput := `apiVersion: apps/v1
kind: StatefulSet
metadata:
  creationTimestamp: null
  labels:
    yelp.com/operator_config_hash: somerandomhash
    yelp.com/rick: andmortyadventures
  name: morty-test-cluster
  namespace: paasta-cassandra
spec:
  replicas: 2
  selector: null
  serviceName: ""
  template:
    metadata:
      creationTimestamp: null
    spec:
      containers: []
  updateStrategy: {}
  volumeClaimTemplates:
  - metadata:
      creationTimestamp: null
    spec:
      resources: {}
    status: {}
status:
  replicas: 0
`
	actualOutput, _ := GetYamlOfObject(someStatefulSet)
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestGetYamlOfHashObjectOfK8sObject(t *testing.T) {
	labels := map[string]string{
		"yelp.com/rick": "andmortyadventures",
		"yelp.com/operator_config_hash": "somerandomhash",
	}
	replicas := int32(2)
	someStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster",
			Namespace: "paasta-cassandra",
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &replicas,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	expectedOutput := `apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    yelp.com/rick: andmortyadventures
  name: morty-test-cluster
  namespace: paasta-cassandra
spec:
  replicas: 2
  selector: null
  serviceName: ""
  template:
    metadata:
      creationTimestamp: null
    spec:
      containers: []
  updateStrategy: {}
  volumeClaimTemplates:
  - metadata:
      creationTimestamp: null
    spec:
      resources: {}
    status: {}
`
	actualOutput, _ := GetYamlOfHashObjectOfK8sObject(someStatefulSet)
	assert.Equal(t, expectedOutput, actualOutput)
}