package controllers

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type JumpToSideWindowController struct {
	baseController
	c *ControllerCommon
}

func NewJumpToSideWindowController(
	c *ControllerCommon,
) *JumpToSideWindowController {
	return &JumpToSideWindowController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *JumpToSideWindowController) Context() types.Context {
	return nil
}

func (self *JumpToSideWindowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	windows := self.c.Helpers().Window.SideWindows()

	if len(opts.Config.Universal.JumpToBlock) != len(windows) {
		log.Fatal("Jump to block keybindings cannot be set. Exactly 5 keybindings must be supplied.")
	}

	return lo.Map(windows, func(window string, index int) *types.Binding {
		return &types.Binding{
			ViewName: "",
			// by default the keys are 1, 2, 3, etc
			Key:      opts.GetKey(opts.Config.Universal.JumpToBlock[index]),
			Modifier: gocui.ModNone,
			Handler:  self.goToSideWindow(window),
		}
	})
}

func (self *JumpToSideWindowController) goToSideWindow(window string) func() error {
	return func() error {
		context := self.c.Helpers().Window.GetContextForWindow(window)

		return self.c.PushContext(context)
	}
}
