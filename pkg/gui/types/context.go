package types

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/i18n"
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
	ClearAllAttachedControllerFunctions()

	// This is a bit of a hack at the moment: we currently only set an onDoubleClick function so
	// that the generic ListController can be specialized by view-specific controllers.
	// We'll need to think of a better way to do this.
	AddOnDoubleClickFn(func() error)
	// Likewise for the focused main view: we need this to communicate between a
	// side panel controller and the focused main view controller.
	AddOnClickFocusedMainViewFn(func(mainViewName string, clickedLineIdx int) error)
	// And for staging the selected line(s) directly from the focused main view
	// (space), delegated to the side panel that owns the diff being shown. The
	// inclusive view-line range is the current selection (a single line, a range, or
	// a hunk).
	AddOnStageFocusedMainViewFn(func(mainViewName string, firstLineIdx int, lastLineIdx int) (focusViewName string, err error))
	// Adding on to the above, this is so that a list-specific handler can register
	// a hook for doing additional click handling
	AddOnClickFn(func(opts gocui.ViewMouseBindingOpts) error)

	AddOnRenderToMainFn(func())
	AddOnFocusFn(func(OnFocusOpts))
	AddOnFocusLostFn(func(OnFocusLostOpts))
	AddOnQuitFn(func())
}

type Context interface {
	IBaseContext

	HandleFocus(opts OnFocusOpts)
	HandleFocusLost(opts OnFocusLostOpts)
	HandleQuit()
	FocusLine(scrollIntoView bool)
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
	FilterPrefix(tr *i18n.TranslationSet) string
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
	OnSearchSelect(selectedLineIdx int)
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

// DiffMainViewContext is implemented by the side panel contexts whose focused
// main view shows a unified diff — files, local commits, sub-commits, reflog,
// stash, and commit files — as opposed to a commit log or other non-diff content
// (branches, tags, status, …). It is distinct from DiffableContext, which is
// about producing a diff between two refs for the diff menu. This is the signal
// for whether to show a selection in the focused main view: a selection is only
// meaningful where there are diff lines to act on (stage, edit, jump by hunk,
// open in a pull request).
type DiffMainViewContext interface {
	Context

	IsDiffMainViewContext()
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
	SetNeedRerenderVisibleLines()

	IndexForGotoBottom() int
}

type IPatchExplorerContext interface {
	Context

	GetState() *patch_exploring.State
	SetState(*patch_exploring.State)
	GetIncludedLineIndices() []int
	RenderAndFocus()
	Render()
	GetContentToRender() string
	NavigateTo(selectedLineIdx int)
	GetMutex() *deadlock.Mutex
	IsPatchExplorerContext() // used for type switch

	// See FocusedMainViewSnapshot. Nil unless this patch explorer was entered
	// from a focused main view.
	GetFocusedMainViewSnapshot() *FocusedMainViewSnapshot
	SetFocusedMainViewSnapshot(*FocusedMainViewSnapshot)
}

// FocusedMainViewSnapshot records where a focused main view was when we dived
// into a patch explorer (staging or patch building) from it, so that escaping
// returns us to the same place with the main view focused again. It is nil when
// the patch explorer was entered the normal way (through a side panel), in which
// case escape just pops to that side panel.
type FocusedMainViewSnapshot struct {
	// The side panel to land on first; pushing it re-renders the original
	// content into the main view. For commits/stash this is the originating side
	// panel (skipping the commit files panel we passed through), preserving the
	// pre-existing "escape all the way out" behavior.
	SidePanel Context
	// The side panel's selected line, to restore before re-rendering it. Diving
	// into staging can change the side panel's selection (e.g. from a directory
	// to a file in the files panel); restoring it makes the main view show the
	// same content again. -1 if the side panel isn't a list.
	SidePanelSelectedLineIdx int
	// The focused main view context to focus afterwards. Where in it to scroll to
	// and select is not captured here: on escape we land on the line the patch
	// explorer ended up selecting, found by its patch identity in the re-rendered
	// content, which survives the diff changing in a way a saved index wouldn't.
	MainView Context
}

type IViewTrait interface {
	FocusPoint(yIdx int, scrollIntoView bool)
	SetRangeSelectStart(yIdx int)
	CancelRangeSelect()
	SetViewPortContent(content string)
	SetViewPortContentAndClearEverythingElse(lineCount int, content string)
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

	// A source line number identifying the line to land on in the patch
	// explorer. If not -1, takes precedence over ClickedViewLineIdx. It is a
	// new-file line number, unless ClickedViewRealLineIsDeletion is set, in which
	// case it is an old-file line number (two consecutive deletions share a
	// new-file line number, so only the old-file number identifies a deletion).
	ClickedViewRealLineIdx int

	// Whether ClickedViewRealLineIdx is an old-file line number for a deletion;
	// see above.
	ClickedViewRealLineIsDeletion bool

	// When entering a patch explorer (staging or patch building) by clicking or
	// pressing enter on a line in a focused main view, we select that line using
	// the default select mode (hunk or line, per the UseHunkModeInStagingView
	// config), the same as when entering through the side panel. Clicking
	// directly on the patch explorer view instead starts a range selection that
	// can be extended by dragging.
	SelectLineInDefaultMode bool

	ScrollSelectionIntoView bool
}

type OnFocusLostOpts struct {
	NewContextKey ContextKey
}

type ContextKey string

type KeybindingsOpts struct {
	GetKeys func(keys config.Keybinding) []gocui.Key
	Config  config.KeybindingConfig
	Guards  KeybindingGuards
}

type (
	KeybindingsFn      func(opts KeybindingsOpts) []*Binding
	MouseKeybindingsFn func(opts KeybindingsOpts) []*gocui.ViewMouseBinding
)

type HasKeybindings interface {
	GetKeybindings(opts KeybindingsOpts) []*Binding
	GetMouseKeybindings(opts KeybindingsOpts) []*gocui.ViewMouseBinding

	// Implement this to get called when there's a double-click on the view. Only supported by list
	// views currently. Will be called after the double-clicked list entry has been selected.
	GetOnDoubleClick() func() error

	// Implement this to get called for any non-double-click in the view. Only supported by list
	// views currently. Will be called after the clicked list entry has been selected, and
	// HandleFocus has already been called (so the main view is up to date). Should return nil if it
	// decides not to do anything with the click.
	GetOnClick() func(opts gocui.ViewMouseBindingOpts) error

	// Implement this in a side-panel controller to get called when there's a click in the main view
	// that belongs to your panel while the main view is already focused.
	GetOnClickFocusedMainView() func(mainViewName string, clickedLineIdx int) error

	// Implement this in a side-panel controller to stage/unstage (or, later, add to
	// the custom patch) the selected diff line(s) when the user presses space in the
	// focused main view. The inclusive view-line range is the current selection (a
	// single line, a range, or a hunk). It returns the name of the focused main view
	// that should hold focus afterwards — staging/unstaging can move the acted-on
	// side to the other pane — or "" when nothing was done. Return a nil func to do
	// nothing.
	GetOnStageFocusedMainView() func(mainViewName string, firstLineIdx int, lastLineIdx int) (focusViewName string, err error)
}

type IController interface {
	HasKeybindings
	Context() Context

	GetOnRenderToMain() func()
	GetOnFocus() func(OnFocusOpts)
	GetOnFocusLost() func(OnFocusLostOpts)

	// Implement this to get called when the app quits, and the controller's context has the focus.
	// Useful for saving state on quit.
	GetOnQuit() func()
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
