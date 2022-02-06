package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*TagsViewModel
	*ListContextTrait
}

var _ types.IListContext = (*TagsContext)(nil)

func NewTagsContext(
	getModel func() []*models.Tag,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *TagsContext {
	viewModel := NewTagsViewModel(getModel)

	return &TagsContext{
		TagsViewModel: viewModel,
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

type TagsViewModel struct {
	*traits.ListCursor
	getModel func() []*models.Tag
}

func NewTagsViewModel(getModel func() []*models.Tag) *TagsViewModel {
	self := &TagsViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *TagsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *TagsViewModel) GetSelected() *models.Tag {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *TagsViewModel) GetSelectedRefName() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.RefName()
}
