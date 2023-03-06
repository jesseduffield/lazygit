package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

var _ types.IController = &SwitchToSubCommitsController{}

type CanSwitchToSubCommits interface {
	types.Context
	GetSelectedRef() types.Ref
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
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return nil
	}

	// need to populate my sub commits
	commits, err := self.git.Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                true,
			FilterPath:           self.modes.Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              ref.FullRefName(),
		},
	)
	if err != nil {
		return err
	}

	self.setSubCommits(commits)

	self.contexts.SubCommits.SetSelectedLineIdx(0)
	self.contexts.SubCommits.SetParentContext(self.context)
	self.contexts.SubCommits.SetWindowName(self.context.GetWindowName())
	self.contexts.SubCommits.SetTitleRef(ref.Description())
	self.contexts.SubCommits.SetRef(ref)
	self.contexts.SubCommits.SetLimitCommits(true)

	err = self.c.PostRefreshUpdate(self.contexts.SubCommits)
	if err != nil {
		return err
	}

	return self.c.PushContext(self.contexts.SubCommits)
}

func (self *SwitchToSubCommitsController) Context() types.Context {
	return self.context
}
