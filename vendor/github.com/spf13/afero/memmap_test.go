package afero

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestNormalizePath(t *testing.T) {
	type test struct {
		input    string
		expected string
	}

	data := []test{
		{".", FilePathSeparator},
		{"./", FilePathSeparator},
		{"..", FilePathSeparator},
		{"../", FilePathSeparator},
		{"./..", FilePathSeparator},
		{"./../", FilePathSeparator},
	}

	for i, d := range data {
		cpath := normalizePath(d.input)
		if d.expected != cpath {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, cpath)
		}
	}
}

func TestPathErrors(t *testing.T) {
	path := filepath.Join(".", "some", "path")
	path2 := filepath.Join(".", "different", "path")
	fs := NewMemMapFs()
	perm := os.FileMode(0755)

	// relevant functions:
	// func (m *MemMapFs) Chmod(name string, mode os.FileMode) error
	// func (m *MemMapFs) Chtimes(name string, atime time.Time, mtime time.Time) error
	// func (m *MemMapFs) Create(name string) (File, error)
	// func (m *MemMapFs) Mkdir(name string, perm os.FileMode) error
	// func (m *MemMapFs) MkdirAll(path string, perm os.FileMode) error
	// func (m *MemMapFs) Open(name string) (File, error)
	// func (m *MemMapFs) OpenFile(name string, flag int, perm os.FileMode) (File, error)
	// func (m *MemMapFs) Remove(name string) error
	// func (m *MemMapFs) Rename(oldname, newname string) error
	// func (m *MemMapFs) Stat(name string) (os.FileInfo, error)

	err := fs.Chmod(path, perm)
	checkPathError(t, err, "Chmod")

	err = fs.Chtimes(path, time.Now(), time.Now())
	checkPathError(t, err, "Chtimes")

	// fs.Create doesn't return an error

	err = fs.Mkdir(path2, perm)
	if err != nil {
		t.Error(err)
	}
	err = fs.Mkdir(path2, perm)
	checkPathError(t, err, "Mkdir")

	err = fs.MkdirAll(path2, perm)
	if err != nil {
		t.Error("MkdirAll:", err)
	}

	_, err = fs.Open(path)
	checkPathError(t, err, "Open")

	_, err = fs.OpenFile(path, os.O_RDWR, perm)
	checkPathError(t, err, "OpenFile")

	err = fs.Remove(path)
	checkPathError(t, err, "Remove")

	err = fs.RemoveAll(path)
	if err != nil {
		t.Error("RemoveAll:", err)
	}

	err = fs.Rename(path, path2)
	checkPathError(t, err, "Rename")

	_, err = fs.Stat(path)
	checkPathError(t, err, "Stat")
}

func checkPathError(t *testing.T, err error, op string) {
	pathErr, ok := err.(*os.PathError)
	if !ok {
		t.Error(op+":", err, "is not a os.PathError")
		return
	}
	_, ok = pathErr.Err.(*os.PathError)
	if ok {
		t.Error(op+":", err, "contains another os.PathError")
	}
}

// Ensure Permissions are set on OpenFile/Mkdir/MkdirAll
func TestPermSet(t *testing.T) {
	const fileName = "/myFileTest"
	const dirPath = "/myDirTest"
	const dirPathAll = "/my/path/to/dir"

	const fileMode = os.FileMode(0765)
	// directories will also have the directory bit set
	const dirMode = fileMode | os.ModeDir

	fs := NewMemMapFs()

	// Test Openfile
	f, err := fs.OpenFile(fileName, os.O_CREATE, fileMode)
	if err != nil {
		t.Errorf("OpenFile Create failed: %s", err)
		return
	}
	f.Close()

	s, err := fs.Stat(fileName)
	if err != nil {
		t.Errorf("Stat failed: %s", err)
		return
	}
	if s.Mode().String() != fileMode.String() {
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), fileMode.String())
		return
	}

	// Test Mkdir
	err = fs.Mkdir(dirPath, dirMode)
	if err != nil {
		t.Errorf("MkDir Create failed: %s", err)
		return
	}
	s, err = fs.Stat(dirPath)
	if err != nil {
		t.Errorf("Stat failed: %s", err)
		return
	}
	// sets File
	if s.Mode().String() != dirMode.String() {
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), dirMode.String())
		return
	}

	// Test MkdirAll
	err = fs.MkdirAll(dirPathAll, dirMode)
	if err != nil {
		t.Errorf("MkDir Create failed: %s", err)
		return
	}
	s, err = fs.Stat(dirPathAll)
	if err != nil {
		t.Errorf("Stat failed: %s", err)
		return
	}
	if s.Mode().String() != dirMode.String() {
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), dirMode.String())
		return
	}
}

// Fails if multiple file objects use the same file.at counter in MemMapFs
func TestMultipleOpenFiles(t *testing.T) {
	defer removeAllTestFiles(t)
	const fileName = "afero-demo2.txt"

	var data = make([][]byte, len(Fss))

	for i, fs := range Fss {
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)
		fh1, err := fs.Create(path)
		if err != nil {
			t.Error("fs.Create failed: " + err.Error())
		}
		_, err = fh1.Write([]byte("test"))
		if err != nil {
			t.Error("fh.Write failed: " + err.Error())
		}
		_, err = fh1.Seek(0, os.SEEK_SET)
		if err != nil {
			t.Error(err)
		}

		fh2, err := fs.OpenFile(path, os.O_RDWR, 0777)
		if err != nil {
			t.Error("fs.OpenFile failed: " + err.Error())
		}
		_, err = fh2.Seek(0, os.SEEK_END)
		if err != nil {
			t.Error(err)
		}
		_, err = fh2.Write([]byte("data"))
		if err != nil {
			t.Error(err)
		}
		err = fh2.Close()
		if err != nil {
			t.Error(err)
		}

		_, err = fh1.Write([]byte("data"))
		if err != nil {
			t.Error(err)
		}
		err = fh1.Close()
		if err != nil {
			t.Error(err)
		}
		// the file now should contain "datadata"
		data[i], err = ReadFile(fs, path)
		if err != nil {
			t.Error(err)
		}
	}

	for i, fs := range Fss {
		if i == 0 {
			continue
		}
		if string(data[0]) != string(data[i]) {
			t.Errorf("%s and %s don't behave the same\n"+
				"%s: \"%s\"\n%s: \"%s\"\n",
				Fss[0].Name(), fs.Name(), Fss[0].Name(), data[0], fs.Name(), data[i])
		}
	}
}

