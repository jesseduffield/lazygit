package ansi

import (
	"bytes"
	"io"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
)

// A HeadingElement is used to render headings.
type HeadingElement struct {
	Level int
	First bool
}

func (e *HeadingElement) Render(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := ctx.options.Styles.Heading

	switch e.Level {
	case 1:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H1)
	case 2:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H2)
	case 3:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H3)
	case 4:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H4)
	case 5:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H5)
	case 6:
		rules = cascadeStyles(true, rules, ctx.options.Styles.H6)
	}

	if !e.First {
		renderText(w, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, "\n")
	}

	be := BlockElement{
		Block: &bytes.Buffer{},
		Style: cascadeStyle(bs.Current().Style, rules, false),
	}
	bs.Push(be)

	renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, rules.BlockPrefix)
	renderText(bs.Current().Block, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, rules.Prefix)
	return nil
}

func (e *HeadingElement) Finish(w io.Writer, ctx RenderContext) error {
	bs := ctx.blockStack
	rules := bs.Current().Style

	var indentation uint
	var margin uint
	if rules.Indent != nil {
		indentation = *rules.Indent
	}
	if rules.Margin != nil {
		margin = *rules.Margin
	}

	iw := indent.NewWriterPipe(w, indentation+margin, func(wr io.Writer) {
		renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, " ")
	})

	flow := wordwrap.NewWriter(int(bs.Width(ctx) - indentation - margin*2))
	_, err := flow.Write(bs.Current().Block.Bytes())
	if err != nil {
		return err
	}
	flow.Close()

	_, err = iw.Write(flow.Bytes())
	if err != nil {
		return err
	}

	renderText(w, ctx.options.ColorProfile, bs.Current().Style.StylePrimitive, rules.Suffix)
	renderText(w, ctx.options.ColorProfile, bs.Parent().Style.StylePrimitive, rules.BlockSuffix)

	bs.Current().Block.Reset()
	bs.Pop()
	return nil
}
