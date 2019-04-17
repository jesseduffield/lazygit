package util_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
)

func TestTempFile(t *testing.T) {
	fs := memfs.New()

	dir, err := util.TempDir(fs, "", "util_test")
	if err != nil {
		t.Fatal(err)
	}
	defer util.RemoveAll(fs, dir)

	f, err := util.TempFile(fs, dir, "foo")
	if f == nil || err != nil {
		t.Errorf("TempFile(%q, `foo`) = %v, %v", dir, f, err)
	}
}

func TestTempDir(t *testing.T) {
	fs := memfs.New()

	dir := os.TempDir()
	name, err := util.TempDir(fs, dir, "util_test")
	if name == "" || err != nil {
		t.Errorf("TempDir(dir, `util_test`) = %v, %v", name, err)
	}
	if name != "" {
		util.RemoveAll(fs, name)
		re := regexp.MustCompile("^" + regexp.QuoteMeta(filepath.Join(dir, "util_test")) + "[0-9]+$")
		if !re.MatchString(name) {
			t.Errorf("TempDir(`"+dir+"`, `util_test`) created bad name %s", name)
		}
	}
}
