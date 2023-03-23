package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ViewHelper struct {
	c *HelperCommon
}

func NewViewHelper(c *HelperCommon, contexts *context.ContextTree) *ViewHelper {
	return &ViewHelper{
		c: c,
	}
}

func (self *ViewHelper) ContextForView(viewName string) (types.Context, bool) {
	view, err := self.c.GocuiGui().View(viewName)
	if err != nil {
		return nil, false
	}

	for _, context := range self.c.Contexts().Flatten() {
		if context.GetViewName() == view.Name() {
			return context, true
		}
	}

	return nil, false
}
