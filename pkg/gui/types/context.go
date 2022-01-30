package types

import "github.com/jesseduffield/lazygit/pkg/config"

type ContextKind int

const (
	SIDE_CONTEXT ContextKind = iota
	MAIN_CONTEXT
	TEMPORARY_POPUP
	PERSISTENT_POPUP
	EXTRAS_CONTEXT
)

type ParentContexter interface {
	SetParentContext(Context)
	// we return a bool here to tell us whether or not the returned value just wraps a nil
	GetParentContext() (Context, bool)
}

type IBaseContext interface {
	ParentContexter

	GetKind() ContextKind
	GetViewName() string
	GetWindowName() string
	SetWindowName(string)
	GetKey() ContextKey

	GetOptionsMap() map[string]string
}

type Context interface {
	IBaseContext

	HandleFocus(opts ...OnFocusOpts) error
	HandleFocusLost() error
	HandleRender() error
	HandleRenderToMain() error
}

type OnFocusOpts struct {
	ClickedViewName    string
	ClickedViewLineIdx int
}

type ContextKey string

type HasKeybindings interface {
	Keybindings(
		getKey func(key string) interface{},
		config config.KeybindingConfig,
		guards KeybindingGuards,
	) []*Binding
}

type IController interface {
	HasKeybindings
	Context() Context
}

type IListContext interface {
	HasKeybindings

	GetSelectedItemId() string
	HandlePrevLine() error
	HandleNextLine() error
	HandleScrollLeft() error
	HandleScrollRight() error
	HandlePrevPage() error
	HandleNextPage() error
	HandleGotoTop() error
	HandleGotoBottom() error
	HandleClick(onClick func() error) error

	OnSearchSelect(selectedLineIdx int) error
	FocusLine()

	GetPanelState() IListPanelState

	Context
}

type IList interface {
	IListCursor
	GetItemsLength() int
}

type IListCursor interface {
	GetSelectedLineIdx() int
	SetSelectedLineIdx(value int)
	MoveSelectedLine(delta int)
	RefreshSelectedIdx()
}

type IListPanelState interface {
	SetSelectedLineIdx(int)
	GetSelectedLineIdx() int
}

type ListItem interface {
	// ID is a SHA when the item is a commit, a filename when the item is a file, 'stash@{4}' when it's a stash entry, 'my_branch' when it's a branch
	ID() string

	// Description is something we would show in a message e.g. '123as14: push blah' for a commit
	Description() string
}
