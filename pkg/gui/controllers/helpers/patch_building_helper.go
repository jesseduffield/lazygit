package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IPatchBuildingHelper interface {
	ValidateNormalWorkingTreeState() (bool, error)
}

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
	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, self.c.ErrorMsg(self.c.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

// takes us from the patch building panel back to the commit files panel
func (self *PatchBuildingHelper) Escape() error {
	return self.c.PopContext()
}

// kills the custom patch and returns us back to the commit files panel if needed
func (self *PatchBuildingHelper) Reset() error {
	self.c.Git().Patch.PatchBuilder.Reset()

	if self.c.CurrentStaticContext().GetKind() != types.SIDE_CONTEXT {
		if err := self.Escape(); err != nil {
			return err
		}
	}

	if err := self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	}); err != nil {
		return err
	}

	// refreshing the current context so that the secondary panel is hidden if necessary.
	return self.c.PostRefreshUpdate(self.c.CurrentContext())
}

func (self *PatchBuildingHelper) RefreshPatchBuildingPanel(opts types.OnFocusOpts) error {
	selectedLineIdx := -1
	if opts.ClickedWindowName == "main" {
		selectedLineIdx = opts.ClickedViewLineIdx
	}

	if !self.c.Git().Patch.PatchBuilder.Active() {
		return self.Escape()
	}

	// get diff from commit file that's currently selected
	path := self.c.Contexts().CommitFiles.GetSelectedPath()
	if path == "" {
		return nil
	}

	ref := self.c.Contexts().CommitFiles.CommitFileTreeViewModel.GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())
	// Passing false for ignoreWhitespace because the patch building panel
	// doesn't work when whitespace is ignored
	diff, err := self.c.Git().WorkingTree.ShowFileDiff(from, to, reverse, path, true, false)
	if err != nil {
		return err
	}

	secondaryDiff := self.c.Git().Patch.PatchBuilder.RenderPatchForFile(path, false, false)
	if err != nil {
		return err
	}

	context := self.c.Contexts().CustomPatchBuilder

	oldState := context.GetState()

	state := patch_exploring.NewState(diff, selectedLineIdx, oldState, self.c.Log)
	context.SetState(state)
	if state == nil {
		return self.Escape()
	}

	mainContent := context.GetContentToRender(true)

	self.c.Contexts().CustomPatchBuilder.FocusSelection()

	return self.c.RenderToMainViews(types.RefreshMainOpts{
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
