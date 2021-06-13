package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// this function is used by both `go test` and from our lazyintegration gui, but
// errors need to be handled differently in each (for example go test is always
// working with *testing.T) so we pass in any differences as args here.
func RunTests(
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	fnWrapper func(test *Test, f func(*testing.T) error),
	updateSnapshots bool,
	record bool,
	speedEnv string,
	onFail func(t *testing.T, expected string, actual string),
	includeSkipped bool,
) error {
	rootDir := GetRootDirectory()
	err := os.Chdir(rootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(rootDir, "test", "integration")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.RunExecutable(
		oscommands.NewCmdObjFromStr(
			fmt.Sprintf("go build -o %s", tempLazygitPath()),
		),
	)
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

		fnWrapper(test, func(t *testing.T) error {
			speeds := getTestSpeeds(test.Speed, updateSnapshots, speedEnv)
			testPath := filepath.Join(testDir, test.Name)
			actualDir := filepath.Join(testPath, "actual")
			expectedDir := filepath.Join(testPath, "expected")
			logf("path: %s", testPath)

			// three retries at normal speed for the sake of flakey tests
			speeds = append(speeds, 1)
			for i, speed := range speeds {
				logf("%s: attempting test at speed %f\n", test.Name, speed)

				findOrCreateDir(testPath)
				prepareIntegrationTestDir(actualDir)
				err := createFixture(testPath, actualDir)
				if err != nil {
					return err
				}

				configDir := filepath.Join(testPath, "used_config")

				cmd, err := getLazygitCommand(testPath, rootDir, record, speed, test.ExtraCmdArgs)
				if err != nil {
					return err
				}

				err = runCmd(cmd)
				if err != nil {
					return err
				}

				if updateSnapshots {
					err = oscommands.CopyDir(actualDir, expectedDir)
					if err != nil {
						return err
					}
					err = os.Rename(
						filepath.Join(expectedDir, ".git"),
						filepath.Join(expectedDir, ".git_keep"),
					)
					if err != nil {
						return err
					}
				}

				actual, expected, err := generateSnapshots(actualDir, expectedDir)
				if err != nil {
					return err
				}

				if expected == actual {
					logf("%s: success at speed %f\n", test.Name, speed)
					break
				}

				// if the snapshots and we haven't tried all playback speeds different we'll retry at a slower speed
				if i == len(speeds)-1 {
					// get the log file and print that
					bytes, err := ioutil.ReadFile(filepath.Join(configDir, "development.log"))
					if err != nil {
						return err
					}
					logf("%s", string(bytes))
					onFail(t, expected, actual)
				}
			}

			return nil
		})
	}

	return nil
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
			panic("must run in lazygit folder or child folder")
		}
	}
}

func createFixture(testPath, actualDir string) error {
	osCommand := oscommands.NewDummyOSCommand()
	bashScriptPath := filepath.Join(testPath, "setup.sh")

	err := osCommand.RunExecutable(
		oscommands.NewCmdObj(
			secureexec.Command("bash", bashScriptPath, actualDir),
		),
	)

	if err != nil {
		return err
	}

	return nil
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
}

func getTestSpeeds(testStartSpeed float64, updateSnapshots bool, speedStr string) []float64 {
	if updateSnapshots {
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
	speeds = append(speeds, 1)

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
			err = os.MkdirAll(path, 0777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func generateSnapshot(dir string) (string, error) {
	osCommand := oscommands.NewDummyOSCommand()

	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return "git directory not found", nil
	}

	snapshot := ""

	cmdStrs := []string{
		fmt.Sprintf(`git -C %s status`, dir),                 // file tree
		fmt.Sprintf(`git -C %s log --pretty=%%B -p -1`, dir), // log
		fmt.Sprintf(`git -C %s tag -n`, dir),                 // tags
	}

	for _, cmdStr := range cmdStrs {
		// ignoring error for now. If there's an error it could be that there are no results
		output, _ := osCommand.RunCommandWithOutput(oscommands.NewCmdObjFromStr(cmdStr))

		snapshot += output + "\n"
	}

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
		snapshot += string(bytes) + "\n"

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

	// git refuses to track .git folders in subdirectories so we need to rename it
	// to git_keep after running a test, and then change it back again
	defer func() {
		err = os.Rename(
			filepath.Join(expectedDir, ".git"),
			filepath.Join(expectedDir, ".git_keep"),
		)

		if err != nil {
			panic(err)
		}
	}()

	// ignoring this error because we might not have a .git_keep file here yet.
	_ = os.Rename(
		filepath.Join(expectedDir, ".git_keep"),
		filepath.Join(expectedDir, ".git"),
	)

	expected, err := generateSnapshot(expectedDir)
	if err != nil {
		return "", "", err
	}

	return actual, expected, nil
}

func getLazygitCommand(testPath string, rootDir string, record bool, speed float64, extraCmdArgs string) (*exec.Cmd, error) {
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

	cmd := osCommand.ExecutableFromString(cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("SPEED=%f", speed))

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

	return cmd, nil
}
