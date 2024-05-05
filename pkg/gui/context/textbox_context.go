package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TextboxContext struct {
	*SimpleContext
	c *ContextCommon

	State TextboxContextState
}

type TextboxContextState struct {
	OnConfirm func() error
	OnClose   func() error
}

var _ types.Context = (*TextboxContext)(nil)

func NewTextboxContext(
	c *ContextCommon,
) *TextboxContext {
	return &TextboxContext{
		c: c,
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:                  c.Views().Textbox,
			WindowName:            "textbox",
			Key:                   TEXTBOX_CONTEXT_KEY,
			Kind:                  types.TEMPORARY_POPUP,
			Focusable:             true,
			HasUncontrolledBounds: true,
		})),
	}
}
