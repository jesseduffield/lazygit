package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (gui *Gui) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
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
					return gui.PopupHandler.Prompt(promptOpts{
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

	return gui.PopupHandler.Menu(createMenuOptions{title: fmt.Sprintf(gui.Tr.CreatePullRequestOptions), items: menuItems})
}

func (gui *Gui) createPullRequest(from string, to string) error {
	hostingServiceMgr := gui.getHostingServiceMgr()
	url, err := hostingServiceMgr.GetPullRequestURL(from, to)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.logAction(gui.Tr.Actions.OpenPullRequest)

	if err := gui.OSCommand.OpenLink(url); err != nil {
		return gui.PopupHandler.Error(err)
	}

	return nil
}

func (gui *Gui) getHostingServiceMgr() *hosting_service.HostingServiceMgr {
	remoteUrl := gui.Git.Config.GetRemoteURL()
	configServices := gui.UserConfig.Services
	return hosting_service.NewHostingServiceMgr(gui.Log, gui.Tr, remoteUrl, configServices)
}
