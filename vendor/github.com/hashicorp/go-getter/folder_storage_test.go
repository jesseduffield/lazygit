package getter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFolderStorage_impl(t *testing.T) {
	var _ Storage = new(FolderStorage)
}

func TestFolderStorage(t *testing.T) {
	s := &FolderStorage{StorageDir: tempDir(t)}

	module := testModule("basic")

	// A module shouldn't exist at first...
	_, ok, err := s.Dir(module)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ok {
		t.Fatal("should not exist")
	}

	key := "foo"

	// We can get it
	err = s.Get(key, module, false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Now the module exists
	dir, ok, err := s.Dir(key)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !ok {
		t.Fatal("should exist")
	}

	mainPath := filepath.Join(dir, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}
