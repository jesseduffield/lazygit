package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
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

func (gui *Gui) handleRemoteSelect() error {
	var task updateTask
	remote := gui.getSelectedRemote()
	if remote == nil {
		task = NewRenderStringTask("No remotes")
	} else {
		task = NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", utils.ColoredString(remote.Name, color.FgGreen), strings.Join(remote.Urls, "\n")))
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

	remotes, err := gui.Git.Remotes().LoadRemotes()
	if err != nil {
		return gui.SurfaceError(err)
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
	return gui.Prompt(PromptOpts{
		Title: gui.Tr.LcNewRemoteName,
		HandleConfirm: func(remoteName string) error {
			return gui.Prompt(PromptOpts{
				Title: gui.Tr.LcNewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					if err := gui.Git.WithSpan(gui.Tr.Spans.AddRemote).Remotes().Add(remoteName, remoteUrl); err != nil {
						return err
					}
					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{REMOTES}})
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

	return gui.Ask(AskOpts{
		Title:  gui.Tr.LcRemoveRemote,
		Prompt: gui.Tr.LcRemoveRemotePrompt + " '" + remote.Name + "'?",
		HandleConfirm: func() error {
			if err := gui.Git.WithSpan(gui.Tr.Spans.RemoveRemote).Remotes().Remove(remote.Name); err != nil {
				return gui.SurfaceError(err)
			}

			return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{BRANCHES, REMOTES}})
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

	gitCommand := gui.Git.WithSpan(gui.Tr.Spans.UpdateRemote)

	return gui.Prompt(PromptOpts{
		Title:          editNameMessage,
		InitialContent: remote.Name,
		HandleConfirm: func(updatedRemoteName string) error {
			if updatedRemoteName != remote.Name {
				if err := gitCommand.Remotes().Rename(remote.Name, updatedRemoteName); err != nil {
					return gui.SurfaceError(err)
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

			return gui.Prompt(PromptOpts{
				Title:          editUrlMessage,
				InitialContent: url,
				HandleConfirm: func(updatedRemoteUrl string) error {
					if err := gitCommand.Remotes().UpdateUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
						return gui.SurfaceError(err)
					}
					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{BRANCHES, REMOTES}})
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

	return gui.WithWaitingStatus(gui.Tr.FetchingRemoteStatus, func() error {
		gui.Mutexes.FetchMutex.Lock()
		defer gui.Mutexes.FetchMutex.Unlock()

		err := gui.Git.Sync().FetchRemote(remote.Name)
		if err != nil {
			return gui.SurfaceError(err)
		}

		return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{BRANCHES, REMOTES}})
	})
}
