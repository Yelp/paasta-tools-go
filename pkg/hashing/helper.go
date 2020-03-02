package hashing

import (
	"encoding/json"
	"fmt"
)

// function to get the relevant information out of k8s object that is needed to compute hash
func GetFilteredK8sObjectForHashing(object interface{}) (map[string]interface{}, error) {
	// By marshaling/unmarshaling the object via JSON we're copying it into a
	// map, which is easier to manipulate in generic way than structs.
	if b, err := json.Marshal(object); err != nil {
		return nil, fmt.Errorf("error while encoding %+v into JSON: %s", object, err)
	} else {
		var v map[string]interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, fmt.Errorf("error while decoding JSON %s into an object: %s", v, err)
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
