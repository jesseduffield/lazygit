package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommonCommitControllerFactory struct {
	controllerCommon *controllerCommon
	viewFiles        func(SwitchToCommitFilesContextOpts) error
}

var _ types.IController = &CommonCommitController{}

type CommitContext interface {
	types.Context
	CanRebase() bool
	GetSelected() *models.Commit
}

type CommonCommitController struct {
	baseController
	*controllerCommon
	context CommitContext

	viewFiles func(SwitchToCommitFilesContextOpts) error
}

func NewCommonCommitControllerFactory(
	common *controllerCommon,
	viewFiles func(SwitchToCommitFilesContextOpts) error,
) *CommonCommitControllerFactory {
	return &CommonCommitControllerFactory{
		controllerCommon: common,
		viewFiles:        viewFiles,
	}
}

func (self *CommonCommitControllerFactory) Create(context CommitContext) *CommonCommitController {
	return &CommonCommitController{
		baseController:   baseController{},
		controllerCommon: self.controllerCommon,
		context:          context,
		viewFiles:        self.viewFiles,
	}
}

func (self *CommonCommitController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcViewCommitFiles,
		},
	}

	return bindings
}

func (self *CommonCommitController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.context.GetSelected()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *CommonCommitController) enter(commit *models.Commit) error {
	return self.viewFiles(SwitchToCommitFilesContextOpts{
		RefName:   commit.Sha,
		CanRebase: self.context.CanRebase(),
		Context:   self.context,
	})
}

func (self *CommonCommitController) Context() types.Context {
	return self.context
}
