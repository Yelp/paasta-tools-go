package containerspec

import (
	"encoding/json"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

const (
	defaultCPU    = KubeResourceQuantity("0.1")
	defaultMemory = KubeResourceQuantity("512")
)

// KubeResourceQuantity : Resource quantity for Kubernetes (e.g.; CPU, mem, disk)
type KubeResourceQuantity string

// UnmarshalJSON : unmarshal the JSON representation of a KubeResourceQuantity
func (n *KubeResourceQuantity) UnmarshalJSON(b []byte) error {
	if len(b) > 1 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	*n = KubeResourceQuantity(string(b))
	return nil
}

// MarshalJSON : marshal the JSON representation of a KubeResourceQuantity
func (n KubeResourceQuantity) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(n))
}

// PaastaContainerSpec : Spec for any paasta container with basic fields and utilities
type PaastaContainerSpec struct {
	CPU    *KubeResourceQuantity `json:"cpus"`
	Memory *KubeResourceQuantity `json:"mem"`
}

// GetContainerResources : get resource requirements based on the container spec
func (spec *PaastaContainerSpec) GetContainerResources() (*corev1.ResourceRequirements, error) {
	var cpu KubeResourceQuantity
	if spec.CPU != nil {
		cpu = *spec.CPU
	} else {
		cpu = defaultCPU
	}
	cpuQuantity, err := resource.ParseQuantity(string(cpu))
	if err != nil {
		return nil, fmt.Errorf("error while parsing cpu request '%s': %s", cpu, err)
	}

	var memory KubeResourceQuantity
	if spec.Memory != nil {
		memory = *spec.Memory
	} else {
		memory = defaultMemory
	}
	if _, err := strconv.Atoi(string(memory)); err == nil {
		// memory value looks like a number, let's treat it as MB according to PaaSTA default
		memory = memory + "M"
	}
	memoryQuantity, err := resource.ParseQuantity(string(memory))
	if err != nil {
		return nil, fmt.Errorf("error while parsing memory '%s': %s", memory, err)
	}

	return &corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpuQuantity,
			corev1.ResourceMemory: memoryQuantity,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpuQuantity,
			corev1.ResourceMemory: memoryQuantity,
		},
	}, nil
}
