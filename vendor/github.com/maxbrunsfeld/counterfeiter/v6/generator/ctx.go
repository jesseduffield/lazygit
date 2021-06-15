// +build go1.14

package generator

import "go/build"

func getBuildContext(workingDir string) build.Context {
	ctx := build.Default
	ctx.Dir = workingDir
	return ctx
}
