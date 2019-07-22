package hashing

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	copy := someStatefulSet.DeepCopyObject()
	hash, err := ComputeHashForKubernetesObject(someStatefulSet)

	assert.True(t, assert.ObjectsAreEqual(copy, someStatefulSet), "expected object to be unchanged")

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
	assert.Equal(t, "7f8447c49c", anotherHash)

	// test the hash changes if we change replicas
	replicas = int32(1)
	someStatefulSet.Spec.Replicas = &replicas
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, "787785c784", hash)

	// test hash changes if we change labels
	someStatefulSet.ObjectMeta.Labels["yelp.com/for"] = "everandever"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, "676658c9dd", hash)

	// test hash ignores yelp.com/operator_config_hash label
	someStatefulSet.ObjectMeta.Labels["yelp.com/operator_config_hash"] = "somehash"
	hash, err = ComputeHashForKubernetesObject(someStatefulSet)
	if err != nil {
		t.Errorf("Failed to calculate hash")
	}
	assert.Equal(t, "676658c9dd", hash)
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

func TestUniqueResourceReqs(t *testing.T) {
	labels := map[string]string{
		"yelp.com/rick": "andmortyadventures",
	}
	oneCPU, _ := resource.ParseQuantity(string("0.1"))
	twoCPU, _ := resource.ParseQuantity(string("0.2"))
	replicas := int32(2)
	deploymentOne := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster",
			Namespace: "paasta-cassandra",
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{},
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: oneCPU,
								},
							},
						},
					},
				},
			},
		},
	}
	deploymentTwo := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster",
			Namespace: "paasta-cassandra",
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{},
					Containers: []corev1.Container{
						corev1.Container{
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: twoCPU,
								},
							},
						},
					},
				},
			},
		},
	}
	hashOne, _ := ComputeHashForKubernetesObject(deploymentOne)
	fmt.Println(hashOne)
	hashTwo, _ := ComputeHashForKubernetesObject(deploymentTwo)
	fmt.Println(hashTwo)
	assert.NotEqual(t, hashOne, hashTwo)
}
