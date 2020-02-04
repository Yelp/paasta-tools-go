package monitoringspec

import (
	"encoding/json"
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
	Components []string      `json:"components,omitempty"`
	Description string       `json:"description,omitempty"`
	AlertAfter string        `json:"alert_after,omitempty"`
	RealertEvery string      `json:"realert_every,omitempty"`
	CheckEvery string        `json:"check_every,omitempty"`
	CheckOomEvents *bool     `json:"check_oom_events,omitempty"`
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
	if in.Components != nil {
		in, out := &in.Components, &out.Components
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

func (in *PaastaMonitoringSpec) DeepCopy() *PaastaMonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(PaastaMonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// GetSensuDefinition : Format the PaastaMonitoringSpec to be understood by yelp-sensu-runner
func (spec *PaastaMonitoringSpec) GetSensuDefinition() []byte {
	bytes, _ := json.Marshal(spec)
	return bytes
}
