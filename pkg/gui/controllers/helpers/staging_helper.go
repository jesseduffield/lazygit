package helpers

import (
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type StagingHelper struct {
	c            *HelperCommon
	windowHelper *WindowHelper
}

func NewStagingHelper(
	c *HelperCommon,
	windowHelper *WindowHelper,
) *StagingHelper {
	return &StagingHelper{
		c:            c,
		windowHelper: windowHelper,
	}
}

// NOTE: used from outside this file
func (self *StagingHelper) RefreshStagingPanel(focusOpts types.OnFocusOpts) {
	secondaryFocused := self.secondaryStagingFocused()
	mainFocused := self.mainStagingFocused()

	// this method could be called when the staging panel is not being used,
	// in which case we don't want to do anything.
	if !mainFocused && !secondaryFocused {
		return
	}

	mainSelectedLineIdx := -1
	mainSelectedRealLineIdx := -1
	secondarySelectedLineIdx := -1
	secondarySelectedRealLineIdx := -1
	if focusOpts.ClickedViewLineIdx > 0 {
		if secondaryFocused {
			secondarySelectedLineIdx = focusOpts.ClickedViewLineIdx
			secondarySelectedRealLineIdx = focusOpts.ClickedViewRealLineIdx
		} else {
			mainSelectedLineIdx = focusOpts.ClickedViewLineIdx
			mainSelectedRealLineIdx = focusOpts.ClickedViewRealLineIdx
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
		self.handleStagingEscape()
		return
	}

	mainDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, false)
	secondaryDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, true)

	// grabbing locks here and releasing before we finish the function
	// because pushing say the secondary context could mean entering this function
	// again, and we don't want to have a deadlock
	mainContext.GetMutex().Lock()
	secondaryContext.GetMutex().Lock()

	mainContext.SetState(
		patch_exploring.NewState(mainDiff, mainSelectedLineIdx, mainSelectedRealLineIdx, mainContext.GetView(), mainContext.GetState()),
	)

	secondaryContext.SetState(
		patch_exploring.NewState(secondaryDiff, secondarySelectedLineIdx, secondarySelectedRealLineIdx, secondaryContext.GetView(), secondaryContext.GetState()),
	)

	mainState := mainContext.GetState()
	secondaryState := secondaryContext.GetState()

	mainContent := mainContext.GetContentToRender()
	secondaryContent := secondaryContext.GetContentToRender()

	mainContext.GetMutex().Unlock()
	secondaryContext.GetMutex().Unlock()

	if mainState == nil && secondaryState == nil {
		self.handleStagingEscape()
		return
	}

	if mainState == nil && !secondaryFocused {
		self.c.Context().Push(secondaryContext, focusOpts)
		return
	}

	if secondaryState == nil && secondaryFocused {
		self.c.Context().Push(mainContext, focusOpts)
		return
	}

	if secondaryFocused {
		self.c.Contexts().StagingSecondary.FocusSelection()
	} else {
		self.c.Contexts().Staging.FocusSelection()
	}

	self.c.RenderToMainViews(types.RefreshMainOpts{
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

func (self *StagingHelper) handleStagingEscape() {
	self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
}

func (self *StagingHelper) secondaryStagingFocused() bool {
	return self.c.Context().CurrentStatic().GetKey() == self.c.Contexts().StagingSecondary.GetKey()
}

func (self *StagingHelper) mainStagingFocused() bool {
	return self.c.Context().CurrentStatic().GetKey() == self.c.Contexts().Staging.GetKey()
}

func (self *StagingHelper) GetFileAndLineForClickedDiffLine(windowName string, lineIdx int) (string, int, bool) {
	v, _ := self.c.GocuiGui().View(self.windowHelper.GetViewNameForWindow(windowName))
	hyperlink, ok := v.HyperLinkInLine(lineIdx, "lazygit-edit:")
	if !ok {
		return "", 0, false
	}

	re := regexp.MustCompile(`^lazygit-edit://(.+?):(\d+)$`)
	matches := re.FindStringSubmatch(hyperlink)
	if matches == nil {
		return "", 0, false
	}
	filepath := matches[1]
	lineNumber := utils.MustConvertToInt(matches[2])
	return filepath, lineNumber, true
}
