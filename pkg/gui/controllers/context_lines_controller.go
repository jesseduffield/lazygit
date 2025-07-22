package controllers

import (
	"errors"
	"fmt"
	"math"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// This controller lets you change the context size for diffs. The 'context' in 'context size' refers to the conventional meaning of the word 'context' in a diff, as opposed to lazygit's own idea of a 'context'.

var CONTEXT_KEYS_SHOWING_DIFFS = []types.ContextKey{
	context.FILES_CONTEXT_KEY,
	context.COMMIT_FILES_CONTEXT_KEY,
	context.STASH_CONTEXT_KEY,
	context.LOCAL_COMMITS_CONTEXT_KEY,
	context.SUB_COMMITS_CONTEXT_KEY,
	context.STAGING_MAIN_CONTEXT_KEY,
	context.STAGING_SECONDARY_CONTEXT_KEY,
	context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
	context.PATCH_BUILDING_SECONDARY_CONTEXT_KEY,
	context.NORMAL_MAIN_CONTEXT_KEY,
	context.NORMAL_SECONDARY_CONTEXT_KEY,
}

type ContextLinesController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &ContextLinesController{}

func NewContextLinesController(
	c *ControllerCommon,
) *ContextLinesController {
	return &ContextLinesController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *ContextLinesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.IncreaseContextInDiffView),
			Handler:     self.Increase,
			Description: self.c.Tr.IncreaseContextInDiffView,
			Tooltip:     self.c.Tr.IncreaseContextInDiffViewTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DecreaseContextInDiffView),
			Handler:     self.Decrease,
			Description: self.c.Tr.DecreaseContextInDiffView,
			Tooltip:     self.c.Tr.DecreaseContextInDiffViewTooltip,
		},
	}

	return bindings
}

func (self *ContextLinesController) Context() types.Context {
	return nil
}

func (self *ContextLinesController) Increase() error {
	if self.isShowingDiff() {
		if err := self.checkCanChangeContext(); err != nil {
			return err
		}

		self.incrementContextSize()
		return self.applyChange()
	}

	return nil
}

func (self *ContextLinesController) Decrease() error {
	if self.isShowingDiff() {
		if err := self.checkCanChangeContext(); err != nil {
			return err
		}

		self.decrementContextSize()
		return self.applyChange()
	}

	return nil
}

// incrementContextSize increases the context size for the current context
func (self *ContextLinesController) incrementContextSize() {
	adaptiveConfig := &self.c.UserConfig().Git.AdaptiveContext
	if !adaptiveConfig.Enabled {
		if self.c.UserConfig().Git.DiffContextSize < math.MaxUint64 {
			self.c.UserConfig().Git.DiffContextSize++
		}
		return
	}

	currentContext := self.currentSidePanel().GetKey()
	switch currentContext {
	case context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY:
		if adaptiveConfig.Files < math.MaxUint64 {
			adaptiveConfig.Files++
		}
	case context.LOCAL_COMMITS_CONTEXT_KEY, context.SUB_COMMITS_CONTEXT_KEY:
		if adaptiveConfig.Commits < math.MaxUint64 {
			adaptiveConfig.Commits++
		}
	case context.STASH_CONTEXT_KEY:
		if adaptiveConfig.Stash < math.MaxUint64 {
			adaptiveConfig.Stash++
		}
	case context.STAGING_MAIN_CONTEXT_KEY, context.STAGING_SECONDARY_CONTEXT_KEY:
		if adaptiveConfig.Staging < math.MaxUint64 {
			adaptiveConfig.Staging++
		}
	case context.PATCH_BUILDING_MAIN_CONTEXT_KEY, context.PATCH_BUILDING_SECONDARY_CONTEXT_KEY:
		if adaptiveConfig.PatchBuilding < math.MaxUint64 {
			adaptiveConfig.PatchBuilding++
		}
	default:
		if self.c.UserConfig().Git.DiffContextSize < math.MaxUint64 {
			self.c.UserConfig().Git.DiffContextSize++
		}
	}
}

// decrementContextSize decreases the context size for the current context
func (self *ContextLinesController) decrementContextSize() {
	adaptiveConfig := &self.c.UserConfig().Git.AdaptiveContext
	if !adaptiveConfig.Enabled {
		if self.c.UserConfig().Git.DiffContextSize > 0 {
			self.c.UserConfig().Git.DiffContextSize--
		}
		return
	}

	currentContext := self.currentSidePanel().GetKey()
	switch currentContext {
	case context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY:
		if adaptiveConfig.Files > 0 {
			adaptiveConfig.Files--
		}
	case context.LOCAL_COMMITS_CONTEXT_KEY, context.SUB_COMMITS_CONTEXT_KEY:
		if adaptiveConfig.Commits > 0 {
			adaptiveConfig.Commits--
		}
	case context.STASH_CONTEXT_KEY:
		if adaptiveConfig.Stash > 0 {
			adaptiveConfig.Stash--
		}
	case context.STAGING_MAIN_CONTEXT_KEY, context.STAGING_SECONDARY_CONTEXT_KEY:
		if adaptiveConfig.Staging > 0 {
			adaptiveConfig.Staging--
		}
	case context.PATCH_BUILDING_MAIN_CONTEXT_KEY, context.PATCH_BUILDING_SECONDARY_CONTEXT_KEY:
		if adaptiveConfig.PatchBuilding > 0 {
			adaptiveConfig.PatchBuilding--
		}
	default:
		if self.c.UserConfig().Git.DiffContextSize > 0 {
			self.c.UserConfig().Git.DiffContextSize--
		}
	}
}

func (self *ContextLinesController) applyChange() error {
	self.c.Toast(fmt.Sprintf(self.c.Tr.DiffContextSizeChanged, self.getContextSizeForCurrentContext()))

	currentContext := self.currentSidePanel()
	switch currentContext.GetKey() {
	// we make an exception for our staging and patch building contexts because they actually need to refresh their state afterwards.
	case context.PATCH_BUILDING_MAIN_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.PATCH_BUILDING}})
	case context.STAGING_MAIN_CONTEXT_KEY, context.STAGING_SECONDARY_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STAGING}})
	default:
		currentContext.HandleRenderToMain()
	}
	return nil
}

func (self *ContextLinesController) checkCanChangeContext() error {
	if self.c.Git().Patch.PatchBuilder.Active() {
		return errors.New(self.c.Tr.CantChangeContextSizeError)
	}

	return nil
}

func (self *ContextLinesController) isShowingDiff() bool {
	return lo.Contains(
		CONTEXT_KEYS_SHOWING_DIFFS,
		self.currentSidePanel().GetKey(),
	)
}

func (self *ContextLinesController) currentSidePanel() types.Context {
	return self.c.currentSidePanel()
}

// getContextSizeForCurrentContext returns the appropriate context size based on the current context
func (self *ContextLinesController) getContextSizeForCurrentContext() uint64 {
	return self.c.getContextSizeForCurrentContext()
}
