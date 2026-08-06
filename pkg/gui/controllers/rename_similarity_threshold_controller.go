package controllers

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller lets you change the similarity threshold for detecting renames.

type RenameSimilarityThresholdController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &RenameSimilarityThresholdController{}

func NewRenameSimilarityThresholdController(
	common *ControllerCommon,
) *RenameSimilarityThresholdController {
	return &RenameSimilarityThresholdController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *RenameSimilarityThresholdController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Universal.IncreaseRenameSimilarityThreshold),
			Handler:     self.Increase,
			Description: self.c.Tr.IncreaseRenameSimilarityThreshold,
			Tooltip:     self.c.Tr.IncreaseRenameSimilarityThresholdTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.DecreaseRenameSimilarityThreshold),
			Handler:     self.Decrease,
			Description: self.c.Tr.DecreaseRenameSimilarityThreshold,
			Tooltip:     self.c.Tr.DecreaseRenameSimilarityThresholdTooltip,
		},
	}

	return bindings
}

func (self *RenameSimilarityThresholdController) Context() types.Context {
	return nil
}

func (self *RenameSimilarityThresholdController) Increase() error {
	if err := self.checkCanChangeThreshold(); err != nil {
		return err
	}

	old_size := self.c.UserConfig().Git.RenameSimilarityThreshold

	if old_size < 100 {
		self.c.UserConfig().Git.RenameSimilarityThreshold = min(100, old_size+5)
	}

	return self.applyChange()
}

func (self *RenameSimilarityThresholdController) Decrease() error {
	if err := self.checkCanChangeThreshold(); err != nil {
		return err
	}

	old_size := self.c.UserConfig().Git.RenameSimilarityThreshold

	if old_size > 5 {
		self.c.UserConfig().Git.RenameSimilarityThreshold = max(5, old_size-5)
	}

	return self.applyChange()
}

func (self *RenameSimilarityThresholdController) applyChange() error {
	self.c.Toast(fmt.Sprintf(self.c.Tr.RenameSimilarityThresholdChanged, self.c.UserConfig().Git.RenameSimilarityThreshold))

	currentContext := self.c.Context().CurrentSide()
	switch currentContext.GetKey() {
	// we make an exception for the files and commit-files contexts, because
	// they actually need to refresh their state afterwards: a changed threshold
	// can turn a rename into a separate delete and add, or vice versa.
	case context.FILES_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	case context.COMMIT_FILES_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMIT_FILES}})
	default:
		currentContext.HandleRenderToMain()
	}
	return nil
}

func (self *RenameSimilarityThresholdController) checkCanChangeThreshold() error {
	if self.c.Git().Patch.PatchBuilder.Active() {
		return errors.New(self.c.Tr.CantChangeRenameThresholdError)
	}

	return nil
}
