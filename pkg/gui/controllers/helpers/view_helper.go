package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ViewHelper struct {
	c        *types.HelperCommon
	contexts *context.ContextTree
}

func NewViewHelper(c *types.HelperCommon, contexts *context.ContextTree) *ViewHelper {
	return &ViewHelper{
		c:        c,
		contexts: contexts,
	}
}

func (self *ViewHelper) ContextForView(viewName string) (types.Context, bool) {
	view, err := self.c.GocuiGui().View(viewName)
	if err != nil {
		return nil, false
	}

	for _, context := range self.contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return context, true
		}
	}

	return nil, false
}
