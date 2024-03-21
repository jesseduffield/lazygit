package types

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

type HelperCommon struct {
	*ContextCommon
}

type ContextCommon struct {
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

	// renders string to a view without resetting its origin
	SetViewContent(view *gocui.View, content string)
	// resets cursor and origin of view. Often used before calling SetViewContent
	ResetViewOrigin(view *gocui.View)

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
	ReplaceContext(context Context) error
	// Removes all given contexts from the stack. If a given context is not in the stack, it is ignored.
	// This is for when you have a group of contexts that are bundled together e.g. with the commit message panel.
	// If you want to remove a single context, you should probably use PopContext instead.
	RemoveContexts([]Context) error
	CurrentContext() Context
	CurrentStaticContext() Context
	CurrentSideContext() Context
	IsCurrentContext(Context) bool
	// TODO: replace the above context-based methods with just using Context() e.g. replace PushContext() with Context().Push()
	Context() IContextMgr
	ContextForKey(key ContextKey) Context

	ActivateContext(context Context) error

	GetConfig() config.AppConfigurer
	GetAppState() *config.AppState
	SaveAppState() error
	SaveAppStateAndLogError()

	// Runs the given function on the UI thread (this is for things like showing a popup asking a user for input).
	// Only necessary to call if you're not already on the UI thread i.e. you're inside a goroutine.
	// All controller handlers are executed on the UI thread.
	OnUIThread(f func() error)
	// Runs a function in a goroutine. Use this whenever you want to run a goroutine and keep track of the fact
	// that lazygit is still busy. See docs/dev/Busy.md
	OnWorker(f func(gocui.Task))
	// Function to call at the end of our 'layout' function which renders views
	// For example, you may want a view's line to be focused only after that view is
	// resized, if in accordion mode.
	AfterLayout(f func() error)

	// Wraps a function, attaching the given operation to the given item while
	// the function is executing, and also causes the given context to be
	// redrawn periodically. This allows the operation to be visualized with a
	// spinning loader animation (e.g. when a branch is being pushed).
	WithInlineStatus(item HasUrn, operation ItemOperation, contextKey ContextKey, f func(gocui.Task) error) error

	// returns the gocui Gui struct. There is a good chance you don't actually want to use
	// this struct and instead want to use another method above
	GocuiGui() *gocui.Gui

	Views() Views

	Git() *commands.GitCommand
	OS() *oscommands.OSCommand
	Model() *Model

	Modes() *Modes

	Mutexes() Mutexes

	State() IStateAccessor

	KeybindingsOpts() KeybindingsOpts
	CallKeybindingHandler(binding *Binding) error

	// hopefully we can remove this once we've moved all our keybinding stuff out of the gui god struct.
	GetInitialKeybindingsWithCustomCommands() ([]*Binding, []*gocui.ViewMouseBinding)

	// Returns true if we're running an integration test
	RunningIntegrationTest() bool

	// Returns true if we're in a demo recording/playback
	InDemo() bool
}

type IModeMgr interface {
	IsAnyModeActive() bool
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
	WithWaitingStatus(message string, f func(gocui.Task) error) error
	WithWaitingStatusSync(message string, f func() error) error
	Menu(opts CreateMenuOptions) error
	Toast(message string)
	ErrorToast(message string)
	SetToastFunc(func(string, ToastKind))
	GetPromptInput() string
}

type ToastKind int

const (
	ToastKindStatus ToastKind = iota
	ToastKindError
)

