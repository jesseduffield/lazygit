package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*TagsContextAux
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

	aux := NewTagsContextAux(getModel)
	viewTrait := NewViewTrait(getView)
	listContextTrait := &ListContextTrait{
		base:      baseContext,
		thing:     aux,
		listTrait: aux.list,
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
	self.TagsContextAux = aux

	return self
}

type TagsContextAux struct {
	list     *ListTrait
	getModel func() []*models.Tag
}

func (self *TagsContextAux) GetItemsLength() int {
	return len(self.getModel())
}

func (self *TagsContextAux) GetSelectedTag() *models.Tag {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.list.GetSelectedLineIdx()]
}

func (self *TagsContextAux) GetSelectedItem() (types.ListItem, bool) {
	tag := self.GetSelectedTag()
	return tag, tag != nil
}

func NewTagsContextAux(getModel func() []*models.Tag) *TagsContextAux {
	self := &TagsContextAux{
		getModel: getModel,
	}

	self.list = &ListTrait{
		selectedIdx: 0,
		HasLength:   self,
	}

	return self
}

func clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
