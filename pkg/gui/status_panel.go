package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) refreshStatus(g *gocui.Gui) error {
	v, err := g.View("status")
	if err != nil {
		panic(err)
	}
	// for some reason if this isn't wrapped in an update the clear seems to
	// be applied after the other things or something like that; the panel's
	// contents end up cleared
	g.Update(func(*gocui.Gui) error {
		v.Clear()
		pushables, pullables := gui.GitCommand.UpstreamDifferenceCount()
		fmt.Fprint(v, "↑"+pushables+"↓"+pullables)
		branches := gui.State.Branches
		if err := gui.updateHasMergeConflictStatus(); err != nil {
			return err
		}
		if gui.State.HasMergeConflicts {
			fmt.Fprint(v, utils.ColoredString(" (merging)", color.FgYellow))
		}

		if len(branches) == 0 {
			return nil
		}
		branch := branches[0]
		name := utils.ColoredString(branch.Name, branch.GetColor())
		repo := utils.GetCurrentRepoName()
		fmt.Fprint(v, " "+repo+" → "+name)
		return nil
	})

	return nil
}
