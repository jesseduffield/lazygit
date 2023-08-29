package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var _ types.IController = &SwitchToSubCommitsController{}

type CanSwitchToSubCommits interface {
	types.Context
	GetSelectedRef() types.Ref
	ShowBranchHeadsInSubCommits() bool
}

type SwitchToSubCommitsController struct {
	baseController
	c       *ControllerCommon
	context CanSwitchToSubCommits
}

func NewSwitchToSubCommitsController(
	controllerCommon *ControllerCommon,
	context CanSwitchToSubCommits,
) *SwitchToSubCommitsController {
	return &SwitchToSubCommitsController{
		baseController: baseController{},
		c:              controllerCommon,
		context:        context,
	}
}

func (self *SwitchToSubCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Handler:     self.viewCommits,
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Description: self.c.Tr.ViewCommits,
		},
	}

	return bindings
}

func (self *SwitchToSubCommitsController) GetOnClick() func() error {
	return self.viewCommits
}

func (self *SwitchToSubCommitsController) viewCommits() error {
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return nil
	}

	return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
		Ref:             ref,
		TitleRef:        ref.RefName(),
		Context:         self.context,
		ShowBranchHeads: self.context.ShowBranchHeadsInSubCommits(),
	})
}

func (self *SwitchToSubCommitsController) Context() types.Context {
	return self.context
}
