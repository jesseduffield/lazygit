package getter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileGetter_impl(t *testing.T) {
	var _ Getter = new(FileGetter)
}

func TestFileGetter(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)

	// With a dir that doesn't exist
	if err := g.Get(dst, testModuleURL("basic")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the destination folder is a symlink
	fi, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("destination is not a symlink")
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestFileGetter_sourceFile(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)

	// With a source URL that is a path to a file
	u := testModuleURL("basic")
	u.Path += "/main.tf"
	if err := g.Get(dst, u); err == nil {
		t.Fatal("should error")
	}
}

func TestFileGetter_sourceNoExist(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)

	// With a source URL that doesn't exist
	u := testModuleURL("basic")
	u.Path += "/main"
	if err := g.Get(dst, u); err == nil {
		t.Fatal("should error")
	}
}

func TestFileGetter_dir(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)

	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatalf("err: %s", err)
	}

	// With a dir that exists that isn't a symlink
	if err := g.Get(dst, testModuleURL("basic")); err == nil {
		t.Fatal("should error")
	}
}

func TestFileGetter_dirSymlink(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)
	dst2 := tempDir(t)

	// Make parents
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := os.MkdirAll(dst2, 0755); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Make a symlink
	if err := os.Symlink(dst2, dst); err != nil {
		t.Fatalf("err: %s", err)
	}

	// With a dir that exists that isn't a symlink
	if err := g.Get(dst, testModuleURL("basic")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestFileGetter_GetFile(t *testing.T) {
	g := new(FileGetter)
	dst := tempFile(t)

	// With a dir that doesn't exist
	if err := g.GetFile(dst, testModuleURL("basic-file/foo.txt")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the destination folder is a symlink
	fi, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatal("destination is not a symlink")
	}

	// Verify the main file exists
	assertContents(t, dst, "Hello\n")
}

func TestFileGetter_GetFile_Copy(t *testing.T) {
	g := new(FileGetter)
	g.Copy = true

	dst := tempFile(t)

	// With a dir that doesn't exist
	if err := g.GetFile(dst, testModuleURL("basic-file/foo.txt")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the destination folder is a symlink
	fi, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal("destination is a symlink")
	}

	// Verify the main file exists
	assertContents(t, dst, "Hello\n")
}

// https://github.com/hashicorp/terraform/issues/8418
func TestFileGetter_percent2F(t *testing.T) {
	g := new(FileGetter)
	dst := tempDir(t)

	// With a dir that doesn't exist
	if err := g.Get(dst, testModuleURL("basic%2Ftest")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestFileGetter_ClientMode_notexist(t *testing.T) {
	g := new(FileGetter)

	u := testURL("nonexistent")
	if _, err := g.ClientMode(u); err == nil {
		t.Fatal("expect source file error")
	}
}

func TestFileGetter_ClientMode_file(t *testing.T) {
	g := new(FileGetter)

	// Check the client mode when pointed at a file.
	mode, err := g.ClientMode(testModuleURL("basic-file/foo.txt"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if mode != ClientModeFile {
		t.Fatal("expect ClientModeFile")
	}
}

func TestFileGetter_ClientMode_dir(t *testing.T) {
	g := new(FileGetter)

	// Check the client mode when pointed at a directory.
	mode, err := g.ClientMode(testModuleURL("basic"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if mode != ClientModeDir {
		t.Fatal("expect ClientModeDir")
	}
}
