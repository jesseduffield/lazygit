package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SimpleContext struct {
	OnFocus     func(opts ...types.OnFocusOpts) error
	OnFocusLost func() error
	OnRender    func() error
	// this is for pushing some content to the main view
	OnRenderToMain func(opts ...types.OnFocusOpts) error

	*BaseContext
}

type ContextCallbackOpts struct {
	OnFocus     func(opts ...types.OnFocusOpts) error
	OnFocusLost func() error
	OnRender    func() error
	// this is for pushing some content to the main view
	OnRenderToMain func(opts ...types.OnFocusOpts) error
}

func NewSimpleContext(baseContext *BaseContext, opts ContextCallbackOpts) *SimpleContext {
	return &SimpleContext{
		OnFocus:        opts.OnFocus,
		OnFocusLost:    opts.OnFocusLost,
		OnRender:       opts.OnRender,
		OnRenderToMain: opts.OnRenderToMain,
		BaseContext:    baseContext,
	}
}

var _ types.Context = &SimpleContext{}

func (self *SimpleContext) HandleFocus(opts ...types.OnFocusOpts) error {
	if self.OnFocus != nil {
		if err := self.OnFocus(opts...); err != nil {
			return err
		}
	}

	if self.OnRenderToMain != nil {
		if err := self.OnRenderToMain(opts...); err != nil {
			return err
		}
	}

	return nil
}

func (self *SimpleContext) HandleFocusLost() error {
	if self.OnFocusLost != nil {
		return self.OnFocusLost()
	}
	return nil
}

func (self *SimpleContext) HandleRender() error {
	if self.OnRender != nil {
		return self.OnRender()
	}
	return nil
}

func (self *SimpleContext) HandleRenderToMain() error {
	if self.OnRenderToMain != nil {
		return self.OnRenderToMain()
	}

	return nil
}
