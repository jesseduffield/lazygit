package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
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

	getDisplayStrings := func(_ int, _ int) [][]string {
		showFileIcons := icons.IsIconEnabled() && c.UserConfig.Gui.ShowFileIcons
		lines := presentation.RenderFileTree(viewModel, c.Model().Submodules, showFileIcons)
		return lo.Map(lines, func(line string, _ int) []string {
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
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(func(selectedLineIdx int) error {
		ctx.GetList().SetSelection(selectedLineIdx)
		return ctx.HandleFocus(types.OnFocusOpts{})
	}))

	return ctx
}
