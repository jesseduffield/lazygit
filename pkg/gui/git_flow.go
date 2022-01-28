package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateGitFlowMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if !gui.Git.Flow.GitFlowEnabled() {
		return gui.PopupHandler.ErrorMsg("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(gui.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			return gui.PopupHandler.Prompt(popup.PromptOpts{
				Title: title,
				HandleConfirm: func(name string) error {
					gui.logAction(gui.Tr.Actions.GitFlowStart)
					return gui.runSubprocessWithSuspenseAndRefresh(
						gui.Git.Flow.StartCmdObj(branchType, name),
					)
				},
			})
		}
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: "git flow",
		Items: []*popup.MenuItem{
			{
				// not localising here because it's one to one with the actual git flow commands
				DisplayString: fmt.Sprintf("finish branch '%s'", branch.Name),
				OnPress: func() error {
					return gui.gitFlowFinishBranch(branch.Name)
				},
			},
			{
				DisplayString: "start feature",
				OnPress:       startHandler("feature"),
			},
			{
				DisplayString: "start hotfix",
				OnPress:       startHandler("hotfix"),
			},
			{
				DisplayString: "start bugfix",
				OnPress:       startHandler("bugfix"),
			},
			{
				DisplayString: "start release",
				OnPress:       startHandler("release"),
			},
		},
	})
}

func (gui *Gui) gitFlowFinishBranch(branchName string) error {
	cmdObj, err := gui.Git.Flow.FinishCmdObj(branchName)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.logAction(gui.Tr.Actions.GitFlowFinish)
	return gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
}
