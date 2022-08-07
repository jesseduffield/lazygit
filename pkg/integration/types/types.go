package types

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type Test interface {
	Name() string
	Description() string
	// this is called before lazygit is run, for the sake of preparing the repo
	SetupRepo(Shell)
	// this gives you the default config and lets you set whatever values on it you like,
	// so that they appear when lazygit runs
	SetupConfig(config *config.AppConfig)
	// this is called upon lazygit starting
	Run(Shell, Input, Assert, config.KeybindingConfig)
	// e.g. '-debug'
	ExtraCmdArgs() string
	// for tests that are flakey and when we don't have time to fix them
	Skip() bool
}

// this is for running shell commands, mostly for the sake of setting up the repo
// but you can also run the commands from within lazygit to emulate things happening
// in the background.
// Implementation is at pkg/integration/shell.go
type Shell interface {
	RunCommand(command string) Shell
	CreateFile(name string, content string) Shell
	NewBranch(branchName string) Shell
	GitAddAll() Shell
	Commit(message string) Shell
	EmptyCommit(message string) Shell
}

// through this interface our test interacts with the lazygit gui
// Implementation is at pkg/gui/input.go
type Input interface {
	// key is something like 'w' or '<space>'. It's best not to pass a direct value,
	// but instead to go through the default user config to get a more meaningful key name
	PushKeys(keys ...string)
	// for typing into a popup prompt
	Type(content string)
	// for when you want to allow lazygit to process something before continuing
	Wait(milliseconds int)
	// going straight to a particular side window
	SwitchToStatusWindow()
	SwitchToFilesWindow()
	SwitchToBranchesWindow()
	SwitchToCommitsWindow()
	SwitchToStashWindow()
	// i.e. pressing enter
	Confirm()
	// i.e. pressing escape
	Cancel()
	// i.e. pressing space
	Select()
	// i.e. pressing down arrow
	NextItem()
	// i.e. pressing up arrow
	PreviousItem()
}

// through this interface we assert on the state of the lazygit gui
type Assert interface {
	WorkingTreeFileCount(int)
	CommitCount(int)
	HeadCommitMessage(string)
	CurrentViewName(expectedViewName string)
	CurrentBranchName(expectedBranchName string)
}

type TestImpl struct {
	name         string
	description  string
	extraCmdArgs string
	skip         bool
	setupRepo    func(shell Shell)
	setupConfig  func(config *config.AppConfig)
	run          func(
		shell Shell,
		input Input,
		assert Assert,
		keys config.KeybindingConfig,
	)
}

type NewTestArgs struct {
	Description  string
	SetupRepo    func(shell Shell)
	SetupConfig  func(config *config.AppConfig)
	Run          func(shell Shell, input Input, assert Assert, keys config.KeybindingConfig)
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

var _ Test = (*TestImpl)(nil)

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

func (self *TestImpl) SetupRepo(shell Shell) {
	self.setupRepo(shell)
}

func (self *TestImpl) Run(
	shell Shell,
	input Input,
	assert Assert,
	keys config.KeybindingConfig,
) {
	self.run(shell, input, assert, keys)
}

func testNameFromFilePath() string {
	path := utils.FilePath(3)
	name := strings.Split(path, "integration/integration_tests/")[1]

	return name[:len(name)-len(".go")]
}
