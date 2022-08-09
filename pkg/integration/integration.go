package integration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/integration/helpers"
	"github.com/jesseduffield/lazygit/pkg/integration/integration_tests"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

// this is the integration runner for the new and improved integration interface

func RunTestsNew(
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	fnWrapper func(test types.Test, f func() error),
	mode Mode,
	onFail func(expected string, actual string, prefix string),
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

	for _, test := range integration_tests.Tests {
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

			err = runLazygit(test, testPath, rootDir)
			if err != nil {
				return err
			}

			if mode == UPDATE_SNAPSHOT {
				// create/update snapshot
				err = oscommands.CopyDir(actualDir, expectedDir)
				if err != nil {
					return err
				}

				if err := renameSpecialPaths(expectedDir); err != nil {
					return err
				}

				logf("%s", "updated snapshot")
			} else {
				if err := validateSameRepos(expectedDir, actualDir); err != nil {
					return err
				}

				// iterate through each repo in the expected dir and comparet to the corresponding repo in the actual dir
				expectedFiles, err := ioutil.ReadDir(expectedDir)
				if err != nil {
					return err
				}

				for _, f := range expectedFiles {
					if !f.IsDir() {
						return errors.New("unexpected file (as opposed to directory) in integration test 'expected' directory")
					}

					// get corresponding file name from actual dir
					actualRepoPath := filepath.Join(actualDir, f.Name())
					expectedRepoPath := filepath.Join(expectedDir, f.Name())

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

						onFail(expectedRepo, actualRepo, f.Name())
					}
				}
			}

			logf("test passed: %s", test.Name())

			return nil
		})
	}

	return nil
}

func createFixtureNew(test types.Test, actualDir string, rootDir string) error {
	if err := os.Chdir(actualDir); err != nil {
		panic(err)
	}

	shell := &helpers.ShellImpl{}
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

func runLazygit(test types.Test, testPath string, rootDir string) error {
	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")
	actualRepoDir := filepath.Join(testPath, "actual", "repo")

	configDir := filepath.Join(testPath, "used_config")

	err := os.RemoveAll(configDir)
	if err != nil {
		return err
	}
	err = oscommands.CopyDir(templateConfigDir, configDir)
	if err != nil {
		return err
	}

	// TODO: support test.ExtraCmdArgs in some form.
	cliArgs := &app.CliArgs{
		Debug:        true,
		UseConfigDir: configDir,
		RepoPath:     actualRepoDir,
	}

	buildInfo := &app.BuildInfo{
		Commit:      "1234abc",
		Date:        "2020-01-01",
		Version:     "1.0.0",
		BuildSource: "unknown",
	}

	return convertPanicToError(func() { app.Start(cliArgs, buildInfo, test) })
}

func convertPanicToError(f func()) (err error) { //nolint: nakedret
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Lazygit panicked. Stacktrace:: \n" + string(debug.Stack()))
		}
	}()

	f()

	return nil
}
