package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/creack/pty"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

// To run an integration test, e.g. for test 'commit', go:
// go test pkg/gui/gui_test.go -run /commit
//
// To record keypresses for an integration test, pass RECORD_EVENTS=true like so:
// RECORD_EVENTS=true go test pkg/gui/gui_test.go -run /commit
//
// To update a snapshot for an integration test, pass UPDATE_SNAPSHOTS=true
// UPDATE_SNAPSHOTS=true go test pkg/gui/gui_test.go -run /commit
//
// When RECORD_EVENTS is true, updates will be updated automatically
//
// integration tests are run in test/integration_test and the final test does
// not clean up that directory so you can cd into it to see for yourself what
// happened when a test failed.
//
// To run tests in parallel pass `PARALLEL=true` as an env var. Tests are run in parallel
// on CI, and are run in a pty so you won't be able to see the stdout of the program
//
// To override speed, pass e.g. `SPEED=1` as an env var. Otherwise we start each test
// at a high speed and then drop down to lower speeds upon each failure until finally
// trying at the original playback speed (speed 1). A speed of 2 represents twice the
// original playback speed. Speed must be an integer.

type integrationTest struct {
	Name        string `json:"name"`
	Speed       int    `json:"speed"`
	Description string `json:"description"`
}

func loadTests(t *testing.T, testDir string) []*integrationTest {
	paths, err := filepath.Glob(filepath.Join(testDir, "/*/test.json"))
	assert.NoError(t, err)

	tests := make([]*integrationTest, len(paths))

	for i, path := range paths {
		data, err := ioutil.ReadFile(path)
		assert.NoError(t, err)

		test := &integrationTest{}

		err = json.Unmarshal(data, test)
		assert.NoError(t, err)

		test.Name = strings.TrimPrefix(filepath.Dir(path), testDir+"/")

		tests[i] = test
	}

	return tests
}

func generateSnapshot(t *testing.T, dir string) string {
	osCommand := oscommands.NewDummyOSCommand()

	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return "git directory not found"
	}

	snapshot := ""

	statusCmd := fmt.Sprintf(`git -C %s status`, dir)
	statusCmdOutput, err := osCommand.RunCommandWithOutput(statusCmd)
	assert.NoError(t, err)

	snapshot += statusCmdOutput + "\n"

	logCmd := fmt.Sprintf(`git -C %s log --pretty=%%B -p -1`, dir)
	logCmdOutput, err := osCommand.RunCommandWithOutput(logCmd)
	assert.NoError(t, err)

	snapshot += logCmdOutput + "\n"

	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		assert.NoError(t, err)

		if f.IsDir() {
			if f.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		bytes, err := ioutil.ReadFile(path)
		assert.NoError(t, err)
		snapshot += string(bytes) + "\n"

		return nil
	})

	assert.NoError(t, err)

	return snapshot
}

