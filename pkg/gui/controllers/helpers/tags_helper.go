package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Helper structs are for defining functionality that could be used by multiple contexts.
// For example, here we have a CreateTagMenu which is applicable to both the tags context
// and the commits context.

type TagsHelper struct {
	c   *types.HelperCommon
	git *commands.GitCommand
}

func NewTagsHelper(c *types.HelperCommon, git *commands.GitCommand) *TagsHelper {
	return &TagsHelper{
		c:   c,
		git: git,
	}
}

func (self *TagsHelper) CreateTagMenu(commitSha string, onCreate func()) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.TagMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.LcLightweightTag,
				OnPress: func() error {
					return self.handleCreateLightweightTag(commitSha, onCreate)
				},
			},
			{
				Label: self.c.Tr.LcAnnotatedTag,
				OnPress: func() error {
					return self.handleCreateAnnotatedTag(commitSha, onCreate)
				},
			},
		},
	})
}

func (self *TagsHelper) afterTagCreate(onCreate func()) error {
	onCreate()
	return self.c.Refresh(types.RefreshOptions{
		Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS},
	})
}

func (self *TagsHelper) handleCreateAnnotatedTag(commitSha string, onCreate func()) error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.TagNameTitle,
		HandleConfirm: func(tagName string) error {
			return self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.TagMessageTitle,
				HandleConfirm: func(msg string) error {
					self.c.LogAction(self.c.Tr.Actions.CreateAnnotatedTag)
					if err := self.git.Tag.CreateAnnotated(tagName, commitSha, msg); err != nil {
						return self.c.Error(err)
					}
					return self.afterTagCreate(onCreate)
				},
			})
		},
	})
}

func (self *TagsHelper) handleCreateLightweightTag(commitSha string, onCreate func()) error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.TagNameTitle,
		HandleConfirm: func(tagName string) error {
			self.c.LogAction(self.c.Tr.Actions.CreateLightweightTag)
			if err := self.git.Tag.CreateLightweight(tagName, commitSha); err != nil {
				return self.c.Error(err)
			}
			return self.afterTagCreate(onCreate)
		},
	})
}
