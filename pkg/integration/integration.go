package integration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/integration/helpers"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	"github.com/stretchr/testify/assert"
)

// this is the integration runner for the new and improved integration interface

var Tests = tests.Tests

type Mode int

const (
	// Default: if a snapshot test fails, the we'll be asked whether we want to update it
	ASK_TO_UPDATE_SNAPSHOT = iota
	// fails the test if the snapshots don't match
	CHECK_SNAPSHOT
	// runs the test and updates the snapshot
	UPDATE_SNAPSHOT
	// This just makes use of the setup step of the test to get you into
	// a lazygit session. Then you'll be able to do whatever you want. Useful
	// when you want to test certain things without needing to manually set
	// up the situation yourself.
	// fails the test if the snapshots don't match
	SANDBOX
)

type (
	logf func(format string, formatArgs ...interface{})
)

func RunTestsNew(
	logf logf,
	runCmd func(cmd *exec.Cmd) error,
	fnWrapper func(test *helpers.Test, f func() error),
	mode Mode,
	includeSkipped bool,
) error {
	rootDir := GetRootDirectory()
	err := os.Chdir(rootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(rootDir, "test", "integration_new")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.Cmd.New("go build -o " + tempLazygitPath()).Run()
	if err != nil {
		return err
	}

	for _, test := range Tests {
		test := test

		fnWrapper(test, func() error { //nolint: thelper
			if test.Skip() && !includeSkipped {
				logf("skipping test: %s", test.Name())
				return nil
			}

			testPath := filepath.Join(testDir, test.Name())

			actualDir := filepath.Join(testPath, "actual")
			expectedDir := filepath.Join(testPath, "expected")
			actualRepoDir := filepath.Join(actualDir, "repo")
			logf("path: %s", testPath)

			findOrCreateDir(testPath)
			prepareIntegrationTestDir(actualDir)
			findOrCreateDir(actualRepoDir)
			err := createFixtureNew(test, actualRepoDir, rootDir)
			if err != nil {
				return err
			}

			configDir := filepath.Join(testPath, "used_config")

			cmd, err := getLazygitCommandNew(test, testPath, rootDir)
			if err != nil {
				return err
			}

			err = runCmd(cmd)
			if err != nil {
				return err
			}

			switch mode {
			case UPDATE_SNAPSHOT:
				if err := updateSnapshot(logf, actualDir, expectedDir); err != nil {
					return err
				}
				logf("Test passed: %s", test.Name())
			case CHECK_SNAPSHOT:
				if err := compareSnapshots(logf, configDir, actualDir, expectedDir, test.Name()); err != nil {
					return err
				}
				logf("Test passed: %s", test.Name())
			case ASK_TO_UPDATE_SNAPSHOT:
				if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
					if err := updateSnapshot(logf, actualDir, expectedDir); err != nil {
						return err
					}
					logf("No existing snapshot found for  %s. Created snapshot.", test.Name())

					return nil
				}

				if err := compareSnapshots(logf, configDir, actualDir, expectedDir, test.Name()); err != nil {
					logf("%s", err)

					// prompt user whether to update the snapshot (Y/N)
					if promptUserToUpdateSnapshot() {
						if err := updateSnapshot(logf, actualDir, expectedDir); err != nil {
							return err
						}
						logf("Snapshot updated: %s", test.Name())
					} else {
						return err
					}
				}

				logf("Test passed: %s", test.Name())
			case SANDBOX:
				logf("Session exited")
			}

			return nil
		})
	}

	return nil
}

func promptUserToUpdateSnapshot() bool {
	fmt.Println("Test failed. Update snapshot? (y/n)")
	var input string
	fmt.Scanln(&input)
	return input == "y"
}

func updateSnapshot(logf logf, actualDir string, expectedDir string) error {
	// create/update snapshot
	err := oscommands.CopyDir(actualDir, expectedDir)
	if err != nil {
		return err
	}

	if err := renameSpecialPaths(expectedDir); err != nil {
		return err
	}

	return err
}

