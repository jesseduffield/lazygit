package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) informationStr() string {
	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		return activeMode.Description()
	}

	if gui.g.Mouse {
		donate := style.FgMagenta.Sprint(style.PrintHyperlink(gui.c.Tr.Donate, constants.Links.Donate))
		askQuestion := style.FgYellow.Sprint(style.PrintHyperlink(gui.c.Tr.AskQuestion, constants.Links.Discussions))
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
	width := view.Width()

	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		if width-cx > utils.StringWidth(gui.c.Tr.ResetInParentheses) {
			return nil
		}
		return activeMode.Reset()
	}

	return nil
}
