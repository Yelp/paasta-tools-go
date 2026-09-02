package framework

import (
	"context"
	"net"
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
	// k3d needs a working Docker daemon — skip if it's not reachable
	conn, err := net.DialTimeout("unix", "/var/run/docker.sock", 2*time.Second)
	if err != nil {
		t.Skip("Docker daemon not available, skipping k3d test")
		return
	}
	conn.Close()
	options := *newOptions(DefaultPrefix("itest"))
	Start(options, nil, nil)
	defer Kube.Close()

	// See if the kubernetes service is available
	err := WaitFor(
		1,
		time.Minute,
		func() (interface{}, error) {
			res := &corev1.ServiceList{}
			err := Kube.Client.List(
				context.TODO(),
				res,
				&client.ListOptions{Namespace: "default"},
			)
			return res, err
		},
	)
	assert.NoError(t, err)
}
