package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/Yelp/paasta-tools-go/pkg/cli"
	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"github.com/logrusorgru/aurora"

	apiclient "github.com/Yelp/paasta-tools-go/pkg/paasta_api/client"
	"github.com/Yelp/paasta-tools-go/pkg/paasta_api/client/operations"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
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

func printDashboards(
	cluster string, dashboards map[string]interface{}, sb *strings.Builder,
) {
	if dashboards == nil {
		sb.WriteString(aurora.Red("No dashboards configured!\n").String())
	} else {
		spacing := 0
		for svc := range dashboards {
			svcLen := len(svc)
			if svcLen > spacing {
				spacing = svcLen
			}
		}
		spacing++
		for svc, dashboard := range dashboards {
			spacer := strings.Repeat(" ", spacing-len(svc))
			sb.WriteString(fmt.Sprintf("  %s:%s", svc, spacer))
			switch d := dashboard.(type) {
			case string:
				sb.WriteString(aurora.Cyan(d).String())
			case []interface{}:
				if len(d) > 1 {
					for _, url := range d {
						sb.WriteString(fmt.Sprintf("\n    %v", aurora.Cyan(url.(string))))
					}
				} else {
					sb.WriteString(aurora.Cyan(d[0].(string)).String())
				}
			}
			sb.WriteString("\n")
		}
	}
}

func printAPIStatus(cluster, endpoint string, sb *strings.Builder) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Failed to parse API endpoint %v: %v", endpoint, err)
	}
	transport := httptransport.New(url.Host, apiclient.DefaultBasePath, []string{url.Scheme})
	client := apiclient.New(transport, strfmt.Default)
	cmdArgs := []string{""}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mp := &operations.MetastatusParams{CmdArgs: cmdArgs, Context: ctx}
	resp, err := client.Operations.Metastatus(mp)
	if err != nil {
		return fmt.Errorf("Failed to get metastatus: %s", aurora.Red(err))
	}
	sb.WriteString(fmt.Sprintf("%s\n", resp.Payload.Output))
	return nil
}

func printClusterStatus(
	cluster, endpoint string,
	dashboards map[string]interface{},
) bool {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("Cluster: %v\n", cluster))
	printDashboards(cluster, dashboards, sb)
	success := true
	err := printAPIStatus(cluster, endpoint, sb)
	if err != nil {
		success = false
		sb.WriteString(fmt.Sprintf("Failed to get status for cluster %v: %v\n", cluster, err))
	}
	fmt.Print(sb.String())
	return success
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

	dashboardLinks := map[string]map[string]interface{}{}
	err = sysStore.Load("dashboard_links", &dashboardLinks)
	if err != nil {
		return false, fmt.Errorf("Failed to load dashboard_links from configs: %v", err)
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
		dashboards, ok := dashboardLinks[cluster]
		wg.Add(1)
		go func(cluster, endpoint string) {
			defer wg.Done()
			if !printClusterStatus(cluster, endpoint, dashboards) {
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
