package controllers

import (
	"errors"

	"github.com/jesseduffield/gocui"
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
			Handler: self.handleTogglePanel,
		},
		{
			Key:     opts.GetKey(opts.Config.CommitMessage.CommitMenu),
			Handler: self.openCommitMenu,
		},
	}

	return bindings
}

func (self *CommitMessageController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.Context().GetViewName(),
			FocusedView: self.c.Contexts().CommitDescription.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClick,
		},
	}
}

func (self *CommitMessageController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		self.c.Views().CommitDescription.Footer = ""
	}
}

func (self *CommitMessageController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.context().RenderSubtitle()
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
	self.c.Context().Replace(self.c.Contexts().CommitDescription)
	return nil
}

func (self *CommitMessageController) handleTogglePanel() error {
	// The default keybinding for this action is "<tab>", which means that we
	// also get here when pasting multi-line text that contains tabs. In that
	// case we don't want to toggle the panel, but insert the tab as a character
	// (somehow, see below).
	//
	// Only do this if the TogglePanel command is actually mapped to "<tab>"
	// (the default). If it's not, we can only hope that it's mapped to some
	// ctrl key or fn key, which is unlikely to occur in pasted text. And if
	// they mapped some *other* command to "<tab>", then we're totally out of
	// luck.
	if self.c.GocuiGui().IsPasting && self.c.UserConfig().Keybinding.Universal.TogglePanel == "<tab>" {
		// It is unlikely that a pasted commit message contains a tab in the
		// subject line, so it shouldn't matter too much how we handle it.
		// Simply insert 4 spaces instead; all that matters is that we don't
		// switch to the description panel.
		view := self.context().GetView()
		for range 4 {
			view.Editor.Edit(view, gocui.KeySpace, ' ', 0)
		}
		return nil
	}

	return self.switchToCommitDescription()
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
		if errors.Is(err, git_commands.ErrInvalidCommitIndex) {
			return false, nil
		}
		return false, errors.New(self.c.Tr.CommitWithoutMessageErr)
	}
	if self.c.UserConfig().Git.Commit.AutoWrapCommitMessage {
		commitMessage = helpers.TryRemoveHardLineBreaks(commitMessage, self.c.UserConfig().Git.Commit.AutoWrapWidth)
	}
	self.c.Helpers().Commits.UpdateCommitPanelView(commitMessage)
	return true, nil
}

func (self *CommitMessageController) confirm() error {
	// The default keybinding for this action is "<enter>", which means that we
	// also get here when pasting multi-line text that contains newlines. In
	// that case we don't want to confirm the commit, but switch to the
	// description panel instead so that the rest of the pasted text goes there.
	//
	// Only do this if the SubmitEditorText command is actually mapped to
	// "<enter>" (the default). If it's not, we can only hope that it's mapped
	// to some ctrl key or fn key, which is unlikely to occur in pasted text.
	// And if they mapped some *other* command to "<enter>", then we're totally
	// out of luck.
	if self.c.GocuiGui().IsPasting && self.c.UserConfig().Keybinding.Universal.SubmitEditorText == "<enter>" {
		return self.switchToCommitDescription()
	}

	return self.c.Helpers().Commits.HandleCommitConfirm()
}

func (self *CommitMessageController) close() error {
	self.c.Helpers().Commits.CloseCommitMessagePanel()
	return nil
}

func (self *CommitMessageController) openCommitMenu() error {
	authorSuggestion := self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc()
	return self.c.Helpers().Commits.OpenCommitMenu(authorSuggestion)
}

func (self *CommitMessageController) onClick(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Replace(self.c.Contexts().CommitMessage)
	return nil
}
