package gui

import (
	"github.com/jesseduffield/gocui"
)

// refreshTags will refresh the tags
func (gui *Gui) refreshTags(g *gocui.Gui) error {

	g.Update(func(g *gocui.Gui) error {
		return nil
	})

	return nil
}
