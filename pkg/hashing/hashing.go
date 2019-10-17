package hashing

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
	"reflect"
)

func ComputeHashForKubernetesObject(object interface{}) (string, error) {
	// By marshaling/unmarshaling the object via JSON we're copying it into a
	// map, which is easier to manipulate in generic way than structs.
	if b, err := json.Marshal(object); err != nil {
		return "", fmt.Errorf("Error while encoding %+v into JSON: %s", object, err)
	} else {
		var v map[string]interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return "", fmt.Errorf("Error while decoding JSON %s into an object: %s", v, err)
		} else {
			// We need only kind/version/spec and labels excluding the label with the
			// current hash value while calculating the hash.  Also Kubernetes adds
			// its own info into `metadata` which we need to ignore.
			meta := v["metadata"].(map[string]interface{})
			labels := meta["labels"]
			if labels != nil {
				delete(labels.(map[string]interface{}), "yelp.com/operator_config_hash")
			} else {
				labels = make(map[string]interface{})
			}
			m := map[string]interface{}{
				"kind":       v["kind"],
				"apiVersion": v["apiVersion"],
				"spec":       v["spec"],
				"metadata": map[string]interface{}{
					"name":      meta["name"],
					"namespace": meta["namespace"],
					"labels":    labels,
				},
			}
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
}

func SetKubernetesObjectHash(configHash string, object interface{}) error {
	labelsToAdd := map[string]string{
		"yelp.com/operator_config_hash": configHash,
	}
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
