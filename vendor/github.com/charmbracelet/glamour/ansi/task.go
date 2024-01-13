package ansi

import (
	"io"
)

// A TaskElement is used to render tasks inside a todo-list.
type TaskElement struct {
	Checked bool
}

func (e *TaskElement) Render(w io.Writer, ctx RenderContext) error {
	var el *BaseElement

	pre := ctx.options.Styles.Task.Unticked
	if e.Checked {
		pre = ctx.options.Styles.Task.Ticked
	}

	el = &BaseElement{
		Prefix: pre,
		Style:  ctx.options.Styles.Task.StylePrimitive,
	}

	return el.Render(w, ctx)
}
