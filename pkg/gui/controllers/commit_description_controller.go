package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitDescriptionController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &CommitMessageController{}

func NewCommitDescriptionController(
	common *controllerCommon,
) *CommitDescriptionController {
	return &CommitDescriptionController{
		baseController:   baseController{},
		controllerCommon: common,
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
	}

	return bindings
}

func (self *CommitDescriptionController) Context() types.Context {
	return self.context()
}

func (self *CommitDescriptionController) context() types.Context {
	return self.contexts.CommitMessage
}

func (self *CommitDescriptionController) switchToCommitMessage() error {
	return self.c.PushContext(self.contexts.CommitMessage)
}

func (self *CommitDescriptionController) close() error {
	return self.helpers.Commits.CloseCommitMessagePanel()
}

func (self *CommitDescriptionController) confirm() error {
	return self.helpers.Commits.HandleCommitConfirm()
}
