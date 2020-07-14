// +build yelp

package main

import (
	"fmt"

	reportermonk "github.com/Yelp/paasta-tools-go/pkg/zipkin/reporter/monk"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
)

const zipkinReporter = "monk"

func noopZipkin(err error) (reporter.Reporter, *zipkin.Tracer, error) {
	rep := reporter.NewNoopReporter()
	tr, _ := zipkin.NewTracer(rep)
	return rep, tr, err
}

func initZipkin(zipkinURL string) (reporter.Reporter, *zipkin.Tracer, error) {
	if zipkinURL == "" {
		return noopZipkin(fmt.Errorf("zipkin url missing"))
	}

	reporter, err := reportermonk.NewReporter(zipkinURL)
	if err != nil {
		return noopZipkin(fmt.Errorf("initializing reporter: %v", err))
	}

	localEndpoint, err := zipkin.NewEndpoint("paasta-cli", "localhost:0")
	if err != nil {
		return noopZipkin(fmt.Errorf("initializing endpoint: %v", err))
	}

	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return noopZipkin(fmt.Errorf("initializing sampler: %v", err))
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return noopZipkin(fmt.Errorf("initializing tracer: %v", err))
	}

	return reporter, tracer, nil
}
