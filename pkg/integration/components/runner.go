package components

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	lazycoreUtils "github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type RunTestArgs struct {
	Tests           []*IntegrationTest
	Logf            func(format string, formatArgs ...interface{})
	RunCmd          func(cmd *exec.Cmd) (int, error)
	TestWrapper     func(test *IntegrationTest, f func() error)
	Sandbox         bool
	WaitForDebugger bool
	RaceDetector    bool
	CodeCoverageDir string
	InputDelay      int
	MaxAttempts     int
}

// This function lets you run tests either from within `go test` or from a regular binary.
// The reason for having two separate ways of testing is that `go test` isn't great at
// showing what's actually happening during the test, but it's still good at running
// tests in telling you about their results.
func RunTests(args RunTestArgs) error {
	projectRootDir := lazycoreUtils.GetLazyRootDirectory()
	err := os.Chdir(projectRootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(projectRootDir, "test", "_results")
	if err := buildLazygit(args); err != nil {
		return err
	}

	gitVersion, err := getGitVersion()
	if err != nil {
		return err
	}

	for _, test := range args.Tests {
		test := test

		args.TestWrapper(test, func() error { //nolint: thelper
			paths := NewPaths(
				filepath.Join(testDir, test.Name()),
			)

			for i := 0; i < args.MaxAttempts; i++ {
				err := runTest(test, args, paths, projectRootDir, gitVersion)
				if err != nil {
					if i == args.MaxAttempts-1 {
						return err
					}
					args.Logf("retrying test %s", test.Name())
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
	args RunTestArgs,
	paths Paths,
	projectRootDir string,
	gitVersion *git_commands.GitVersion,
) error {
	if test.Skip() {
		args.Logf("Skipping test %s", test.Name())
		return nil
	}

	if !test.ShouldRunForGitVersion(gitVersion) {
		args.Logf("Skipping test %s for git version %d.%d.%d", test.Name(), gitVersion.Major, gitVersion.Minor, gitVersion.Patch)
		return nil
	}

	workingDir, err := prepareTestDir(test, paths, projectRootDir)
	if err != nil {
		return err
	}

	cmd, err := getLazygitCommand(test, args, paths, projectRootDir, workingDir)
	if err != nil {
		return err
	}

	pid, err := args.RunCmd(cmd)

	// Print race detector log regardless of the command's exit status
	if args.RaceDetector {
		logPath := fmt.Sprintf("%s.%d", raceDetectorLogsPath(), pid)
		if bytes, err := os.ReadFile(logPath); err == nil {
			args.Logf("Race detector log:\n" + string(bytes))
		}
	}

	return err
}

func prepareTestDir(
	test *IntegrationTest,
	paths Paths,
	rootDir string,
) (string, error) {
	findOrCreateDir(paths.Root())
	deleteAndRecreateEmptyDir(paths.Actual())

	err := os.Mkdir(paths.ActualRepo(), 0o777)
	if err != nil {
		return "", err
	}

	workingDir := createFixture(test, paths, rootDir)

	return workingDir, nil
}

func buildLazygit(testArgs RunTestArgs) error {
	args := []string{"go", "build"}
	if testArgs.WaitForDebugger {
		// Disable compiler optimizations (-N) and inlining (-l) because this
		// makes debugging work better
		args = append(args, "-gcflags=all=-N -l")
	}
	if testArgs.RaceDetector {
		args = append(args, "-race")
	}
	if testArgs.CodeCoverageDir != "" {
		args = append(args, "-cover")
	}
	args = append(args, "-o", tempLazygitPath(), filepath.FromSlash("pkg/integration/clients/injector/main.go"))
	osCommand := oscommands.NewDummyOSCommand()
	return osCommand.Cmd.New(args).Run()
}

// Sets up the fixture for test and returns the working directory to invoke
// lazygit in.
func createFixture(test *IntegrationTest, paths Paths, rootDir string) string {
	env := NewTestEnvironment(rootDir)

	env = append(env, fmt.Sprintf("%s=%s", PWD, paths.ActualRepo()))
	shell := NewShell(
		paths.ActualRepo(),
		env,
		func(errorMsg string) { panic(errorMsg) },
	)
	shell.Init()

	test.SetupRepo(shell)

	return shell.dir
}

func testPath(rootdir string) string {
	return filepath.Join(rootdir, "test")
}

func globalGitConfigPath(rootDir string) string {
	return filepath.Join(testPath(rootDir), "global_git_config")
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

func getLazygitCommand(
	test *IntegrationTest,
	args RunTestArgs,
	paths Paths,
	rootDir string,
	workingDir string,
) (*exec.Cmd, error) {
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

	cmdArgs := []string{tempLazygitPath(), "-debug", "--use-config-dir=" + paths.Config()}

	resolvedExtraArgs := lo.Map(test.ExtraCmdArgs(), func(arg string, _ int) string {
		return utils.ResolvePlaceholderString(arg, map[string]string{
			"actualPath":     paths.Actual(),
			"actualRepoPath": paths.ActualRepo(),
		})
	})
	cmdArgs = append(cmdArgs, resolvedExtraArgs...)

	// Use a limited environment for test isolation, including pass through
	// of just allowed host environment variables
	cmdObj := osCommand.Cmd.NewWithEnviron(cmdArgs, NewTestEnvironment(rootDir))

	// Integration tests related to symlink behavior need a PWD that
	// preserves symlinks. By default, SetWd will set a symlink-resolved
	// value for PWD. Here, we override that with the path (that may)
	// contain a symlink to simulate behavior in a user's shell correctly.
	cmdObj.SetWd(workingDir)
	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", PWD, workingDir))

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", LAZYGIT_ROOT_DIR, rootDir))

	if args.CodeCoverageDir != "" {
		// We set this explicitly here rather than inherit it from the test runner's
		// environment because the test runner has its own coverage directory that
		// it writes to and so if we pass GOCOVERDIR to that, it will be overwritten.
		cmdObj.AddEnvVars("GOCOVERDIR=" + args.CodeCoverageDir)
	}

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", TEST_NAME_ENV_VAR, test.Name()))
	if args.Sandbox {
		cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", SANDBOX_ENV_VAR, "true"))
	}
	if args.WaitForDebugger {
		cmdObj.AddEnvVars(fmt.Sprintf("%s=true", WAIT_FOR_DEBUGGER_ENV_VAR))
	}
	// Set a race detector log path only to avoid spamming the terminal with the
	// logs. We are not showing this anywhere yet.
	cmdObj.AddEnvVars(fmt.Sprintf("GORACE=log_path=%s", raceDetectorLogsPath()))
	if test.ExtraEnvVars() != nil {
		for key, value := range test.ExtraEnvVars() {
			cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", key, value))
		}
	}

	if args.InputDelay > 0 {
		cmdObj.AddEnvVars(fmt.Sprintf("INPUT_DELAY=%d", args.InputDelay))
	}

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", GIT_CONFIG_GLOBAL_ENV_VAR, globalGitConfigPath(rootDir)))

	return cmdObj.GetCmd(), nil
}

func tempLazygitPath() string {
	return filepath.Join("/tmp", "lazygit", "test_lazygit")
}

func raceDetectorLogsPath() string {
	return filepath.Join("/tmp", "lazygit", "race_log")
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
	dir, err := os.ReadDir(path)
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
