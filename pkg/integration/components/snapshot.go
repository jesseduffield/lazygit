package components

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

// This creates and compares integration test snapshots.

type (
	logf func(format string, formatArgs ...interface{})
)

func HandleSnapshots(paths Paths, logf logf, test *IntegrationTest, mode Mode) error {
	return NewSnapshotter(paths, logf, test, mode).
		handleSnapshots()
}

type Snapshotter struct {
	paths Paths
	logf  logf
	test  *IntegrationTest
	mode  Mode
}

func NewSnapshotter(
	paths Paths,
	logf logf,
	test *IntegrationTest,
	mode Mode,
) *Snapshotter {
	return &Snapshotter{
		paths: paths,
		logf:  logf,
		test:  test,
		mode:  mode,
	}
}

func (self *Snapshotter) handleSnapshots() error {
	switch self.mode {
	case UPDATE_SNAPSHOT:
		return self.handleUpdate()
	case CHECK_SNAPSHOT:
		return self.handleCheck()
	case ASK_TO_UPDATE_SNAPSHOT:
		return self.handleAskToUpdate()
	case SANDBOX:
		self.logf("Sandbox session exited")
	}
	return nil
}

func (self *Snapshotter) handleUpdate() error {
	if err := self.updateSnapshot(); err != nil {
		return err
	}
	self.logf("Test passed: %s", self.test.Name())
	return nil
}

func (self *Snapshotter) handleCheck() error {
	self.logf("Comparing snapshots")
	if err := self.compareSnapshots(); err != nil {
		return err
	}
	self.logf("Test passed: %s", self.test.Name())
	return nil
}

func (self *Snapshotter) handleAskToUpdate() error {
	if _, err := os.Stat(self.paths.Expected()); os.IsNotExist(err) {
		if err := self.updateSnapshot(); err != nil {
			return err
		}
		self.logf("No existing snapshot found for  %s. Created snapshot.", self.test.Name())

		return nil
	}

	self.logf("Comparing snapshots...")
	if err := self.compareSnapshots(); err != nil {
		self.logf("%s", err)

		// prompt user whether to update the snapshot (Y/N)
		if promptUserToUpdateSnapshot() {
			if err := self.updateSnapshot(); err != nil {
				return err
			}
			self.logf("Snapshot updated: %s", self.test.Name())
		} else {
			return err
		}
	}

	self.logf("Test passed: %s", self.test.Name())
	return nil
}

func (self *Snapshotter) updateSnapshot() error {
	// create/update snapshot
	err := oscommands.CopyDir(self.paths.Actual(), self.paths.Expected())
	if err != nil {
		return err
	}

	if err := renameSpecialPaths(self.paths.Expected()); err != nil {
		return err
	}

	return nil
}

func (self *Snapshotter) compareSnapshots() error {
	// there are a couple of reasons we're not generating the snapshot in expectedDir directly:
	// Firstly we don't want to have to revert our .git file back to .git_keep.
	// Secondly, the act of calling git commands like 'git status' actually changes the index
	// for some reason, and we don't want to leave your lazygit working tree dirty as a result.
	expectedDirCopy := filepath.Join(os.TempDir(), "expected_dir_test", self.test.Name())
	err := oscommands.CopyDir(self.paths.Expected(), expectedDirCopy)
	if err != nil {
		return err
	}

	defer func() {
		err := os.RemoveAll(expectedDirCopy)
		if err != nil {
			panic(err)
		}
	}()

	if err := restoreSpecialPaths(expectedDirCopy); err != nil {
		return err
	}

	err = validateSameRepos(expectedDirCopy, self.paths.Actual())
	if err != nil {
		return err
	}

	// iterate through each repo in the expected dir and comparet to the corresponding repo in the actual dir
	expectedFiles, err := ioutil.ReadDir(expectedDirCopy)
	if err != nil {
		return err
	}

	for _, f := range expectedFiles {
		if !f.IsDir() {
			return errors.New("unexpected file (as opposed to directory) in integration test 'expected' directory")
		}

		// get corresponding file name from actual dir
		actualRepoPath := filepath.Join(self.paths.Actual(), f.Name())
		expectedRepoPath := filepath.Join(expectedDirCopy, f.Name())

		actualRepo, expectedRepo, err := generateSnapshots(actualRepoPath, expectedRepoPath)
		if err != nil {
			return err
		}

		if expectedRepo != actualRepo {
			// get the log file and print it
			bytes, err := os.ReadFile(filepath.Join(self.paths.Config(), "development.log"))
			if err != nil {
				return err
			}
			self.logf("%s", string(bytes))

			return errors.New(getDiff(f.Name(), expectedRepo, actualRepo))
		}
	}

	return nil
}

