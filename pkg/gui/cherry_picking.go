package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// you can only copy from one context at a time, because the order and position of commits matter

func (gui *Gui) handleCopyCommit() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	// get currently selected commit, add the sha to state.
	context := gui.currentSideContext()
	if context == nil {
		return nil
	}

	commit, ok := context.SelectedItem().(*commands.Commit)
	if !ok {
		gui.Log.Error("type cast failed for handling copy commit")
	}
	if commit == nil {
		return nil
	}

	// we will un-copy it if it's already copied
	for index, cherryPickedCommit := range gui.State.Modes.CherryPicking.CherryPickedCommits {
		if commit.Sha == cherryPickedCommit.Sha {
			gui.State.Modes.CherryPicking.CherryPickedCommits = append(gui.State.Modes.CherryPicking.CherryPickedCommits[0:index], gui.State.Modes.CherryPicking.CherryPickedCommits[index+1:]...)
			return context.HandleRender()
		}
	}

	gui.addCommitToCherryPickedCommits(gui.State.Panels.Commits.SelectedLineIdx)
	return context.HandleRender()
}

func (gui *Gui) CherryPickedCommitShaMap() map[string]bool {
	commitShaMap := map[string]bool{}
	for _, commit := range gui.State.Modes.CherryPicking.CherryPickedCommits {
		commitShaMap[commit.Sha] = true
	}
	return commitShaMap
}

func (gui *Gui) addCommitToCherryPickedCommits(index int) {
	commitShaMap := gui.CherryPickedCommitShaMap()
	commitShaMap[gui.State.Commits[index].Sha] = true

	newCommits := []*commands.Commit{}
	for _, commit := range gui.State.Commits {
		if commitShaMap[commit.Sha] {
			// duplicating just the things we need to put in the rebase TODO list
			newCommits = append(newCommits, &commands.Commit{Name: commit.Name, Sha: commit.Sha})
		}
	}

	gui.State.Modes.CherryPicking.CherryPickedCommits = newCommits
}

func (gui *Gui) handleCopyCommitRange() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	commitShaMap := gui.CherryPickedCommitShaMap()

	// find the last commit that is copied that's above our position
	// if there are none, startIndex = 0
	startIndex := 0
	for index, commit := range gui.State.Commits[0:gui.State.Panels.Commits.SelectedLineIdx] {
		if commitShaMap[commit.Sha] {
			startIndex = index
		}
	}

	for index := startIndex; index <= gui.State.Panels.Commits.SelectedLineIdx; index++ {
		gui.addCommitToCherryPickedCommits(index)
	}

	return gui.Contexts.BranchCommits.Context.HandleRender()
}

// HandlePasteCommits begins a cherry-pick rebase with the commits the user has copied
func (gui *Gui) HandlePasteCommits() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.ask(askOpts{
		returnToView:       gui.getCommitsView(),
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("CherryPick"),
		prompt:             gui.Tr.SLocalize("SureCherryPick"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("CherryPickingStatus"), func() error {
				err := gui.GitCommand.CherryPickCommits(gui.State.Modes.CherryPicking.CherryPickedCommits)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}
