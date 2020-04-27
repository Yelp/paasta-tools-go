package mesostools

import (
	"fmt"
	"log"
	"time"

	"github.com/mesos/mesos-go/api/v0/detector"
	// ...
	_ "github.com/mesos/mesos-go/api/v0/detector/zoo"
	mesos "github.com/mesos/mesos-go/api/v0/mesosproto"
)

// MesosConfig ...
type MesosConfig struct {
	Path            string        `mapstructure:"path,omitempty"`
	Debug           bool          `mapstructure:"debug,omitempty"`
	LogFile         string        `mapstructure:"log_file,omitempty"`
	LogLevel        string        `mapstructure:"log_level,omitempty"`
	Master          string        `mapstructure:"master,omitempty"`
	MaxWorkers      int           `mapstructure:"max_workers,omitempty"`
	Scheme          string        `mapstructure:"scheme,omitempty"`
	ResponseTimeout time.Duration `mapstructure:"response_timeout,omitempty"`
}

// DefaultMesosConfig ...
var DefaultMesosConfig = MesosConfig{
	Path:            "./tmp/mesos-cli.json",
	Debug:           false,
	LogFile:         "",
	LogLevel:        "warning",
	Master:          "localhost:5050",
	MaxWorkers:      5,
	Scheme:          "http",
	ResponseTimeout: 5,
}

// MasterDetector ...
type MasterDetector struct {
	MasterInfo *mesos.MasterInfo
	Updated    chan bool
}

// OnMasterChanged ...
func (l *MasterDetector) OnMasterChanged(info *mesos.MasterInfo) {
	log.Printf("New master info: %+v\n", info)
	l.MasterInfo = info
	l.Updated <- true
}

// Wait ...
func (l *MasterDetector) Wait() {
	<-l.Updated
}

// NewMasterDetector will fetch mesos master info from zookeeper
func NewMasterDetector(mesosConfig *MesosConfig) (*MasterDetector, error) {
	log.Printf("creating ZK detector for %q\n", mesosConfig.Master)

	m, err := detector.New(mesosConfig.Master)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZK listener for Mesos masters: %v", err)
	}

	log.Println("created ZK detector")
	d := &MasterDetector{
		Updated: make(chan bool),
	}
	err = m.Detect(d)
	if err != nil {
		return nil, fmt.Errorf("failed to register ZK listener: %v", err)
	}

	return d, nil
}
