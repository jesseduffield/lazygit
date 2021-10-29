package gui

import (
	"fmt"
	"strconv"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (gui *Gui) createOrOpenPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*menuItem, 0, 4)

	fromToDisplayStrings := func(from string, to string) []string {
		return []string{fmt.Sprintf("%s → %s", from, to)}
	}

	menuItemsForBranch := func(branch *models.Branch) []*menuItem {
		return []*menuItem{
			{
				displayStrings: fromToDisplayStrings(branch.Name, gui.Tr.LcDefaultBranch),
				onPress: func() error {
					return gui.createPullRequest(branch.Name, "")
				},
			},
			{
				displayStrings: fromToDisplayStrings(branch.Name, gui.Tr.LcSelectBranch),
				onPress: func() error {
					return gui.prompt(promptOpts{
						title:               branch.Name + " →",
						findSuggestionsFunc: gui.getBranchNameSuggestionsFunc(),
						handleConfirm: func(targetBranchName string) error {
							return gui.createPullRequest(branch.Name, targetBranchName)
						}},
					)
				},
			},
		}
	}

	if selectedBranch.PR != nil {
		menuItems = append(menuItems, &menuItem{
			displayString: "open #" + strconv.Itoa(selectedBranch.PR.Number),
			onPress: func() error {
				return gui.OSCommand.OpenLink(selectedBranch.PR.Url)
			},
		})
	}

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems,
			&menuItem{
				displayStrings: fromToDisplayStrings(checkedOutBranch.Name, selectedBranch.Name),
				onPress: func() error {
					return gui.createPullRequest(checkedOutBranch.Name, selectedBranch.Name)
				},
			},
		)
		menuItems = append(menuItems, menuItemsForBranch(checkedOutBranch)...)
	}

	menuItems = append(menuItems, menuItemsForBranch(selectedBranch)...)

	return gui.createMenu(fmt.Sprintf(gui.Tr.CreateOrOpenPullRequestOptions), menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) createPullRequest(from string, to string) error {
	pullRequest := commands.NewPullRequest(gui.GitCommand)
	url, err := pullRequest.Create(from, to)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.OnRunCommand(oscommands.NewCmdLogEntry(fmt.Sprintf(gui.Tr.CreatingPullRequestAtUrl, url), gui.Tr.CreateOrShowPullRequest, false))

	return nil
}
