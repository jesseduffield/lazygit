package gui

import (
	"bufio"
	"strings"
	"fmt"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

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

func (gui *Gui) RenderCommitLength() {
	if !gui.c.UserConfig.Gui.CommitLength.Show {
		return
	}

	commitLength := getBufferLength(gui.Views.CommitMessage)
	gui.Views.CommitMessage.Subtitle = fmt.Sprintf(" %d ", commitLength)
	gui.checkCommitLengthWarning(commitLength)
}

func getBufferLength(view *gocui.View) int {
	return strings.Count(view.TextArea.GetContent(), "") - 1
}

func (gui *Gui) checkCommitLengthWarning(commitLength int) {
	if gui.c.UserConfig.Gui.CommitLength.WarningThreshold == 0 {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(gui.Views.CommitMessage.TextArea.GetContent()))

	for scanner.Scan() {
		if len(scanner.Text()) > gui.c.UserConfig.Gui.CommitLength.WarningThreshold {
			gui.Views.CommitMessage.FgColor = gocui.ColorRed
		} else {
			gui.Views.CommitMessage.FgColor = theme.GocuiDefaultTextColor
		}
	}
}