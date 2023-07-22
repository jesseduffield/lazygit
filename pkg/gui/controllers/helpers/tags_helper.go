package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
	onConfirm := func(tagName string, description string) error {
		return self.c.WithWaitingStatus(self.c.Tr.CreatingTag, func(gocui.Task) error {
			if description != "" {
				self.c.LogAction(self.c.Tr.Actions.CreateAnnotatedTag)
				if err := self.c.Git().Tag.CreateAnnotated(tagName, ref, description); err != nil {
					return self.c.Error(err)
				}
			} else {
				self.c.LogAction(self.c.Tr.Actions.CreateLightweightTag)
				if err := self.c.Git().Tag.CreateLightweight(tagName, ref); err != nil {
					return self.c.Error(err)
				}
			}

			self.commitsHelper.OnCommitSuccess()

			return self.c.Refresh(types.RefreshOptions{
				Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS},
			})
		})
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
