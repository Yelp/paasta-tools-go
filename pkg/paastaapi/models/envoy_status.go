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

// EnvoyStatus envoy status
//
// swagger:model EnvoyStatus
type EnvoyStatus struct {

	// Number of backends expected to be present in each location
	ExpectedBackendsPerLocation int32 `json:"expected_backends_per_location,omitempty"`

	// Locations the service is deployed
	Locations []*EnvoyLocation `json:"locations"`

	// Registration name of the service in Smartstack
	Registration string `json:"registration,omitempty"`
}

// Validate validates this envoy status
func (m *EnvoyStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLocations(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EnvoyStatus) validateLocations(formats strfmt.Registry) error {

	if swag.IsZero(m.Locations) { // not required
		return nil
	}

	for i := 0; i < len(m.Locations); i++ {
		if swag.IsZero(m.Locations[i]) { // not required
			continue
		}

		if m.Locations[i] != nil {
			if err := m.Locations[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("locations" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *EnvoyStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EnvoyStatus) UnmarshalBinary(b []byte) error {
	var res EnvoyStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
