package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TagsHelper struct {
	c             *HelperCommon
	commitsHelper *CommitsHelper
}

func NewTagsHelper(c *HelperCommon, commitsHelper *CommitsHelper) *TagsHelper {
	return &TagsHelper{
		c:             c,
		commitsHelper: commitsHelper,
	}
}

func (self *TagsHelper) OpenCreateTagPrompt(ref string, onCreate func()) error {
	doCreateTag := func(tagName string, description string, force bool) error {
		return self.c.WithWaitingStatus(self.c.Tr.CreatingTag, func(gocui.Task) error {
			if description != "" {
				self.c.LogAction(self.c.Tr.Actions.CreateAnnotatedTag)
				if err := self.c.Git().Tag.CreateAnnotated(tagName, ref, description, force); err != nil {
					return err
				}
			} else {
				self.c.LogAction(self.c.Tr.Actions.CreateLightweightTag)
				if err := self.c.Git().Tag.CreateLightweight(tagName, ref, force); err != nil {
					return err
				}
			}

			self.commitsHelper.OnCommitSuccess()

			return self.c.Refresh(types.RefreshOptions{
				Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS},
			})
		})
	}

	onConfirm := func(tagName string, description string) error {
		if self.c.Git().Tag.HasTag(tagName) {
			prompt := utils.ResolvePlaceholderString(
				self.c.Tr.ForceTagPrompt,
				map[string]string{
					"tagName":    tagName,
					"cancelKey":  self.c.UserConfig.Keybinding.Universal.Return,
					"confirmKey": self.c.UserConfig.Keybinding.Universal.Confirm,
				},
			)
			return self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.ForceTag,
				Prompt: prompt,
				HandleConfirm: func() error {
					return doCreateTag(tagName, description, true)
				},
			})
		} else {
			return doCreateTag(tagName, description, false)
		}
	}

	return self.commitsHelper.OpenCommitMessagePanel(
		&OpenCommitMessagePanelOpts{
			CommitIndex:      context.NoCommitIndex,
			InitialMessage:   "",
			SummaryTitle:     self.c.Tr.TagNameTitle,
			DescriptionTitle: self.c.Tr.TagMessageTitle,
			PreserveMessage:  false,
			OnConfirm:        onConfirm,
		},
	)
}
