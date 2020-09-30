package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func (gui *Gui) handleRemoveSubmodule() error {
	submodule := gui.getSelectedSubmodule()
	if submodule == nil {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("RemoveSubmodule"),
		prompt: gui.Tr.SLocalizef("RemoveSubmodulePrompt", submodule.Name),
		handleConfirm: func() error {
			if err := gui.GitCommand.SubmoduleDelete(submodule); err != nil {
				return gui.surfaceError(err)
			}

			return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES, FILES}})
		},
	})
}

func (gui *Gui) handleResetSubmodule() error {
	return gui.WithWaitingStatus(gui.Tr.SLocalize("resettingSubmoduleStatus"), func() error {
		submodule := gui.getSelectedSubmodule()
		if submodule == nil {
			return nil
		}

		return gui.resetSubmodule(submodule)
	})
}

func (gui *Gui) fileForSubmodule(submodule *models.SubmoduleConfig) *models.File {
	for _, file := range gui.State.Files {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}

func (gui *Gui) resetSubmodule(submodule *models.SubmoduleConfig) error {
	file := gui.fileForSubmodule(submodule)
	if file != nil {
		if err := gui.GitCommand.UnStageFile(file.Name, file.Tracked); err != nil {
			return gui.surfaceError(err)
		}
	}

	if err := gui.GitCommand.SubmoduleStash(submodule); err != nil {
		return gui.surfaceError(err)
	}
	if err := gui.GitCommand.SubmoduleReset(submodule); err != nil {
		return gui.surfaceError(err)
	}

	return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES, SUBMODULES}})
}

func (gui *Gui) handleAddSubmodule() error {
	return gui.prompt(gui.Tr.SLocalize("newSubmoduleUrl"), "", func(submoduleUrl string) error {
		nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

		return gui.prompt(gui.Tr.SLocalize("newSubmoduleName"), nameSuggestion, func(submoduleName string) error {
			return gui.prompt(gui.Tr.SLocalize("newSubmodulePath"), submoduleName, func(submodulePath string) error {
				return gui.WithWaitingStatus(gui.Tr.SLocalize("addingSubmoduleStatus"), func() error {
					err := gui.GitCommand.AddSubmodule(submoduleName, submodulePath, submoduleUrl)
					gui.handleCredentialsPopup(err)

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})

				// go func() {
				// 	err := gui.GitCommand.AddSubmodule(submoduleName, submodulePath, submoduleUrl)
				// 	gui.handleCredentialsPopup(err)

				// 	_ = gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				// }()
				return nil
			})
		})
	})
}

// func (gui *Gui) handleEditsubmodule(g *gocui.Gui, v *gocui.View) error {
// 	submodule := gui.getSelectedSubmodule()
// 	if submodule == nil {
// 		return nil
// 	}

// 	editNameMessage := gui.Tr.TemplateLocalize(
// 		"editsubmoduleName",
// 		Teml{
// 			"submoduleName": submodule.Name,
// 		},
// 	)

// 	return gui.prompt(editNameMessage, submodule.Name, func(updatedsubmoduleName string) error {
// 		if updatedsubmoduleName != submodule.Name {
// 			if err := gui.GitCommand.Renamesubmodule(submodule.Name, updatedsubmoduleName); err != nil {
// 				return gui.surfaceError(err)
// 			}
// 		}

// 		editUrlMessage := gui.Tr.TemplateLocalize(
// 			"editsubmoduleUrl",
// 			Teml{
// 				"submoduleName": updatedsubmoduleName,
// 			},
// 		)

// 		urls := submodule.Urls
// 		url := ""
// 		if len(urls) > 0 {
// 			url = urls[0]
// 		}

// 		return gui.prompt(editUrlMessage, url, func(updatedsubmoduleUrl string) error {
// 			if err := gui.GitCommand.UpdatesubmoduleUrl(updatedsubmoduleName, updatedsubmoduleUrl); err != nil {
// 				return gui.surfaceError(err)
// 			}
// 			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, submoduleS}})
// 		})
// 	})
// }
