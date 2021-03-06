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

// MarathonMesosNonrunningTask struct for MarathonMesosNonrunningTask
type MarathonMesosNonrunningTask struct {
	// The unix timestamp at which the task was deployed
	DeployedTimestamp *float32 `json:"deployed_timestamp,omitempty"`
	// Name of the Mesos agent on which this task is running
	Hostname *string `json:"hostname,omitempty"`
	// The ID of the task in Mesos
	Id *string `json:"id,omitempty"`
	// The current state of the task
	State *string `json:"state,omitempty"`
	TailLines *TaskTailLines `json:"tail_lines,omitempty"`
}

// NewMarathonMesosNonrunningTask instantiates a new MarathonMesosNonrunningTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMarathonMesosNonrunningTask() *MarathonMesosNonrunningTask {
	this := MarathonMesosNonrunningTask{}
	return &this
}

// NewMarathonMesosNonrunningTaskWithDefaults instantiates a new MarathonMesosNonrunningTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMarathonMesosNonrunningTaskWithDefaults() *MarathonMesosNonrunningTask {
	this := MarathonMesosNonrunningTask{}
	return &this
}

// GetDeployedTimestamp returns the DeployedTimestamp field value if set, zero value otherwise.
func (o *MarathonMesosNonrunningTask) GetDeployedTimestamp() float32 {
	if o == nil || o.DeployedTimestamp == nil {
		var ret float32
		return ret
	}
	return *o.DeployedTimestamp
}

// GetDeployedTimestampOk returns a tuple with the DeployedTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonMesosNonrunningTask) GetDeployedTimestampOk() (*float32, bool) {
	if o == nil || o.DeployedTimestamp == nil {
		return nil, false
	}
	return o.DeployedTimestamp, true
}

// HasDeployedTimestamp returns a boolean if a field has been set.
func (o *MarathonMesosNonrunningTask) HasDeployedTimestamp() bool {
	if o != nil && o.DeployedTimestamp != nil {
		return true
	}

	return false
}

// SetDeployedTimestamp gets a reference to the given float32 and assigns it to the DeployedTimestamp field.
func (o *MarathonMesosNonrunningTask) SetDeployedTimestamp(v float32) {
	o.DeployedTimestamp = &v
}

// GetHostname returns the Hostname field value if set, zero value otherwise.
func (o *MarathonMesosNonrunningTask) GetHostname() string {
	if o == nil || o.Hostname == nil {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonMesosNonrunningTask) GetHostnameOk() (*string, bool) {
	if o == nil || o.Hostname == nil {
		return nil, false
	}
	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *MarathonMesosNonrunningTask) HasHostname() bool {
	if o != nil && o.Hostname != nil {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *MarathonMesosNonrunningTask) SetHostname(v string) {
	o.Hostname = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *MarathonMesosNonrunningTask) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonMesosNonrunningTask) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *MarathonMesosNonrunningTask) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *MarathonMesosNonrunningTask) SetId(v string) {
	o.Id = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *MarathonMesosNonrunningTask) GetState() string {
	if o == nil || o.State == nil {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonMesosNonrunningTask) GetStateOk() (*string, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *MarathonMesosNonrunningTask) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *MarathonMesosNonrunningTask) SetState(v string) {
	o.State = &v
}

// GetTailLines returns the TailLines field value if set, zero value otherwise.
func (o *MarathonMesosNonrunningTask) GetTailLines() TaskTailLines {
	if o == nil || o.TailLines == nil {
		var ret TaskTailLines
		return ret
	}
	return *o.TailLines
}

// GetTailLinesOk returns a tuple with the TailLines field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonMesosNonrunningTask) GetTailLinesOk() (*TaskTailLines, bool) {
	if o == nil || o.TailLines == nil {
		return nil, false
	}
	return o.TailLines, true
}

// HasTailLines returns a boolean if a field has been set.
func (o *MarathonMesosNonrunningTask) HasTailLines() bool {
	if o != nil && o.TailLines != nil {
		return true
	}

	return false
}

// SetTailLines gets a reference to the given TaskTailLines and assigns it to the TailLines field.
func (o *MarathonMesosNonrunningTask) SetTailLines(v TaskTailLines) {
	o.TailLines = &v
}

func (o MarathonMesosNonrunningTask) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DeployedTimestamp != nil {
		toSerialize["deployed_timestamp"] = o.DeployedTimestamp
	}
	if o.Hostname != nil {
		toSerialize["hostname"] = o.Hostname
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.State != nil {
		toSerialize["state"] = o.State
	}
	if o.TailLines != nil {
		toSerialize["tail_lines"] = o.TailLines
	}
	return json.Marshal(toSerialize)
}

type NullableMarathonMesosNonrunningTask struct {
	value *MarathonMesosNonrunningTask
	isSet bool
}

func (v NullableMarathonMesosNonrunningTask) Get() *MarathonMesosNonrunningTask {
	return v.value
}

func (v *NullableMarathonMesosNonrunningTask) Set(val *MarathonMesosNonrunningTask) {
	v.value = val
	v.isSet = true
}

func (v NullableMarathonMesosNonrunningTask) IsSet() bool {
	return v.isSet
}

func (v *NullableMarathonMesosNonrunningTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMarathonMesosNonrunningTask(val *MarathonMesosNonrunningTask) *NullableMarathonMesosNonrunningTask {
	return &NullableMarathonMesosNonrunningTask{value: val, isSet: true}
}

func (v NullableMarathonMesosNonrunningTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMarathonMesosNonrunningTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


