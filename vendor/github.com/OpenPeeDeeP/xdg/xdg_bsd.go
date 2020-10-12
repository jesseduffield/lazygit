// +build freebsd openbsd netbsd

// Copyright (c) 2017, OpenPeeDeeP. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdg

import (
	"os"
	"path/filepath"
)

func (o *osDefaulter) defaultDataHome() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "share")
}

func (o *osDefaulter) defaultDataDirs() []string {
	return []string{"/usr/local/share/", "/usr/share/"}
}

func (o *osDefaulter) defaultConfigHome() string {
	return filepath.Join(os.Getenv("HOME"), ".config")
}

func (o *osDefaulter) defaultConfigDirs() []string {
	return []string{"/etc/xdg"}
}

func (o *osDefaulter) defaultCacheHome() string {
	return filepath.Join(os.Getenv("HOME"), ".cache")
}
