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
	onConfirm := func(tagName string, description string) error {
		prompt := utils.ResolvePlaceholderString(
			self.c.Tr.ForceTagPrompt,
			map[string]string{
				"tagName":    tagName,
				"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return,
				"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm,
			},
		)
		force := self.c.Git().Tag.HasTag(tagName)
		return self.c.ConfirmIf(force, types.ConfirmOpts{
			Title:  self.c.Tr.ForceTag,
			Prompt: prompt,
			HandleConfirm: func() error {
				var command *oscommands.CmdObj
				if description != "" || self.c.Git().Config.GetGpgTagSign() {
					self.c.LogAction(self.c.Tr.Actions.CreateAnnotatedTag)
					command = self.c.Git().Tag.CreateAnnotatedObj(tagName, ref, description, force)
				} else {
					self.c.LogAction(self.c.Tr.Actions.CreateLightweightTag)
					command = self.c.Git().Tag.CreateLightweightObj(tagName, ref, force)
				}

				return self.gpg.WithGpgHandling(command, git_commands.TagGpgSign, self.c.Tr.CreatingTag, func() error {
					return nil
				}, []types.RefreshableView{types.COMMITS, types.TAGS})
			},
		})
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
