package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsSwitchControllerFactory struct {
	controllerCommon *controllerCommon
	setSubCommits    func([]*models.Commit)
}

var _ types.IController = &SubCommitsSwitchController{}

type ContextWithRefName interface {
	types.Context
	GetSelectedRefName() string
}

type SubCommitsSwitchController struct {
	baseController
	*controllerCommon
	context ContextWithRefName

	setSubCommits func([]*models.Commit)
}

func NewSubCommitsSwitchControllerFactory(
	common *controllerCommon,
	setSubCommits func([]*models.Commit),
) *SubCommitsSwitchControllerFactory {
	return &SubCommitsSwitchControllerFactory{
		controllerCommon: common,
		setSubCommits:    setSubCommits,
	}
}

func (self *SubCommitsSwitchControllerFactory) Create(context ContextWithRefName) *SubCommitsSwitchController {
	return &SubCommitsSwitchController{
		baseController:   baseController{},
		controllerCommon: self.controllerCommon,
		context:          context,
		setSubCommits:    self.setSubCommits,
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

func (self *SubCommitsSwitchController) GetOnClick() func() error {
	return self.viewCommits
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
	self.contexts.SubCommits.SetSelectedLineIdx(0)
	self.contexts.SubCommits.SetParentContext(self.context)

	return self.c.PushContext(self.contexts.SubCommits)
}

func (self *SubCommitsSwitchController) Context() types.Context {
	return self.context
}
