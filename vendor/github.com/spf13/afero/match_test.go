// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2009 The Go Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package afero

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// contains returns true if vector contains the string s.
func contains(vector []string, s string) bool {
	for _, elem := range vector {
		if elem == s {
			return true
		}
	}
	return false
}

func setupGlobDirRoot(t *testing.T, fs Fs) string {
	path := testDir(fs)
	setupGlobFiles(t, fs, path)
	return path
}

func setupGlobDirReusePath(t *testing.T, fs Fs, path string) string {
	testRegistry[fs] = append(testRegistry[fs], path)
	return setupGlobFiles(t, fs, path)
}

func setupGlobFiles(t *testing.T, fs Fs, path string) string {
	testSubDir := filepath.Join(path, "globs", "bobs")
	err := fs.MkdirAll(testSubDir, 0700)
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}

	f, err := fs.Create(filepath.Join(testSubDir, "/matcher"))
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Testfile 1 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "/../submatcher"))
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Testfile 2 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "/../../match"))
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Testfile 3 content")
	f.Close()

	return testSubDir
}

func TestGlob(t *testing.T) {
	defer removeAllTestFiles(t)
	var testDir string
	for i, fs := range Fss {
		if i == 0 {
			testDir = setupGlobDirRoot(t, fs)
		} else {
			setupGlobDirReusePath(t, fs, testDir)
		}
	}

	var globTests = []struct {
		pattern, result string
	}{
		{testDir + "/globs/bobs/matcher", testDir + "/globs/bobs/matcher"},
		{testDir + "/globs/*/mat?her", testDir + "/globs/bobs/matcher"},
		{testDir + "/globs/bobs/../*", testDir + "/globs/submatcher"},
		{testDir + "/match", testDir + "/match"},
	}

	for _, fs := range Fss {

		for _, tt := range globTests {
			pattern := tt.pattern
			result := tt.result
			if runtime.GOOS == "windows" {
				pattern = filepath.Clean(pattern)
				result = filepath.Clean(result)
			}
			matches, err := Glob(fs, pattern)
			if err != nil {
				t.Errorf("Glob error for %q: %s", pattern, err)
				continue
			}
			if !contains(matches, result) {
				t.Errorf("Glob(%#q) = %#v want %v", pattern, matches, result)
			}
		}
		for _, pattern := range []string{"no_match", "../*/no_match"} {
			matches, err := Glob(fs, pattern)
			if err != nil {
				t.Errorf("Glob error for %q: %s", pattern, err)
				continue
			}
			if len(matches) != 0 {
				t.Errorf("Glob(%#q) = %#v want []", pattern, matches)
			}
		}

	}
}

func TestGlobSymlink(t *testing.T) {
	defer removeAllTestFiles(t)

	fs := &OsFs{}
	testDir := setupGlobDirRoot(t, fs)

	err := os.Symlink("target", filepath.Join(testDir, "symlink"))
	if err != nil {
		t.Skipf("skipping on %s", runtime.GOOS)
	}

	var globSymlinkTests = []struct {
		path, dest string
		brokenLink bool
	}{
		{"test1", "link1", false},
		{"test2", "link2", true},
	}

	for _, tt := range globSymlinkTests {
		path := filepath.Join(testDir, tt.path)
		dest := filepath.Join(testDir, tt.dest)
		f, err := fs.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		err = os.Symlink(path, dest)
		if err != nil {
			t.Fatal(err)
		}
		if tt.brokenLink {
			// Break the symlink.
			fs.Remove(path)
		}
		matches, err := Glob(fs, dest)
		if err != nil {
			t.Errorf("GlobSymlink error for %q: %s", dest, err)
		}
		if !contains(matches, dest) {
			t.Errorf("Glob(%#q) = %#v want %v", dest, matches, dest)
		}
	}
}


func TestGlobError(t *testing.T) {
	for _, fs := range Fss {
		_, err := Glob(fs, "[7]")
		if err != nil {
			t.Error("expected error for bad pattern; got none")
		}
	}
}
