package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) gitFlowFinishBranch(gitFlowConfig string, branchName string) error {
	// need to find out what kind of branch this is
	prefix := strings.SplitAfterN(branchName, "/", 2)[0]
	suffix := strings.Replace(branchName, prefix, "", 1)

	branchType := ""
	for _, line := range strings.Split(strings.TrimSpace(gitFlowConfig), "\n") {
		if strings.HasPrefix(line, "gitflow.prefix.") && strings.HasSuffix(line, prefix) {
			// now I just need to how do you say
			regex := regexp.MustCompile("gitflow.prefix.([^ ]*) .*")
			matches := regex.FindAllStringSubmatch(line, 1)

			if len(matches) > 0 && len(matches[0]) > 1 {
				branchType = matches[0][1]
				break
			}
		}
	}

	if branchType == "" {
		return gui.createErrorPanel(gui.Tr.NotAGitFlowBranch)
	}

	return gui.runSubprocessWithSuspenseAndRefresh(
		gui.GitCommand.WithSpan(gui.Tr.Spans.GitFlowFinish).NewCmdObjWithLog("git flow " + branchType + " finish " + suffix),
	)
}

func (gui *Gui) handleCreateGitFlowMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	// get config
	gitFlowConfig, err := gui.GitCommand.RunWithOutput(gui.GitCommand.NewCmdObj("git config --local --get-regexp gitflow"))
	if err != nil {
		return gui.createErrorPanel("You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := utils.ResolvePlaceholderString(gui.Tr.NewGitFlowBranchPrompt, map[string]string{"branchType": branchType})

			return gui.prompt(promptOpts{
				title: title,
				handleConfirm: func(name string) error {
					return gui.runSubprocessWithSuspenseAndRefresh(
						gui.GitCommand.WithSpan(gui.Tr.Spans.GitFlowStart).NewCmdObjWithLog("git flow " + branchType + " start " + name),
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
				return gui.gitFlowFinishBranch(gitFlowConfig, branch.Name)
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
