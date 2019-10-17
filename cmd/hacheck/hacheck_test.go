package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestHadownRequest(t *testing.T) {
	opts := &HadownOptions{
		ServiceName: "foo",
		ServicePort: 12345,
		ServiceIP:   "10.10.1.2",
		Reason:      "testing",
		HacheckHost: "localhost",
		HacheckPort: 3333,
		Expiration:  1565214667,
	}
	request, err := hadownRequest(opts)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	if request.Method != "POST" {
		t.Errorf("Incorrect method, got: %s, want: POST", request.Method)
	}

	expectedURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:3333",
		Path:   "/spool/foo/12345/",
	}
	if !reflect.DeepEqual(request.URL, expectedURL) {
		t.Errorf("Incorrect URL, got: %v, want: %v", request.URL, expectedURL)
	}

	request.ParseForm()
	expectedForm := url.Values{
		"status":     {"down"},
		"reason":     {"testing"},
		"expiration": {"1565214667"},
	}
	if !reflect.DeepEqual(request.PostForm, expectedForm) {
		t.Errorf("Incorrect request body, got: %v, want: %v", request.PostForm, expectedForm)
	}

	ip_header := request.Header.Get("X-Nerve-Check-IP")
	if ip_header != opts.ServiceIP {
		t.Errorf("Incorrect IP header value, got: %v, want: %v", ip_header, opts.ServiceIP)
	}
}
