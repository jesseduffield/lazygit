package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/git"
)

// list panel functions

func (gui *Gui) getSelectedBranch() *commands.Branch {
	selectedLine := gui.State.Panels.Branches.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Branches[selectedLine]
}

// may want to standardise how these select methods work
func (gui *Gui) handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoBranchesThisRepo"))
	}
	branch := gui.getSelectedBranch()
	if err := gui.focusPoint(0, gui.State.Panels.Branches.SelectedLine, v); err != nil {
		return err
	}
	go func() {
		_ = gui.RenderSelectedBranchUpstreamDifferences()
	}()
	go func() {
		graph, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(graph, "fatal: ambiguous argument") {
			graph = gui.Tr.SLocalize("NoTrackingThisBranch")
		}
		_ = gui.renderString(g, "main", graph)
	}()
	return nil
}

func (gui *Gui) RenderSelectedBranchUpstreamDifferences() error {
	// here we tell the selected branch that it is selected.
	// this is necessary for showing stats on a branch that is selected, because
	// the displaystring function doesn't have access to gui state to tell if it's selected
	for i, branch := range gui.State.Branches {
		branch.Selected = i == gui.State.Panels.Branches.SelectedLine
	}

	branch := gui.getSelectedBranch()
	branch.Pushables, branch.Pullables = gui.GitCommand.GetBranchUpstreamDifferenceCount(branch.Name)
	return gui.renderListPanel(gui.getBranchesView(), gui.State.Branches)
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		builder, err := git.NewBranchListBuilder(gui.Log, gui.GitCommand)
		if err != nil {
			return err
		}
		gui.State.Branches = builder.Build()

		gui.refreshSelectedLine(&gui.State.Panels.Branches.SelectedLine, len(gui.State.Branches))
		if err := gui.RenderSelectedBranchUpstreamDifferences(); err != nil {
			return err
		}

		return gui.refreshStatus(g)
	})
	return nil
}

func (gui *Gui) handleBranchesNextLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Branches
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Branches), false)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleBranchSelect(gui.g, v)
}

func (gui *Gui) handleBranchesPrevLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Branches
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Branches), true)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleBranchSelect(gui.g, v)
}

// specific functions

func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Branches.SelectedLine == -1 {
		return nil
	}
	if gui.State.Panels.Branches.SelectedLine == 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("AlreadyCheckedOutBranch"))
	}
	branch := gui.getSelectedBranch()
	if err := gui.GitCommand.Checkout(branch.Name, false); err != nil {
		if err := gui.createErrorPanel(g, err.Error()); err != nil {
			return err
		}
	} else {
		gui.State.Panels.Branches.SelectedLine = 0
	}

	return gui.refreshSidePanels(g)
}

func (gui *Gui) handleCreatePullRequestPress(g *gocui.Gui, v *gocui.View) error {
	pullRequest := commands.NewPullRequest(gui.GitCommand)

	branch := gui.getSelectedBranch()
	if err := pullRequest.Create(branch); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}

	return nil
}

func (gui *Gui) handleGitFetch(g *gocui.Gui, v *gocui.View) error {
	if err := gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("FetchWait")); err != nil {
		return err
	}
	go func() {
		unamePassOpend, err := gui.fetch(g, v, true)
		gui.HandleCredentialsPopup(g, unamePassOpend, err)
	}()
	return nil
}

func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	message := gui.Tr.SLocalize("SureForceCheckout")
	title := gui.Tr.SLocalize("ForceCheckoutBranch")
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.Checkout(branch.Name, true); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	gui.createPromptPanel(g, v, gui.Tr.SLocalize("BranchName")+":", func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.Checkout(gui.trimmedContent(v), false); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshSidePanels(g)
	})
	return nil
}

func (gui *Gui) handleNewBranch(g *gocui.Gui, v *gocui.View) error {
	branch := gui.State.Branches[0]
	message := gui.Tr.TemplateLocalize(
		"NewBranchNameBranchOff",
		Teml{
			"branchName": branch.Name,
		},
	)
	gui.createPromptPanel(g, v, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.NewBranch(gui.trimmedContent(v)); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		gui.refreshSidePanels(g)
		return gui.handleBranchSelect(g, v)
	})
	return nil
}

func (gui *Gui) handleDeleteBranch(g *gocui.Gui, v *gocui.View) error {
	return gui.deleteBranch(g, v, false)
}

func (gui *Gui) handleForceDeleteBranch(g *gocui.Gui, v *gocui.View) error {
	return gui.deleteBranch(g, v, true)
}

func (gui *Gui) deleteBranch(g *gocui.Gui, v *gocui.View, force bool) error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}
	checkedOutBranch := gui.State.Branches[0]
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantDeleteCheckOutBranch"))
	}
	return gui.deleteNamedBranch(g, v, selectedBranch, force)
}

func (gui *Gui) deleteNamedBranch(g *gocui.Gui, v *gocui.View, selectedBranch *commands.Branch, force bool) error {
	title := gui.Tr.SLocalize("DeleteBranch")
	var messageID string
	if force {
		messageID = "ForceDeleteBranchMessage"
	} else {
		messageID = "DeleteBranchMessage"
	}
	message := gui.Tr.TemplateLocalize(
		messageID,
		Teml{
			"selectedBranchName": selectedBranch.Name,
		},
	)
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.DeleteBranch(selectedBranch.Name, force); err != nil {
			errMessage := err.Error()
			if !force && strings.Contains(errMessage, "is not fully merged") {
				return gui.deleteNamedBranch(g, v, selectedBranch, true)
			}
			return gui.createErrorPanel(g, errMessage)
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0].Name
	selectedBranch := gui.getSelectedBranch().Name
	if checkedOutBranch == selectedBranch {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantMergeBranchIntoItself"))
	}
	prompt := gui.Tr.TemplateLocalize(
		"ConfirmMerge",
		Teml{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranch,
		},
	)
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("MergingTitle"), prompt,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.Merge(selectedBranch)
			return gui.handleGenericMergeCommandResult(err)
		}, nil)
}

func (gui *Gui) handleRebase(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0].Name
	selectedBranch := gui.getSelectedBranch().Name
	if selectedBranch == checkedOutBranch {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantRebaseOntoSelf"))
	}
	prompt := gui.Tr.TemplateLocalize(
		"ConfirmRebase",
		Teml{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranch,
		},
	)
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("RebasingTitle"), prompt,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.RebaseBranch(selectedBranch)
			return gui.handleGenericMergeCommandResult(err)
		}, nil)
}

func (gui *Gui) handleFastForward(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}
	if branch.Pushables == "" {
		return nil
	}
	if branch.Pushables == "?" {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("FwdNoUpstream"))
	}
	if branch.Pushables != "0" {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("FwdCommitsToPush"))
	}
	upstream := "origin" // hardcoding for now
	message := gui.Tr.TemplateLocalize(
		"Fetching",
		Teml{
			"from": fmt.Sprintf("%s/%s", upstream, branch.Name),
			"to":   branch.Name,
		},
	)
	go func() {
		_ = gui.createLoaderPanel(gui.g, v, message)
		if err := gui.GitCommand.FastForward(branch.Name); err != nil {
			_ = gui.createErrorPanel(gui.g, err.Error())
		} else {
			_ = gui.closeConfirmationPrompt(gui.g)
			_ = gui.RenderSelectedBranchUpstreamDifferences()
		}
	}()
	return nil
}
