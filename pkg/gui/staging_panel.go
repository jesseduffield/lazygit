package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) refreshStagingPanel(forceSecondaryFocused bool, selectedLineIdx int) error {
	gui.State.SplitMainPanel = true

	state := gui.State.Panels.LineByLine

	// We need to force focus here because the confirmation panel for safely staging lines does not return focus automatically.
	// This is because if we tell it to return focus it will unconditionally return it to the main panel which may not be what we want
	// e.g. in the event that there's nothing left to stage.
	if err := gui.switchFocus(gui.g, nil, gui.getMainView()); err != nil {
		return err
	}

	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return gui.handleStagingEscape(gui.g, nil)
	}

	if !file.HasUnstagedChanges && !file.HasStagedChanges {
		return gui.handleStagingEscape(gui.g, nil)
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
	diff := gui.GitCommand.Diff(file, true, secondaryFocused)
	secondaryDiff := gui.GitCommand.Diff(file, true, !secondaryFocused)

	// if we have e.g. a deleted file with nothing else to the diff will have only
	// 4-5 lines in which case we'll swap panels
	if len(strings.Split(diff, "\n")) < 5 {
		if len(strings.Split(secondaryDiff, "\n")) < 5 {
			return gui.handleStagingEscape(gui.g, nil)
		}
		secondaryFocused = !secondaryFocused
		diff, secondaryDiff = secondaryDiff, diff
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, secondaryFocused, selectedLineIdx)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleStagingEscape(gui.g, nil)
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

func (gui *Gui) handleStagingEscape(g *gocui.Gui, v *gocui.View) error {
	gui.handleEscapeLineByLinePanel()

	return gui.switchFocus(gui.g, nil, gui.getFilesView())
}

func (gui *Gui) handleStageSelection(g *gocui.Gui, v *gocui.View) error {
	return gui.applySelectionWithPrompt(false)
}

func (gui *Gui) handleResetSelection(g *gocui.Gui, v *gocui.View) error {
	return gui.applySelectionWithPrompt(true)
}

func (gui *Gui) applySelectionWithPrompt(reverse bool) error {
	state := gui.State.Panels.LineByLine

	if !reverse && state.SecondaryFocused {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("CantStageStaged"))
	} else if reverse && !state.SecondaryFocused && !gui.Config.GetUserConfig().GetBool("gui.skipUnstageLineWarning") {
		return gui.createConfirmationPanel(gui.g, gui.getMainView(), false, "unstage lines", "Are you sure you want to unstage these lines? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipUnstageLineWarning' to true", func(*gocui.Gui, *gocui.View) error {
			return gui.applySelection(reverse)
		}, nil)
	}

	return gui.applySelection(reverse)
}

func (gui *Gui) applySelection(reverse bool) error {
	state := gui.State.Panels.LineByLine

	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		return err
	}

	patch := commands.ModifiedPatchForRange(gui.Log, file.Name, state.Diff, state.FirstLineIdx, state.LastLineIdx, reverse, false)

	if patch == "" {
		return nil
	}

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	applyFlags := []string{}
	if !reverse || state.SecondaryFocused {
		applyFlags = append(applyFlags, "cached")
	}
	err = gui.GitCommand.ApplyPatch(patch, applyFlags...)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	if state.SelectMode == RANGE {
		state.SelectMode = LINE
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}
	if err := gui.refreshStagingPanel(false, -1); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handleMouseDownSecondaryWhileStaging(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	state.SecondaryFocused = !state.SecondaryFocused

	return gui.refreshStagingPanel(false, -1)
}
