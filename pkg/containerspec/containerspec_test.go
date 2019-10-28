package containerspec

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
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

func TestAllDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"cpus":"0.25","mem":"2048","disk":"10240"}`,
	)
}

func TestJSONRoundTrip(t *testing.T) {
	in := `{"cpus":"0.2","mem":"1024","disk":"4096"}`
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
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512M"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"512M"}}`,
	)
}

func TestOnlyCPUResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.5"}`,
		`{"limits":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512M"},"requests":{"cpu":"500m","ephemeral-storage":"1Gi","memory":"512M"}}`,
	)
}

func TestOnlyMemResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"mem":"1024"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1024M"},"requests":{"cpu":"100m","ephemeral-storage":"1Gi","memory":"1024M"}}`,
	)
}

func TestOnlyDiskResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"disk":"2000"}`,
		`{"limits":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512M"},"requests":{"cpu":"100m","ephemeral-storage":"2000Mi","memory":"512M"}}`,
	)
}

func TestBothMemCPUResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.2","mem":"1024"}`,
		`{"limits":{"cpu":"200m","ephemeral-storage":"1Gi","memory":"1024M"},"requests":{"cpu":"200m","ephemeral-storage":"1Gi","memory":"1024M"}}`,
	)
}

func TestAllResources(t *testing.T) {
	checkEqualResources(
		t,
		`{"cpus":"0.2","mem":"1024","disk":"10Gi"}`,
		`{"limits":{"cpu":"200m","ephemeral-storage":"10Gi","memory":"1024M"},"requests":{"cpu":"200m","ephemeral-storage":"10Gi","memory":"1024M"}}`,
	)
}
