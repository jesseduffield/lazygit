package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) exitDiffMode() error {
	gui.State.Modes.Diffing = Diffing{}
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) renderDiff() error {
	cmd := gui.OSCommand.ExecutableFromString(
		fmt.Sprintf("git diff --submodule --no-ext-diff --color %s", gui.diffStr()),
	)
	task := gui.createRunPtyTask(cmd)

	return gui.refreshMainViews(refreshMainOpts{
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
	switch gui.currentContext().GetKey() {
	case "":
		return nil
	case FILES_CONTEXT_KEY:
		return []string{""}
	case COMMIT_FILES_CONTEXT_KEY:
		return []string{gui.State.Panels.CommitFiles.refName}
	case LOCAL_BRANCHES_CONTEXT_KEY:
		// for our local branches we want to include both the branch and its upstream
		branch := gui.getSelectedBranch()
		if branch != nil {
			names := []string{branch.ID()}
			if branch.UpstreamName != "" {
				names = append(names, branch.UpstreamName)
			}
			return names
		}
		return nil
	default:
		context := gui.currentSideContext()
		if context == nil {
			return nil
		}
		item, ok := context.GetSelectedItem()
		if !ok {
			return nil
		}
		return []string{item.ID()}
	}
}

func (gui *Gui) currentDiffTerminal() string {
	names := gui.currentDiffTerminals()
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func (gui *Gui) currentlySelectedFilename() string {
	switch gui.currentContext().GetKey() {
	case FILES_CONTEXT_KEY, COMMIT_FILES_CONTEXT_KEY:
		return gui.getSideContextSelectedItemId()
	default:
		return ""
	}
}

func (gui *Gui) diffStr() string {
	output := gui.State.Modes.Diffing.Ref

	right := gui.currentDiffTerminal()
	if right != "" {
		output += " " + right
	}

	if gui.State.Modes.Diffing.Reverse {
		output += " -R"
	}

	file := gui.currentlySelectedFilename()
	if file != "" {
		output += " -- " + file
	} else if gui.State.Modes.Filtering.Active() {
		output += " -- " + gui.State.Modes.Filtering.Path
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
					gui.State.Modes.Diffing.Ref = name
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
				return gui.prompt(gui.Tr.SLocalize("enteRefName"), "", func(response string) error {
					gui.State.Modes.Diffing.Ref = strings.TrimSpace(response)
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				})
			},
		},
	}...)

	if gui.State.Modes.Diffing.Active() {
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: gui.Tr.SLocalize("swapDiff"),
				onPress: func() error {
					gui.State.Modes.Diffing.Reverse = !gui.State.Modes.Diffing.Reverse
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				},
			},
			{
				displayString: gui.Tr.SLocalize("exitDiffMode"),
				onPress: func() error {
					gui.State.Modes.Diffing = Diffing{}
					return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
				},
			},
		}...)
	}

	return gui.createMenu(gui.Tr.SLocalize("DiffingMenuTitle"), menuItems, createMenuOptions{showCancel: true})
}
