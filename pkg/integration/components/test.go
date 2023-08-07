package components

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// IntegrationTest describes an integration test that will be run against the lazygit gui.

// our unit tests will use this description to avoid a panic caused by attempting
// to get the test's name via it's file's path.
const unitTestDescription = "test test"

const (
	defaultWidth  = 100
	defaultHeight = 100
)

type IntegrationTest struct {
	name         string
	description  string
	extraCmdArgs []string
	extraEnvVars map[string]string
	skip         bool
	setupRepo    func(shell *Shell)
	setupConfig  func(config *config.AppConfig)
	run          func(
		testDriver *TestDriver,
		keys config.KeybindingConfig,
	)
	gitVersion    GitVersionRestriction
	width         int
	height        int
	isDemo        bool
	useCustomPath bool
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
	ExtraCmdArgs []string
	// for when a test is flakey
	ExtraEnvVars map[string]string
	Skip         bool
	// to run a test only on certain git versions
	GitVersion GitVersionRestriction
	// width and height when running in headless mode, for testing
	// the UI in different sizes.
	// If these are set, the test must be run in headless mode
	Width  int
	Height int
	// If true, this is not a test but a demo to be added to our docs
	IsDemo bool
	// If true, the test won't invoke lazygit with the --path arg.
	// Useful for when we're passing --git-dir and --work-tree (because --path is
	// incompatible with those args)
	UseCustomPath bool
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
		return lo.SomeBy(self.includes, func(str string) bool {
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
		name:          name,
		description:   args.Description,
		extraCmdArgs:  args.ExtraCmdArgs,
		extraEnvVars:  args.ExtraEnvVars,
		skip:          args.Skip,
		setupRepo:     args.SetupRepo,
		setupConfig:   args.SetupConfig,
		run:           args.Run,
		gitVersion:    args.GitVersion,
		width:         args.Width,
		height:        args.Height,
		isDemo:        args.IsDemo,
		useCustomPath: args.UseCustomPath,
	}
}

func (self *IntegrationTest) Name() string {
	return self.name
}

func (self *IntegrationTest) Description() string {
	return self.description
}

func (self *IntegrationTest) ExtraCmdArgs() []string {
	return self.extraCmdArgs
}

func (self *IntegrationTest) ExtraEnvVars() map[string]string {
	return self.extraEnvVars
}

func (self *IntegrationTest) Skip() bool {
	return self.skip
}

func (self *IntegrationTest) IsDemo() bool {
	return self.isDemo
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
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	shell := NewShell(pwd, func(errorMsg string) { gui.Fail(errorMsg) })
	keys := gui.Keys()
	testDriver := NewTestDriver(gui, shell, keys, KeyPressDelay())

	if KeyPressDelay() > 0 {
		// Setting caption to clear the options menu from whatever it starts with
		testDriver.SetCaption("")
		testDriver.SetCaptionPrefix("")
	}

	self.run(testDriver, keys)

	if KeyPressDelay() > 0 {
		// Clear whatever caption there was so it doesn't linger
		testDriver.SetCaption("")
		testDriver.SetCaptionPrefix("")
		// the dev would want to see the final state if they're running in slow mode
		testDriver.Wait(2000)
	}
}

func (self *IntegrationTest) HeadlessDimensions() (int, int) {
	if self.width == 0 && self.height == 0 {
		return defaultWidth, defaultHeight
	}

	return self.width, self.height
}

func (self *IntegrationTest) RequiresHeadless() bool {
	return self.width != 0 && self.height != 0
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
