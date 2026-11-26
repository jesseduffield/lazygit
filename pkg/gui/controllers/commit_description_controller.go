package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type CommitDescriptionController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &CommitMessageController{}

func NewCommitDescriptionController(
	c *ControllerCommon,
) *CommitDescriptionController {
	return &CommitDescriptionController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *CommitDescriptionController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:    opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler: self.switchToCommitMessage,
		},
		{
			Keys:    opts.GetKeys(opts.Config.Universal.Return),
			Handler: self.close,
		},
		{
			Keys:    opts.GetKeys(opts.Config.Universal.ConfirmInEditor),
			Handler: self.confirm,
		},
		{
			Keys:    opts.GetKeys(opts.Config.CommitMessage.CommitMenu),
			Handler: self.openCommitMenu,
		},
	}

	return bindings
}

func (self *CommitDescriptionController) Context() types.Context {
	return self.c.Contexts().CommitDescription
}

func (self *CommitDescriptionController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.Context().GetViewName(),
			FocusedView: self.c.Contexts().CommitMessage.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClick,
		},
	}
}

func (self *CommitDescriptionController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		footer := ""
		keys := self.c.UserConfig().Keybinding.Universal.ConfirmInEditor
		if len(keys) > 0 {
			footer = utils.ResolvePlaceholderString(self.c.Tr.CommitDescriptionFooter,
				map[string]string{
					"confirmInEditorKeybinding": keys.String(),
				})
		}
		self.c.Views().CommitDescription.Footer = footer
	}
}

func (self *CommitDescriptionController) GetOnQuit() func() {
	return func() {
		self.c.Helpers().Commits.PreserveCommitMessage()
	}
}

func (self *CommitDescriptionController) switchToCommitMessage() error {
	self.c.Context().Replace(self.c.Contexts().CommitMessage)
	return nil
}

func (self *CommitDescriptionController) close() error {
	self.c.Helpers().Commits.CloseCommitMessagePanel()
	return nil
}

func (self *CommitDescriptionController) confirm() error {
	return self.c.Helpers().Commits.HandleCommitConfirm()
}

func (self *CommitDescriptionController) openCommitMenu() error {
	authorSuggestion := self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc()
	return self.c.Helpers().Commits.OpenCommitMenu(authorSuggestion)
}

func (self *CommitDescriptionController) onClick(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Replace(self.c.Contexts().CommitDescription)
	return nil
}
