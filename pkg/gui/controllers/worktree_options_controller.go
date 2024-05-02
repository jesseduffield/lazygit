package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that have items you can create a worktree from

var _ types.IController = &WorktreeOptionsController{}

type CanViewWorktreeOptions interface {
	types.IListContext
}

type WorktreeOptionsController struct {
	baseController
	*ListControllerTrait[string]
	c       *ControllerCommon
	context CanViewWorktreeOptions
}

func NewWorktreeOptionsController(c *ControllerCommon, context CanViewWorktreeOptions) *WorktreeOptionsController {
	return &WorktreeOptionsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[string](
			c,
			context,
			context.GetSelectedItemId,
			context.GetSelectedItemIds,
		),
		c:       c,
		context: context,
	}
}

func (self *WorktreeOptionsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Worktrees.ViewWorktreeOptions),
			Handler:     self.withItem(self.viewWorktreeOptions),
			Description: self.c.Tr.ViewWorktreeOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *WorktreeOptionsController) viewWorktreeOptions(ref string) error {
	return self.c.Helpers().Worktree.ViewWorktreeOptions(self.context, ref)
}
