package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemote() *commands.Remote {
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
		task = gui.createRenderStringTask("No remotes")
	} else {
		task = gui.createRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", utils.ColoredString(remote.Name, color.FgGreen), strings.Join(remote.Urls, "\n")))
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

	remotes, err := gui.GitCommand.GetRemotes()
	if err != nil {
		return gui.surfaceError(err)
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

	return gui.postRefreshUpdate(gui.contextForContextKey(gui.getBranchesView().Context))
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

	return gui.switchContext(gui.Contexts.Remotes.Branches.Context)
}

func (gui *Gui) handleAddRemote(g *gocui.Gui, v *gocui.View) error {
	return gui.prompt(gui.Tr.SLocalize("newRemoteName"), "", func(remoteName string) error {
		return gui.prompt(gui.Tr.SLocalize("newRemoteUrl"), "", func(remoteUrl string) error {
			if err := gui.GitCommand.AddRemote(remoteName, remoteUrl); err != nil {
				return err
			}
			return gui.refreshSidePanels(refreshOptions{scope: []int{REMOTES}})
		})
	})
}

func (gui *Gui) handleRemoveRemote(g *gocui.Gui, v *gocui.View) error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("removeRemote"),
		prompt: gui.Tr.SLocalize("removeRemotePrompt") + " '" + remote.Name + "'?",
		handleConfirm: func() error {
			if err := gui.GitCommand.RemoveRemote(remote.Name); err != nil {
				return err
			}

			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
		},
	})
}

func (gui *Gui) handleEditRemote(g *gocui.Gui, v *gocui.View) error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	editNameMessage := gui.Tr.TemplateLocalize(
		"editRemoteName",
		Teml{
			"remoteName": remote.Name,
		},
	)

	return gui.prompt(editNameMessage, remote.Name, func(updatedRemoteName string) error {
		if updatedRemoteName != remote.Name {
			if err := gui.GitCommand.RenameRemote(remote.Name, updatedRemoteName); err != nil {
				return gui.surfaceError(err)
			}
		}

		editUrlMessage := gui.Tr.TemplateLocalize(
			"editRemoteUrl",
			Teml{
				"remoteName": updatedRemoteName,
			},
		)

		urls := remote.Urls
		url := ""
		if len(urls) > 0 {
			url = urls[0]
		}

		return gui.prompt(editUrlMessage, url, func(updatedRemoteUrl string) error {
			if err := gui.GitCommand.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
		})
	})
}

func (gui *Gui) handleFetchRemote(g *gocui.Gui, v *gocui.View) error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("FetchingRemoteStatus"), func() error {
		gui.State.FetchMutex.Lock()
		defer gui.State.FetchMutex.Unlock()

		if err := gui.GitCommand.FetchRemote(remote.Name); err != nil {
			return err
		}

		return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
	})
}
