package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
}

func NewMainContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,
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
	}

	return ctx
}
