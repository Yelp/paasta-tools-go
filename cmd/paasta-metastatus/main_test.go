package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Yelp/paasta-tools-go/pkg/cli"

	"github.com/stretchr/testify/assert"
)

type recordingTransport struct {
	req []*http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.req = append(t.req, req)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func makeTestContext() context.Context {
	t := &recordingTransport{}
	t.req = make([]*http.Request, 0)
	cl := &http.Client{Transport: t}
	fmt.Printf("%p\n", cl)
	ctx := context.WithValue(context.Background(), ctxKeyHTTPClient, cl)
	ctx = context.WithValue(ctx, ctxKeyOut, &strings.Builder{})
	ctx = context.WithValue(ctx, ctxKeyErr, &strings.Builder{})
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

	cl := ctx.Value(ctxKeyHTTPClient).(*http.Client)
	tr := cl.Transport.(*recordingTransport)
	if len(tr.req) != 1 {
		test.Logf("requests: %v", tr.req)
		test.Errorf("expected number of requests: 1, got: %v", len(tr.req))
	}

	request := tr.req[0]
	if request.URL.Path != "/v1/metastatus" {
		test.Errorf("unexpected path: %v", request.URL.Path)
	}

	query := request.URL.RawQuery
	if query != "cmd_args=foo%2Cbar" {
		test.Errorf("expected mock args: %v, actual: %v", mockCmdArgs, query)
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
