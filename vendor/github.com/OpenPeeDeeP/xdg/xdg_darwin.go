// Copyright (c) 2017, OpenPeeDeeP. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdg

import (
	"os"
	"path/filepath"
)

func (o *osDefaulter) defaultDataHome() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
}

func (o *osDefaulter) defaultDataDirs() []string {
	return []string{filepath.Join("/Library", "Application Support")}
}

func (o *osDefaulter) defaultConfigHome() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
}

func (o *osDefaulter) defaultConfigDirs() []string {
	return []string{filepath.Join("/Library", "Application Support")}
}

func (o *osDefaulter) defaultCacheHome() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
}
