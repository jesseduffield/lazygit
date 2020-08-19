package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) inDiffMode() bool {
	return gui.State.Diff.Ref != ""
}

func (gui *Gui) exitDiffMode() error {
	gui.State.Diff = DiffState{}
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) renderDiff() error {
	filterArg := ""

	if gui.inFilterMode() {
		filterArg = fmt.Sprintf(" -- %s", gui.State.FilterPath)
	}

	cmd := gui.OSCommand.ExecutableFromString(
		fmt.Sprintf("git diff --color %s %s", gui.diffStr(), filterArg),
	)
	task := gui.createRunPtyTask(cmd)

	return gui.refreshMain(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Diff",
			task:  task,
		},
	})
}

// currentDiffTerminals returns the current diff terminals of the currently selected item.
// in the case of a branch it returns both the branch and it's upstream name,
// which becomes an option when you bring up the diff menu, but when you're just
// flicking through branches it will be using the local branch name.
func (gui *Gui) currentDiffTerminals() []string {
	switch gui.currentContextKey() {
	case "files":
		// not supporting files for now
	case "commit-files":
		// not supporting commit files for now
	case "branch-commits":
		item := gui.getSelectedCommit()
		if item != nil {
			return []string{item.RefName()}
		}
	case "reflog-commits":
		item := gui.getSelectedReflogCommit()
		if item != nil {
			return []string{item.RefName()}
		}
	case "stash":
		item := gui.getSelectedStashEntry()
		if item != nil {
			return []string{item.RefName()}
		}

	case "local-branches":
		branch := gui.getSelectedBranch()
		if branch != nil {
			names := []string{branch.RefName()}
			if branch.UpstreamName != "" {
				names = append(names, branch.UpstreamName)
			}
			return names
		}
		return nil
	case "remotes":
		item := gui.getSelectedRemote()
		if item != nil {
			return []string{item.RefName()}
		}
	case "remote-branches":
		item := gui.getSelectedRemoteBranch()
		if item != nil {
			return []string{item.RefName()}
		}
	case "tags":
		item := gui.getSelectedTag()
		if item != nil {
			return []string{item.RefName()}
		}
	}

	return nil
}

func (gui *Gui) currentDiffTerminal() string {
	names := gui.currentDiffTerminals()
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func (gui *Gui) diffStr() string {
	output := gui.State.Diff.Ref

	right := gui.currentDiffTerminal()
	if right != "" {
		output += " " + right
	}
	if gui.State.Diff.Reverse {
		output += " -R"
	}
	return output
}

func (gui *Gui) handleCreateDiffingMenuPanel(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	names := gui.currentDiffTerminals()

	menuItems := []*menuItem{}
	for _, name := range names {
		name := name
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: fmt.Sprintf("%s %s", gui.Tr.SLocalize("diff"), name),
				onPress: func() error {
					gui.State.Diff.Ref = name
					// can scope this down based on current view but too lazy right now
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				},
			},
		}...)
	}

	menuItems = append(menuItems, []*menuItem{
		{
			displayString: gui.Tr.SLocalize("enterRefToDiff"),
			onPress: func() error {
				return gui.prompt(v, gui.Tr.SLocalize("enteRefName"), "", func(response string) error {
					gui.State.Diff.Ref = strings.TrimSpace(response)
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				})
			},
		},
	}...)

	if gui.inDiffMode() {
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: gui.Tr.SLocalize("swapDiff"),
				onPress: func() error {
					gui.State.Diff.Reverse = !gui.State.Diff.Reverse
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				},
			},
			{
				displayString: gui.Tr.SLocalize("exitDiffMode"),
				onPress: func() error {
					gui.State.Diff = DiffState{}
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				},
			},
		}...)
	}

	return gui.createMenu(gui.Tr.SLocalize("DiffingMenuTitle"), menuItems, createMenuOptions{showCancel: true})
}
