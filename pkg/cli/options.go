package cli

import flag "github.com/spf13/pflag"

type PaastaOptions struct {
	Help bool
	Verbosity int
	SoaDir string
	SysDir string
	UseMesosCache bool
}

func (o *PaastaOptions) Setup() {
	flag.BoolVarP(&o.Help, "help", "h", false, "")
	flag.CountVarP(&o.Verbosity, "verbose", "v", "")
	flag.StringVarP(&o.SoaDir, "soa-dir", "d", "/nail/etc/services", "")
	flag.StringVarP(&o.SysDir, "sys-dir", "C", "/etc/paasta", "")
	flag.BoolVarP(&o.UseMesosCache, "use-mesos-cache", "", false, "")
}

type CSIOptions struct {
	Cluster string
	Service string
	Instance string
}

func (o *CSIOptions) Setup() {
	flag.StringVarP(&o.Cluster, "cluster", "c", "", "")
	flag.StringVarP(&o.Service, "service", "s", "", "")
	flag.StringVarP(&o.Instance, "instance", "i", "", "")
}
