package hashing

import (
	"encoding/json"
	"fmt"
	utils "github.com/Yelp/paasta-tools-go/pkg/utils"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
	"reflect"
)

func ComputeHashForKubernetesObject(object interface{}) (string, error) {
	if m, err := utils.GetHashObjectOfKubernetes(object); err != nil {
		return "", err
	} else {
		// By using serialized JSON for hashing we're making the hashing process
		// a bit easier (like having maps always being sorted by keys).
		if b, err := json.Marshal(m); err != nil {
			return "", fmt.Errorf("Error while encoding %+v into JSON: %s", m, err)
		} else {
			hasher := fnv.New32a()
			hasher.Write(b)
			return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32())), nil
		}
	}
}

func SetKubernetesObjectHash(configHash string, object interface{}) error {
	labelsToAdd := map[string]string{}
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
			// Set the configHash label last, so that any prior configHash label is overwritten with the new one
			labelsToAdd["yelp.com/operator_config_hash"] = configHash
			labels.Set(reflect.ValueOf(labelsToAdd))
		}
	}
	return nil
}
