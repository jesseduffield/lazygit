package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// this is the integration runner for the new and improved integration interface

const (
	TEST_NAME_ENV_VAR = "TEST_NAME"
	SANDBOX_ENV_VAR   = "SANDBOX"
)

type Mode int

const (
	// Default: if a snapshot test fails, the we'll be asked whether we want to update it
	ASK_TO_UPDATE_SNAPSHOT Mode = iota
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

func RunTests(
	tests []*IntegrationTest,
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	testWrapper func(test *IntegrationTest, f func() error),
	mode Mode,
	keyPressDelay int,
	maxAttempts int,
) error {
	projectRootDir := utils.GetLazyRootDirectory()
	err := os.Chdir(projectRootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(projectRootDir, "test", "integration_new")

	if err := buildLazygit(); err != nil {
		return err
	}

	for _, test := range tests {
		test := test

		testWrapper(test, func() error { //nolint: thelper
			paths := NewPaths(
				filepath.Join(testDir, test.Name()),
			)

			for i := 0; i < maxAttempts; i++ {
				err := runTest(test, paths, projectRootDir, logf, runCmd, mode, keyPressDelay)
				if err != nil {
					if i == maxAttempts-1 {
						return err
					}
					logf("retrying test %s", test.Name())
				} else {
					break
				}
			}

			return nil
		})
	}

	return nil
}

func runTest(
	test *IntegrationTest,
	paths Paths,
	projectRootDir string,
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	mode Mode,
	keyPressDelay int,
) error {
	if test.Skip() {
		logf("Skipping test %s", test.Name())
		return nil
	}

	logf("path: %s", paths.Root())

	if err := prepareTestDir(test, paths); err != nil {
		return err
	}

	cmd, err := getLazygitCommand(test, paths, projectRootDir, mode, keyPressDelay)
	if err != nil {
		return err
	}

	err = runCmd(cmd)
	if err != nil {
		return err
	}

	return HandleSnapshots(paths, logf, test, mode)
}

func prepareTestDir(
	test *IntegrationTest,
	paths Paths,
) error {
	findOrCreateDir(paths.Root())
	deleteAndRecreateEmptyDir(paths.Actual())

	err := os.Mkdir(paths.ActualRepo(), 0o777)
	if err != nil {
		return err
	}

	return createFixture(test, paths)
}

func buildLazygit() error {
	osCommand := oscommands.NewDummyOSCommand()
	return osCommand.Cmd.New(fmt.Sprintf(
		"go build -o %s pkg/integration/clients/injector/main.go", tempLazygitPath(),
	)).Run()
}

func createFixture(test *IntegrationTest, paths Paths) error {
	shell := NewShell(paths.ActualRepo())
	shell.RunCommand("git init -b master")
	shell.RunCommand(`git config user.email "CI@example.com"`)
	shell.RunCommand(`git config user.name "CI"`)
	shell.RunCommand(`git config commit.gpgSign false`)
	shell.RunCommand(`git config protocol.file.allow always`)

	test.SetupRepo(shell)

	return nil
}

func getLazygitCommand(test *IntegrationTest, paths Paths, rootDir string, mode Mode, keyPressDelay int) (*exec.Cmd, error) {
	osCommand := oscommands.NewDummyOSCommand()

	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")

	err := os.RemoveAll(paths.Config())
	if err != nil {
		return nil, err
	}
	err = oscommands.CopyDir(templateConfigDir, paths.Config())
	if err != nil {
		return nil, err
	}

	cmdStr := fmt.Sprintf("%s -debug --use-config-dir=%s --path=%s %s", tempLazygitPath(), paths.Config(), paths.ActualRepo(), test.ExtraCmdArgs())

	cmdObj := osCommand.Cmd.New(cmdStr)

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", TEST_NAME_ENV_VAR, test.Name()))
	if mode == SANDBOX {
		cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", "SANDBOX", "true"))
	}

	if keyPressDelay > 0 {
		cmdObj.AddEnvVars(fmt.Sprintf("KEY_PRESS_DELAY=%d", keyPressDelay))
	}

	return cmdObj.GetCmd(), nil
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
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

func deleteAndRecreateEmptyDir(path string) {
	// remove contents of integration test directory
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0o777)
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
