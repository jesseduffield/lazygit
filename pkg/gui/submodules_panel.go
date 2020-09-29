package gui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getSelectedSubmodule() *models.SubmoduleConfig {
	selectedLine := gui.State.Panels.Submodules.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.Submodules) == 0 {
		return nil
	}

	return gui.State.Submodules[selectedLine]
}

func (gui *Gui) handleSubmoduleSelect() error {
	var task updateTask
	submodule := gui.getSelectedSubmodule()
	if submodule == nil {
		task = gui.createRenderStringTask("No submodules")
	} else {
		// TODO: we want to display the path, name, url, and a diff. We really need to be able to pipe commands together. We can always pipe commands together and just not do it asynchronously, but what if it's an expensive diff to obtain? I think that makes the most sense now though.
		task = gui.createRenderStringTask(
			fmt.Sprintf(
				"Name: %s\nPath: %s\nUrl:  %s\n",
				utils.ColoredString(submodule.Name, color.FgGreen),
				utils.ColoredString(submodule.Path, color.FgYellow),
				utils.ColoredString(submodule.Url, color.FgCyan),
			),
		)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Submodule",
			task:  task,
		},
	})
}

func (gui *Gui) handleSubmoduleEnter() error {
	submodule := gui.getSelectedSubmodule()
	if submodule == nil {
		return nil
	}

	return gui.enterSubmodule(submodule)
}

func (gui *Gui) enterSubmodule(submodule *models.SubmoduleConfig) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gui.State.RepoPathStack = append(gui.State.RepoPathStack, wd)

	return gui.dispatchSwitchToRepo(submodule.Path)
}

// func (gui *Gui) handleAddRemote(g *gocui.Gui, v *gocui.View) error {
// 	return gui.prompt(gui.Tr.SLocalize("newRemoteName"), "", func(remoteName string) error {
// 		return gui.prompt(gui.Tr.SLocalize("newRemoteUrl"), "", func(remoteUrl string) error {
// 			if err := gui.GitCommand.AddRemote(remoteName, remoteUrl); err != nil {
// 				return err
// 			}
// 			return gui.refreshSidePanels(refreshOptions{scope: []int{REMOTES}})
// 		})
// 	})
// }

// func (gui *Gui) handleRemoveRemote(g *gocui.Gui, v *gocui.View) error {
// 	remote := gui.getSelectedSubmodule()
// 	if remote == nil {
// 		return nil
// 	}

// 	return gui.ask(askOpts{
// 		title:  gui.Tr.SLocalize("removeRemote"),
// 		prompt: gui.Tr.SLocalize("removeRemotePrompt") + " '" + remote.Name + "'?",
// 		handleConfirm: func() error {
// 			if err := gui.GitCommand.RemoveRemote(remote.Name); err != nil {
// 				return err
// 			}

// 			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
// 		},
// 	})
// }

// func (gui *Gui) handleEditRemote(g *gocui.Gui, v *gocui.View) error {
// 	remote := gui.getSelectedSubmodule()
// 	if remote == nil {
// 		return nil
// 	}

// 	editNameMessage := gui.Tr.TemplateLocalize(
// 		"editRemoteName",
// 		Teml{
// 			"remoteName": remote.Name,
// 		},
// 	)

// 	return gui.prompt(editNameMessage, remote.Name, func(updatedRemoteName string) error {
// 		if updatedRemoteName != remote.Name {
// 			if err := gui.GitCommand.RenameRemote(remote.Name, updatedRemoteName); err != nil {
// 				return gui.surfaceError(err)
// 			}
// 		}

// 		editUrlMessage := gui.Tr.TemplateLocalize(
// 			"editRemoteUrl",
// 			Teml{
// 				"remoteName": updatedRemoteName,
// 			},
// 		)

// 		urls := remote.Urls
// 		url := ""
// 		if len(urls) > 0 {
// 			url = urls[0]
// 		}

// 		return gui.prompt(editUrlMessage, url, func(updatedRemoteUrl string) error {
// 			if err := gui.GitCommand.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
// 				return gui.surfaceError(err)
// 			}
// 			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
// 		})
// 	})
// }
