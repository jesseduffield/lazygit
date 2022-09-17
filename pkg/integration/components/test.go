package components

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Test describes an integration tests that will be run against the lazygit gui.

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
		shell *Shell,
		input *Input,
		assert *Assert,
		keys config.KeybindingConfig,
	)
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
	Run func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig)
	// additional args passed to lazygit
	ExtraCmdArgs string
	// for when a test is flakey
	Skip bool
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

func (self *IntegrationTest) SetupConfig(config *config.AppConfig) {
	self.setupConfig(config)
}

func (self *IntegrationTest) SetupRepo(shell *Shell) {
	self.setupRepo(shell)
}

// I want access to all contexts, the model, the ability to press a key, the ability to log,
func (self *IntegrationTest) Run(gui integrationTypes.GuiDriver) {
	shell := NewShell("/tmp/lazygit-test")
	assert := NewAssert(gui)
	keys := gui.Keys()
	input := NewInput(gui, keys, assert, KeyPressDelay())

	self.run(shell, input, assert, keys)

	if KeyPressDelay() > 0 {
		// the dev would want to see the final state if they're running in slow mode
		input.Wait(2000)
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
