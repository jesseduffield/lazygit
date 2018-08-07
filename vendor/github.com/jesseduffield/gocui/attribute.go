// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/nsf/termbox-go"

// Attribute represents a terminal attribute, like color, font style, etc. They
// can be combined using bitwise OR (|). Note that it is not possible to
// combine multiple color attributes.
type Attribute termbox.Attribute

// Color attributes.
const (
	ColorDefault Attribute = Attribute(termbox.ColorDefault)
	ColorBlack             = Attribute(termbox.ColorBlack)
	ColorRed               = Attribute(termbox.ColorRed)
	ColorGreen             = Attribute(termbox.ColorGreen)
	ColorYellow            = Attribute(termbox.ColorYellow)
	ColorBlue              = Attribute(termbox.ColorBlue)
	ColorMagenta           = Attribute(termbox.ColorMagenta)
	ColorCyan              = Attribute(termbox.ColorCyan)
	ColorWhite             = Attribute(termbox.ColorWhite)
)

// Text style attributes.
const (
	AttrBold      Attribute = Attribute(termbox.AttrBold)
	AttrUnderline           = Attribute(termbox.AttrUnderline)
	AttrReverse             = Attribute(termbox.AttrReverse)
)
