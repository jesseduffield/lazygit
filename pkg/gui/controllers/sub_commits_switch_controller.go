package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsSwitchControllerFactory struct {
	c                 *types.ControllerCommon
	subCommitsContext *context.SubCommitsContext
	git               *commands.GitCommand
	modes             *types.Modes
	setSubCommits     func([]*models.Commit)
}

var _ types.IController = &SubCommitsSwitchController{}

type ContextWithRefName interface {
	types.Context
	GetSelectedRefName() string
}

type SubCommitsSwitchController struct {
	baseController

	c                 *types.ControllerCommon
	context           ContextWithRefName
	subCommitsContext *context.SubCommitsContext
	git               *commands.GitCommand
	modes             *types.Modes
	setSubCommits     func([]*models.Commit)
}

func NewSubCommitsSwitchControllerFactory(
	c *types.ControllerCommon,
	subCommitsContext *context.SubCommitsContext,
	git *commands.GitCommand,
	modes *types.Modes,
	setSubCommits func([]*models.Commit),
) *SubCommitsSwitchControllerFactory {
	return &SubCommitsSwitchControllerFactory{
		c:                 c,
		subCommitsContext: subCommitsContext,
		git:               git,
		modes:             modes,
		setSubCommits:     setSubCommits,
	}
}

func (self *SubCommitsSwitchControllerFactory) Create(context ContextWithRefName) *SubCommitsSwitchController {
	return &SubCommitsSwitchController{
		baseController:    baseController{},
		c:                 self.c,
		context:           context,
		subCommitsContext: self.subCommitsContext,
		git:               self.git,
		modes:             self.modes,
		setSubCommits:     self.setSubCommits,
	}
}

func (self *SubCommitsSwitchController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Handler:     self.viewCommits,
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Description: self.c.Tr.LcViewCommits,
		},
	}

	return bindings
}

func (self *SubCommitsSwitchController) viewCommits() error {
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
	self.subCommitsContext.SetSelectedLineIdx(0)
	self.subCommitsContext.SetParentContext(self.context)

	return self.c.PushContext(self.subCommitsContext)
}

func (self *SubCommitsSwitchController) Context() types.Context {
	return self.context
}
