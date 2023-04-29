package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StagingHelper struct {
	c *HelperCommon
}

func NewStagingHelper(
	c *HelperCommon,
) *StagingHelper {
	return &StagingHelper{
		c: c,
	}
}

// NOTE: used from outside this file
func (self *StagingHelper) RefreshStagingPanel(focusOpts types.OnFocusOpts) error {
	secondaryFocused := self.secondaryStagingFocused()
	mainFocused := self.mainStagingFocused()

	// this method could be called when the staging panel is not being used,
	// in which case we don't want to do anything.
	if !mainFocused && !secondaryFocused {
		return nil
	}

	mainSelectedLineIdx := -1
	secondarySelectedLineIdx := -1
	if focusOpts.ClickedViewLineIdx > 0 {
		if secondaryFocused {
			secondarySelectedLineIdx = focusOpts.ClickedViewLineIdx
		} else {
			mainSelectedLineIdx = focusOpts.ClickedViewLineIdx
		}
	}

	mainContext := self.c.Contexts().Staging
	secondaryContext := self.c.Contexts().StagingSecondary

	var file *models.File
	node := self.c.Contexts().Files.GetSelected()
	if node != nil {
		file = node.File
	}

	if file == nil || (!file.HasUnstagedChanges && !file.HasStagedChanges) {
		return self.handleStagingEscape()
	}

	mainDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, false, false)
	secondaryDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, true, false)

	// grabbing locks here and releasing before we finish the function
	// because pushing say the secondary context could mean entering this function
	// again, and we don't want to have a deadlock
	mainContext.GetMutex().Lock()
	secondaryContext.GetMutex().Lock()

	mainContext.SetState(
		patch_exploring.NewState(mainDiff, mainSelectedLineIdx, mainContext.GetState(), self.c.Log),
	)

	secondaryContext.SetState(
		patch_exploring.NewState(secondaryDiff, secondarySelectedLineIdx, secondaryContext.GetState(), self.c.Log),
	)

	mainState := mainContext.GetState()
	secondaryState := secondaryContext.GetState()

	mainContent := mainContext.GetContentToRender(!secondaryFocused)
	secondaryContent := secondaryContext.GetContentToRender(secondaryFocused)

	mainContext.GetMutex().Unlock()
	secondaryContext.GetMutex().Unlock()

	if mainState == nil && secondaryState == nil {
		return self.handleStagingEscape()
	}

	if mainState == nil && !secondaryFocused {
		return self.c.PushContext(secondaryContext, focusOpts)
	}

	if secondaryState == nil && secondaryFocused {
		return self.c.PushContext(mainContext, focusOpts)
	}

	if secondaryFocused {
		self.c.Contexts().StagingSecondary.FocusSelection()
	} else {
		self.c.Contexts().Staging.FocusSelection()
	}

	return self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Staging,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(mainContent),
			Title: self.c.Tr.UnstagedChanges,
		},
		Secondary: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(secondaryContent),
			Title: self.c.Tr.StagedChanges,
		},
	})
}

func (self *StagingHelper) handleStagingEscape() error {
	return self.c.PushContext(self.c.Contexts().Files)
}

func (self *StagingHelper) secondaryStagingFocused() bool {
	return self.c.CurrentStaticContext().GetKey() == self.c.Contexts().StagingSecondary.GetKey()
}

func (self *StagingHelper) mainStagingFocused() bool {
	return self.c.CurrentStaticContext().GetKey() == self.c.Contexts().Staging.GetKey()
}
