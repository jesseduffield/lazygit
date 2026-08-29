package types

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	// The view that keyboard input goes to while this context is focused. That is
	// the context's own view, unless the context has an editable view embedded in
	// it which takes the keyboard instead, like the menu's filter input.
	GetInputViewName() string
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

	// true if the context holds something for a selection to sit on. Contexts that
	// don't show a selection at all say false, and so do lists with nothing in them.
	HasSelectableContent() bool
	SetHasSelectableContent(bool)

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
	// Adding on to the above, this is so that a list-specific handler can register
	// a hook for doing additional click handling
	AddOnClickFn(func(opts gocui.ViewMouseBindingOpts) error)
	// Likewise for the focused main view, which acts on the diff of whichever panel
	// is beneath it and so has to reach that panel's controller. nil for a panel
	// that shows no diff.
	AddFocusedMainViewDiffSource(FocusedMainViewDiffSource)

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
// about producing a diff between two refs for the diff menu. Implementing it is
// what makes the focused main view show a selection: a selection is only
// meaningful where there are diff lines to act on (edit one, copy some, jump by
// hunk or file). The returned type additionally classifies what acting on that
// selection means.
type DiffMainViewContext interface {
	Context

	GetDiffMainViewType() DiffMainViewType
}

// DiffMainViewType classifies what the focused main view's diff belongs to, which
// decides what acting on a selection in it means.
type DiffMainViewType int

const (
	// DiffMainViewTypeNone: the main view holds no diff, so there is nothing to
	// select. It is what a side panel that doesn't implement DiffMainViewContext
	// amounts to, not a value any panel returns.
	DiffMainViewTypeNone DiffMainViewType = iota
	// DiffMainViewTypeStaging: the diff is the working tree's, so the selection can
	// be staged or unstaged (the files panel).
	DiffMainViewTypeStaging
	// DiffMainViewTypePatchBuilding: the diff belongs to a commit, so the selection
	// can be taken into a custom patch (the commit files / commits / sub-commits /
	// reflog / stash panels).
	DiffMainViewTypePatchBuilding
)

// DiffPaneContext is one of the two panes the main section can show, as the thing
// that holds a diff with a selection in it. The panels that act on such a selection
// are handed the pane it was made in, and speak to it through this.
type DiffPaneContext interface {
	Context

	DiffSelectState() *DiffSelectState
}

// FocusedMainViewDiffSource is how a side panel hands out the diff behind what it
// renders into the focused main view: the diff of the given files as git writes it,
// with no colour and no diff renderer in the way. What the main view shows is that
// same diff after a renderer has had it, which may have restructured, reordered or
// dropped parts of it — so anything that needs the diff itself, rather than a picture
// of it, asks the panel that produced it.
//
// paths are repo-relative, and are asked for rather than assumed so that a few lines
// of a commit's diff can be had without fetching the whole thing. pane says which of
// the two main panes is asking, since a panel can show a different diff in each — the
// files panel shows the unstaged changes in one and the staged ones in the other.
type FocusedMainViewDiffSource interface {
	PlainDiff(pane DiffPaneContext, paths []string) string
}

// FocusedMainViewActions is what a side panel does when the user acts on a selection
// of diff lines in the focused main view. The main view owns the selection and the
// keys; what acting on it means is the panel's business — the working tree panel
// stages and unstages.
//
// It extends the diff source rather than standing beside it, because acting on a
// selection needs the diff behind the rendering just as reading it does; a panel that
// implements only the source offers a diff to read and copy but nothing to do to it.
type FocusedMainViewActions interface {
	FocusedMainViewDiffSource

	// PrimaryAction acts on the diff lines in the inclusive view-line range, which is
	// the current selection in the given pane: a single line, a range, or a hunk. The
	// panel re-renders the diff itself, being the one that knows what it did to it.
	PrimaryAction(pane DiffPaneContext, firstLineIdx int, lastLineIdx int) error

	// DiscardSelection takes the selected diff lines back out of whatever they are part
	// of: the working tree for the files panel, the commit itself for the panels showing
	// a commit's diff.
	DiscardSelection(pane DiffPaneContext, firstLineIdx int, lastLineIdx int) error

	// DiscardSelectionDisabledReason says why the selection can't be discarded where it
	// is, and nil when it can. Taking lines out of a commit means rewriting it, which
	// isn't always something we may do; the working tree has no such condition.
	DiscardSelectionDisabledReason(pane DiffPaneContext) *DisabledReason

	// PatchInclusion says which lines of the diff this panel shows are in the custom
	// patch being built from it, which is what the marks over those lines are drawn
	// from. nil where nothing about this diff is being built into a patch — which is
	// always so for a diff that can't be.
	PatchInclusion() func(info DiffLineInfo) bool
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
}

type OnFocusOpts struct {
	ClickedWindowName  string
	ClickedViewLineIdx int

	// Focusing a list context scrolls its selection into view. Set this to leave
	// the view's scroll position alone instead; only for callers that maintain
	// it themselves, e.g. by keeping the selection at the edge of the viewport.
	KeepScrollPosition bool

	// Set this when the focused item hasn't changed and the main view's current
	// content is still valid.
	SkipMainViewUpdate bool
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

	// Implement this in a side-panel controller to hand out the diff behind what your
	// panel renders into the focused main view, for the commands that act on a
	// selection in it. nil for a controller whose panel shows no diff.
	GetFocusedMainViewDiffSource() FocusedMainViewDiffSource
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
	IsInStack(context Context) bool
	UpdateSelectionHighlights()
	IsCurrent(c Context) bool
	IsCurrentOrParent(c Context) bool
	ForEach(func(Context))
	AllList() []IListContext
	AllFilterable() []IFilterableContext
	AllSearchable() []ISearchableContext
}
