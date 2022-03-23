package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
)

type Test struct {
	Name         string  `json:"name"`
	Speed        float64 `json:"speed"`
	Description  string  `json:"description"`
	ExtraCmdArgs string  `json:"extraCmdArgs"`
	Skip         bool    `json:"skip"`
}

type Mode int

const (
	// default: for when we're just running a test and comparing to the snapshot
	TEST = iota
	// for when we want to record a test and set the snapshot based on the result
	RECORD
	// when we just want to use the setup of the test for our own sandboxing purposes.
	// This does not record the session and does not create/update snapshots
	SANDBOX
	// running a test but updating the snapshot
	UPDATE_SNAPSHOT
)

func GetModeFromEnv() Mode {
	switch os.Getenv("MODE") {
	case "record":
		return RECORD
	case "", "test":
		return TEST
	case "updateSnapshot":
		return UPDATE_SNAPSHOT
	case "sandbox":
		return SANDBOX
	default:
		log.Fatalf("unknown test mode: %s, must be one of [test, record, update, sandbox]", os.Getenv("MODE"))
		panic("unreachable")
	}
}

// this function is used by both `go test` and from our lazyintegration gui, but
// errors need to be handled differently in each (for example go test is always
// working with *testing.T) so we pass in any differences as args here.
func RunTests(
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	fnWrapper func(test *Test, f func(*testing.T) error),
	mode Mode,
	speedEnv string,
	onFail func(t *testing.T, expected string, actual string, prefix string),
	includeSkipped bool,
) error {
	rootDir := GetRootDirectory()
	err := os.Chdir(rootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(rootDir, "test", "integration")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.Cmd.New("go build -o " + tempLazygitPath()).Run()
	if err != nil {
		return err
	}

	tests, err := LoadTests(testDir)
	if err != nil {
		return err
	}

	for _, test := range tests {
		test := test

		if test.Skip && !includeSkipped {
			logf("skipping test: %s", test.Name)
			continue
		}

		fnWrapper(test, func(t *testing.T) error { //nolint: thelper
			speeds := getTestSpeeds(test.Speed, mode, speedEnv)
			testPath := filepath.Join(testDir, test.Name)
			actualRepoDir := filepath.Join(testPath, "actual")
			expectedRepoDir := filepath.Join(testPath, "expected")
			actualRemoteDir := filepath.Join(testPath, "actual_remote")
			expectedRemoteDir := filepath.Join(testPath, "expected_remote")
			otherRepoDir := filepath.Join(testPath, "other_repo")
			logf("path: %s", testPath)

			for i, speed := range speeds {
				if mode != SANDBOX && mode != RECORD {
					logf("%s: attempting test at speed %f\n", test.Name, speed)
				}

				findOrCreateDir(testPath)
				prepareIntegrationTestDir(actualRepoDir)
				removeDir(otherRepoDir)
				removeDir(actualRemoteDir)
				err := createFixture(testPath, actualRepoDir)
				if err != nil {
					return err
				}

				configDir := filepath.Join(testPath, "used_config")

				cmd, err := getLazygitCommand(testPath, rootDir, mode, speed, test.ExtraCmdArgs)
				if err != nil {
					return err
				}

				err = runCmd(cmd)
				if err != nil {
					return err
				}

				// submodule tests currently make use of a repo called 'other_repo' but we don't want that
				// to stick around. Long-term we should have an 'actual' folder which itself contains
				// repos, and there we can put the 'repo' repo which is the main one, alongside
				// any others that we use as part of the test (including remotes). Then we'll do snapshots for
				// each of them.
				removeDir(otherRepoDir)

				if mode == UPDATE_SNAPSHOT || mode == RECORD {
					// create/update snapshot
					err = oscommands.CopyDir(actualRepoDir, expectedRepoDir)
					if err != nil {
						return err
					}

					if err := renameGitDirs(expectedRepoDir); err != nil {
						return err
					}

					// see if we have a remote dir and if so, copy it over. Otherwise, delete the expected dir because we have no remote folder.
					if folderExists(actualRemoteDir) {
						err = oscommands.CopyDir(actualRemoteDir, expectedRemoteDir)
						if err != nil {
							return err
						}
					} else {
						removeDir(expectedRemoteDir)
					}

					logf("%s", "updated snapshot")
				} else {
					// compare result to snapshot
					actualRepo, expectedRepo, err := generateSnapshots(actualRepoDir, expectedRepoDir)
					if err != nil {
						return err
					}

					actualRemote := "remote folder does not exist"
					expectedRemote := "remote folder does not exist"
					if folderExists(expectedRemoteDir) {
						actualRemote, expectedRemote, err = generateSnapshotsForRemote(actualRemoteDir, expectedRemoteDir)
						if err != nil {
							return err
						}
					} else if folderExists(actualRemoteDir) {
						actualRemote = "remote folder exists"
					}

					if expectedRepo == actualRepo && expectedRemote == actualRemote {
						logf("%s: success at speed %f\n", test.Name, speed)
						break
					}

					// if the snapshot doesn't match and we haven't tried all playback speeds different we'll retry at a slower speed
					if i == len(speeds)-1 {
						// get the log file and print that
						bytes, err := ioutil.ReadFile(filepath.Join(configDir, "development.log"))
						if err != nil {
							return err
						}
						logf("%s", string(bytes))
						if expectedRepo != actualRepo {
							onFail(t, expectedRepo, actualRepo, "repo")
						} else {
							onFail(t, expectedRemote, actualRemote, "remote")
						}
					}
				}
			}

			return nil
		})
	}

	return nil
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
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

func createFixture(testPath, actualDir string) error {
	bashScriptPath := filepath.Join(testPath, "setup.sh")
	cmd := secureexec.Command("bash", bashScriptPath, actualDir)

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}

	return nil
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
}

func getTestSpeeds(testStartSpeed float64, mode Mode, speedStr string) []float64 {
	if mode != TEST {
		// have to go at original speed if updating snapshots in case we go to fast and create a junk snapshot
		return []float64{1.0}
	}

	if speedStr != "" {
		speed, err := strconv.ParseFloat(speedStr, 64)
		if err != nil {
			panic(err)
		}
		return []float64{speed}
	}

	// default is 10, 5, 1
	startSpeed := 10.0
	if testStartSpeed != 0 {
		startSpeed = testStartSpeed
	}
	speeds := []float64{startSpeed}
	if startSpeed > 5 {
		speeds = append(speeds, 5)
	}
	speeds = append(speeds, 1, 1)

	return speeds
}

func LoadTests(testDir string) ([]*Test, error) {
	paths, err := filepath.Glob(filepath.Join(testDir, "/*/test.json"))
	if err != nil {
		return nil, err
	}

	tests := make([]*Test, len(paths))

	for i, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		test := &Test{}

		err = json.Unmarshal(data, test)
		if err != nil {
			return nil, err
		}

		test.Name = strings.TrimPrefix(filepath.Dir(path), testDir+"/")

		tests[i] = test
	}

	return tests, nil
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
		`status`,                         // file tree
		`log --pretty=%B -p -1`,          // log
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

func generateSnapshots(actualDir string, expectedDir string) (string, string, error) {
	actual, err := generateSnapshot(actualDir)
	if err != nil {
		return "", "", err
	}

	// there are a couple of reasons we're not generating the snapshot in expectedDir directly:
	// Firstly we don't want to have to revert our .git file back to .git_keep.
	// Secondly, the act of calling git commands like 'git status' actually changes the index
	// for some reason, and we don't want to leave your lazygit working tree dirty as a result.
	expectedDirCopyDir := filepath.Join(filepath.Dir(expectedDir), "expected_dir_test")
	err = oscommands.CopyDir(expectedDir, expectedDirCopyDir)
	if err != nil {
		return "", "", err
	}

	if err := restoreGitDirs(expectedDirCopyDir); err != nil {
		return "", "", err
	}

	expected, err := generateSnapshot(expectedDirCopyDir)
	if err != nil {
		return "", "", err
	}

	err = os.RemoveAll(expectedDirCopyDir)
	if err != nil {
		return "", "", err
	}

	return actual, expected, nil
}

func getPathsToRename(dir string, needle string) []string {
	pathsToRename := []string{}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Name() == needle {
			pathsToRename = append(pathsToRename, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return pathsToRename
}

// Git refuses to track .git and .gitmodules folders in subdirectories so we need to rename it
// to git_keep after running a test, and then change it back again
var untrackedGitDirs []string = []string{".git", ".gitmodules"}

func renameGitDirs(dir string) error {
	for _, untrackedGitDir := range untrackedGitDirs {
		for _, path := range getPathsToRename(dir, untrackedGitDir) {
			err := os.Rename(path, path+"_keep")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func restoreGitDirs(dir string) error {
	for _, untrackedGitDir := range untrackedGitDirs {
		for _, path := range getPathsToRename(dir, untrackedGitDir+"_keep") {
			err := os.Rename(path, strings.TrimSuffix(path, "_keep"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateSnapshotsForRemote(actualDir string, expectedDir string) (string, string, error) {
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

func getLazygitCommand(testPath string, rootDir string, mode Mode, speed float64, extraCmdArgs string) (*exec.Cmd, error) {
	osCommand := oscommands.NewDummyOSCommand()

	replayPath := filepath.Join(testPath, "recording.json")
	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")
	actualDir := filepath.Join(testPath, "actual")

	exists, err := osCommand.FileExists(filepath.Join(testPath, "config"))
	if err != nil {
		return nil, err
	}

	if exists {
		templateConfigDir = filepath.Join(testPath, "config")
	}

	configDir := filepath.Join(testPath, "used_config")

	err = os.RemoveAll(configDir)
	if err != nil {
		return nil, err
	}
	err = oscommands.CopyDir(templateConfigDir, configDir)
	if err != nil {
		return nil, err
	}

	cmdStr := fmt.Sprintf("%s -debug --use-config-dir=%s --path=%s %s", tempLazygitPath(), configDir, actualDir, extraCmdArgs)

	cmdObj := osCommand.Cmd.New(cmdStr)
	cmdObj.AddEnvVars(fmt.Sprintf("SPEED=%f", speed))

	switch mode {
	case RECORD:
		cmdObj.AddEnvVars(fmt.Sprintf("RECORD_EVENTS_TO=%s", replayPath))
	case TEST, UPDATE_SNAPSHOT:
		cmdObj.AddEnvVars(fmt.Sprintf("REPLAY_EVENTS_FROM=%s", replayPath))
	}

	return cmdObj.GetCmd(), nil
}

func folderExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
