//go:build !js
// +build !js

/*
   Copyright 2022 The Flux authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package osfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/go-git/go-billy/v5"
)

// BoundOS is a fs implementation based on the OS filesystem which is bound to
// a base dir.
// Prefer this fs implementation over ChrootOS.
//
// Behaviours of note:
//  1. Read and write operations can only be directed to files which descends
//     from the base dir.
//  2. Symlinks don't have their targets modified, and therefore can point
//     to locations outside the base dir or to non-existent paths.
//  3. Readlink and Lstat ensures that the link file is located within the base
//     dir, evaluating any symlinks that file or base dir may contain.
type BoundOS struct {
	baseDir         string
	deduplicatePath bool
}

func newBoundOS(d string, deduplicatePath bool) billy.Filesystem {
	return &BoundOS{baseDir: d, deduplicatePath: deduplicatePath}
}

func (fs *BoundOS) Create(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, defaultCreateMode)
}

func (fs *BoundOS) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	fn, err := fs.abs(filename)
	if err != nil {
		return nil, err
	}
	return openFile(fn, flag, perm, fs.createDir)
}

func (fs *BoundOS) ReadDir(path string) ([]os.FileInfo, error) {
	dir, err := fs.abs(path)
	if err != nil {
		return nil, err
	}

	return readDir(dir)
}

func (fs *BoundOS) Rename(from, to string) error {
	f, err := fs.abs(from)
	if err != nil {
		return err
	}
	t, err := fs.abs(to)
	if err != nil {
		return err
	}

	// MkdirAll for target name.
	if err := fs.createDir(t); err != nil {
		return err
	}

	return os.Rename(f, t)
}

func (fs *BoundOS) MkdirAll(path string, perm os.FileMode) error {
	dir, err := fs.abs(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, perm)
}

func (fs *BoundOS) Open(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDONLY, 0)
}

func (fs *BoundOS) Stat(filename string) (os.FileInfo, error) {
	filename, err := fs.abs(filename)
	if err != nil {
		return nil, err
	}
	return os.Stat(filename)
}

func (fs *BoundOS) Remove(filename string) error {
	fn, err := fs.abs(filename)
	if err != nil {
		return err
	}
	return os.Remove(fn)
}

// TempFile creates a temporary file. If dir is empty, the file
// will be created within the OS Temporary dir. If dir is provided
// it must descend from the current base dir.
func (fs *BoundOS) TempFile(dir, prefix string) (billy.File, error) {
	if dir != "" {
		var err error
		dir, err = fs.abs(dir)
		if err != nil {
			return nil, err
		}
	}

	return tempFile(dir, prefix)
}

func (fs *BoundOS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (fs *BoundOS) RemoveAll(path string) error {
	dir, err := fs.abs(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func (fs *BoundOS) Symlink(target, link string) error {
	ln, err := fs.abs(link)
	if err != nil {
		return err
	}
	// MkdirAll for containing dir.
	if err := fs.createDir(ln); err != nil {
		return err
	}
	return os.Symlink(target, ln)
}

func (fs *BoundOS) Lstat(filename string) (os.FileInfo, error) {
	filename = filepath.Clean(filename)
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(fs.baseDir, filename)
	}
	if ok, err := fs.insideBaseDirEval(filename); !ok {
		return nil, err
	}
	return os.Lstat(filename)
}

func (fs *BoundOS) Readlink(link string) (string, error) {
	if !filepath.IsAbs(link) {
		link = filepath.Clean(filepath.Join(fs.baseDir, link))
	}
	if ok, err := fs.insideBaseDirEval(link); !ok {
		return "", err
	}
	return os.Readlink(link)
}

// Chroot returns a new OS filesystem, with the base dir set to the
// result of joining the provided path with the underlying base dir.
func (fs *BoundOS) Chroot(path string) (billy.Filesystem, error) {
	joined, err := securejoin.SecureJoin(fs.baseDir, path)
	if err != nil {
		return nil, err
	}
	return New(joined), nil
}

// Root returns the current base dir of the billy.Filesystem.
// This is required in order for this implementation to be a drop-in
// replacement for other upstream implementations (e.g. memory and osfs).
func (fs *BoundOS) Root() string {
	return fs.baseDir
}

func (fs *BoundOS) createDir(fullpath string) error {
	dir := filepath.Dir(fullpath)
	if dir != "." {
		if err := os.MkdirAll(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}

	return nil
}

// abs transforms filename to an absolute path, taking into account the base dir.
// Relative paths won't be allowed to ascend the base dir, so `../file` will become
// `/working-dir/file`.
//
// Note that if filename is a symlink, the returned address will be the target of the
// symlink.
func (fs *BoundOS) abs(filename string) (string, error) {
	if filename == fs.baseDir {
		filename = string(filepath.Separator)
	}

	path, err := securejoin.SecureJoin(fs.baseDir, filename)
	if err != nil {
		return "", nil
	}

	if fs.deduplicatePath {
		vol := filepath.VolumeName(fs.baseDir)
		dup := filepath.Join(fs.baseDir, fs.baseDir[len(vol):])
		if strings.HasPrefix(path, dup+string(filepath.Separator)) {
			return fs.abs(path[len(dup):])
		}
	}
	return path, nil
}

// insideBaseDir checks whether filename is located within
// the fs.baseDir.
func (fs *BoundOS) insideBaseDir(filename string) (bool, error) {
	if filename == fs.baseDir {
		return true, nil
	}
	if !strings.HasPrefix(filename, fs.baseDir+string(filepath.Separator)) {
		return false, fmt.Errorf("path outside base dir")
	}
	return true, nil
}

// insideBaseDirEval checks whether filename is contained within
// a dir that is within the fs.baseDir, by first evaluating any symlinks
// that either filename or fs.baseDir may contain.
func (fs *BoundOS) insideBaseDirEval(filename string) (bool, error) {
	// "/" contains all others.
	if fs.baseDir == "/" {
		return true, nil
	}
	dir, err := filepath.EvalSymlinks(filepath.Dir(filename))
	if dir == "" || os.IsNotExist(err) {
		dir = filepath.Dir(filename)
	}
	wd, err := filepath.EvalSymlinks(fs.baseDir)
	if wd == "" || os.IsNotExist(err) {
		wd = fs.baseDir
	}
	if filename != wd && dir != wd && !strings.HasPrefix(dir, wd+string(filepath.Separator)) {
		return false, fmt.Errorf("%q: path outside base dir %q: %w", filename, fs.baseDir, os.ErrNotExist)
	}
	return true, nil
}
