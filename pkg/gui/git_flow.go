package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
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
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("NotAGitFlowBranch"))
	}

	subProcess := gui.OSCommand.PrepareSubProcess("git", "flow", branchType, "finish", suffix)
	gui.SubProcess = subProcess
	return gui.Errors.ErrSubProcess
}

func (gui *Gui) handleCreateGitFlowMenu(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	// get config
	gitFlowConfig, err := gui.OSCommand.RunCommandWithOutput("git config --local --get-regexp gitflow")
	if err != nil {
		return gui.createErrorPanel(gui.g, "You need to install git-flow and enable it in this repo to use git-flow features")
	}

	startHandler := func(branchType string) func() error {
		return func() error {
			title := gui.Tr.TemplateLocalize("NewBranchNamePrompt", map[string]interface{}{"branchType": branchType})
			return gui.createPromptPanel(gui.g, gui.getMenuView(), title, "", func(g *gocui.Gui, v *gocui.View) error {
				name := gui.trimmedContent(v)
				subProcess := gui.OSCommand.PrepareSubProcess("git", "flow", branchType, "start", name)
				gui.SubProcess = subProcess
				return gui.Errors.ErrSubProcess
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
			displayString: "start release",
			onPress:       startHandler("release"),
		},
	}

	return gui.createMenu("git flow", menuItems, createMenuOptions{})
}
