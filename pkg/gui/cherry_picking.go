package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// you can only copy from one context at a time, because the order and position of commits matter

func (gui *Gui) resetCherryPickingIfNecessary(context types.Context) error {
	oldContextKey := types.ContextKey(gui.State.Modes.CherryPicking.ContextKey)

	if oldContextKey != context.GetKey() {
		// need to reset the cherry picking mode
		gui.State.Modes.CherryPicking.ContextKey = string(context.GetKey())
		gui.State.Modes.CherryPicking.CherryPickedCommits = make([]*models.Commit, 0)

		return gui.rerenderContextViewIfPresent(oldContextKey)
	}

	return nil
}

func (gui *Gui) handleCopyCommit() error {
	// get currently selected commit, add the sha to state.
	context := gui.currentSideListContext()
	if context == nil {
		return nil
	}

	if err := gui.resetCherryPickingIfNecessary(context); err != nil {
		return err
	}

	item, ok := context.GetSelectedItem()
	if !ok {
		return nil
	}
	commit, ok := item.(*models.Commit)
	if !ok {
		return nil
	}

	// we will un-copy it if it's already copied
	for index, cherryPickedCommit := range gui.State.Modes.CherryPicking.CherryPickedCommits {
		if commit.Sha == cherryPickedCommit.Sha {
			gui.State.Modes.CherryPicking.CherryPickedCommits = append(gui.State.Modes.CherryPicking.CherryPickedCommits[0:index], gui.State.Modes.CherryPicking.CherryPickedCommits[index+1:]...)
			return context.HandleRender()
		}
	}

	gui.addCommitToCherryPickedCommits(context.GetPanelState().GetSelectedLineIdx())
	return context.HandleRender()
}

func (gui *Gui) cherryPickedCommitShaMap() map[string]bool {
	commitShaMap := map[string]bool{}
	for _, commit := range gui.State.Modes.CherryPicking.CherryPickedCommits {
		commitShaMap[commit.Sha] = true
	}
	return commitShaMap
}

func (gui *Gui) commitsListForContext() []*models.Commit {
	context := gui.currentSideListContext()
	if context == nil {
		return nil
	}

	// using a switch statement, but we should use polymorphism
	switch context.GetKey() {
	case BRANCH_COMMITS_CONTEXT_KEY:
		return gui.State.Commits
	case REFLOG_COMMITS_CONTEXT_KEY:
		return gui.State.FilteredReflogCommits
	case SUB_COMMITS_CONTEXT_KEY:
		return gui.State.SubCommits
	default:
		gui.c.Log.Errorf("no commit list for context %s", context.GetKey())
		return nil
	}
}

func (gui *Gui) addCommitToCherryPickedCommits(index int) {
	commitShaMap := gui.cherryPickedCommitShaMap()
	commitsList := gui.commitsListForContext()
	commitShaMap[commitsList[index].Sha] = true

	newCommits := []*models.Commit{}
	for _, commit := range commitsList {
		if commitShaMap[commit.Sha] {
			// duplicating just the things we need to put in the rebase TODO list
			newCommits = append(newCommits, &models.Commit{Name: commit.Name, Sha: commit.Sha})
		}
	}

	gui.State.Modes.CherryPicking.CherryPickedCommits = newCommits
}

func (gui *Gui) handleCopyCommitRange() error {
	// get currently selected commit, add the sha to state.
	context := gui.currentSideListContext()
	if context == nil {
		return nil
	}

	if err := gui.resetCherryPickingIfNecessary(context); err != nil {
		return err
	}

	commitShaMap := gui.cherryPickedCommitShaMap()
	commitsList := gui.commitsListForContext()
	selectedLineIdx := context.GetPanelState().GetSelectedLineIdx()

	if selectedLineIdx > len(commitsList)-1 {
		return nil
	}

	// find the last commit that is copied that's above our position
	// if there are none, startIndex = 0
	startIndex := 0
	for index, commit := range commitsList[0:selectedLineIdx] {
		if commitShaMap[commit.Sha] {
			startIndex = index
		}
	}

	for index := startIndex; index <= selectedLineIdx; index++ {
		gui.addCommitToCherryPickedCommits(index)
	}

	return context.HandleRender()
}

// HandlePasteCommits begins a cherry-pick rebase with the commits the user has copied
func (gui *Gui) HandlePasteCommits() error {
	return gui.c.Ask(popup.AskOpts{
		Title:  gui.c.Tr.CherryPick,
		Prompt: gui.c.Tr.SureCherryPick,
		HandleConfirm: func() error {
			return gui.c.WithWaitingStatus(gui.c.Tr.CherryPickingStatus, func() error {
				gui.c.LogAction(gui.c.Tr.Actions.CherryPick)
				err := gui.git.Rebase.CherryPickCommits(gui.State.Modes.CherryPicking.CherryPickedCommits)
				return gui.checkMergeOrRebase(err)
			})
		},
	})
}

func (gui *Gui) exitCherryPickingMode() error {
	contextKey := types.ContextKey(gui.State.Modes.CherryPicking.ContextKey)

	gui.State.Modes.CherryPicking.ContextKey = ""
	gui.State.Modes.CherryPicking.CherryPickedCommits = nil

	if contextKey == "" {
		gui.c.Log.Warn("context key blank when trying to exit cherry picking mode")
		return nil
	}

	return gui.rerenderContextViewIfPresent(contextKey)
}

func (gui *Gui) rerenderContextViewIfPresent(contextKey types.ContextKey) error {
	if contextKey == "" {
		return nil
	}

	context := gui.mustContextForContextKey(contextKey)

	viewName := context.GetViewName()

	view, err := gui.g.View(viewName)
	if err != nil {
		gui.c.Log.Error(err)
		return nil
	}

	if types.ContextKey(view.Context) == contextKey {
		if err := context.HandleRender(); err != nil {
			return err
		}
	}

	return nil
}
