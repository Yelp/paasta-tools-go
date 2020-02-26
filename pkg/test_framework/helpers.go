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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

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
			delete(section, v)
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
	panic("Unreachable code in DeleteValue")
}

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

func LoadInto(r io.Reader, into interface{}) error {
	if err := yaml.NewYAMLOrJSONDecoder(r, bytes.MinRead).Decode(into); err != nil {
		return err
	}
	return nil
}

type WaitSource func()(runtime.Object, error)

func WaitForNone(timeout time.Duration, from WaitSource) error {
	none := int32(0)
	return WaitFor(&none, timeout, from)
}

func WaitFor(reps* int32, timeout time.Duration, from WaitSource) error {
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


func getReady(obj runtime.Object) (int32, error) {
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
	default:
		log.Panicf("Unsupported type %v", t.GetObjectKind())
	}
	return 0, nil
}
