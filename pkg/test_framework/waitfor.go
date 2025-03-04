package framework

import (
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type WaitForFn func() (interface{}, error)

// WaitFor waits until "from" function returns a given number of instances
//
// If the "from" function returns a strongly-typed K8s object of a type which can own pods, this number of instances
// refers to the number of ready pods. If the "from" function returns any of K8s List objects (including
// unstructured.UnstructuredList) this number refers to the number of items in the list. For Unstructured object, the
// returned number is a hardcoded 1. Function "from" can also return an integer number, which is useful if it needs
// to perform some specific check e.g. on data inside an object read from the test cluster. Specifics for individual
// K8s types are inside getReady function below.
//
// If the "from" function does not return an expected number of instances within "timeout", this function will
// panic, hence failing test.
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
