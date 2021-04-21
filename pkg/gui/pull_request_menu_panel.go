package gui

import (
	"fmt"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (gui *Gui) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*menuItem, 0, 2)

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems, &menuItem{
			displayStrings: []string{
				fmt.Sprintf("%s -> default branch", selectedBranch.Name),
			},
			onPress: func() error {
				return createPullRequest(selectedBranch, nil, gui)
			},
		}, &menuItem{
			displayStrings: []string{
				fmt.Sprintf("%s -> %s", checkedOutBranch.Name, selectedBranch.Name),
			},
			onPress: func() error {
				return createPullRequest(checkedOutBranch, selectedBranch, gui)
			},
		})
	}

	menuItems = append(menuItems, &menuItem{
		displayStrings: []string{
			fmt.Sprintf("%s -> default branch", checkedOutBranch.Name),
		},
		onPress: func() error {
			return createPullRequest(checkedOutBranch, nil, gui)
		},
	})

	return gui.createMenu(fmt.Sprintf(gui.Tr.CreatePullRequest), menuItems, createMenuOptions{showCancel: true})
}

func createPullRequest(checkedOutBranch *models.Branch, selectedBranch *models.Branch, gui *Gui) error {
	pullRequest := commands.NewPullRequest(gui.GitCommand)
	url, err := pullRequest.Create(checkedOutBranch, selectedBranch)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.OnRunCommand(oscommands.NewCmdLogEntry(fmt.Sprintf("Creating pull request at URL: %s", url), "Create pull request", false))

	return nil
}
