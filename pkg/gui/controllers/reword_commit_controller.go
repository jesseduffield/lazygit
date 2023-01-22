package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RewordCommitController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &RewordCommitController{}

func NewRewordCommitController(
	common *controllerCommon,
) *RewordCommitController {
	return &RewordCommitController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *RewordCommitController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.SubmitEditorText),
			Handler: self.confirm,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.close,
		},
	}

	return bindings
}

func (self *RewordCommitController) Context() types.Context {
	return self.context()
}

// this method is pointless in this context but I'm keeping it consistent
// with other contexts so that when generics arrive it's easier to refactor
func (self *RewordCommitController) context() types.Context {
	return self.contexts.RewordCommitMessage
}

func (self *RewordCommitController) confirm() error {
	message := self.context().GetView().TextArea.GetContent()

	if message == "" {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}

	self.c.LogAction(self.c.Tr.Actions.RewordCommit)
	if err := self.git.Rebase.RewordCommit(
		self.model.Commits,
		self.contexts.LocalCommits.GetSelectedLineIdx(),
		message,
	); err != nil {
		return self.c.Error(err)
	}

	_ = self.c.PopContext()
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *RewordCommitController) close() error {
	return self.c.PopContext()
}
