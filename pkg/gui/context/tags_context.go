package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*FilteredListViewModel[*models.Tag]
	*ListContextTrait
}

var _ types.IListContext = (*TagsContext)(nil)

func NewTagsContext(
	getItems func() []*models.Tag,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *TagsContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.Tag) string {
		return item.Name
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetTagListDisplayStrings(viewModel.getModel(), guiContextState.Modes().Diffing.Ref)
	}

	return &TagsContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "branches",
				WindowName: "branches",
				Key:        TAGS_CONTEXT_KEY,
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

func (self *TagsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *TagsContext) GetSelectedRef() types.Ref {
	tag := self.GetSelected()
	if tag == nil {
		return nil
	}
	return tag
}
