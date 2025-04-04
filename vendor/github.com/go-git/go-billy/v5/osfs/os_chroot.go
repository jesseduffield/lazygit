//go:build !js
// +build !js

package osfs

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/chroot"
)

// ChrootOS is a legacy filesystem based on a "soft chroot" of the os filesystem.
// Although this is still the default os filesystem, consider using BoundOS instead.
//
// Behaviours of note:
//  1. A "soft chroot" translates the base dir to "/" for the purposes of the
//     fs abstraction.
//  2. Symlinks targets may be modified to be kept within the chroot bounds.
//  3. Some file modes does not pass-through the fs abstraction.
//  4. The combination of 1 and 2 may cause go-git to think that a Git repository
//     is dirty, when in fact it isn't.
type ChrootOS struct{}

func newChrootOS(baseDir string) billy.Filesystem {
	return chroot.New(&ChrootOS{}, baseDir)
}

func (fs *ChrootOS) Create(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, defaultCreateMode)
}

func (fs *ChrootOS) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	return openFile(filename, flag, perm, fs.createDir)
}

func (fs *ChrootOS) createDir(fullpath string) error {
	dir := filepath.Dir(fullpath)
	if dir != "." {
		if err := os.MkdirAll(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}

	return nil
}

func (fs *ChrootOS) ReadDir(dir string) ([]os.FileInfo, error) {
	return readDir(dir)
}

func (fs *ChrootOS) Rename(from, to string) error {
	if err := fs.createDir(to); err != nil {
		return err
	}

	return rename(from, to)
}

func (fs *ChrootOS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, defaultDirectoryMode)
}

func (fs *ChrootOS) Open(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDONLY, 0)
}

func (fs *ChrootOS) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

func (fs *ChrootOS) Remove(filename string) error {
	return os.Remove(filename)
}

func (fs *ChrootOS) TempFile(dir, prefix string) (billy.File, error) {
	if err := fs.createDir(dir + string(os.PathSeparator)); err != nil {
		return nil, err
	}

	return tempFile(dir, prefix)
}

func (fs *ChrootOS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (fs *ChrootOS) RemoveAll(path string) error {
	return os.RemoveAll(filepath.Clean(path))
}

func (fs *ChrootOS) Lstat(filename string) (os.FileInfo, error) {
	return os.Lstat(filepath.Clean(filename))
}

func (fs *ChrootOS) Symlink(target, link string) error {
	if err := fs.createDir(link); err != nil {
		return err
	}

	return os.Symlink(target, link)
}

func (fs *ChrootOS) Readlink(link string) (string, error) {
	return os.Readlink(link)
}

// Capabilities implements the Capable interface.
func (fs *ChrootOS) Capabilities() billy.Capability {
	return billy.DefaultCapabilities
}
