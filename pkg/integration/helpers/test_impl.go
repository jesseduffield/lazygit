package helpers

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	guiTypes "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TestImpl struct {
	name         string
	description  string
	extraCmdArgs string
	skip         bool
	setupRepo    func(shell types.Shell)
	setupConfig  func(config *config.AppConfig)
	run          func(
		shell types.Shell,
		input types.Input,
		assert types.Assert,
		keys config.KeybindingConfig,
	)
}

type NewTestArgs struct {
	Description  string
	SetupRepo    func(shell types.Shell)
	SetupConfig  func(config *config.AppConfig)
	Run          func(shell types.Shell, input types.Input, assert types.Assert, keys config.KeybindingConfig)
	ExtraCmdArgs string
	Skip         bool
}

func NewTest(args NewTestArgs) *TestImpl {
	return &TestImpl{
		name:         testNameFromFilePath(),
		description:  args.Description,
		extraCmdArgs: args.ExtraCmdArgs,
		skip:         args.Skip,
		setupRepo:    args.SetupRepo,
		setupConfig:  args.SetupConfig,
		run:          args.Run,
	}
}

var _ types.Test = (*TestImpl)(nil)

func (self *TestImpl) Name() string {
	return self.name
}

func (self *TestImpl) Description() string {
	return self.description
}

func (self *TestImpl) ExtraCmdArgs() string {
	return self.extraCmdArgs
}

func (self *TestImpl) Skip() bool {
	return self.skip
}

func (self *TestImpl) SetupConfig(config *config.AppConfig) {
	self.setupConfig(config)
}

func (self *TestImpl) SetupRepo(shell types.Shell) {
	self.setupRepo(shell)
}

// I want access to all contexts, the model, the ability to press a key, the ability to log,
func (self *TestImpl) Run(gui guiTypes.GuiAdapter) {
	shell := &ShellImpl{}
	assert := &AssertImpl{gui: gui}
	keys := gui.Keys()
	input := NewInputImpl(gui, keys, assert, KeyPressDelay())

	self.run(shell, input, assert, keys)
}

func testNameFromFilePath() string {
	path := utils.FilePath(3)
	name := strings.Split(path, "integration/integration_tests/")[1]

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
