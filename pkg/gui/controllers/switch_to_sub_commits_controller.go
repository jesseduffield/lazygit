package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var _ types.IController = &SwitchToSubCommitsController{}

type CanSwitchToSubCommits interface {
	types.Context
	GetSelectedRefName() string
}

type SwitchToSubCommitsController struct {
	baseController
	*controllerCommon
	context CanSwitchToSubCommits

	setSubCommits func([]*models.Commit)
}

func NewSwitchToSubCommitsController(
	controllerCommon *controllerCommon,
	setSubCommits func([]*models.Commit),
	context CanSwitchToSubCommits,
) *SwitchToSubCommitsController {
	return &SwitchToSubCommitsController{
		baseController:   baseController{},
		controllerCommon: controllerCommon,
		context:          context,
		setSubCommits:    setSubCommits,
	}
}

func (self *SwitchToSubCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Handler:     self.viewCommits,
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Description: self.c.Tr.LcViewCommits,
		},
	}

	return bindings
}

func (self *SwitchToSubCommitsController) GetOnClick() func() error {
	return self.viewCommits
}

func (self *SwitchToSubCommitsController) viewCommits() error {
	refName := self.context.GetSelectedRefName()
	if refName == "" {
		return nil
	}

	// need to populate my sub commits
	commits, err := self.git.Loaders.Commits.GetCommits(
		loaders.GetCommitsOptions{
			Limit:                true,
			FilterPath:           self.modes.Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              refName,
		},
	)
	if err != nil {
		return err
	}

	self.setSubCommits(commits)
	self.contexts.SubCommits.SetSelectedLineIdx(0)
	self.contexts.SubCommits.SetParentContext(self.context)

	return self.c.PushContext(self.contexts.SubCommits)
}

func (self *SwitchToSubCommitsController) Context() types.Context {
	return self.context
}
