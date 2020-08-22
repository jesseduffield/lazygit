package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) handleCopyCommit(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	// get currently selected commit, add the sha to state.
	commit := gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx]

	// we will un-copy it if it's already copied
	for index, cherryPickedCommit := range gui.State.CherryPickedCommits {
		if commit.Sha == cherryPickedCommit.Sha {
			gui.State.CherryPickedCommits = append(gui.State.CherryPickedCommits[0:index], gui.State.CherryPickedCommits[index+1:]...)
			return gui.Contexts.BranchCommits.Context.HandleRender()
		}
	}

	gui.addCommitToCherryPickedCommits(gui.State.Panels.Commits.SelectedLineIdx)
	return gui.Contexts.BranchCommits.Context.HandleRender()
}

func (gui *Gui) cherryPickedCommitShaMap() map[string]bool {
	commitShaMap := map[string]bool{}
	for _, commit := range gui.State.CherryPickedCommits {
		commitShaMap[commit.Sha] = true
	}
	return commitShaMap
}

func (gui *Gui) addCommitToCherryPickedCommits(index int) {
	commitShaMap := gui.cherryPickedCommitShaMap()
	commitShaMap[gui.State.Commits[index].Sha] = true

	newCommits := []*commands.Commit{}
	for _, commit := range gui.State.Commits {
		if commitShaMap[commit.Sha] {
			// duplicating just the things we need to put in the rebase TODO list
			newCommits = append(newCommits, &commands.Commit{Name: commit.Name, Sha: commit.Sha})
		}
	}

	gui.State.CherryPickedCommits = newCommits
}

func (gui *Gui) handleCopyCommitRange(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	commitShaMap := gui.cherryPickedCommitShaMap()

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
func (gui *Gui) HandlePasteCommits(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("CherryPick"),
		prompt:             gui.Tr.SLocalize("SureCherryPick"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("CherryPickingStatus"), func() error {
				err := gui.GitCommand.CherryPickCommits(gui.State.CherryPickedCommits)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}
