package types

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
)

type ContextKind int

const (
	// this is your files, branches, commits, contexts etc. They're all on the left hand side
	// and you can cycle through them.
	SIDE_CONTEXT ContextKind = iota
	// This is either the left or right 'main' contexts that appear to the right of the side contexts
	MAIN_CONTEXT
	// A persistent popup is one that has its own identity e.g. the commit message context.
	// When you open a popup over it, we'll let you return to it upon pressing escape
	PERSISTENT_POPUP
	// A temporary popup is one that could be used for various things (e.g. a generic menu or confirmation popup).
	// Because we reuse these contexts, they're temporary in that you can't return to them after you've switched from them
	// to some other context, because the context you switched to might actually be the same context but rendering different content.
	// We should really be able to spawn new contexts for menus/prompts so that we can actually return to old ones.
	TEMPORARY_POPUP
	// This contains the command log, underneath the main contexts.
	EXTRAS_CONTEXT
	// only used by the one global context, purely for the sake of defining keybindings globally
	GLOBAL_CONTEXT
	// a display context only renders a view. It has no keybindings associated and
	// it cannot receive focus.
	DISPLAY_CONTEXT
)

type ParentContexter interface {
	SetParentContext(Context)
	GetParentContext() Context
}

type NeedsRerenderOnWidthChangeLevel int

const (
	// view doesn't render differently when its width changes
	NEEDS_RERENDER_ON_WIDTH_CHANGE_NONE NeedsRerenderOnWidthChangeLevel = iota
	// view renders differently when its width changes. An example is a view
	// that truncates long lines to the view width, e.g. the branches view
	NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_WIDTH_CHANGES
	// view renders differently only when the screen mode changes
	NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_SCREEN_MODE_CHANGES
)

type IBaseContext interface {
	HasKeybindings
	ParentContexter

	GetKind() ContextKind
	GetViewName() string
	GetView() *gocui.View
	GetViewTrait() IViewTrait
	GetWindowName() string
	SetWindowName(string)
	GetKey() ContextKey
	IsFocusable() bool
	// if a context is transient, then it only appears via some keybinding on another
	// context. Until we add support for having multiple of the same context, no two
	// of the same transient context can appear at once meaning one might be 'stolen'
	// from another window.
	IsTransient() bool
	// this tells us if the view's bounds are determined by its window or if they're
	// determined independently.
	HasControlledBounds() bool

	// the total height of the content that the view is currently showing
	TotalContentHeight() int

	// to what extent the view needs to be rerendered when its width changes
	NeedsRerenderOnWidthChange() NeedsRerenderOnWidthChangeLevel

	// true if the view needs to be rerendered when its height changes
	NeedsRerenderOnHeightChange() bool

	// returns the desired title for the view upon activation. If there is no desired title (returns empty string), then
	// no title will be set
	Title() string

	GetOptionsMap() map[string]string

	AddKeybindingsFn(KeybindingsFn)
	AddMouseKeybindingsFn(MouseKeybindingsFn)
	ClearAllBindingsFn()

	// This is a bit of a hack at the moment: we currently only set an onclick function so that
	// our list controller can come along and wrap it in a list-specific click handler.
	// We'll need to think of a better way to do this.
	AddOnClickFn(func() error)
	// Likewise for the focused main view: we need this to communicate between a
	// side panel controller and the focused main view controller.
	AddOnClickFocusedMainViewFn(func(mainViewName string, clickedLineIdx int) error)

	AddOnRenderToMainFn(func())
	AddOnFocusFn(func(OnFocusOpts))
	AddOnFocusLostFn(func(OnFocusLostOpts))
}

type Context interface {
	IBaseContext

	HandleFocus(opts OnFocusOpts)
	HandleFocusLost(opts OnFocusLostOpts)
	FocusLine()
	HandleRender()
	HandleRenderToMain()
}

type ISearchHistoryContext interface {
	Context

	GetSearchHistory() *utils.HistoryBuffer[string]
}

type IFilterableContext interface {
	Context
	IListPanelState
	ISearchHistoryContext

	SetFilter(string, bool)
	GetFilter() string
	ClearFilter()
	ReApplyFilter(bool)
	IsFiltering() bool
	IsFilterableContext()
}

type ISearchableContext interface {
	Context
	ISearchHistoryContext

	// These are all implemented by SearchTrait
	SetSearchString(string)
	GetSearchString() string
	ClearSearchString()
	IsSearching() bool
	IsSearchableContext()
	RenderSearchStatus(int, int)

	// This must be implemented by each concrete context. Return nil if not searching the model.
	ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition
}

