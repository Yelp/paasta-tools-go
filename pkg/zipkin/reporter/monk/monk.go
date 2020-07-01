package monk

import (
	"fmt"
	"net/url"
	"strconv"
	"unicode/utf8"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"k8s.io/klog/klogr"

	monkClient "github.yelpcorp.com/go-packages/monk/client"
	monkProducer "github.yelpcorp.com/go-packages/monk/producer"
)

var logger = klogr.New().WithName("paasta-go")

type monkReporter struct {
	producer   monkProducer.Producer
	stream     chan []byte
	serializer reporter.SpanSerializer
}

func (r *monkReporter) Send(m model.SpanModel) {
	// must report array of spans
	bytes, err := r.serializer.Serialize([]*model.SpanModel{&m})
	if err != nil {
		logger.Error(err, "failed to send zipkin span")
	}
	r.stream <- bytes
}

func (r *monkReporter) Close() error {
	r.producer.Close()
	return nil
}

// NewReporter creates Monk reporter for Zipkin
func NewReporter(zipkinURL string) (reporter.Reporter, error) {
	url, err := url.Parse(zipkinURL)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "monk" {
		return nil, fmt.Errorf("scheme must be monk, was %v", url.Scheme)
	}
	streamName := url.Path
	if streamName == "" {
		return nil, fmt.Errorf("stream name is missing")
	}
	_, i := utf8.DecodeRuneInString(streamName)
	streamName = streamName[i:]
	reporter := &monkReporter{
		serializer: reporter.JSONSerializer{},
	}
	port, err := strconv.ParseInt(url.Port(), 0, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing port: %v", err)
	}
	factory := monkClient.NewMonkConnectionFactory(url.Hostname(), uint64(port), "paasta-go")
	reporter.producer = monkProducer.New(factory)
	reporter.stream = reporter.producer.LogChannel(streamName, 10, 1024)
	return reporter, nil
}
