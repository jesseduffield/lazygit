package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesContext struct {
	*SubmodulesViewModel
	*ListContextTrait
}

var _ types.IListContext = (*SubmodulesContext)(nil)

func NewSubmodulesContext(
	getModel func() []*models.SubmoduleConfig,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *SubmodulesContext {
	viewModel := NewSubmodulesViewModel(getModel)

	return &SubmodulesContext{
		SubmodulesViewModel: viewModel,
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

type SubmodulesViewModel struct {
	*traits.ListCursor
	getModel func() []*models.SubmoduleConfig
}

func NewSubmodulesViewModel(getModel func() []*models.SubmoduleConfig) *SubmodulesViewModel {
	self := &SubmodulesViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *SubmodulesViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *SubmodulesViewModel) GetSelected() *models.SubmoduleConfig {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}
