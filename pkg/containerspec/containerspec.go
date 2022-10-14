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
	defaultMemory = KubeResourceQuantity("512Mi")
	defaultDisk   = KubeResourceQuantity("1024Mi")
)

var defaultSpec = newDefaultSpec()

// KubeResourceQuantity : Resource quantity for Kubernetes (e.g.; CPU, mem, disk)
type KubeResourceQuantity string

func (n *KubeResourceQuantity) withSuffix() KubeResourceQuantity {
	if _, err := strconv.Atoi(string(*n)); err == nil {
		// value looks like a number, let's treat it as MB according to PaaSTA default
		return *n + "Mi"
	}
	return *n
}

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

func (n KubeResourceQuantity) copy() KubeResourceQuantity {
	return KubeResourceQuantity(string(n))
}

func (n KubeResourceQuantity) Cmp(quantityToCompare KubeResourceQuantity) (int, error) {
	quantityResource, err := resource.ParseQuantity(string(n))
	if err != nil {
		return 0, fmt.Errorf("error while parsing KubeResourceQuantity quantity '%s': %s", string(n), err)
	}

	quantityToCompareResource, err := resource.ParseQuantity(string(quantityToCompare))
	if err != nil {
		return 0, fmt.Errorf("error while parsing KubeResourceQuantity quantity '%s': %s", string(quantityToCompare), err)
	}

	// Cmp returns 0 if the quantityResource is equal to quantityToCompareResource
	// -1 if the quantityResource is less than quantityToCompareResource
	// 1 if the quantityResource is greater than quantityToCompareResource.
	return quantityResource.Cmp(quantityToCompareResource), nil
}

// PaastaContainerSpec : Spec for any paasta container with basic fields and utilities
type PaastaContainerSpec struct {
	CPU         *KubeResourceQuantity `json:"cpus"`
	CPULimit    *KubeResourceQuantity `json:"cpus_limit,omitempty"`
	Memory      *KubeResourceQuantity `json:"mem"`
	MemoryLimit *KubeResourceQuantity `json:"mem_limit,omitempty"`
	Disk        *KubeResourceQuantity `json:"disk"`
	DiskLimit   *KubeResourceQuantity `json:"disk_limit,omitempty"`
}

// GetContainerResources : get resource requirements based on the container spec
func (spec *PaastaContainerSpec) GetContainerResources() (*corev1.ResourceRequirements, error) {
	return spec.GetContainerResourcesWithDefaults(defaultSpec)
}

