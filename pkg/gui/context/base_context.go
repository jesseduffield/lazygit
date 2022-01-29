package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type BaseContext struct {
	Kind            types.ContextKind
	Key             types.ContextKey
	ViewName        string
	WindowName      string
	OnGetOptionsMap func() map[string]string

	*ParentContextMgr
}

func (self *BaseContext) GetOptionsMap() map[string]string {
	if self.OnGetOptionsMap != nil {
		return self.OnGetOptionsMap()
	}
	return nil
}

func (self *BaseContext) SetWindowName(windowName string) {
	self.WindowName = windowName
}

func (self *BaseContext) GetWindowName() string {
	windowName := self.WindowName

	if windowName != "" {
		return windowName
	}

	// TODO: actually set this for everything so we don't default to the view name
	return self.ViewName
}

func (self *BaseContext) GetViewName() string {
	return self.ViewName
}

func (self *BaseContext) GetKind() types.ContextKind {
	return self.Kind
}

func (self *BaseContext) GetKey() types.ContextKey {
	return self.Key
}

type NewBaseContextOpts struct {
	Kind       types.ContextKind
	Key        types.ContextKey
	ViewName   string
	WindowName string

	OnGetOptionsMap func() map[string]string
}

func NewBaseContext(opts NewBaseContextOpts) *BaseContext {
	return &BaseContext{
		Kind:             opts.Kind,
		Key:              opts.Key,
		ViewName:         opts.ViewName,
		WindowName:       opts.WindowName,
		OnGetOptionsMap:  opts.OnGetOptionsMap,
		ParentContextMgr: &ParentContextMgr{},
	}
}
