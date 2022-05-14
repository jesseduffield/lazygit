package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

	gui.Views.CommitMessage.Subtitle = getBufferLength(gui.Views.CommitMessage)
}

func getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, true)

	// This function is called again on refresh as part of the general resize popup call,
	// but we need to call it here so that when we go to render the text area it's not
	// considered out of bounds to add a newline, meaning we can avoid unnecessary scrolling.
	err := gui.resizePopupPanel(v, v.TextArea.GetContent())
	if err != nil {
		gui.c.Log.Error(err)
	}
	v.RenderTextArea()
	gui.RenderCommitLength()

	return matched
}
