package controllers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type JumpToSideWindowController struct {
	baseController
	c           *ControllerCommon
	nextTabFunc func() error
}

func NewJumpToSideWindowController(
	c *ControllerCommon,
	nextTabFunc func() error,
) *JumpToSideWindowController {
	return &JumpToSideWindowController{
		baseController: baseController{},
		c:              c,
		nextTabFunc:    nextTabFunc,
	}
}

func (self *JumpToSideWindowController) Context() types.Context {
	return nil
}

func (self *JumpToSideWindowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	windows := self.c.Helpers().Window.SideWindows()
	jumpBindings := opts.Config.Universal.JumpToBlock

	// Auto-extend jump bindings if there are more windows than bindings
	for len(jumpBindings) < len(windows) {
		jumpBindings = append(jumpBindings, fmt.Sprintf("%d", len(jumpBindings)+1))
	}

	return lo.Map(windows, func(window string, index int) *types.Binding {
		return &types.Binding{
			ViewName: "",
			// by default the keys are 1, 2, 3, etc
			Key:      opts.GetKey(jumpBindings[index]),
			Modifier: gocui.ModNone,
			Handler:  opts.Guards.NoPopupPanel(self.goToSideWindow(window)),
		}
	})
}

func (self *JumpToSideWindowController) goToSideWindow(window string) func() error {
	return func() error {
		sideWindowAlreadyActive := self.c.Helpers().Window.CurrentWindow() == window
		if sideWindowAlreadyActive && self.c.UserConfig().Gui.SwitchTabsWithPanelJumpKeys {
			return self.nextTabFunc()
		}

		context := self.c.Helpers().Window.GetContextForWindow(window)

		self.c.Context().Push(context, types.OnFocusOpts{})
		return nil
	}
}
