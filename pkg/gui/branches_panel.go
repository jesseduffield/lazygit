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
		return gui.createErrorPanel(g, "You have already checked out this branch")
	}
	branch := gui.getSelectedBranch(v)
	if output, err := gui.GitCommand.Checkout(branch.Name, false); err != nil {
		gui.createErrorPanel(g, output)
	}
	return gui.refreshSidePanels(g)
}

func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch(v)
	return gui.createConfirmationPanel(g, v, "Force Checkout Branch", "Are you sure you want force checkout? You will lose all local changes", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gui.GitCommand.Checkout(branch.Name, true); err != nil {
			gui.createErrorPanel(g, output)
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	gui.createPromptPanel(g, v, "Branch Name:", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gui.GitCommand.Checkout(gui.trimmedContent(v), false); err != nil {
			return gui.createErrorPanel(g, output)
		}
		return gui.refreshSidePanels(g)
	})
	return nil
}

func (gui *Gui) handleNewBranch(g *gocui.Gui, v *gocui.View) error {
	branch := gui.State.Branches[0]
	gui.createPromptPanel(g, v, "New Branch Name (Branch is off of "+branch.Name+")", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gui.GitCommand.NewBranch(gui.trimmedContent(v)); err != nil {
			return gui.createErrorPanel(g, output)
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
		return gui.createErrorPanel(g, "You cannot delete the checked out branch!")
	}
	return gui.createConfirmationPanel(g, v, "Delete Branch", "Are you sure you want delete the branch "+selectedBranch.Name+" ?", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gui.GitCommand.DeleteBranch(selectedBranch.Name); err != nil {
			return gui.createErrorPanel(g, output)
		}
		return gui.refreshSidePanels(g)
	}, nil)
}

func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	defer gui.refreshSidePanels(g)
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(g, "You cannot merge a branch into itself")
	}
	if output, err := gui.GitCommand.Merge(selectedBranch.Name); err != nil {
		return gui.createErrorPanel(g, output)
	}
	return nil
}

func (gui *Gui) getSelectedBranch(v *gocui.View) commands.Branch {
	lineNumber := gui.getItemPosition(v)
	return gui.State.Branches[lineNumber]
}

func (gui *Gui) renderBranchesOptions(g *gocui.Gui) error {
	return gui.renderOptionsMap(g, map[string]string{
		"space":   "checkout",
		"f":       "force checkout",
		"m":       "merge",
		"c":       "checkout by name",
		"n":       "new branch",
		"d":       "delete branch",
		"← → ↑ ↓": "navigate",
	})
}

// may want to standardise how these select methods work
func (gui *Gui) handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if err := gui.renderBranchesOptions(g); err != nil {
		return err
	}
	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {
		return gui.renderString(g, "main", "No branches for this repo")
	}
	go func() {
		branch := gui.getSelectedBranch(v)
		diff, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(diff, "fatal: ambiguous argument") {
			diff = "There is no tracking for this branch"
		}
		gui.renderString(g, "main", diff)
	}()
	return nil
}

// refreshStatus is called at the end of this because that's when we can
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
		return refreshStatus(g)
	})
	return nil
}
