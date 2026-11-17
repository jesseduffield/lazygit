package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) informationStr() string {
	if activeMode, ok := gui.helpers.Mode.GetActiveMode(); ok {
		return activeMode.InfoLabel()
	}

	version := displayVersion(gui.Config.GetVersion())
	if note := gui.helpers.Update.Note(); note != "" {
		version = note + " " + version
	}

	return version
}

// displayVersion formats a release version as "v1.2.3", the way the update
// note and the docs links in the status panel show versions, no matter how the
// build set it. Other version strings, e.g. those of dev builds, stay as they
// are.
func displayVersion(versionStr string) string {
	if version, err := types.ParseVersionNumber(versionStr); err == nil {
		return version.String()
	}
	return versionStr
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
