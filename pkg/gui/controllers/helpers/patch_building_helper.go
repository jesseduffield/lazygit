package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"

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

// takes us from the patch building panel back to the commit files panel
func (self *PatchBuildingHelper) Escape() {
	self.c.Context().Pop()
}

// kills the custom patch and returns us back to the commit files panel if needed
func (self *PatchBuildingHelper) Reset() error {
	self.c.Git().Patch.PatchBuilder.Reset()

	// The patch-building view is the one thing with nothing left to show once the patch is
	// gone, so it is the one thing left behind. Everywhere else the user is looking at
	// something of their own — a commit's diff, the working tree's — which is still there.
	if self.c.Context().CurrentStatic().GetKey() == self.c.Contexts().CustomPatchBuilder.GetKey() {
		self.Escape()
	}

	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	})

	// Render again so that the pane that was previewing the patch goes with it. The
	// panel asked to do that is the side panel rather than whichever context has the
	// focus. Both main panes are rendered by the panel beneath them, so a reset from
	// within the focused main view has to go through that panel too.
	self.c.PostRefreshUpdate(self.c.Context().CurrentSide())
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
	file := self.c.Contexts().CommitFiles.GetSelectedFile()
	if file == nil {
		return
	}

	from, to := self.c.Contexts().CommitFiles.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	diff, err := self.c.Git().WorkingTree.ShowFileDiff(from, to, reverse, file.Path, file.PreviousPath, git_commands.DiffModePlain)
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
