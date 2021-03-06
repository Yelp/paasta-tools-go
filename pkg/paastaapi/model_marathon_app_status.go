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

// MarathonAppStatus struct for MarathonAppStatus
type MarathonAppStatus struct {
	// Backoff in seconds before launching next task
	BackoffSeconds *int32 `json:"backoff_seconds,omitempty"`
	// Unix timestamp when this app was created
	CreateTimestamp *float32 `json:"create_timestamp,omitempty"`
	// Marathon dashboard URL for this app
	DashboardUrl *string `json:"dashboard_url,omitempty"`
	// Deploy status of this app
	DeployStatus *string `json:"deploy_status,omitempty"`
	// Tasks associated to this app
	Tasks *[]MarathonTask `json:"tasks,omitempty"`
	// Number of healthy tasks for this app
	TasksHealthy *int32 `json:"tasks_healthy,omitempty"`
	// Number running tasks for this app
	TasksRunning *int32 `json:"tasks_running,omitempty"`
	// Number of staged tasks for this app
	TasksStaged *int32 `json:"tasks_staged,omitempty"`
	// Total number of tasks for this app
	TasksTotal *int32 `json:"tasks_total,omitempty"`
	// Mapping of reason offer was refused to the number of times that type of refusal was seen
	UnusedOfferReasonCounts *map[string]interface{} `json:"unused_offer_reason_counts,omitempty"`
	UnusedOffers *map[string]interface{} `json:"unused_offers,omitempty"`
}

// NewMarathonAppStatus instantiates a new MarathonAppStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMarathonAppStatus() *MarathonAppStatus {
	this := MarathonAppStatus{}
	return &this
}

// NewMarathonAppStatusWithDefaults instantiates a new MarathonAppStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMarathonAppStatusWithDefaults() *MarathonAppStatus {
	this := MarathonAppStatus{}
	return &this
}

// GetBackoffSeconds returns the BackoffSeconds field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetBackoffSeconds() int32 {
	if o == nil || o.BackoffSeconds == nil {
		var ret int32
		return ret
	}
	return *o.BackoffSeconds
}

// GetBackoffSecondsOk returns a tuple with the BackoffSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetBackoffSecondsOk() (*int32, bool) {
	if o == nil || o.BackoffSeconds == nil {
		return nil, false
	}
	return o.BackoffSeconds, true
}

// HasBackoffSeconds returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasBackoffSeconds() bool {
	if o != nil && o.BackoffSeconds != nil {
		return true
	}

	return false
}

// SetBackoffSeconds gets a reference to the given int32 and assigns it to the BackoffSeconds field.
func (o *MarathonAppStatus) SetBackoffSeconds(v int32) {
	o.BackoffSeconds = &v
}

// GetCreateTimestamp returns the CreateTimestamp field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetCreateTimestamp() float32 {
	if o == nil || o.CreateTimestamp == nil {
		var ret float32
		return ret
	}
	return *o.CreateTimestamp
}

// GetCreateTimestampOk returns a tuple with the CreateTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetCreateTimestampOk() (*float32, bool) {
	if o == nil || o.CreateTimestamp == nil {
		return nil, false
	}
	return o.CreateTimestamp, true
}

// HasCreateTimestamp returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasCreateTimestamp() bool {
	if o != nil && o.CreateTimestamp != nil {
		return true
	}

	return false
}

// SetCreateTimestamp gets a reference to the given float32 and assigns it to the CreateTimestamp field.
func (o *MarathonAppStatus) SetCreateTimestamp(v float32) {
	o.CreateTimestamp = &v
}

// GetDashboardUrl returns the DashboardUrl field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetDashboardUrl() string {
	if o == nil || o.DashboardUrl == nil {
		var ret string
		return ret
	}
	return *o.DashboardUrl
}

// GetDashboardUrlOk returns a tuple with the DashboardUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetDashboardUrlOk() (*string, bool) {
	if o == nil || o.DashboardUrl == nil {
		return nil, false
	}
	return o.DashboardUrl, true
}

// HasDashboardUrl returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasDashboardUrl() bool {
	if o != nil && o.DashboardUrl != nil {
		return true
	}

	return false
}

// SetDashboardUrl gets a reference to the given string and assigns it to the DashboardUrl field.
func (o *MarathonAppStatus) SetDashboardUrl(v string) {
	o.DashboardUrl = &v
}

