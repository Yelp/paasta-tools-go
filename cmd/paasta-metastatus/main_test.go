package main

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Yelp/paasta-tools-go/pkg/cli"

	"github.com/Yelp/paasta-tools-go/pkg/paasta_api/models"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"

	operations "github.com/Yelp/paasta-tools-go/pkg/paasta_api/client/operations"
)

type MockTransport struct {
	Ops []*runtime.ClientOperation
}

func (m *MockTransport) Submit(co *runtime.ClientOperation) (interface{}, error) {
	m.Ops = append(m.Ops, co)
	return &operations.MetastatusOK{
		Payload: &models.MetaStatus{
			Output: "foo",
		},
	}, nil
}

func (m *MockTransport) Reset() {
	m.Ops = []*runtime.ClientOperation{}
}

func makeTestContext() context.Context {
	t := &MockTransport{}
	err := &strings.Builder{}
	out := &strings.Builder{}
	ctx := context.WithValue(context.Background(), ctxKeyTransport, t)
	ctx = context.WithValue(ctx, ctxKeyOut, out)
	ctx = context.WithValue(ctx, ctxKeyErr, err)
	return ctx
}

func TestMetastatus(test *testing.T) {
	ctx := makeTestContext()
	mockCmdArgs := []string{"foo", "bar"}
	metastatus(
		ctx,
		[]string{"cluster-foo"},
		map[string]string{"cluster-foo": "endpoint-foo"},
		map[string]map[string]interface{}{
			"cluster-foo": {
				"dashboard-foo-1": "dashboard-foo-1-content",
			},
		},
		mockCmdArgs,
	)

	transport := ctx.Value(ctxKeyTransport).(*MockTransport)
	if len(transport.Ops) != 1 {
		test.Logf("opes: %v", transport.Ops)
		test.Errorf("expected number of operations: 1, got: %v", len(transport.Ops))
	}

	metastatusOp := transport.Ops[0]
	if metastatusOp.PathPattern != "/metastatus" {
		test.Errorf("unexpected path: %v", metastatusOp.PathPattern)
	}

	params := metastatusOp.Params.(*operations.MetastatusParams)
	if !reflect.DeepEqual(params.CmdArgs, mockCmdArgs) {
		test.Errorf("expected mock args: %v, actual: %v", mockCmdArgs, params.CmdArgs)
	}

	errf := ctx.Value(ctxKeyErr).(*strings.Builder)
	if errf.String() != "" {
		test.Errorf("error stream not empty: %v", errf)
	}

	outf := ctx.Value(ctxKeyOut).(*strings.Builder)
	outs := outf.String()
	ok, _ := regexp.MatchString(`Cluster: cluster-foo`, outs)
	if !ok {
		test.Errorf("out doesn't match `Cluster: cluster-foo`:\n%v", outs)
	}
}

func Test_writeDashboards(test *testing.T) {
	sb := &strings.Builder{}
	writeDashboards("cluster-foo", nil, sb)
	assert.Regexp(test, `No dashboards configured`, sb.String())

	sb = &strings.Builder{}
	writeDashboards(
		"cluster-foo",
		map[string]interface{}{
			"one":   "two",
			"three": []string{"four"},
			"five":  []string{"six", "seven"},
		},
		sb,
	)
	assert.Regexp(test, `one:.*two`, sb.String())
	assert.Regexp(test, `three:.*four`, sb.String())
	assert.Regexp(test, `five:.*\n.*six.*\n.*seven`, sb.String())
}
func Test_buildMetastatusCmdArgs(test *testing.T) {
	args, timeout := buildMetastatusCmdArgs(&PaastaMetastatusOptions{})
	assert.Equal(test, args, []string{})
	assert.Equal(test, timeout, time.Duration(20))

	args, timeout = buildMetastatusCmdArgs(&PaastaMetastatusOptions{
		PaastaOptions: cli.PaastaOptions{Verbosity: 5},
	})
	assert.Equal(test, args, []string{"-vvvvv"})
	assert.Equal(test, timeout, time.Duration(120))

	args, timeout = buildMetastatusCmdArgs(&PaastaMetastatusOptions{
		AutoscalingInfo: true,
	})
	assert.Equal(test, args, []string{"-a", "-vv"})
	assert.Equal(test, timeout, time.Duration(120))

	args, _ = buildMetastatusCmdArgs(&PaastaMetastatusOptions{
		Groupings: []string{"foo", "bar"},
	})
	assert.Equal(test, args, []string{"-g", "foo", "bar"})

	args, _ = buildMetastatusCmdArgs(&PaastaMetastatusOptions{
		PaastaOptions: cli.PaastaOptions{UseMesosCache: true},
	})
	assert.Equal(test, args, []string{"--use-mesos-cache"})
}

func Test_writeAPIStatus(test *testing.T) {
	// TODO: more tests
}

func Test_getClusterStatus(test *testing.T) {
	// TODO: more tests
}
