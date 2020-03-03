package framework

import (
	"bytes"
	"context"
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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// Wait until from returns a given number of instances. In the context of strongly
// typed non-List k8s objects this typically means number of ready pods; for List objects
// (including UnstructuredList) this means number of items; for Unstructured object
// this means hardcoded 1. For details, refer to getReady function.
func WaitFor(wanted int, timeout time.Duration, from WaitForFn) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := from()
		ready := 0
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

func getReady(obj interface{}) (int, error) {
	switch t := obj.(type) {
	case *appsv1.StatefulSet:
		return int((*t).Status.ReadyReplicas), nil
	case *appsv1.Deployment:
		return int((*t).Status.ReadyReplicas), nil
	case *appsv1.DaemonSet:
		return int((*t).Status.NumberReady), nil
	case *appsv1.ReplicaSet:
		return int((*t).Status.ReadyReplicas), nil
	case *batchv1.Job:
		return int((*t).Status.Active), nil
	case *appsv1.StatefulSetList:
		return len((*t).Items), nil
	case *appsv1.DeploymentList:
		return len((*t).Items), nil
	case *appsv1.DaemonSetList:
		return len((*t).Items), nil
	case *appsv1.ReplicaSetList:
		return len((*t).Items), nil
	case *appsv1.ControllerRevisionList:
		return len((*t).Items), nil
	case *batchv1.JobList:
		return len((*t).Items), nil
	case *corev1.PersistentVolumeList:
		return len((*t).Items), nil
	case *corev1.PersistentVolumeClaimList:
		return len((*t).Items), nil
	case *corev1.PodList:
		return len((*t).Items), nil
	case *corev1.ServiceList:
		return len((*t).Items), nil
	case *corev1.ServiceAccountList:
		return len((*t).Items), nil
	case *corev1.EndpointsList:
		return len((*t).Items), nil
	case *corev1.NodeList:
		return len((*t).Items), nil
	case *corev1.NamespaceList:
		return len((*t).Items), nil
	case *corev1.EventList:
		return len((*t).Items), nil
	case *corev1.SecretList:
		return len((*t).Items), nil
	case *corev1.ConfigMapList:
		return len((*t).Items), nil
	case *corev1.ComponentStatusList:
		return len((*t).Items), nil
	case *rbacv1.RoleBindingList:
		return len((*t).Items), nil
	case *rbacv1.RoleList:
		return len((*t).Items), nil
	case *rbacv1.ClusterRoleBindingList:
		return len((*t).Items), nil
	case *rbacv1.ClusterRoleList:
		return len((*t).Items), nil
	case *unstructured.UnstructuredList:
		return len((*t).Items), nil
	case *unstructured.Unstructured:
		if !t.IsList() {
			// Consider single object an equivalent for a list of 1
			return 1, nil
		}
		list, err := t.ToList()
		if err != nil {
			log.Panic(err)
		}
		return len(list.Items), nil
	case int64:
		return int(t), nil
	case int32:
		return int(t), nil
	case int:
		return t, nil
	default:
		log.Panicf("Unsupported type %t", obj)
	}
	return 0, nil
}

type WaitForResourceFn func(obj unstructured.Unstructured)(int, error)

// Wait for a singular resource like the one provided. Will use APIVersion and Kind from
// like parameter. If name parameter is provided, will use Name and Namespace from name
// parameter, otherwise will use Name and Namespace from like as well.
// Optionally takes a function which can be used to transform a single resource into
// a number (e.g. acting on some part of the object's spec or status section).
func WaitForResource(wanted int, timeout time.Duration, client1 client.Client, like runtime.Object, name runtime.Object, fn WaitForResourceFn) error {
	if name == nil {
		name = like
	}
	apiver, kind := like.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
	objname, objns := getSingularNameNs(name)

	return WaitFor(
		wanted,
		timeout,
		func() (interface{}, error) {
			res := &unstructured.Unstructured{}
			res.SetAPIVersion(apiver)
			res.SetKind(kind)
			err := client1.Get(context.TODO(), client.ObjectKey{
				Namespace: objns,
				Name: objname,
			}, res)
			if err != nil {
				return nil, err
			}
			if fn == nil {
				return res, nil
			}
			return fn(*res)
		},
	)
}

type WaitForResourcesFn func(obj unstructured.UnstructuredList)(int, error)

// Wait for a list of resources like the one provided. Will use APIVersion and Kind from
// like parameter. If name parameter is provided, will use Namespace (only!) from name
// parameter, otherwise will use Namespace from like as well.
// Optionally takes a function which can be used to transform the list into
// a number (e.g. filtering by some part of the objects' name or status).
//
// Note: lists of kubernetes resources do not carry a singular Namespace, so if name
// is nil AND like is a list of objects, you will get a panic from getSingularNameNs.
// This can be avoided either by passing a name parameter, or by passing a singular
// (rather than a list) resource to like parameter, with namespace set.
// Example waiting for a number of pods in someNamespace:
//   pod := corev1.Pod{}
//   pod.SetNamespace(someNamespace)
//   err := WaitForResources(howManyPodsWanted, someTimeout, someClient, pod, nil, nil)
func WaitForResources(wanted int, timeout time.Duration, client1 client.Client, like runtime.Object, name runtime.Object, fn WaitForResourcesFn) error {
	if name == nil {
		name = like
	}
	apiver, kind := like.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
	_, objns := getSingularNameNs(name)

	return WaitFor(
		wanted,
		timeout,
		func() (interface{}, error) {
			res := &unstructured.UnstructuredList{}
			res.SetAPIVersion(apiver)
			res.SetKind(kind)
			err := client1.List(context.TODO(), &client.ListOptions{
				Namespace: objns,
			}, res)
			if err != nil {
				return nil, err
			}
			if fn == nil {
				return res, nil
			}
			return fn(*res)
		},
	)
}

func getSingularNameNs(obj runtime.Object) (string, string) {
	switch t := obj.(type) {
	case *appsv1.StatefulSet:
		return (*t).GetName(), (*t).GetNamespace()
	case *appsv1.Deployment:
		return (*t).GetName(), (*t).GetNamespace()
	case *appsv1.DaemonSet:
		return (*t).GetName(), (*t).GetNamespace()
	case *appsv1.ReplicaSet:
		return (*t).GetName(), (*t).GetNamespace()
	case *batchv1.Job:
		return (*t).GetName(), (*t).GetNamespace()
	case *appsv1.ControllerRevision:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.PersistentVolume:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.PersistentVolumeClaim:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Pod:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Service:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.ServiceAccount:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Endpoints:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Node:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Namespace:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Event:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.Secret:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.ConfigMap:
		return (*t).GetName(), (*t).GetNamespace()
	case *corev1.ComponentStatus:
		return (*t).GetName(), (*t).GetNamespace()
	case *rbacv1.RoleBinding:
		return (*t).GetName(), (*t).GetNamespace()
	case *rbacv1.Role:
		return (*t).GetName(), (*t).GetNamespace()
	case *rbacv1.ClusterRoleBinding:
		return (*t).GetName(), (*t).GetNamespace()
	case *rbacv1.ClusterRole:
		return (*t).GetName(), (*t).GetNamespace()
	case *unstructured.Unstructured:
		return (*t).GetName(), (*t).GetNamespace()
	default:
		log.Panicf("Unsupported type %t", obj)
	}
	return "", ""
}
