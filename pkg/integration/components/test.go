package components

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// IntegrationTest describes an integration test that will be run against the lazygit gui.

// our unit tests will use this description to avoid a panic caused by attempting
// to get the test's name via it's file's path.
const unitTestDescription = "test test"

type IntegrationTest struct {
	name         string
	description  string
	extraCmdArgs string
	skip         bool
	setupRepo    func(shell *Shell)
	setupConfig  func(config *config.AppConfig)
	run          func(
		testDriver *TestDriver,
		keys config.KeybindingConfig,
	)
	gitVersion GitVersionRestriction
}

var _ integrationTypes.IntegrationTest = &IntegrationTest{}

type NewIntegrationTestArgs struct {
	// Briefly describes what happens in the test and what it's testing for
	Description string
	// prepares a repo for testing
	SetupRepo func(shell *Shell)
	// takes a config and mutates. The mutated context will end up being passed to the gui
	SetupConfig func(config *config.AppConfig)
	// runs the test
	Run func(t *TestDriver, keys config.KeybindingConfig)
	// additional args passed to lazygit
	ExtraCmdArgs string
	// for when a test is flakey
	Skip bool
	// to run a test only on certain git versions
	GitVersion GitVersionRestriction
}

type GitVersionRestriction struct {
	// Only one of these fields can be non-empty; use functions below to construct
	from     string
	before   string
	includes []string
}

// Verifies the version is at least the given version (inclusive)
func AtLeast(version string) GitVersionRestriction {
	return GitVersionRestriction{from: version}
}

// Verifies the version is before the given version (exclusive)
func Before(version string) GitVersionRestriction {
	return GitVersionRestriction{before: version}
}

func Includes(versions ...string) GitVersionRestriction {
	return GitVersionRestriction{includes: versions}
}

func (self GitVersionRestriction) shouldRunOnVersion(version *git_commands.GitVersion) bool {
	if self.from != "" {
		from, err := git_commands.ParseGitVersion(self.from)
		if err != nil {
			panic("Invalid git version string: " + self.from)
		}
		return !version.IsOlderThanVersion(from)
	}
	if self.before != "" {
		before, err := git_commands.ParseGitVersion(self.before)
		if err != nil {
			panic("Invalid git version string: " + self.before)
		}
		return version.IsOlderThanVersion(before)
	}
	if len(self.includes) != 0 {
		return slices.Some(self.includes, func(str string) bool {
			v, err := git_commands.ParseGitVersion(str)
			if err != nil {
				panic("Invalid git version string: " + str)
			}
			return version.Major == v.Major && version.Minor == v.Minor && version.Patch == v.Patch
		})
	}
	return true
}

func NewIntegrationTest(args NewIntegrationTestArgs) *IntegrationTest {
	name := ""
	if args.Description != unitTestDescription {
		// this panics if we're in a unit test for our integration tests,
		// so we're using "test test" as a sentinel value
		name = testNameFromCurrentFilePath()
	}

	return &IntegrationTest{
		name:         name,
		description:  args.Description,
		extraCmdArgs: args.ExtraCmdArgs,
		skip:         args.Skip,
		setupRepo:    args.SetupRepo,
		setupConfig:  args.SetupConfig,
		run:          args.Run,
		gitVersion:   args.GitVersion,
	}
}

func (self *IntegrationTest) Name() string {
	return self.name
}

func (self *IntegrationTest) Description() string {
	return self.description
}

func (self *IntegrationTest) ExtraCmdArgs() string {
	return self.extraCmdArgs
}

func (self *IntegrationTest) Skip() bool {
	return self.skip
}

func (self *IntegrationTest) ShouldRunForGitVersion(version *git_commands.GitVersion) bool {
	return self.gitVersion.shouldRunOnVersion(version)
}

func (self *IntegrationTest) SetupConfig(config *config.AppConfig) {
	self.setupConfig(config)
}

func (self *IntegrationTest) SetupRepo(shell *Shell) {
	self.setupRepo(shell)
}

func (self *IntegrationTest) Run(gui integrationTypes.GuiDriver) {
	// we pass the --pass arg to lazygit when running an integration test, and that
	// ends up stored in the following env var
	repoPath := env.GetGitWorkTreeEnv()

	shell := NewShell(repoPath, func(errorMsg string) { gui.Fail(errorMsg) })
	keys := gui.Keys()
	testDriver := NewTestDriver(gui, shell, keys, KeyPressDelay())

	self.run(testDriver, keys)

	if KeyPressDelay() > 0 {
		// the dev would want to see the final state if they're running in slow mode
		testDriver.Wait(2000)
	}
}

func testNameFromCurrentFilePath() string {
	path := utils.FilePath(3)
	return TestNameFromFilePath(path)
}

func TestNameFromFilePath(path string) string {
	name := strings.Split(path, "integration/tests/")[1]

	return name[:len(name)-len(".go")]
}

// this is the delay in milliseconds between keypresses
// defaults to zero
func KeyPressDelay() int {
	delayStr := os.Getenv("KEY_PRESS_DELAY")
	if delayStr == "" {
		return 0
	}

	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		panic(err)
	}
	return delay
}
