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
	originalHash, err := ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}

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
	assert.Equal(t, originalHash, anotherHash)

	// test the hash changes if we change replicas
	replicas = int32(1)
	someStatefulSet.Spec.Replicas = &replicas
	var changedReplicasHash string
	changedReplicasHash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.NotEqual(t, originalHash, changedReplicasHash)

	// test hash changes if we change labels
	someStatefulSet.ObjectMeta.Labels["yelp.com/for"] = "everandever"
	var changedLabelHash string
	changedLabelHash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.NotEqual(t, originalHash, changedLabelHash)
	assert.NotEqual(t, changedReplicasHash, changedLabelHash)

	// test hash ignores yelp.com/operator_config_hash label
	someStatefulSet.ObjectMeta.Labels["yelp.com/operator_config_hash"] = "somehash"
	var changedSpecialLabelHash string
	changedSpecialLabelHash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, changedLabelHash, changedSpecialLabelHash)

	// test we get same hash if we pass the actual struct not a pointer
	var ptrHash string
	ptrHash, err = ComputeHashForKubernetesObject(*someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, changedSpecialLabelHash, ptrHash)
}

func TestSetKubernetesObjectHash(t *testing.T) {
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

	err := SetKubernetesObjectHash("abc1234", someStatefulSet)
	if err != nil {
		t.Errorf("Failed to add label")
	}

	// the new label and existing label are present in ObjectMeta
	assert.Equal(t, someStatefulSet.ObjectMeta.Labels["yelp.com/operator_config_hash"], "abc1234")
	assert.Equal(t, someStatefulSet.ObjectMeta.Labels["yelp.com/rick"], "andmortyadventures")

	// the new label is *not* present on other parts of the kubernetes object
	assert.NotEqual(t, someStatefulSet.Spec.Selector.MatchLabels["yelp.com/operator_config_hash"], "abc1234")
	assert.NotEqual(t, someStatefulSet.Spec.Template.ObjectMeta.Labels["yelp.com/operator_config_hash"], "abc1234")
	// but the existing labels are
	assert.Equal(t, someStatefulSet.Spec.Selector.MatchLabels["yelp.com/rick"], "andmortyadventures")
}
