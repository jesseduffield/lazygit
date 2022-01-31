package types

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

type ControllerCommon struct {
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
	RunSubprocessAndRefresh(oscommands.ICmdObj) error
	PushContext(context Context, opts ...OnFocusOpts) error
	PopContext() error
	CurrentContext() Context
	// enters search mode for the current view
	OpenSearch()

	GetAppState() *config.AppState
	SaveAppState() error
}

type IPopupHandler interface {
	ErrorMsg(message string) error
	Error(err error) error
	Ask(opts AskOpts) error
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

	// when HandlersManageFocus is true, do not return from the confirmation context automatically. It's expected that the handlers will manage focus, whether that means switching to another context, or manually returning the context.
	HandlersManageFocus bool

	FindSuggestionsFunc func(string) []*Suggestion
}

type AskOpts struct {
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleClose         func() error
	HandlersManageFocus bool
}

type PromptOpts struct {
	Title               string
	InitialContent      string
	FindSuggestionsFunc func(string) []*Suggestion
	HandleConfirm       func(string) error
}

type MenuItem struct {
	DisplayString  string
	DisplayStrings []string
	OnPress        func() error
	// only applies when displayString is used
	OpensMenu bool
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

	BisectInfo     *git_commands.BisectInfo
	RemoteBranches []*models.RemoteBranch
	Tags           []*models.Tag

	// for displaying suggestions while typing in a file name
	FilesTrie *patricia.Trie
}
