package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// This controller lets you change the similarity threshold for detecting renames.

var CONTEXT_KEYS_SHOWING_RENAMES = []types.ContextKey{
	context.FILES_CONTEXT_KEY,
	context.SUB_COMMITS_CONTEXT_KEY,
	context.LOCAL_COMMITS_CONTEXT_KEY,
	context.STASH_CONTEXT_KEY,
}

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
			Key:         opts.GetKey(opts.Config.Universal.IncreaseRenameSimilarityThreshold),
			Handler:     self.Increase,
			Description: self.c.Tr.IncreaseRenameSimilarityThreshold,
			Tooltip:     self.c.Tr.IncreaseRenameSimilarityThresholdTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DecreaseRenameSimilarityThreshold),
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
	old_size := self.c.AppState.RenameSimilarityThreshold

	if self.isShowingRenames() && old_size < 100 {
		self.c.AppState.RenameSimilarityThreshold = min(100, old_size+5)
		return self.applyChange()
	}

	return nil
}

func (self *RenameSimilarityThresholdController) Decrease() error {
	old_size := self.c.AppState.RenameSimilarityThreshold

	if self.isShowingRenames() && old_size > 5 {
		self.c.AppState.RenameSimilarityThreshold = max(5, old_size-5)
		return self.applyChange()
	}

	return nil
}

func (self *RenameSimilarityThresholdController) applyChange() error {
	self.c.Toast(fmt.Sprintf(self.c.Tr.RenameSimilarityThresholdChanged, self.c.AppState.RenameSimilarityThreshold))
	self.c.SaveAppStateAndLogError()

	currentContext := self.c.Context().CurrentStatic()
	switch currentContext.GetKey() {
	// we make an exception for our files context, because it actually need to refresh its state afterwards.
	case context.FILES_CONTEXT_KEY:
		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	default:
		currentContext.HandleRenderToMain()
		return nil
	}
}

func (self *RenameSimilarityThresholdController) isShowingRenames() bool {
	return lo.Contains(
		CONTEXT_KEYS_SHOWING_RENAMES,
		self.c.Context().CurrentStatic().GetKey(),
	)
}
