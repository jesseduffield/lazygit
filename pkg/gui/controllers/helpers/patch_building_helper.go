package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PatchBuildingHelper struct {
	c             *HelperCommon
	stagingHelper *StagingHelper
}

func NewPatchBuildingHelper(
	c *HelperCommon,
	stagingHelper *StagingHelper,
) *PatchBuildingHelper {
	return &PatchBuildingHelper{
		c:             c,
		stagingHelper: stagingHelper,
	}
}

func (self *PatchBuildingHelper) ShowHunkStagingHint() {
	if !self.c.AppState.DidShowHunkStagingHint && self.c.UserConfig().Gui.UseHunkModeInStagingView {
		self.c.AppState.DidShowHunkStagingHint = true
		self.c.SaveAppStateAndLogError()

		message := fmt.Sprintf(self.c.Tr.HunkStagingHint, self.c.UserConfig().Keybinding.Main.ToggleSelectHunk)
		self.c.Confirm(types.ConfirmOpts{
			Prompt: message,
		})
	}
}

// takes us from the patch building panel back to the commit files panel, or to
// the focused main view if that's where we entered it from
func (self *PatchBuildingHelper) Escape() {
	EscapeFromPatchExplorer(self.c, self.stagingHelper, self.c.Contexts().CustomPatchBuilder)
}

// EscapeFromPatchExplorer returns from a patch explorer context (staging or
// patch building). If we entered it from a focused main view, we go back to
// where we came from (re-rendering the side panel's content into the main view,
// like the plain escape does), then focus the main view and land on the line the
// explorer currently has selected. Otherwise we just pop to the side panel.
func EscapeFromPatchExplorer(c *HelperCommon, stagingHelper *StagingHelper, context types.IPatchExplorerContext) {
	snapshot := context.GetFocusedMainViewSnapshot()
	if snapshot == nil {
		c.Context().Pop()
		return
	}

	// Clear the snapshot wherever it was set. The staging view records it on both
	// its halves (see FilesController.EnterFile) so escape works after the selection
	// crosses between them; clear both so a stale one can't linger. For patch
	// building it's only on the context we're escaping from.
	context.SetFocusedMainViewSnapshot(nil)
	c.Contexts().Staging.SetFocusedMainViewSnapshot(nil)
	c.Contexts().StagingSecondary.SetFocusedMainViewSnapshot(nil)

	// Restore the side panel's selection before we render it, so it shows the
	// same content the main view had (diving into staging can change it, e.g.
	// from a directory to a file in the files panel).
	if listContext, ok := snapshot.SidePanel.(types.IListContext); ok && snapshot.SidePanelSelectedLineIdx >= 0 {
		listContext.GetList().SetSelectedLineIdx(snapshot.SidePanelSelectedLineIdx)
	}

	// Ask the upcoming re-render of the main view to land on the line the explorer
	// currently has selected. Read that identity now, before the pushes: pushing
	// the side panel re-renders its content into the main view, and the restore
	// rides that re-render — finding the matching row as it loads and scrolling to
	// and selecting it as the content first appears. See RestoreFocusedMainViewOnEscape.
	//
	// Anchor on the *first* line of the explorer's selection: in hunk or range mode
	// the selection spans several lines and its cursor sits at the last one, but
	// returning to the start of the hunk is what reads as "the same place".
	selectedViewLine := context.GetView().SelectedLineIdx()
	if state := context.GetState(); state != nil {
		selectedViewLine, _ = state.SelectedViewRange()
	}
	stagingHelper.RestoreFocusedMainViewOnEscape(
		context.GetView(), snapshot.MainView.GetView(), selectedViewLine)

	// Land on the side panel first (this re-renders the original content into the
	// main view), then focus the main view on top of it.
	c.Context().Push(snapshot.SidePanel, types.OnFocusOpts{})
	c.Context().Push(snapshot.MainView, types.OnFocusOpts{})
}

// kills the custom patch and returns us back to the commit files panel if needed
func (self *PatchBuildingHelper) Reset() error {
	self.c.Git().Patch.PatchBuilder.Reset()

	if self.c.Context().CurrentStatic().GetKind() != types.SIDE_CONTEXT {
		self.Escape()
	}

	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	})

	// refreshing the current context so that the secondary panel is hidden if necessary.
	self.c.PostRefreshUpdate(self.c.Context().Current())
	return nil
}

func (self *PatchBuildingHelper) RefreshPatchBuildingPanel(opts types.OnFocusOpts) {
	selectedLineIdx := -1
	selectedRealLineIdx := -1
	if opts.ClickedWindowName == "main" {
		selectedLineIdx = opts.ClickedViewLineIdx
		selectedRealLineIdx = opts.ClickedViewRealLineIdx
	}

	if !self.c.Git().Patch.PatchBuilder.Active() {
		self.Escape()
		return
	}

	// get diff from commit file that's currently selected
	path := self.c.Contexts().CommitFiles.GetSelectedPath()
	if path == "" {
		return
	}

	from, to := self.c.Contexts().CommitFiles.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	diff, err := self.c.Git().WorkingTree.ShowFileDiff(from, to, reverse, path, true)
	if err != nil {
		return
	}

	secondaryDiff := self.c.Git().Patch.PatchBuilder.RenderPatchForFile(patch.RenderPatchForFileOpts{
		Filename:                               path,
		Plain:                                  false,
		Reverse:                                false,
		TurnAddedFilesIntoDiffAgainstEmptyFile: true,
	})

	context := self.c.Contexts().CustomPatchBuilder

	oldState := context.GetState()

	state := patch_exploring.NewState(diff, selectedLineIdx, selectedRealLineIdx, opts.ClickedViewRealLineIsDeletion, context.GetView(), oldState, self.c.UserConfig().Gui.UseHunkModeInStagingView, opts.SelectLineInDefaultMode)
	context.SetState(state)
	if state == nil {
		self.Escape()
		return
	}

	mainContent := context.GetContentToRender()

	self.c.Contexts().CustomPatchBuilder.FocusSelection()

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().PatchBuilding,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(mainContent),
			Title: self.c.Tr.Patch,
		},
		Secondary: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(secondaryDiff),
			Title: self.c.Tr.CustomPatch,
		},
	})
}
