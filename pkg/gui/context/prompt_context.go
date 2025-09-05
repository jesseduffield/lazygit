package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PromptContext struct {
	*SimpleContext
	c *ContextCommon

	State ConfirmationContextState
}

var _ types.Context = (*PromptContext)(nil)

func NewPromptContext(
	c *ContextCommon,
) *PromptContext {
	return &PromptContext{
		c: c,
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:                  c.Views().Prompt,
			WindowName:            "prompt",
			Key:                   PROMPT_CONTEXT_KEY,
			Kind:                  types.TEMPORARY_POPUP,
			Focusable:             true,
			HasUncontrolledBounds: true,
		})),
	}
}
