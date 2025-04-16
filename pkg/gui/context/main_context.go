package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
	*SearchTrait
}

var _ types.ISearchableContext = (*MainContext)(nil)

func NewMainContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,
	c *ContextCommon,
) *MainContext {
	ctx := &MainContext{
		SimpleContext: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:             types.MAIN_CONTEXT,
				View:             view,
				WindowName:       windowName,
				Key:              key,
				Focusable:        true,
				HighlightOnFocus: false,
			})),
		SearchTrait: NewSearchTrait(c),
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(func(int) error { return nil }))

	return ctx
}

func (self *MainContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}
