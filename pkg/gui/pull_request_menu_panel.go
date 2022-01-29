package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*types.MenuItem, 0, 4)

	fromToDisplayStrings := func(from string, to string) []string {
		return []string{fmt.Sprintf("%s → %s", from, to)}
	}

	menuItemsForBranch := func(branch *models.Branch) []*types.MenuItem {
		return []*types.MenuItem{
			{
				DisplayStrings: fromToDisplayStrings(branch.Name, gui.c.Tr.LcDefaultBranch),
				OnPress: func() error {
					return gui.createPullRequest(branch.Name, "")
				},
			},
			{
				DisplayStrings: fromToDisplayStrings(branch.Name, gui.c.Tr.LcSelectBranch),
				OnPress: func() error {
					return gui.c.Prompt(types.PromptOpts{
						Title:               branch.Name + " →",
						FindSuggestionsFunc: gui.helpers.suggestions.GetBranchNameSuggestionsFunc(),
						HandleConfirm: func(targetBranchName string) error {
							return gui.createPullRequest(branch.Name, targetBranchName)
						}},
					)
				},
			},
		}
	}

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems,
			&types.MenuItem{
				DisplayStrings: fromToDisplayStrings(checkedOutBranch.Name, selectedBranch.Name),
				OnPress: func() error {
					return gui.createPullRequest(checkedOutBranch.Name, selectedBranch.Name)
				},
			},
		)
		menuItems = append(menuItems, menuItemsForBranch(checkedOutBranch)...)
	}

	menuItems = append(menuItems, menuItemsForBranch(selectedBranch)...)

	return gui.c.Menu(types.CreateMenuOptions{Title: fmt.Sprintf(gui.c.Tr.CreatePullRequestOptions), Items: menuItems})
}

func (gui *Gui) createPullRequest(from string, to string) error {
	hostingServiceMgr := gui.getHostingServiceMgr()
	url, err := hostingServiceMgr.GetPullRequestURL(from, to)
	if err != nil {
		return gui.c.Error(err)
	}

	gui.c.LogAction(gui.c.Tr.Actions.OpenPullRequest)

	if err := gui.OSCommand.OpenLink(url); err != nil {
		return gui.c.Error(err)
	}

	return nil
}

func (gui *Gui) getHostingServiceMgr() *hosting_service.HostingServiceMgr {
	remoteUrl := gui.git.Config.GetRemoteURL()
	configServices := gui.c.UserConfig.Services
	return hosting_service.NewHostingServiceMgr(gui.Log, gui.Tr, remoteUrl, configServices)
}
