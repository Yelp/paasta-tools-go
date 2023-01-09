package containerspec

import (
	"encoding/json"
	"fmt"
	"testing"

    "github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDeepCopy(t *testing.T) {
	in := `{"cpus":"0.25","cpus_limit":"0.25","mem":"2048","mem_limit":"2048","disk":"10240","disk_limit":"102400"}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	spec2 := spec.DeepCopy()
	*spec.CPU = "x"
	if *spec2.CPU == "x" {
		t.Errorf("Detected shallow copy of CPU")
	}
	*spec.CPULimit = "x"
	if *spec2.CPULimit == "x" {
		t.Errorf("Detected shallow copy of CPULimit")
	}
	*spec.Memory = "x"
	if *spec2.Memory == "x" {
		t.Errorf("Detected shallow copy of Memory")
	}
	*spec.MemoryLimit = "x"
	if *spec2.MemoryLimit == "x" {
		t.Errorf("Detected shallow copy of MemoryLimit")
	}
	*spec.Disk = "x"
	if *spec2.Disk == "x" {
		t.Errorf("Detected shallow copy of Disk")
	}
	*spec.DiskLimit = "x"
	if *spec2.DiskLimit == "x" {
		t.Errorf("Detected shallow copy of DiskLimit")
	}
	out, err := json.Marshal(spec2)
	if err != nil {
		t.Errorf("Failed to marshal: %s", err)
	}
	if string(out) != in {
		t.Errorf("%s != %s", out, in)
	}
}

func TestUnmarshal(t *testing.T) {
	cpu := KubeResourceQuantity("0.2")
	cpuLimit := KubeResourceQuantity("0.25")
	mem := KubeResourceQuantity("1024")
	memLimit := KubeResourceQuantity("3072")
	disk := KubeResourceQuantity("2048")
	diskLimit := KubeResourceQuantity("10240")
	in := fmt.Sprintf(`{"cpus":"%s","cpus_limit":"%s","mem":"%s","mem_limit":"%s","disk":"%s","disk_limit":"%s"}`,
		string(cpu),
		string(cpuLimit),
		string(mem),
		string(memLimit),
		string(disk),
		string(diskLimit),
	)
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	if *spec.CPU != cpu {
		t.Errorf("%s != %s", *spec.CPU, cpu)
	}
	if *spec.CPULimit != cpuLimit {
		t.Errorf("%s != %s", *spec.CPULimit, cpuLimit)
	}
	if *spec.Memory != mem {
		t.Errorf("%s != %s", *spec.Memory, mem)
	}
	if *spec.MemoryLimit != memLimit {
		t.Errorf("%s != %s", *spec.MemoryLimit, memLimit)
	}
	if *spec.Disk != disk {
		t.Errorf("%s != %s", *spec.Disk, disk)
	}
	if *spec.DiskLimit != diskLimit {
		t.Errorf("%s != %s", *spec.DiskLimit, diskLimit)
	}
}

func TestUnmarshalNull(t *testing.T) {
	in := `{"cpus":null,"cpus_limit":null,"mem":null,"mem_limit":null,"disk":null,"disk_limit":null}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	if spec.CPU != nil {
		t.Errorf("%s != nil", *spec.CPU)
	}
	if spec.CPULimit != nil {
		t.Errorf("%s != nil", *spec.CPULimit)
	}
	if spec.Memory != nil {
		t.Errorf("%s != nil", *spec.Memory)
	}
	if spec.MemoryLimit != nil {
		t.Errorf("%s != nil", *spec.MemoryLimit)
	}
	if spec.Disk != nil {
		t.Errorf("%s != nil", *spec.Disk)
	}
	if spec.DiskLimit != nil {
		t.Errorf("%s != nil", *spec.DiskLimit)
	}
}

func checkDeepCopy(t *testing.T, input string) {
	in := []byte(input)
	var spec PaastaContainerSpec
	if err := json.Unmarshal(in, &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	spec2 := spec.DeepCopy()
	out, err := json.Marshal(spec2)
	if err != nil {
		t.Errorf("Failed to marshal: %s", err)
	}
	if string(out) != input {
		t.Errorf("%s != %s", out, in)
	}
}

func TestEmptyDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":null,"disk":null}`,
	)
}

func TestOnlyCPUDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":"0.5","mem":null,"disk":null}`,
	)
}

func TestOnlyCPULimitDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"cpus_limit":"0.5","mem":null,"disk":null}`,
	)
}

func TestOnlyMemDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":"2048","disk":null}`,
	)
}

func TestOnlyMemLimitDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":null,"mem_limit":"2048","disk":null}`,
	)
}

func TestOnlyDiskDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":null,"disk":"10240"}`,
	)
}

func TestOnlyDiskLimitDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":null,"disk":null,"disk_limit":"10240"}`,
	)
}

func TestAllDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":"0.25","cpus_limit":"0.5","mem":"2048","mem_limit":"3072","disk":"10240","disk_limit":"102400"}`,
	)
}

func TestJSONRoundTrip(t *testing.T) {
	in := `{"cpus":"0.2","cpus_limit":"0.5","mem":"1024","mem_limit":"1.2Gi","disk":"4096","disk_limit":"4Gi"}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	out, err := json.Marshal(spec)
	if err != nil {
		t.Errorf("Failed to marshal: %s", err)
	}
	if string(out) != in {
		t.Errorf("%s != %s", out, in)
	}
}

func checkEqualResources(t *testing.T, input string, exp string) {
	in := []byte(input)
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	res, err := spec.GetContainerResources()
	if err != nil {
		t.Errorf("Failed to build resource requirements: %s", err)
	}
	out, err := json.Marshal(res)
	if err != nil {
		t.Errorf("Failed to marshal resource requirements: %s", err)
	}
	if string(out) != exp {
		t.Errorf("%s != %s", out, exp)
	}
}

func TestEmptyResources(t *testing.T) {
	checkEqualResources(
		t,
		"{}",
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestOnlyCPUResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.5"}`,
		`{"limits":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512Mi"},"requests":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestOnlyCPULimitResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus_limit":"0.5"}`,
		`{"limits":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestBothCPUResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.4","cpus_limit":"0.5"}`,
		`{"limits":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512Mi"},"requests":{"cpu":"400m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestOnlyMemResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"mem":"1024"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"}}`,
	)
}

func TestOnlyMemLimitResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"mem_limit":"1024"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestBothMemResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"mem":"768","mem_limit":"1024"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"768Mi"}}`,
	)
}

func TestOnlyDiskResourcesBin(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk":"2000"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512Mi"}}`,
	)
}

func TestOnlyDiskResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk":"2048"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"2Gi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"2Gi","memory":"512Mi"}}`,
	)
}

func TestBothDiskLimitDiskResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk":"2000Mi","disk_limit":"20Gi"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"20Gi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512Mi"}}`,
	)
}

func TestLimitDiskResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk_limit":"20480"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"20Gi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512Mi"}}`,
	)
}

func TestBothMemCPUResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.2","mem":"1024"}`,
		`{"limits":{"cpu":"200m","ephemeral-storage":"1Gi","memory":"1Gi"},"requests":{"cpu":"200m","ephemeral-storage":"1Gi","memory":"1Gi"}}`,
	)
}

func TestAllResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.2","cpus_limit":"0.25","mem":"1000","mem_limit":"1024","disk":"10Gi","disk_limit":"2048Gi"}`,
		`{"limits":{"cpu":"250m","ephemeral-storage":"2Ti","memory":"1Gi"},"requests":{"cpu":"200m","ephemeral-storage":"10Gi","memory":"1000Mi"}}`,
	)
}

func checkResources(t *testing.T, input string, req *corev1.ResourceRequirements) error {
	in := []byte(input)
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	r, err := spec.GetContainerResources()
	if err == nil && req != nil {
		*req = *r
	}
	return err
}

func checkResourcesError(t *testing.T, input string) error {
	req := corev1.ResourceRequirements{}
	err := checkResources(t, input, &req)
	if err == nil {
		t.Logf("Got: %v", req)
		return err
	}
	return err
}

func TestTooSmallDiskLimit(t *testing.T) {
	err := checkResourcesError(t, `{"disk":"201","disk_limit":"200"}`)
	if err == nil {
		t.Errorf("Detection of a too small disk limit has failed")
	}
}

func TestTooSmallDiskLimitDefaultMiSuffix(t *testing.T) {
	err := checkResourcesError(t, `{"disk":"200","disk_limit":"200M"}`)
	if err == nil {
		t.Errorf("Detection of a too small disk limit has failed")
	}
}

