// meshregistrator has multiple goroutines:
// - fetchPods will take a snapshot of pods running in local kubelet and
//   turn them into a map of registrations
// - fetchLocal will gather locally running backends according to
//   configuration in ...
// - fetchAWS will download existing registrations in AWS cloudmap
//   to find potential zombies for cleanup
// - writeAWS will merge those and execute relevant cloudmap
//   register/deregister calls
package main

import (
	"context"
	"os"
	"sync"

	origflag "flag"

	flag "github.com/spf13/pflag"

	"k8s.io/klog/v2"
)

type MeshregistratorOptions struct {
	SystemPaastaDir  string
	YelpSoaDir       string
	LocalServicesDir string
	TrackLocal       bool
	TrackKubelet     bool
}

// Setup ...
func (o *MeshregistratorOptions) Setup() {
	flag.StringVarP(&o.SystemPaastaDir, "systempaastadir", "", "/etc/paasta", "")
	flag.StringVarP(&o.YelpSoaDir, "yelpsoadir", "", "/nail/etc/services", "")
	flag.StringVarP(&o.LocalServicesDir, "localservicesdir", "", "/etc/nerve/puppet_services.d", "")
	flag.BoolVarP(&o.TrackLocal, "tracklocal", "", true, "")
	flag.BoolVarP(&o.TrackKubelet, "trackkubelet", "", true, "")
}

func parseFlags(opts *MeshregistratorOptions) error {
	opts.Setup()
	flag.Parse()
	return nil
}

// A subprocess keeps track of ysoa-configs and some local configuration coming from puppet and other sources to understand Yelp’s service topology
func main() {
	klogFlags := origflag.NewFlagSet("klog", origflag.ExitOnError)
	klog.InitFlags(klogFlags)
	debug, _ := os.LookupEnv("MESHREGISTRATOR_DEBUG")
	v := klogFlags.Lookup("v")
	if v != nil {
		if debug != "" {
			v.Value.Set("10")
		} else {
			v.Value.Set("0")
		}
	}

	var options MeshregistratorOptions
	parseFlags(&options)

	klog.Infof("starting meshregistrator: %+v", options)

	// sysStore := configstore.NewStore(options.SystemPaastaDir, nil)

	var wg sync.WaitGroup
	pods := make(chan []ServiceRegistration, 1)
	local := make(chan []ServiceRegistration, 1)
	aws := make(chan []ServiceRegistration, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if options.TrackKubelet {
		wg.Add(1)
		go func() { defer wg.Done(); fetchPods(ctx, pods); cancel() }()
	}

	if options.TrackLocal {
		wg.Add(1)
		go func() { defer wg.Done(); fetchLocal(ctx, local, options.LocalServicesDir); cancel() }()
	}

	wg.Add(1)
	go func() { defer wg.Done(); fetchAWS(ctx, aws); cancel() }()

	wg.Add(1)
	go func() { defer wg.Done(); writeAWS(ctx, pods, local, aws); cancel() }()

	go signalLoop(ctx, cancel)
	wg.Wait()

	klog.Info("meshregistrator out")
}
