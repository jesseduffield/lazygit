package components

import "path/filepath"

// convenience struct for easily getting directories within our test directory.
// We have one test directory for each test, found in test/integration_new.
type Paths struct {
	// e.g. test/integration/test_name
	root string
}

func NewPaths(root string) Paths {
	return Paths{root: root}
}

// when a test first runs, it's situated in a repo called 'repo' within this
// directory. In its setup step, the test is allowed to create other repos
// alongside the 'repo' repo in this directory, for example, creating remotes
// or repos to add as submodules.
func (self Paths) Actual() string {
	return filepath.Join(self.root, "actual")
}

// this is the 'repo' directory within the 'actual' directory,
// where a lazygit test will start within.
func (self Paths) ActualRepo() string {
	return filepath.Join(self.Actual(), "repo")
}

// When an integration test first runs, we copy everything in the 'actual' directory,
// and copy it into the 'expected' directory so that future runs can be compared
// against what we expect.
func (self Paths) Expected() string {
	return filepath.Join(self.root, "expected")
}

func (self Paths) Config() string {
	return filepath.Join(self.root, "used_config")
}

func (self Paths) Root() string {
	return self.root
}
