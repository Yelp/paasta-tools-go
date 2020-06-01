package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const DEFAULT_HACHECK_PORT = 6666
const DEFAULT_HACHECK_HOST = "169.254.255.254"

type HadownOptions struct {
	ServiceName string
	ServicePort int
	ServiceIP   string
	Reason      string
	Expiration  float64
	HacheckHost string
	HacheckPort int
}

func hadownRequest(opts *HadownOptions) (*http.Request, error) {
	req_url := fmt.Sprintf("http://%s:%d/spool/%s/%d/", opts.HacheckHost, opts.HacheckPort, opts.ServiceName, opts.ServicePort)

	data := url.Values{}
	data.Set("status", "down")
	data.Set("reason", opts.Reason)
	if opts.Expiration > 0 {
		data.Set("expiration", fmt.Sprintf("%.f", opts.Expiration))
	}
	req, err := http.NewRequest("POST", req_url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if opts.ServiceIP != "" {
		req.Header.Set("X-Nerve-Check-IP", opts.ServiceIP)
	}
	return req, nil
}

func parseFlags(opts *HadownOptions) error {
	flag.StringVar(&opts.ServiceName, "service", "", "Service to down")
	flag.StringVar(&opts.Reason, "reason", "", "Reason for downing service")
	flag.IntVar(&opts.ServicePort, "servicePort", 0, "Port to set status for")
	flag.StringVar(&opts.ServiceIP, "serviceIP", "", "IP to set status for")
	flag.StringVar(&opts.HacheckHost, "host", DEFAULT_HACHECK_HOST, "Host that hacheck is running on")
	flag.IntVar(&opts.HacheckPort, "port", DEFAULT_HACHECK_PORT, "Port that hacheck is running on")
	flag.Float64Var(&opts.Expiration, "expiration", 0, "Expiration of down status (unix time)")
	flag.Parse()

	if opts.ServiceName == "" {
		return fmt.Errorf("Service name required")
	}
	if opts.Reason == "" {
		return fmt.Errorf("Reason is required")
	}
	return nil
}

func main() {
	options := &HadownOptions{}
	err := parseFlags(options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req, err := hadownRequest(options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("Marked %s on %s:%d as DOWN for reason \"%s\"\n", options.ServiceName, options.ServiceIP, options.ServicePort, options.Reason)
	}
	defer resp.Body.Close()
}
