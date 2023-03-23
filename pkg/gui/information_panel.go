package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/mattn/go-runewidth"
)

func (gui *Gui) informationStr() string {
	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		return activeMode.Description()
	}

	if gui.g.Mouse {
		donate := style.FgMagenta.SetUnderline().Sprint(gui.c.Tr.Donate)
		askQuestion := style.FgYellow.SetUnderline().Sprint(gui.c.Tr.AskQuestion)
		return fmt.Sprintf("%s %s %s", donate, askQuestion, gui.Config.GetVersion())
	} else {
		return gui.Config.GetVersion()
	}
}

func (gui *Gui) handleInfoClick() error {
	if !gui.g.Mouse {
		return nil
	}

	view := gui.Views.Information

	cx, _ := view.Cursor()
	width, _ := view.Size()

	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		if width-cx > runewidth.StringWidth(gui.c.Tr.ResetInParentheses) {
			return nil
		}
		return activeMode.Reset()
	}

	// if we're not in an active mode we show the donate button
	if cx <= runewidth.StringWidth(gui.c.Tr.Donate) {
		return gui.os.OpenLink(constants.Links.Donate)
	} else if cx <= runewidth.StringWidth(gui.c.Tr.Donate)+1+runewidth.StringWidth(gui.c.Tr.AskQuestion) {
		return gui.os.OpenLink(constants.Links.Discussions)
	}
	return nil
}
