package framework

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartQuick(t *testing.T) {
	options := *newOptions()
	kube := startHarness(options)
	assert.NotNil(t, kube)
	test := kube.NewTest(t)
	defer test.Close()
	err := test.StartOperator()
	// error because make tests-operator-start is not blocking
	assert.NotNil(t, err)
}

func TestStartSlow(t *testing.T) {
	options := *newOptions()
	options.Prefix = "test-sleep25-"
	kube := startHarness(options)
	assert.NotNil(t, kube)
	test := kube.NewTest(t)
	defer test.Close()
	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)
}
