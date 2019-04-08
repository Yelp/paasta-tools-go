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
	hash, err := ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "76ffc95c66")

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
			Replicas:             &theSameReplicas,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
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
	assert.Equal(t, anotherHash, "76ffc95c66")

	// test the hash changes if we change replicas
	replicas = int32(1)
	someStatefulSet.Spec.Replicas = &replicas
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "59cb75c79c")

	// test hash changes if we change labels
	someStatefulSet.ObjectMeta.Labels["yelp.com/for"] = "everandever"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "9c45f99")

	// test hash ignores yelp.com/operator_config_hash label
	someStatefulSet.ObjectMeta.Labels["yelp.com/operator_config_hash"] = "somehash"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "9c45f99")

	// test we get same hash if we pass the actual struct not a pointer
	hash, err = ComputeHashForKubernetesObject(*someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, hash, "9c45f99")
}

func TestAddLabelToMetadata(t *testing.T) {
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

	labelToAdd := map[string]string{"yelp.com/malcom": "tucker"}
	err := AddLabelsToMetadata(labelToAdd, someStatefulSet)
	if err != nil {
		t.Errorf("Failed to add label")
	}

	// the new label and existing label are present in ObjectMeta
	assert.Equal(t, someStatefulSet.ObjectMeta.Labels["yelp.com/malcom"], "tucker")
	assert.Equal(t, someStatefulSet.ObjectMeta.Labels["yelp.com/rick"], "andmortyadventures")

	// the new label is *not* present on other parts of the kubernetes object
	assert.NotEqual(t, someStatefulSet.Spec.Selector.MatchLabels["yelp.com/malcom"], "tucker")
	assert.NotEqual(t, someStatefulSet.Spec.Template.ObjectMeta.Labels["yelp.com/malcom"], "tucker")
	// but the existing labels are
	assert.Equal(t, someStatefulSet.Spec.Selector.MatchLabels["yelp.com/rick"], "andmortyadventures")
}
