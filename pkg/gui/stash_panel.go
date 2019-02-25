package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedStashEntry(v *gocui.View) *commands.StashEntry {
	selectedLine := gui.State.Panels.Stash.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.StashEntries[selectedLine]
}

func (gui *Gui) handleStashEntrySelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	stashEntry := gui.getSelectedStashEntry(v)
	if stashEntry == nil {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoStashEntries"))
	}
	if err := gui.focusPoint(0, gui.State.Panels.Stash.SelectedLine, v); err != nil {
		return err
	}
	go func() {
		// doing this asynchronously cos it can take time
		diff, _ := gui.GitCommand.GetStashEntryDiff(stashEntry.Index)
		_ = gui.renderString(g, "main", diff)
	}()
	return nil
}

func (gui *Gui) refreshStashEntries(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		gui.State.StashEntries = gui.GitCommand.GetStashEntries()

		gui.refreshSelectedLine(&gui.State.Panels.Stash.SelectedLine, len(gui.State.StashEntries))

		isFocused := gui.g.CurrentView().Name() == "stash"
		list, err := utils.RenderList(gui.State.StashEntries, isFocused)
		if err != nil {
			return err
		}

		v := gui.getStashView()
		v.Clear()
		fmt.Fprint(v, list)

		if err := gui.resetOrigin(v); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (gui *Gui) handleStashNextLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Stash
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.StashEntries), false)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleStashEntrySelect(gui.g, v)
}

func (gui *Gui) handleStashPrevLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Stash
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.StashEntries), true)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleStashEntrySelect(gui.g, v)
}

// specific functions

func (gui *Gui) handleStashApply(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo(g, v, "apply")
}

func (gui *Gui) handleStashPop(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo(g, v, "pop")
}

func (gui *Gui) handleStashDrop(g *gocui.Gui, v *gocui.View) error {
	title := gui.Tr.SLocalize("StashDrop")
	message := gui.Tr.SLocalize("SureDropStashEntry")
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		return gui.stashDo(g, v, "drop")
	}, nil)
}

func (gui *Gui) stashDo(g *gocui.Gui, v *gocui.View, method string) error {
	stashEntry := gui.getSelectedStashEntry(v)
	if stashEntry == nil {
		errorMessage := gui.Tr.TemplateLocalize(
			"NoStashTo",
			Teml{
				"method": method,
			},
		)
		return gui.createErrorPanel(g, errorMessage)
	}
	if err := gui.GitCommand.StashDo(stashEntry.Index, method); err != nil {
		gui.createErrorPanel(g, err.Error())
	}
	gui.refreshStashEntries(g)
	return gui.refreshFiles()
}

func (gui *Gui) handleStashSave(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoTrackedStagedFilesStash"))
	}
	gui.createPromptPanel(g, filesView, gui.Tr.SLocalize("StashChanges"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.StashSave(gui.trimmedContent(v)); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		gui.refreshStashEntries(g)
		return gui.refreshFiles()
	})
	return nil
}
