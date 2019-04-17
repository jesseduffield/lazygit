// Copyright ©2018 Steve Francia <spf@spf13.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"testing"
)

func TestLstatIfPossible(t *testing.T) {
	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()

	osFs := &OsFs{}

	workDir, err := TempDir(osFs, "", "afero-lstate")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		osFs.RemoveAll(workDir)
	}()

	memWorkDir := "/lstate"

	memFs := NewMemMapFs()
	overlayFs1 := &CopyOnWriteFs{base: osFs, layer: memFs}
	overlayFs2 := &CopyOnWriteFs{base: memFs, layer: osFs}
	overlayFsMemOnly := &CopyOnWriteFs{base: memFs, layer: NewMemMapFs()}
	basePathFs := &BasePathFs{source: osFs, path: workDir}
	basePathFsMem := &BasePathFs{source: memFs, path: memWorkDir}
	roFs := &ReadOnlyFs{source: osFs}
	roFsMem := &ReadOnlyFs{source: memFs}

	pathFileMem := filepath.Join(memWorkDir, "aferom.txt")

	WriteFile(osFs, filepath.Join(workDir, "afero.txt"), []byte("Hi, Afero!"), 0777)
	WriteFile(memFs, filepath.Join(pathFileMem), []byte("Hi, Afero!"), 0777)

	os.Chdir(workDir)
	if err := os.Symlink("afero.txt", "symafero.txt"); err != nil {
		t.Fatal(err)
	}

	pathFile := filepath.Join(workDir, "afero.txt")
	pathSymlink := filepath.Join(workDir, "symafero.txt")

	checkLstat := func(l Lstater, name string, shouldLstat bool) os.FileInfo {
		statFile, isLstat, err := l.LstatIfPossible(name)
		if err != nil {
			t.Fatalf("Lstat check failed: %s", err)
		}
		if isLstat != shouldLstat {
			t.Fatalf("Lstat status was %t for %s", isLstat, name)
		}
		return statFile
	}

	testLstat := func(l Lstater, pathFile, pathSymlink string) {
		shouldLstat := pathSymlink != ""
		statRegular := checkLstat(l, pathFile, shouldLstat)
		statSymlink := checkLstat(l, pathSymlink, shouldLstat)
		if statRegular == nil || statSymlink == nil {
			t.Fatal("got nil FileInfo")
		}

		symSym := statSymlink.Mode()&os.ModeSymlink == os.ModeSymlink
		if symSym == (pathSymlink == "") {
			t.Fatal("expected the FileInfo to describe the symlink")
		}

		_, _, err := l.LstatIfPossible("this-should-not-exist.txt")
		if err == nil || !os.IsNotExist(err) {
			t.Fatalf("expected file to not exist, got %s", err)
		}
	}

	testLstat(osFs, pathFile, pathSymlink)
	testLstat(overlayFs1, pathFile, pathSymlink)
	testLstat(overlayFs2, pathFile, pathSymlink)
	testLstat(basePathFs, "afero.txt", "symafero.txt")
	testLstat(overlayFsMemOnly, pathFileMem, "")
	testLstat(basePathFsMem, "aferom.txt", "")
	testLstat(roFs, pathFile, pathSymlink)
	testLstat(roFsMem, pathFileMem, "")
}
