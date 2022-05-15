package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesContext struct {
	*FilteredListViewModel[*models.SubmoduleConfig]
	*ListContextTrait
}

var _ types.IListContext = (*SubmodulesContext)(nil)

func NewSubmodulesContext(
	getItems func() []*models.SubmoduleConfig,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *SubmodulesContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.SubmoduleConfig) string {
		return item.Name
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetSubmoduleListDisplayStrings(viewModel.getModel())
	}

	return &SubmodulesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "files",
				WindowName: "files",
				Key:        SUBMODULES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			}), ContextCallbackOpts{
				OnFocus:        onFocus,
				OnFocusLost:    onFocusLost,
				OnRenderToMain: onRenderToMain,
			}),
			list:              viewModel,
			viewTrait:         NewViewTrait(view),
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *SubmodulesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
