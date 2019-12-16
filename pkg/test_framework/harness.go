package framework

import (
	harness "github.com/dlespiau/kube-test-harness"
)

type Harness struct {
	harness.Harness
	Options Options
}

type Options struct {
	harness.Options
	Makefile string
	MakeDir  string
}