// GetDeployStatus returns the DeployStatus field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetDeployStatus() string {
	if o == nil || o.DeployStatus == nil {
		var ret string
		return ret
	}
	return *o.DeployStatus
}

// GetDeployStatusOk returns a tuple with the DeployStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetDeployStatusOk() (*string, bool) {
	if o == nil || o.DeployStatus == nil {
		return nil, false
	}
	return o.DeployStatus, true
}

// HasDeployStatus returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasDeployStatus() bool {
	if o != nil && o.DeployStatus != nil {
		return true
	}

	return false
}

// SetDeployStatus gets a reference to the given string and assigns it to the DeployStatus field.
func (o *MarathonAppStatus) SetDeployStatus(v string) {
	o.DeployStatus = &v
}

// GetTasks returns the Tasks field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetTasks() []MarathonTask {
	if o == nil || o.Tasks == nil {
		var ret []MarathonTask
		return ret
	}
	return *o.Tasks
}

// GetTasksOk returns a tuple with the Tasks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetTasksOk() (*[]MarathonTask, bool) {
	if o == nil || o.Tasks == nil {
		return nil, false
	}
	return o.Tasks, true
}

// HasTasks returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasTasks() bool {
	if o != nil && o.Tasks != nil {
		return true
	}

	return false
}

// SetTasks gets a reference to the given []MarathonTask and assigns it to the Tasks field.
func (o *MarathonAppStatus) SetTasks(v []MarathonTask) {
	o.Tasks = &v
}

// GetTasksHealthy returns the TasksHealthy field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetTasksHealthy() int32 {
	if o == nil || o.TasksHealthy == nil {
		var ret int32
		return ret
	}
	return *o.TasksHealthy
}

// GetTasksHealthyOk returns a tuple with the TasksHealthy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetTasksHealthyOk() (*int32, bool) {
	if o == nil || o.TasksHealthy == nil {
		return nil, false
	}
	return o.TasksHealthy, true
}

// HasTasksHealthy returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasTasksHealthy() bool {
	if o != nil && o.TasksHealthy != nil {
		return true
	}

	return false
}

// SetTasksHealthy gets a reference to the given int32 and assigns it to the TasksHealthy field.
func (o *MarathonAppStatus) SetTasksHealthy(v int32) {
	o.TasksHealthy = &v
}

// GetTasksRunning returns the TasksRunning field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetTasksRunning() int32 {
	if o == nil || o.TasksRunning == nil {
		var ret int32
		return ret
	}
	return *o.TasksRunning
}

// GetTasksRunningOk returns a tuple with the TasksRunning field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetTasksRunningOk() (*int32, bool) {
	if o == nil || o.TasksRunning == nil {
		return nil, false
	}
	return o.TasksRunning, true
}

// HasTasksRunning returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasTasksRunning() bool {
	if o != nil && o.TasksRunning != nil {
		return true
	}

	return false
}

// SetTasksRunning gets a reference to the given int32 and assigns it to the TasksRunning field.
func (o *MarathonAppStatus) SetTasksRunning(v int32) {
	o.TasksRunning = &v
}

// GetTasksStaged returns the TasksStaged field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetTasksStaged() int32 {
	if o == nil || o.TasksStaged == nil {
		var ret int32
		return ret
	}
	return *o.TasksStaged
}

// GetTasksStagedOk returns a tuple with the TasksStaged field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetTasksStagedOk() (*int32, bool) {
	if o == nil || o.TasksStaged == nil {
		return nil, false
	}
	return o.TasksStaged, true
}

// HasTasksStaged returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasTasksStaged() bool {
	if o != nil && o.TasksStaged != nil {
		return true
	}

	return false
}

// SetTasksStaged gets a reference to the given int32 and assigns it to the TasksStaged field.
func (o *MarathonAppStatus) SetTasksStaged(v int32) {
	o.TasksStaged = &v
}

// GetTasksTotal returns the TasksTotal field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetTasksTotal() int32 {
	if o == nil || o.TasksTotal == nil {
		var ret int32
		return ret
	}
	return *o.TasksTotal
}

// GetTasksTotalOk returns a tuple with the TasksTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetTasksTotalOk() (*int32, bool) {
	if o == nil || o.TasksTotal == nil {
		return nil, false
	}
	return o.TasksTotal, true
}

