package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BaseContext struct {
	kind            types.ContextKind
	key             types.ContextKey
	view            *gocui.View
	viewTrait       types.IViewTrait
	windowName      string
	onGetOptionsMap func() map[string]string

	keybindingsFns      []types.KeybindingsFn
	mouseKeybindingsFns []types.MouseKeybindingsFn
	onClickFn           func() error
	onRenderToMainFn    func() error
	onFocusFn           onFocusFn
	onFocusLostFn       onFocusLostFn

	focusable           bool
	transient           bool
	hasControlledBounds bool
	highlightOnFocus    bool

	*ParentContextMgr
}

type (
	onFocusFn     = func(types.OnFocusOpts) error
	onFocusLostFn = func(types.OnFocusLostOpts) error
)

var _ types.IBaseContext = &BaseContext{}

type NewBaseContextOpts struct {
	Kind                  types.ContextKind
	Key                   types.ContextKey
	View                  *gocui.View
	WindowName            string
	Focusable             bool
	Transient             bool
	HasUncontrolledBounds bool // negating for the sake of making false the default
	HighlightOnFocus      bool

	OnGetOptionsMap func() map[string]string
}

func NewBaseContext(opts NewBaseContextOpts) *BaseContext {
	viewTrait := NewViewTrait(opts.View)

	hasControlledBounds := !opts.HasUncontrolledBounds

	return &BaseContext{
		kind:                opts.Kind,
		key:                 opts.Key,
		view:                opts.View,
		windowName:          opts.WindowName,
		onGetOptionsMap:     opts.OnGetOptionsMap,
		focusable:           opts.Focusable,
		transient:           opts.Transient,
		hasControlledBounds: hasControlledBounds,
		highlightOnFocus:    opts.HighlightOnFocus,
		ParentContextMgr:    &ParentContextMgr{},
		viewTrait:           viewTrait,
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
	// for the sake of the global context which has no view
	if self.view == nil {
		return ""
	}

	return self.view.Name()
}

func (self *BaseContext) GetView() *gocui.View {
	return self.view
}

func (self *BaseContext) GetViewTrait() types.IViewTrait {
	return self.viewTrait
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

func (self *BaseContext) AddOnRenderToMainFn(fn func() error) {
	if fn != nil {
		self.onRenderToMainFn = fn
	}
}

func (self *BaseContext) GetOnRenderToMain() func() error {
	return self.onRenderToMainFn
}

func (self *BaseContext) AddOnFocusFn(fn onFocusFn) {
	if fn != nil {
		self.onFocusFn = fn
	}
}

func (self *BaseContext) GetOnFocus() onFocusFn {
	return self.onFocusFn
}

func (self *BaseContext) AddOnFocusLostFn(fn onFocusLostFn) {
	if fn != nil {
		self.onFocusLostFn = fn
	}
}

func (self *BaseContext) GetOnFocusLost() onFocusLostFn {
	return self.onFocusLostFn
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

func (self *BaseContext) HasControlledBounds() bool {
	return self.hasControlledBounds
}

func (self *BaseContext) Title() string {
	return ""
}
