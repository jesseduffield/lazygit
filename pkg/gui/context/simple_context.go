package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SimpleContext struct {
	*BaseContext
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

func (self *SimpleContext) HandleFocus(opts types.OnFocusOpts) error {
	if self.highlightOnFocus {
		self.GetViewTrait().SetHighlight(true)
	}

	if self.onFocusFn != nil {
		if err := self.onFocusFn(opts); err != nil {
			return err
		}
	}

	if self.onRenderToMainFn != nil {
		if err := self.onRenderToMainFn(); err != nil {
			return err
		}
	}

	return nil
}

func (self *SimpleContext) HandleFocusLost(opts types.OnFocusLostOpts) error {
	self.GetViewTrait().SetHighlight(false)
	_ = self.view.SetOriginX(0)
	if self.onFocusLostFn != nil {
		return self.onFocusLostFn(opts)
	}
	return nil
}

func (self *SimpleContext) HandleRender() error {
	return nil
}

func (self *SimpleContext) HandleRenderToMain() error {
	if self.onRenderToMainFn != nil {
		return self.onRenderToMainFn()
	}

	return nil
}
