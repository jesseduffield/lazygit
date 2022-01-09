package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateGitFlowMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if !gui.Git.Flow.GitFlowEnabled() {
		return gui.createErrorPanel("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(gui.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			return gui.prompt(promptOpts{
				title: title,
				handleConfirm: func(name string) error {
					gui.logAction(gui.Tr.Actions.GitFlowStart)
					return gui.runSubprocessWithSuspenseAndRefresh(
						gui.Git.Flow.StartCmdObj(branchType, name),
					)
				},
			})
		}
	}

	menuItems := []*menuItem{
		{
			// not localising here because it's one to one with the actual git flow commands
			displayString: fmt.Sprintf("finish branch '%s'", branch.Name),
			onPress: func() error {
				return gui.gitFlowFinishBranch(branch.Name)
			},
		},
		{
			displayString: "start feature",
			onPress:       startHandler("feature"),
		},
		{
			displayString: "start hotfix",
			onPress:       startHandler("hotfix"),
		},
		{
			displayString: "start bugfix",
			onPress:       startHandler("bugfix"),
		},
		{
			displayString: "start release",
			onPress:       startHandler("release"),
		},
	}

	return gui.createMenu("git flow", menuItems, createMenuOptions{})
}

func (gui *Gui) gitFlowFinishBranch(branchName string) error {
	cmdObj, err := gui.Git.Flow.FinishCmdObj(branchName)
	if err != nil {
		return gui.surfaceError(err)
	}

	gui.logAction(gui.Tr.Actions.GitFlowFinish)
	return gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
}
