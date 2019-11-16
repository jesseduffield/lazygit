package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) refreshStatus(g *gocui.Gui) error {
	state := gui.State.Panels.Status

	v, err := g.View("status")
	if err != nil {
		panic(err)
	}
	// for some reason if this isn't wrapped in an update the clear seems to
	// be applied after the other things or something like that; the panel's
	// contents end up cleared
	g.Update(func(*gocui.Gui) error {
		v.Clear()
		state.pushables, state.pullables = gui.GitCommand.GetCurrentBranchUpstreamDifferenceCount()
		if err := gui.updateWorkTreeState(); err != nil {
			return err
		}
		status := fmt.Sprintf("↑%s↓%s", state.pushables, state.pullables)
		branches := gui.State.Branches

		if gui.State.WorkingTreeState != "normal" {
			status += utils.ColoredString(fmt.Sprintf(" (%s)", gui.State.WorkingTreeState), color.FgYellow)
		}

		if len(branches) > 0 {
			branch := branches[0]
			name := utils.ColoredString(branch.Name, commands.GetBranchColor(branch.Name))
			repoName := utils.GetCurrentRepoName()
			status += fmt.Sprintf(" %s → %s", repoName, name)
		}

		fmt.Fprint(v, status)
		return nil
	})

	return nil
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
	state := gui.State.Panels.Status

	cx, _ := v.Cursor()
	upstreamStatus := fmt.Sprintf("↑%s↓%s", state.pushables, state.pullables)
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

	return gui.renderString(g, "main", dashboardString)
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
	merging, err := gui.GitCommand.IsInMergeState()
	if err != nil {
		return err
	}
	if merging {
		gui.State.WorkingTreeState = "merging"
		return nil
	}
	rebaseMode, err := gui.GitCommand.RebaseMode()
	if err != nil {
		return err
	}
	if rebaseMode != "" {
		gui.State.WorkingTreeState = "rebasing"
		return nil
	}
	gui.State.WorkingTreeState = "normal"
	return nil
}