// HasTasksTotal returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasTasksTotal() bool {
	if o != nil && o.TasksTotal != nil {
		return true
	}

	return false
}

// SetTasksTotal gets a reference to the given int32 and assigns it to the TasksTotal field.
func (o *MarathonAppStatus) SetTasksTotal(v int32) {
	o.TasksTotal = &v
}

// GetUnusedOfferReasonCounts returns the UnusedOfferReasonCounts field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetUnusedOfferReasonCounts() map[string]interface{} {
	if o == nil || o.UnusedOfferReasonCounts == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.UnusedOfferReasonCounts
}

// GetUnusedOfferReasonCountsOk returns a tuple with the UnusedOfferReasonCounts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetUnusedOfferReasonCountsOk() (*map[string]interface{}, bool) {
	if o == nil || o.UnusedOfferReasonCounts == nil {
		return nil, false
	}
	return o.UnusedOfferReasonCounts, true
}

// HasUnusedOfferReasonCounts returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasUnusedOfferReasonCounts() bool {
	if o != nil && o.UnusedOfferReasonCounts != nil {
		return true
	}

	return false
}

// SetUnusedOfferReasonCounts gets a reference to the given map[string]interface{} and assigns it to the UnusedOfferReasonCounts field.
func (o *MarathonAppStatus) SetUnusedOfferReasonCounts(v map[string]interface{}) {
	o.UnusedOfferReasonCounts = &v
}

// GetUnusedOffers returns the UnusedOffers field value if set, zero value otherwise.
func (o *MarathonAppStatus) GetUnusedOffers() map[string]interface{} {
	if o == nil || o.UnusedOffers == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.UnusedOffers
}

// GetUnusedOffersOk returns a tuple with the UnusedOffers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MarathonAppStatus) GetUnusedOffersOk() (*map[string]interface{}, bool) {
	if o == nil || o.UnusedOffers == nil {
		return nil, false
	}
	return o.UnusedOffers, true
}

// HasUnusedOffers returns a boolean if a field has been set.
func (o *MarathonAppStatus) HasUnusedOffers() bool {
	if o != nil && o.UnusedOffers != nil {
		return true
	}

	return false
}

// SetUnusedOffers gets a reference to the given map[string]interface{} and assigns it to the UnusedOffers field.
func (o *MarathonAppStatus) SetUnusedOffers(v map[string]interface{}) {
	o.UnusedOffers = &v
}

func (o MarathonAppStatus) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.BackoffSeconds != nil {
		toSerialize["backoff_seconds"] = o.BackoffSeconds
	}
	if o.CreateTimestamp != nil {
		toSerialize["create_timestamp"] = o.CreateTimestamp
	}
	if o.DashboardUrl != nil {
		toSerialize["dashboard_url"] = o.DashboardUrl
	}
	if o.DeployStatus != nil {
		toSerialize["deploy_status"] = o.DeployStatus
	}
	if o.Tasks != nil {
		toSerialize["tasks"] = o.Tasks
	}
	if o.TasksHealthy != nil {
		toSerialize["tasks_healthy"] = o.TasksHealthy
	}
	if o.TasksRunning != nil {
		toSerialize["tasks_running"] = o.TasksRunning
	}
	if o.TasksStaged != nil {
		toSerialize["tasks_staged"] = o.TasksStaged
	}
	if o.TasksTotal != nil {
		toSerialize["tasks_total"] = o.TasksTotal
	}
	if o.UnusedOfferReasonCounts != nil {
		toSerialize["unused_offer_reason_counts"] = o.UnusedOfferReasonCounts
	}
	if o.UnusedOffers != nil {
		toSerialize["unused_offers"] = o.UnusedOffers
	}
	return json.Marshal(toSerialize)
}

type NullableMarathonAppStatus struct {
	value *MarathonAppStatus
	isSet bool
}

func (v NullableMarathonAppStatus) Get() *MarathonAppStatus {
	return v.value
}

func (v *NullableMarathonAppStatus) Set(val *MarathonAppStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableMarathonAppStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableMarathonAppStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMarathonAppStatus(val *MarathonAppStatus) *NullableMarathonAppStatus {
	return &NullableMarathonAppStatus{value: val, isSet: true}
}

func (v NullableMarathonAppStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMarathonAppStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


