package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) refreshStashEntries(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("stash")
		if err != nil {
			panic(err)
		}
		gui.State.StashEntries = gui.GitCommand.GetStashEntries()
		v.Clear()
		for _, stashEntry := range gui.State.StashEntries {
			fmt.Fprintln(v, stashEntry.DisplayString)
		}
		return gui.resetOrigin(v)
	})
	return nil
}

func (gui *Gui) getSelectedStashEntry(v *gocui.View) *commands.StashEntry {
	if len(gui.State.StashEntries) == 0 {
		return nil
	}
	lineNumber := gui.getItemPosition(v)
	return &gui.State.StashEntries[lineNumber]
}

func (gui *Gui) renderStashOptions(g *gocui.Gui) error {
	return gui.renderOptionsMap(g, map[string]string{
		"space":   "apply",
		"g":       "pop",
		"d":       "drop",
		"← → ↑ ↓": "navigate",
	})
}

func (gui *Gui) handleStashEntrySelect(g *gocui.Gui, v *gocui.View) error {
	if err := gui.renderStashOptions(g); err != nil {
		return err
	}
	go func() {
		stashEntry := gui.getSelectedStashEntry(v)
		if stashEntry == nil {
			gui.renderString(g, "main", "No stash entries")
			return
		}
		diff, _ := gui.GitCommand.GetStashEntryDiff(stashEntry.Index)
		gui.renderString(g, "main", diff)
	}()
	return nil
}

func (gui *Gui) handleStashApply(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo(g, v, "apply")
}

func (gui *Gui) handleStashPop(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo(g, v, "pop")
}

func (gui *Gui) handleStashDrop(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, "Stash drop", "Are you sure you want to drop this stash entry?", func(g *gocui.Gui, v *gocui.View) error {
		return gui.stashDo(g, v, "drop")
	}, nil)
}

func (gui *Gui) stashDo(g *gocui.Gui, v *gocui.View, method string) error {
	stashEntry := gui.getSelectedStashEntry(v)
	if stashEntry == nil {
		return gui.createErrorPanel(g, "No stash to "+method)
	}
	if err := gui.GitCommand.StashDo(stashEntry.Index, method); err != nil {
		gui.createErrorPanel(g, err.Error())
	}
	gui.refreshStashEntries(g)
	return gui.refreshFiles(g)
}

func (gui *Gui) handleStashSave(g *gocui.Gui, filesView *gocui.View) error {
	gui.createPromptPanel(g, filesView, "Stash changes", func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.StashSave(gui.trimmedContent(v)); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		gui.refreshStashEntries(g)
		return gui.refreshFiles(g)
	})
	return nil
}