func promptUserToUpdateSnapshot() bool {
	fmt.Println("Test failed. Update snapshot? (y/n)")
	var input string
	fmt.Scanln(&input)
	return input == "y"
}

func generateSnapshots(actualDir string, expectedDir string) (string, string, error) {
	actual, err := generateSnapshot(actualDir)
	if err != nil {
		return "", "", err
	}

	expected, err := generateSnapshot(expectedDir)
	if err != nil {
		return "", "", err
	}

	return actual, expected, nil
}

// note that we don't actually store this snapshot in the lazygit repo.
// Instead we store the whole expected git repo of our test, so that
// we can easily change what we want to compare without needing to regenerate
// snapshots for each test.
func generateSnapshot(dir string) (string, error) {
	osCommand := oscommands.NewDummyOSCommand()

	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return "git directory not found", nil
	}

	snapshot := ""

	cmdStrs := []string{
		`remote show -n origin`, // remote branches
		// TODO: find a way to bring this back without breaking tests
		// `ls-remote origin`,
		`status`,                         // file tree
		`log --pretty=%B|%an|%ae -p -1`,  // log
		`tag -n`,                         // tags
		`stash list`,                     // stash
		`submodule foreach 'git status'`, // submodule status
		`submodule foreach 'git log --pretty=%B -p -1'`, // submodule log
		`submodule foreach 'git tag -n'`,                // submodule tags
		`submodule foreach 'git stash list'`,            // submodule stash
	}

	for _, cmdStr := range cmdStrs {
		// ignoring error for now. If there's an error it could be that there are no results
		output, _ := osCommand.Cmd.New(fmt.Sprintf("git -C %s %s", dir, cmdStr)).RunWithOutput()

		snapshot += fmt.Sprintf("git %s:\n%s\n", cmdStr, output)
	}

	snapshot += "files in repo:\n"
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			if f.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		snapshot += fmt.Sprintf("path: %s\ncontent:\n%s\n", relativePath, string(bytes))

		return nil
	})

	if err != nil {
		return "", err
	}

	return snapshot, nil
}

func getPathsToRename(dir string, needle string, contains string) []string {
	pathsToRename := []string{}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Name() == needle && (contains == "" || strings.Contains(path, contains)) {
			pathsToRename = append(pathsToRename, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return pathsToRename
}

var specialPathMappings = []struct{ original, new, contains string }{
	// git refuses to track .git or .gitmodules in subdirectories so we need to rename them
	{".git", ".git_keep", ""},
	{".gitmodules", ".gitmodules_keep", ""},
	// we also need git to ignore the contents of our test gitignore files so that
	// we actually commit files that are ignored within the test.
	{".gitignore", "lg_ignore_file", ""},
	// this is the .git/info/exclude file. We're being a little more specific here
	// so that we don't accidentally mess with some other file named 'exclude' in the test.
	{"exclude", "lg_exclude_file", ".git/info/exclude"},
}

func renameSpecialPaths(dir string) error {
	for _, specialPath := range specialPathMappings {
		for _, path := range getPathsToRename(dir, specialPath.original, specialPath.contains) {
			err := os.Rename(path, filepath.Join(filepath.Dir(path), specialPath.new))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func restoreSpecialPaths(dir string) error {
	for _, specialPath := range specialPathMappings {
		for _, path := range getPathsToRename(dir, specialPath.new, specialPath.contains) {
			err := os.Rename(path, filepath.Join(filepath.Dir(path), specialPath.original))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// validates that the actual and expected dirs have the same repo names (doesn't actually check the contents of the repos)
func validateSameRepos(expectedDir string, actualDir string) error {
	// iterate through each repo in the expected dir and compare to the corresponding repo in the actual dir
	expectedFiles, err := ioutil.ReadDir(expectedDir)
	if err != nil {
		return err
	}

	var actualFiles []os.FileInfo
	actualFiles, err = ioutil.ReadDir(actualDir)
	if err != nil {
		return err
	}

	expectedFileNames := slices.Map(expectedFiles, getFileName)
	actualFileNames := slices.Map(actualFiles, getFileName)
	if !slices.Equal(expectedFileNames, actualFileNames) {
		return fmt.Errorf("expected and actual repo dirs do not match: expected: %s, actual: %s", expectedFileNames, actualFileNames)
	}

	return nil
}

func getFileName(f os.FileInfo) string {
	return f.Name()
}

func getDiff(prefix string, expected string, actual string) string {
	mockT := &MockTestingT{}
	assert.Equal(mockT, expected, actual, fmt.Sprintf("Unexpected %s. Expected:\n%s\nActual:\n%s\n", prefix, expected, actual))
	return mockT.err
}

type MockTestingT struct {
	err string
}

func (self *MockTestingT) Errorf(format string, args ...interface{}) {
	self.err += fmt.Sprintf(format, args...)
}
