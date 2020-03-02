package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/yaml.v2"
)

// function to get the yaml output of an object
// marshal/unmarshal to and from json in order to maintain consistency between json and yaml
// currently k8s objects have json field tags, so this helps in skipping empty fields
// and give consistent output yaml similar to json
func getYamlOfObject(object interface{}) (string, error) {
	if b, err := json.Marshal(object); err != nil {
		return "", fmt.Errorf("error while encoding %+v into JSON: %s", object, err)
	} else {
		var v map[string]interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return "", fmt.Errorf("error while decoding JSON %s into an object: %s", v, err)
		} else {
			if y, err := yaml.Marshal(&v); err != nil {
				return "", fmt.Errorf("error while encoding Map %s into YAML: %s", v, err)
			} else {
				return string(y), nil
			}
		}
	}
}

func generateYamlDiff(yaml1 string, yaml2 string, context int) string {
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(yaml1),
		B:        difflib.SplitLines(yaml2),
		FromFile: "Old",
		FromDate: "",
		ToFile:   "New",
		ToDate:   "",
		Context:  context,
	})
	return diff
}

// func to generate yaml diff of objects used for hashing
// context : number of context lines to use for generating diff
// for more reference on context : https://github.com/pmezard/go-difflib/blob/5d4384ee4fb2527b0a1256a821ebfc92f91efefc/difflib/difflib.go#L559
func GetYamlDiffForObjects(objectOld interface{}, objectNew interface{}, context int) (string, error) {
	yamlOld, err := getYamlOfObject(objectOld)
	if err != nil {
		return "", err
	}
	yamlNew, err := getYamlOfObject(objectNew)
	if err != nil {
		return "", err
	}
	return generateYamlDiff(yamlOld, yamlNew, context), nil
}
