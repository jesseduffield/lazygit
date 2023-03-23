package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ScreenModeActions struct {
	c *ControllerCommon
}

func (self *ScreenModeActions) Next() error {
	self.c.State().GetRepoState().SetScreenMode(
		nextIntInCycle(
			[]types.WindowMaximisation{types.SCREEN_NORMAL, types.SCREEN_HALF, types.SCREEN_FULL},
			self.c.State().GetRepoState().GetScreenMode(),
		),
	)

	return self.rerenderViewsWithScreenModeDependentContent()
}

func (self *ScreenModeActions) Prev() error {
	self.c.State().GetRepoState().SetScreenMode(
		prevIntInCycle(
			[]types.WindowMaximisation{types.SCREEN_NORMAL, types.SCREEN_HALF, types.SCREEN_FULL},
			self.c.State().GetRepoState().GetScreenMode(),
		),
	)

	return self.rerenderViewsWithScreenModeDependentContent()
}

// these views need to be re-rendered when the screen mode changes. The commits view,
// for example, will show authorship information in half and full screen mode.
func (self *ScreenModeActions) rerenderViewsWithScreenModeDependentContent() error {
	// for now we re-render all list views.
	for _, context := range self.c.Context().AllList() {
		if err := self.rerenderView(context.GetView()); err != nil {
			return err
		}
	}

	return nil
}

func (self *ScreenModeActions) rerenderView(view *gocui.View) error {
	context, ok := self.c.Helpers().View.ContextForView(view.Name())
	if !ok {
		self.c.Log.Errorf("no context found for view %s", view.Name())
		return nil
	}

	return context.HandleRender()
}

func nextIntInCycle(sl []types.WindowMaximisation, current types.WindowMaximisation) types.WindowMaximisation {
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

func prevIntInCycle(sl []types.WindowMaximisation, current types.WindowMaximisation) types.WindowMaximisation {
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
