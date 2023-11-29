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

const (
	TEST_NAME_ENV_VAR         = "TEST_NAME"
	SANDBOX_ENV_VAR           = "SANDBOX"
	WAIT_FOR_DEBUGGER_ENV_VAR = "WAIT_FOR_DEBUGGER"
	GIT_CONFIG_GLOBAL_ENV_VAR = "GIT_CONFIG_GLOBAL"
)

type RunTestArgs struct {
	Tests           []*IntegrationTest
	Logf            func(format string, formatArgs ...interface{})
	RunCmd          func(cmd *exec.Cmd) (int, error)
	TestWrapper     func(test *IntegrationTest, f func() error)
	Sandbox         bool
	WaitForDebugger bool
	RaceDetector    bool
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
	if err := buildLazygit(args.WaitForDebugger, args.RaceDetector); err != nil {
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
				err := runTest(test, paths, projectRootDir, args.Logf, args.RunCmd, args.Sandbox, args.WaitForDebugger, args.RaceDetector, args.InputDelay, gitVersion)
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
	paths Paths,
	projectRootDir string,
	logf func(format string, formatArgs ...interface{}),
	runCmd func(cmd *exec.Cmd) (int, error),
	sandbox bool,
	waitForDebugger bool,
	raceDetector bool,
	inputDelay int,
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

	cmd, err := getLazygitCommand(test, paths, projectRootDir, sandbox, waitForDebugger, inputDelay)
	if err != nil {
		return err
	}

	pid, err := runCmd(cmd)

	// Print race detector log regardless of the command's exit status
	if raceDetector {
		logPath := fmt.Sprintf("%s.%d", raceDetectorLogsPath(), pid)
		if bytes, err := os.ReadFile(logPath); err == nil {
			logf("Race detector log:\n" + string(bytes))
		}
	}

	return err
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

func buildLazygit(debug bool, raceDetector bool) error {
	// // TODO: remove this line!
	// // skipping this because I'm not making changes to the app code atm.
	// return nil

	args := []string{"go", "build"}
	if debug {
		// Disable compiler optimizations (-N) and inlining (-l) because this
		// makes debugging work better
		args = append(args, "-gcflags=all=-N -l")
	}
	if raceDetector {
		args = append(args, "-race")
	}
	args = append(args, "-o", tempLazygitPath(), filepath.FromSlash("pkg/integration/clients/injector/main.go"))
	osCommand := oscommands.NewDummyOSCommand()
	return osCommand.Cmd.New(args).Run()
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

func getLazygitCommand(test *IntegrationTest, paths Paths, rootDir string, sandbox bool, waitForDebugger bool, inputDelay int) (*exec.Cmd, error) {
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
	if !test.useCustomPath {
		cmdArgs = append(cmdArgs, "--path="+paths.ActualRepo())
	}
	resolvedExtraArgs := lo.Map(test.ExtraCmdArgs(), func(arg string, _ int) string {
		return utils.ResolvePlaceholderString(arg, map[string]string{
			"actualPath":     paths.Actual(),
			"actualRepoPath": paths.ActualRepo(),
		})
	})
	cmdArgs = append(cmdArgs, resolvedExtraArgs...)

	cmdObj := osCommand.Cmd.New(cmdArgs)

	cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", TEST_NAME_ENV_VAR, test.Name()))
	if sandbox {
		cmdObj.AddEnvVars(fmt.Sprintf("%s=%s", SANDBOX_ENV_VAR, "true"))
	}
	if waitForDebugger {
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

	if inputDelay > 0 {
		cmdObj.AddEnvVars(fmt.Sprintf("INPUT_DELAY=%d", inputDelay))
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
