package monitoringspec

import (
	"encoding/json"
	"fmt"
	"github.com/Yelp/paasta-tools-go/pkg/config"
	"path"
)

const (
	SoaDir string = "/nail/etc/services"
	ServiceMonitoringFile string = "monitoring.yaml"
)

// PaastaMonitoringSpec : Spec for PaaSTA service monitoring configuration
type PaastaMonitoringSpec struct {
	Team string              `json:"team,omitempty"`
	Runbook string           `json:"runbook,omitempty"`
	Tip string               `json:"tip,omitempty"`
	Page *bool               `json:"page,omitempty"`
	Ticket *bool             `json:"ticket,omitempty"`
	NotificationEmail string `json:"notification_email,omitempty"`
	SlackChannels []string   `json:"slack_channels,omitempty"`
	Project string           `json:"project,omitempty"`
	Priority string          `json:"priority,omitempty"`
	Tags []string            `json:"tags,omitempty"`
	Componenents []string    `json:"components,omitempty"`
	Description string       `json:"description,omitempty"`
	AlertAfter string        `json:"alert_after,omitempty"`
	RealertEvery string      `json:"realert_every,omitempty"`
	CheckEvery string        `json:"check_every,omitempty"`
	CheckOomEvents *bool     `json:"check_oom_events,omitempty"`
}

type Monitoring interface {
	GetMonitoringSpec() (*PaastaMonitoringSpec, error)
}

type ServiceMonitoring struct {
	service string
}

type MultiMonitoring struct {
	monitors []Monitoring
}

func (in *PaastaMonitoringSpec) DeepCopyInto(out *PaastaMonitoringSpec) {
	if in.Team != "" {
		out.Team = in.Team
	}
	if in.Runbook != "" {
		out.Runbook = in.Runbook
	}
	if in.Tip != "" {
		out.Tip = in.Tip
	}
	if in.Page != nil {
		in, out := &in.Page, &out.Page
		*out = new(bool)
		**out = **in
	}
	if in.Ticket != nil {
		in, out := &in.Ticket, &out.Ticket
		*out = new(bool)
		**out = **in
	}
	if in.NotificationEmail != "" {
		out.NotificationEmail = in.NotificationEmail
	}
	if in.SlackChannels != nil {
		in, out := &in.SlackChannels, &out.SlackChannels
		*out = make([]string, len(*in))
		for i, elem := range *in {
			(*out)[i] = elem
		}
	}
	if in.Project != "" {
		out.Project = in.Project
	}
	if in.Priority != "" {
		out.Priority = in.Priority
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]string, len(*in))
		for i, elem := range *in {
			(*out)[i] = elem
		}
	}
	if in.Componenents != nil {
		in, out := &in.Componenents, &out.Componenents
		*out = make([]string, len(*in))
		for i, elem := range *in {
			(*out)[i] = elem
		}
	}
	if in.Description != "" {
		out.Description = in.Description
	}
	if in.AlertAfter != "" {
		out.AlertAfter = in.AlertAfter
	}
	if in.RealertEvery != "" {
		out.CheckEvery = in.CheckEvery
	}
	if in.CheckEvery != "" {
		out.CheckEvery = in.CheckEvery
	}
	if in.CheckOomEvents != nil {
		in, out := &in.CheckOomEvents, &out.CheckOomEvents
		*out = new(bool)
		**out = **in
	}
}

func (spec *PaastaMonitoringSpec) GetMonitoringSpec() (*PaastaMonitoringSpec, error) {
	return spec, nil
}

// ToSensu : Format the PaastaMonitoringSpec to be understood by yelp-sensu-runner
func (s *PaastaMonitoringSpec) ToSensu() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}

func NewServiceMonitoring(service string) *ServiceMonitoring {
	return &ServiceMonitoring{
		service: service,
	}
}

func (s *ServiceMonitoring) GetMonitoringSpec() (*PaastaMonitoringSpec, error) {
	configReader := config.ConfigFileReader{
		Basedir:  path.Join(SoaDir, s.service),
		Filename: ServiceMonitoringFile,
	}

	monitoring, err := monitoringFromConfig(configReader)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading Monitoring information for service %s: %s", s.service, err,
		)
	}

	return monitoring, nil
}

func (m *MultiMonitoring) GetMonitoringSpec() (*PaastaMonitoringSpec, error) {
	var err error
	spec := &PaastaMonitoringSpec{}
	for _, monitor := range m.monitors {
		monitoring, err := monitor.GetMonitoringSpec()
		if err == nil {
			monitoring.DeepCopyInto(spec)
		}
	}
	// if the last PaastaMonitoringSpec is valid, err will be nil
	return spec, err
}

func monitoringFromConfig(cr config.ConfigReader) (*PaastaMonitoringSpec, error) {
	spec := &PaastaMonitoringSpec{}
	err := cr.Read(spec)
	return spec, err
}

// ServiceMonitoringSpec returns an updated version of a PaastaMonitoringSpec, adding defaults
// potentially found in the monitoring configuration of the given service.
//
// This reads /nail/etc/services/<service>/monitoring.yaml to fill the gaps in `spec`.
func ServiceMonitoringSpec(service string, spec *PaastaMonitoringSpec) (*PaastaMonitoringSpec, error) {
	m := MultiMonitoring{
		monitors: []Monitoring{
			NewServiceMonitoring(service),
			spec,
		},
	}
	return m.GetMonitoringSpec()
}
