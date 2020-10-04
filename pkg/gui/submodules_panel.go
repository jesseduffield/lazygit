package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
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
		prefix := fmt.Sprintf(
			"Name: %s\nPath: %s\nUrl:  %s\n\n",
			utils.ColoredString(submodule.Name, color.FgGreen),
			utils.ColoredString(submodule.Path, color.FgYellow),
			utils.ColoredString(submodule.Url, color.FgCyan),
		)

		file := gui.fileForSubmodule(submodule)
		if file == nil {
			task = gui.createRenderStringTask(prefix)
		} else {
			cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(file, false, !file.HasUnstagedChanges && file.HasStagedChanges)
			cmd := gui.OSCommand.ExecutableFromString(cmdStr)
			task = gui.createRunCommandTaskWithPrefix(cmd, prefix)
		}
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Submodule",
			task:  task,
		},
	})
}

func (gui *Gui) refreshStateSubmoduleConfigs() error {
	configs, err := gui.GitCommand.GetSubmoduleConfigs()
	if err != nil {
		return err
	}

	gui.State.Submodules = configs

	return nil
}

func (gui *Gui) handleSubmoduleEnter(submodule *models.SubmoduleConfig) error {
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

func (gui *Gui) removeSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.ask(askOpts{
		title:  gui.Tr.RemoveSubmodule,
		prompt: fmt.Sprintf(gui.Tr.RemoveSubmodulePrompt, submodule.Name),
		handleConfirm: func() error {
			if err := gui.GitCommand.SubmoduleDelete(submodule); err != nil {
				return gui.surfaceError(err)
			}

			return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES, FILES}})
		},
	})
}

func (gui *Gui) handleResetSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcResettingSubmoduleStatus, func() error {
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
	return gui.prompt(gui.Tr.LcNewSubmoduleUrl, "", func(submoduleUrl string) error {
		nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

		return gui.prompt(gui.Tr.LcNewSubmoduleName, nameSuggestion, func(submoduleName string) error {
			return gui.prompt(gui.Tr.LcNewSubmodulePath, submoduleName, func(submodulePath string) error {
				return gui.WithWaitingStatus(gui.Tr.LcAddingSubmoduleStatus, func() error {
					err := gui.GitCommand.SubmoduleAdd(submoduleName, submodulePath, submoduleUrl)
					gui.handleCredentialsPopup(err)

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})
			})
		})
	})
}

func (gui *Gui) handleEditSubmoduleUrl(submodule *models.SubmoduleConfig) error {
	return gui.prompt(fmt.Sprintf(gui.Tr.LcUpdateSubmoduleUrl, submodule.Name), submodule.Url, func(newUrl string) error {
		return gui.WithWaitingStatus(gui.Tr.LcUpdatingSubmoduleUrlStatus, func() error {
			err := gui.GitCommand.SubmoduleUpdateUrl(submodule.Name, submodule.Path, newUrl)
			gui.handleCredentialsPopup(err)

			return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
		})
	})
}

func (gui *Gui) handleSubmoduleInit(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcInitializingSubmoduleStatus, func() error {
		err := gui.GitCommand.SubmoduleInit(submodule.Path)
		gui.handleCredentialsPopup(err)

		return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
	})
}

func (gui *Gui) forSubmodule(callback func(*models.SubmoduleConfig) error) func(g *gocui.Gui, v *gocui.View) error {
	return gui.wrappedHandler(
		func() error {
			submodule := gui.getSelectedSubmodule()
			if submodule == nil {
				return nil
			}

			return callback(submodule)
		},
	)
}

func (gui *Gui) handleResetRemoveSubmodule(submodule *models.SubmoduleConfig) error {
	menuItems := []*menuItem{
		{
			displayString: gui.Tr.LcSubmoduleStashAndReset,
			onPress: func() error {
				return gui.resetSubmodule(submodule)
			},
		},
		{
			displayString: gui.Tr.LcRemoveSubmodule,
			onPress: func() error {
				return gui.removeSubmodule(submodule)
			},
		},
	}

	return gui.createMenu(submodule.Name, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleBulkSubmoduleActionsMenu() error {
	menuItems := []*menuItem{
		{
			displayStrings: []string{gui.Tr.LcBulkInitSubmodules, utils.ColoredString(gui.GitCommand.SubmoduleBulkInitCmdStr(), color.FgGreen)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.RunCommand(gui.GitCommand.SubmoduleBulkInitCmdStr()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcBulkUpdateSubmodules, utils.ColoredString(gui.GitCommand.SubmoduleBulkUpdateCmdStr(), color.FgYellow)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.RunCommand(gui.GitCommand.SubmoduleBulkUpdateCmdStr()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcSubmoduleStashAndReset, utils.ColoredString(fmt.Sprintf("git stash in each submodule && %s", gui.GitCommand.SubmoduleForceBulkUpdateCmdStr()), color.FgRed)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.GitCommand.ResetSubmodules(gui.State.Submodules); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcBulkDeinitSubmodules, utils.ColoredString(gui.GitCommand.SubmoduleBulkDeinitCmdStr(), color.FgRed)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.RunCommand(gui.GitCommand.SubmoduleBulkDeinitCmdStr()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
				})
			},
		},
	}

	return gui.createMenu(gui.Tr.LcBulkSubmoduleOptions, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleUpdateSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcUpdatingSubmoduleStatus, func() error {
		err := gui.GitCommand.SubmoduleUpdate(submodule.Path)
		gui.handleCredentialsPopup(err)

		return gui.refreshSidePanels(refreshOptions{scope: []int{SUBMODULES}})
	})
}
