package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// never call this on its own, it should only be called from within refreshCommits()
func (gui *Gui) refreshStatus() {
	gui.State.RefreshingStatusMutex.Lock()
	defer gui.State.RefreshingStatusMutex.Unlock()

	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.Pushables != "" && currentBranch.Pullables != "" {
		trackColor := color.FgYellow
		if currentBranch.Pushables == "0" && currentBranch.Pullables == "0" {
			trackColor = color.FgGreen
		} else if currentBranch.Pushables == "?" && currentBranch.Pullables == "?" {
			trackColor = color.FgRed
		}

		status = utils.ColoredString(fmt.Sprintf("↑%s↓%s ", currentBranch.Pushables, currentBranch.Pullables), trackColor)
	}

	if gui.GitCommand.WorkingTreeState() != "normal" {
		status += utils.ColoredString(fmt.Sprintf("(%s) ", gui.GitCommand.WorkingTreeState()), color.FgYellow)
	}

	name := utils.ColoredString(currentBranch.Name, presentation.GetBranchColor(currentBranch.Name))
	repoName := utils.GetCurrentRepoName()
	status += fmt.Sprintf("%s → %s ", repoName, name)

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.getStatusView(), status)
		return nil
	})
}

func runeCount(str string) int {
	return len([]rune(str))
}

func cursorInSubstring(cx int, prefix string, substring string) bool {
	return cx >= runeCount(prefix) && cx < runeCount(prefix+substring)
}

func (gui *Gui) handleCheckForUpdate(g *gocui.Gui, v *gocui.View) error {
	gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)
	return gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("CheckingForUpdates"))
}

func (gui *Gui) handleStatusClick(g *gocui.Gui, v *gocui.View) error {
	currentBranch := gui.currentBranch()

	cx, _ := v.Cursor()
	upstreamStatus := fmt.Sprintf("↑%s↓%s", currentBranch.Pushables, currentBranch.Pullables)
	repoName := utils.GetCurrentRepoName()
	switch gui.GitCommand.WorkingTreeState() {
	case "rebasing", "merging":
		workingTreeStatus := fmt.Sprintf("(%s)", gui.GitCommand.WorkingTreeState())
		if cursorInSubstring(cx, upstreamStatus+" ", workingTreeStatus) {
			return gui.handleCreateRebaseOptionsMenu(gui.g, v)
		}
		if cursorInSubstring(cx, upstreamStatus+" "+workingTreeStatus+" ", repoName) {
			return gui.handleCreateRecentReposMenu(gui.g, v)
		}
	default:
		if cursorInSubstring(cx, upstreamStatus+" ", repoName) {
			return gui.handleCreateRecentReposMenu(gui.g, v)
		}
	}

	return gui.handleStatusSelect(gui.g, v)
}

func (gui *Gui) handleStatusSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = ""

	if gui.inDiffMode() {
		return gui.renderDiff()
	}

	magenta := color.New(color.FgMagenta)

	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			"Copyright (c) 2018 Jesse Duffield",
			"Keybindings: https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings",
			"Config Options: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md",
			"Tutorial: https://youtu.be/VDXvbHZYeKY",
			"Raise an Issue: https://github.com/jesseduffield/lazygit/issues",
			magenta.Sprint("Become a sponsor (github is matching all donations for 12 months): https://github.com/sponsors/jesseduffield"), // caffeine ain't free
		}, "\n\n")

	return gui.newStringTask("main", dashboardString)
}

func (gui *Gui) handleOpenConfig(g *gocui.Gui, v *gocui.View) error {
	return gui.openFile(gui.Config.GetUserConfig().ConfigFileUsed())
}

func (gui *Gui) handleEditConfig(g *gocui.Gui, v *gocui.View) error {
	filename := gui.Config.GetUserConfig().ConfigFileUsed()
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

func (gui *Gui) workingTreeState() string {
	rebaseMode, _ := gui.GitCommand.RebaseMode()
	if rebaseMode != "" {
		return "rebasing"
	}
	merging, _ := gui.GitCommand.IsInMergeState()
	if merging {
		return "merging"
	}
	return "normal"
}
