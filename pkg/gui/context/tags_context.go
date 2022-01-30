package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*TagsViewModel
	*BaseContext
	*ListContextTrait
}

var _ types.IListContext = (*TagsContext)(nil)

func NewTagsContext(
	getModel func() []*models.Tag,
	getView func() *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *TagsContext {
	baseContext := NewBaseContext(NewBaseContextOpts{
		ViewName:   "branches",
		WindowName: "branches",
		Key:        TAGS_CONTEXT_KEY,
		Kind:       types.SIDE_CONTEXT,
	})

	self := &TagsContext{}
	takeFocus := func() error { return c.PushContext(self) }

	list := NewTagsViewModel(getModel)
	viewTrait := NewViewTrait(getView)
	listContextTrait := &ListContextTrait{
		base:      baseContext,
		list:      list,
		viewTrait: viewTrait,

		GetDisplayStrings: getDisplayStrings,
		OnFocus:           onFocus,
		OnRenderToMain:    onRenderToMain,
		OnFocusLost:       onFocusLost,
		takeFocus:         takeFocus,

		// TODO: handle this in a trait
		RenderSelection: false,

		c: c,
	}

	self.BaseContext = baseContext
	self.ListContextTrait = listContextTrait
	self.TagsViewModel = list

	return self
}

func (self *TagsContext) GetSelectedItemId() string {
	item := self.GetSelectedTag()
	if item == nil {
		return ""
	}

	return item.ID()
}

type TagsViewModel struct {
	*traits.ListCursor
	getModel func() []*models.Tag
}

func (self *TagsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *TagsViewModel) GetSelectedTag() *models.Tag {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func NewTagsViewModel(getModel func() []*models.Tag) *TagsViewModel {
	self := &TagsViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}