func compareSnapshots(logf logf, configDir string, actualDir string, expectedDir string, testName string) error {
	// there are a couple of reasons we're not generating the snapshot in expectedDir directly:
	// Firstly we don't want to have to revert our .git file back to .git_keep.
	// Secondly, the act of calling git commands like 'git status' actually changes the index
	// for some reason, and we don't want to leave your lazygit working tree dirty as a result.
	expectedDirCopy := filepath.Join(os.TempDir(), "expected_dir_test", testName)
	err := oscommands.CopyDir(expectedDir, expectedDirCopy)
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

	err = validateSameRepos(expectedDirCopy, actualDir)
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
		actualRepoPath := filepath.Join(actualDir, f.Name())
		expectedRepoPath := filepath.Join(expectedDirCopy, f.Name())

		actualRepo, expectedRepo, err := generateSnapshots(actualRepoPath, expectedRepoPath)
		if err != nil {
			return err
		}

		if expectedRepo != actualRepo {
			// get the log file and print it
			bytes, err := ioutil.ReadFile(filepath.Join(configDir, "development.log"))
			if err != nil {
				return err
			}
			logf("%s", string(bytes))

			return errors.New(getDiff(f.Name(), actualRepo, expectedRepo))
		}
	}

	return nil
}

func createFixtureNew(test *helpers.Test, actualDir string, rootDir string) error {
	if err := os.Chdir(actualDir); err != nil {
		panic(err)
	}

	shell := helpers.NewShell()
	shell.RunCommand("git init")
	shell.RunCommand(`git config user.email "CI@example.com"`)
	shell.RunCommand(`git config user.name "CI"`)

	test.SetupRepo(shell)

	// changing directory back to rootDir after the setup is done
	if err := os.Chdir(rootDir); err != nil {
		panic(err)
	}

	return nil
}

func getLazygitCommandNew(test *helpers.Test, testPath string, rootDir string) (*exec.Cmd, error) {
	osCommand := oscommands.NewDummyOSCommand()

	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")
	actualRepoDir := filepath.Join(testPath, "actual", "repo")

	configDir := filepath.Join(testPath, "used_config")

	err := os.RemoveAll(configDir)
	if err != nil {
		return nil, err
	}
	err = oscommands.CopyDir(templateConfigDir, configDir)
	if err != nil {
		return nil, err
	}

	cmdStr := fmt.Sprintf("%s -debug --use-config-dir=%s --path=%s %s", tempLazygitPath(), configDir, actualRepoDir, test.ExtraCmdArgs())

	cmdObj := osCommand.Cmd.New(cmdStr)

	cmdObj.AddEnvVars(fmt.Sprintf("LAZYGIT_TEST_NAME=%s", test.Name()))

	return cmdObj.GetCmd(), nil
}

func GetModeFromEnv() Mode {
	switch os.Getenv("MODE") {
	case "", "ask":
		return ASK_TO_UPDATE_SNAPSHOT
	case "check":
		return CHECK_SNAPSHOT
	case "updateSnapshot":
		return UPDATE_SNAPSHOT
	case "sandbox":
		return SANDBOX
	default:
		log.Fatalf("unknown test mode: %s, must be one of [test, record, updateSnapshot, sandbox]", os.Getenv("MODE"))
		panic("unreachable")
	}
}

func GetRootDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		_, err := os.Stat(filepath.Join(path, ".git"))

		if err == nil {
			return path
		}

		if !os.IsNotExist(err) {
			panic(err)
		}

		path = filepath.Dir(path)

		if path == "/" {
			log.Fatal("must run in lazygit folder or child folder")
		}
	}
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
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

		bytes, err := ioutil.ReadFile(path)
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

func findOrCreateDir(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0o777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func prepareIntegrationTestDir(actualDir string) {
	// remove contents of integration test directory
	dir, err := ioutil.ReadDir(actualDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(actualDir, 0o777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	for _, d := range dir {
		os.RemoveAll(filepath.Join(actualDir, d.Name()))
	}
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
