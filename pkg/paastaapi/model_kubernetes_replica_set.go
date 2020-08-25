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

// KubernetesReplicaSet struct for KubernetesReplicaSet
type KubernetesReplicaSet struct {
	// Time at which the replicaset was created
	CreateTimestamp *float32 `json:"create_timestamp,omitempty"`
	// name of the replicaset in Kubernetes
	Name *string `json:"name,omitempty"`
	// number of ready replicas for the replicaset
	ReadyReplicas *int32 `json:"ready_replicas,omitempty"`
	// number of desired replicas for the replicaset
	Replicas *int32 `json:"replicas,omitempty"`
}

// NewKubernetesReplicaSet instantiates a new KubernetesReplicaSet object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKubernetesReplicaSet() *KubernetesReplicaSet {
	this := KubernetesReplicaSet{}
	return &this
}

// NewKubernetesReplicaSetWithDefaults instantiates a new KubernetesReplicaSet object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKubernetesReplicaSetWithDefaults() *KubernetesReplicaSet {
	this := KubernetesReplicaSet{}
	return &this
}

// GetCreateTimestamp returns the CreateTimestamp field value if set, zero value otherwise.
func (o *KubernetesReplicaSet) GetCreateTimestamp() float32 {
	if o == nil || o.CreateTimestamp == nil {
		var ret float32
		return ret
	}
	return *o.CreateTimestamp
}

// GetCreateTimestampOk returns a tuple with the CreateTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KubernetesReplicaSet) GetCreateTimestampOk() (*float32, bool) {
	if o == nil || o.CreateTimestamp == nil {
		return nil, false
	}
	return o.CreateTimestamp, true
}

// HasCreateTimestamp returns a boolean if a field has been set.
func (o *KubernetesReplicaSet) HasCreateTimestamp() bool {
	if o != nil && o.CreateTimestamp != nil {
		return true
	}

	return false
}

// SetCreateTimestamp gets a reference to the given float32 and assigns it to the CreateTimestamp field.
func (o *KubernetesReplicaSet) SetCreateTimestamp(v float32) {
	o.CreateTimestamp = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *KubernetesReplicaSet) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KubernetesReplicaSet) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *KubernetesReplicaSet) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *KubernetesReplicaSet) SetName(v string) {
	o.Name = &v
}

// GetReadyReplicas returns the ReadyReplicas field value if set, zero value otherwise.
func (o *KubernetesReplicaSet) GetReadyReplicas() int32 {
	if o == nil || o.ReadyReplicas == nil {
		var ret int32
		return ret
	}
	return *o.ReadyReplicas
}

// GetReadyReplicasOk returns a tuple with the ReadyReplicas field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KubernetesReplicaSet) GetReadyReplicasOk() (*int32, bool) {
	if o == nil || o.ReadyReplicas == nil {
		return nil, false
	}
	return o.ReadyReplicas, true
}

// HasReadyReplicas returns a boolean if a field has been set.
func (o *KubernetesReplicaSet) HasReadyReplicas() bool {
	if o != nil && o.ReadyReplicas != nil {
		return true
	}

	return false
}

// SetReadyReplicas gets a reference to the given int32 and assigns it to the ReadyReplicas field.
func (o *KubernetesReplicaSet) SetReadyReplicas(v int32) {
	o.ReadyReplicas = &v
}

// GetReplicas returns the Replicas field value if set, zero value otherwise.
func (o *KubernetesReplicaSet) GetReplicas() int32 {
	if o == nil || o.Replicas == nil {
		var ret int32
		return ret
	}
	return *o.Replicas
}

// GetReplicasOk returns a tuple with the Replicas field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KubernetesReplicaSet) GetReplicasOk() (*int32, bool) {
	if o == nil || o.Replicas == nil {
		return nil, false
	}
	return o.Replicas, true
}

// HasReplicas returns a boolean if a field has been set.
func (o *KubernetesReplicaSet) HasReplicas() bool {
	if o != nil && o.Replicas != nil {
		return true
	}

	return false
}

// SetReplicas gets a reference to the given int32 and assigns it to the Replicas field.
func (o *KubernetesReplicaSet) SetReplicas(v int32) {
	o.Replicas = &v
}

func (o KubernetesReplicaSet) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CreateTimestamp != nil {
		toSerialize["create_timestamp"] = o.CreateTimestamp
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.ReadyReplicas != nil {
		toSerialize["ready_replicas"] = o.ReadyReplicas
	}
	if o.Replicas != nil {
		toSerialize["replicas"] = o.Replicas
	}
	return json.Marshal(toSerialize)
}

type NullableKubernetesReplicaSet struct {
	value *KubernetesReplicaSet
	isSet bool
}

func (v NullableKubernetesReplicaSet) Get() *KubernetesReplicaSet {
	return v.value
}

func (v *NullableKubernetesReplicaSet) Set(val *KubernetesReplicaSet) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesReplicaSet) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesReplicaSet) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesReplicaSet(val *KubernetesReplicaSet) *NullableKubernetesReplicaSet {
	return &NullableKubernetesReplicaSet{value: val, isSet: true}
}

func (v NullableKubernetesReplicaSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesReplicaSet) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


