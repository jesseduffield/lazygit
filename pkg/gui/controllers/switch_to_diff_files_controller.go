package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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
	c                *ControllerCommon
	context          CanSwitchToDiffFiles
	diffFilesContext *context.CommitFilesContext
}

func NewSwitchToDiffFilesController(
	c *ControllerCommon,
	context CanSwitchToDiffFiles,
	diffFilesContext *context.CommitFilesContext,
) *SwitchToDiffFilesController {
	return &SwitchToDiffFilesController{
		baseController:   baseController{},
		c:                c,
		context:          context,
		diffFilesContext: diffFilesContext,
	}
}

func (self *SwitchToDiffFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.ViewItemFiles,
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

func (self *SwitchToDiffFilesController) viewFiles(opts SwitchToCommitFilesContextOpts) error {
	diffFilesContext := self.diffFilesContext

	diffFilesContext.SetSelectedLineIdx(0)
	diffFilesContext.SetRef(opts.Ref)
	diffFilesContext.SetTitleRef(opts.Ref.Description())
	diffFilesContext.SetCanRebase(opts.CanRebase)
	diffFilesContext.SetParentContext(opts.Context)
	diffFilesContext.SetWindowName(opts.Context.GetWindowName())
	diffFilesContext.ClearSearchString()

	if err := self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	}); err != nil {
		return err
	}

	return self.c.PushContext(diffFilesContext)
}
