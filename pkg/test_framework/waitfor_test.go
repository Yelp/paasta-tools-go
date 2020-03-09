package framework

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestWaitFor_Basic(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("This test is meant to run on Linux only")
		return
	}
	options := *newOptions("itest")
	Start(options, nil, nil)
	defer Kube.Close()

	// See if the kubernetes service is available
	err := WaitFor(
		1,
		time.Minute,
		func()(interface{}, error) {
			res := &corev1.ServiceList{}
			err := Kube.Client.List(
				context.TODO(),
				&client.ListOptions{Namespace: "default"},
				res,
			)
			return res, err
		},
	)
	assert.NoError(t, err)
}
