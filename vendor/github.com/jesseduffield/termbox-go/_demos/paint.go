package main

import (
	"github.com/nsf/termbox-go"
)

var curCol = 0
var curRune = 0
var backbuf []termbox.Cell
var bbw, bbh int

var runes = []rune{' ', '░', '▒', '▓', '█'}
var colors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorRed,
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorWhite,
}

type attrFunc func(int) (rune, termbox.Attribute, termbox.Attribute)

func updateAndDrawButtons(current *int, x, y int, mx, my int, n int, attrf attrFunc) {
	lx, ly := x, y
	for i := 0; i < n; i++ {
		if lx <= mx && mx <= lx+3 && ly <= my && my <= ly+1 {
			*current = i
		}
		r, fg, bg := attrf(i)
		termbox.SetCell(lx+0, ly+0, r, fg, bg)
		termbox.SetCell(lx+1, ly+0, r, fg, bg)
		termbox.SetCell(lx+2, ly+0, r, fg, bg)
		termbox.SetCell(lx+3, ly+0, r, fg, bg)
		termbox.SetCell(lx+0, ly+1, r, fg, bg)
		termbox.SetCell(lx+1, ly+1, r, fg, bg)
		termbox.SetCell(lx+2, ly+1, r, fg, bg)
		termbox.SetCell(lx+3, ly+1, r, fg, bg)
		lx += 4
	}
	lx, ly = x, y
	for i := 0; i < n; i++ {
		if *current == i {
			fg := termbox.ColorRed | termbox.AttrBold
			bg := termbox.ColorDefault
			termbox.SetCell(lx+0, ly+2, '^', fg, bg)
			termbox.SetCell(lx+1, ly+2, '^', fg, bg)
			termbox.SetCell(lx+2, ly+2, '^', fg, bg)
			termbox.SetCell(lx+3, ly+2, '^', fg, bg)
		}
		lx += 4
	}
}

func update_and_redraw_all(mx, my int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if mx != -1 && my != -1 {
		backbuf[bbw*my+mx] = termbox.Cell{Ch: runes[curRune], Fg: colors[curCol]}
	}
	copy(termbox.CellBuffer(), backbuf)
	_, h := termbox.Size()
	updateAndDrawButtons(&curRune, 0, 0, mx, my, len(runes), func(i int) (rune, termbox.Attribute, termbox.Attribute) {
		return runes[i], termbox.ColorDefault, termbox.ColorDefault
	})
	updateAndDrawButtons(&curCol, 0, h-3, mx, my, len(colors), func(i int) (rune, termbox.Attribute, termbox.Attribute) {
		return ' ', termbox.ColorDefault, colors[i]
	})
	termbox.Flush()
}

func reallocBackBuffer(w, h int) {
	bbw, bbh = w, h
	backbuf = make([]termbox.Cell, w*h)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	reallocBackBuffer(termbox.Size())
	update_and_redraw_all(-1, -1)

mainloop:
	for {
		mx, my := -1, -1
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break mainloop
			}
		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				mx, my = ev.MouseX, ev.MouseY
			}
		case termbox.EventResize:
			reallocBackBuffer(ev.Width, ev.Height)
		}
		update_and_redraw_all(mx, my)
	}
}
