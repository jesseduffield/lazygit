package main

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	index := getItemPosition(v)
	if index == 0 {
		return createErrorPanel(g, "You have already checked out this branch")
	}
	branch := getSelectedBranch(v)
	if output, err := gitCheckout(branch.Name, false); err != nil {
		createErrorPanel(g, output)
	}
	return refreshSidePanels(g)
}

func handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := getSelectedBranch(v)
	return createConfirmationPanel(g, v, "Force Checkout Branch", "Are you sure you want force checkout? You will lose all local changes", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gitCheckout(branch.Name, true); err != nil {
			createErrorPanel(g, output)
		}
		return refreshSidePanels(g)
	}, nil)
}

func handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	createPromptPanel(g, v, "Branch Name:", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gitCheckout(trimmedContent(v), false); err != nil {
			return createErrorPanel(g, output)
		}
		return refreshSidePanels(g)
	})
	return nil
}

func handleNewBranch(g *gocui.Gui, v *gocui.View) error {
	branch := state.Branches[0]
	createPromptPanel(g, v, "New Branch Name (Branch is off of "+branch.Name+")", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gitNewBranch(trimmedContent(v)); err != nil {
			return createErrorPanel(g, output)
		}
		refreshSidePanels(g)
		return handleBranchSelect(g, v)
	})
	return nil
}

func handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := state.Branches[0]
	selectedBranch := getSelectedBranch(v)
	defer refreshSidePanels(g)
	if checkedOutBranch.Name == selectedBranch.Name {
		return createErrorPanel(g, "You cannot merge a branch into itself")
	}
	if output, err := gitMerge(selectedBranch.Name); err != nil {
		return createErrorPanel(g, output)
	}
	return nil
}

func getSelectedBranch(v *gocui.View) Branch {
	lineNumber := getItemPosition(v)
	return state.Branches[lineNumber]
}

func renderBranchesOptions(g *gocui.Gui) error {
	return renderOptionsMap(g, map[string]string{
		"space":   "checkout",
		"f":       "force checkout",
		"m":       "merge",
		"c":       "checkout by name",
		"n":       "new branch",
		"← → ↑ ↓": "navigate",
	})
}

// may want to standardise how these select methods work
func handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if err := renderBranchesOptions(g); err != nil {
		return err
	}
	// This really shouldn't happen: there should always be a master branch
	if len(state.Branches) == 0 {
		return renderString(g, "main", "No branches for this repo")
	}
	go func() {
		branch := getSelectedBranch(v)
		diff, _ := getBranchGraph(branch.Name, branch.BaseBranch)
		renderString(g, "main", diff)
	}()
	return nil
}

// refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func refreshBranches(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("branches")
		if err != nil {
			panic(err)
		}
		state.Branches = getGitBranches()
		v.Clear()
		for _, branch := range state.Branches {
			fmt.Fprintln(v, branch.DisplayString)
		}
		resetOrigin(v)
		return refreshStatus(g)
	})
	return nil
}
