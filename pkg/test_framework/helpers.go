package framework

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
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

type WaitForFn func()(interface{}, error)

// Wait until WaitSource returns 0 elements or NotFound error - useful for waiting until
// something is deleted
func WaitForNone(timeout time.Duration, from WaitForFn) error {
	none := int32(0)
	return WaitFor(&none, timeout, from)
}

// Wait until from returns a given number of instances. In the context of strongly
// typed non-List k8s objects this typically means number of ready pods; for List objects
// (including UnstructuredList) this means number of items; for Unstructured object
// this means hardcoded 1. For details, refer to getReady function.
// We use a pointer int32 to allow the use of ...Spec.Replicas; for the same reason we
// fallback to hardcoded 1 if nil is provided
func WaitFor(reps* int32, timeout time.Duration, from WaitForFn) error {
	wanted := int32(1)
	if reps != nil {
		wanted = *reps
	}

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := from()
		ready := int32(0)
		if err != nil {
			if !errors.IsNotFound(err) {
				return false, err
			}
			// else let's stick with ready = 0
		} else {
			ready, err = getReady(current)
			if err != nil {
				return false, err
			}
		}
		if ready == wanted {
			return true, nil
		}
		return false, nil
	})
}

func getReady(obj interface{}) (int32, error) {
	switch t := obj.(type) {
	case *appsv1.StatefulSet:
		return (*t).Status.ReadyReplicas, nil
	case *appsv1.Deployment:
		return (*t).Status.ReadyReplicas, nil
	case *appsv1.DaemonSet:
		return (*t).Status.NumberReady, nil
	case *appsv1.ReplicaSet:
		return (*t).Status.ReadyReplicas, nil
	case *batchv1.Job:
		return (*t).Status.Active, nil
	case *appsv1.StatefulSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.DeploymentList:
		return int32(len((*t).Items)), nil
	case *appsv1.DaemonSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.ReplicaSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.ControllerRevisionList:
		return int32(len((*t).Items)), nil
	case *batchv1.JobList:
		return int32(len((*t).Items)), nil
	case *corev1.PersistentVolumeList:
		return int32(len((*t).Items)), nil
	case *corev1.PersistentVolumeClaimList:
		return int32(len((*t).Items)), nil
	case *corev1.PodList:
		return int32(len((*t).Items)), nil
	case *corev1.ServiceList:
		return int32(len((*t).Items)), nil
	case *corev1.ServiceAccountList:
		return int32(len((*t).Items)), nil
	case *corev1.EndpointsList:
		return int32(len((*t).Items)), nil
	case *corev1.NodeList:
		return int32(len((*t).Items)), nil
	case *corev1.NamespaceList:
		return int32(len((*t).Items)), nil
	case *corev1.EventList:
		return int32(len((*t).Items)), nil
	case *corev1.SecretList:
		return int32(len((*t).Items)), nil
	case *corev1.ConfigMapList:
		return int32(len((*t).Items)), nil
	case *corev1.ComponentStatusList:
		return int32(len((*t).Items)), nil
	case *rbacv1.RoleBindingList:
		return int32(len((*t).Items)), nil
	case *rbacv1.RoleList:
		return int32(len((*t).Items)), nil
	case *rbacv1.ClusterRoleBindingList:
		return int32(len((*t).Items)), nil
	case *rbacv1.ClusterRoleList:
		return int32(len((*t).Items)), nil
	case *unstructured.UnstructuredList:
		return int32(len((*t).Items)), nil
	case *unstructured.Unstructured:
		if !t.IsList() {
			// Consider single object an equivalent for a list of 1
			return 1, nil
		}
		list, err := t.ToList()
		if err != nil {
			log.Panic(err)
		}
		return int32(len(list.Items)), nil
	case int64:
		return int32(t), nil
	case int32:
		return t, nil
	case int:
		return int32(t),nil
	default:
		log.Panicf("Unsupported type %t", obj)
	}
	return 0, nil
}
