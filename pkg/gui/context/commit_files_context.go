package context

import (
	"fmt"

	"github.com/jesseduffield/gocui"
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
	*SearchTrait
}

var (
	_ types.IListContext       = (*CommitFilesContext)(nil)
	_ types.DiffableContext    = (*CommitFilesContext)(nil)
	_ types.ISearchableContext = (*CommitFilesContext)(nil)
)

func NewCommitFilesContext(c *ContextCommon) *CommitFilesContext {
	viewModel := filetree.NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return c.Model().CommitFiles },
		c.Log,
		c.UserConfig().Gui.ShowFileTree,
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		if viewModel.Len() == 0 {
			return [][]string{{style.FgRed.Sprint("(none)")}}
		}

		showFileIcons := icons.IsIconEnabled() && c.UserConfig().Gui.ShowFileIcons
		lines := presentation.RenderCommitFileTree(viewModel, c.Git().Patch.PatchBuilder, showFileIcons)
		return lo.Map(lines, func(line string, _ int) []string {
			return []string{line}
		})
	}

	ctx := &CommitFilesContext{
		CommitFileTreeViewModel: viewModel,
		DynamicTitleBuilder:     NewDynamicTitleBuilder(c.Tr.CommitFilesDynamicTitle),
		SearchTrait:             NewSearchTrait(c),
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

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(ctx.OnSearchSelect))

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
	if refs := self.GetRefRange(); refs != nil {
		return refs.From.ParentRefName(), refs.To.RefName()
	}
	ref := self.GetRef()
	return ref.ParentRefName(), ref.RefName()
}

func (self *CommitFilesContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}

func (self *CommitFilesContext) ReInit(ref types.Ref, refRange *types.RefRange) {
	self.SetRef(ref)
	self.SetRefRange(refRange)
	if refRange != nil {
		self.SetTitleRef(fmt.Sprintf("%s-%s", refRange.From.ShortRefName(), refRange.To.ShortRefName()))
	} else {
		self.SetTitleRef(ref.Description())
	}
	self.GetView().Title = self.Title()
}
