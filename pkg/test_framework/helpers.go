package framework

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ReadValue (element, array or section) from an unstructured object
//
// Requires a path consisting of nested element names e.g. "metadata", "labels", "some_label". Returns an error if
// element is not found or if a nesting element is not a section. By "section" I mean map[string]interface{} i.e. map
// capable of nesting elements. No support for addressing individual elements in an array (slice).
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

// WriteValue (element, array or section) from an unstructured object
//
// Requires a path consisting of nested element names e.g. "metadata", "labels", "some_label". Returns an error if
// nesting element is not found or is not a section. By "section" I mean map[string]interface{} i.e. map
// capable of nesting elements. No support for addressing individual elements in an array (slice).
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

// DeleteValue (element, array or section) from the unstructured object
//
// Requires a path consisting of nested element names e.g. "metadata", "labels", "some_label". Returns an error if
// nesting element is not a section. There is no error if the element being deleted does not exist, because that would
// be a no-op anyway. By "section" I mean map[string]interface{} i.e. map capable of nesting elements. No support for
// addressing individual elements in an array (slice).
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

// LoadUnstructured load an object from the textual data
//
// Using standard Reader interface (i.e. could be a file, hardcoded text, something else). Note: Unstructured object
// supports conversion to UnstructuredList for any object which contains "items" map directly under the object root.
// Both JSON and YAML formats are supported.
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

// LoadInto load an object from the textual data into an existing K8s object
//
// This function supports both strong-typed object and Unstructured object. For example to load a deployment file
// the following can be used:
//
//    dep := appsv1.Deployment{}
//    if err := framework.LoadInto(bufio.NewReader(file), &dep); err != nil { . . .
//
// The benefits of using strong-typed objects (rather than unstructured.Unstructured) are:
//   * K8s will validate the input file while loading it
//   * Individual elements of the service definition will be easier to work with
//   * Function framework.WaitFor() provides special handling waiting on pods for types which can own them
//
// Both JSON and YAML formats are supported.
func LoadInto(r io.Reader, into interface{}) error {
	if err := yaml.NewYAMLOrJSONDecoder(r, bytes.MinRead).Decode(into); err != nil {
		return err
	}
	return nil
}
