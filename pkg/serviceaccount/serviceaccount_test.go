package serviceaccount

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// NewFakeClient creates a new fake Kubernetes client.
func NewFakeClient(initObjs ...runtime.Object) client.Client {
	return fake.NewClientBuilder().WithRuntimeObjects(initObjs...).WithScheme(clientgoscheme.Scheme).Build()
}

func TestEnsureServiceAccountForIamRoleAddsMissing(t *testing.T) {

	serviceAccountName := "fake-iam-role"
	namespace := "fake-namespace"

	client := NewFakeClient()
	//first, assert that there is nothing there
	result := &corev1.ServiceAccount{}
	err := client.Get(context.Background(),
		types.NamespacedName{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
		result,
	)
	if err != nil {
		if !errors.IsNotFound(err) {
			t.Failed()
		}
	}

	err = EnsureServiceAccountForIamRole(context.TODO(), serviceAccountName, namespace, client)

	err = client.Get(context.Background(),
		types.NamespacedName{
			Name:      serviceAccountName,
			Namespace: namespace,
		},
		result,
	)
	if err != nil {
		if errors.IsNotFound(err) {
			t.Failed()
		}
	}
}
