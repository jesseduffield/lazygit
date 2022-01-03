package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCommitConfirm() error {
	message := strings.TrimSpace(gui.Views.CommitMessage.TextArea.GetContent())
	if message == "" {
		return gui.createErrorPanel(gui.Tr.CommitWithoutMessageErr)
	}
	flags := []string{}
	skipHookPrefix := gui.UserConfig.Git.SkipHookPrefix
	if skipHookPrefix != "" && strings.HasPrefix(message, skipHookPrefix) {
		flags = append(flags, "--no-verify")
	}

	if gui.UserConfig.Git.Commit.SignOff {
		flags = append(flags, "--signoff")
	}

	cmdObj := gui.GitCommand.CommitCmdObj(message, strings.Join(flags, " "))
	gui.OnRunCommand(oscommands.NewCmdLogEntry(cmdObj.ToString(), gui.Tr.Spans.Commit, true))

	_ = gui.returnFromContext()
	return gui.withGpgHandling(cmdObj, gui.Tr.CommittingStatus, func() error {
		gui.Views.CommitMessage.ClearTextArea()
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

	gui.renderString(gui.Views.Options, message)
	return nil
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
