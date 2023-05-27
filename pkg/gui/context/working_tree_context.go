package context

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorkingTreeContext struct {
	*FilteredList[*models.File]
	*filetree.FileTreeViewModel
	*ListContextTrait
}

var (
	_ types.IListContext       = (*WorkingTreeContext)(nil)
	_ types.IFilterableContext = (*WorkingTreeContext)(nil)
)

func NewWorkingTreeContext(c *ContextCommon) *WorkingTreeContext {
	filteredList := NewFilteredList(
		func() []*models.File { return c.Model().Files },
		func(file *models.File) []string { return []string{file.GetPath()} },
	)

	viewModel := filetree.NewFileTreeViewModel(
		func() []*models.File { return filteredList.GetFilteredList() },
		c.Log,
		c.UserConfig.Gui.ShowFileTree,
	)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		lines := presentation.RenderFileTree(viewModel, c.Modes().Diffing.Ref, c.Model().Submodules)
		return slices.Map(lines, func(line string) []string {
			return []string{line}
		})
	}

	return &WorkingTreeContext{
		FilteredList:      filteredList,
		FileTreeViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Files,
				WindowName: "files",
				Key:        FILES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *WorkingTreeContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

// used for type switch
func (self *WorkingTreeContext) IsFilterableContext() {}

// TODO: see if we can just call SetTree() within HandleRender(). It doesn't seem
// right that we need to imperatively refresh the view model like this
func (self *WorkingTreeContext) SetFilter(filter string) {
	self.FilteredList.SetFilter(filter)
	self.SetTree()
}

func (self *WorkingTreeContext) ClearFilter() {
	self.SetFilter("")
}
