package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type ParentContextMgr struct {
	ParentContext types.Context
}

var _ types.ParentContexter = (*ParentContextMgr)(nil)

func (self *ParentContextMgr) SetParentContext(context types.Context) {
	self.ParentContext = context
}

func (self *ParentContextMgr) GetParentContext() types.Context {
	return self.ParentContext
}
