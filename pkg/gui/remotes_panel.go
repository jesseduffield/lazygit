package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemote() *models.Remote {
	selectedLine := gui.State.Panels.Remotes.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.Remotes) == 0 {
		return nil
	}

	return gui.State.Remotes[selectedLine]
}

func (gui *Gui) remotesRenderToMain() error {
	var task updateTask
	remote := gui.getSelectedRemote()
	if remote == nil {
		task = NewRenderStringTask("No remotes")
	} else {
		task = NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", style.FgGreen.Sprint(remote.Name), strings.Join(remote.Urls, "\n")))
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote",
			task:  task,
		},
	})
}

func (gui *Gui) refreshRemotes() error {
	prevSelectedRemote := gui.getSelectedRemote()

	remotes, err := gui.Git.Loaders.Remotes.GetRemotes()
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.State.Remotes = remotes

	// we need to ensure our selected remote branches aren't now outdated
	if prevSelectedRemote != nil && gui.State.RemoteBranches != nil {
		// find remote now
		for _, remote := range remotes {
			if remote.Name == prevSelectedRemote.Name {
				gui.State.RemoteBranches = remote.Branches
			}
		}
	}

	return gui.postRefreshUpdate(gui.mustContextForContextKey(ContextKey(gui.Views.Branches.Context)))
}

func (gui *Gui) handleRemoteEnter() error {
	// naive implementation: get the branches and render them to the list, change the context
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	gui.State.RemoteBranches = remote.Branches

	newSelectedLine := 0
	if len(remote.Branches) == 0 {
		newSelectedLine = -1
	}
	gui.State.Panels.RemoteBranches.SelectedLineIdx = newSelectedLine

	return gui.pushContext(gui.State.Contexts.RemoteBranches)
}

func (gui *Gui) handleAddRemote() error {
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title: gui.Tr.LcNewRemoteName,
		HandleConfirm: func(remoteName string) error {
			return gui.PopupHandler.Prompt(popup.PromptOpts{
				Title: gui.Tr.LcNewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					gui.logAction(gui.Tr.Actions.AddRemote)
					if err := gui.Git.Remote.AddRemote(remoteName, remoteUrl); err != nil {
						return err
					}
					return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.REMOTES}})
				},
			})
		},
	})

}

func (gui *Gui) handleRemoveRemote() error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.LcRemoveRemote,
		Prompt: gui.Tr.LcRemoveRemotePrompt + " '" + remote.Name + "'?",
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.RemoveRemote)
			if err := gui.Git.Remote.RemoveRemote(remote.Name); err != nil {
				return gui.PopupHandler.Error(err)
			}

			return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (gui *Gui) handleEditRemote() error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	editNameMessage := utils.ResolvePlaceholderString(
		gui.Tr.LcEditRemoteName,
		map[string]string{
			"remoteName": remote.Name,
		},
	)

	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:          editNameMessage,
		InitialContent: remote.Name,
		HandleConfirm: func(updatedRemoteName string) error {
			if updatedRemoteName != remote.Name {
				gui.logAction(gui.Tr.Actions.UpdateRemote)
				if err := gui.Git.Remote.RenameRemote(remote.Name, updatedRemoteName); err != nil {
					return gui.PopupHandler.Error(err)
				}
			}

			editUrlMessage := utils.ResolvePlaceholderString(
				gui.Tr.LcEditRemoteUrl,
				map[string]string{
					"remoteName": updatedRemoteName,
				},
			)

			urls := remote.Urls
			url := ""
			if len(urls) > 0 {
				url = urls[0]
			}

			return gui.PopupHandler.Prompt(popup.PromptOpts{
				Title:          editUrlMessage,
				InitialContent: url,
				HandleConfirm: func(updatedRemoteUrl string) error {
					gui.logAction(gui.Tr.Actions.UpdateRemote)
					if err := gui.Git.Remote.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
						return gui.PopupHandler.Error(err)
					}
					return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				},
			})
		},
	})
}

func (gui *Gui) handleFetchRemote() error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.FetchingRemoteStatus, func() error {
		gui.Mutexes.FetchMutex.Lock()
		defer gui.Mutexes.FetchMutex.Unlock()

		err := gui.Git.Sync.FetchRemote(remote.Name)
		if err != nil {
			_ = gui.PopupHandler.Error(err)
		}

		return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
	})
}
