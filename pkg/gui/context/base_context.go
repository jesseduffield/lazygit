package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BaseContext struct {
	kind            types.ContextKind
	key             types.ContextKey
	ViewName        string
	windowName      string
	onGetOptionsMap func() map[string]string

	keybindingsFns      []types.KeybindingsFn
	mouseKeybindingsFns []types.MouseKeybindingsFn
	onClickFn           func() error

	focusable bool
	transient bool

	*ParentContextMgr
}

var _ types.IBaseContext = &BaseContext{}

type NewBaseContextOpts struct {
	Kind       types.ContextKind
	Key        types.ContextKey
	ViewName   string
	WindowName string
	Focusable  bool
	Transient  bool

	OnGetOptionsMap func() map[string]string
}

func NewBaseContext(opts NewBaseContextOpts) *BaseContext {
	return &BaseContext{
		kind:             opts.Kind,
		key:              opts.Key,
		ViewName:         opts.ViewName,
		windowName:       opts.WindowName,
		onGetOptionsMap:  opts.OnGetOptionsMap,
		focusable:        opts.Focusable,
		transient:        opts.Transient,
		ParentContextMgr: &ParentContextMgr{},
	}
}

func (self *BaseContext) GetOptionsMap() map[string]string {
	if self.onGetOptionsMap != nil {
		return self.onGetOptionsMap()
	}
	return nil
}

func (self *BaseContext) SetWindowName(windowName string) {
	self.windowName = windowName
}

func (self *BaseContext) GetWindowName() string {
	return self.windowName
}

func (self *BaseContext) GetViewName() string {
	return self.ViewName
}

func (self *BaseContext) GetKind() types.ContextKind {
	return self.kind
}

func (self *BaseContext) GetKey() types.ContextKey {
	return self.key
}

func (self *BaseContext) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{}
	for i := range self.keybindingsFns {
		// the first binding in the bindings array takes precedence but we want the
		// last keybindingsFn to take precedence to we add them in reverse
		bindings = append(bindings, self.keybindingsFns[len(self.keybindingsFns)-1-i](opts)...)
	}

	return bindings
}

func (self *BaseContext) AddKeybindingsFn(fn types.KeybindingsFn) {
	self.keybindingsFns = append(self.keybindingsFns, fn)
}

func (self *BaseContext) AddMouseKeybindingsFn(fn types.MouseKeybindingsFn) {
	self.mouseKeybindingsFns = append(self.mouseKeybindingsFns, fn)
}

func (self *BaseContext) AddOnClickFn(fn func() error) {
	if fn != nil {
		self.onClickFn = fn
	}
}

func (self *BaseContext) GetOnClick() func() error {
	return self.onClickFn
}

func (self *BaseContext) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	bindings := []*gocui.ViewMouseBinding{}
	for i := range self.mouseKeybindingsFns {
		// the first binding in the bindings array takes precedence but we want the
		// last keybindingsFn to take precedence to we add them in reverse
		bindings = append(bindings, self.mouseKeybindingsFns[len(self.mouseKeybindingsFns)-1-i](opts)...)
	}

	return bindings
}

func (self *BaseContext) IsFocusable() bool {
	return self.focusable
}

func (self *BaseContext) IsTransient() bool {
	return self.transient
}

func (self *BaseContext) Title() string {
	return ""
}
