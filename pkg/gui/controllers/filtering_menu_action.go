package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FilteringMenuAction struct {
	c *ControllerCommon
}

func (self *FilteringMenuAction) Call() error {
	fileName := ""
	author := ""
	switch self.c.Context().CurrentSide() {
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
	case self.c.Contexts().LocalCommits:
		commit := self.c.Contexts().LocalCommits.GetSelected()
		if commit != nil {
			author = fmt.Sprintf("%s <%s>", commit.AuthorName, commit.AuthorEmail)
		}
	}

	menuItems := []*types.MenuItem{}
	tooltip := ""
	if self.c.Modes().Filtering.Active() {
		tooltip = self.c.Tr.WillCancelExistingFilterTooltip
	}

	if fileName != "" {
		menuItems = append(menuItems, &types.MenuItem{
			Label: fmt.Sprintf("%s '%s'", self.c.Tr.FilterBy, fileName),
			OnPress: func() error {
				return self.c.Helpers().Mode.SetFilteringPath(fileName)
			},
			Tooltip: tooltip,
		})
	}

	if author != "" {
		menuItems = append(menuItems, &types.MenuItem{
			Label: fmt.Sprintf("%s '%s'", self.c.Tr.FilterBy, author),
			OnPress: func() error {
				return self.c.Helpers().Mode.SetFilteringAuthor(author)
			},
			Tooltip: tooltip,
		})
	}

	menuItems = append(menuItems, &types.MenuItem{
		Label: self.c.Tr.FilterPathOption,
		OnPress: func() error {
			self.c.Prompt(types.PromptOpts{
				FindSuggestionsFunc: self.c.Helpers().Suggestions.GetFilePathSuggestionsFunc(),
				Title:               self.c.Tr.EnterFileName,
				HandleConfirm: func(response string) error {
					return self.c.Helpers().Mode.SetFilteringPath(response)
				},
			})

			return nil
		},
		Tooltip: tooltip,
	})

	menuItems = append(menuItems, &types.MenuItem{
		Label: self.c.Tr.FilterAuthorOption,
		OnPress: func() error {
			self.c.Prompt(types.PromptOpts{
				FindSuggestionsFunc: self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc(),
				Title:               self.c.Tr.EnterAuthor,
				HandleConfirm: func(response string) error {
					return self.c.Helpers().Mode.SetFilteringAuthor(response)
				},
			})

			return nil
		},
		Tooltip: tooltip,
	})

	if self.c.Modes().Filtering.Active() {
		menuItems = append(menuItems, &types.MenuItem{
			Label:   self.c.Tr.ExitFilterMode,
			OnPress: self.c.Helpers().Mode.ClearFiltering,
		})
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.FilteringMenuTitle, Items: menuItems})
}
