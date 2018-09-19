package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/git"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	index := gui.getItemPosition(gui.getBranchesView(g))
	if index == 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("AlreadyCheckedOutBranch"))
	}
	branch := gui.getSelectedBranch(gui.getBranchesView(g))
	if err := gui.GitCommand.Checkout(branch.Name, false); err != nil {
		gui.createErrorPanel(g, err.Error())
	}
	return gui.refreshSidePanels(g)
}

func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch(v)
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
	checkedOutBranch := gui.State.Branches[0]
	selectedBranch := gui.getSelectedBranch(v)
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantDeleteCheckOutBranch"))
	}
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
	return gui.createConfirmationPanel(g, v, title, message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.DeleteBranch(selectedBranch.Name, force); err != nil {
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
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantMergeBranchIntoItself"))
	}
	if err := gui.GitCommand.Merge(selectedBranch.Name); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	return nil
}

func (gui *Gui) getSelectedBranch(v *gocui.View) *commands.Branch {
	lineNumber := gui.getItemPosition(v)
	return gui.State.Branches[lineNumber]
}

func (gui *Gui) renderBranchesOptions(g *gocui.Gui) error {
	return gui.renderGlobalOptions(g)
}

// may want to standardise how these select methods work
func (gui *Gui) handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if err := gui.renderBranchesOptions(g); err != nil {
		return err
	}
	// This really shouldn't happen: there should always be a master branch
	if len(gui.State.Branches) == 0 {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoBranchesThisRepo"))
	}
	go func() {
		branch := gui.getSelectedBranch(v)
		diff, err := gui.GitCommand.GetBranchGraph(branch.Name)
		if err != nil && strings.HasPrefix(diff, "fatal: ambiguous argument") {
			diff = gui.Tr.SLocalize("NoTrackingThisBranch")
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
		list, err := utils.RenderList(gui.State.Branches)
		if err != nil {
			return err
		}

		fmt.Fprint(v, list)

		gui.resetOrigin(v)
		return gui.refreshStatus(g)
	})
	return nil
}
