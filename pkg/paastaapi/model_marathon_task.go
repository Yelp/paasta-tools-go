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

// MarathonTask struct for MarathonTask
type MarathonTask struct {
	// Time at which the task was deployed
	DeployedTimestamp *float32 `json:"deployed_timestamp,omitempty"`
	// Name of the host on which the task is running
	Host NullableString `json:"host,omitempty"`
	// ID of the task in Mesos
	Id *string `json:"id,omitempty"`
	// Whether Marathon thinks the task is healthy
	IsHealthy NullableBool `json:"is_healthy,omitempty"`
	// Port on which the task is listening
	Port *int32 `json:"port,omitempty"`
}

// NewMarathonTask instantiates a new MarathonTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMarathonTask() *MarathonTask {
	this := MarathonTask{}
	return &this
}

// NewMarathonTaskWithDefaults instantiates a new MarathonTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMarathonTaskWithDefaults() *MarathonTask {
	this := MarathonTask{}
	return &this
}

// GetDeployedTimestamp returns the DeployedTimestamp field value if set, zero value otherwise.
func (o *MarathonTask) GetDeployedTimestamp() float32 {
	if o == nil || o.DeployedTimestamp == nil {
		var ret float32
		return ret
	}
	return *o.DeployedTimestamp
}

// GetDeployedTimestampOk returns a tuple with the DeployedTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonTask) GetDeployedTimestampOk() (*float32, bool) {
	if o == nil || o.DeployedTimestamp == nil {
		return nil, false
	}
	return o.DeployedTimestamp, true
}

// HasDeployedTimestamp returns a boolean if a field has been set.
func (o *MarathonTask) HasDeployedTimestamp() bool {
	if o != nil && o.DeployedTimestamp != nil {
		return true
	}

	return false
}

// SetDeployedTimestamp gets a reference to the given float32 and assigns it to the DeployedTimestamp field.
func (o *MarathonTask) SetDeployedTimestamp(v float32) {
	o.DeployedTimestamp = &v
}

// GetHost returns the Host field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *MarathonTask) GetHost() string {
	if o == nil || o.Host.Get() == nil {
		var ret string
		return ret
	}
	return *o.Host.Get()
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MarathonTask) GetHostOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Host.Get(), o.Host.IsSet()
}

// HasHost returns a boolean if a field has been set.
func (o *MarathonTask) HasHost() bool {
	if o != nil && o.Host.IsSet() {
		return true
	}

	return false
}

// SetHost gets a reference to the given NullableString and assigns it to the Host field.
func (o *MarathonTask) SetHost(v string) {
	o.Host.Set(&v)
}
// SetHostNil sets the value for Host to be an explicit nil
func (o *MarathonTask) SetHostNil() {
	o.Host.Set(nil)
}

// UnsetHost ensures that no value is present for Host, not even an explicit nil
func (o *MarathonTask) UnsetHost() {
	o.Host.Unset()
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *MarathonTask) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonTask) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *MarathonTask) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *MarathonTask) SetId(v string) {
	o.Id = &v
}

// GetIsHealthy returns the IsHealthy field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *MarathonTask) GetIsHealthy() bool {
	if o == nil || o.IsHealthy.Get() == nil {
		var ret bool
		return ret
	}
	return *o.IsHealthy.Get()
}

// GetIsHealthyOk returns a tuple with the IsHealthy field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MarathonTask) GetIsHealthyOk() (*bool, bool) {
	if o == nil  {
		return nil, false
	}
	return o.IsHealthy.Get(), o.IsHealthy.IsSet()
}

// HasIsHealthy returns a boolean if a field has been set.
func (o *MarathonTask) HasIsHealthy() bool {
	if o != nil && o.IsHealthy.IsSet() {
		return true
	}

	return false
}

// SetIsHealthy gets a reference to the given NullableBool and assigns it to the IsHealthy field.
func (o *MarathonTask) SetIsHealthy(v bool) {
	o.IsHealthy.Set(&v)
}
// SetIsHealthyNil sets the value for IsHealthy to be an explicit nil
func (o *MarathonTask) SetIsHealthyNil() {
	o.IsHealthy.Set(nil)
}

// UnsetIsHealthy ensures that no value is present for IsHealthy, not even an explicit nil
func (o *MarathonTask) UnsetIsHealthy() {
	o.IsHealthy.Unset()
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *MarathonTask) GetPort() int32 {
	if o == nil || o.Port == nil {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonTask) GetPortOk() (*int32, bool) {
	if o == nil || o.Port == nil {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *MarathonTask) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *MarathonTask) SetPort(v int32) {
	o.Port = &v
}

func (o MarathonTask) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DeployedTimestamp != nil {
		toSerialize["deployed_timestamp"] = o.DeployedTimestamp
	}
	if o.Host.IsSet() {
		toSerialize["host"] = o.Host.Get()
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.IsHealthy.IsSet() {
		toSerialize["is_healthy"] = o.IsHealthy.Get()
	}
	if o.Port != nil {
		toSerialize["port"] = o.Port
	}
	return json.Marshal(toSerialize)
}

type NullableMarathonTask struct {
	value *MarathonTask
	isSet bool
}

func (v NullableMarathonTask) Get() *MarathonTask {
	return v.value
}

func (v *NullableMarathonTask) Set(val *MarathonTask) {
	v.value = val
	v.isSet = true
}

func (v NullableMarathonTask) IsSet() bool {
	return v.isSet
}

func (v *NullableMarathonTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMarathonTask(val *MarathonTask) *NullableMarathonTask {
	return &NullableMarathonTask{value: val, isSet: true}
}

func (v NullableMarathonTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMarathonTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


