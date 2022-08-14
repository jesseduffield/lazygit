package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain commit files.

var _ types.IController = &SwitchToDiffFilesController{}

type CanSwitchToDiffFiles interface {
	types.Context
	CanRebase() bool
	GetSelectedRef() types.Ref
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

func (self *SwitchToDiffFilesController) checkSelected(callback func(types.Ref) error) func() error {
	return func() error {
		ref := self.context.GetSelectedRef()
		if ref == nil {
			return nil
		}

		return callback(ref)
	}
}

func (self *SwitchToDiffFilesController) enter(ref types.Ref) error {
	return self.viewFiles(SwitchToCommitFilesContextOpts{
		Ref:       ref,
		CanRebase: self.context.CanRebase(),
		Context:   self.context,
	})
}

func (self *SwitchToDiffFilesController) Context() types.Context {
	return self.context
}
