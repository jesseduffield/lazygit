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
	GitAdd(path string) Shell
	GitAddAll() Shell
	Commit(message string) Shell
	EmptyCommit(message string) Shell
	// convenience method for creating a file and adding it
	CreateFileAndAdd(fileName string, fileContents string) Shell
	// creates commits 01, 02, 03, ..., n with a new file in each
	// The reason for padding with zeroes is so that it's easier to do string
	// matches on the commit messages when there are many of them
	CreateNCommits(n int) Shell
}

// through this interface our test interacts with the lazygit gui
// Implementation is at pkg/gui/input.go
type Input interface {
	// key is something like 'w' or '<space>'. It's best not to pass a direct value,
	// but instead to go through the default user config to get a more meaningful key name
	PressKeys(keys ...string)
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
	// this will look for a list item in the current panel and if it finds it, it will
	// enter the keypresses required to navigate to it.
	// The test will fail if:
	//  - the user is not in a list item
	//  - no list item is found containing the given text
	//  - multiple list items are found containing the given text in the initial page of items
	NavigateToListItemContainingText(text string)
	ContinueRebase()
	ContinueMerge()
}

// through this interface we assert on the state of the lazygit gui
type Assert interface {
	WorkingTreeFileCount(int)
	CommitCount(int)
	HeadCommitMessage(string)
	CurrentViewName(expectedViewName string)
	CurrentBranchName(expectedBranchName string)
	InListContext()
	SelectedLineContains(text string)
	// for when you just want to fail the test yourself
	Fail(errorMessage string)
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
