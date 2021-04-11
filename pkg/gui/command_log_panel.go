package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) GetOnRunCommand() func(entry oscommands.CmdLogEntry) {
	// closing over this so that nobody else can modify it
	currentSpan := ""

	return func(entry oscommands.CmdLogEntry) {
		if gui.Views.Extras == nil {
			return
		}

		gui.Views.Extras.Autoscroll = true

		if entry.GetSpan() != currentSpan {
			fmt.Fprint(gui.Views.Extras, "\n"+utils.ColoredString(entry.GetSpan(), color.FgYellow))
			currentSpan = entry.GetSpan()
		}

		clrAttr := theme.DefaultTextColor
		if !entry.GetCommandLine() {
			clrAttr = color.FgMagenta
		}
		gui.CmdLog = append(gui.CmdLog, entry.GetCmdStr())
		indentedCmdStr := "  " + strings.Replace(entry.GetCmdStr(), "\n", "\n  ", -1)
		fmt.Fprint(gui.Views.Extras, "\n"+utils.ColoredString(indentedCmdStr, clrAttr))
	}
}
