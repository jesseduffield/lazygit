package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCommitConfirm() error {
	message := strings.TrimSpace(gui.Views.CommitMessage.TextArea.GetContent())
	gui.State.failedCommitMessage = message
	if message == "" {
		return gui.PopupHandler.ErrorMsg(gui.Tr.CommitWithoutMessageErr)
	}

	cmdObj := gui.Git.Commit.CommitCmdObj(message)
	gui.logAction(gui.Tr.Actions.Commit)

	_ = gui.returnFromContext()
	return gui.withGpgHandling(cmdObj, gui.Tr.CommittingStatus, func() error {
		gui.Views.CommitMessage.ClearTextArea()
		gui.State.failedCommitMessage = ""
		return nil
	})
}

func (gui *Gui) handleCommitClose() error {
	return gui.returnFromContext()
}

func (gui *Gui) handleCommitMessageFocused() error {
	message := utils.ResolvePlaceholderString(
		gui.Tr.CommitMessageConfirm,
		map[string]string{
			"keyBindClose":   gui.getKeyDisplay(gui.UserConfig.Keybinding.Universal.Return),
			"keyBindConfirm": gui.getKeyDisplay(gui.UserConfig.Keybinding.Universal.Confirm),
			"keyBindNewLine": gui.getKeyDisplay(gui.UserConfig.Keybinding.Universal.AppendNewline),
		},
	)

	return gui.renderString(gui.Views.Options, message)
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.UserConfig.Gui.CommitLength.Show {
		return
	}

	gui.Views.CommitMessage.Subtitle = gui.getBufferLength(gui.Views.CommitMessage)
}
