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
	// if a context is transient, then it only appears via some keybinding on another
	// context. Until we add support for having multiple of the same context, no two
	// of the same transient context can appear at once meaning one might be 'stolen'
	// from another window.
	IsTransient() bool

	// returns the desired title for the view upon activation. If there is no desired title (returns empty string), then
	// no title will be set
	Title() string

	GetOptionsMap() map[string]string

	AddKeybindingsFn(KeybindingsFn)
	AddMouseKeybindingsFn(MouseKeybindingsFn)

	// This is a bit of a hack at the moment: we currently only set an onclick function so that
	// our list controller can come along and wrap it in a list-specific click handler.
	// We'll need to think of a better way to do this.
	AddOnClickFn(func() error)
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
	ScrollUp()
	ScrollDown()
	PageDelta() int
	SelectedLineIdx() int
}

type OnFocusOpts struct {
	ClickedViewName    string
	ClickedViewLineIdx int
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
}

type IController interface {
	HasKeybindings
	Context() Context
}

type IList interface {
	IListCursor
	Len() int
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
