package types

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// TODO: refactor this so that we don't have code spread around so much. We want
// our TestImpl struct to take the dependencies it needs from the gui and then
// create the input, assert, shell structs itself. That way, we can potentially
// ditch these interfaces so that we don't need to keep updating them every time
// we add a method to the concrete struct.

type Test interface {
	Name() string
	Description() string
	// this is called before lazygit is run, for the sake of preparing the repo
	SetupRepo(Shell)
	// this gives you the default config and lets you set whatever values on it you like,
	// so that they appear when lazygit runs
	SetupConfig(config *config.AppConfig)
	// this is called upon lazygit starting
	Run(GuiAdapter)
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
// implementation is at pkg/gui/assert.go
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

type GuiAdapter interface {
	PressKey(string)
	Keys() config.KeybindingConfig
	CurrentContext() types.Context
	Model() *types.Model
	Fail(message string)
	// These two log methods are for the sake of debugging while testing. There's no need to actually
	// commit any logging.
	// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
	Log(message string)
	// logs in the actual UI (in the commands panel)
	LogUI(message string)
	CheckedOutRef() *models.Branch
}
