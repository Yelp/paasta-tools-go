package main

import (
	"fmt"
	"os"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	"github.com/Yelp/paasta-tools-go/pkg/mesostools"
	"github.com/mitchellh/mapstructure"
)

// DirectMode ...
type DirectMode struct{}

func (d *DirectMode) isKubernetesAvailable() (bool, error) {
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		return false, nil
	}
	info, err := os.Stat(kubeConfig)
	return err != nil && !info.IsDir(), err
}

// CheckResult ...
type CheckResult struct {
	Message string
}

func (d *DirectMode) runMesosChecks(opts *PaastaMetastatusOptions, config *configstore.Store) ([]CheckResult, error) {
	mesosConfig := mesostools.DefaultMesosConfig
	_, err := config.Load("mesos_config", mesosConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to load mesos_config: %v", err)
	}
	if mesosConfig.Path == "" {
		return nil, fmt.Errorf("Failed to find mesos_config in %v", config.Dir)
	}
	info, err := os.Stat(mesosConfig.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to stat mesos config file: %v", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("Mesos config path must be a file: %v", mesosConfig.Path)
	}

	var mesosConfigJSON map[string]interface{}
	err = config.ParseFile(mesosConfig.Path, &mesosConfigJSON)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %v", mesosConfig.Path, err)
	}

	profile := "default"
	profileFromConfig, ok := mesosConfigJSON["profile"]
	if ok {
		profile = profileFromConfig.(string)
	}
	mesosConfigProfile, ok := mesosConfigJSON[profile]
	if ok {
		err = mapstructure.Decode(mesosConfigProfile, &mesosConfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to destructure %s: %v", mesosConfigJSON, err)
		}
	}

	fmt.Printf("Mesos config: %+v\n", mesosConfig)
	masterDetector, err := mesostools.NewMasterDetector(&mesosConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed initialize master detector: %v")
	}
	masterDetector.Wait()
	fmt.Printf("Mesos master: %+v\n", masterDetector.MasterInfo)

	return []CheckResult{}, nil
	// if mesos_available:
	//     master_kwargs = {}
	//     # we don't want to be passing False to not override a possible True
	//     # value from system config
	//     if args.use_mesos_cache:
	//         master_kwargs["use_mesos_cache"] = True

	//     master = get_mesos_master(**master_kwargs)

	//     marathon_servers = get_marathon_servers(system_paasta_config)
	//     marathon_clients = all_marathon_clients(get_marathon_clients(marathon_servers))

	//     try:
	//         mesos_state = a_sync.block(master.state)
	//         all_mesos_results = _run_mesos_checks(
	//             mesos_master=master, mesos_state=mesos_state
	//         )
	//     except MasterNotAvailableException as e:
	//         # if we can't connect to master at all,
	//         # then bomb out early
	//         paasta_print(PaastaColors.red("CRITICAL:  %s" % "\n".join(e.args)))
	//         raise FatalError(2)

	//     marathon_results = _run_marathon_checks(marathon_clients)
	// else:
	//     marathon_results = [
	//         metastatus_lib.HealthCheckResult(
	//             message="Marathon is not configured to run here", healthy=True
	//         )
	//     ]
	//     all_mesos_results = [
	//         metastatus_lib.HealthCheckResult(
	//             message="Mesos is not configured to run here", healthy=True
	//         )
	//     ]
}

func (d *DirectMode) runKubernetesChecks(opts *PaastaMetastatusOptions, config *configstore.Store) []CheckResult {
	kubernetesAvailable, err := d.isKubernetesAvailable()
	if err != nil {
		fmt.Printf("Error while loading kubernetes config: %v\n", err)
	}
	fmt.Printf("Kubernetes available: %v\n", kubernetesAvailable)
	if kubernetesAvailable {
		return []CheckResult{}
	}
	// if kube_available:
	//     kube_client = KubeClient()
	//     kube_results = _run_kube_checks(kube_client)
	// else:
	//     kube_results = [
	//         metastatus_lib.HealthCheckResult(
	//             message="Kubernetes is not configured to run here", healthy=True
	//         )
	//     ]
	return nil
}

