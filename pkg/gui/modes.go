package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type modeStatus struct {
	isActive    func() bool
	description func() string
	onReset     func() error
}

func (gui *Gui) modeStatuses() []modeStatus {
	return []modeStatus{
		{
			isActive: gui.State.Modes.Diffing.Active,
			description: func() string {
				return utils.ColoredString(fmt.Sprintf("%s %s %s", gui.Tr.SLocalize("showingGitDiff"), "git diff "+gui.diffStr(), utils.ColoredString(gui.Tr.SLocalize("(reset)"), color.Underline)), color.FgMagenta)
			},
			onReset: gui.exitDiffMode,
		},
		{
			isActive: gui.State.Modes.Filtering.Active,
			description: func() string {
				return utils.ColoredString(fmt.Sprintf("%s '%s' %s", gui.Tr.SLocalize("filteringBy"), gui.State.Modes.Filtering.Path, utils.ColoredString(gui.Tr.SLocalize("(reset)"), color.Underline)), color.FgRed, color.Bold)
			},
			onReset: gui.exitFilterMode,
		},
		{
			isActive: gui.GitCommand.PatchManager.Active,
			description: func() string {
				return utils.ColoredString(fmt.Sprintf("%s %s", gui.Tr.SLocalize("buildingPatch"), utils.ColoredString(gui.Tr.SLocalize("(reset)"), color.Underline)), color.FgYellow, color.Bold)
			},
			onReset: gui.handleResetPatch,
		},
		{
			isActive: gui.State.Modes.CherryPicking.Active,
			description: func() string {
				return utils.ColoredString(fmt.Sprintf("%d commits copied", len(gui.State.Modes.CherryPicking.CherryPickedCommits)), color.FgCyan)
			},
			onReset: gui.exitCherryPickingMode,
		},
	}
}
