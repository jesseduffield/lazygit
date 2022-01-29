package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateGitFlowMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if !gui.git.Flow.GitFlowEnabled() {
		return gui.c.ErrorMsg("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(gui.c.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			return gui.c.Prompt(types.PromptOpts{
				Title: title,
				HandleConfirm: func(name string) error {
					gui.c.LogAction(gui.c.Tr.Actions.GitFlowStart)
					return gui.runSubprocessWithSuspenseAndRefresh(
						gui.git.Flow.StartCmdObj(branchType, name),
					)
				},
			})
		}
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "git flow",
		Items: []*types.MenuItem{
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
	cmdObj, err := gui.git.Flow.FinishCmdObj(branchName)
	if err != nil {
		return gui.c.Error(err)
	}

	gui.c.LogAction(gui.c.Tr.Actions.GitFlowFinish)
	return gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
}
