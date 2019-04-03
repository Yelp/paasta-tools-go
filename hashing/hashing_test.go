package hashing

import (
	assert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetHashForKubernetesObject(t *testing.T) {
	labels := map[string]string{
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
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	hash, err := ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "f968dcd9f")

	// to test that a new pointer to a semantically matching object has the same hash
	theSameReplicas := int32(2)
	theSameStatefulSet := &appsv1.StatefulSet{
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
			Replicas: &theSameReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	anotherHash, err := ComputeHashForKubernetesObject(theSameStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, anotherHash, "f968dcd9f")

	// test the hash changes if we change replicas
	replicas = int32(1)
	someStatefulSet.Spec.Replicas = &replicas
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "569c7fc7d4")

	// test hash changes if we change labels
	someStatefulSet.ObjectMeta.Labels["yelp.com/for"] = "everandever"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "85445d55bb")

	// test hash ignores yelp.com/operator_config_hash label
	someStatefulSet.ObjectMeta.Labels["yelp.com/operator_config_hash"] = "somehash"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "85445d55bb")

	// test we get same hash if we pass the actual struct not a pointer
	hash, err = ComputeHashForKubernetesObject(*someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "85445d55bb")
}
