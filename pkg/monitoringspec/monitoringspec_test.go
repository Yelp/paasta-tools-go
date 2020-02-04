package monitoringspec

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	in := `{"team":"some-team","page":true,"tags":["some-tag"]}`
	var spec PaastaMonitoringSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	spec2 := spec.DeepCopy()
	spec.Team = "x"
	if spec2.Team == "x" {
		t.Errorf("Detected shallow copy of Team")
	}
	*spec.Page = false
	if *spec2.Page == false {
		t.Errorf("Detected shallow copy of Page")
	}
	spec.Tags = append(spec.Tags, "other-tag")
	if len(spec2.Tags) != 1 {
		t.Errorf("Detected shallow copy of Tags")
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
	team := "some-team"
	page := true
	tags := []string{"some-tag"}
	jsonTags, _ := json.Marshal(tags)
	in := fmt.Sprintf(`{"team":"%s","page":%t,"tags":%s}`,
		team,
		page,
		string(jsonTags),
	)
	var spec PaastaMonitoringSpec
	if err := json.Unmarshal([]byte(in), &spec); err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	if spec.Team != team {
		t.Errorf("%s != %s", spec.Team, team)
	}
	if *spec.Page != page {
		t.Errorf("%t != %t", *spec.Page, page)
	}
	if len(spec.Tags) != len(tags) {
		t.Errorf("%s != %s", spec.Tags, tags)
	}
}

func checkDeepCopy(t *testing.T, input string) {
	in := []byte(input)
	var spec PaastaMonitoringSpec
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
		`{}`,
	)
}

func TestOnlyTeamDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"team":"some-team"}`,
	)
}

func TestOnlyPageDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"page":true}`,
	)
}

func TestOnlyTagsDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"tags":["some-tag"]}`,
	)
}

func TestMultiDeepCopy(t *testing.T) {
	checkDeepCopy(
		t,
		`{"team":"some-team","page":true,"tags":["some-tag"]}`,
	)
}
