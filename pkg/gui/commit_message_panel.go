package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCommitMessageFocused() error {
	message := utils.ResolvePlaceholderString(
		gui.c.Tr.CommitMessageConfirm,
		map[string]string{
			"keyBindClose":   keybindings.Label(gui.c.UserConfig.Keybinding.Universal.Return),
			"keyBindConfirm": keybindings.Label(gui.c.UserConfig.Keybinding.Universal.Confirm),
			"keyBindNewLine": keybindings.Label(gui.c.UserConfig.Keybinding.Universal.AppendNewline),
		},
	)

	gui.RenderCommitLength()

	return gui.renderString(gui.Views.Options, message)
}

func (gui *Gui) RenderCommitLength() {
	if !gui.c.UserConfig.Gui.CommitLength.Show {
		return
	}

	gui.Views.CommitMessage.Subtitle = getBufferLength(gui.Views.CommitMessage)
}

func getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}
