package ansi

import (
	"bytes"
	"io"
)

// StyleWriter is a Writer that applies styling on whatever you write to it.
type StyleWriter struct {
	ctx   RenderContext
	w     io.Writer
	buf   bytes.Buffer
	rules StylePrimitive
}

// NewStyleWriter returns a new StyleWriter.
func NewStyleWriter(ctx RenderContext, w io.Writer, rules StylePrimitive) *StyleWriter {
	return &StyleWriter{
		ctx:   ctx,
		w:     w,
		rules: rules,
	}
}

func (w *StyleWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// Close must be called when you're finished writing to a StyleWriter.
func (w *StyleWriter) Close() error {
	renderText(w.w, w.ctx.options.ColorProfile, w.rules, w.buf.String())
	return nil
}
