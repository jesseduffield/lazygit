package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain commit files.

type CommitishControllerFactory struct {
	controllerCommon *controllerCommon
	viewFiles        func(SwitchToCommitFilesContextOpts) error
}

var _ types.IController = &CommitishController{}

type Commitish interface {
	types.Context
	CanRebase() bool
	GetSelectedRefName() string
}

type CommitishController struct {
	baseController
	*controllerCommon
	context Commitish

	viewFiles func(SwitchToCommitFilesContextOpts) error
}

func NewCommitishControllerFactory(
	common *controllerCommon,
	viewFiles func(SwitchToCommitFilesContextOpts) error,
) *CommitishControllerFactory {
	return &CommitishControllerFactory{
		controllerCommon: common,
		viewFiles:        viewFiles,
	}
}

func (self *CommitishControllerFactory) Create(context Commitish) *CommitishController {
	return &CommitishController{
		baseController:   baseController{},
		controllerCommon: self.controllerCommon,
		context:          context,
		viewFiles:        self.viewFiles,
	}
}

func (self *CommitishController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcViewItemFiles,
		},
	}

	return bindings
}

func (self *CommitishController) checkSelected(callback func(string) error) func() error {
	return func() error {
		refName := self.context.GetSelectedRefName()
		if refName == "" {
			return nil
		}

		return callback(refName)
	}
}

func (self *CommitishController) enter(refName string) error {
	return self.viewFiles(SwitchToCommitFilesContextOpts{
		RefName:   refName,
		CanRebase: self.context.CanRebase(),
		Context:   self.context,
	})
}

func (self *CommitishController) Context() types.Context {
	return self.context
}
