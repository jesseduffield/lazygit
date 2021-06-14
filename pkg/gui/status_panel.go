package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// never call this on its own, it should only be called from within refreshCommits()
func (gui *Gui) refreshStatus() {
	gui.Mutexes.RefreshingStatusMutex.Lock()
	defer gui.Mutexes.RefreshingStatusMutex.Unlock()

	currentBranch := gui.CurrentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.IsRealBranch() {
		status += presentation.ColoredBranchStatus(currentBranch) + " "
	}

	if gui.Git.WorkingTreeState() != commands.REBASE_MODE_NORMAL {
		status += utils.ColoredString(fmt.Sprintf("(%s) ", gui.Git.WorkingTreeState()), color.FgYellow)
	}

	name := utils.ColoredString(currentBranch.Name, presentation.GetBranchColor(currentBranch.Name))
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
	return gui.WithPopupWaitingStatus(gui.Tr.CheckingForUpdates, func() error {
		gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish)

		return nil
	})
}

func (gui *Gui) handleStatusClick() error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	if gui.PopupPanelFocused() {
		return nil
	}

	currentBranch := gui.CurrentBranch()
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
	switch gui.Git.WorkingTreeState() {
	case commands.REBASE_MODE_REBASING, commands.REBASE_MODE_MERGING:
		workingTreeStatus := fmt.Sprintf("(%s)", gui.Git.WorkingTreeState())
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

	return gui.handleStatusSelect()
}

func (gui *Gui) handleStatusSelect() error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	if gui.PopupPanelFocused() {
		return nil
	}

	magenta := color.New(color.FgMagenta)

	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			"Copyright (c) 2018 Jesse Duffield",
			fmt.Sprintf("Keybindings: %s", constants.Links.Docs.Keybindings),
			fmt.Sprintf("Config Options: %s", constants.Links.Docs.Config),
			fmt.Sprintf("Tutorial: %s", constants.Links.Docs.Tutorial),
			fmt.Sprintf("Raise an Issue: %s", constants.Links.Issues),
			fmt.Sprintf("Release Notes: %s", constants.Links.Releases),
			magenta.Sprintf("Become a sponsor (github is matching all donations for 12 months): %s", constants.Links.Donate), // caffeine ain't free
		}, "\n\n")

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "",
			task:  NewRenderStringTask(dashboardString),
		},
	})
}

func (gui *Gui) handleOpenConfig() error {
	return gui.openFile(gui.Config.GetUserConfigPath())
}

func (gui *Gui) handleEditConfig() error {
	filename := gui.Config.GetUserConfigPath()
	return gui.editFile(filename)
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

func (gui *Gui) workingTreeState() commands.WorkingTreeState {
	rebaseMode, _ := gui.Git.RebaseMode()
	if rebaseMode != "" {
		return commands.REBASE_MODE_REBASING
	}
	merging, _ := gui.Git.IsInMergeState()
	if merging {
		return commands.REBASE_MODE_MERGING
	}
	return commands.REBASE_MODE_NORMAL
}
