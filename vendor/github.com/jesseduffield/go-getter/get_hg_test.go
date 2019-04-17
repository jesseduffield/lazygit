package getter

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testHasHg bool

func init() {
	if _, err := exec.LookPath("hg"); err == nil {
		testHasHg = true
	}
}

func TestHgGetter_impl(t *testing.T) {
	var _ Getter = new(HgGetter)
}

func TestHgGetter(t *testing.T) {
	if !testHasHg {
		t.Log("hg not found, skipping")
		t.Skip()
	}

	g := new(HgGetter)
	dst := tempDir(t)

	// With a dir that doesn't exist
	if err := g.Get(dst, testModuleURL("basic-hg")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHgGetter_branch(t *testing.T) {
	if !testHasHg {
		t.Log("hg not found, skipping")
		t.Skip()
	}

	g := new(HgGetter)
	dst := tempDir(t)

	url := testModuleURL("basic-hg")
	q := url.Query()
	q.Add("rev", "test-branch")
	url.RawQuery = q.Encode()

	if err := g.Get(dst, url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main_branch.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get again should work
	if err := g.Get(dst, url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath = filepath.Join(dst, "main_branch.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHgGetter_GetFile(t *testing.T) {
	if !testHasHg {
		t.Log("hg not found, skipping")
		t.Skip()
	}

	g := new(HgGetter)
	dst := tempFile(t)

	// Download
	if err := g.GetFile(dst, testModuleURL("basic-hg/foo.txt")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("err: %s", err)
	}
	assertContents(t, dst, "Hello\n")
}
