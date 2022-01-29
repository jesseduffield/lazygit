package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type ParentContextMgr struct {
	ParentContext types.Context
	// we can't know on the calling end whether a Context is actually a nil value without reflection, so we're storing this flag here to tell us. There has got to be a better way around this
	hasParent bool
}

var _ types.ParentContexter = (*ParentContextMgr)(nil)

func (self *ParentContextMgr) SetParentContext(context types.Context) {
	self.ParentContext = context
	self.hasParent = true
}

func (self *ParentContextMgr) GetParentContext() (types.Context, bool) {
	return self.ParentContext, self.hasParent
}
