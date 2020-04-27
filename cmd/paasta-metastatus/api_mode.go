package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"github.com/Yelp/paasta-tools-go/pkg/paasta_api/client/operations"
	"github.com/go-openapi/strfmt"
	"github.com/logrusorgru/aurora"

	apiclient "github.com/Yelp/paasta-tools-go/pkg/paasta_api/client"
	httptransport "github.com/go-openapi/runtime/client"
)

// APIMode ...
type APIMode struct{}

func (a *APIMode) printDashboards(
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

func (a *APIMode) getMetastatusCmdArgs(opts *PaastaMetastatusOptions) ([]string, time.Duration) {
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

func (a *APIMode) printAPIStatus(cluster, endpoint string, opts *PaastaMetastatusOptions, sb *strings.Builder) error {
	url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Failed to parse API endpoint %v: %v", endpoint, err)
	}

	var (
		transport        = httptransport.New(url.Host, apiclient.DefaultBasePath, []string{url.Scheme})
		client           = apiclient.New(transport, strfmt.Default)
		cmdArgs, timeout = a.getMetastatusCmdArgs(opts)
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	mp := &operations.MetastatusParams{CmdArgs: cmdArgs, Context: ctx}
	resp, err := client.Operations.Metastatus(mp)
	if err != nil {
		return fmt.Errorf("Failed to get metastatus: %s", aurora.Red(err))
	}
	sb.WriteString(fmt.Sprintf("%s\n", resp.Payload.Output))
	return nil
}

func (a *APIMode) printClusterStatus(
	cluster, endpoint string,
	dashboards map[string]interface{},
	opts *PaastaMetastatusOptions,
) bool {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("Cluster: %v\n", cluster))
	a.printDashboards(cluster, dashboards, sb)
	success := true
	err := a.printAPIStatus(cluster, endpoint, opts, sb)
	if err != nil {
		success = false
		sb.WriteString(fmt.Sprintf("Failed to get status for cluster %v: %v\n", cluster, err))
	}
	fmt.Print(sb.String())
	return success
}

func (a *APIMode) metastatus(opts *PaastaMetastatusOptions) (bool, error) {
	if opts.AutoscalingInfo {
		if opts.Verbosity < 2 {
			opts.Verbosity = 2
		}
	}
	sysStore := configstore.NewStore(opts.SysDir, nil)

	apiEndpoints := map[string]string{}
	ok, err := sysStore.Load("api_endpoints", &apiEndpoints)
	if err != nil {
		return false, fmt.Errorf("Failed to load api_endpoints from configs: %v", err)
	}
	if !ok {
		return false, fmt.Errorf("No api_endpoints configured in %v", sysStore.Dir)
	}

	dashboardLinks := map[string]map[string]interface{}{}
	_, err = sysStore.Load("dashboard_links", &dashboardLinks)
	if err != nil {
		return false, fmt.Errorf("Failed to load dashboard_links from configs: %v", err)
	}

	var clusters []string
	if opts.Cluster != "" {
		clusters = []string{opts.Cluster}
	} else {
		ok, err := sysStore.Load("clusters", &clusters)
		if err != nil {
			return false, fmt.Errorf("Failed to load clusters from configs: %v", err)
		}
		if !ok {
			return false, fmt.Errorf("No clusters configured")
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
			if !a.printClusterStatus(cluster, endpoint, dashboards, opts) {
				success = false
			}
		}(cluster, endpoint)
	}
	wg.Wait()

	return success, nil
}
