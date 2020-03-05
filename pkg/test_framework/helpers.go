package framework

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Read an element (value or section) from the unstructured object, given a path
// consisting of nested element names e.g. "metadata", "labels", "some_label". Returns
// an error if element is not found or if a nesting element is not a section
// By "section" I mean map[string]interface{} i.e. map capable of nesting elements
func ReadValue(obj *unstructured.Unstructured, path ...string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	errpath := ""
	errjson, _ := obj.MarshalJSON()
	section := obj.UnstructuredContent()
	for i, v := range path {
		errpath += "[" + v + "]"

		if i+1 == len(path) {
			result, ok := section[v]
			if !ok {
				return nil, fmt.Errorf("could not find value %s in:\n%s", errpath, string(errjson))
			} else {
				return result, nil
			}
		} else {
			tmp, ok := section[v]
			if !ok {
				return nil, fmt.Errorf("could not find section %s in:\n%s", errpath, string(errjson))
			}
			section, ok = tmp.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("found %s value, but section was expected in:\n%s", errpath, string(errjson))
			}
		}
	}
	panic("Unreachable code in ReadValue")
}

// Write an element (value or section) from the unstructured object, given a path
// consisting of nested element names e.g. "metadata", "labels", "some_label". Returns
// an error if nesting element is not found or is not a section
// By "section" I mean map[string]interface{} i.e. map capable of nesting elements
func WriteValue(obj *unstructured.Unstructured, value interface{}, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}
	errpath := ""
	errjson, _ := obj.MarshalJSON()
	section := obj.UnstructuredContent()
	for i, v := range path {
		errpath += "[" + v + "]"

		if i+1 == len(path) {
			section[v] = value
			return nil
		} else {
			tmp, ok := section[v]
			if !ok {
				return fmt.Errorf("could not find section %s in:\n%s", errpath, string(errjson))
			}
			section, ok = tmp.(map[string]interface{})
			if !ok {
				return fmt.Errorf("found %s value, but section was expected in:\n%s", errpath, string(errjson))
			}
		}
	}
	panic("Unreachable code in WriteValue")
}

// Delete an element (value or section) from the unstructured object, given a path
// consisting of nested element names e.g. "metadata", "labels", "some_label". Returns
// an error if nesting element is not a section. There is no error if the element
// being deleted does not exist, because that would be a no-op anyway.
// By "section" I mean map[string]interface{} i.e. map capable of nesting elements
func DeleteValue(obj *unstructured.Unstructured, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("empty path")
	}
	errpath := ""
	errjson, _ := obj.MarshalJSON()
	section := obj.UnstructuredContent()
	for i, v := range path {
		errpath += "[" + v + "]"

		if i+1 == len(path) {
			// Note: delete is a no-op if v is not in section
			delete(section, v)
			return nil
		} else {
			tmp, ok := section[v]
			if !ok {
				// The element we wanted to delete does not exist, so no-op
				return nil
			}
			section, ok = tmp.(map[string]interface{})
			if !ok {
				return fmt.Errorf("found %s value, but section was expected in:\n%s", errpath, string(errjson))
			}
		}
	}
	panic("Unreachable code in DeleteValue")
}

// Load an Unstructured object from the textual data, either in JSON or YAML format
// using standard Reader interface (i.e. could be a file, hardcoded text, something else).
// Note: Unstructured object supports conversion to UnstructuredList for any object
// which contains "items" directly under the root
func LoadUnstructured(r io.Reader) (*unstructured.Unstructured, error) {
	reader, _, isJson := yaml.GuessJSONStream(r, bytes.MinRead)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if !isJson {
		tmp, err := yaml.ToJSON(data)
		if err != nil {
			return nil, err
		}
		data = tmp
	}

	result := unstructured.Unstructured{}
	err = result.UnmarshalJSON(data)
	return &result, err
}

// Load either strongly typed k8s object, or Unstructured object, from textual data
// either in JSON or YAML format.
func LoadInto(r io.Reader, into interface{}) error {
	if err := yaml.NewYAMLOrJSONDecoder(r, bytes.MinRead).Decode(into); err != nil {
		return err
	}
	return nil
}
