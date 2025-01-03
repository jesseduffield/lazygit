package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PanelSizeActions struct {
	c *ControllerCommon
}

func (self *PanelSizeActions) Next() error {
	self.c.State().GetRepoState().SetPanelSize(
		nextIntInCycle(
			[]types.PanelSize{types.PANEL_SIZE_NORMAL, types.PANEL_SIZE_HALF, types.PANEL_SIZE_FULL},
			self.c.State().GetRepoState().GetPanelSize(),
		),
	)

	self.rerenderViewsWithPanelSizeDependentContent()
	return nil
}

func (self *PanelSizeActions) Prev() error {
	self.c.State().GetRepoState().SetPanelSize(
		prevIntInCycle(
			[]types.PanelSize{types.PANEL_SIZE_NORMAL, types.PANEL_SIZE_HALF, types.PANEL_SIZE_FULL},
			self.c.State().GetRepoState().GetPanelSize(),
		),
	)

	self.rerenderViewsWithPanelSizeDependentContent()
	return nil
}

// these views need to be re-rendered when the panel size changes. The commits view,
// for example, will show authorship information in half and full panel size.
func (self *PanelSizeActions) rerenderViewsWithPanelSizeDependentContent() {
	for _, context := range self.c.Context().AllList() {
		if context.NeedsRerenderOnWidthChange() == types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_PANEL_SIZE_CHANGES {
			self.rerenderView(context.GetView())
		}
	}
}

func (self *PanelSizeActions) rerenderView(view *gocui.View) {
	context, ok := self.c.Helpers().View.ContextForView(view.Name())
	if !ok {
		self.c.Log.Errorf("no context found for view %s", view.Name())
		return
	}

	context.HandleRender()
}

func nextIntInCycle(sl []types.PanelSize, current types.PanelSize) types.PanelSize {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

func prevIntInCycle(sl []types.PanelSize, current types.PanelSize) types.PanelSize {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}
