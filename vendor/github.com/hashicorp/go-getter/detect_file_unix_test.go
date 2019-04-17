// +build test unix

package getter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// If a relative symlink is passed in as the pwd to Detect, the resulting URL
// can have an invalid path.
func TestFileDetector_relativeSymlink(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "go-getter")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	// We may have a symlinked tmp dir,
	// e.g. OSX uses /var -> /private/var
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(tmpDir, "realPWD"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	subdir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subdir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	prevDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevDir)

	err = os.Chdir(subdir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink("../realPWD", "linkedPWD")
	if err != nil {
		t.Fatal(err)
	}

	// if detech doesn't fully resolve the pwd symlink, the output will be the
	// invalid path: "file:///../modules/foo"
	f := new(FileDetector)
	out, ok, err := f.Detect("../modules/foo", "./linkedPWD")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !ok {
		t.Fatal("not ok")
	}
	if out != "file://"+filepath.Join(tmpDir, "modules/foo") {
		t.Logf("expected: %v", "file://"+filepath.Join(tmpDir, "modules/foo"))
		t.Fatalf("bad:      %v", out)
	}
}
