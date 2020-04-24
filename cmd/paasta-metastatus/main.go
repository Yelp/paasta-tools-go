package main

import (
	"fmt"
	"os"
	"sync"

	flag "github.com/spf13/pflag"

	"github.com/Yelp/paasta-tools-go/pkg/cli"
	"github.com/Yelp/paasta-tools-go/pkg/configstore"
)

// PaastaMetastatusOptions ...
type PaastaMetastatusOptions struct {
	cli.PaastaOptions
	cli.CSIOptions
	AutoscalingInfo bool
	Groupings       string
}

// Setup ...
func (o *PaastaMetastatusOptions) Setup() {
	(&o.PaastaOptions).Setup()
	(&o.CSIOptions).Setup()
	flag.BoolVar(&o.AutoscalingInfo, "autoscaling-info", false, "")
	flag.StringVarP(&o.Groupings, "groupings", "g", "", "")
}

func parseFlags(opts *PaastaMetastatusOptions) error {
	opts.Setup()
	flag.CommandLine.SortFlags = false
	flag.Parse()
	return nil
}

func printClusterStatus(
	cluster, endpoint string,
	opts *PaastaMetastatusOptions,
	conf *configstore.Store,
) error {
	fmt.Printf("cluster %v endpoint %v\n", cluster, endpoint)
	return nil
}

func metastatus(opts *PaastaMetastatusOptions) (bool, error) {
	if opts.AutoscalingInfo {
		if opts.Verbosity < 2 {
			opts.Verbosity = 2
		}
	}
	sysStore := configstore.NewStore(opts.SysDir, nil)

	apiEndpoints := map[string]string{}
	err := sysStore.Load("api_endpoints", &apiEndpoints)
	if err != nil {
		return false, fmt.Errorf("Failed to load api_endpoints from configs: %v", err)
	}

	var clusters []string
	if opts.Cluster != "" {
		clusters = []string{opts.Cluster}
	} else {
		err := sysStore.Load("clusters", &clusters)
		if err != nil {
			return false, fmt.Errorf("Failed to load clusters from configs: %v", err)
		}
	}

	var wg sync.WaitGroup
	var success bool = true
	for _, cluster := range clusters {
		endpoint, ok := apiEndpoints[cluster]
		if !ok {
			fmt.Printf("WARN: api endpoint not found for %v\n", cluster)
			continue
		}
		wg.Add(1)
		go func(cluster, endpoint string) {
			defer wg.Done()
			err := printClusterStatus(cluster, endpoint, opts, sysStore)
			if err != nil {
				fmt.Printf(
					"ERR: couldn't get status for cluster %v: %v", cluster, err,
				)
				success = false
			}
		}(cluster, endpoint)
	}
	wg.Wait()

	return success, nil
}

func main() {
	options := &PaastaMetastatusOptions{}
	err := parseFlags(options)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	if options.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	success, err := metastatus(options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !success {
		os.Exit(1)
	}
}
