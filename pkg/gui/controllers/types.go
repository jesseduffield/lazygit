package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IGuiCommon interface {
	popup.IPopupHandler

	LogAction(action string)
	LogCommand(cmdStr string, isCommandLine bool)
	// we call this when we want to refetch some models and render the result. Internally calls PostRefreshUpdate
	Refresh(types.RefreshOptions) error
	// we call this when we've changed something in the view model but not the actual model,
	// e.g. expanding or collapsing a folder in a file view. Calling 'Refresh' in this
	// case would be overkill, although refresh will internally call 'PostRefreshUpdate'
	PostRefreshUpdate(types.Context) error
	RunSubprocessAndRefresh(oscommands.ICmdObj) error
	PushContext(context types.Context, opts ...types.OnFocusOpts) error
	PopContext() error

	GetAppState() *config.AppState
	SaveAppState() error
}

type IRefHelper interface {
	CheckoutRef(ref string, options types.CheckoutRefOptions) error
	CreateGitResetMenu(ref string) error
	ResetToRef(ref string, strength string, envVars []string) error
}

type ISuggestionsHelper interface {
	GetRemoteSuggestionsFunc() func(string) []*types.Suggestion
	GetBranchNameSuggestionsFunc() func(string) []*types.Suggestion
	GetFilePathSuggestionsFunc() func(string) []*types.Suggestion
	GetRemoteBranchesSuggestionsFunc(separator string) func(string) []*types.Suggestion
	GetRefsSuggestionsFunc() func(string) []*types.Suggestion
	GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion
}

type IFileHelper interface {
	EditFile(filename string) error
	EditFileAtLine(filename string, lineNumber int) error
	OpenFile(filename string) error
}

type IWorkingTreeHelper interface {
	AnyStagedFiles() bool
	AnyTrackedFiles() bool
	IsWorkingTreeDirty() bool
	FileForSubmodule(submodule *models.SubmoduleConfig) *models.File
}

// all fields mandatory (except `CanRebase` because it's boolean)
type SwitchToCommitFilesContextOpts struct {
	RefName    string
	CanRebase  bool
	Context    types.Context
	WindowName string
}
