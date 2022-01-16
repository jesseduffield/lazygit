package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BasicContext struct {
	OnFocus     func(opts ...types.OnFocusOpts) error
	OnFocusLost func() error
	OnRender    func() error
	// this is for pushing some content to the main view
	OnRenderToMain  func(opts ...types.OnFocusOpts) error
	Kind            types.ContextKind
	Key             types.ContextKey
	ViewName        string
	WindowName      string
	OnGetOptionsMap func() map[string]string

	ParentContext types.Context
	// we can't know on the calling end whether a Context is actually a nil value without reflection, so we're storing this flag here to tell us. There has got to be a better way around this
	hasParent bool
}

var _ types.Context = &BasicContext{}

func (self *BasicContext) GetOptionsMap() map[string]string {
	if self.OnGetOptionsMap != nil {
		return self.OnGetOptionsMap()
	}
	return nil
}

func (self *BasicContext) SetParentContext(context types.Context) {
	self.ParentContext = context
	self.hasParent = true
}

func (self *BasicContext) GetParentContext() (types.Context, bool) {
	return self.ParentContext, self.hasParent
}

func (self *BasicContext) SetWindowName(windowName string) {
	self.WindowName = windowName
}

func (self *BasicContext) GetWindowName() string {
	windowName := self.WindowName

	if windowName != "" {
		return windowName
	}

	// TODO: actually set this for everything so we don't default to the view name
	return self.ViewName
}

func (self *BasicContext) HandleRender() error {
	if self.OnRender != nil {
		return self.OnRender()
	}
	return nil
}

func (self *BasicContext) GetViewName() string {
	return self.ViewName
}

func (self *BasicContext) HandleFocus(opts ...types.OnFocusOpts) error {
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

func (self *BasicContext) HandleFocusLost() error {
	if self.OnFocusLost != nil {
		return self.OnFocusLost()
	}
	return nil
}

func (self *BasicContext) HandleRenderToMain() error {
	if self.OnRenderToMain != nil {
		return self.OnRenderToMain()
	}

	return nil
}

func (self *BasicContext) GetKind() types.ContextKind {
	return self.Kind
}

func (self *BasicContext) GetKey() types.ContextKey {
	return self.Key
}
