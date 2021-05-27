package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getSelectedTag() *models.Tag {
	selectedLine := gui.State.Panels.Tags.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.Tags) == 0 {
		return nil
	}

	return gui.State.Tags[selectedLine]
}

func (gui *Gui) handleCreateTag() error {
	return gui.prompt(promptOpts{
		title: gui.Tr.CreateTagTitle,
		handleConfirm: func(tagName string) error {
			// leaving commit SHA blank so that we're just creating the tag for the current commit
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.CreateLightweightTag).CreateLightweightTag(tagName, ""); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{COMMITS, TAGS}, then: func() {
				// find the index of the tag and set that as the currently selected line
				for i, tag := range gui.State.Tags {
					if tag.Name == tagName {
						gui.State.Panels.Tags.SelectedLineIdx = i
						if err := gui.State.Contexts.Tags.HandleRender(); err != nil {
							gui.Log.Error(err)
						}

						return
					}
				}
			},
			})
		},
	})
}

// tag-specific handlers
// view model would need to raise an event called 'tag selected', perhaps containing a tag. The listener would _be_ the main view, or the main context, and it would be able to render to itself.
func (gui *Gui) handleTagSelect() error {
	var task updateTask
	tag := gui.getSelectedTag()
	if tag == nil {
		task = NewRenderStringTask("No tags")
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.GetBranchGraphCmdStr(tag.Name),
		)
		task = NewRunCommandTask(cmd)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Tag",
			task:  task,
		},
	})
}

// this is a controller: it can't access tags directly. Or can it? It should be able to get but not set. But that's exactly what I'm doing here, setting it. but through a mutator which encapsulates the event.
func (gui *Gui) refreshTags() error {
	tags, err := gui.GitCommand.GetTags()
	if err != nil {
		return gui.surfaceError(err)
	}

	gui.State.Tags = tags

	return gui.postRefreshUpdate(gui.State.Contexts.Tags)
}

func (gui *Gui) withSelectedTag(f func(tag *models.Tag) error) func() error {
	return func() error {
		tag := gui.getSelectedTag()
		if tag == nil {
			return nil
		}

		return f(tag)
	}
}

func (gui *Gui) handleCheckoutTag(tag *models.Tag) error {
	if err := gui.handleCheckoutRef(tag.Name, handleCheckoutRefOptions{span: gui.Tr.Spans.CheckoutTag}); err != nil {
		return err
	}
	return gui.pushContext(gui.State.Contexts.Branches)
}

func (gui *Gui) handleDeleteTag(tag *models.Tag) error {
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
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.DeleteTag).DeleteTag(tag.Name); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{COMMITS, TAGS}})
		},
	})
}

func (gui *Gui) handlePushTag(tag *models.Tag) error {
	title := utils.ResolvePlaceholderString(
		gui.Tr.PushTagTitle,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return gui.prompt(promptOpts{
		title:          title,
		initialContent: "origin",
		handleConfirm: func(response string) error {
			return gui.WithWaitingStatus(gui.Tr.PushingTagStatus, func() error {
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.PushTag).PushTag(response, tag.Name, gui.promptUserForCredential)
				gui.handleCredentialsPopup(err)

				return nil
			})
		},
	})
}

func (gui *Gui) handleCreateResetToTagMenu(tag *models.Tag) error {
	return gui.createResetMenu(tag.Name)
}