// Test if file.Write() fails when opened as read only
func TestReadOnly(t *testing.T) {
	defer removeAllTestFiles(t)
	const fileName = "afero-demo.txt"

	for _, fs := range Fss {
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)

		f, err := fs.Create(path)
		if err != nil {
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		}
		_, err = f.Write([]byte("test"))
		if err != nil {
			t.Error(fs.Name()+":", "Write failed: "+err.Error())
		}
		f.Close()

		f, err = fs.Open(path)
		if err != nil {
			t.Error("fs.Open failed: " + err.Error())
		}
		_, err = f.Write([]byte("data"))
		if err == nil {
			t.Error(fs.Name()+":", "No write error")
		}
		f.Close()

		f, err = fs.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			t.Error("fs.Open failed: " + err.Error())
		}
		_, err = f.Write([]byte("data"))
		if err == nil {
			t.Error(fs.Name()+":", "No write error")
		}
		f.Close()
	}
}

func TestWriteCloseTime(t *testing.T) {
	defer removeAllTestFiles(t)
	const fileName = "afero-demo.txt"

	for _, fs := range Fss {
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)

		f, err := fs.Create(path)
		if err != nil {
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		}
		f.Close()

		f, err = fs.Create(path)
		if err != nil {
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		}
		fi, err := f.Stat()
		if err != nil {
			t.Error(fs.Name()+":", "Stat failed: "+err.Error())
		}
		timeBefore := fi.ModTime()

		// sorry for the delay, but we have to make sure time advances,
		// also on non Un*x systems...
		switch runtime.GOOS {
		case "windows":
			time.Sleep(2 * time.Second)
		case "darwin":
			time.Sleep(1 * time.Second)
		default: // depending on the FS, this may work with < 1 second, on my old ext3 it does not
			time.Sleep(1 * time.Second)
		}

		_, err = f.Write([]byte("test"))
		if err != nil {
			t.Error(fs.Name()+":", "Write failed: "+err.Error())
		}
		f.Close()
		fi, err = fs.Stat(path)
		if err != nil {
			t.Error(fs.Name()+":", "fs.Stat failed: "+err.Error())
		}
		if fi.ModTime().Equal(timeBefore) {
			t.Error(fs.Name()+":", "ModTime was not set on Close()")
		}
	}
}

// This test should be run with the race detector on:
// go test -race -v -timeout 10s -run TestRacingDeleteAndClose
func TestRacingDeleteAndClose(t *testing.T) {
	fs := NewMemMapFs()
	pathname := "testfile"
	f, err := fs.Create(pathname)
	if err != nil {
		t.Fatal(err)
	}

	in := make(chan bool)

	go func() {
		<-in
		f.Close()
	}()
	go func() {
		<-in
		fs.Remove(pathname)
	}()
	close(in)
}

// This test should be run with the race detector on:
// go test -run TestMemFsDataRace -race
func TestMemFsDataRace(t *testing.T) {
	const dir = "test_dir"
	fs := NewMemMapFs()

	if err := fs.MkdirAll(dir, 0777); err != nil {
		t.Fatal(err)
	}

	const n = 1000
	done := make(chan struct{})

	go func() {
		defer close(done)
		for i := 0; i < n; i++ {
			fname := filepath.Join(dir, fmt.Sprintf("%d.txt", i))
			if err := WriteFile(fs, fname, []byte(""), 0777); err != nil {
				panic(err)
			}
			if err := fs.Remove(fname); err != nil {
				panic(err)
			}
		}
	}()

loop:
	for {
		select {
		case <-done:
			break loop
		default:
			_, err := ReadDir(fs, dir)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestMemFsDirMode(t *testing.T) {
	fs := NewMemMapFs()
	err := fs.Mkdir("/testDir1", 0644)
	if err != nil {
		t.Error(err)
	}
	err = fs.MkdirAll("/sub/testDir2", 0644)
	if err != nil {
		t.Error(err)
	}
	info, err := fs.Stat("/testDir1")
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Error("should be a directory")
	}
	if !info.Mode().IsDir() {
		t.Error("FileMode is not directory")
	}
	info, err = fs.Stat("/sub/testDir2")
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Error("should be a directory")
	}
	if !info.Mode().IsDir() {
		t.Error("FileMode is not directory")
	}
}

func TestMemFsUnexpectedEOF(t *testing.T) {
	t.Parallel()

	fs := NewMemMapFs()

	if err := WriteFile(fs, "file.txt", []byte("abc"), 0777); err != nil {
		t.Fatal(err)
	}

	f, err := fs.Open("file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Seek beyond the end.
	_, err = f.Seek(512, 0)
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 256)
	_, err = io.ReadAtLeast(f, buff, 256)

	if err != io.ErrUnexpectedEOF {
		t.Fatal("Expected ErrUnexpectedEOF")
	}
}
