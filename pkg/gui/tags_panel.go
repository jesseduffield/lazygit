package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedTag() *models.Tag {
	selectedLine := gui.State.Panels.Tags.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.Tags) == 0 {
		return nil
	}

	return gui.State.Tags[selectedLine]
}

func (gui *Gui) handleTagSelect() error {
	var task updateTask
	tag := gui.getSelectedTag()
	if tag == nil {
		task = gui.createRenderStringTask("No tags")
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.GetBranchGraphCmdStr(tag.Name),
		)
		task = gui.createRunCommandTask(cmd)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Tag",
			task:  task,
		},
	})
}

func (gui *Gui) refreshTags() error {
	tags, err := gui.GitCommand.GetTags()
	if err != nil {
		return gui.surfaceError(err)
	}

	gui.State.Tags = tags

	return gui.postRefreshUpdate(gui.Contexts.Tags.Context)
}

func (gui *Gui) handleCheckoutTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}
	if err := gui.handleCheckoutRef(tag.Name, handleCheckoutRefOptions{}); err != nil {
		return err
	}
	return gui.switchContext(gui.Contexts.Branches.Context)
}

func (gui *Gui) handleDeleteTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	prompt := utils.ResolvePlaceholderString(
		gui.Tr.DeleteTagPrompt,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.DeleteTagTitle,
		prompt: prompt,
		handleConfirm: func() error {
			if err := gui.GitCommand.DeleteTag(tag.Name); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{COMMITS, TAGS}})
		},
	})
}

func (gui *Gui) handlePushTag(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	title := utils.ResolvePlaceholderString(
		gui.Tr.PushTagTitle,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return gui.prompt(title, "origin", func(response string) error {
		return gui.WithWaitingStatus(gui.Tr.PushingTagStatus, func() error {
			err := gui.GitCommand.PushTag(response, tag.Name, gui.promptUserForCredential)
			gui.handleCredentialsPopup(err)

			return nil
		})
	})
}

func (gui *Gui) handleCreateTag(g *gocui.Gui, v *gocui.View) error {
	return gui.prompt(gui.Tr.CreateTagTitle, "", func(tagName string) error {
		// leaving commit SHA blank so that we're just creating the tag for the current commit
		if err := gui.GitCommand.CreateLightweightTag(tagName, ""); err != nil {
			return gui.surfaceError(err)
		}
		return gui.refreshSidePanels(refreshOptions{scope: []int{COMMITS, TAGS}, then: func() {
			// find the index of the tag and set that as the currently selected line
			for i, tag := range gui.State.Tags {
				if tag.Name == tagName {
					gui.State.Panels.Tags.SelectedLineIdx = i
					if err := gui.Contexts.Tags.Context.HandleRender(); err != nil {
						gui.Log.Error(err)
					}

					return
				}
			}
		},
		})
	})
}

func (gui *Gui) handleCreateResetToTagMenu(g *gocui.Gui, v *gocui.View) error {
	tag := gui.getSelectedTag()
	if tag == nil {
		return nil
	}

	return gui.createResetMenu(tag.Name)
}
