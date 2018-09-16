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
			gui.Log.Errorf("Failed to get branches view at refreshbranches %s\n", err)
			return err
		}

		builder, err := git.NewBranchListBuilder(gui.Log, gui.GitCommand)
		if err != nil {
			gui.Log.Errorf("Failed to create branchbuilder at refreshBranches: %s\n", err)
			return err
		}

		gui.State.Branches = builder.Build()

		v.Clear()

		for _, branch := range gui.State.Branches {
			fmt.Fprintln(v, branch.GetDisplayString())
		}

		err = gui.resetOrigin(v)
		if err != nil {
			gui.Log.Errorf("Failed to reset origin at refreshBranches: %s\n", err)
			return err
		}

		err = gui.refreshStatus()
		if err != nil {
			gui.Log.Errorf("Failed to refresh statsu at refreshBranches: %s\n", err)
			return err
		}

		return nil
	})

	return nil
}

// handleBranchPress is called when the user selects a branch.
// g and v are passed by the gocui library, but are not used.
// In case something goes wrong it returns an error
func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {

	v, err := gui.g.View("branches")
	if err != nil {
		gui.Log.Errorf("Failed to get branch view at handleBranchPress: %s\n", err.Error())
	}

	index := gui.getItemPosition(v)
	if index == 0 {

		err := gui.createErrorPanel(gui.Tr.SLocalize("AlreadyCheckedOutBranch"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleBranchPress: %s\n", err)
			return err
		}

		return nil
	}

	branch := gui.getSelectedBranch(v)

	err = gui.GitCommand.Checkout(branch.Name, false)
	if err != nil {
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleBranchPress: %s\n", err)
			return err
		}
	}

	gui.refresh()

	return nil
}

// handleForceCheckout is called when the user wants to force checkout a branch
// g and v are passed by the gocui library, but are not used.
// In case something goes wrong it returns an error
func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {

	branch := gui.getSelectedBranch(v)
	message := gui.Tr.SLocalize("SureForceCheckout")
	title := gui.Tr.SLocalize("ForceCheckoutBranch")

	err := gui.createConfirmationPanel(v, title, message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.Checkout(branch.Name, true)
			if err != nil {
				err = gui.createErrorPanel(err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleForceCheckout: %s\n", err)
					return err
				}
			}

			gui.refresh()

			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to create confirmation panel at handleForceCheckout: %s\n", err)
		return err
	}

	return nil
}

// handleCheckoutByName gets called when a user presses the key
// to checkout a branch by name.
// g and v are passed by the gocui library, but only v is used.
// If something goes wrong it returns an error
func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {

	err := gui.createPromptPanel(v, gui.Tr.SLocalize("BranchName")+":",
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.Checkout(gui.trimmedContent(v), false)
			if err != nil {
				err = gui.createErrorPanel(err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleCheckoutByName: %s\n", err)
					return err
				}

				return nil
			}

			gui.refresh()

			return nil
		})
	if err != nil {
		gui.Log.Errorf("Failed to create prompt panel at handleCheckoutByName: %s\n", err)
		return err
	}

	return nil
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

	err := gui.createPromptPanel(v, message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.NewBranch(gui.trimmedContent(v))
			if err != nil {

				err = gui.createErrorPanel(err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleNewBranch: %s\n", err)
					return err
				}

				return nil
			}

			gui.refresh()

			err = gui.handleBranchSelect(v)
			if err != nil {
				gui.Log.Errorf("Failed to handleBranchSelect at handleNewBranch: %s\n", err)
				return err
			}

			return nil
		})
	if err != nil {
		gui.Log.Errorf("Failed to create prompt panel at handleNewBranch: %s\n", err)
		return err
	}

	return nil
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

	var messageId string

	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	title := gui.Tr.SLocalize("DeleteBranch")

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
		err := gui.createErrorPanel(gui.Tr.SLocalize("CantDeleteCheckOutBranch"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at deleteBranch: %s\n", err)
			return err
		}

		return nil
	}

	err := gui.createConfirmationPanel(v, title, message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.DeleteBranch(selectedBranch.Name, force)
			if err != nil {

				err = gui.createErrorPanel(err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at deleteBranch: %s\n", err)
					return err
				}

				return nil
			}

			gui.refresh()

			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to create confirmation panel at checkedOutBranch: %s\n", err)
		return err
	}

	return nil
}

// handleMerge is called when the user wants to merge.
// g and v are passed by the gocui library, but only v is used.
// If something goes wrong it returns an error
func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)

	defer gui.refresh()

	if checkedOutBranch.Name == selectedBranch.Name {
		err := gui.createErrorPanel(gui.Tr.SLocalize("CantMergeBranchIntoItself"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleMerge: %s\n", err)
			return err
		}

		return nil
	}

	err := gui.GitCommand.Merge(selectedBranch.Name)
	if err != nil {
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to error confirmation panel at handleMerge: %s\n", err)
			return err
		}

		return nil
	}

	return nil
}

// getSelectedBranch returns the selected branch
func (gui *Gui) getSelectedBranch(v *gocui.View) commands.Branch {
	lineNumber := gui.getItemPosition(v)
	return gui.State.Branches[lineNumber]
}

// handleBranchSelect gets called when the user selects a branch
func (gui *Gui) handleBranchSelect(v *gocui.View) error {

	err := gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Errorf("Failed to renderGlobalOptions at handleBranchSelect: %s\n", err)
		return err
	}

	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {

		err = gui.renderString("main", gui.Tr.SLocalize("NoBranchesThisRepo"))
		if err != nil {
			gui.Log.Errorf("Failed to render string at handleBranchSelect: %s\n", err)
			return err
		}

		return nil
	}

	go func() {
		branch := gui.getSelectedBranch(v)

		diff, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(diff, "fatal: ambiguous argument") {
			diff = gui.Tr.SLocalize("NoTrackingThisBranch")
		}

		err = gui.renderString("main", diff)
		if err != nil {
			gui.Log.Errorf("Failed to render string at handleBranchSelect: %s\n", err)
			return
		}
	}()
	return nil
}
