package utils

import (
	assert "github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetYamlDiffForObjects(t *testing.T) {
	labels1 := map[string]string{
		"yelp.com/rick": "andmortyadventures1",
		"yelp.com/operator_config_hash": "somerandomhash1",
	}
	replicas1 := int32(2)
	someStatefulSet1 := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster1",
			Namespace: "paasta-cassandra",
			Labels:    labels1,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &replicas1,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{},
					Containers: []corev1.Container{},
				},
			},
		},
	}

	labels2 := map[string]string{
		"yelp.com/rick": "andmortyadventures2",
		"yelp.com/operator_config_hash": "somerandomhash2",
	}
	replicas2 := int32(2)
	someStatefulSet2 := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "morty-test-cluster2",
			Namespace: "paasta-cassandra",
			Labels:    labels2,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &replicas2,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{corev1.PersistentVolumeClaim{}},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:    []corev1.Volume{
						{
							Name:         "volume1",
							VolumeSource: corev1.VolumeSource{},
						},
					},
					Containers: []corev1.Container{},
				},
			},
		},
	}
	expectedDiff := `--- Old
+++ New
@@ -5,5 +5,5 @@
   labels:
-    yelp.com/operator_config_hash: somerandomhash1
-    yelp.com/rick: andmortyadventures1
-  name: morty-test-cluster1
+    yelp.com/operator_config_hash: somerandomhash2
+    yelp.com/rick: andmortyadventures2
+  name: morty-test-cluster2
   namespace: paasta-cassandra
@@ -18,2 +18,4 @@
       containers: []
+      volumes:
+      - name: volume1
   updateStrategy: {}
`
	actualDiff, _ := GetYamlDiffForObjects(someStatefulSet1, someStatefulSet2)
	assert.Equal(t, expectedDiff, actualDiff)
}
