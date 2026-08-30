// Copyright (c) 2020, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package version

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

const ourModulePath = "mvdan.cc/gofumpt"

func findModule(info *debug.BuildInfo, modulePath string) *debug.Module {
	if info.Main.Path == modulePath {
		return &info.Main
	}
	for _, dep := range info.Deps {
		if dep.Path == modulePath {
			return dep
		}
	}
	return nil
}

func gofumptVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(no build info)"
	}
	// Note that gofumpt may be used as a library via the format package,
	// so we cannot assume it is the main module in the build.
	mod := findModule(info, ourModulePath)
	if mod == nil {
		return "(module not found)"
	}
	if mod.Replace != nil {
		mod = mod.Replace
	}
	return mod.Version
}

func goVersion() string {
	// For the tests, as we don't want the Go version to change over time.
	if testVersion := os.Getenv("GO_VERSION_TEST"); testVersion != "" {
		return testVersion
	}
	return runtime.Version()
}

func String(injected string) string {
	if injected != "" {
		return fmt.Sprintf("%s (%s)", injected, goVersion())
	}
	return fmt.Sprintf("%s (%s)", gofumptVersion(), goVersion())
}
