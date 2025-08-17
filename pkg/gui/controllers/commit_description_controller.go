package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
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
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: self.handleTogglePanel,
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
			Key:     opts.GetKey(opts.Config.Universal.ConfirmInEditorAlt),
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
		if self.c.UserConfig().Keybinding.Universal.ConfirmInEditor != "<disabled>" || self.c.UserConfig().Keybinding.Universal.ConfirmInEditorAlt != "<disabled>" {
			if self.c.UserConfig().Keybinding.Universal.ConfirmInEditor == "<disabled>" {
				footer = utils.ResolvePlaceholderString(self.c.Tr.CommitDescriptionFooter,
					map[string]string{
						"confirmInEditorKeybinding": keybindings.Label(self.c.UserConfig().Keybinding.Universal.ConfirmInEditorAlt),
					})
			} else if self.c.UserConfig().Keybinding.Universal.ConfirmInEditorAlt == "<disabled>" {
				footer = utils.ResolvePlaceholderString(self.c.Tr.CommitDescriptionFooter,
					map[string]string{
						"confirmInEditorKeybinding": keybindings.Label(self.c.UserConfig().Keybinding.Universal.ConfirmInEditor),
					})
			} else {
				footer = utils.ResolvePlaceholderString(self.c.Tr.CommitDescriptionFooterTwoBindings,
					map[string]string{
						"confirmInEditorKeybinding1": keybindings.Label(self.c.UserConfig().Keybinding.Universal.ConfirmInEditor),
						"confirmInEditorKeybinding2": keybindings.Label(self.c.UserConfig().Keybinding.Universal.ConfirmInEditorAlt),
					})
			}
		}
		self.c.Views().CommitDescription.Footer = footer
	}
}

func (self *CommitDescriptionController) switchToCommitMessage() error {
	self.c.Context().Replace(self.c.Contexts().CommitMessage)
	return nil
}

func (self *CommitDescriptionController) handleTogglePanel() error {
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
		// Handling tabs in pasted commit messages is not optimal, but hopefully
		// good enough for now. We simply insert 4 spaces without worrying about
		// column alignment. This works well enough for leading indentation,
		// which is common in pasted code snippets.
		view := self.Context().GetView()
		for range 4 {
			view.Editor.Edit(view, gocui.KeySpace, ' ', 0)
		}
		return nil
	}

	return self.switchToCommitMessage()
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
