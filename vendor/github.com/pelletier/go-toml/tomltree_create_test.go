package toml

import (
	"strconv"
	"testing"
	"time"
)

type customString string

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}

func validate(t *testing.T, path string, object interface{}) {
	switch o := object.(type) {
	case *Tree:
		for key, tree := range o.values {
			validate(t, path+"."+key, tree)
		}
	case []*Tree:
		for index, tree := range o {
			validate(t, path+"."+strconv.Itoa(index), tree)
		}
	case *tomlValue:
		switch o.value.(type) {
		case int64, uint64, bool, string, float64, time.Time,
			[]int64, []uint64, []bool, []string, []float64, []time.Time:
		default:
			t.Fatalf("tomlValue at key %s containing incorrect type %T", path, o.value)
		}
	default:
		t.Fatalf("value at key %s is of incorrect type %T", path, object)
	}
	t.Logf("validation ok %s as %T", path, object)
}

func validateTree(t *testing.T, tree *Tree) {
	validate(t, "", tree)
}

func TestTreeCreateToTree(t *testing.T) {
	data := map[string]interface{}{
		"a_string": "bar",
		"an_int":   42,
		"time":     time.Now(),
		"int8":     int8(2),
		"int16":    int16(2),
		"int32":    int32(2),
		"uint8":    uint8(2),
		"uint16":   uint16(2),
		"uint32":   uint32(2),
		"float32":  float32(2),
		"a_bool":   false,
		"stringer": stringer{},
		"nested": map[string]interface{}{
			"foo": "bar",
		},
		"array":                 []string{"a", "b", "c"},
		"array_uint":            []uint{uint(1), uint(2)},
		"array_table":           []map[string]interface{}{{"sub_map": 52}},
		"array_times":           []time.Time{time.Now(), time.Now()},
		"map_times":             map[string]time.Time{"now": time.Now()},
		"custom_string_map_key": map[customString]interface{}{customString("custom"): "custom"},
	}
	tree, err := TreeFromMap(data)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	validateTree(t, tree)
}

func TestTreeCreateToTreeInvalidLeafType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": t})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidMapKeyType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": map[int]interface{}{2: 1}})
	expected := "map key needs to be a string, not int (int)"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidArrayMemberType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": []*testing.T{t}})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestTreeCreateToTreeInvalidTableGroupType(t *testing.T) {
	_, err := TreeFromMap(map[string]interface{}{"foo": []map[string]interface{}{{"hello": t}}})
	expected := "cannot convert type *testing.T to Tree"
	if err.Error() != expected {
		t.Fatalf("expected error %s, got %s", expected, err.Error())
	}
}

func TestRoundTripArrayOfTables(t *testing.T) {
	orig := "\n[[stuff]]\n  name = \"foo\"\n  things = [\"a\",\"b\"]\n"
	tree, err := Load(orig)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	m := tree.ToMap()

	tree, err = TreeFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := orig
	got := tree.String()

	if got != want {
		t.Errorf("want:\n%s\ngot:\n%s", want, got)
	}
}
