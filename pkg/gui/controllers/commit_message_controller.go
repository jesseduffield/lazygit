package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &CommitMessageController{}

func NewCommitMessageController(
	c *ControllerCommon,
) *CommitMessageController {
	return &CommitMessageController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *CommitMessageController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.SubmitEditorText),
			Handler:     self.confirm,
			Description: self.c.Tr.Confirm,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.close,
			Description: self.c.Tr.Close,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevItem),
			Handler: self.handlePreviousCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextItem),
			Handler: self.handleNextCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: self.switchToCommitDescription,
		},
		{
			Key:     opts.GetKey(opts.Config.CommitMessage.CommitMenu),
			Handler: self.openCommitMenu,
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

func (self *CommitMessageController) context() *context.CommitMessageContext {
	return self.c.Contexts().CommitMessage
}

func (self *CommitMessageController) handlePreviousCommit() error {
	return self.handleCommitIndexChange(1)
}

func (self *CommitMessageController) handleNextCommit() error {
	if self.context().GetSelectedIndex() == context.NoCommitIndex {
		return nil
	}
	return self.handleCommitIndexChange(-1)
}

func (self *CommitMessageController) switchToCommitDescription() error {
	if err := self.c.PushContext(self.c.Contexts().CommitDescription); err != nil {
		return err
	}
	return nil
}

func (self *CommitMessageController) handleCommitIndexChange(value int) error {
	currentIndex := self.context().GetSelectedIndex()
	newIndex := currentIndex + value
	if newIndex == context.NoCommitIndex {
		self.context().SetSelectedIndex(newIndex)
		self.c.Helpers().Commits.SetMessageAndDescriptionInView(self.context().GetHistoryMessage())
		return nil
	} else if currentIndex == context.NoCommitIndex {
		self.context().SetHistoryMessage(self.c.Helpers().Commits.JoinCommitMessageAndUnwrappedDescription())
	}

	validCommit, err := self.setCommitMessageAtIndex(newIndex)
	if validCommit {
		self.context().SetSelectedIndex(newIndex)
	}
	return err
}

// returns true if the given index is for a valid commit
func (self *CommitMessageController) setCommitMessageAtIndex(index int) (bool, error) {
	commitMessage, err := self.c.Git().Commit.GetCommitMessageFromHistory(index)
	if err != nil {
		if err == git_commands.ErrInvalidCommitIndex {
			return false, nil
		}
		return false, self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}
	if self.c.UserConfig.Git.Commit.AutoWrapCommitMessage {
		commitMessage = helpers.TryRemoveHardLineBreaks(commitMessage, self.c.UserConfig.Git.Commit.AutoWrapWidth)
	}
	self.c.Helpers().Commits.UpdateCommitPanelView(commitMessage)
	return true, nil
}

func (self *CommitMessageController) confirm() error {
	return self.c.Helpers().Commits.HandleCommitConfirm()
}

func (self *CommitMessageController) close() error {
	return self.c.Helpers().Commits.CloseCommitMessagePanel()
}

func (self *CommitMessageController) openCommitMenu() error {
	authorSuggestion := self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc()
	return self.c.Helpers().Commits.OpenCommitMenu(authorSuggestion)
}
