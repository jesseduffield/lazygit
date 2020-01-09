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
	selectedLine := gui.State.Panels.Remotes.SelectedLine
	if selectedLine == -1 || len(gui.State.Remotes) == 0 {
		return nil
	}

	return gui.State.Remotes[selectedLine]
}

func (gui *Gui) handleRemoteSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Remote"

	remote := gui.getSelectedRemote()
	if remote == nil {
		return gui.renderString(g, "main", "No remotes")
	}
	if err := gui.focusPoint(0, gui.State.Panels.Remotes.SelectedLine, len(gui.State.Remotes), v); err != nil {
		return err
	}

	return gui.renderString(g, "main", fmt.Sprintf("%s\nUrls:\n%s", utils.ColoredString(remote.Name, color.FgGreen), strings.Join(remote.Urls, "\n")))
}

func (gui *Gui) refreshRemotes() error {
	prevSelectedRemote := gui.getSelectedRemote()

	remotes, err := gui.GitCommand.GetRemotes()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
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
	if err := gui.renderListPanel(branchesView, gui.State.Remotes); err != nil {
		return err
	}
	if gui.g.CurrentView() == branchesView && branchesView.Context == "remotes" {
		if err := gui.handleRemoteSelect(gui.g, branchesView); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleRemoteEnter(g *gocui.Gui, v *gocui.View) error {
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
	return gui.createPromptPanel(g, branchesView, gui.Tr.SLocalize("newRemoteName"), "", func(g *gocui.Gui, v *gocui.View) error {
		remoteName := gui.trimmedContent(v)
		return gui.createPromptPanel(g, branchesView, gui.Tr.SLocalize("newRemoteUrl"), "", func(g *gocui.Gui, v *gocui.View) error {
			remoteUrl := gui.trimmedContent(v)
			if err := gui.GitCommand.AddRemote(remoteName, remoteUrl); err != nil {
				return err
			}
			return gui.refreshRemotes()
		})
	})
}

func (gui *Gui) handleRemoveRemote(g *gocui.Gui, v *gocui.View) error {
	remote := gui.getSelectedRemote()
	if remote == nil {
		return nil
	}
	return gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("removeRemote"), gui.Tr.SLocalize("removeRemotePrompt")+" '"+remote.Name+"'?", func(*gocui.Gui, *gocui.View) error {
		if err := gui.GitCommand.RemoveRemote(remote.Name); err != nil {
			return err
		}

		return gui.refreshRemotes()

	}, nil)
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

	return gui.createPromptPanel(g, branchesView, editNameMessage, "", func(g *gocui.Gui, v *gocui.View) error {
		updatedRemoteName := gui.trimmedContent(v)

		if updatedRemoteName != remote.Name {
			if err := gui.GitCommand.RenameRemote(remote.Name, updatedRemoteName); err != nil {
				return gui.createErrorPanel(gui.g, err.Error())
			}
		}

		editUrlMessage := gui.Tr.TemplateLocalize(
			"editRemoteUrl",
			Teml{
				"remoteName": updatedRemoteName,
			},
		)

		return gui.createPromptPanel(g, branchesView, editUrlMessage, "", func(g *gocui.Gui, v *gocui.View) error {
			updatedRemoteUrl := gui.trimmedContent(v)
			if err := gui.GitCommand.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
				return gui.createErrorPanel(gui.g, err.Error())
			}
			return gui.refreshRemotes()
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

		return gui.refreshRemotes()
	})
}
