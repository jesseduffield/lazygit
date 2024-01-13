package ansi

import (
	"io"
	"strconv"
)

// An ItemElement is used to render items inside a list.
type ItemElement struct {
	IsOrdered   bool
	Enumeration uint
}

func (e *ItemElement) Render(w io.Writer, ctx RenderContext) error {
	var el *BaseElement
	if e.IsOrdered {
		el = &BaseElement{
			Style:  ctx.options.Styles.Enumeration,
			Prefix: strconv.FormatInt(int64(e.Enumeration), 10),
		}
	} else {
		el = &BaseElement{
			Style: ctx.options.Styles.Item,
		}
	}

	return el.Render(w, ctx)
}
