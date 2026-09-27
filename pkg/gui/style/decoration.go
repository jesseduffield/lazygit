package style

import "github.com/gookit/color"

type Decoration struct {
	bold          bool
	dim           bool
	underline     bool
	reverse       bool
	strikethrough bool
}

func (d *Decoration) SetBold() {
	d.bold = true
}

func (d *Decoration) SetDim() {
	d.dim = true
}

func (d *Decoration) SetUnderline() {
	d.underline = true
}

func (d *Decoration) SetReverse() {
	d.reverse = true
}

func (d *Decoration) SetStrikethrough() {
	d.strikethrough = true
}

func (d Decoration) ToOpts() color.Opts {
	opts := make([]color.Color, 0, 5)

	if d.bold {
		opts = append(opts, color.OpBold)
	}

	if d.dim {
		opts = append(opts, color.OpFuzzy)
	}

	if d.underline {
		opts = append(opts, color.OpUnderscore)
	}

	if d.reverse {
		opts = append(opts, color.OpReverse)
	}

	if d.strikethrough {
		opts = append(opts, color.OpStrikethrough)
	}

	return opts
}

func (d Decoration) Merge(other Decoration) Decoration {
	if other.bold {
		d.bold = true
	}

	if other.dim {
		d.dim = true
	}

	if other.underline {
		d.underline = true
	}

	if other.reverse {
		d.reverse = true
	}

	if other.strikethrough {
		d.strikethrough = true
	}

	return d
}
