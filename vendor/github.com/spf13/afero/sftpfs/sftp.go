// Copyright © 2015 Jerry Jacobs <jerry.jacobs@xor-gate.org>.
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

package sftpfs

import (
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
)

// Fs is a afero.Fs implementation that uses functions provided by the sftp package.
//
// For details in any method, check the documentation of the sftp package
// (github.com/pkg/sftp).
type Fs struct {
	client *sftp.Client
}

func New(client *sftp.Client) afero.Fs {
	return &Fs{client: client}
}

func (s Fs) Name() string { return "sftpfs" }

func (s Fs) Create(name string) (afero.File, error) {
	return FileCreate(s.client, name)
}

func (s Fs) Mkdir(name string, perm os.FileMode) error {
	err := s.client.Mkdir(name)
	if err != nil {
		return err
	}
	return s.client.Chmod(name, perm)
}

func (s Fs) MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := s.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return err
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = s.MkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = s.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := s.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

func (s Fs) Open(name string) (afero.File, error) {
	return FileOpen(s.client, name)
}

func (s Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return nil, nil
}

func (s Fs) Remove(name string) error {
	return s.client.Remove(name)
}

func (s Fs) RemoveAll(path string) error {
	// TODO have a look at os.RemoveAll
	// https://github.com/golang/go/blob/master/src/os/path.go#L66
	return nil
}

func (s Fs) Rename(oldname, newname string) error {
	return s.client.Rename(oldname, newname)
}

func (s Fs) Stat(name string) (os.FileInfo, error) {
	return s.client.Stat(name)
}

func (s Fs) Lstat(p string) (os.FileInfo, error) {
	return s.client.Lstat(p)
}

func (s Fs) Chmod(name string, mode os.FileMode) error {
	return s.client.Chmod(name, mode)
}

func (s Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return s.client.Chtimes(name, atime, mtime)
}