type CreateMenuOptions struct {
	Title           string
	Items           []*MenuItem
	HideCancel      bool
	ColumnAlignment []utils.Alignment
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

type MenuSection struct {
	Title  string
	Column int // The column that this section title should be aligned with
}

type DisabledReason struct {
	Text string

	// When trying to invoke a disabled key binding or menu item, we normally
	// show the disabled reason as a toast; setting this to true shows it as an
	// error panel instead. This is useful if the text is very long, or if it is
	// important enough to show it more prominently, or both.
	ShowErrorInPanel bool
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

	// If non-nil, show this in a tooltip, style the menu item as disabled,
	// and refuse to invoke the command
	DisabledReason *DisabledReason

	// Can be used to group menu items into sections with headers. MenuItems
	// with the same Section should be contiguous, and will automatically get a
	// section header. If nil, the item is not part of a section.
	// Note that pointer comparison is used to determine whether two menu items
	// belong to the same section, so make sure all your items in a given
	// section point to the same MenuSection instance.
	Section *MenuSection
}

// Defining this for the sake of conforming to the HasID interface, which is used
// in list contexts.
func (self *MenuItem) ID() string {
	return self.Label
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
	Worktrees    []*models.Worktree

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

	// Name of the currently checked out branch. This will be set even when
	// we're on a detached head because we're rebasing or bisecting.
	CheckedOutBranch string

	// for displaying suggestions while typing in a file name
	FilesTrie *patricia.Trie

	Authors map[string]*models.Author
}

// if you add a new mutex here be sure to instantiate it. We're using pointers to
// mutexes so that we can pass the mutexes to controllers.
type Mutexes struct {
	RefreshingFilesMutex    *deadlock.Mutex
	RefreshingBranchesMutex *deadlock.Mutex
	RefreshingStatusMutex   *deadlock.Mutex
	LocalCommitsMutex       *deadlock.Mutex
	SubCommitsMutex         *deadlock.Mutex
	AuthorsMutex            *deadlock.Mutex
	SubprocessMutex         *deadlock.Mutex
	PopupMutex              *deadlock.Mutex
	PtyMutex                *deadlock.Mutex
}

// A long-running operation associated with an item. For example, we'll show
// that a branch is being pushed from so that there's visual feedback about
// what's happening and so that you can see multiple branches' concurrent
// operations
type ItemOperation int

const (
	ItemOperationNone ItemOperation = iota
	ItemOperationPushing
	ItemOperationPulling
	ItemOperationFastForwarding
	ItemOperationDeleting
	ItemOperationFetching
	ItemOperationCheckingOut
)

type HasUrn interface {
	URN() string
}

type IStateAccessor interface {
	GetRepoPathStack() *utils.StringStack
	GetRepoState() IRepoStateAccessor
	// tells us whether we're currently updating lazygit
	GetUpdating() bool
	SetUpdating(bool)
	SetIsRefreshingFiles(bool)
	GetIsRefreshingFiles() bool
	GetShowExtrasWindow() bool
	SetShowExtrasWindow(bool)
	GetRetainOriginalDir() bool
	SetRetainOriginalDir(bool)
	GetItemOperation(item HasUrn) ItemOperation
	SetItemOperation(item HasUrn, operation ItemOperation)
	ClearItemOperation(item HasUrn)
}

type IRepoStateAccessor interface {
	GetViewsSetup() bool
	GetWindowViewNameMap() *utils.ThreadSafeMap[string, string]
	GetStartupStage() StartupStage
	SetStartupStage(stage StartupStage)
	GetCurrentPopupOpts() *CreatePopupPanelOpts
	SetCurrentPopupOpts(*CreatePopupPanelOpts)
	GetScreenMode() WindowMaximisation
	SetScreenMode(WindowMaximisation)
	InSearchPrompt() bool
	GetSearchState() *SearchState
	SetSplitMainPanel(bool)
	GetSplitMainPanel() bool
}

// startup stages so we don't need to load everything at once
type StartupStage int

const (
	INITIAL StartupStage = iota
	COMPLETE
)

// screen sizing determines how much space your selected window takes up (window
// as in panel, not your terminal's window). Sometimes you want a bit more space
// to see the contents of a panel, and this keeps track of how much maximisation
// you've set
type WindowMaximisation int

const (
	SCREEN_NORMAL WindowMaximisation = iota
	SCREEN_HALF
	SCREEN_FULL
)
