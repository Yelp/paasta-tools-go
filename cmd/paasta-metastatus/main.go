package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/Yelp/paasta-tools-go/pkg/cli"
	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"github.com/logrusorgru/aurora"

	"github.com/Yelp/paasta-tools-go/pkg/paastaapi"
)

// PaastaMetastatusOptions ...
type PaastaMetastatusOptions struct {
	cli.PaastaOptions
	cli.CSIOptions
	AutoscalingInfo bool
	Groupings       []string
}

// Setup ...
func (o *PaastaMetastatusOptions) Setup() {
	(&o.PaastaOptions).Setup()
	(&o.CSIOptions).Setup()
	flag.BoolVar(&o.AutoscalingInfo, "autoscaling-info", false, "")
	flag.StringArrayVarP(&o.Groupings, "groupings", "g", []string{"region"}, "")
}

func parseFlags(opts *PaastaMetastatusOptions) error {
	opts.Setup()
	flag.CommandLine.SortFlags = false
	flag.Parse()
	return nil
}

func writeDashboards(
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
			case []string:
				if len(d) > 1 {
					for _, url := range d {
						sb.WriteString(
							fmt.Sprintf("\n    %v", aurora.Cyan(url)),
						)
					}
				} else if len(d) == 1 {
					sb.WriteString(aurora.Cyan(d[0]).String())
				}
			}
			sb.WriteString("\n")
		}
	}
}

func buildMetastatusCmdArgs(opts *PaastaMetastatusOptions) ([]string, time.Duration) {
	cmdArgs := []string{}
	verbosity := 0
	timeout := time.Duration(20)

	if opts.Verbosity > 0 {
		verbosity = opts.Verbosity
		timeout = 120
	}
	if opts.AutoscalingInfo {
		cmdArgs = append(cmdArgs, "-a")
		if verbosity < 2 {
			verbosity = 2
		}
		timeout = 120
	}
	if len(opts.Groupings) > 0 {
		cmdArgs = append(cmdArgs, "-g")
		cmdArgs = append(cmdArgs, opts.Groupings...)
	}
	if opts.UseMesosCache {
		cmdArgs = append(cmdArgs, "--use-mesos-cache")
	}
	if verbosity > 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("-%s", strings.Repeat("v", verbosity)))
	}

	return cmdArgs, timeout
}

type ctxKey int

const (
	ctxKeyHTTPClient ctxKey = iota
	ctxKeyOut
	ctxKeyErr
)

func writeAPIStatus(
	ctx context.Context,
	cluster, endpoint string,
	cmdArgs []string,
	sb *strings.Builder,
) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Failed to parse API endpoint %v: %v", endpoint, err)
	}

	config := paastaapi.NewConfiguration()
	config.Host = url.Host
	config.Scheme = url.Scheme

	httpClient := ctx.Value(ctxKeyHTTPClient)
	if httpClient != nil {
		config.HTTPClient = httpClient.(*http.Client)
	}

	client := paastaapi.NewAPIClient(config)
	mr := client.DefaultApi.Metastatus(ctx).CmdArgs(cmdArgs)
	metastatusResp, _, err := mr.Execute()
	if err != nil {
		return fmt.Errorf("Failed to get metastatus: %s", aurora.Red(err))
	}
	sb.WriteString(fmt.Sprintf("%s\n", metastatusResp.GetOutput()))
	return nil
}

func getClusterStatus(
	ctx context.Context,
	cluster, endpoint string,
	dashboards map[string]interface{},
	cmdArgs []string,
) (*strings.Builder, error) {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("Cluster: %v\n", cluster))
	writeDashboards(cluster, dashboards, sb)
	err := writeAPIStatus(ctx, cluster, endpoint, cmdArgs, sb)
	if err != nil {
		return sb, fmt.Errorf("Failed to get status for cluster %v: %v", cluster, err)
	}
	return sb, nil
}

func metastatus(
	ctx context.Context,
	clusters []string,
	apiEndpoints map[string]string,
	dashboardLinks map[string]map[string]interface{},
	cmdArgs []string,
) (bool, error) {
	outf := ctx.Value(ctxKeyOut).(io.Writer)
	errf := ctx.Value(ctxKeyErr).(io.Writer)

	var wg sync.WaitGroup
	var success bool = true
	for _, cluster := range clusters {
		endpoint, ok := apiEndpoints[cluster]
		if !ok {
			fmt.Fprintf(errf, "WARN: api endpoint not found for %v\n", cluster)
			continue
		}
		dashboards, _ := dashboardLinks[cluster]
		wg.Add(1)
		go func(cluster, endpoint string) {
			defer wg.Done()
			sb, err := getClusterStatus(ctx, cluster, endpoint, dashboards, cmdArgs)
			fmt.Fprint(outf, sb)
			if err != nil {
				fmt.Fprint(errf, err.Error())
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

	if options.AutoscalingInfo {
		if options.Verbosity < 2 {
			options.Verbosity = 2
		}
	}
	sysStore := configstore.NewStore(options.SysDir, nil)

	apiEndpoints := map[string]string{}
	ok, err := sysStore.Load("api_endpoints", &apiEndpoints)
	if !ok || err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load api_endpoints from configs: found=%v, error=%v", ok, err)
		os.Exit(1)
	}

	dashboardLinks := map[string]map[string]interface{}{}
	ok, err = sysStore.Load("dashboard_links", &dashboardLinks)
	if !ok || err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load dashboard_links from configs: found=%v, error=%v", ok, err)
		os.Exit(1)
	}

	var clusters []string
	if options.Cluster != "" {
		clusters = []string{options.Cluster}
	} else {
		ok, err := sysStore.Load("clusters", &clusters)
		if !ok || err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load clusters from configs: found=%v, error=%v", ok, err)
			os.Exit(1)
		}
	}

	cmdArgs, timeout := buildMetastatusCmdArgs(options)
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, ctxKeyOut, os.Stdout)
	ctx = context.WithValue(ctx, ctxKeyErr, os.Stderr)

	success, err := metastatus(ctx, clusters, apiEndpoints, dashboardLinks, cmdArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !success {
		os.Exit(1)
	}
}
