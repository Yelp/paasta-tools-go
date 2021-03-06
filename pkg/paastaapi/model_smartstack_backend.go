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

// SmartstackBackend struct for SmartstackBackend
type SmartstackBackend struct {
	// Check code reported by HAProxy
	CheckCode *string `json:"check_code,omitempty"`
	// Duration in ms of the last health check performed by HAProxy
	CheckDuration *int32 `json:"check_duration,omitempty"`
	// Status of last health check of the backend
	CheckStatus *string `json:"check_status,omitempty"`
	// Whether this backend has an associated task running
	HasAssociatedTask *bool `json:"has_associated_task,omitempty"`
	// Name of the host on which the backend is running
	Hostname *string `json:"hostname,omitempty"`
	// Seconds since last change in backend status
	LastChange *int32 `json:"last_change,omitempty"`
	// Port number on which the backend responds
	Port *int32 `json:"port,omitempty"`
	// Status of the backend in HAProxy
	Status *string `json:"status,omitempty"`
}

// NewSmartstackBackend instantiates a new SmartstackBackend object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSmartstackBackend() *SmartstackBackend {
	this := SmartstackBackend{}
	return &this
}

// NewSmartstackBackendWithDefaults instantiates a new SmartstackBackend object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSmartstackBackendWithDefaults() *SmartstackBackend {
	this := SmartstackBackend{}
	return &this
}

// GetCheckCode returns the CheckCode field value if set, zero value otherwise.
func (o *SmartstackBackend) GetCheckCode() string {
	if o == nil || o.CheckCode == nil {
		var ret string
		return ret
	}
	return *o.CheckCode
}

// GetCheckCodeOk returns a tuple with the CheckCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetCheckCodeOk() (*string, bool) {
	if o == nil || o.CheckCode == nil {
		return nil, false
	}
	return o.CheckCode, true
}

// HasCheckCode returns a boolean if a field has been set.
func (o *SmartstackBackend) HasCheckCode() bool {
	if o != nil && o.CheckCode != nil {
		return true
	}

	return false
}

// SetCheckCode gets a reference to the given string and assigns it to the CheckCode field.
func (o *SmartstackBackend) SetCheckCode(v string) {
	o.CheckCode = &v
}

// GetCheckDuration returns the CheckDuration field value if set, zero value otherwise.
func (o *SmartstackBackend) GetCheckDuration() int32 {
	if o == nil || o.CheckDuration == nil {
		var ret int32
		return ret
	}
	return *o.CheckDuration
}

// GetCheckDurationOk returns a tuple with the CheckDuration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetCheckDurationOk() (*int32, bool) {
	if o == nil || o.CheckDuration == nil {
		return nil, false
	}
	return o.CheckDuration, true
}

// HasCheckDuration returns a boolean if a field has been set.
func (o *SmartstackBackend) HasCheckDuration() bool {
	if o != nil && o.CheckDuration != nil {
		return true
	}

	return false
}

// SetCheckDuration gets a reference to the given int32 and assigns it to the CheckDuration field.
func (o *SmartstackBackend) SetCheckDuration(v int32) {
	o.CheckDuration = &v
}

// GetCheckStatus returns the CheckStatus field value if set, zero value otherwise.
func (o *SmartstackBackend) GetCheckStatus() string {
	if o == nil || o.CheckStatus == nil {
		var ret string
		return ret
	}
	return *o.CheckStatus
}

// GetCheckStatusOk returns a tuple with the CheckStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetCheckStatusOk() (*string, bool) {
	if o == nil || o.CheckStatus == nil {
		return nil, false
	}
	return o.CheckStatus, true
}

// HasCheckStatus returns a boolean if a field has been set.
func (o *SmartstackBackend) HasCheckStatus() bool {
	if o != nil && o.CheckStatus != nil {
		return true
	}

	return false
}

// SetCheckStatus gets a reference to the given string and assigns it to the CheckStatus field.
func (o *SmartstackBackend) SetCheckStatus(v string) {
	o.CheckStatus = &v
}