// GetContainerResources : get resource requirements based on the container spec using defaults when necessary
func (spec *PaastaContainerSpec) GetContainerResourcesWithDefaults(defaults *PaastaContainerSpec) (*corev1.ResourceRequirements, error) {
	var cpu KubeResourceQuantity
	if spec.CPU != nil {
		cpu = *spec.CPU
	} else {
		cpu = *defaults.CPU
	}
	cpuQuantity, err := resource.ParseQuantity(string(cpu))
	if err != nil {
		return nil, fmt.Errorf("error while parsing cpu request '%s': %s", cpu, err)
	}

	var cpuLimit KubeResourceQuantity
	if spec.CPULimit != nil {
		cpuLimit = *spec.CPULimit
	} else {
		cpuLimit = cpu
	}

	// Note, there is a gotcha in resource.Quantity.Value(): it rounds to full numbers.
	// Since CPUs resources can be fractional, we have to rely on MilliValue() instead.
	cpuLimitQuantity, err := resource.ParseQuantity(string(cpuLimit))
	if err != nil {
		return nil, fmt.Errorf("error while parsing cpu limit '%s': %s", cpuLimit, err)
	} else if cpuLimitQuantity.MilliValue() < cpuQuantity.MilliValue() {
		return nil, fmt.Errorf("cpu limit '%s' must not be smaller than cpu '%s'", cpuLimit, cpu)
	}

	var memory KubeResourceQuantity
	if spec.Memory != nil {
		memory = spec.Memory.withSuffix()
	} else {
		memory = defaults.Memory.withSuffix()
	}
	memoryQuantity, err := resource.ParseQuantity(string(memory))
	if err != nil {
		return nil, fmt.Errorf("error while parsing memory '%s': %s", memory, err)
	}

	var memoryLimit KubeResourceQuantity
	if spec.MemoryLimit != nil {
		memoryLimit = spec.MemoryLimit.withSuffix()
	} else {
		memoryLimit = memory
	}
	memoryLimitQuantity, err := resource.ParseQuantity(string(memoryLimit))
	if err != nil {
		return nil, fmt.Errorf("error while parsing memory limit '%s': %s", memoryLimit, err)
	} else if memoryLimitQuantity.Value() < memoryQuantity.Value() {
		return nil, fmt.Errorf("memory limit '%s' must not be smaller than memory '%s'", memoryLimit, memory)
	}

	var disk KubeResourceQuantity
	if spec.Disk != nil {
		disk = spec.Disk.withSuffix()
	} else {
		disk = defaults.Disk.withSuffix()
	}
	diskQuantity, err := resource.ParseQuantity(string(disk))
	if err != nil {
		return nil, fmt.Errorf("error while parsing disk '%s': %s", disk, err)
	}

	var diskLimit KubeResourceQuantity
	if spec.DiskLimit != nil {
		diskLimit = spec.DiskLimit.withSuffix()
	} else {
		diskLimit = disk
	}
	diskLimitQuantity, err := resource.ParseQuantity(string(diskLimit))
	if err != nil {
		return nil, fmt.Errorf("error while parsing disk limit '%s': %s", diskLimit, err)
	} else if diskLimitQuantity.Value() < diskQuantity.Value() {
		return nil, fmt.Errorf("disk limit '%s' must not be smaller than disk '%s'", diskLimit, disk)
	}

	return &corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:              cpuQuantity,
			corev1.ResourceMemory:           memoryQuantity,
			corev1.ResourceEphemeralStorage: diskQuantity,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:              cpuLimitQuantity,
			corev1.ResourceMemory:           memoryLimitQuantity,
			corev1.ResourceEphemeralStorage: diskLimitQuantity,
		},
	}, nil
}

func (in *PaastaContainerSpec) DeepCopyInto(out *PaastaContainerSpec) {
	*out = *in
	if in.CPU != nil {
		in, out := &in.CPU, &out.CPU
		*out = new(KubeResourceQuantity)
		**out = **in
	}
	if in.CPULimit != nil {
		in, out := &in.CPULimit, &out.CPULimit
		*out = new(KubeResourceQuantity)
		**out = **in
	}
	if in.Memory != nil {
		in, out := &in.Memory, &out.Memory
		*out = new(KubeResourceQuantity)
		**out = **in
	}
	if in.MemoryLimit != nil {
		in, out := &in.MemoryLimit, &out.MemoryLimit
		*out = new(KubeResourceQuantity)
		**out = **in
	}
	if in.Disk != nil {
		in, out := &in.Disk, &out.Disk
		*out = new(KubeResourceQuantity)
		**out = **in
	}
	if in.DiskLimit != nil {
		in, out := &in.DiskLimit, &out.DiskLimit
		*out = new(KubeResourceQuantity)
		**out = **in
	}
}

func (in *PaastaContainerSpec) DeepCopy() *PaastaContainerSpec {
	if in == nil {
		return nil
	}
	out := new(PaastaContainerSpec)
	in.DeepCopyInto(out)
	return out
}

func newDefaultSpec() *PaastaContainerSpec {
	cpu := defaultCPU.copy()
	memory := defaultMemory.copy()
	disk := defaultDisk.copy()
	return &PaastaContainerSpec{
		CPU:         &cpu,
		CPULimit:    nil,
		Memory:      &memory,
		MemoryLimit: nil,
		Disk:        &disk,
		DiskLimit:   nil,
	}
}
