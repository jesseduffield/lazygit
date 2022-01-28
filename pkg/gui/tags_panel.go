package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
	// leaving commit SHA blank so that we're just creating the tag for the current commit
	return gui.createTagMenu("")
}

func (gui *Gui) tagsRenderToMain() error {
	var task updateTask
	tag := gui.getSelectedTag()
	if tag == nil {
		task = NewRenderStringTask("No tags")
	} else {
		cmdObj := gui.Git.Branch.GetGraphCmdObj(tag.Name)
		task = NewRunCommandTask(cmdObj.GetCmd())
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
	tags, err := gui.Git.Loaders.Tags.GetTags()
	if err != nil {
		return gui.PopupHandler.Error(err)
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

// tag-specific handlers

func (gui *Gui) handleCheckoutTag(tag *models.Tag) error {
	gui.logAction(gui.Tr.Actions.CheckoutTag)
	if err := gui.handleCheckoutRef(tag.Name, handleCheckoutRefOptions{}); err != nil {
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

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.DeleteTagTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.DeleteTag)
			if err := gui.Git.Tag.Delete(tag.Name); err != nil {
				return gui.PopupHandler.Error(err)
			}
			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
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

	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: gui.getRemoteSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.PushingTagStatus, func() error {
				gui.logAction(gui.Tr.Actions.PushTag)
				err := gui.Git.Tag.Push(response, tag.Name)
				if err != nil {
					_ = gui.PopupHandler.Error(err)
				}

				return nil
			})
		},
	})
}

func (gui *Gui) handleCreateResetToTagMenu(tag *models.Tag) error {
	return gui.createResetMenu(tag.Name)
}
