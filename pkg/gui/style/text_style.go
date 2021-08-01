package style

import (
	"github.com/gookit/color"
)

// A TextStyle contains a foreground color, background color, and
// decorations (bold/underline/reverse).
//
// Colors may each be either 16-bit or 24-bit RGB colors. When
// we need to produce a string with a TextStyle, if either foreground or
// background color is RGB, we'll promote the other color component to RGB as well.
// We could simplify this code by forcing everything to be RGB, but we're not
// sure how compatible or efficient that would be with various terminals.
// Lazygit will typically stick to 16-bit colors, but users may configure RGB colors.
//
// TextStyles are value objects, not entities, so for example if you want to
// add the bold decoration to a TextStyle, we'll create a new TextStyle with
// that decoration applied.
//
// Decorations are additive, so when we merge two TextStyles, if either is bold
// then the resulting style will also be bold.
//
// So that we aren't rederiving the underlying style each time we want to print
// a string, we derive it when a new TextStyle is created and store it in the
// `style` field.

type TextStyle struct {
	fg         *Color
	bg         *Color
	decoration Decoration

	style Sprinter
}

type Sprinter interface {
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
}

func New() TextStyle {
	s := TextStyle{}
	s.style = s.deriveStyle()
	return s
}

func (b TextStyle) Sprint(a ...interface{}) string {
	return b.style.Sprint(a...)
}

func (b TextStyle) Sprintf(format string, a ...interface{}) string {
	return b.style.Sprintf(format, a...)
}

// note that our receiver here is not a pointer which means we're receiving a
// copy of the original TextStyle. This allows us to mutate and return that
// TextStyle receiver without actually modifying the original.
func (b TextStyle) SetBold() TextStyle {
	b.decoration.SetBold()
	b.style = b.deriveStyle()
	return b
}

func (b TextStyle) SetUnderline() TextStyle {
	b.decoration.SetUnderline()
	b.style = b.deriveStyle()
	return b
}

func (b TextStyle) SetReverse() TextStyle {
	b.decoration.SetReverse()
	b.style = b.deriveStyle()
	return b
}

func (b TextStyle) SetBg(color Color) TextStyle {
	b.bg = &color
	b.style = b.deriveStyle()
	return b
}

func (b TextStyle) SetFg(color Color) TextStyle {
	b.fg = &color
	b.style = b.deriveStyle()
	return b
}

func (b TextStyle) MergeStyle(other TextStyle) TextStyle {
	b.decoration = b.decoration.Merge(other.decoration)

	if other.fg != nil {
		b.fg = other.fg
	}

	if other.bg != nil {
		b.bg = other.bg
	}

	b.style = b.deriveStyle()

	return b
}

func (b TextStyle) deriveStyle() Sprinter {
	if b.fg == nil && b.bg == nil {
		return color.Style(b.decoration.ToOpts())
	}

	isRgb := (b.fg != nil && b.fg.IsRGB()) || (b.bg != nil && b.bg.IsRGB())
	if isRgb {
		return b.deriveRGBStyle()
	}

	return b.deriveBasicStyle()
}

func (b TextStyle) deriveBasicStyle() color.Style {
	style := make([]color.Color, 0, 5)

	if b.fg != nil {
		style = append(style, *b.fg.basic)
	}

	if b.bg != nil {
		style = append(style, *b.bg.basic)
	}

	style = append(style, b.decoration.ToOpts()...)

	return color.Style(style)
}

func (b TextStyle) deriveRGBStyle() *color.RGBStyle {
	style := &color.RGBStyle{}

	if b.fg != nil {
		style.SetFg(*b.fg.ToRGB(false).rgb)
	}

	if b.bg != nil {
		// We need to convert the bg firstly to a foreground color,
		// For more info see
		style.SetBg(*b.bg.ToRGB(true).rgb)
	}

	style.SetOpts(b.decoration.ToOpts())

	return style
}
