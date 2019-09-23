package containerspec

import (
    "encoding/json"
    "fmt"
    "testing"
)

func TestUnmarshal(t *testing.T) {
    cpu := KubeResourceQuantity("0.2")
    mem := KubeResourceQuantity("1024")
    in := fmt.Sprintf(`{"cpus":"%s","mem":"%s"}`, string(cpu), string(mem))
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
}

func TestJSONRoundTrip(t *testing.T) {
    in := `{"cpus":"0.2","mem":"1024"}`
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

func checkEqual(t *testing.T, input string, exp string) {
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

func TestEmpty(t *testing.T) {
    checkEqual(
        t,
        "{}",
        `{"limits":{"cpu":"100m","memory":"512M"},"requests":{"cpu":"100m","memory":"512M"}}`,
    )
}

func TestOnlyCPU(t *testing.T) {
    checkEqual(
        t,
        `{"cpus":"0.5"}`,
        `{"limits":{"cpu":"500m","memory":"512M"},"requests":{"cpu":"500m","memory":"512M"}}`,
    )
}

func TestOnlyMem(t *testing.T) {
    checkEqual(
        t,
        `{"mem":"1024"}`,
        `{"limits":{"cpu":"100m","memory":"1024M"},"requests":{"cpu":"100m","memory":"1024M"}}`,
    )
}

func TestBothMemCPU(t *testing.T) {
    checkEqual(
        t,
        `{"cpus":"0.2","mem":"1024"}`,
        `{"limits":{"cpu":"200m","memory":"1024M"},"requests":{"cpu":"200m","memory":"1024M"}}`,
    )
}
