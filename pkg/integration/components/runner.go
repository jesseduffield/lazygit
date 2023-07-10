package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

const (
	TEST_NAME_ENV_VAR         = "TEST_NAME"
	SANDBOX_ENV_VAR           = "SANDBOX"
	GIT_CONFIG_GLOBAL_ENV_VAR = "GIT_CONFIG_GLOBAL"
)

// This function lets you run tests either from within `go test` or from a regular binary.
// The reason for having two separate ways of testing is that `go test` isn't great at
// showing what's actually happening during the test, but it's still good at running
// tests in telling you about their results.
func RunTests(
	tests []*IntegrationTest,
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) error,
	testWrapper func(test *IntegrationTest, f func() error),
	sandbox bool,
	keyPressDelay int,
	maxAttempts int,
) error {
	projectRootDir := utils.GetLazyRootDirectory()
	err := os.Chdir(projectRootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(projectRootDir, "test", "results")

	if err := buildLazygit(); err != nil {
		return err
	}

	gitVersion, err := getGitVersion()
	if err != nil {
		return err
	}

	for _, test := range tests {
		test := test

		testWrapper(test, func() error { //nolint: thelper
			paths := NewPaths(
				filepath.Join(testDir, test.Name()),
			)

			for i := 0; i < maxAttempts; i++ {
				err := runTest(test, paths, projectRootDir, logf, runCmd, sandbox, keyPressDelay, gitVersion)
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
	sandbox bool,
	keyPressDelay int,
	gitVersion *git_commands.GitVersion,
) error {
	if test.Skip() {
		logf("Skipping test %s", test.Name())
		return nil
	}

	if !test.ShouldRunForGitVersion(gitVersion) {
		logf("Skipping test %s for git version %d.%d.%d", test.Name(), gitVersion.Major, gitVersion.Minor, gitVersion.Patch)
		return nil
	}

	if err := prepareTestDir(test, paths, projectRootDir); err != nil {
		return err
	}

	cmd, err := getLazygitCommand(test, paths, projectRootDir, sandbox, keyPressDelay)
	if err != nil {
		return err
	}

	err = runCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

func prepareTestDir(
	test *IntegrationTest,
	paths Paths,
	rootDir string,
) error {
	findOrCreateDir(paths.Root())
	deleteAndRecreateEmptyDir(paths.Actual())

	err := os.Mkdir(paths.ActualRepo(), 0o777)
	if err != nil {
		return err
	}

	return createFixture(test, paths, rootDir)
}

func buildLazygit() error {
	// // TODO: remove this line!
	// // skipping this because I'm not making changes to the app code atm.
	// return nil

	osCommand := oscommands.NewDummyOSCommand()
	return osCommand.Cmd.New([]string{
		"go", "build", "-o", tempLazygitPath(), filepath.FromSlash("pkg/integration/clients/injector/main.go"),
	}).Run()
}

func createFixture(test *IntegrationTest, paths Paths, rootDir string) error {
	shell := NewShell(paths.ActualRepo(), func(errorMsg string) { panic(errorMsg) })
	shell.Init()

	os.Setenv(GIT_CONFIG_GLOBAL_ENV_VAR, globalGitConfigPath(rootDir))

	test.SetupRepo(shell)

	return nil
}

func globalGitConfigPath(rootDir string) string {
	return filepath.Join(rootDir, "test", "global_git_config")
}

func getGitVersion() (*git_commands.GitVersion, error) {
	osCommand := oscommands.NewDummyOSCommand()
	cmdObj := osCommand.Cmd.New([]string{"git", "--version"})
	versionStr, err := cmdObj.RunWithOutput()
	if err != nil {
		return nil, err
	}
	return git_commands.ParseGitVersion(versionStr)
}

func getLazygitCommand(test *IntegrationTest, paths Paths, rootDir string, sandbox bool, keyPressDelay int) (*exec.Cmd, error) {
	osCommand := oscommands.NewDummyOSCommand()

	err := os.RemoveAll(paths.Config())
	if err != nil {
		return nil, err
	}

	templateConfigDir := filepath.Join(rootDir, "test", "default_test_config")
	err = oscommands.CopyDir(templateConfigDir, paths.Config())
	if err != nil {
		return nil, err
	}

	cmdArgs := []string{tempLazygitPath(), "-debug", "--use-config-dir=" + paths.Config(), "--path=" + paths.ActualRepo()}
	cmdArgs = append(cmdArgs, test.ExtraCmdArgs()...)

	cmdObj := osCommand.Cmd.New(cmdArgs)

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", TEST_NAME_ENV_VAR, test.Name()))
	if sandbox {
		cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", SANDBOX_ENV_VAR, "true"))
	}
	if test.ExtraEnvVars() != nil {
		for key, value := range test.ExtraEnvVars() {
			cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", key, value))
		}
	}

	if keyPressDelay > 0 {
		cmdObj.AddEnvVars(fmt.Sprintf("KEY_PRESS_DELAY=%d", keyPressDelay))
	}

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", GIT_CONFIG_GLOBAL_ENV_VAR, globalGitConfigPath(rootDir)))

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
