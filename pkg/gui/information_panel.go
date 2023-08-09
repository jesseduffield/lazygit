package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

	var title, url string

	// if we're not in an active mode we show the donate button
	if cx <= runewidth.StringWidth(gui.c.Tr.Donate) {
		url = constants.Links.Donate
		title = gui.c.Tr.Donate
	} else if cx <= runewidth.StringWidth(gui.c.Tr.Donate)+1+runewidth.StringWidth(gui.c.Tr.AskQuestion) {
		url = constants.Links.Discussions
		title = gui.c.Tr.AskQuestion
	}
	err := gui.os.OpenLink(url)
	if err != nil {
		// Opening the link via the OS failed for some reason. (For example, this
		// can happen if the `os.openLink` config key references a command that
		// doesn't exist, or that errors when called.)
		//
		// In that case, rather than crash the app, fall back to simply showing a
		// dialog asking the user to visit the URL.
		placeholders := map[string]string{"url": url}
		message := utils.ResolvePlaceholderString(gui.c.Tr.PleaseGoToURL, placeholders)
		return gui.c.Alert(title, message)
	}

	return nil
}
