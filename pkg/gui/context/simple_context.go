package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SimpleContext struct {
	*BaseContext
	handleRenderFunc func()
}

func NewSimpleContext(baseContext *BaseContext) *SimpleContext {
	return &SimpleContext{
		BaseContext: baseContext,
	}
}

var _ types.Context = &SimpleContext{}

// A Display context only renders a view. It has no keybindings and is not focusable.
func NewDisplayContext(key types.ContextKey, view *gocui.View, windowName string) types.Context {
	return NewSimpleContext(
		NewBaseContext(NewBaseContextOpts{
			Kind:       types.DISPLAY_CONTEXT,
			Key:        key,
			View:       view,
			WindowName: windowName,
			Focusable:  false,
			Transient:  false,
		}),
	)
}

func (self *SimpleContext) HandleFocus(opts types.OnFocusOpts) {
	if self.highlightOnFocus {
		self.GetViewTrait().SetHighlight(true)
	}

	for _, fn := range self.onFocusFns {
		fn(opts)
	}

	if self.onRenderToMainFn != nil {
		self.onRenderToMainFn()
	}
}

func (self *SimpleContext) HandleFocusLost(opts types.OnFocusLostOpts) {
	self.GetViewTrait().SetHighlight(false)
	self.view.SetOriginX(0)
	for _, fn := range self.onFocusLostFns {
		fn(opts)
	}
}

func (self *SimpleContext) FocusLine() {
}

func (self *SimpleContext) HandleRender() {
	if self.handleRenderFunc != nil {
		self.handleRenderFunc()
	}
}

func (self *SimpleContext) SetHandleRenderFunc(f func()) {
	self.handleRenderFunc = f
}

func (self *SimpleContext) HandleRenderToMain() {
	if self.onRenderToMainFn != nil {
		self.onRenderToMainFn()
	}
}
