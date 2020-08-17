package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

// list panel functions

func (gui *Gui) getSelectedStashEntry() *commands.StashEntry {
	selectedLine := gui.State.Panels.Stash.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.StashEntries[selectedLine]
}

func (gui *Gui) handleStashEntrySelect() error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.splitMainPanel(false)

	gui.getMainView().Title = "Stash"

	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return gui.newStringTask("main", gui.Tr.SLocalize("NoStashEntries"))
	}

	if gui.inDiffMode() {
		return gui.renderDiff()
	}

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowStashEntryCmdStr(stashEntry.Index),
	)
	if err := gui.newPtyTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) refreshStashEntries(g *gocui.Gui) error {
	gui.State.StashEntries = gui.GitCommand.GetStashEntries(gui.State.FilterPath)

	gui.refreshSelectedLine(&gui.State.Panels.Stash.SelectedLine, len(gui.State.StashEntries))

	stashView := gui.getStashView()

	displayStrings := presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Diff.Ref)
	gui.renderDisplayStrings(stashView, displayStrings)

	return gui.resetOrigin(stashView)
}

// specific functions

func (gui *Gui) handleStashApply(g *gocui.Gui, v *gocui.View) error {
	skipStashWarning := gui.Config.GetUserConfig().GetBool("gui.skipStashWarning")

	apply := func() error {
		return gui.stashDo("apply")
	}

	if skipStashWarning {
		return apply()
	}

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("StashApply"),
		prompt:             gui.Tr.SLocalize("SureApplyStashEntry"),
		handleConfirm: func() error {
			return apply()
		},
	})
}

func (gui *Gui) handleStashPop(g *gocui.Gui, v *gocui.View) error {
	skipStashWarning := gui.Config.GetUserConfig().GetBool("gui.skipStashWarning")

	pop := func() error {
		return gui.stashDo("pop")
	}

	if skipStashWarning {
		return pop()
	}

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("StashPop"),
		prompt:             gui.Tr.SLocalize("SurePopStashEntry"),
		handleConfirm: func() error {
			return pop()
		},
	})
}

func (gui *Gui) handleStashDrop(g *gocui.Gui, v *gocui.View) error {
	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("StashDrop"),
		prompt:             gui.Tr.SLocalize("SureDropStashEntry"),
		handleConfirm: func() error {
			return gui.stashDo("drop")
		},
	})
}

func (gui *Gui) stashDo(method string) error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		errorMessage := gui.Tr.TemplateLocalize(
			"NoStashTo",
			Teml{
				"method": method,
			},
		)
		return gui.createErrorPanel(errorMessage)
	}
	if err := gui.GitCommand.StashDo(stashEntry.Index, method); err != nil {
		return gui.surfaceError(err)
	}
	return gui.refreshSidePanels(refreshOptions{scope: []int{STASH, FILES}})
}

func (gui *Gui) handleStashSave(stashFunc func(message string) error) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoTrackedStagedFilesStash"))
	}
	return gui.prompt(gui.getFilesView(), gui.Tr.SLocalize("StashChanges"), "", func(stashComment string) error {
		if err := stashFunc(stashComment); err != nil {
			return gui.surfaceError(err)
		}
		return gui.refreshSidePanels(refreshOptions{scope: []int{STASH, FILES}})
	})
}
