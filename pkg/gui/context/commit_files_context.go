package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesContext struct {
	*filetree.CommitFileTreeViewModel
	*ListContextTrait
	*DynamicTitleBuilder
}

var (
	_ types.IListContext    = (*CommitFilesContext)(nil)
	_ types.DiffableContext = (*CommitFilesContext)(nil)
)

func NewCommitFilesContext(
	getDisplayStrings func(startIdx int, length int) [][]string,

	c *types.HelperCommon,
) *CommitFilesContext {
	viewModel := filetree.NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return c.Model().CommitFiles },
		c.Log,
		c.UserConfig.Gui.ShowFileTree,
	)

	return &CommitFilesContext{
		CommitFileTreeViewModel: viewModel,
		DynamicTitleBuilder:     NewDynamicTitleBuilder(c.Tr.CommitFilesDynamicTitle),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(
				NewBaseContext(NewBaseContextOpts{
					View:       c.Views().CommitFiles,
					WindowName: "commits",
					Key:        COMMIT_FILES_CONTEXT_KEY,
					Kind:       types.SIDE_CONTEXT,
					Focusable:  true,
					Transient:  true,
				}),
			),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *CommitFilesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *CommitFilesContext) GetDiffTerminals() []string {
	return []string{self.GetRef().RefName()}
}

func (self *CommitFilesContext) renderToMain() error {
	node := self.GetSelected()
	if node == nil {
		return nil
	}

	ref := self.GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(
		from, to, reverse, node.GetPath(), false, self.c.State().GetIgnoreWhitespaceInDiffView(),
	)
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	pair := self.c.MainViewPairs().Normal
	if node.File != nil {
		pair = self.c.MainViewPairs().PatchBuilding
	}

	return self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: pair,
		Main: &types.ViewUpdateOpts{
			Title: self.c.Tr.Patch,
			Task:  task,
		},
		Secondary: secondaryPatchPanelUpdateOpts(self.c),
	})
}

func secondaryPatchPanelUpdateOpts(c *types.HelperCommon) *types.ViewUpdateOpts {
	if c.Git().Patch.PatchBuilder.Active() {
		patch := c.Git().Patch.PatchBuilder.RenderAggregatedPatch(false)

		return &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(patch),
			Title: c.Tr.CustomPatch,
		}
	}

	return nil
}
