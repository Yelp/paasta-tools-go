package containerspec

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	in := `{"cpus":"0.25","mem":"2048","disk":"10240","disk_limit":"102400"}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	spec2 := spec.DeepCopy()
	*spec.CPU = "x"
	if *spec2.CPU == "x" {
		t.Errorf("Detected shallow copy of CPU")
	}
	*spec.Memory = "x"
	if *spec2.Memory == "x" {
		t.Errorf("Detected shallow copy of Memory")
	}
	*spec.Disk = "x"
	if *spec2.Disk == "x" {
		t.Errorf("Detected shallow copy of Disk")
	}
	*spec.DiskLimit = "x"
	if *spec2.DiskLimit == "x" {
		t.Errorf("Detected shallow copy of DiskLimit")
	}
	out, err := json.Marshal(spec2);
	if err != nil {
		t.Errorf("Failed to marshal: %s", err)
	}
	if string(out) != in {
		t.Errorf("%s != %s", out, in)
	}
}

func TestUnmarshal(t *testing.T) {
	cpu := KubeResourceQuantity("0.2")
	mem := KubeResourceQuantity("1024")
	disk := KubeResourceQuantity("2048")
	in := fmt.Sprintf(`{"cpus":"%s","mem":"%s","disk":"%s"}`,
		string(cpu),
		string(mem),
		string(disk),
	)
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	if *spec.CPU != cpu {
		t.Errorf("%s != %s", *spec.CPU, cpu)
	}
	if *spec.Memory != mem {
		t.Errorf("%s != %s", *spec.Memory, mem)
	}
	if *spec.Disk != disk {
		t.Errorf("%s != %s", *spec.Disk, disk)
	}
}

func TestUnmarshalNull(t *testing.T) {
	in := `{"cpus":null,"mem":null,"disk":null}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	if spec.CPU != nil {
		t.Errorf("%s != nil", *spec.CPU)
	}
	if spec.Memory != nil {
		t.Errorf("%s != nil", *spec.Memory)
	}
	if spec.Disk != nil {
		t.Errorf("%s != nil", *spec.Disk)
	}
}

func checkDeepCopy(t *testing.T, input string) {
	in := []byte(input)
	var spec PaastaContainerSpec
	if err := json.Unmarshal(in, &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	spec2 := spec.DeepCopy()
	out, err := json.Marshal(spec2);
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

func TestOnlyMemDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":null,"mem":"2048","disk":null}`,
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
		`{"cpus":"0.25","mem":"2048","disk":"10240","disk_limit":"102400"}`,
	)
}

func TestJSONRoundTrip(t *testing.T) {
	in := `{"cpus":"0.2","mem":"1024","disk":"4096","disk_limit":"4Gi"}`
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	out, err := json.Marshal(spec);
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

func TestOnlyMemResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"mem":"1024"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1Gi"}}`,
	)
}

func TestOnlyDiskResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk":"2000"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512Mi"},"requests":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512Mi"}}`,
	)
}

func TestOnlyDiskResourcesBin(t *testing.T) {
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
		`{"cpus":"0.2","mem":"1024","disk":"10Gi","disk_limit":"2048Gi"}`,
		`{"limits":{"cpu":"200m","ephemeral-storage":"2Ti","memory":"1Gi"},"requests":{"cpu":"200m","ephemeral-storage":"10Gi","memory":"1Gi"}}`,
	)
}

func checkResourcesError(t *testing.T, input string) error {
	in := []byte(input)
	var spec PaastaContainerSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	_, err := spec.GetContainerResources()
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

func TestDiskLimitSameAsDisk(t *testing.T) {
	err := checkResourcesError(t, `{"disk":"2Gi","disk_limit":"2048Mi"}`)
	if err != nil {
		t.Errorf("Detection of a too small disk limit wrongly triggered")
	}
}