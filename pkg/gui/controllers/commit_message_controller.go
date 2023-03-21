package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageController struct {
	baseController
	*controllerCommon

	getCommitMessage func() string
	onCommitAttempt  func(message string)
	onCommitSuccess  func()
}

var _ types.IController = &CommitMessageController{}

func NewCommitMessageController(
	common *controllerCommon,
	getCommitMessage func() string,
	onCommitAttempt func(message string),
	onCommitSuccess func(),
) *CommitMessageController {
	return &CommitMessageController{
		baseController:   baseController{},
		controllerCommon: common,

		getCommitMessage: getCommitMessage,
		onCommitAttempt:  onCommitAttempt,
		onCommitSuccess:  onCommitSuccess,
	}
}

// TODO: merge that commit panel PR because we're not currently showing how to add a newline as it's
// handled by the editor func rather than by the controller here.
func (self *CommitMessageController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.SubmitEditorText),
			Handler:     self.confirm,
			Description: self.c.Tr.LcConfirm,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.close,
			Description: self.c.Tr.LcClose,
		},
	}

	return bindings
}

func (self *CommitMessageController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		self.context().RenderCommitLength()
		return nil
	}
}

func (self *CommitMessageController) Context() types.Context {
	return self.context()
}

// this method is pointless in this context but I'm keeping it consistent
// with other contexts so that when generics arrive it's easier to refactor
func (self *CommitMessageController) context() *context.CommitMessageContext {
	return self.contexts.CommitMessage
}

func (self *CommitMessageController) confirm() error {
	message := self.getCommitMessage()
	self.onCommitAttempt(message)

	if message == "" {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}

	cmdObj := self.git.Commit.CommitCmdObj(message)
	self.c.LogAction(self.c.Tr.Actions.Commit)

	_ = self.c.PopContext()
	return self.helpers.GPG.WithGpgHandling(cmdObj, self.c.Tr.CommittingStatus, func() error {
		self.onCommitSuccess()
		return nil
	})
}

func (self *CommitMessageController) close() error {
	return self.c.PopContext()
}
