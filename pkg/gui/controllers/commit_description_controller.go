package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: self.switchToCommitMessage,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.close,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.ConfirmInEditor),
			Handler: self.confirm,
		},
		{
			Key:     opts.GetKey(opts.Config.CommitMessage.CommitMenu),
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
			ViewName: self.Context().GetViewName(),
			Key:      gocui.MouseLeft,
			Handler:  self.onClick,
		},
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
	// Activate the description panel when the commit message panel is currently active
	if self.c.Context().Current().GetKey() == context.COMMIT_MESSAGE_CONTEXT_KEY {
		self.c.Context().Replace(self.c.Contexts().CommitDescription)
	}

	return nil
}
