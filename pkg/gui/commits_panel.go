package gui

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/go-errors/errors"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/git"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedCommit(g *gocui.Gui) *commands.Commit {
	selectedLine := gui.State.Panels.Commits.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Commits[selectedLine]
}

func (gui *Gui) handleCommitSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	if err := gui.focusPoint(0, gui.State.Panels.Commits.SelectedLine, len(gui.State.Commits), v); err != nil {
		return err
	}

	// if specific diff mode is on, don't show diff
	if gui.State.Panels.Commits.SpecificDiffMode {
		return nil
	}

	commitText, err := gui.GitCommand.Show(commit.Sha)
	if err != nil {
		return err
	}
	return gui.renderString(g, "main", commitText)
}

func (gui *Gui) refreshCommits(g *gocui.Gui) error {
	g.Update(func(*gocui.Gui) error {
		builder, err := git.NewCommitListBuilder(gui.Log, gui.GitCommand, gui.OSCommand, gui.Tr, gui.State.CherryPickedCommits, gui.State.DiffEntries)
		if err != nil {
			return err
		}
		commits, err := builder.GetCommits()
		if err != nil {
			return err
		}
		gui.State.Commits = commits

		gui.refreshSelectedLine(&gui.State.Panels.Commits.SelectedLine, len(gui.State.Commits))

		isFocused := gui.g.CurrentView().Name() == "commits"
		list, err := utils.RenderList(gui.State.Commits, isFocused)
		if err != nil {
			return err
		}

		v := gui.getCommitsView()
		v.Clear()
		fmt.Fprint(v, list)

		gui.refreshStatus(g)
		if g.CurrentView() == v {
			gui.handleCommitSelect(g, v)
		}
		if g.CurrentView() == gui.getCommitFilesView() {
			return gui.refreshCommitFilesView()
		}
		return nil
	})
	return nil
}

func (gui *Gui) handleCommitsNextLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Commits
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Commits), false)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleCommitSelect(gui.g, v)
}

func (gui *Gui) handleCommitsPrevLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Commits
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Commits), true)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleCommitSelect(gui.g, v)
}

// specific functions

func (gui *Gui) handleResetToCommit(g *gocui.Gui, commitView *gocui.View) error {
	return gui.createConfirmationPanel(g, commitView, gui.Tr.SLocalize("ResetToCommit"), gui.Tr.SLocalize("SureResetThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		commit := gui.getSelectedCommit(g)
		if commit == nil {
			panic(errors.New(gui.Tr.SLocalize("NoCommitsThisBranch")))
		}

		if err := gui.GitCommand.ResetToCommit(commit.Sha, "mixed"); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			panic(err)
		}
		if err := gui.refreshFiles(); err != nil {
			panic(err)
		}
		gui.resetOrigin(commitView)
		gui.State.Panels.Commits.SelectedLine = 0
		return gui.handleCommitSelect(g, commitView)
	}, nil)
}

func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	applied, err := gui.handleMidRebaseCommand("squash")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Squash"), gui.Tr.SLocalize("SureSquashThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("SquashingStatus"), func() error {
			err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, "squash")
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
	return nil
}

// TODO: move to files panel
func (gui *Gui) anyUnStagedChanges(files []*commands.File) bool {
	for _, file := range files {
		if file.Tracked && file.HasUnstagedChanges {
			return true
		}
	}
	return false
}

func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {
	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	applied, err := gui.handleMidRebaseCommand("fixup")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Fixup"), gui.Tr.SLocalize("SureFixupThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("FixingStatus"), func() error {
			err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, "fixup")
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
	return nil
}

func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	applied, err := gui.handleMidRebaseCommand("reword")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	if gui.State.Panels.Commits.SelectedLine != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}
	return gui.createPromptPanel(g, v, gui.Tr.SLocalize("renameCommit"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.RenameCommit(v.Buffer()); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			panic(err)
		}
		return gui.handleCommitSelect(g, v)
	})
}

func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {
	applied, err := gui.handleMidRebaseCommand("reword")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	subProcess, err := gui.GitCommand.RewordCommit(gui.State.Commits, gui.State.Panels.Commits.SelectedLine)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	if subProcess != nil {
		gui.SubProcess = subProcess
		return gui.Errors.ErrSubProcess
	}

	return nil
}

// handleMidRebaseCommand sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (gui *Gui) handleMidRebaseCommand(action string) (bool, error) {
	selectedCommit := gui.State.Commits[gui.State.Panels.Commits.SelectedLine]
	if selectedCommit.Status != "rebasing" {
		return false, nil
	}

	// for now we do not support setting 'reword' because it requires an editor
	// and that means we either unconditionally wait around for the subprocess to ask for
	// our input or we set a lazygit client as the EDITOR env variable and have it
	// request us to edit the commit message when prompted.
	if action == "reword" {
		return true, gui.createErrorPanel(gui.g, gui.Tr.SLocalize("rewordNotSupported"))
	}

	if err := gui.GitCommand.EditRebaseTodo(gui.State.Panels.Commits.SelectedLine, action); err != nil {
		return false, gui.createErrorPanel(gui.g, err.Error())
	}
	return true, gui.refreshCommits(gui.g)
}

