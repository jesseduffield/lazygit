package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IPatchBuildingHelper interface {
	ValidateNormalWorkingTreeState() (bool, error)
}

type PatchBuildingHelper struct {
	c        *HelperCommon
	git      *commands.GitCommand
	contexts *context.ContextTree
}

func NewPatchBuildingHelper(
	c *HelperCommon,
	git *commands.GitCommand,
	contexts *context.ContextTree,
) *PatchBuildingHelper {
	return &PatchBuildingHelper{
		c:        c,
		git:      git,
		contexts: contexts,
	}
}

func (self *PatchBuildingHelper) ValidateNormalWorkingTreeState() (bool, error) {
	if self.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
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
	self.git.Patch.PatchBuilder.Reset()

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

	if !self.git.Patch.PatchBuilder.Active() {
		return self.Escape()
	}

	// get diff from commit file that's currently selected
	path := self.contexts.CommitFiles.GetSelectedPath()
	if path == "" {
		return nil
	}

	ref := self.contexts.CommitFiles.CommitFileTreeViewModel.GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())
	diff, err := self.git.WorkingTree.ShowFileDiff(from, to, reverse, path, true, self.c.State().GetIgnoreWhitespaceInDiffView())
	if err != nil {
		return err
	}

	secondaryDiff := self.git.Patch.PatchBuilder.RenderPatchForFile(path, false, false)
	if err != nil {
		return err
	}

	context := self.contexts.CustomPatchBuilder

	oldState := context.GetState()

	state := patch_exploring.NewState(diff, selectedLineIdx, oldState, self.c.Log)
	context.SetState(state)
	if state == nil {
		return self.Escape()
	}

	mainContent := context.GetContentToRender(true)

	self.contexts.CustomPatchBuilder.FocusSelection()

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
