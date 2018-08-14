package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/git"
)

func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	index := gui.getItemPosition(v)
	if index == 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("AlreadyCheckedOutBranch", "You have already checked out this branch"))
	}
	branch := gui.getSelectedBranch(v)
	if err := gui.GitCommand.Checkout(branch.Name, false); err != nil {
		gui.createErrorPanel(g, err.Error())
	}
	return gui.refreshSidePanels(g)
}

func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch(v)
	message := gui.Tr.SLocalize("SureForceCheckout", "Are you sure you want force checkout? You will lose all local changes")
	title := gui.Tr.SLocalize("ForceCheckoutBranch", "Force Checkout Branch")
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.Checkout(branch.Name, true); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	gui.createPromptPanel(g, v, gui.Tr.SLocalize("BranchName", "Branch Name")+":", func(g *gocui.Gui, v *gocui.View) error {
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
		"New Branch Name (Branch is off of {{.branchName}})",
		map[string]interface{}{
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
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantDeleteCheckOutBranch", "You cannot delete the checked out branch!"))
	}
	message := gui.Tr.TemplateLocalize(
		"DeleteBranchMessage",
		"Are you sure you want delete the branch {{.selectedBranchName}} ?",
		map[string]interface{}{
			"selectedBranchName": selectedBranch.Name,
		},
	)
	title := gui.Tr.SLocalize("DeleteBranch", "Delete Branch")
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.DeleteBranch(selectedBranch.Name); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	defer gui.refreshSidePanels(g)
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantMergeBranchIntoItself", "You cannot merge a branch into itself"))
	}
	if err := gui.GitCommand.Merge(selectedBranch.Name); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	return nil
}

func (gui *Gui) getSelectedBranch(v *gocui.View) commands.Branch {
	lineNumber := gui.getItemPosition(v)
	return gui.State.Branches[lineNumber]
}

func (gui *Gui) renderBranchesOptions(g *gocui.Gui) error {
	return gui.renderOptionsMap(g, map[string]string{
		"space":   gui.Tr.SLocalize("checkout", "checkout"),
		"f":       gui.Tr.SLocalize("forceCheckout", "force checkout"),
		"m":       gui.Tr.SLocalize("merge", "merge"),
		"c":       gui.Tr.SLocalize("checkoutByName", "checkout by name"),
		"n":       gui.Tr.SLocalize("newBranch", "new branch"),
		"d":       gui.Tr.SLocalize("deleteBranch", "delete branch"),
		"← → ↑ ↓": gui.Tr.SLocalize("navigate", "navigate"),
	})
}

// may want to standardise how these select methods work
func (gui *Gui) handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if err := gui.renderBranchesOptions(g); err != nil {
		return err
	}
	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoBranchesThisRepo", "No branches for this repo"))
	}
	go func() {
		branch := gui.getSelectedBranch(v)
		diff, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(diff, "fatal: ambiguous argument") {
			diff = gui.Tr.SLocalize("NoTrackingThisRepo", "There is no tracking for this branch")
		}
		gui.renderString(g, "main", diff)
	}()
	return nil
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("branches")
		if err != nil {
			panic(err)
		}
		builder, err := git.NewBranchListBuilder(gui.Log, gui.GitCommand)
		if err != nil {
			return err
		}
		gui.State.Branches = builder.Build()
		v.Clear()
		for _, branch := range gui.State.Branches {
			fmt.Fprintln(v, branch.GetDisplayString())
		}
		gui.resetOrigin(v)
		return gui.refreshStatus(g)
	})
	return nil
}
