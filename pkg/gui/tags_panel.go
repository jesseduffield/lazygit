package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

// list panel functions

func (gui *Gui) getSelectedTag() *commands.Tag {
	selectedLine := gui.State.Panels.Tags.SelectedLine
	if selectedLine == -1 || len(gui.State.Tags) == 0 {
		return nil
	}

	return gui.State.Tags[selectedLine]
}

func (gui *Gui) handleTagSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Tag"

	tag := gui.getSelectedTag()
	if tag == nil {
		return gui.newStringTask("main", "No tags")
	}
	v.FocusPoint(0, gui.State.Panels.Tags.SelectedLine)

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.GetBranchGraphCmdStr(tag.Name),
	)
	if err := gui.newCmdTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) refreshTags() error {
	tags, err := gui.GitCommand.GetTags()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.Tags = tags

	if gui.getBranchesView().Context == "tags" {
		return gui.renderTagsWithSelection()
	}

	return nil
}

func (gui *Gui) renderTagsWithSelection() error {
	branchesView := gui.getBranchesView()

	gui.refreshSelectedLine(&gui.State.Panels.Tags.SelectedLine, len(gui.State.Tags))
	displayStrings := presentation.GetTagListDisplayStrings(gui.State.Tags)
	gui.renderDisplayStrings(branchesView, displayStrings)
	if gui.g.CurrentView() == branchesView && branchesView.Context == "tags" {
		if err := gui.handleTagSelect(gui.g, branchesView); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleCheckoutTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}
	if err := gui.handleCheckoutRef(tag.Name, handleCheckoutRefOptions{}); err != nil {
		return err
	}
	return gui.switchBranchesPanelContext("local-branches")
}

func (gui *Gui) handleDeleteTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	prompt := gui.Tr.TemplateLocalize(
		"DeleteTagPrompt",
		Teml{
			"tagName": tag.Name,
		},
	)

	return gui.createConfirmationPanel(gui.g, v, true, gui.Tr.SLocalize("DeleteTagTitle"), prompt, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.DeleteTag(tag.Name); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshTags(); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		return nil
	}, nil)
}

func (gui *Gui) handlePushTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	title := gui.Tr.TemplateLocalize(
		"PushTagTitle",
		Teml{
			"tagName": tag.Name,
		},
	)

	return gui.createPromptPanel(gui.g, v, title, "origin", func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.PushTag(v.Buffer(), tag.Name); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		return gui.refreshTags()
	})
}

func (gui *Gui) handleCreateTag(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel(gui.g, v, gui.Tr.SLocalize("CreateTagTitle"), "", func(g *gocui.Gui, v *gocui.View) error {
		// leaving commit SHA blank so that we're just creating the tag for the current commit
		if err := gui.GitCommand.CreateLightweightTag(v.Buffer(), ""); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshTags(); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		return nil
	})
}

func (gui *Gui) handleCreateResetToTagMenu(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	return gui.createResetMenu(tag.Name)
}
