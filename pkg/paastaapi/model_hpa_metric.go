/*
 * Paasta API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package paastaapi

import (
	"encoding/json"
)

// HPAMetric struct for HPAMetric
type HPAMetric struct {
	// setpoint/target_value as specified in yelpsoa_configs
	CurrentValue *string `json:"current_value,omitempty"`
	// name of the metric
	Name *string `json:"name,omitempty"`
	// setpoint/target_value as specified in yelpsoa_configs
	TargetValue *string `json:"target_value,omitempty"`
}

// NewHPAMetric instantiates a new HPAMetric object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHPAMetric() *HPAMetric {
	this := HPAMetric{}
	return &this
}

// NewHPAMetricWithDefaults instantiates a new HPAMetric object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHPAMetricWithDefaults() *HPAMetric {
	this := HPAMetric{}
	return &this
}

// GetCurrentValue returns the CurrentValue field value if set, zero value otherwise.
func (o *HPAMetric) GetCurrentValue() string {
	if o == nil || o.CurrentValue == nil {
		var ret string
		return ret
	}
	return *o.CurrentValue
}

// GetCurrentValueOk returns a tuple with the CurrentValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HPAMetric) GetCurrentValueOk() (*string, bool) {
	if o == nil || o.CurrentValue == nil {
		return nil, false
	}
	return o.CurrentValue, true
}

// HasCurrentValue returns a boolean if a field has been set.
func (o *HPAMetric) HasCurrentValue() bool {
	if o != nil && o.CurrentValue != nil {
		return true
	}

	return false
}

// SetCurrentValue gets a reference to the given string and assigns it to the CurrentValue field.
func (o *HPAMetric) SetCurrentValue(v string) {
	o.CurrentValue = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *HPAMetric) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HPAMetric) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *HPAMetric) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *HPAMetric) SetName(v string) {
	o.Name = &v
}

// GetTargetValue returns the TargetValue field value if set, zero value otherwise.
func (o *HPAMetric) GetTargetValue() string {
	if o == nil || o.TargetValue == nil {
		var ret string
		return ret
	}
	return *o.TargetValue
}

// GetTargetValueOk returns a tuple with the TargetValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HPAMetric) GetTargetValueOk() (*string, bool) {
	if o == nil || o.TargetValue == nil {
		return nil, false
	}
	return o.TargetValue, true
}

// HasTargetValue returns a boolean if a field has been set.
func (o *HPAMetric) HasTargetValue() bool {
	if o != nil && o.TargetValue != nil {
		return true
	}

	return false
}

// SetTargetValue gets a reference to the given string and assigns it to the TargetValue field.
func (o *HPAMetric) SetTargetValue(v string) {
	o.TargetValue = &v
}

func (o HPAMetric) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CurrentValue != nil {
		toSerialize["current_value"] = o.CurrentValue
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.TargetValue != nil {
		toSerialize["target_value"] = o.TargetValue
	}
	return json.Marshal(toSerialize)
}

type NullableHPAMetric struct {
	value *HPAMetric
	isSet bool
}

func (v NullableHPAMetric) Get() *HPAMetric {
	return v.value
}

func (v *NullableHPAMetric) Set(val *HPAMetric) {
	v.value = val
	v.isSet = true
}

func (v NullableHPAMetric) IsSet() bool {
	return v.isSet
}

func (v *NullableHPAMetric) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHPAMetric(val *HPAMetric) *NullableHPAMetric {
	return &NullableHPAMetric{value: val, isSet: true}
}

func (v NullableHPAMetric) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHPAMetric) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


