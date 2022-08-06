package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MergeConflictsController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &MergeConflictsController{}

func NewMergeConflictsController(
	common *controllerCommon,
) *MergeConflictsController {
	return &MergeConflictsController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *MergeConflictsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.EditFile,
			Description: self.c.Tr.LcEditFile,
		},
	}

	return bindings
}

func (self *MergeConflictsController) Context() types.Context {
	return self.context()
}

func (self *MergeConflictsController) context() *context.MergeConflictsContext {
	return self.contexts.MergeConflicts
}

func (self *MergeConflictsController) EditFile() error {
	lineNumber := self.context().State().GetSelectedLine()
	return self.helpers.Files.EditFileAtLine(self.context().State().GetPath(), lineNumber)
}

func (self *MergeConflictsController) withMergeConflictLock(f func() error) error {
	self.context().State().Lock()
	defer self.context().State().Unlock()

	return f()
}
