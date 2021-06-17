package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) exitDiffMode() error {
	gui.State.Modes.Diffing = diffing.New()
	return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
}

func (gui *Gui) renderDiff() error {
	cmdObj := gui.Git.Diff().ShowFileDiffCmdObj(
		gui.State.Modes.Diffing.Ref,
		gui.currentDiffTerminal(),
		gui.State.Modes.Diffing.Reverse,
		gui.fileToDiff(),
		false,
		true, // not sure if it actually matters much whether I show renames here
	)

	task := NewRunPtyTask(cmdObj)

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
	case FILES_CONTEXT_KEY, SUBMODULES_CONTEXT_KEY:
		// TODO: should we just return nil here?
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
		context := gui.currentSideListContext()
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
	return gui.Git.Diff().DiffEndArgs(
		gui.State.Modes.Diffing.Ref,
		gui.currentDiffTerminal(),
		gui.State.Modes.Diffing.Reverse,
		gui.fileToDiff(),
	)
}

func (gui *Gui) fileToDiff() string {
	file := gui.currentlySelectedFilename()
	if file != "" {
		return file
	}

	return gui.State.Modes.Filtering.GetPath()
}

func (gui *Gui) handleCreateDiffingMenuPanel() error {
	names := gui.currentDiffTerminals()

	menuItems := []*menuItem{}
	for _, name := range names {
		name := name
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: fmt.Sprintf("%s %s", gui.Tr.LcDiff, name),
				onPress: func() error {
					gui.State.Modes.Diffing.Ref = name
					// can scope this down based on current view but too lazy right now
					return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
				},
			},
		}...)
	}

	menuItems = append(menuItems, []*menuItem{
		{
			displayString: gui.Tr.LcEnterRefToDiff,
			onPress: func() error {
				return gui.Prompt(PromptOpts{
					Title: gui.Tr.LcEnteRefName,
					HandleConfirm: func(response string) error {
						gui.State.Modes.Diffing.Ref = strings.TrimSpace(response)
						return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
					},
				})
			},
		},
	}...)

	if gui.State.Modes.Diffing.Active() {
		menuItems = append(menuItems, []*menuItem{
			{
				displayString: gui.Tr.LcSwapDiff,
				onPress: func() error {
					gui.State.Modes.Diffing.Reverse = !gui.State.Modes.Diffing.Reverse
					return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
				},
			},
			{
				displayString: gui.Tr.LcExitDiffMode,
				onPress: func() error {
					gui.State.Modes.Diffing = diffing.New()
					return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
				},
			},
		}...)
	}

	return gui.createMenu(gui.Tr.DiffingMenuTitle, menuItems, createMenuOptions{showCancel: true})
}