func TestTooSmallDiskLimitMixedSuffixes(t *testing.T) {
	err := checkResourcesError(t, `{"disk":"2Gi","disk_limit":"2048M"}`)
	if err == nil {
		t.Errorf("Detection of a too small disk limit has failed")
	}
}

func TestTooSmallMemoryLimit(t *testing.T) {
	err := checkResourcesError(t, `{"mem":"201","mem_limit":"200"}`)
	if err == nil {
		t.Errorf("Detection of a too small memory limit has failed")
	}
}

func TestTooCPULimit(t *testing.T) {
	err := checkResourcesError(t, `{"cpus":"505m","cpus_limit":"0.5"}`)
	if err == nil {
		t.Errorf("Detection of a too small cpu limit has failed")
	}
}

func checkResourcesEqual(t *testing.T, input string, wanted string, requestFn func(req *corev1.ResourceRequirements) *resource.Quantity, limitFn func(req *corev1.ResourceRequirements) *resource.Quantity) {
	req := corev1.ResourceRequirements{}
	err := checkResources(t, input, &req)
	if err != nil {
		t.Errorf("Failed to generate resource requirements from %s: %s", input, err)
	}
	expected, err := resource.ParseQuantity(wanted)
	if err != nil {
		t.Errorf("Failed to parse from %s: %s", wanted, err)
	}
	request := requestFn(&req)
	limit := limitFn(&req)
	if request.MilliValue() != expected.MilliValue() || request.MilliValue() != limit.MilliValue() {
		t.Errorf("Detection of explicitiy equal requirements has failed: %d != %d != %d", expected.MilliValue(), request.MilliValue(), limit.MilliValue())
	}
}

func TestEqualDiskResources(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"disk":"8192Mi","disk_limit":"8Gi"}`,
		"8.0Gi",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.StorageEphemeral()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.StorageEphemeral()
		},
	)
}

func TestImpliedEqualDiskLimit(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"disk":"8Gi"}`,
		"8.0Gi",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.StorageEphemeral()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.StorageEphemeral()
		},
	)
}

func TestEqualMemoryResources(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"mem":"3072Mi",",mem_limit":"3Gi"}`,
		"3.0Gi",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.Memory()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.Memory()
		},
	)
}

func TestImpliedEqualMemoryLimit(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"mem":"3072Mi"}`,
		"3.0Gi",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.Memory()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.Memory()
		},
	)
}

func TestEqualCPUResources(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"cpus":"0.3","cpus_limit":"300m"}`,
		"0.300",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.Cpu()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.Cpu()
		},
	)
}

func TestImpliedEqualCPULimit(t *testing.T) {
	checkResourcesEqual(
		t,
		`{"cpus":"0.3"}`,
		"0.300",
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Requests.Cpu()
		},
		func(req *corev1.ResourceRequirements) *resource.Quantity {
			return req.Limits.Cpu()
		},
	)
}


func TestCmp_QuantityComparingToIsEqual_ReturnZero(t *testing.T) {
	kubeResourceQuantity := KubeResourceQuantity("10")
	kubeResourceQuantityToCompare := KubeResourceQuantity("10")

	result, err := kubeResourceQuantity.Cmp(kubeResourceQuantityToCompare)
	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	}

	expected := 0

	if !assert.Equal(t, expected, result) {
		t.Fatalf("Not equal: Expected: %d, Actual: %d", expected, result)
	}
}

func TestCmp_QuantityComparingToIsGreater_ReturnNegativeOne(t *testing.T) {
	kubeResourceQuantity := KubeResourceQuantity("10")
	kubeResourceQuantityToCompare := KubeResourceQuantity("20")

	result, err := kubeResourceQuantity.Cmp(kubeResourceQuantityToCompare)
	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	}

	expected := -1

	if !assert.Equal(t, expected, result) {
		t.Fatalf("Not equal: Expected: %d, Actual: %d", expected, result)
	}
}

func TestCmp_QuantityComparingToIsGreater_ReturnOne(t *testing.T) {
	kubeResourceQuantity := KubeResourceQuantity("10")
	kubeResourceQuantityToCompare := KubeResourceQuantity("5")

	result, err := kubeResourceQuantity.Cmp(kubeResourceQuantityToCompare)
	if err != nil {
		t.Errorf("Unexpected error, %s", err.Error())
	}

	expected := 1

	if !assert.Equal(t, expected, result) {
		t.Fatalf("Not equal: Expected: %d, Actual: %d", expected, result)
	}
}

