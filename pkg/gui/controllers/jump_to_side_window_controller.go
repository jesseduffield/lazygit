package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
	jumpKeys := opts.Config.Universal.JumpToBlock

	// Assign jump keys to panels positionally (by default 1 to the first panel,
	// 2 to the second, etc.), for as many panels as there are keys. If there are
	// more panels than keys the extra panels just have no jump key, and if there
	// are more keys than panels the extra keys are unused; either way panels stay
	// reachable via the next/previous-panel keys.
	count := min(len(windows), len(jumpKeys))
	bindings := make([]*types.Binding, 0, count)
	for i := range count {
		bindings = append(bindings, &types.Binding{
			ViewName: "",
			Keys:     opts.GetKeys(jumpKeys[i]),
			Handler:  opts.Guards.NoPopupPanel(self.goToSideWindow(windows[i])),
		})
	}
	return bindings
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
