package safetemp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	d, c, err := Dir("", "test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, err := os.Stat(d); err == nil || !os.IsNotExist(err) {
		t.Fatalf("directory %q should not exist", d)
	}

	parent := filepath.Dir(d)
	fi, err := os.Stat(parent)
	if err != nil {
		t.Fatalf("parent directory error: %s", err)
	}
	if v := fi.Mode().Perm();v != 0700 {
		t.Fatalf("parent directory should be 0700: %s", v)
	}

	// Create the directory
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := os.Stat(d); err != nil {
		t.Fatalf("directory %q should exist", d)
	}

	// Close should remove it
	if err := c.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := os.Stat(d); err == nil || !os.IsNotExist(err) {
		t.Fatalf("directory %q should not exist", d)
	}
	if _, err := os.Stat(parent); err == nil || !os.IsNotExist(err) {
		t.Fatalf("directory %q should not exist", parent)
	}
}
