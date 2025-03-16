package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TagsHelper struct {
	c             *HelperCommon
	commitsHelper *CommitsHelper
	gpg           *GpgHelper
}

func NewTagsHelper(c *HelperCommon, commitsHelper *CommitsHelper, gpg *GpgHelper) *TagsHelper {
	return &TagsHelper{
		c:             c,
		commitsHelper: commitsHelper,
		gpg:           gpg,
	}
}

func (self *TagsHelper) OpenCreateTagPrompt(ref string, onCreate func()) error {
	doCreateTag := func(tagName string, description string, force bool) error {
		var command oscommands.ICmdObj
		if description != "" || self.c.Git().Config.GetGpgTagSign() {
			self.c.LogAction(self.c.Tr.Actions.CreateAnnotatedTag)
			command = self.c.Git().Tag.CreateAnnotatedObj(tagName, ref, description, force)
		} else {
			self.c.LogAction(self.c.Tr.Actions.CreateLightweightTag)
			command = self.c.Git().Tag.CreateLightweightObj(tagName, ref, force)
		}

		return self.gpg.WithGpgHandling(command, git_commands.TagGpgSign, self.c.Tr.CreatingTag, func() error {
			self.commitsHelper.OnCommitSuccess()
			return nil
		}, []types.RefreshableView{types.COMMITS, types.TAGS})
	}

	onConfirm := func(tagName string, description string) error {
		if self.c.Git().Tag.HasTag(tagName) {
			prompt := utils.ResolvePlaceholderString(
				self.c.Tr.ForceTagPrompt,
				map[string]string{
					"tagName":    tagName,
					"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return,
					"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm,
				},
			)
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.ForceTag,
				Prompt: prompt,
				HandleConfirm: func() error {
					return doCreateTag(tagName, description, true)
				},
			})

			return nil
		}

		return doCreateTag(tagName, description, false)
	}

	self.commitsHelper.OpenCommitMessagePanel(
		&OpenCommitMessagePanelOpts{
			CommitIndex:      context.NoCommitIndex,
			InitialMessage:   "",
			SummaryTitle:     self.c.Tr.TagNameTitle,
			DescriptionTitle: self.c.Tr.TagMessageTitle,
			PreserveMessage:  false,
			OnConfirm:        onConfirm,
		},
	)

	return nil
}
