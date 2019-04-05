package hashing

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/mohae/deepcopy"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
	"reflect"
	"strings"
)

func ComputeHashForKubernetesObject(object interface{}) (string, error) {
	// lets be sure we don't mutate the kubernetes object
	copyOfObject := deepcopy.Copy(object)
	value := reflect.ValueOf(copyOfObject)
	var objectMeta reflect.Value
	if value.Kind() == reflect.Ptr {
		objectMeta = value.Elem().FieldByName("ObjectMeta")
	} else if value.Kind() == reflect.Struct {
		objectMeta = value.FieldByName("ObjectMeta")
	} else {
		return "", fmt.Errorf("Must pass Kubernetes Object or pointer to Kubernetes Objcect")
	}
	// recreate the labels map so we can pass it to the
	// DeepHashObject function
	labelsToHash := make(map[string]string)
	if objectMeta.Kind() == reflect.Struct {
		labels := objectMeta.FieldByName("Labels")
		if labels.Kind() == reflect.Map {
			for _, k := range labels.MapKeys() {
				v := labels.MapIndex(k)
				key := k.Interface().(string)
				// we ignore this so that the already hashed versions will match
				// the newly calculated versions
				if key == "yelp.com/operator_config_hash" {
					continue
				}
				labelsToHash[key] = v.Interface().(string)
			}

		}
	}
	// handy library to turn any struct into map[string]interface{}
	// so we can easily get the SomethingSpec
	mapOfObject := structs.Map(copyOfObject)
	mapToHash := make(map[string]interface{})
	mapToHash["Labels"] = labelsToHash
	for k, v := range mapOfObject {
		// we match on Spec suffix since we don't know if this has
		// DeploymentSpec PodSpec StatefulSetSpec...
		if strings.HasSuffix(k, "Spec") {
			mapToHash[k] = v
		}
	}
	hasher := fnv.New32a()
	hashutil.DeepHashObject(hasher, mapToHash)
	return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32())), nil
}

func AddLabelsToMetadata(labelsToAdd map[string]string, object interface{}) error {
	value := reflect.ValueOf(object)
	var objectMeta reflect.Value
	if value.Kind() == reflect.Ptr {
		objectMeta = value.Elem().FieldByName("ObjectMeta")
	} else {
		return fmt.Errorf("Must pass pointer to AddLabelsToMetadata so we can update labels using reflection.")
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
