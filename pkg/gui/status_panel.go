package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// never call this on its own, it should only be called from within refreshCommits()
func (gui *Gui) refreshStatus() {
	gui.Mutexes.RefreshingStatusMutex.Lock()
	defer gui.Mutexes.RefreshingStatusMutex.Unlock()

	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.IsRealBranch() {
		status += presentation.ColoredBranchStatus(currentBranch) + " "
	}

	if gui.GitCommand.WorkingTreeState() != enums.REBASE_MODE_NONE {
		status += style.FgYellow.Sprintf("(%s) ", gui.GitCommand.WorkingTreeState())
	}

	name := presentation.GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	repoName := utils.GetCurrentRepoName()
	status += fmt.Sprintf("%s → %s ", repoName, name)

	gui.setViewContent(gui.Views.Status, status)
}

func runeCount(str string) int {
	return len([]rune(str))
}

func cursorInSubstring(cx int, prefix string, substring string) bool {
	return cx >= runeCount(prefix) && cx < runeCount(prefix+substring)
}

func (gui *Gui) handleCheckForUpdate() error {
	gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)
	return gui.createLoaderPanel(gui.Tr.CheckingForUpdates)
}

func (gui *Gui) handleStatusClick() error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	if gui.popupPanelFocused() {
		return nil
	}

	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	if err := gui.pushContext(gui.State.Contexts.Status); err != nil {
		return err
	}

	cx, _ := gui.Views.Status.Cursor()
	upstreamStatus := presentation.BranchStatus(currentBranch)
	repoName := utils.GetCurrentRepoName()
	workingTreeState := gui.GitCommand.WorkingTreeState()
	switch workingTreeState {
	case enums.REBASE_MODE_REBASING, enums.REBASE_MODE_MERGING:
		var formattedState string
		if workingTreeState == enums.REBASE_MODE_REBASING {
			formattedState = "rebasing"
		} else {
			formattedState = "merging"
		}
		workingTreeStatus := fmt.Sprintf("(%s)", formattedState)
		if cursorInSubstring(cx, upstreamStatus+" ", workingTreeStatus) {
			return gui.handleCreateRebaseOptionsMenu()
		}
		if cursorInSubstring(cx, upstreamStatus+" "+workingTreeStatus+" ", repoName) {
			return gui.handleCreateRecentReposMenu()
		}
	default:
		if cursorInSubstring(cx, upstreamStatus+" ", repoName) {
			return gui.handleCreateRecentReposMenu()
		}
	}

	return nil
}

func (gui *Gui) statusRenderToMain() error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	if gui.popupPanelFocused() {
		return nil
	}

	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			"Copyright (c) 2018 Jesse Duffield",
			fmt.Sprintf("Keybindings: %s", constants.Links.Docs.Keybindings),
			fmt.Sprintf("Config Options: %s", constants.Links.Docs.Config),
			fmt.Sprintf("Tutorial: %s", constants.Links.Docs.Tutorial),
			fmt.Sprintf("Raise an Issue: %s", constants.Links.Issues),
			fmt.Sprintf("Release Notes: %s", constants.Links.Releases),
			style.FgMagenta.Sprintf("Become a sponsor (github is matching all donations for 12 months): %s", constants.Links.Donate), // caffeine ain't free
		}, "\n\n")

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "",
			task:  NewRenderStringTask(dashboardString),
		},
	})
}

func (gui *Gui) askForConfigFile(action func(file string) error) error {
	confPaths := gui.Config.GetUserConfigPaths()
	switch len(confPaths) {
	case 0:
		return errors.New(gui.Tr.NoConfigFileFoundErr)
	case 1:
		return action(confPaths[0])
	default:
		menuItems := make([]*menuItem, len(confPaths))
		for i, file := range confPaths {
			i := i
			menuItems[i] = &menuItem{
				displayString: file,
				onPress: func() error {
					return action(confPaths[i])
				},
			}
		}
		return gui.createMenu(gui.Tr.SelectConfigFile, menuItems, createMenuOptions{})
	}
}

func (gui *Gui) handleOpenConfig() error {
	return gui.askForConfigFile(gui.openFile)
}

func (gui *Gui) handleEditConfig() error {
	return gui.askForConfigFile(gui.editFile)
}

func lazygitTitle() string {
	return `
   _                       _ _
  | |                     (_) |
  | | __ _ _____   _  __ _ _| |_
  | |/ _` + "`" + ` |_  / | | |/ _` + "`" + ` | | __|
  | | (_| |/ /| |_| | (_| | | |_
  |_|\__,_/___|\__, |\__, |_|\__|
                __/ | __/ |
               |___/ |___/       `
}
