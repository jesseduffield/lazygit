package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemote() *commands.Remote {
	selectedLine := gui.State.Panels.Remotes.SelectedLine
	if selectedLine == -1 || len(gui.State.Remotes) == 0 {
		return nil
	}

	return gui.State.Remotes[selectedLine]
}

func (gui *Gui) handleRemoteSelect() error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	gui.getMainView().Title = "Remote"

	remote := gui.getSelectedRemote()
	if remote == nil {
		return gui.newStringTask("main", "No remotes")
	}
	if gui.inDiffMode() {
		return gui.renderDiff()
	}

	return gui.newStringTask("main", fmt.Sprintf("%s\nUrls:\n%s", utils.ColoredString(remote.Name, color.FgGreen), strings.Join(remote.Urls, "\n")))
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

	// TODO: see if this works for deleting remote branches
	switch gui.getBranchesView().Context {
	case "remotes":
		return gui.renderRemotesWithSelection()
	case "remote-branches":
		return gui.renderRemoteBranchesWithSelection()
	}

	return nil
}

func (gui *Gui) renderRemotesWithSelection() error {
	branchesView := gui.getBranchesView()

	gui.refreshSelectedLine(&gui.State.Panels.Remotes.SelectedLine, len(gui.State.Remotes))

	displayStrings := presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Diff.Ref)
	gui.renderDisplayStrings(branchesView, displayStrings)

	if gui.g.CurrentView() == branchesView && branchesView.Context == "remotes" {
		if err := gui.handleRemoteSelect(); err != nil {
			return err
		}
	}

	return nil
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
	gui.State.Panels.RemoteBranches.SelectedLine = newSelectedLine

	return gui.switchBranchesPanelContext("remote-branches")
}

func (gui *Gui) handleAddRemote(g *gocui.Gui, v *gocui.View) error {
	branchesView := gui.getBranchesView()
	return gui.prompt(branchesView, gui.Tr.SLocalize("newRemoteName"), "", func(remoteName string) error {
		return gui.prompt(branchesView, gui.Tr.SLocalize("newRemoteUrl"), "", func(remoteUrl string) error {
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
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("removeRemote"),
		prompt:             gui.Tr.SLocalize("removeRemotePrompt") + " '" + remote.Name + "'?",
		handleConfirm: func() error {
			if err := gui.GitCommand.RemoveRemote(remote.Name); err != nil {
				return err
			}

			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
		},
	})
}

func (gui *Gui) handleEditRemote(g *gocui.Gui, v *gocui.View) error {
	branchesView := gui.getBranchesView()
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

	return gui.prompt(branchesView, editNameMessage, "", func(updatedRemoteName string) error {
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

		return gui.prompt(branchesView, editUrlMessage, "", func(updatedRemoteUrl string) error {
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
		if err := gui.GitCommand.FetchRemote(remote.Name); err != nil {
			return err
		}

		return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
	})
}
