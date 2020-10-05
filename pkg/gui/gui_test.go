package gui

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	name       string
	fixture    string
	startSpeed int
}

func tests() []integrationTest {
	return []integrationTest{
		{
			name:       "commit",
			fixture:    "newFile",
			startSpeed: 20,
		},
		{
			name:       "squash",
			fixture:    "manyCommits",
			startSpeed: 6,
		},
		{
			name:       "patchBuilding",
			fixture:    "updatedFile",
			startSpeed: 3,
		},
		{
			name:       "patchBuilding2",
			fixture:    "updatedFile",
			startSpeed: 3,
		},
		{
			name:    "mergeConflicts",
			fixture: "mergeConflicts",
		},
		{
			name:    "searching",
			fixture: "newFile",
		},
		{
			name:    "searchingInStagingPanel",
			fixture: "newFile2",
		},
	}
}

func generateSnapshot(t *testing.T, actualDir string) string {
	osCommand := oscommands.NewDummyOSCommand()
	cmd := fmt.Sprintf(`bash -c "cd %s && git status; cat ./*; git log --pretty=%%B -p"`, actualDir)

	// need to copy from current directory to

	snapshot, err := osCommand.RunCommandWithOutput(cmd)
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

func Test(t *testing.T) {
	tests := tests()

	rootDir := getRootDirectory()

	record := os.Getenv("RECORD_EVENTS") != ""
	updateSnapshots := record || os.Getenv("UPDATE_SNAPSHOTS") != ""

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			if runInParallel() {
				t.Parallel()
			}

			speeds := getTestSpeeds(test.startSpeed, updateSnapshots)

			for i, speed := range speeds {
				t.Logf("%s: attempting test at speed %d\n", test.name, speed)

				testPath := filepath.Join(rootDir, "test", "integration", test.name)
				actualDir := filepath.Join(testPath, "actual")
				expectedDir := filepath.Join(testPath, "expected")
				findOrCreateDir(testPath)

				prepareIntegrationTestDir(testPath)

				err := createFixture(rootDir, test.fixture, actualDir)
				assert.NoError(t, err)

				runLazygit(t, testPath, rootDir, record, speed)

				actual := generateSnapshot(t, actualDir)

				if updateSnapshots {
					err = oscommands.CopyDir(actualDir, expectedDir)
					assert.NoError(t, err)
				}

				expected := generateSnapshot(t, expectedDir)

				if expected == actual {
					t.Logf("%s: success at speed %d\n", test.name, speed)
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

func createFixture(rootDir string, name string, actualDir string) error {
	osCommand := oscommands.NewDummyOSCommand()
	cmd := exec.Command("bash", filepath.Join(rootDir, "test", "fixtures", fmt.Sprintf("%s.sh", name)), actualDir)

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
	cmdStr := fmt.Sprintf("go run %s", filepath.Join(rootDir, "main.go"))
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

	cmdStr = fmt.Sprintf("%s --use-config-dir=%s --path=%s", cmdStr, configDir, actualDir)

	cmd := osCommand.ExecutableFromString(cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("REPLAY_SPEED=%d", speed))

	if record {
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
	if runInParallel() {
		cmd.Env = append(cmd.Env, "TERM=xterm")

		f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 100, Cols: 100})
		assert.NoError(t, err)

		_, _ = io.Copy(ioutil.Discard, f)

		assert.NoError(t, err)

		_ = f.Close()
	} else {
		err := osCommand.RunExecutable(cmd)
		assert.NoError(t, err)
	}
}

func runInParallel() bool {
	return os.Getenv("PARALLEL") != "" || os.Getenv("TERM") == ""
}

func prepareIntegrationTestDir(testPath string) {
	path := filepath.Join(testPath, "actual")

	// remove contents of integration test directory
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	for _, d := range dir {
		os.RemoveAll(filepath.Join(path, d.Name()))
	}
}
