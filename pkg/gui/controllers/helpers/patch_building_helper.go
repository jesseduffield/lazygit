package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PatchBuildingHelper struct {
	c *HelperCommon
}

func NewPatchBuildingHelper(
	c *HelperCommon,
) *PatchBuildingHelper {
	return &PatchBuildingHelper{
		c: c,
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
	EscapeFromPatchExplorer(self.c, self.c.Contexts().CustomPatchBuilder)
}

// EscapeFromPatchExplorer returns from a patch explorer context (staging or
// patch building). If we entered it from a focused main view, we go back to
// where we came from (re-rendering the side panel's content into the main view,
// like the plain escape does), then focus the main view and restore its scroll
// position and selection. Otherwise we just pop to the side panel.
func EscapeFromPatchExplorer(c *HelperCommon, context types.IPatchExplorerContext) {
	snapshot := context.GetFocusedMainViewSnapshot()
	if snapshot == nil {
		c.Context().Pop()
		return
	}

	context.SetFocusedMainViewSnapshot(nil)

	// Restore the side panel's selection before we render it, so it shows the
	// same content the main view had (diving into staging can change it, e.g.
	// from a directory to a file in the files panel).
	if listContext, ok := snapshot.SidePanel.(types.IListContext); ok && snapshot.SidePanelSelectedLineIdx >= 0 {
		listContext.GetList().SetSelectedLineIdx(snapshot.SidePanelSelectedLineIdx)
	}

	view := snapshot.MainView.GetView()

	// Ask the upcoming re-render to restore the scroll position. Pushing the side
	// panel re-renders its content into the main view via a cmd/pty task. Until
	// that content is ready, the main view keeps showing the placeholder that
	// CopyContent left in it (the view we're leaving) at its current scroll; the
	// task then scrolls to the saved position as part of the first paint that
	// shows the real content. Doing it this way (rather than setting the origin up
	// front) avoids both a jump to the top and a misplaced placeholder frame.
	if manager := c.GetViewBufferManagerForView(view); manager != nil {
		manager.ScrollToOriginYForNextTask(snapshot.OriginY)
	}

	// Land on the side panel first (this re-renders the original content into the
	// main view), then focus the main view on top of it.
	c.Context().Push(snapshot.SidePanel, types.OnFocusOpts{})
	c.Context().Push(snapshot.MainView, types.OnFocusOpts{})

	restore := func() {
		view.FocusPoint(0, snapshot.SelectedLineIdx, false)
		view.Highlight = true
		view.HighlightInactive = false
	}

	// The scroll position is handled by the re-render above, but the selection
	// still needs the content loaded down to the selected line, which happens
	// asynchronously. Wait until the diff has been fully read before restoring
	// it. We do this on the next UI tick, by which point the re-render task is
	// live and ReadToEnd can hook into it.
	c.OnUIThread(func() error {
		manager := c.GetViewBufferManagerForView(view)
		if manager == nil {
			restore()
			return nil
		}
		manager.ReadToEnd(func() {
			c.OnUIThread(func() error {
				restore()
				return nil
			})
		})
		return nil
	})
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
	file := self.c.Contexts().CommitFiles.GetSelectedFile()
	if file == nil {
		return
	}

	from, to := self.c.Contexts().CommitFiles.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	diff, err := self.c.Git().WorkingTree.ShowFileDiff(from, to, reverse, file.Path, file.PreviousPath, true)
	if err != nil {
		return
	}

	secondaryDiff := self.c.Git().Patch.PatchBuilder.RenderPatchForFile(patch.RenderPatchForFileOpts{
		Filename:                               file.Path,
		PreviousPath:                           file.PreviousPath,
		Plain:                                  false,
		Reverse:                                false,
		TurnAddedFilesIntoDiffAgainstEmptyFile: true,
	})

	context := self.c.Contexts().CustomPatchBuilder

	oldState := context.GetState()

	state := patch_exploring.NewState(diff, selectedLineIdx, selectedRealLineIdx, context.GetView(), oldState, self.c.UserConfig().Gui.UseHunkModeInStagingView, opts.SelectLineInDefaultMode)
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
