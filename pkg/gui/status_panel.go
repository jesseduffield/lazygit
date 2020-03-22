package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getCurrentBranchTrack() (string, string) {
	currentBranch := gui.currentBranch()
	if currentBranch != nil {
		return currentBranch.Pushables, currentBranch.Pullables
	}
	return "?", "?"
}

// refreshStatus is dependent on state that's set in the refreshCommits and refreshBranches methods.
// It needs to know the pushable/pullable changes for the current branch (determines in refreshBranches)
// and it needs to know the current worktree state (determined in refreshCommits).
// refreshStatus should never be called on its own: it should only ever be called from within one of those
// two other methods. Because the two other methods can be called at roughly the same time we use a mutex here
// so that we're never rendering old information
func (gui *Gui) refreshStatus() {
	gui.State.RefreshingStatusMutex.Lock()
	defer gui.State.RefreshingStatusMutex.Unlock()

	pushables, pullables := gui.getCurrentBranchTrack()

	status := ""

	if pushables != "" && pullables != "" {
		trackColor := color.FgYellow
		if pushables == "0" && pullables == "0" {
			trackColor = color.FgGreen
		} else if pushables == "?" && pullables == "?" {
			trackColor = color.FgRed
		}

		status = utils.ColoredString(fmt.Sprintf("↑%s↓%s ", pushables, pullables), trackColor)
	}

	branches := gui.State.Branches

	if gui.State.WorkingTreeState != "normal" {
		status += utils.ColoredString(fmt.Sprintf("(%s) ", gui.State.WorkingTreeState), color.FgYellow)
	}

	if len(branches) > 0 {
		branch := branches[0]
		name := utils.ColoredString(branch.Name, presentation.GetBranchColor(branch.Name))
		repoName := utils.GetCurrentRepoName()
		status += fmt.Sprintf("%s → %s ", repoName, name)
	}

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.g, gui.getStatusView(), status)
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
	pushables, pullables := gui.getCurrentBranchTrack()

	cx, _ := v.Cursor()
	upstreamStatus := fmt.Sprintf("↑%s↓%s", pushables, pullables)
	repoName := utils.GetCurrentRepoName()
	gui.Log.Warn(gui.State.WorkingTreeState)
	switch gui.State.WorkingTreeState {
	case "rebasing", "merging":
		workingTreeStatus := fmt.Sprintf("(%s)", gui.State.WorkingTreeState)
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

func (gui *Gui) updateWorkTreeState() error {
	rebaseMode, err := gui.GitCommand.RebaseMode()
	if err != nil {
		return err
	}
	if rebaseMode != "" {
		gui.State.WorkingTreeState = "rebasing"
		return nil
	}
	merging, err := gui.GitCommand.IsInMergeState()
	if err != nil {
		return err
	}
	if merging {
		gui.State.WorkingTreeState = "merging"
		return nil
	}
	gui.State.WorkingTreeState = "normal"
	return nil
}
