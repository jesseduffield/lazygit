package types

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/sasha-s/go-deadlock"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

type HelperCommon struct {
	*common.Common
	IGuiCommon
}

type IGuiCommon interface {
	IPopupHandler

	LogAction(action string)
	LogCommand(cmdStr string, isCommandLine bool)
	// we call this when we want to refetch some models and render the result. Internally calls PostRefreshUpdate
	Refresh(RefreshOptions) error
	// we call this when we've changed something in the view model but not the actual model,
	// e.g. expanding or collapsing a folder in a file view. Calling 'Refresh' in this
	// case would be overkill, although refresh will internally call 'PostRefreshUpdate'
	PostRefreshUpdate(Context) error
	// this just re-renders the screen
	Render()
	// allows rendering to main views (i.e. the ones to the right of the side panel)
	// in such a way that avoids concurrency issues when there are slow commands
	// to display the output of
	RenderToMainViews(opts RefreshMainOpts) error
	// used purely for the sake of RenderToMainViews to provide the pair of main views we want to render to
	MainViewPairs() MainViewPairs

	// returns true if command completed successfully
	RunSubprocess(cmdObj oscommands.ICmdObj) (bool, error)
	RunSubprocessAndRefresh(oscommands.ICmdObj) error

	PushContext(context Context, opts ...OnFocusOpts) error
	PopContext() error
	// Removes all given contexts from the stack. If a given context is not in the stack, it is ignored.
	// This is for when you have a group of contexts that are bundled together e.g. with the commit message panel.
	// If you want to remove a single context, you should probably use PopContext instead.
	RemoveContexts([]Context) error
	CurrentContext() Context
	CurrentStaticContext() Context
	IsCurrentContext(Context) bool
	ActivateContext(context Context) error
	// enters search mode for the current view
	OpenSearch()

	GetAppState() *config.AppState
	SaveAppState() error

	// Runs the given function on the UI thread (this is for things like showing a popup asking a user for input).
	// Only necessary to call if you're not already on the UI thread i.e. you're inside a goroutine.
	// All controller handlers are executed on the UI thread.
	OnUIThread(f func() error)
}

type IPopupHandler interface {
	// Shows a popup with a (localized) "Error" caption and the given error message (in red).
	//
	// This is a convenience wrapper around Alert().
	ErrorMsg(message string) error
	Error(err error) error
	// Shows a notification popup with the given title and message to the user.
	//
	// This is a convenience wrapper around Confirm(), thus the popup can be closed using both 'Enter' and 'ESC'.
	Alert(title string, message string) error
	// Shows a popup asking the user for confirmation.
	Confirm(opts ConfirmOpts) error
	// Shows a popup prompting the user for input.
	Prompt(opts PromptOpts) error
	WithLoaderPanel(message string, f func() error) error
	WithWaitingStatus(message string, f func() error) error
	Menu(opts CreateMenuOptions) error
	Toast(message string)
	GetPromptInput() string
}

type CreateMenuOptions struct {
	Title      string
	Items      []*MenuItem
	HideCancel bool
}

type CreatePopupPanelOpts struct {
	HasLoader           bool
	Editable            bool
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleConfirmPrompt func(string) error
	HandleClose         func() error

	FindSuggestionsFunc func(string) []*Suggestion
	Mask                bool
}

type ConfirmOpts struct {
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleClose         func() error
	HasLoader           bool
	FindSuggestionsFunc func(string) []*Suggestion
	Editable            bool
	Mask                bool
}

type PromptOpts struct {
	Title               string
	InitialContent      string
	FindSuggestionsFunc func(string) []*Suggestion
	HandleConfirm       func(string) error
	// CAPTURE THIS
	HandleClose func() error
	Mask        bool
}

type MenuItem struct {
	Label string

	// alternative to Label. Allows specifying columns which will be auto-aligned
	LabelColumns []string

	OnPress func() error

	// Only applies when Label is used
	OpensMenu bool

	// If Key is defined it allows the user to press the key to invoke the menu
	// item, as opposed to having to navigate to it
	Key Key

	// The tooltip will be displayed upon highlighting the menu item
	Tooltip string
}

type Model struct {
	CommitFiles  []*models.CommitFile
	Files        []*models.File
	Submodules   []*models.SubmoduleConfig
	Branches     []*models.Branch
	Commits      []*models.Commit
	StashEntries []*models.StashEntry
	SubCommits   []*models.Commit
	Remotes      []*models.Remote

	// FilteredReflogCommits are the ones that appear in the reflog panel.
	// when in filtering mode we only include the ones that match the given path
	FilteredReflogCommits []*models.Commit
	// ReflogCommits are the ones used by the branches panel to obtain recency values
	// if we're not in filtering mode, CommitFiles and FilteredReflogCommits will be
	// one and the same
	ReflogCommits []*models.Commit

	BisectInfo                          *git_commands.BisectInfo
	WorkingTreeStateAtLastCommitRefresh enums.RebaseMode
	RemoteBranches                      []*models.RemoteBranch
	Tags                                []*models.Tag

	// for displaying suggestions while typing in a file name
	FilesTrie *patricia.Trie
}

// if you add a new mutex here be sure to instantiate it. We're using pointers to
// mutexes so that we can pass the mutexes to controllers.
type Mutexes struct {
	RefreshingFilesMutex  *deadlock.Mutex
	RefreshingStatusMutex *deadlock.Mutex
	SyncMutex             *deadlock.Mutex
	LocalCommitsMutex     *deadlock.Mutex
	SubCommitsMutex       *deadlock.Mutex
	SubprocessMutex       *deadlock.Mutex
	PopupMutex            *deadlock.Mutex
	PtyMutex              *deadlock.Mutex
}
