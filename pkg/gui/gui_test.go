package gui

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
// TODO: support passing an env var for playback speed, given it's currently pretty fast

type integrationTest struct {
	name    string
	fixture string
}

func tests() []integrationTest {
	return []integrationTest{
		{
			name:    "commit",
			fixture: "newFile",
		},
		{
			name:    "squash",
			fixture: "manyCommits",
		},
		{
			name:    "patchBuilding",
			fixture: "updatedFile",
		},
		{
			name:    "patchBuilding2",
			fixture: "updatedFile",
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

func generateSnapshot(t *testing.T) string {
	osCommand := oscommands.NewDummyOSCommand()
	cmd := `bash -c "git status; cat ./*; git log --pretty=%B -p"`

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

func Test(t *testing.T) {
	tests := tests()

	gotoRootDirectory()

	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	record := os.Getenv("RECORD_EVENTS") != ""
	updateSnapshots := record || os.Getenv("UPDATE_SNAPSHOTS") != ""

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			speeds := []int{10, 5, 1}
			if updateSnapshots {
				// have to go at original speed if updating snapshots in case we go to fast and create a junk snapshot
				speeds = []int{1}
			}

			for i, speed := range speeds {
				t.Logf("%s: attempting test at speed %d\n", test.name, speed)

				testPath := filepath.Join(rootDir, "test", "integration", test.name)
				findOrCreateDir(testPath)

				snapshotPath := filepath.Join(testPath, "snapshot.txt")

				err := os.Chdir(rootDir)
				assert.NoError(t, err)

				prepareIntegrationTestDir()

				err = createFixture(rootDir, test.fixture)
				assert.NoError(t, err)

				runLazygit(t, testPath, rootDir, record, speed)

				actual := generateSnapshot(t)

				if updateSnapshots {
					err := ioutil.WriteFile(snapshotPath, []byte(actual), 0600)
					assert.NoError(t, err)
				}

				expectedBytes, err := ioutil.ReadFile(snapshotPath)
				assert.NoError(t, err)
				expected := string(expectedBytes)

				if expected == actual {
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

func createFixture(rootDir string, name string) error {
	osCommand := oscommands.NewDummyOSCommand()
	cmd := exec.Command("bash", filepath.Join(rootDir, "test", "fixtures", fmt.Sprintf("%s.sh", name)))

	if err := osCommand.RunExecutable(cmd); err != nil {
		return err
	}

	return nil
}

func gotoRootDirectory() {
	for {
		_, err := os.Stat(".git")

		if err == nil {
			return
		}

		if !os.IsNotExist(err) {
			panic(err)
		}

		if err = os.Chdir(".."); err != nil {
			panic(err)
		}
	}
}

func runLazygit(t *testing.T, testPath string, rootDir string, record bool, speed int) {
	osCommand := oscommands.NewDummyOSCommand()

	replayPath := filepath.Join(testPath, "recording.json")
	cmdStr := fmt.Sprintf("go run %s", filepath.Join(rootDir, "main.go"))
	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")

	exists, err := osCommand.FileExists(filepath.Join(testPath, "config"))
	assert.NoError(t, err)

	if exists {
		templateConfigDir = filepath.Join(testPath, "config")
	}

	configDir := filepath.Join(rootDir, "test", "integration_test_config")

	err = os.RemoveAll(configDir)
	assert.NoError(t, err)
	err = oscommands.CopyDir(templateConfigDir, configDir)
	assert.NoError(t, err)

	cmdStr = fmt.Sprintf("%s --use-config-dir=%s", cmdStr, configDir)

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
	if os.Getenv("TERM") == "" {
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

func prepareIntegrationTestDir() {
	path := filepath.Join("test", "integration_test")

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

	if err := os.Chdir(path); err != nil {
		panic(err)
	}
}