type DiffableContext interface {
	Context

	// Returns the current diff terminals of the currently selected item.
	// in the case of a branch it returns both the branch and it's upstream name,
	// which becomes an option when you bring up the diff menu, but when you're just
	// flicking through branches it will be using the local branch name.
	GetDiffTerminals() []string

	// Returns the ref that should be used for creating a diff of what's
	// currently shown in the main view against the working directory, in order
	// to adjust line numbers in the diff to match the current state of the
	// shown file. For example, if the main view shows a range diff of commits,
	// we need to pass the first commit of the range. This is used by
	// DiffHelper.AdjustLineNumber.
	RefForAdjustingLineNumberInDiff() string
}

type IListContext interface {
	Context

	GetSelectedItemId() string
	GetSelectedItemIds() ([]string, int, int)
	IsItemVisible(item HasUrn) bool

	GetList() IList
	ViewIndexToModelIndex(int) int
	ModelIndexToViewIndex(int) int

	IsListContext() // used for type switch
	RangeSelectEnabled() bool
	RenderOnlyVisibleLines() bool

	IndexForGotoBottom() int
}

type IPatchExplorerContext interface {
	Context

	GetState() *patch_exploring.State
	SetState(*patch_exploring.State)
	GetIncludedLineIndices() []int
	RenderAndFocus()
	Render()
	Focus()
	GetContentToRender() string
	NavigateTo(selectedLineIdx int)
	GetMutex() *deadlock.Mutex
	IsPatchExplorerContext() // used for type switch
}

type IViewTrait interface {
	FocusPoint(yIdx int)
	SetRangeSelectStart(yIdx int)
	CancelRangeSelect()
	SetViewPortContent(content string)
	SetViewPortContentAndClearEverythingElse(content string)
	SetContentLineCount(lineCount int)
	SetContent(content string)
	SetFooter(value string)
	SetOriginX(value int)
	ViewPortYBounds() (int, int)
	ScrollLeft()
	ScrollRight()
	ScrollUp(value int)
	ScrollDown(value int)
	PageDelta() int
	SelectedLineIdx() int
	SetHighlight(bool)
}

type OnFocusOpts struct {
	ClickedWindowName  string
	ClickedViewLineIdx int
}

type OnFocusLostOpts struct {
	NewContextKey ContextKey
}

type ContextKey string

type KeybindingsOpts struct {
	GetKey func(key string) Key
	Config config.KeybindingConfig
	Guards KeybindingGuards
}

type (
	KeybindingsFn      func(opts KeybindingsOpts) []*Binding
	MouseKeybindingsFn func(opts KeybindingsOpts) []*gocui.ViewMouseBinding
)

type HasKeybindings interface {
	GetKeybindings(opts KeybindingsOpts) []*Binding
	GetMouseKeybindings(opts KeybindingsOpts) []*gocui.ViewMouseBinding
	GetOnClick() func() error
	GetOnClickFocusedMainView() func(mainViewName string, clickedLineIdx int) error
	GetOnRenderToMain() func()
	GetOnFocus() func(OnFocusOpts)
	GetOnFocusLost() func(OnFocusLostOpts)
}

type IController interface {
	HasKeybindings
	Context() Context
}

type IList interface {
	IListCursor
	Len() int
	GetItem(index int) HasUrn
}

type IListCursor interface {
	GetSelectedLineIdx() int
	SetSelectedLineIdx(value int)
	SetSelection(value int)
	MoveSelectedLine(delta int)
	ClampSelection()
	CancelRangeSelect()
	GetRangeStartIdx() (int, bool)
	GetSelectionRange() (int, int)
	IsSelectingRange() bool
	AreMultipleItemsSelected() bool
	ToggleStickyRange()
	ExpandNonStickyRange(int)
}

type IListPanelState interface {
	SetSelectedLineIdx(int)
	SetSelection(int)
	GetSelectedLineIdx() int
}

type ListItem interface {
	// ID is a hash when the item is a commit, a filename when the item is a file, 'stash@{4}' when it's a stash entry, 'my_branch' when it's a branch
	ID() string

	// Description is something we would show in a message e.g. '123as14: push blah' for a commit
	Description() string
}

type IContextMgr interface {
	Push(context Context, opts OnFocusOpts)
	Pop()
	Replace(context Context)
	Activate(context Context, opts OnFocusOpts)
	Current() Context
	CurrentStatic() Context
	CurrentSide() Context
	CurrentPopup() []Context
	NextInStack(context Context) Context
	IsCurrent(c Context) bool
	IsCurrentOrParent(c Context) bool
	ForEach(func(Context))
	AllList() []IListContext
	AllFilterable() []IFilterableContext
	AllSearchable() []ISearchableContext
	AllPatchExplorer() []IPatchExplorerContext
}
