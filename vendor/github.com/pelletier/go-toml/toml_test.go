// Testing support for go-toml

package toml

import (
	"testing"
)

func TestTomlHas(t *testing.T) {
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if !tree.Has("test.key") {
		t.Errorf("Has - expected test.key to exists")
	}

	if tree.Has("") {
		t.Errorf("Should return false if the key is not provided")
	}
}

func TestTomlGet(t *testing.T) {
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if tree.Get("") != tree {
		t.Errorf("Get should return the tree itself when given an empty path")
	}

	if tree.Get("test.key") != "value" {
		t.Errorf("Get should return the value")
	}
	if tree.Get(`\`) != nil {
		t.Errorf("should return nil when the key is malformed")
	}
}

func TestTomlGetDefault(t *testing.T) {
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if tree.GetDefault("", "hello") != tree {
		t.Error("GetDefault should return the tree itself when given an empty path")
	}

	if tree.GetDefault("test.key", "hello") != "value" {
		t.Error("Get should return the value")
	}

	if tree.GetDefault("whatever", "hello") != "hello" {
		t.Error("GetDefault should return the default value if the key does not exist")
	}
}

func TestTomlHasPath(t *testing.T) {
	tree, _ := Load(`
		[test]
		key = "value"
	`)

	if !tree.HasPath([]string{"test", "key"}) {
		t.Errorf("HasPath - expected test.key to exists")
	}
}

func TestTomlDelete(t *testing.T) {
	tree, _ := Load(`
        key = "value"
    `)
	err := tree.Delete("key")
	if err != nil {
		t.Errorf("Delete - unexpected error while deleting key: %s", err.Error())
	}

	if tree.Get("key") != nil {
		t.Errorf("Delete should have removed key but did not.")
	}

}

func TestTomlDeleteUnparsableKey(t *testing.T) {
	tree, _ := Load(`
        key = "value"
    `)
	err := tree.Delete(".")
	if err == nil {
		t.Errorf("Delete should error")
	}
}

func TestTomlDeleteNestedKey(t *testing.T) {
	tree, _ := Load(`
		[foo]
        [foo.bar]
        key = "value"
    `)
	err := tree.Delete("foo.bar.key")
	if err != nil {
		t.Errorf("Error while deleting nested key: %s", err.Error())
	}

	if tree.Get("key") != nil {
		t.Errorf("Delete should have removed nested key but did not.")
	}

}

func TestTomlDeleteNonexistentNestedKey(t *testing.T) {
	tree, _ := Load(`
		[foo]
        [foo.bar]
        key = "value"
    `)
	err := tree.Delete("foo.not.there.key")
	if err == nil {
		t.Errorf("Delete should have thrown an error trying to delete key in nonexistent tree")
	}
}

func TestTomlGetPath(t *testing.T) {
	node := newTree()
	//TODO: set other node data

	for idx, item := range []struct {
		Path     []string
		Expected *Tree
	}{
		{ // empty path test
			[]string{},
			node,
		},
	} {
		result := node.GetPath(item.Path)
		if result != item.Expected {
			t.Errorf("GetPath[%d] %v - expected %v, got %v instead.", idx, item.Path, item.Expected, result)
		}
	}

	tree, _ := Load("[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6")
	if tree.GetPath([]string{"whatever"}) != nil {
		t.Error("GetPath should return nil when the key does not exist")
	}
}

func TestTomlFromMap(t *testing.T) {
	simpleMap := map[string]interface{}{"hello": 42}
	tree, err := TreeFromMap(simpleMap)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if tree.Get("hello") != int64(42) {
		t.Fatal("hello should be 42, not", tree.Get("hello"))
	}
}

func TestLoadBytesBOM(t *testing.T) {
	payloads := [][]byte{
		[]byte("\xFE\xFFhello=1"),
		[]byte("\xFF\xFEhello=1"),
		[]byte("\xEF\xBB\xBFhello=1"),
		[]byte("\x00\x00\xFE\xFFhello=1"),
		[]byte("\xFF\xFE\x00\x00hello=1"),
	}
	for _, data := range payloads {
		tree, err := LoadBytes(data)
		if err != nil {
			t.Fatal("unexpected error:", err, "for:", data)
		}
		v := tree.Get("hello")
		if v != int64(1) {
			t.Fatal("hello should be 1, not", v)
		}
	}
}
