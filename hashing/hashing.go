package hashing

import (
	"fmt"
	"hash/fnv"
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/hash"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/rand"
)

func ComputeHashForKubernetesObject(object runtime.Object) (string, error) {
	copy := object.DeepCopyObject()
	accessor := meta.NewAccessor()
	labels, err := accessor.Labels(copy)
	if err != nil {
		return "", err
	}
	delete(labels, "yelp.com/operator_config_hash")
	accessor.SetLabels(copy, labels)

	h := fnv.New32a()
	hash.DeepHashObject(h, copy)
	return rand.SafeEncodeString(fmt.Sprint(h.Sum32())), nil
}

func SetKubernetesObjectHash(configHash string, object runtime.Object) error {
	value := reflect.ValueOf(object)
	labelsToAdd := map[string]string{"yelp.com/operator_config_hash": configHash}
	var objectMeta reflect.Value
	if value.Kind() == reflect.Ptr {
		objectMeta = value.Elem().FieldByName("ObjectMeta")
	} else {
		return fmt.Errorf("must pass pointer to AddLabelsToMetadata so we can update labels using reflection")
	}
	if objectMeta.Kind() == reflect.Struct {
		labels := objectMeta.FieldByName("Labels")
		if labels.Kind() == reflect.Map {
			for _, k := range labels.MapKeys() {
				v := labels.MapIndex(k)
				labelsToAdd[k.Interface().(string)] = v.Interface().(string)
			}
			labels.Set(reflect.ValueOf(labelsToAdd))
		}
	}
	return nil
}
