// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// InstanceStatus instance status
//
// swagger:model InstanceStatus
type InstanceStatus struct {

	// Adhoc instance status
	Adhoc InstanceStatusAdhoc `json:"adhoc,omitempty"`

	// flink
	Flink *InstanceStatusFlink `json:"flink,omitempty"`

	// Git sha of a service
	GitSha string `json:"git_sha,omitempty"`

	// Instance name
	Instance string `json:"instance,omitempty"`

	// Kubernetes instance status
	Kubernetes *InstanceStatusKubernetes `json:"kubernetes,omitempty"`

	// Marathon instance status
	Marathon *InstanceStatusMarathon `json:"marathon,omitempty"`

	// Service name
	Service string `json:"service,omitempty"`

	// Tron instance status
	Tron *InstanceStatusTron `json:"tron,omitempty"`
}

// Validate validates this instance status
func (m *InstanceStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdhoc(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFlink(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKubernetes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMarathon(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTron(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InstanceStatus) validateAdhoc(formats strfmt.Registry) error {

	if swag.IsZero(m.Adhoc) { // not required
		return nil
	}

	if err := m.Adhoc.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("adhoc")
		}
		return err
	}

	return nil
}

func (m *InstanceStatus) validateFlink(formats strfmt.Registry) error {

	if swag.IsZero(m.Flink) { // not required
		return nil
	}

	if m.Flink != nil {
		if err := m.Flink.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("flink")
			}
			return err
		}
	}

	return nil
}

func (m *InstanceStatus) validateKubernetes(formats strfmt.Registry) error {

	if swag.IsZero(m.Kubernetes) { // not required
		return nil
	}

	if m.Kubernetes != nil {
		if err := m.Kubernetes.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("kubernetes")
			}
			return err
		}
	}

	return nil
}

func (m *InstanceStatus) validateMarathon(formats strfmt.Registry) error {

	if swag.IsZero(m.Marathon) { // not required
		return nil
	}

	if m.Marathon != nil {
		if err := m.Marathon.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("marathon")
			}
			return err
		}
	}

	return nil
}

func (m *InstanceStatus) validateTron(formats strfmt.Registry) error {

	if swag.IsZero(m.Tron) { // not required
		return nil
	}

	if m.Tron != nil {
		if err := m.Tron.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("tron")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InstanceStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InstanceStatus) UnmarshalBinary(b []byte) error {
	var res InstanceStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// InstanceStatusFlink Nullable Flink instance status and metadata
//
// swagger:model InstanceStatusFlink
type InstanceStatusFlink struct {

	// metadata
	Metadata InstanceMetadataFlink `json:"metadata,omitempty"`

	// status
	Status *InstanceStatusFlink `json:"status,omitempty"`
}

// Validate validates this instance status flink
func (m *InstanceStatusFlink) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InstanceStatusFlink) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("flink" + "." + "status")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InstanceStatusFlink) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InstanceStatusFlink) UnmarshalBinary(b []byte) error {
	var res InstanceStatusFlink
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
