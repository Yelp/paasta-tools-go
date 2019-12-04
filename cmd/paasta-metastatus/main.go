package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/Yelp/paasta-tools-go/pkg/cli"
	"github.com/Yelp/paasta-tools-go/pkg/config_store"
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

func metastatus(opts *PaastaMetastatusOptions) error {
	if opts.AutoscalingInfo {
		if opts.Verbosity < 2 {
			opts.Verbosity = 2
		}
	}
	sysStore := config_store.NewStore(opts.SysDir, nil)
	apiEndpoints, err := sysStore.Get("api_endpoints")
	if err != nil {
		return err
	}

	var clusters []interface{}
	if opts.Cluster != "" {
		clusters = []interface{}{opts.Cluster}
	} else {
		clusters, err := sysStore.Get("clusters")
		if err != nil {
			return err
		}
	}

	return nil
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
	err = metastatus(options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