func findOrCreateDir(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func getTestSpeeds(testStartSpeed int, updateSnapshots bool) []int {
	if updateSnapshots {
		// have to go at original speed if updating snapshots in case we go to fast and create a junk snapshot
		return []int{1}
	}

	speedEnv := os.Getenv("SPEED")
	if speedEnv != "" {
		speed, err := strconv.Atoi(speedEnv)
		if err != nil {
			panic(err)
		}
		return []int{speed}
	}

	// default is 10, 5, 1
	startSpeed := 10
	if testStartSpeed != 0 {
		startSpeed = testStartSpeed
	}
	speeds := []int{startSpeed}
	if startSpeed > 5 {
		speeds = append(speeds, 5)
	}
	speeds = append(speeds, 1)

	return speeds
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
}

func Test(t *testing.T) {
	rootDir := getRootDirectory()
	err := os.Chdir(rootDir)
	assert.NoError(t, err)

	testDir := filepath.Join(rootDir, "test", "integration")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.RunCommand("go build -o %s", tempLazygitPath())
	assert.NoError(t, err)

	tests := loadTests(t, testDir)

	record := os.Getenv("RECORD_EVENTS") != ""
	updateSnapshots := record || os.Getenv("UPDATE_SNAPSHOTS") != ""

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			if runInParallel() {
				t.Parallel()
			}

			speeds := getTestSpeeds(test.Speed, updateSnapshots)

			for i, speed := range speeds {
				t.Logf("%s: attempting test at speed %d\n", test.Name, speed)

				testPath := filepath.Join(testDir, test.Name)
				actualDir := filepath.Join(testPath, "actual")
				expectedDir := filepath.Join(testPath, "expected")
				t.Logf("testPath: %s, actualDir: %s, expectedDir: %s", testPath, actualDir, expectedDir)
				findOrCreateDir(testPath)

				prepareIntegrationTestDir(actualDir)

				err := createFixture(testPath, actualDir)
				assert.NoError(t, err)

				runLazygit(t, testPath, rootDir, record, speed)

				if updateSnapshots {
					err = oscommands.CopyDir(actualDir, expectedDir)
					assert.NoError(t, err)
				}

				actual := generateSnapshot(t, actualDir)

				expected := ""

				func() {
					// git refuses to track .git folders in subdirectories so we need to rename it
					// to git_keep after running a test

					defer func() {
						err = os.Rename(
							filepath.Join(expectedDir, ".git"),
							filepath.Join(expectedDir, ".git_keep"),
						)

						assert.NoError(t, err)
					}()

					// ignoring this error because we might not have a .git_keep file here yet.
					_ = os.Rename(
						filepath.Join(expectedDir, ".git_keep"),
						filepath.Join(expectedDir, ".git"),
					)

					expected = generateSnapshot(t, expectedDir)
				}()

				if expected == actual {
					t.Logf("%s: success at speed %d\n", test.Name, speed)
					break
				}

				// if the snapshots and we haven't tried all playback speeds different we'll retry at a slower speed
				if i == len(speeds)-1 {
					assert.Equal(t, expected, actual, fmt.Sprintf("expected:\n%s\nactual:\n%s\n", expected, actual))
				}
			}
		})
	}
}

func createFixture(testPath, actualDir string) error {
	osCommand := oscommands.NewDummyOSCommand()
	bashScriptPath := filepath.Join(testPath, "setup.sh")
	cmd := exec.Command("bash", bashScriptPath, actualDir)

	if err := osCommand.RunExecutable(cmd); err != nil {
		return err
	}

	return nil
}

func getRootDirectory() string {
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
			panic("must run in lazygit folder or child folder")
		}
	}
}

func runLazygit(t *testing.T, testPath string, rootDir string, record bool, speed int) {
	osCommand := oscommands.NewDummyOSCommand()

	replayPath := filepath.Join(testPath, "recording.json")
	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")
	actualDir := filepath.Join(testPath, "actual")

	exists, err := osCommand.FileExists(filepath.Join(testPath, "config"))
	assert.NoError(t, err)

	if exists {
		templateConfigDir = filepath.Join(testPath, "config")
	}

	configDir := filepath.Join(testPath, "used_config")

	err = os.RemoveAll(configDir)
	assert.NoError(t, err)
	err = oscommands.CopyDir(templateConfigDir, configDir)
	assert.NoError(t, err)

	cmdStr := fmt.Sprintf("%s --use-config-dir=%s --path=%s", tempLazygitPath(), configDir, actualDir)

	cmd := osCommand.ExecutableFromString(cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("REPLAY_SPEED=%d", speed))

	if record {
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("RECORD_EVENTS_TO=%s", replayPath),
		)
	} else {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("REPLAY_EVENTS_FROM=%s", replayPath),
		)
	}

	// if we're on CI we'll need to use a PTY. We can work that out by seeing if the 'TERM' env is defined.
	if runInPTY() {
		cmd.Env = append(cmd.Env, "TERM=xterm")

		f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 100, Cols: 100})
		assert.NoError(t, err)

		_, _ = io.Copy(ioutil.Discard, f)

		assert.NoError(t, err)

		_ = f.Close()
	} else {
		err := cmd.Run()
		assert.NoError(t, err)
	}
}

func runInParallel() bool {
	return os.Getenv("PARALLEL") != ""
}

func runInPTY() bool {
	return runInParallel() || os.Getenv("TERM") == ""
}

func prepareIntegrationTestDir(actualDir string) {
	// remove contents of integration test directory
	dir, err := ioutil.ReadDir(actualDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(actualDir, 0777)
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
