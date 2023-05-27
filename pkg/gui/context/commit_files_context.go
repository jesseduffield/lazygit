package context

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesContext struct {
	*FilteredList[*models.CommitFile]
	*filetree.CommitFileTreeViewModel
	*ListContextTrait
	*DynamicTitleBuilder
}

var (
	_ types.IListContext    = (*CommitFilesContext)(nil)
	_ types.DiffableContext = (*CommitFilesContext)(nil)
)

func NewCommitFilesContext(c *ContextCommon) *CommitFilesContext {
	filteredList := NewFilteredList(
		func() []*models.CommitFile { return c.Model().CommitFiles },
		func(file *models.CommitFile) []string { return []string{file.GetPath()} },
	)

	viewModel := filetree.NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return filteredList.GetFilteredList() },
		c.Log,
		c.UserConfig.Gui.ShowFileTree,
	)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		if viewModel.Len() == 0 {
			return [][]string{{style.FgRed.Sprint("(none)")}}
		}

		lines := presentation.RenderCommitFileTree(viewModel, c.Modes().Diffing.Ref, c.Git().Patch.PatchBuilder)
		return slices.Map(lines, func(line string) []string {
			return []string{line}
		})
	}

	return &CommitFilesContext{
		FilteredList:            filteredList,
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

// used for type switch
func (self *CommitFilesContext) IsFilterableContext() {}

// TODO: see if we can just call SetTree() within HandleRender(). It doesn't seem
// right that we need to imperatively refresh the view model like this
func (self *CommitFilesContext) SetFilter(filter string) {
	self.FilteredList.SetFilter(filter)
	self.SetTree()
}

func (self *CommitFilesContext) ClearFilter() {
	self.SetFilter("")
}
