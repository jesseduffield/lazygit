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
		return gui.c.ErrorMsg(gui.c.Tr.CommitWithoutMessageErr)
	}

	cmdObj := gui.git.Commit.CommitCmdObj(message)
	gui.c.LogAction(gui.c.Tr.Actions.Commit)

	_ = gui.returnFromContext()
	return gui.withGpgHandling(cmdObj, gui.c.Tr.CommittingStatus, func() error {
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
		gui.c.Tr.CommitMessageConfirm,
		map[string]string{
			"keyBindClose":   gui.getKeyDisplay(gui.c.UserConfig.Keybinding.Universal.Return),
			"keyBindConfirm": gui.getKeyDisplay(gui.c.UserConfig.Keybinding.Universal.Confirm),
			"keyBindNewLine": gui.getKeyDisplay(gui.c.UserConfig.Keybinding.Universal.AppendNewline),
		},
	)

	gui.RenderCommitLength()

	return gui.renderString(gui.Views.Options, message)
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.c.UserConfig.Gui.CommitLength.Show {
		return
	}

	gui.Views.CommitMessage.Subtitle = gui.getBufferLength(gui.Views.CommitMessage)
}
