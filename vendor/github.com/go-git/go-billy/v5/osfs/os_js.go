//go:build js
// +build js

package osfs

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-billy/v5/memfs"
)

// globalMemFs is the global memory fs
var globalMemFs = memfs.New()

// Default Filesystem representing the root of in-memory filesystem for a
// js/wasm environment.
var Default = memfs.New()

// New returns a new OS filesystem.
func New(baseDir string, _ ...Option) billy.Filesystem {
	return chroot.New(Default, Default.Join("/", baseDir))
}

type options struct {
}
