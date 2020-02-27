package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

// function to get the relevant information out of k8s object that is needed to compute hash
func GetHashObjectOfKubernetes(object interface{}) (map[string]interface{}, error) {
	// By marshaling/unmarshaling the object via JSON we're copying it into a
	// map, which is easier to manipulate in generic way than structs.
	if b, err := json.Marshal(object); err != nil {
		return nil, fmt.Errorf("Error while encoding %+v into JSON: %s", object, err)
	} else {
		var v map[string]interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, fmt.Errorf("Error while decoding JSON %s into an object: %s", v, err)
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
			return m, nil
		}
	}
}

// function to get the yaml output of an object
func GetYamlOfObject(object interface{}) (string, error) {
	if b, err := json.Marshal(object); err != nil {
		return "", fmt.Errorf("Error while encoding %+v into JSON: %s", object, err)
	} else {
		var v map[string]interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return "", fmt.Errorf("Error while decoding JSON %s into an object: %s", v, err)
		} else {
			if y, err := yaml.Marshal(&v); err != nil {
				return "", fmt.Errorf("Error while encoding Map %s into YAML: %s", v, err)
			} else {
				return string(y), nil
			}
		}
	}
}

// function to get yaml output of hash object for k8s objects
func GetYamlOfHashObjectOfK8sObject(object interface{}) (string, error) {
	if hashObject, err := GetHashObjectOfKubernetes(object); err != nil {
		return "", err
	} else {
		if yamlOutput, err := GetYamlOfObject(hashObject); err != nil {
			return "", err
		} else {
			return yamlOutput, nil
		}
	}
}