func (d *DirectMode) printClusterStatus(
	cluster string,
	opts *PaastaMetastatusOptions,
	config *configstore.Store,
) (bool, error) {
	mesosResults, mesosErr := d.runMesosChecks(opts, config)
	kubernetesResults := d.runKubernetesChecks(opts, config)

	fmt.Printf("mesos results: %v %v\nkubernetes results: %v\n", mesosResults, mesosErr, kubernetesResults)

	// mesos_ok = all(metastatus_lib.status_for_results(all_mesos_results))
	// marathon_ok = all(metastatus_lib.status_for_results(marathon_results))
	// kube_ok = all(metastatus_lib.status_for_results(kube_results))

	// mesos_summary = metastatus_lib.generate_summary_for_check("Mesos", mesos_ok)
	// marathon_summary = metastatus_lib.generate_summary_for_check(
	//     "Marathon", marathon_ok
	// )
	// kube_summary = metastatus_lib.generate_summary_for_check("Kubernetes", kube_ok)

	// healthy_exit = True if all([mesos_ok, marathon_ok]) else False

	// paasta_print(f"Master paasta_tools version: {__version__}")
	// paasta_print("Mesos leader: %s" % get_mesos_leader())
	// metastatus_lib.print_results_for_healthchecks(
	//     mesos_summary, mesos_ok, all_mesos_results, args.verbose
	// )
	// if args.verbose > 1 and mesos_available:
	//     print_with_indent("Resources Grouped by %s" % ", ".join(args.groupings), 2)
	//     all_rows, healthy_exit = utilization_table_by_grouping_from_mesos_state(
	//         groupings=args.groupings, threshold=args.threshold, mesos_state=mesos_state
	//     )
	//     for line in format_table(all_rows):
	//         print_with_indent(line, 4)

	//     if args.autoscaling_info:
	//         print_with_indent("Autoscaling resources:", 2)
	//         headers = [
	//             field.replace("_", " ").capitalize()
	//             for field in AutoscalingInfo._fields
	//         ]
	//         table = [headers] + [
	//             [str(x) for x in asi]
	//             for asi in get_autoscaling_info_for_all_resources(mesos_state)
	//         ]

	//         for line in format_table(table):
	//             print_with_indent(line, 4)

	//     if args.verbose >= 3:
	//         print_with_indent("Per Slave Utilization", 2)
	//         cluster = system_paasta_config.get_cluster()
	//         service_instance_stats = get_service_instance_stats(
	//             args.service, args.instance, cluster
	//         )
	//         if service_instance_stats:
	//             print_with_indent(
	//                 "Service-Instance stats:" + str(service_instance_stats), 2
	//             )
	//         # print info about slaves here. Note that we don't make modifications to
	//         # the healthy_exit variable here, because we don't care about a single slave
	//         # having high usage.
	//         all_rows, _ = utilization_table_by_grouping_from_mesos_state(
	//             groupings=args.groupings + ["hostname"],
	//             threshold=args.threshold,
	//             mesos_state=mesos_state,
	//             service_instance_stats=service_instance_stats,
	//         )
	//         # The last column from utilization_table_by_grouping_from_mesos_state is "Agent count", which will always be
	//         # 1 for per-slave resources, so delete it.
	//         for row in all_rows:
	//             row.pop()

	//         for line in format_table(all_rows):
	//             print_with_indent(line, 4)
	// metastatus_lib.print_results_for_healthchecks(
	//     marathon_summary, marathon_ok, marathon_results, args.verbose
	// )
	// metastatus_lib.print_results_for_healthchecks(
	//     kube_summary, kube_ok, kube_results, args.verbose
	// )
	// if args.verbose > 1 and kube_available:
	//     print_with_indent("Resources Grouped by %s" % ", ".join(args.groupings), 2)
	//     all_rows, healthy_exit = utilization_table_by_grouping_from_kube(
	//         groupings=args.groupings, threshold=args.threshold, kube_client=kube_client
	//     )
	//     for line in format_table(all_rows):
	//         print_with_indent(line, 4)

	//     if args.autoscaling_info:
	//         print_with_indent("No autoscaling resources for Kubernetes", 2)

	//     if args.verbose >= 3:
	//         print_with_indent("Per Node Utilization", 2)
	//         cluster = system_paasta_config.get_cluster()
	//         service_instance_stats = get_service_instance_stats(
	//             args.service, args.instance, cluster
	//         )
	//         if service_instance_stats:
	//             print_with_indent(
	//                 "Service-Instance stats:" + str(service_instance_stats), 2
	//             )
	//         # print info about nodes here. Note that we don't make
	//         # modifications to the healthy_exit variable here, because we don't
	//         # care about a single node having high usage.
	//         all_rows, _ = utilization_table_by_grouping_from_kube(
	//             groupings=args.groupings + ["hostname"],
	//             threshold=args.threshold,
	//             kube_client=kube_client,
	//             service_instance_stats=service_instance_stats,
	//         )
	//         # The last column from utilization_table_by_grouping_from_kube is "Agent count", which will always be
	//         # 1 for per-node resources, so delete it.
	//         for row in all_rows:
	//             row.pop()

	//         for line in format_table(all_rows):
	//             print_with_indent(line, 4)

	// if not healthy_exit:
	//     raise FatalError(2)

	return true, nil
}

func (d *DirectMode) metastatus(opts *PaastaMetastatusOptions) (bool, error) {
	if opts.AutoscalingInfo {
		if opts.Verbosity < 2 {
			opts.Verbosity = 2
		}
	}
	sysStore := configstore.NewStore(
		opts.SysDir,
		// prevent loading all configs
		map[string]string{"mesos_config": "mesos_config"},
	)

	var cluster string
	if opts.Cluster != "" {
		cluster = opts.Cluster
	} else {
		// TODO: load cluster from /nail/etc/paasta_cluster
		cluster = "norcal-stagef"
	}

	return d.printClusterStatus(cluster, opts, sysStore)
}
