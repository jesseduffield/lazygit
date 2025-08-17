package helpers

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
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

func (self *PatchBuildingHelper) ValidateNormalWorkingTreeState() (bool, error) {
	if self.c.Git().Status.WorkingTreeState().Any() {
		return false, errors.New(self.c.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

func (self *PatchBuildingHelper) ShowHunkStagingHint() {
	if !self.c.AppState.DidShowHunkStagingHint && self.c.UserConfig().Gui.UseHunkModeInStagingView {
		self.c.AppState.DidShowHunkStagingHint = true
		self.c.SaveAppStateAndLogError()

		message := fmt.Sprintf(self.c.Tr.HunkStagingHint,
			keybindings.Label(self.c.UserConfig().Keybinding.Main.ToggleSelectHunk))
		self.c.Confirm(types.ConfirmOpts{
			Prompt: message,
		})
	}
}

// takes us from the patch building panel back to the commit files panel
func (self *PatchBuildingHelper) Escape() {
	self.c.Context().Pop()
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
	if opts.ClickedWindowName == "main" {
		selectedLineIdx = opts.ClickedViewLineIdx
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

	state := patch_exploring.NewState(diff, selectedLineIdx, context.GetView(), oldState, self.c.UserConfig().Gui.UseHunkModeInStagingView)
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
