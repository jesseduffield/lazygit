package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) exitDiffMode() error {
	gui.State.Modes.Diffing = diffing.New()
	return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) renderDiff() error {
	cmdObj := gui.os.Cmd.New(
		fmt.Sprintf("git diff --submodule --no-ext-diff --color %s", gui.diffStr()),
	)
	task := NewRunPtyTask(cmdObj.GetCmd())

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
	c := gui.currentSideContext()

	if c.GetKey() == "" {
		return nil
	}

	switch v := c.(type) {
	case *context.WorkingTreeContext, *context.SubmodulesContext:
		// TODO: should we just return nil here?
		return []string{""}
	case *context.CommitFilesContext:
		return []string{v.GetRef().RefName()}
	case *context.BranchesContext:
		// for our local branches we want to include both the branch and its upstream
		branch := gui.State.Contexts.Branches.GetSelected()
		if branch != nil {
			names := []string{branch.ID()}
			if branch.IsTrackingRemote() {
				names = append(names, branch.ID()+"@{u}")
			}
			return names
		}
		return nil
	case types.IListContext:
		itemId := v.GetSelectedItemId()

		return []string{itemId}
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

func (gui *Gui) currentlySelectedFilename() string {
	switch gui.currentContext().GetKey() {
	case context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY:
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
		output += " -- " + gui.State.Modes.Filtering.GetPath()
	}

	return output
}

func (gui *Gui) handleCreateDiffingMenuPanel() error {
	names := gui.currentDiffTerminals()

	menuItems := []*types.MenuItem{}
	for _, name := range names {
		name := name
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: fmt.Sprintf("%s %s", gui.c.Tr.LcDiff, name),
				OnPress: func() error {
					gui.State.Modes.Diffing.Ref = name
					// can scope this down based on current view but too lazy right now
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	menuItems = append(menuItems, []*types.MenuItem{
		{
			Label: gui.c.Tr.LcEnterRefToDiff,
			OnPress: func() error {
				return gui.c.Prompt(types.PromptOpts{
					Title:               gui.c.Tr.LcEnteRefName,
					FindSuggestionsFunc: gui.helpers.Suggestions.GetRefsSuggestionsFunc(),
					HandleConfirm: func(response string) error {
						gui.State.Modes.Diffing.Ref = strings.TrimSpace(response)
						return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
					},
				})
			},
		},
	}...)

	if gui.State.Modes.Diffing.Active() {
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: gui.c.Tr.LcSwapDiff,
				OnPress: func() error {
					gui.State.Modes.Diffing.Reverse = !gui.State.Modes.Diffing.Reverse
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
			{
				Label: gui.c.Tr.LcExitDiffMode,
				OnPress: func() error {
					gui.State.Modes.Diffing = diffing.New()
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.DiffingMenuTitle, Items: menuItems})
}
