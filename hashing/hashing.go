package hashing

import (
	"fmt"
	"github.com/fatih/structs"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
	"reflect"
	"strings"
)

func ComputeHashForKubernetesObject(object interface{}) (string, error) {
	value := reflect.ValueOf(object)
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
			// equivalent of delete(labels, "yelp.com/operator_config_hash")
			// we remove this so that the already hashed versions will match
			// the newly calculated versions
			labels.SetMapIndex(reflect.ValueOf("yelp.com/operator_config_hash"), reflect.Value{})
			for _, k := range labels.MapKeys() {
				v := labels.MapIndex(k)
				labelsToHash[k.Interface().(string)] = v.Interface().(string)
			}

		}
	}
	// handy library to turn any struct into map[string]interface{}
	// so we can easily get the SomethingSpec
	mapOfObject := structs.Map(object)
	mapToHash := make(map[string]interface{})
	mapToHash["Labels"] = labelsToHash
	for k, v := range mapOfObject {
		// we match on Spec suffix since we don't know if this has
		// DeploymentSpec PodSpec StatefulSetSpec...
		if strings.HasSuffix(k, "Spec") {
			mapToHash["Spec"] = v
		}
	}
	hasher := fnv.New32a()
	hashutil.DeepHashObject(hasher, mapToHash)
	return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32())), nil
}
