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
	c       *ControllerCommon
	context CanViewWorktreeOptions
}

func NewWorktreeOptionsController(controllerCommon *ControllerCommon, context CanViewWorktreeOptions) *WorktreeOptionsController {
	return &WorktreeOptionsController{
		baseController: baseController{},
		c:              controllerCommon,
		context:        context,
	}
}

func (self *WorktreeOptionsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Worktrees.ViewWorktreeOptions),
			Handler:     self.checkSelected(self.viewWorktreeOptions),
			Description: self.c.Tr.ViewWorktreeOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *WorktreeOptionsController) checkSelected(callback func(string) error) func() error {
	return func() error {
		ref := self.context.GetSelectedItemId()
		if ref == "" {
			return nil
		}

		return callback(ref)
	}
}

func (self *WorktreeOptionsController) Context() types.Context {
	return self.context
}

func (self *WorktreeOptionsController) viewWorktreeOptions(ref string) error {
	return self.c.Helpers().Worktree.ViewWorktreeOptions(self.context, ref)
}
