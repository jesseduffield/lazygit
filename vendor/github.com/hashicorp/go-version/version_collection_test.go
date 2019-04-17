package version

import (
	"reflect"
	"sort"
	"testing"
)

func TestCollection(t *testing.T) {
	versionsRaw := []string{
		"1.1.1",
		"1.0",
		"1.2",
		"2",
		"0.7.1",
	}

	versions := make([]*Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, err := NewVersion(raw)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		versions[i] = v
	}

	sort.Sort(Collection(versions))

	actual := make([]string, len(versions))
	for i, v := range versions {
		actual[i] = v.String()
	}

	expected := []string{
		"0.7.1",
		"1.0.0",
		"1.1.1",
		"1.2.0",
		"2.0.0",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
