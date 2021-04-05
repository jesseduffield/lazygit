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

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
)

type Test struct {
	Name        string `json:"name"`
	Speed       int    `json:"speed"`
	Description string `json:"description"`
}

func PrepareIntegrationTestDir(actualDir string) {
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

func CreateFixture(testPath, actualDir string) error {
	osCommand := oscommands.NewDummyOSCommand()
	bashScriptPath := filepath.Join(testPath, "setup.sh")
	cmd := secureexec.Command("bash", bashScriptPath, actualDir)

	if err := osCommand.RunExecutable(cmd); err != nil {
		return err
	}

	return nil
}

func TempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
}

func GetTestSpeeds(testStartSpeed int, updateSnapshots bool) []int {
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

func FindOrCreateDir(path string) {
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

func GenerateSnapshot(dir string) (string, error) {
	osCommand := oscommands.NewDummyOSCommand()

	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err != nil {
		return "git directory not found", nil
	}

	snapshot := ""

	statusCmd := fmt.Sprintf(`git -C %s status`, dir)
	statusCmdOutput, err := osCommand.RunCommandWithOutput(statusCmd)
	if err != nil {
		return "", err
	}

	snapshot += statusCmdOutput + "\n"

	logCmd := fmt.Sprintf(`git -C %s log --pretty=%%B -p -1`, dir)
	logCmdOutput, err := osCommand.RunCommandWithOutput(logCmd)
	if err != nil {
		return "", err
	}

	snapshot += logCmdOutput + "\n"

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

func GenerateSnapshots(actualDir string, expectedDir string) (string, string, error) {
	actual, err := GenerateSnapshot(actualDir)
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

	expected, err := GenerateSnapshot(expectedDir)
	if err != nil {
		return "", "", err
	}

	return actual, expected, nil
}

func GetLazygitCommand(testPath string, rootDir string, record bool, speed int) (*exec.Cmd, error) {
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

	cmdStr := fmt.Sprintf("%s -debug --use-config-dir=%s --path=%s", TempLazygitPath(), configDir, actualDir)

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

	return cmd, nil
}
