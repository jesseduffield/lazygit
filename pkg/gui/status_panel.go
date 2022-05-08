package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func runeCount(str string) int {
	return len([]rune(str))
}

func cursorInSubstring(cx int, prefix string, substring string) bool {
	return cx >= runeCount(prefix) && cx < runeCount(prefix+substring)
}

func (gui *Gui) handleCheckForUpdate() error {
	return gui.c.WithWaitingStatus(gui.c.Tr.CheckingForUpdates, func() error {
		gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)
		return nil
	})
}

func (gui *Gui) handleStatusClick() error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	currentBranch := gui.helpers.Refs.GetCheckedOutRef()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	if err := gui.c.PushContext(gui.State.Contexts.Status); err != nil {
		return err
	}

	cx, _ := gui.Views.Status.Cursor()
	upstreamStatus := presentation.BranchStatus(currentBranch, gui.Tr)
	repoName := utils.GetCurrentRepoName()
	workingTreeState := gui.git.Status.WorkingTreeState()
	switch workingTreeState {
	case enums.REBASE_MODE_REBASING, enums.REBASE_MODE_MERGING:
		workingTreeStatus := fmt.Sprintf("(%s)", formatWorkingTreeState(workingTreeState))
		if cursorInSubstring(cx, upstreamStatus+" ", workingTreeStatus) {
			return gui.helpers.MergeAndRebase.CreateRebaseOptionsMenu()
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

func formatWorkingTreeState(rebaseMode enums.RebaseMode) string {
	switch rebaseMode {
	case enums.REBASE_MODE_REBASING:
		return "rebasing"
	case enums.REBASE_MODE_MERGING:
		return "merging"
	default:
		return "none"
	}
}

func (gui *Gui) statusRenderToMain() error {
	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			"Copyright 2022 Jesse Duffield",
			fmt.Sprintf("Keybindings: %s", constants.Links.Docs.Keybindings),
			fmt.Sprintf("Config Options: %s", constants.Links.Docs.Config),
			fmt.Sprintf("Tutorial: %s", constants.Links.Docs.Tutorial),
			fmt.Sprintf("Raise an Issue: %s", constants.Links.Issues),
			fmt.Sprintf("Release Notes: %s", constants.Links.Releases),
			style.FgMagenta.Sprintf("Become a sponsor: %s", constants.Links.Donate), // caffeine ain't free
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
		return errors.New(gui.c.Tr.NoConfigFileFoundErr)
	case 1:
		return action(confPaths[0])
	default:
		menuItems := slices.Map(confPaths, func(path string) *types.MenuItem {
			return &types.MenuItem{
				Label: path,
				OnPress: func() error {
					return action(path)
				},
			}
		})

		return gui.c.Menu(types.CreateMenuOptions{
			Title: gui.c.Tr.SelectConfigFile,
			Items: menuItems,
		})
	}
}

func (gui *Gui) handleOpenConfig() error {
	return gui.askForConfigFile(gui.helpers.Files.OpenFile)
}

func (gui *Gui) handleEditConfig() error {
	return gui.askForConfigFile(gui.helpers.Files.EditFile)
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
