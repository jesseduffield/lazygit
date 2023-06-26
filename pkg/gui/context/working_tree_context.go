package context

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorkingTreeContext struct {
	*filetree.FileTreeViewModel
	*ListContextTrait
	*SearchTrait
}

var _ types.IListContext = (*WorkingTreeContext)(nil)

func NewWorkingTreeContext(c *ContextCommon) *WorkingTreeContext {
	viewModel := filetree.NewFileTreeViewModel(
		func() []*models.File { return c.Model().Files },
		c.Log,
		c.UserConfig.Gui.ShowFileTree,
	)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		lines := presentation.RenderFileTree(viewModel, c.Modes().Diffing.Ref, c.Model().Submodules)
		return slices.Map(lines, func(line string) []string {
			return []string{line}
		})
	}

	ctx := &WorkingTreeContext{
		SearchTrait:       NewSearchTrait(c),
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

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(func(selectedLineIdx int) error {
		ctx.GetList().SetSelectedLineIdx(selectedLineIdx)
		return ctx.HandleFocus(types.OnFocusOpts{})
	}))

	return ctx
}

func (self *WorkingTreeContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
