// Copyright (c) 2017, OpenPeeDeeP. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xdg

import "os"

func (o *osDefaulter) defaultDataHome() string {
	return os.Getenv("APPDATA")
}

func (o *osDefaulter) defaultDataDirs() []string {
	return []string{os.Getenv("PROGRAMDATA")}
}

func (o *osDefaulter) defaultConfigHome() string {
	return os.Getenv("APPDATA")
}

func (o *osDefaulter) defaultConfigDirs() []string {
	return []string{os.Getenv("PROGRAMDATA")}
}

func (o *osDefaulter) defaultCacheHome() string {
	return os.Getenv("LOCALAPPDATA")
}
