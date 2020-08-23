package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
)

func (gui *Gui) refreshStagingPanel(forceSecondaryFocused bool, selectedLineIdx int) error {
	gui.splitMainPanel(true)

	state := gui.State.Panels.LineByLine

	file := gui.getSelectedFile()
	if file == nil {
		return gui.handleStagingEscape()
	}

	if !file.HasUnstagedChanges && !file.HasStagedChanges {
		return gui.handleStagingEscape()
	}

	secondaryFocused := false
	if forceSecondaryFocused {
		secondaryFocused = true
	} else {
		if state != nil {
			secondaryFocused = state.SecondaryFocused
		}
	}

	if (secondaryFocused && !file.HasStagedChanges) || (!secondaryFocused && !file.HasUnstagedChanges) {
		secondaryFocused = !secondaryFocused
	}

	if secondaryFocused {
		gui.getMainView().Title = gui.Tr.SLocalize("StagedChanges")
		gui.getSecondaryView().Title = gui.Tr.SLocalize("UnstagedChanges")
	} else {
		gui.getMainView().Title = gui.Tr.SLocalize("UnstagedChanges")
		gui.getSecondaryView().Title = gui.Tr.SLocalize("StagedChanges")
	}

	// note for custom diffs, we'll need to send a flag here saying not to use the custom diff
	diff := gui.GitCommand.WorktreeFileDiff(file, true, secondaryFocused)
	secondaryDiff := gui.GitCommand.WorktreeFileDiff(file, true, !secondaryFocused)

	// if we have e.g. a deleted file with nothing else to the diff will have only
	// 4-5 lines in which case we'll swap panels
	if len(strings.Split(diff, "\n")) < 5 {
		if len(strings.Split(secondaryDiff, "\n")) < 5 {
			return gui.handleStagingEscape()
		}
		secondaryFocused = !secondaryFocused
		diff, secondaryDiff = secondaryDiff, diff
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, secondaryFocused, selectedLineIdx)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleStagingEscape()
	}

	return nil
}

func (gui *Gui) handleTogglePanelClick(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	state.SecondaryFocused = !state.SecondaryFocused

	return gui.refreshStagingPanel(false, v.SelectedLineIdx())
}

func (gui *Gui) handleTogglePanel(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	state.SecondaryFocused = !state.SecondaryFocused
	return gui.refreshStagingPanel(false, -1)
}

func (gui *Gui) handleStagingEscape() error {
	gui.handleEscapeLineByLinePanel()

	return gui.switchContext(gui.Contexts.Files.Context)
}

func (gui *Gui) handleToggleStagedSelection(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	return gui.applySelection(state.SecondaryFocused)
}

func (gui *Gui) handleResetSelection(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	if state.SecondaryFocused {
		// for backwards compatibility
		return gui.applySelection(true)
	}

	if !gui.Config.GetUserConfig().GetBool("gui.skipUnstageLineWarning") {
		return gui.ask(askOpts{
			title:               gui.Tr.SLocalize("UnstageLinesTitle"),
			prompt:              gui.Tr.SLocalize("UnstageLinesPrompt"),
			handlersManageFocus: true,
			handleConfirm: func() error {
				if err := gui.switchContext(gui.Contexts.Staging.Context); err != nil {
					return err
				}

				return gui.applySelection(true)
			},
			handleClose: func() error {
				return gui.switchContext(gui.Contexts.Staging.Context)
			},
		})
	} else {
		return gui.applySelection(true)
	}
}

func (gui *Gui) applySelection(reverse bool) error {
	state := gui.State.Panels.LineByLine

	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	patch := patch.ModifiedPatchForRange(gui.Log, file.Name, state.Diff, state.FirstLineIdx, state.LastLineIdx, reverse, false)

	if patch == "" {
		return nil
	}

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	applyFlags := []string{}
	if !reverse || state.SecondaryFocused {
		applyFlags = append(applyFlags, "cached")
	}
	err := gui.GitCommand.ApplyPatch(patch, applyFlags...)
	if err != nil {
		return gui.surfaceError(err)
	}

	if state.SelectMode == RANGE {
		state.SelectMode = LINE
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []int{FILES}}); err != nil {
		return err
	}
	if err := gui.refreshStagingPanel(false, -1); err != nil {
		return err
	}
	return nil
}
