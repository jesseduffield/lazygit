package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain commit files.

var _ types.IController = &SwitchToDiffFilesController{}

type CanSwitchToDiffFiles interface {
	types.Context
	CanRebase() bool
	GetSelectedRefName() string
}

type SwitchToDiffFilesController struct {
	baseController
	*controllerCommon
	context   CanSwitchToDiffFiles
	viewFiles func(SwitchToCommitFilesContextOpts) error
}

func NewSwitchToDiffFilesController(
	controllerCommon *controllerCommon,
	viewFiles func(SwitchToCommitFilesContextOpts) error,
	context CanSwitchToDiffFiles,
) *SwitchToDiffFilesController {
	return &SwitchToDiffFilesController{
		baseController:   baseController{},
		controllerCommon: controllerCommon,
		context:          context,
		viewFiles:        viewFiles,
	}
}

func (self *SwitchToDiffFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcViewItemFiles,
		},
	}

	return bindings
}

func (self *SwitchToDiffFilesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *SwitchToDiffFilesController) checkSelected(callback func(string) error) func() error {
	return func() error {
		refName := self.context.GetSelectedRefName()
		if refName == "" {
			return nil
		}

		return callback(refName)
	}
}

func (self *SwitchToDiffFilesController) enter(refName string) error {
	return self.viewFiles(SwitchToCommitFilesContextOpts{
		RefName:   refName,
		CanRebase: self.context.CanRebase(),
		Context:   self.context,
	})
}

func (self *SwitchToDiffFilesController) Context() types.Context {
	return self.context
}
