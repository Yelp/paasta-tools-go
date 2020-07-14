package zipkin

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
)

type zipkinInitializer interface {
	zipkinInitialize(string) (reporter.Reporter, *zipkin.Tracer, error)
}

var zipkinInitializers map[string]zipkinInitializer

type noopInitializer struct{}

func (*noopInitializer) zipkinInitialize(_ string) (reporter.Reporter, *zipkin.Tracer, error) {
	rep := reporter.NewNoopReporter()
	tr, _ := zipkin.NewTracer(rep)
	return rep, tr, nil
}

func init() {
	if zipkinInitializers == nil {
		zipkinInitializers = map[string]zipkinInitializer{}
	}
	zipkinInitializers["noop"] = &noopInitializer{}
}

// InitZipkin returns the reporter and tracer for zipkinURL
func InitZipkin(zipkinURL string) (reporter.Reporter, *zipkin.Tracer, error) {
	var initializer string
	var errors []string

	if zipkinURL == "" {
		errors = append(errors, fmt.Sprintf("zipkin URL empty"))
		initializer = "noop"
	} else {
		url, err := url.Parse(zipkinURL)
		if err != nil {
			errors = append(errors, fmt.Sprintf("parsing zipkin url: %v", err))
			initializer = "noop"
		}
		if url.Scheme != "" {
			initializer = url.Scheme
		}
	}
	initializerF, ok := zipkinInitializers[initializer]
	if !ok {
		errors = append(errors, fmt.Sprintf("zipkin initializer for %s not found", initializer))
		initializerF, _ = zipkinInitializers["noop"]
	}
	rep, tr, err := initializerF.zipkinInitialize(zipkinURL)
	if err != nil {
		errors = append(errors, fmt.Sprintf("initializing %T: %v", initializerF, err))
	}
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, ", "))
	}
	return rep, tr, err
}

// Initializers returns a list of registered zipkin initializers
func Initializers() []string {
	res := []string{}
	for k := range zipkinInitializers {
		res = append(res, k)
	}
	return res
}
