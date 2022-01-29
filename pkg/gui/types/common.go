package types

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
)

// if Go let me do private struct embedding of structs with public fields (which it should)
// I would just do that. But alas.
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
