/*
 * Paasta API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package paastaapi

import (
	"encoding/json"
)

// InstanceStatusFlink Nullable Flink instance status and metadata
type InstanceStatusFlink struct {
	// Flink instance metadata
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	// Flink instance status
	Status *map[string]interface{} `json:"status,omitempty"`
}

// NewInstanceStatusFlink instantiates a new InstanceStatusFlink object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInstanceStatusFlink() *InstanceStatusFlink {
	this := InstanceStatusFlink{}
	return &this
}

// NewInstanceStatusFlinkWithDefaults instantiates a new InstanceStatusFlink object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInstanceStatusFlinkWithDefaults() *InstanceStatusFlink {
	this := InstanceStatusFlink{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *InstanceStatusFlink) GetMetadata() map[string]interface{} {
	if o == nil || o.Metadata == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InstanceStatusFlink) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *InstanceStatusFlink) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *InstanceStatusFlink) SetMetadata(v map[string]interface{}) {
	o.Metadata = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *InstanceStatusFlink) GetStatus() map[string]interface{} {
	if o == nil || o.Status == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InstanceStatusFlink) GetStatusOk() (*map[string]interface{}, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *InstanceStatusFlink) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given map[string]interface{} and assigns it to the Status field.
func (o *InstanceStatusFlink) SetStatus(v map[string]interface{}) {
	o.Status = &v
}

func (o InstanceStatusFlink) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	return json.Marshal(toSerialize)
}

type NullableInstanceStatusFlink struct {
	value *InstanceStatusFlink
	isSet bool
}

func (v NullableInstanceStatusFlink) Get() *InstanceStatusFlink {
	return v.value
}

func (v *NullableInstanceStatusFlink) Set(val *InstanceStatusFlink) {
	v.value = val
	v.isSet = true
}

func (v NullableInstanceStatusFlink) IsSet() bool {
	return v.isSet
}

func (v *NullableInstanceStatusFlink) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInstanceStatusFlink(val *InstanceStatusFlink) *NullableInstanceStatusFlink {
	return &NullableInstanceStatusFlink{value: val, isSet: true}
}

func (v NullableInstanceStatusFlink) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInstanceStatusFlink) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


