// +build windows

package gocui

import "github.com/jesseduffield/termbox-go"

func (g *Gui) getTermWindowSize() (int, int, error) {
	x, y := termbox.Size()
	return x, y, nil
}
