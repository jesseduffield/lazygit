package getter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSourceDirSubdir(t *testing.T) {
	cases := []struct {
		Input    string
		Dir, Sub string
	}{
		{
			"hashicorp.com",
			"hashicorp.com", "",
		},
		{
			"hashicorp.com//foo",
			"hashicorp.com", "foo",
		},
		{
			"hashicorp.com//foo?bar=baz",
			"hashicorp.com?bar=baz", "foo",
		},
		{
			"https://hashicorp.com/path//*?archive=foo",
			"https://hashicorp.com/path?archive=foo", "*",
		},
		{
			"file://foo//bar",
			"file://foo", "bar",
		},
	}

	for i, tc := range cases {
		adir, asub := SourceDirSubdir(tc.Input)
		if adir != tc.Dir {
			t.Fatalf("%d: bad dir: %#v", i, adir)
		}
		if asub != tc.Sub {
			t.Fatalf("%d: bad sub: %#v", i, asub)
		}
	}
}

func TestSourceSubdirGlob(t *testing.T) {
	td, err := ioutil.TempDir("", "subdir-glob")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	if err := os.Mkdir(filepath.Join(td, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(filepath.Join(td, "subdir/one"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(filepath.Join(td, "subdir/two"), 0755); err != nil {
		t.Fatal(err)
	}

	subdir := filepath.Join(td, "subdir")

	// match the exact directory
	res, err := SubdirGlob(td, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	if res != subdir {
		t.Fatalf(`expected "subdir", got: %q`, subdir)
	}

	// single match from a wildcard
	res, err = SubdirGlob(td, "*")
	if err != nil {
		t.Fatal(err)
	}
	if res != subdir {
		t.Fatalf(`expected "subdir", got: %q`, subdir)
	}

	// multiple matches
	res, err = SubdirGlob(td, "subdir/*")
	if err == nil {
		t.Fatalf("expected multiple matches, got %q", res)
	}

	// non-existent
	res, err = SubdirGlob(td, "foo")
	if err == nil {
		t.Fatalf("expected no matches, got %q", res)
	}
}
