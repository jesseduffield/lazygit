package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
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
		task = NewRenderStringTask("No submodules")
	} else {
		prefix := fmt.Sprintf(
			"Name: %s\nPath: %s\nUrl:  %s\n\n",
			utils.ColoredString(submodule.Name, color.FgGreen),
			utils.ColoredString(submodule.Path, color.FgYellow),
			utils.ColoredString(submodule.Url, color.FgCyan),
		)

		file := gui.fileForSubmodule(submodule)
		if file == nil {
			task = NewRenderStringTask(prefix)
		} else {
			cmdObj := gui.GitCommand.WorktreeFileDiffCmdObj(file, false, !file.HasUnstagedChanges && file.HasStagedChanges)
			task = NewRunCommandTaskWithPrefix(cmdObj, prefix)
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
	gui.RepoPathStack = append(gui.RepoPathStack, wd)

	return gui.dispatchSwitchToRepo(submodule.Path, true)
}

func (gui *Gui) removeSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.Ask(AskOpts{
		Title:  gui.Tr.RemoveSubmodule,
		Prompt: fmt.Sprintf(gui.Tr.RemoveSubmodulePrompt, submodule.Name),
		HandleConfirm: func() error {
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.RemoveSubmodule).SubmoduleDelete(submodule); err != nil {
				return gui.SurfaceError(err)
			}

			return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES, FILES}})
		},
	})
}

func (gui *Gui) handleResetSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcResettingSubmoduleStatus, func() error {
		return gui.resetSubmodule(submodule)
	})
}

func (gui *Gui) fileForSubmodule(submodule *models.SubmoduleConfig) *models.File {
	for _, file := range gui.State.FileManager.GetAllFiles() {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}

func (gui *Gui) resetSubmodule(submodule *models.SubmoduleConfig) error {
	gitCommand := gui.GitCommand.WithSpan(gui.Tr.Spans.ResetSubmodule)

	file := gui.fileForSubmodule(submodule)
	if file != nil {
		if err := gitCommand.UnStageFile(file.Names(), file.Tracked); err != nil {
			return gui.SurfaceError(err)
		}
	}

	if err := gitCommand.SubmoduleStash(submodule); err != nil {
		return gui.SurfaceError(err)
	}
	if err := gitCommand.SubmoduleReset(submodule); err != nil {
		return gui.SurfaceError(err)
	}

	return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES, SUBMODULES}})
}

func (gui *Gui) handleAddSubmodule() error {
	return gui.Prompt(PromptOpts{
		Title: gui.Tr.LcNewSubmoduleUrl,
		HandleConfirm: func(submoduleUrl string) error {
			nameSuggestion := filepath.Base(strings.TrimSuffix(submoduleUrl, filepath.Ext(submoduleUrl)))

			return gui.Prompt(PromptOpts{
				Title:          gui.Tr.LcNewSubmoduleName,
				InitialContent: nameSuggestion,
				HandleConfirm: func(submoduleName string) error {

					return gui.Prompt(PromptOpts{
						Title:          gui.Tr.LcNewSubmodulePath,
						InitialContent: submoduleName,
						HandleConfirm: func(submodulePath string) error {
							return gui.WithWaitingStatus(gui.Tr.LcAddingSubmoduleStatus, func() error {
								err := gui.GitCommand.WithSpan(gui.Tr.Spans.AddSubmodule).SubmoduleAdd(submoduleName, submodulePath, submoduleUrl)
								if err != nil {
									return gui.SurfaceError(err)
								}

								return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
							})
						},
					})
				},
			})
		},
	})

}

func (gui *Gui) handleEditSubmoduleUrl(submodule *models.SubmoduleConfig) error {
	return gui.Prompt(PromptOpts{
		Title:          fmt.Sprintf(gui.Tr.LcUpdateSubmoduleUrl, submodule.Name),
		InitialContent: submodule.Url,
		HandleConfirm: func(newUrl string) error {
			return gui.WithWaitingStatus(gui.Tr.LcUpdatingSubmoduleUrlStatus, func() error {
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.UpdateSubmoduleUrl).SubmoduleUpdateUrl(submodule.Name, submodule.Path, newUrl)
				if err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
			})
		},
	})
}

func (gui *Gui) handleSubmoduleInit(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcInitializingSubmoduleStatus, func() error {
		err := gui.GitCommand.WithSpan(gui.Tr.Spans.InitialiseSubmodule).SubmoduleInit(submodule.Path)
		if err != nil {
			return gui.SurfaceError(err)
		}

		return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
	})
}

func (gui *Gui) forSubmodule(callback func(*models.SubmoduleConfig) error) func() error {
	return func() error {
		submodule := gui.getSelectedSubmodule()
		if submodule == nil {
			return nil
		}

		return callback(submodule)
	}
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
	bulkInitCmdObj := gui.GitCommand.SubmoduleBulkInitCmdObj()
	bulkUpdateCmdObj := gui.GitCommand.SubmoduleBulkUpdateCmdObj()
	bulkDeinitCmdObj := gui.GitCommand.SubmoduleBulkDeinitCmdObj()
	bulkForceUpdateCmdObj := gui.GitCommand.SubmoduleForceBulkUpdateCmdObj()

	menuItems := []*menuItem{
		{
			displayStrings: []string{gui.Tr.LcBulkInitSubmodules, utils.ColoredString(bulkInitCmdObj.ToString(), color.FgGreen)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.WithSpan(gui.Tr.Spans.BulkInitialiseSubmodules).RunExecutable(bulkInitCmdObj); err != nil {
						return gui.SurfaceError(err)
					}

					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcBulkUpdateSubmodules, utils.ColoredString(bulkUpdateCmdObj.ToString(), color.FgYellow)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.WithSpan(gui.Tr.Spans.BulkUpdateSubmodules).RunExecutable(bulkUpdateCmdObj); err != nil {
						return gui.SurfaceError(err)
					}

					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcSubmoduleStashAndReset, utils.ColoredString(fmt.Sprintf("git stash in each submodule && %s", bulkForceUpdateCmdObj.ToString()), color.FgRed)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.GitCommand.WithSpan(gui.Tr.Spans.BulkStashAndResetSubmodules).ResetSubmodules(gui.State.Submodules); err != nil {
						return gui.SurfaceError(err)
					}

					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
				})
			},
		},
		{
			displayStrings: []string{gui.Tr.LcBulkDeinitSubmodules, utils.ColoredString(bulkDeinitCmdObj.ToString(), color.FgRed)},
			onPress: func() error {
				return gui.WithWaitingStatus(gui.Tr.LcRunningCommand, func() error {
					if err := gui.OSCommand.WithSpan(gui.Tr.Spans.BulkDeinitialiseSubmodules).RunExecutable(bulkDeinitCmdObj); err != nil {
						return gui.SurfaceError(err)
					}

					return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
				})
			},
		},
	}

	return gui.createMenu(gui.Tr.LcBulkSubmoduleOptions, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleUpdateSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.WithWaitingStatus(gui.Tr.LcUpdatingSubmoduleStatus, func() error {
		err := gui.GitCommand.WithSpan(gui.Tr.Spans.UpdateSubmodule).SubmoduleUpdate(submodule.Path)
		if err != nil {
			return gui.SurfaceError(err)
		}

		return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{SUBMODULES}})
	})
}