// handleMoveTodoDown like handleMidRebaseCommand but for moving an item up in the todo list
func (gui *Gui) handleMoveTodoDown(index int) (bool, error) {
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status != "rebasing" {
		return false, nil
	}
	if gui.State.Commits[index+1].Status != "rebasing" {
		return true, nil
	}
	if err := gui.GitCommand.MoveTodoDown(index); err != nil {
		return true, gui.createErrorPanel(gui.g, err.Error())
	}
	return true, gui.refreshCommits(gui.g)
}

func (gui *Gui) handleCommitDelete(g *gocui.Gui, v *gocui.View) error {
	applied, err := gui.handleMidRebaseCommand("drop")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.createConfirmationPanel(gui.g, v, gui.Tr.SLocalize("DeleteCommitTitle"), gui.Tr.SLocalize("DeleteCommitPrompt"), func(*gocui.Gui, *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("DeletingStatus"), func() error {
			err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, "drop")
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
}

func (gui *Gui) handleCommitMoveDown(g *gocui.Gui, v *gocui.View) error {
	index := gui.State.Panels.Commits.SelectedLine
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if gui.State.Commits[index+1].Status != "rebasing" {
			return nil
		}
		if err := gui.GitCommand.MoveTodoDown(index); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		gui.State.Panels.Commits.SelectedLine++
		return gui.refreshCommits(gui.g)
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("MovingStatus"), func() error {
		err := gui.GitCommand.MoveCommitDown(gui.State.Commits, index)
		if err == nil {
			gui.State.Panels.Commits.SelectedLine++
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitMoveUp(g *gocui.Gui, v *gocui.View) error {
	index := gui.State.Panels.Commits.SelectedLine
	if index == 0 {
		return nil
	}
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if err := gui.GitCommand.MoveTodoDown(index - 1); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		gui.State.Panels.Commits.SelectedLine--
		return gui.refreshCommits(gui.g)
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("MovingStatus"), func() error {
		err := gui.GitCommand.MoveCommitDown(gui.State.Commits, index-1)
		if err == nil {
			gui.State.Panels.Commits.SelectedLine--
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitEdit(g *gocui.Gui, v *gocui.View) error {
	applied, err := gui.handleMidRebaseCommand("edit")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		err = gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, "edit")
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitAmendTo(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(gui.g, v, gui.Tr.SLocalize("AmendCommitTitle"), gui.Tr.SLocalize("AmendCommitPrompt"), func(*gocui.Gui, *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("AmendingStatus"), func() error {
			err := gui.GitCommand.AmendTo(gui.State.Commits[gui.State.Panels.Commits.SelectedLine].Sha)
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
}

func (gui *Gui) handleCommitPick(g *gocui.Gui, v *gocui.View) error {
	applied, err := gui.handleMidRebaseCommand("pick")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	// at this point we aren't actually rebasing so we will interpret this as an
	// attempt to pull. We might revoke this later after enabling configurable keybindings
	return gui.pullFiles(g, v)
}

func (gui *Gui) handleCommitRevert(g *gocui.Gui, v *gocui.View) error {
	if err := gui.GitCommand.Revert(gui.State.Commits[gui.State.Panels.Commits.SelectedLine].Sha); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	gui.State.Panels.Commits.SelectedLine++
	return gui.refreshCommits(gui.g)
}

func (gui *Gui) handleCopyCommit(g *gocui.Gui, v *gocui.View) error {
	// get currently selected commit, add the sha to state.
	commit := gui.State.Commits[gui.State.Panels.Commits.SelectedLine]

	// we will un-copy it if it's already copied
	for index, cherryPickedCommit := range gui.State.CherryPickedCommits {
		if commit.Sha == cherryPickedCommit.Sha {
			gui.State.CherryPickedCommits = append(gui.State.CherryPickedCommits[0:index], gui.State.CherryPickedCommits[index+1:]...)
			return gui.refreshCommits(gui.g)
		}
	}

	gui.addCommitToCherryPickedCommits(gui.State.Panels.Commits.SelectedLine)
	return gui.refreshCommits(gui.g)
}

func (gui *Gui) addCommitToCherryPickedCommits(index int) {
	// not super happy with modifying the state of the Commits array here
	// but the alternative would be very tricky
	gui.State.Commits[index].Copied = true

	newCommits := []*commands.Commit{}
	for _, commit := range gui.State.Commits {
		if commit.Copied {
			// duplicating just the things we need to put in the rebase TODO list
			newCommits = append(newCommits, &commands.Commit{Name: commit.Name, Sha: commit.Sha})
		}
	}

	gui.State.CherryPickedCommits = newCommits
}

func (gui *Gui) handleCopyCommitRange(g *gocui.Gui, v *gocui.View) error {
	// whenever I add a commit, I need to make sure I retain its order

	// find the last commit that is copied that's above our position
	// if there are none, startIndex = 0
	startIndex := 0
	for index, commit := range gui.State.Commits[0:gui.State.Panels.Commits.SelectedLine] {
		if commit.Copied {
			startIndex = index
		}
	}

	gui.Log.Info("commit copy start index: " + strconv.Itoa(startIndex))

	for index := startIndex; index <= gui.State.Panels.Commits.SelectedLine; index++ {
		gui.addCommitToCherryPickedCommits(index)
	}

	return gui.refreshCommits(gui.g)
}

// HandlePasteCommits begins a cherry-pick rebase with the commits the user has copied
func (gui *Gui) HandlePasteCommits(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("CherryPick"), gui.Tr.SLocalize("SureCherryPick"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("CherryPickingStatus"), func() error {
			err := gui.GitCommand.CherryPickCommits(gui.State.CherryPickedCommits)
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
}

func (gui *Gui) handleSwitchToCommitFilesPanel(g *gocui.Gui, v *gocui.View) error {
	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.switchFocus(g, v, gui.getCommitFilesView())
}

func (gui *Gui) handleToggleDiffCommit(g *gocui.Gui, v *gocui.View) error {
	selectLimit := 2

	// get selected commit
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	// if already selected commit delete
	if idx, has := gui.hasCommit(gui.State.DiffEntries, commit.Sha); has {
		gui.State.DiffEntries = gui.unchooseCommit(gui.State.DiffEntries, idx)
	} else {
		if len(gui.State.DiffEntries) == selectLimit {
			gui.State.DiffEntries = gui.unchooseCommit(gui.State.DiffEntries, 0)
		}
		gui.State.DiffEntries = append(gui.State.DiffEntries, commit)
	}

	gui.setDiffMode()

	// if selected two commits, display diff between
	if len(gui.State.DiffEntries) == selectLimit {
		commitText, err := gui.GitCommand.DiffCommits(gui.State.DiffEntries[0].Sha, gui.State.DiffEntries[1].Sha)

		if err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}

		return gui.renderString(g, "main", commitText)
	}

	return nil
}

func (gui *Gui) setDiffMode() {
	v := gui.getCommitsView()
	if len(gui.State.DiffEntries) != 0 {
		gui.State.Panels.Commits.SpecificDiffMode = true
		v.Title = gui.Tr.SLocalize("CommitsDiffTitle")
	} else {
		gui.State.Panels.Commits.SpecificDiffMode = false
		v.Title = gui.Tr.SLocalize("CommitsTitle")
	}

	gui.refreshCommits(gui.g)
}

func (gui *Gui) hasCommit(commits []*commands.Commit, target string) (int, bool) {
	for idx, commit := range commits {
		if commit.Sha == target {
			return idx, true
		}
	}
	return -1, false
}

func (gui *Gui) unchooseCommit(commits []*commands.Commit, i int) []*commands.Commit {
	return append(commits[:i], commits[i+1:]...)
}

func (gui *Gui) handleCreateFixupCommit(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return nil
	}

	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("CreateFixupCommit"), gui.Tr.TemplateLocalize(
		"SureCreateFixupCommit",
		Teml{
			"commit": commit.Sha,
		},
	), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.CreateFixupCommit(commit.Sha); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}

		return gui.refreshSidePanels(gui.g)
	}, nil)
}

func (gui *Gui) handleSquashAllAboveFixupCommits(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return nil
	}

	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("SquashAboveCommits"), gui.Tr.TemplateLocalize(
		"SureSquashAboveCommits",
		Teml{
			"commit": commit.Sha,
		},
	), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("SquashingStatus"), func() error {
			err := gui.GitCommand.SquashAllAboveFixupCommits(commit.Sha)
			return gui.handleGenericMergeCommandResult(err)
		})
	}, nil)
}

type resetOption struct {
	description string
	command     string
}

// GetDisplayStrings is a function.
func (r *resetOption) GetDisplayStrings(isFocused bool) []string {
	return []string{r.description, color.New(color.FgRed).Sprint(r.command)}
}

func (gui *Gui) handleCreateCommitResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	strengths := []string{"soft", "mixed", "hard"}
	options := make([]*resetOption, len(strengths))
	for i, strength := range strengths {
		options[i] = &resetOption{
			description: fmt.Sprintf("%s reset", strength),
			command:     fmt.Sprintf("reset --%s %s", strength, commit.Sha),
		}
	}

	handleMenuPress := func(index int) error {
		if err := gui.GitCommand.ResetToCommit(commit.Sha, strengths[index]); err != nil {
			return err
		}

		if err := gui.refreshCommits(g); err != nil {
			return err
		}
		if err := gui.refreshFiles(); err != nil {
			return err
		}
		if err := gui.resetOrigin(gui.getCommitsView()); err != nil {
			return err
		}

		gui.State.Panels.Commits.SelectedLine = 0
		return gui.handleCommitSelect(g, gui.getCommitsView())
	}

	return gui.createMenu(fmt.Sprintf("%s %s", gui.Tr.SLocalize("resetTo"), commit.Sha), options, len(options), handleMenuPress)
}
