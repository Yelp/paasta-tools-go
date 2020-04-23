// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// KubernetesPod kubernetes pod
//
// swagger:model KubernetesPod
type KubernetesPod struct {

	// containers
	Containers []*KubernetesContainer `json:"containers"`

	// Time at which the pod was deployed
	DeployedTimestamp float32 `json:"deployed_timestamp,omitempty"`

	// name of the pod's host
	Host string `json:"host,omitempty"`

	// long message explaining the pod's state
	Message *string `json:"message,omitempty"`

	// name of the pod in Kubernetes
	Name string `json:"name,omitempty"`

	// The status of the pod
	Phase string `json:"phase,omitempty"`

	// Whether or not the pod is ready (i.e. all containers up)
	Ready bool `json:"ready,omitempty"`

	// short message explaining the pod's state
	Reason *string `json:"reason,omitempty"`
}

// Validate validates this kubernetes pod
func (m *KubernetesPod) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContainers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *KubernetesPod) validateContainers(formats strfmt.Registry) error {

	if swag.IsZero(m.Containers) { // not required
		return nil
	}

	for i := 0; i < len(m.Containers); i++ {
		if swag.IsZero(m.Containers[i]) { // not required
			continue
		}

		if m.Containers[i] != nil {
			if err := m.Containers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("containers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *KubernetesPod) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *KubernetesPod) UnmarshalBinary(b []byte) error {
	var res KubernetesPod
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}