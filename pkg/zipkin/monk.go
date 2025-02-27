//go:build yelp
// +build yelp

package zipkin

import (
	"fmt"

	reportermonk "github.com/Yelp/paasta-tools-go/pkg/zipkin/reporter/monk"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
)

type monkInitializer struct{}

func (*monkInitializer) zipkinInitialize(zipkinURL string) (reporter.Reporter, *zipkin.Tracer, error) {
	if zipkinURL == "" {
		return nil, nil, fmt.Errorf("zipkin url missing")
	}

	reporter, err := reportermonk.NewReporter(zipkinURL)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing reporter: %v", err)
	}

	localEndpoint, err := zipkin.NewEndpoint("paasta-cli", "localhost:0")
	if err != nil {
		return nil, nil, fmt.Errorf("initializing endpoint: %v", err)
	}

	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing sampler: %v", err)
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing tracer: %v", err)
	}

	return reporter, tracer, nil
}

func init() {
	if zipkinInitializers == nil {
		zipkinInitializers = map[string]zipkinInitializer{}
	}
	zipkinInitializers["monk"] = &monkInitializer{}
}
