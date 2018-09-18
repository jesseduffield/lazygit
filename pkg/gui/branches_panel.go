package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/git"
)

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches() error {
	gui.g.Update(func(g *gocui.Gui) error {

		v, err := gui.g.View("branches")
		if err != nil {
			return err
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

		if err := gui.resetOrigin(v); err != nil {
			return err
		}

		return gui.refreshStatus()
	})

	return nil
}

// handleBranchPress is called when the user selects a branch.
// g and v are passed by the gocui library, but are not used.
// In case something goes wrong it returns an error
func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	branchesView, err := gui.g.View("branches")
	if err != nil {
		return err
	}

	index := gui.getItemPosition(branchesView)
	if index == 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("AlreadyCheckedOutBranch"))
	}

	branch := gui.getSelectedBranch(branchesView)

	if err := gui.GitCommand.Checkout(branch.Name, false); err != nil {
		return gui.createErrorPanel(err.Error())
	}

	return gui.refresh()
}

// handleForceCheckout is called when the user wants to force checkout a branch
// g and v are passed by the gocui library, but are not used.
// In case something goes wrong it returns an error
func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch(v)
	message := gui.Tr.SLocalize("SureForceCheckout")
	title := gui.Tr.SLocalize("ForceCheckoutBranch")

	return gui.createConfirmationPanel(v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.Checkout(branch.Name, true); err != nil {
			return gui.createErrorPanel(err.Error())
		}

		return gui.refresh()
	}, nil)
}

// handleCheckoutByName gets called when a user presses the key
// to checkout a branch by name.
// g and v are passed by the gocui library, but only v is used.
// If something goes wrong it returns an error
func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel(v, gui.Tr.SLocalize("BranchName")+":",
		func(g *gocui.Gui, v *gocui.View) error {
			if err := gui.GitCommand.Checkout(gui.trimmedContent(v), false); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refresh()
		})
}

// handleNewBranch is called when a user wants to create a new branch.
// g and v are passed by the gocui library but only v is used.
// If something goes wrong it returns an error.
func (gui *Gui) handleNewBranch(g *gocui.Gui, v *gocui.View) error {
	branch := gui.State.Branches[0]
	message := gui.Tr.TemplateLocalize(
		"NewBranchNameBranchOff",
		Teml{
			"branchName": branch.Name,
		},
	)

	return gui.createPromptPanel(v, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.NewBranch(gui.trimmedContent(v)); err != nil {
			return gui.createErrorPanel(err.Error())
		}

		if err := gui.refresh(); err != nil {
			return err
		}

		return gui.handleBranchSelect(v)
	})
}

// handleDeleteBranch gets called when the user wants to normally delete a
// branch.
func (gui *Gui) handleDeleteBranch(g *gocui.Gui, v *gocui.View) error {
	return gui.deleteBranch(v, false)
}

// handleForceDeleteBranch gets called when the user wants to force delete a
// branch
func (gui *Gui) handleForceDeleteBranch(g *gocui.Gui, v *gocui.View) error {
	return gui.deleteBranch(v, true)
}

// deleteBranch gets called when the user wants to delete a branch.
// v is passed for ease sake, force is to indicate if it must be force
// deleted.
// If anything goes wrong, it returns an error.
func (gui *Gui) deleteBranch(v *gocui.View, force bool) error {
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	title := gui.Tr.SLocalize("DeleteBranch")

	var messageId string
	if force {
		messageId = "ForceDeleteBranchMessage"
	} else {
		messageId = "DeleteBranchMessage"
	}

	message := gui.Tr.TemplateLocalize(
		messageId,
		Teml{
			"selectedBranchName": selectedBranch.Name,
		},
	)

	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(gui.Tr.SLocalize("CantDeleteCheckOutBranch"))
	}

	return gui.createConfirmationPanel(v, title, message,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := gui.GitCommand.DeleteBranch(selectedBranch.Name, force); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refresh()
		}, nil)
}

// handleMerge is called when the user wants to merge.
// g and v are passed by the gocui library, but only v is used.
// If something goes wrong it returns an error
func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)

	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(gui.Tr.SLocalize("CantMergeBranchIntoItself"))
	}

	if err := gui.GitCommand.Merge(selectedBranch.Name); err != nil {
		return gui.createErrorPanel(err.Error())
	}

	return gui.refresh()
}

// getSelectedBranch returns the selected branch
func (gui *Gui) getSelectedBranch(v *gocui.View) commands.Branch {
	lineNumber := gui.getItemPosition(v)
	return gui.State.Branches[lineNumber]
}

// handleBranchSelect gets called when the user selects a branch
func (gui *Gui) handleBranchSelect(v *gocui.View) error {
	if err := gui.renderGlobalOptions(); err != nil {
		return err
	}

	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {
		return gui.renderString("main", gui.Tr.SLocalize("NoBranchesThisRepo"))
	}

	go func() {
		branch := gui.getSelectedBranch(v)

		diff, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(diff, "fatal: ambiguous argument") {
			diff = gui.Tr.SLocalize("NoTrackingThisBranch")
		}

		if err = gui.renderString("main", diff); err != nil {
			gui.Log.Errorf("Failed to render string at handleBranchSelect: %s\n", err)
		}
	}()
	return nil
}
