// +build !yelp

package main

import (
	"fmt"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const zipkinReporter = "http"

func initZipkin(zipkinURL string) (reporter.Reporter, *zipkin.Tracer, error) {
	reporter := reporterhttp.NewReporter(zipkinURL)

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

	return reporter, tracer, err
}
