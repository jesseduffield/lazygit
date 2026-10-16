package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/samber/lo"
)

// FindLazygitRootDirectory returns the root directory of the lazygit source
// tree, by searching the working directory and its parents for the go.mod file
// that declares lazygit's module. Only development tools use it: the
// integration test runner, the cheatsheet generator, and the JSON schema
// generator. Not to be confused with finding the root directory of the
// repository that lazygit is being run in.
//
// We search upwards rather than expect to be called from the root directory,
// because `go test` runs each test binary in the source directory of its
// package, not in the directory that `go test` was invoked from.
func FindLazygitRootDirectory() (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := startDir
	for {
		if declaresLazygitModule(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(
				"failed to find the lazygit root directory: there is no go.mod for a lazygit module in %s or any of its parent directories",
				startDir)
		}
		dir = parent
	}
}

// MustFindLazygitRootDirectory is FindLazygitRootDirectory for tools that can't
// do anything useful if the directory isn't found.
func MustFindLazygitRootDirectory() string {
	dir, err := FindLazygitRootDirectory()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// A fork is free to rename the module, so we accept any module path with a
// "lazygit" element in it, e.g. github.com/jesseduffield/lazygit.
func declaresLazygitModule(goModPath string) bool {
	contents, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	return lo.SomeBy(strings.Split(string(contents), "\n"), func(line string) bool {
		fields := strings.Fields(line)
		return len(fields) >= 2 && fields[0] == "module" &&
			slices.Contains(strings.Split(fields[1], "/"), "lazygit")
	})
}
