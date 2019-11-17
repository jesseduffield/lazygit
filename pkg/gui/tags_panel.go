package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
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
		return gui.renderString(g, "main", "No tags")
	}
	if err := gui.focusPoint(0, gui.State.Panels.Tags.SelectedLine, len(gui.State.Tags), v); err != nil {
		return err
	}

	go func() {
		show, err := gui.GitCommand.ShowTag(tag.Name)
		if err != nil {
			show = ""
		}

		graph, err := gui.GitCommand.GetBranchGraph(tag.Name)
		if err != nil {
			graph = "No graph for tag " + tag.Name
		}

		_ = gui.renderString(g, "main", fmt.Sprintf("%s\n%s", show, graph))
	}()

	return nil
}

func (gui *Gui) refreshTags() error {
	tags, err := gui.GitCommand.GetTags()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.Tags = tags

	if gui.getBranchesView().Context == "tags" {
		gui.renderTagsWithSelection()
	}

	return nil
}

func (gui *Gui) renderTagsWithSelection() error {
	branchesView := gui.getBranchesView()

	gui.refreshSelectedLine(&gui.State.Panels.Tags.SelectedLine, len(gui.State.Tags))
	if err := gui.renderListPanel(branchesView, gui.State.Tags); err != nil {
		return err
	}
	if err := gui.handleTagSelect(gui.g, branchesView); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleCheckoutTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}
	if err := gui.handleCheckoutBranch(tag.Name); err != nil {
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
		return gui.refreshTags()
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
		return gui.refreshTags()
	})
}
