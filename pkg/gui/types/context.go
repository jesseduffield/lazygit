package types

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
)

type ContextKind int

const (
	SIDE_CONTEXT ContextKind = iota
	MAIN_CONTEXT
	TEMPORARY_POPUP
	PERSISTENT_POPUP
	EXTRAS_CONTEXT
	// only used by the one global context
	GLOBAL_CONTEXT
)

type ParentContexter interface {
	SetParentContext(Context)
	// we return a bool here to tell us whether or not the returned value just wraps a nil
	GetParentContext() (Context, bool)
}

type IBaseContext interface {
	HasKeybindings
	ParentContexter

	GetKind() ContextKind
	GetViewName() string
	GetWindowName() string
	SetWindowName(string)
	GetKey() ContextKey
	IsFocusable() bool

	GetOptionsMap() map[string]string

	AddKeybindingsFn(KeybindingsFn)
	AddMouseKeybindingsFn(MouseKeybindingsFn)
}

type Context interface {
	IBaseContext

	HandleFocus(opts ...OnFocusOpts) error
	HandleFocusLost() error
	HandleRender() error
	HandleRenderToMain() error
}

type IListContext interface {
	Context

	GetSelectedItemId() string

	GetList() IList

	OnSearchSelect(selectedLineIdx int) error
	FocusLine()

	GetPanelState() IListPanelState
	GetViewTrait() IViewTrait
}

type IViewTrait interface {
	FocusPoint(yIdx int)
	SetViewPortContent(content string)
	SetContent(content string)
	SetFooter(value string)
	SetOriginX(value int)
	ViewPortYBounds() (int, int)
	ScrollLeft()
	ScrollRight()
	PageDelta() int
	SelectedLineIdx() int
}

type OnFocusOpts struct {
	ClickedViewName    string
	ClickedViewLineIdx int
}

type ContextKey string

type KeybindingsOpts struct {
	GetKey func(key string) interface{}
	Config config.KeybindingConfig
	Guards KeybindingGuards
}

type KeybindingsFn func(opts KeybindingsOpts) []*Binding
type MouseKeybindingsFn func(opts KeybindingsOpts) []*gocui.ViewMouseBinding

type HasKeybindings interface {
	GetKeybindings(opts KeybindingsOpts) []*Binding
	GetMouseKeybindings(opts KeybindingsOpts) []*gocui.ViewMouseBinding
}

type IController interface {
	HasKeybindings
	Context() Context
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