// GetHasAssociatedTask returns the HasAssociatedTask field value if set, zero value otherwise.
func (o *SmartstackBackend) GetHasAssociatedTask() bool {
	if o == nil || o.HasAssociatedTask == nil {
		var ret bool
		return ret
	}
	return *o.HasAssociatedTask
}

// GetHasAssociatedTaskOk returns a tuple with the HasAssociatedTask field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetHasAssociatedTaskOk() (*bool, bool) {
	if o == nil || o.HasAssociatedTask == nil {
		return nil, false
	}
	return o.HasAssociatedTask, true
}

// HasHasAssociatedTask returns a boolean if a field has been set.
func (o *SmartstackBackend) HasHasAssociatedTask() bool {
	if o != nil && o.HasAssociatedTask != nil {
		return true
	}

	return false
}

// SetHasAssociatedTask gets a reference to the given bool and assigns it to the HasAssociatedTask field.
func (o *SmartstackBackend) SetHasAssociatedTask(v bool) {
	o.HasAssociatedTask = &v
}

// GetHostname returns the Hostname field value if set, zero value otherwise.
func (o *SmartstackBackend) GetHostname() string {
	if o == nil || o.Hostname == nil {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetHostnameOk() (*string, bool) {
	if o == nil || o.Hostname == nil {
		return nil, false
	}
	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *SmartstackBackend) HasHostname() bool {
	if o != nil && o.Hostname != nil {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *SmartstackBackend) SetHostname(v string) {
	o.Hostname = &v
}

// GetLastChange returns the LastChange field value if set, zero value otherwise.
func (o *SmartstackBackend) GetLastChange() int32 {
	if o == nil || o.LastChange == nil {
		var ret int32
		return ret
	}
	return *o.LastChange
}

// GetLastChangeOk returns a tuple with the LastChange field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetLastChangeOk() (*int32, bool) {
	if o == nil || o.LastChange == nil {
		return nil, false
	}
	return o.LastChange, true
}

// HasLastChange returns a boolean if a field has been set.
func (o *SmartstackBackend) HasLastChange() bool {
	if o != nil && o.LastChange != nil {
		return true
	}

	return false
}

// SetLastChange gets a reference to the given int32 and assigns it to the LastChange field.
func (o *SmartstackBackend) SetLastChange(v int32) {
	o.LastChange = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *SmartstackBackend) GetPort() int32 {
	if o == nil || o.Port == nil {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetPortOk() (*int32, bool) {
	if o == nil || o.Port == nil {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *SmartstackBackend) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *SmartstackBackend) SetPort(v int32) {
	o.Port = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *SmartstackBackend) GetStatus() string {
	if o == nil || o.Status == nil {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SmartstackBackend) GetStatusOk() (*string, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *SmartstackBackend) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *SmartstackBackend) SetStatus(v string) {
	o.Status = &v
}

func (o SmartstackBackend) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CheckCode != nil {
		toSerialize["check_code"] = o.CheckCode
	}
	if o.CheckDuration != nil {
		toSerialize["check_duration"] = o.CheckDuration
	}
	if o.CheckStatus != nil {
		toSerialize["check_status"] = o.CheckStatus
	}
	if o.HasAssociatedTask != nil {
		toSerialize["has_associated_task"] = o.HasAssociatedTask
	}
	if o.Hostname != nil {
		toSerialize["hostname"] = o.Hostname
	}
	if o.LastChange != nil {
		toSerialize["last_change"] = o.LastChange
	}
	if o.Port != nil {
		toSerialize["port"] = o.Port
	}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	return json.Marshal(toSerialize)
}

type NullableSmartstackBackend struct {
	value *SmartstackBackend
	isSet bool
}

func (v NullableSmartstackBackend) Get() *SmartstackBackend {
	return v.value
}

func (v *NullableSmartstackBackend) Set(val *SmartstackBackend) {
	v.value = val
	v.isSet = true
}

func (v NullableSmartstackBackend) IsSet() bool {
	return v.isSet
}

func (v *NullableSmartstackBackend) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSmartstackBackend(val *SmartstackBackend) *NullableSmartstackBackend {
	return &NullableSmartstackBackend{value: val, isSet: true}
}

func (v NullableSmartstackBackend) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSmartstackBackend) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


