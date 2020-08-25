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

// AdhocLaunchHistory A single run
type AdhocLaunchHistory struct {
	// framework id
	FrameworkId *string `json:"framework_id,omitempty"`
	// when the job was launched
	LaunchTime *string `json:"launch_time,omitempty"`
	// id of the single run
	RunId *string `json:"run_id,omitempty"`
}

// NewAdhocLaunchHistory instantiates a new AdhocLaunchHistory object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAdhocLaunchHistory() *AdhocLaunchHistory {
	this := AdhocLaunchHistory{}
	return &this
}

// NewAdhocLaunchHistoryWithDefaults instantiates a new AdhocLaunchHistory object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAdhocLaunchHistoryWithDefaults() *AdhocLaunchHistory {
	this := AdhocLaunchHistory{}
	return &this
}

// GetFrameworkId returns the FrameworkId field value if set, zero value otherwise.
func (o *AdhocLaunchHistory) GetFrameworkId() string {
	if o == nil || o.FrameworkId == nil {
		var ret string
		return ret
	}
	return *o.FrameworkId
}

// GetFrameworkIdOk returns a tuple with the FrameworkId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdhocLaunchHistory) GetFrameworkIdOk() (*string, bool) {
	if o == nil || o.FrameworkId == nil {
		return nil, false
	}
	return o.FrameworkId, true
}

// HasFrameworkId returns a boolean if a field has been set.
func (o *AdhocLaunchHistory) HasFrameworkId() bool {
	if o != nil && o.FrameworkId != nil {
		return true
	}

	return false
}

// SetFrameworkId gets a reference to the given string and assigns it to the FrameworkId field.
func (o *AdhocLaunchHistory) SetFrameworkId(v string) {
	o.FrameworkId = &v
}

// GetLaunchTime returns the LaunchTime field value if set, zero value otherwise.
func (o *AdhocLaunchHistory) GetLaunchTime() string {
	if o == nil || o.LaunchTime == nil {
		var ret string
		return ret
	}
	return *o.LaunchTime
}

// GetLaunchTimeOk returns a tuple with the LaunchTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdhocLaunchHistory) GetLaunchTimeOk() (*string, bool) {
	if o == nil || o.LaunchTime == nil {
		return nil, false
	}
	return o.LaunchTime, true
}

// HasLaunchTime returns a boolean if a field has been set.
func (o *AdhocLaunchHistory) HasLaunchTime() bool {
	if o != nil && o.LaunchTime != nil {
		return true
	}

	return false
}

// SetLaunchTime gets a reference to the given string and assigns it to the LaunchTime field.
func (o *AdhocLaunchHistory) SetLaunchTime(v string) {
	o.LaunchTime = &v
}

// GetRunId returns the RunId field value if set, zero value otherwise.
func (o *AdhocLaunchHistory) GetRunId() string {
	if o == nil || o.RunId == nil {
		var ret string
		return ret
	}
	return *o.RunId
}

// GetRunIdOk returns a tuple with the RunId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdhocLaunchHistory) GetRunIdOk() (*string, bool) {
	if o == nil || o.RunId == nil {
		return nil, false
	}
	return o.RunId, true
}

// HasRunId returns a boolean if a field has been set.
func (o *AdhocLaunchHistory) HasRunId() bool {
	if o != nil && o.RunId != nil {
		return true
	}

	return false
}

// SetRunId gets a reference to the given string and assigns it to the RunId field.
func (o *AdhocLaunchHistory) SetRunId(v string) {
	o.RunId = &v
}

func (o AdhocLaunchHistory) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.FrameworkId != nil {
		toSerialize["framework_id"] = o.FrameworkId
	}
	if o.LaunchTime != nil {
		toSerialize["launch_time"] = o.LaunchTime
	}
	if o.RunId != nil {
		toSerialize["run_id"] = o.RunId
	}
	return json.Marshal(toSerialize)
}

type NullableAdhocLaunchHistory struct {
	value *AdhocLaunchHistory
	isSet bool
}

func (v NullableAdhocLaunchHistory) Get() *AdhocLaunchHistory {
	return v.value
}

func (v *NullableAdhocLaunchHistory) Set(val *AdhocLaunchHistory) {
	v.value = val
	v.isSet = true
}

func (v NullableAdhocLaunchHistory) IsSet() bool {
	return v.isSet
}

func (v *NullableAdhocLaunchHistory) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAdhocLaunchHistory(val *AdhocLaunchHistory) *NullableAdhocLaunchHistory {
	return &NullableAdhocLaunchHistory{value: val, isSet: true}
}

func (v NullableAdhocLaunchHistory) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAdhocLaunchHistory) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


