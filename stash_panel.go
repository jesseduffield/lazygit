package main

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func refreshStashEntries(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("stash")
		if err != nil {
			panic(err)
		}
		state.StashEntries = getGitStashEntries()
		v.Clear()
		for _, stashEntry := range state.StashEntries {
			fmt.Fprintln(v, stashEntry.DisplayString)
		}
		return resetOrigin(v)
	})
	return nil
}

func getSelectedStashEntry(v *gocui.View) *StashEntry {
	if len(state.StashEntries) == 0 {
		return nil
	}
	lineNumber := getItemPosition(v)
	return &state.StashEntries[lineNumber]
}

func renderStashOptions(g *gocui.Gui) error {
	return renderOptionsMap(g, map[string]string{
		"space":   "apply",
		"g":       "pop",
		"d":       "drop",
		"← → ↑ ↓": "navigate",
	})
}

func handleStashEntrySelect(g *gocui.Gui, v *gocui.View) error {
	if err := renderStashOptions(g); err != nil {
		return err
	}
	go func() {
		stashEntry := getSelectedStashEntry(v)
		if stashEntry == nil {
			renderString(g, "main", "No stash entries")
			return
		}
		diff, _ := getStashEntryDiff(stashEntry.Index)
		renderString(g, "main", diff)
	}()
	return nil
}

func handleStashApply(g *gocui.Gui, v *gocui.View) error {
	return stashDo(g, v, "apply")
}

func handleStashPop(g *gocui.Gui, v *gocui.View) error {
	return stashDo(g, v, "pop")
}

func handleStashDrop(g *gocui.Gui, v *gocui.View) error {
	return createConfirmationPanel(g, v, "Stash drop", "Are you sure you want to drop this stash entry?", func(g *gocui.Gui, v *gocui.View) error {
		return stashDo(g, v, "drop")
	}, nil)
}

func stashDo(g *gocui.Gui, v *gocui.View, method string) error {
	stashEntry := getSelectedStashEntry(v)
	if stashEntry == nil {
		return createErrorPanel(g, "No stash to "+method)
	}
	if output, err := gitStashDo(stashEntry.Index, method); err != nil {
		createErrorPanel(g, output)
	}
	refreshStashEntries(g)
	return refreshFiles(g)
}

func handleStashSave(g *gocui.Gui, filesView *gocui.View) error {
	createPromptPanel(g, filesView, "Stash changes", nil, func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gitStashSave(trimmedContent(v)); err != nil {
			createErrorPanel(g, output)
		}
		refreshStashEntries(g)
		return refreshFiles(g)
	})
	return nil
}
