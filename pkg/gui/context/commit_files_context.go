package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type CommitFilesContext struct {
	*filetree.CommitFileTreeViewModel
	*ListContextTrait
	*DynamicTitleBuilder
}

var (
	_ types.IListContext        = (*CommitFilesContext)(nil)
	_ types.DiffableContext     = (*CommitFilesContext)(nil)
	_ types.IFilterableContext  = (*CommitFilesContext)(nil)
	_ types.DiffMainViewContext = (*CommitFilesContext)(nil)
)

func (self *CommitFilesContext) GetDiffMainViewType() types.DiffMainViewType {
	return types.DiffMainViewTypePatchBuilding
}

func NewCommitFilesContext(c *ContextCommon) *CommitFilesContext {
	viewModel := filetree.NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return c.Model().CommitFiles },
		c.Common,
		c.UserConfig().Gui.ShowFileTree,
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		if viewModel.Len() == 0 {
			return [][]string{{style.FgRed.Sprint("(none)")}}
		}

		showFileIcons := icons.IsIconEnabled() && c.UserConfig().Gui.ShowFileIcons
		lines := presentation.RenderCommitFileTree(viewModel, c.Git().Patch.PatchBuilder, showFileIcons, &c.UserConfig().Gui.CustomIcons)
		return lo.Map(lines, func(line string, _ int) []string {
			return []string{line}
		})
	}

	ctx := &CommitFilesContext{
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
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}

	return ctx
}

func (self *CommitFilesContext) GetDiffTerminals() []string {
	return []string{self.GetRef().RefName()}
}

func (self *CommitFilesContext) RefForAdjustingLineNumberInDiff() string {
	if refs := self.GetRefRange(); refs != nil {
		return refs.To.RefName()
	}
	return self.GetRef().RefName()
}

func (self *CommitFilesContext) GetFromAndToForDiff() (string, string) {
	return FromAndToForDiff(self.GetRef(), self.GetRefRange())
}

// FromAndToForDiff derives the diff endpoints for a ref (or a range of refs): a range
// diffs its parent-of-from against to, a single ref its parent against itself. It's
// shared by the commit files context and by patch building straight from the commits /
// sub-commits / stash main views, which build a patch for the panel's selected ref.
func FromAndToForDiff(ref models.Ref, refRange *types.RefRange) (string, string) {
	if refRange != nil {
		return refRange.From.ParentRefName(), refRange.To.RefName()
	}
	return ref.ParentRefName(), ref.RefName()
}

func (self *CommitFilesContext) ReInit(ref models.Ref, refRange *types.RefRange) {
	self.SetRef(ref)
	self.SetRefRange(refRange)
	if refRange != nil {
		self.SetTitleRef(fmt.Sprintf("%s-%s", refRange.From.ShortRefName(), refRange.To.ShortRefName()))
	} else {
		self.SetTitleRef(ref.Description())
	}
	self.GetView().Title = self.Title()
}
