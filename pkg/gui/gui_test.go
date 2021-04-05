package gui

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/creack/pty"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/stretchr/testify/assert"
)

// This file is quite similar to integration/main.go. The main difference is that this file is
// run via `go test` whereas the other is run via `test/lazyintegration/main.go` which provides
//  a convenient gui wrapper around our integration tests. The `go test` approach is better
// for CI and for running locally in the background to ensure you haven't broken
// anything while making changes. If you want to visually see what's happening when a test is run,
// you'll need to take the other approach
//
// As for this file, to run an integration test, e.g. for test 'commit', go:
// go test pkg/gui/gui_test.go -run /commit
//
// To update a snapshot for an integration test, pass UPDATE_SNAPSHOTS=true
// UPDATE_SNAPSHOTS=true go test pkg/gui/gui_test.go -run /commit
//
// integration tests are run in test/integration/<test_name>/actual and the final test does
// not clean up that directory so you can cd into it to see for yourself what
// happened when a test fails.
//
// To override speed, pass e.g. `SPEED=1` as an env var. Otherwise we start each test
// at a high speed and then drop down to lower speeds upon each failure until finally
// trying at the original playback speed (speed 1). A speed of 2 represents twice the
// original playback speed. Speed may be a decimal.

func Test(t *testing.T) {
	rootDir := integration.GetRootDirectory()
	err := os.Chdir(rootDir)
	assert.NoError(t, err)

	testDir := filepath.Join(rootDir, "test", "integration")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.RunCommand("go build -o %s", integration.TempLazygitPath())
	assert.NoError(t, err)

	tests, err := integration.LoadTests(testDir)
	assert.NoError(t, err)

	record := false
	updateSnapshots := record || os.Getenv("UPDATE_SNAPSHOTS") != ""

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			speeds := integration.GetTestSpeeds(test.Speed, updateSnapshots)
			testPath := filepath.Join(testDir, test.Name)
			actualDir := filepath.Join(testPath, "actual")
			expectedDir := filepath.Join(testPath, "expected")
			t.Logf("testPath: %s, actualDir: %s, expectedDir: %s", testPath, actualDir, expectedDir)

			// three retries at normal speed for the sake of flakey tests
			speeds = append(speeds, 1, 1, 1)
			for i, speed := range speeds {
				t.Logf("%s: attempting test at speed %f\n", test.Name, speed)

				integration.FindOrCreateDir(testPath)
				integration.PrepareIntegrationTestDir(actualDir)
				err := integration.CreateFixture(testPath, actualDir)
				assert.NoError(t, err)

				configDir := filepath.Join(testPath, "used_config")

				cmd, err := integration.GetLazygitCommand(testPath, rootDir, record, speed)
				assert.NoError(t, err)

				cmd.Env = append(
					cmd.Env,
					"HEADLESS=true",
					"TERM=xterm",
				)

				f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 100, Cols: 100})
				assert.NoError(t, err)

				_, _ = io.Copy(ioutil.Discard, f)

				assert.NoError(t, err)

				_ = f.Close()

				if updateSnapshots {
					err = oscommands.CopyDir(actualDir, expectedDir)
					assert.NoError(t, err)
				}

				actual, expected, err := integration.GenerateSnapshots(actualDir, expectedDir)
				assert.NoError(t, err)

				if expected == actual {
					t.Logf("%s: success at speed %f\n", test.Name, speed)
					break
				}

				// if the snapshots and we haven't tried all playback speeds different we'll retry at a slower speed
				if i == len(speeds)-1 {
					// get the log file and print that
					bytes, err := ioutil.ReadFile(filepath.Join(configDir, "development.log"))
					assert.NoError(t, err)
					t.Log(string(bytes))
					assert.Equal(t, expected, actual, fmt.Sprintf("expected:\n%s\nactual:\n%s\n", expected, actual))
				}
			}
		})
	}
}
