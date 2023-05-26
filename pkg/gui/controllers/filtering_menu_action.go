package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FilteringMenuAction struct {
	c *ControllerCommon
}

func (self *FilteringMenuAction) Call() error {
	fileName := ""
	switch self.c.CurrentSideContext() {
	case self.c.Contexts().Files:
		node := self.c.Contexts().Files.GetSelected()
		if node != nil {
			fileName = node.GetPath()
		}
	case self.c.Contexts().CommitFiles:
		node := self.c.Contexts().CommitFiles.GetSelected()
		if node != nil {
			fileName = node.GetPath()
		}
	}

	menuItems := []*types.MenuItem{}

	if fileName != "" {
		menuItems = append(menuItems, &types.MenuItem{
			Label: fmt.Sprintf("%s '%s'", self.c.Tr.FilterBy, fileName),
			OnPress: func() error {
				return self.setFiltering(fileName)
			},
		})
	}

	menuItems = append(menuItems, &types.MenuItem{
		Label: self.c.Tr.FilterPathOption,
		OnPress: func() error {
			return self.c.Prompt(types.PromptOpts{
				FindSuggestionsFunc: self.c.Helpers().Suggestions.GetFilePathSuggestionsFunc(),
				Title:               self.c.Tr.EnterFileName,
				HandleConfirm: func(response string) error {
					return self.setFiltering(strings.TrimSpace(response))
				},
			})
		},
	})

	if self.c.Modes().Filtering.Active() {
		menuItems = append(menuItems, &types.MenuItem{
			Label:   self.c.Tr.ExitFilterMode,
			OnPress: self.c.Helpers().Mode.ClearFiltering,
		})
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.FilteringMenuTitle, Items: menuItems})
}

func (self *FilteringMenuAction) setFiltering(path string) error {
	self.c.Modes().Filtering.SetPath(path)

	repoState := self.c.State().GetRepoState()
	if repoState.GetScreenMode() == types.SCREEN_NORMAL {
		repoState.SetScreenMode(types.SCREEN_HALF)
	}

	if err := self.c.PushContext(self.c.Contexts().LocalCommits); err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}, Then: func() {
		self.c.Contexts().LocalCommits.SetSelectedLineIdx(0)
		self.c.Contexts().LocalCommits.FocusLine()
	}})
}
