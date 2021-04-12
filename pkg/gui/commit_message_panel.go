package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCommitConfirm() error {
	message := gui.trimmedContent(gui.Views.CommitMessage)
	if message == "" {
		return gui.createErrorPanel(gui.Tr.CommitWithoutMessageErr)
	}
	flags := ""
	skipHookPrefix := gui.Config.GetUserConfig().Git.SkipHookPrefix
	if skipHookPrefix != "" && strings.HasPrefix(message, skipHookPrefix) {
		flags = "--no-verify"
	}

	cmdStr := gui.GitCommand.CommitCmdStr(message, flags)
	gui.OnRunCommand(oscommands.NewCmdLogEntry(cmdStr, gui.Tr.Spans.Commit, true))
	return gui.withGpgHandling(cmdStr, gui.Tr.CommittingStatus, func() error {
		_ = gui.returnFromContext()
		gui.clearEditorView(gui.Views.CommitMessage)
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
			"keyBindClose":   gui.getKeyDisplay(gui.Config.GetUserConfig().Keybinding.Universal.Return),
			"keyBindConfirm": gui.getKeyDisplay(gui.Config.GetUserConfig().Keybinding.Universal.Confirm),
			"keyBindNewLine": gui.getKeyDisplay(gui.Config.GetUserConfig().Keybinding.Universal.AppendNewline),
		},
	)

	gui.renderString(gui.Views.Options, message)
	return nil
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.Buffer(), "")-1) + " "
}

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().Gui.CommitLength.Show {
		return
	}

	gui.Views.CommitMessage.Subtitle = gui.getBufferLength(gui.Views.CommitMessage)
}
