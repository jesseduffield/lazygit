package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) inDiffMode() bool {
	return gui.State.Diff.Ref != ""
}

func (gui *Gui) exitDiffMode() error {
	gui.State.Diff = DiffState{}
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) renderDiff() error {
	gui.getMainView().Title = "Diff"
	gui.State.SplitMainPanel = false
	filterArg := ""
	if gui.inFilterMode() {
		filterArg = fmt.Sprintf(" -- %s", gui.State.FilterPath)
	}

	cmd := gui.OSCommand.ExecutableFromString(
		fmt.Sprintf("git diff --color %s %s", gui.diffStr(), filterArg),
	)
	if err := gui.newPtyTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}
	return nil
}

// currentDiffTerminals returns the current diff terminals of the currently selected item.
// in the case of a branch it returns both the branch and it's upstream name,
// which becomes an option when you bring up the diff menu, but when you're just
// flicking through branches it will be using the local branch name.
func (gui *Gui) currentDiffTerminals() []string {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return nil
	}
	names := []string{}
	switch currentView.Name() {
	case "files":
	// not supporting files for now
	// file, err := gui.getSelectedFile()
	// if err == nil {
	// 	names = append(names, file.Name)
	// }
	case "commitFiles":
		// not supporting commit files for now
		// file := gui.getSelectedCommitFile()
		// if file != nil {
		// 	names = append(names, file.Name)
		// }
	case "commits":
		var commit *commands.Commit
		switch gui.getCommitsView().Context {
		case "reflog-commits":
			commit = gui.getSelectedReflogCommit()
		case "branch-commits":
			commit = gui.getSelectedCommit()
		}
		if commit != nil {
			names = append(names, commit.Sha)
		}
	case "stash":
		entry := gui.getSelectedStashEntry()
		if entry != nil {
			names = append(names, entry.RefName())
		}
	case "branches":
		switch gui.getBranchesView().Context {
		case "local-branches":
			branch := gui.getSelectedBranch()
			if branch != nil {
				names = append(names, branch.Name)
				if branch.UpstreamName != "" {
					names = append(names, branch.UpstreamName)
				}
			}
		case "remotes":
			remote := gui.getSelectedRemote()
			if remote != nil {
				names = append(names, remote.Name)
			}
		case "remote-branches":
			remoteBranch := gui.getSelectedRemoteBranch()
			if remoteBranch != nil {
				names = append(names, remoteBranch.FullName())
			}
		case "tags":
			tag := gui.getSelectedTag()
			if tag != nil {
				names = append(names, tag.Name)
			}
		}
	}
	return names
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
				return gui.createPromptPanel(gui.g, v, gui.Tr.SLocalize("enteRefName"), "", func(g *gocui.Gui, promptView *gocui.View) error {
					gui.State.Diff.Ref = strings.TrimSpace(promptView.Buffer())
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
