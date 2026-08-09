package controllers

import (
	"fmt"
	"math"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller lets you change the context size for diffs. The 'context' in 'context size' refers to the conventional meaning of the word 'context' in a diff, as opposed to lazygit's own idea of a 'context'.

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
			Keys:        opts.GetKeys(opts.Config.Universal.IncreaseContextInDiffView),
			Handler:     self.Increase,
			Description: self.c.Tr.IncreaseContextInDiffView,
			Tooltip:     self.c.Tr.IncreaseContextInDiffViewTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.DecreaseContextInDiffView),
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
	if self.c.UserConfig().Git.DiffContextSize < math.MaxUint64 {
		self.c.UserConfig().Git.DiffContextSize++
	}
	return self.applyChange()
}

func (self *ContextLinesController) Decrease() error {
	if self.c.UserConfig().Git.DiffContextSize > 0 {
		self.c.UserConfig().Git.DiffContextSize--
	}
	return self.applyChange()
}

func (self *ContextLinesController) applyChange() error {
	self.c.Toast(fmt.Sprintf(self.c.Tr.DiffContextSizeChanged, self.c.UserConfig().Git.DiffContextSize))

	currentContext := self.c.Context().CurrentSide()
	switch currentContext.GetKey() {
	// we make an exception for our staging and patch building contexts because they actually need to refresh their state afterwards.
	case context.PATCH_BUILDING_MAIN_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.PATCH_BUILDING}})
	case context.STAGING_MAIN_CONTEXT_KEY, context.STAGING_SECONDARY_CONTEXT_KEY:
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STAGING}})
	default:
		// The new context size re-renders the diff with a different git command, which
		// would otherwise reset the main view to the top. Preserve where it was
		// scrolled (and its selection, if the focused main view is showing one)
		// instead. This covers both the focused main view and the side panels, since
		// they all render their diff into the same "main" view.
		self.c.Helpers().Staging.PreserveDiffPositionOnRerender(self.c.Contexts().Normal.GetView())
		currentContext.HandleRenderToMain()
	}
	return nil
}
