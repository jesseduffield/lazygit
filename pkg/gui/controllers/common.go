package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ControllerCommon struct {
	*helpers.HelperCommon
	IGetHelpers
}

type IGetHelpers interface {
	Helpers() *helpers.Helpers
}

func NewControllerCommon(
	c *helpers.HelperCommon,
	IGetHelpers IGetHelpers,
) *ControllerCommon {
	return &ControllerCommon{
		HelperCommon: c,
		IGetHelpers:  IGetHelpers,
	}
}

// getContextSizeForCurrentContext returns the appropriate context size based on the current context
func (self *ControllerCommon) getContextSizeForCurrentContext() uint64 {
	adaptiveConfig := self.UserConfig().Git.AdaptiveContext
	if !adaptiveConfig.Enabled {
		return self.UserConfig().Git.DiffContextSize
	}

	currentContext := self.currentSidePanel().GetKey()
	switch currentContext {
	case context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY:
		return adaptiveConfig.Files
	case context.LOCAL_COMMITS_CONTEXT_KEY, context.SUB_COMMITS_CONTEXT_KEY:
		return adaptiveConfig.Commits
	case context.STASH_CONTEXT_KEY:
		return adaptiveConfig.Stash
	case context.STAGING_MAIN_CONTEXT_KEY, context.STAGING_SECONDARY_CONTEXT_KEY:
		return adaptiveConfig.Staging
	case context.PATCH_BUILDING_MAIN_CONTEXT_KEY, context.PATCH_BUILDING_SECONDARY_CONTEXT_KEY:
		return adaptiveConfig.PatchBuilding
	default:
		return self.UserConfig().Git.DiffContextSize
	}
}

func (self *ControllerCommon) currentSidePanel() types.Context {
	currentContext := self.Context().CurrentStatic()
	if currentContext.GetKey() == context.NORMAL_MAIN_CONTEXT_KEY ||
		currentContext.GetKey() == context.NORMAL_SECONDARY_CONTEXT_KEY {
		if sidePanelContext := self.Context().NextInStack(currentContext); sidePanelContext != nil {
			return sidePanelContext
		}
	}

	return currentContext
}
