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
	c       *ControllerCommon
	context CanSwitchToSubCommits

	setSubCommits func([]*models.Commit)
}

func NewSwitchToSubCommitsController(
	controllerCommon *ControllerCommon,
	setSubCommits func([]*models.Commit),
	context CanSwitchToSubCommits,
) *SwitchToSubCommitsController {
	return &SwitchToSubCommitsController{
		baseController: baseController{},
		c:              controllerCommon,
		context:        context,
		setSubCommits:  setSubCommits,
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

	// need to populate my sub commits
	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                true,
			FilterPath:           self.c.Modes().Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              ref.FullRefName(),
		},
	)
	if err != nil {
		return err
	}

	self.setSubCommits(commits)

	subCommitsContext := self.c.Contexts().SubCommits
	subCommitsContext.SetSelectedLineIdx(0)
	subCommitsContext.SetParentContext(self.context)
	subCommitsContext.SetWindowName(self.context.GetWindowName())
	subCommitsContext.SetTitleRef(ref.Description())
	subCommitsContext.SetRef(ref)
	subCommitsContext.SetLimitCommits(true)
	subCommitsContext.ClearSearchString()
	subCommitsContext.GetView().ClearSearch()

	err = self.c.PostRefreshUpdate(self.c.Contexts().SubCommits)
	if err != nil {
		return err
	}

	return self.c.PushContext(self.c.Contexts().SubCommits)
}

func (self *SwitchToSubCommitsController) Context() types.Context {
	return self.context
}
