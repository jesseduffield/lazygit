package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type modeStatus struct {
	isActive    func() bool
	description func() string
	reset       func() error
}

func (gui *Gui) modeStatuses() []modeStatus {
	return []modeStatus{
		{
			isActive: gui.State.Modes.Diffing.Active,
			description: func() string {
				return utils.ColoredString(
					fmt.Sprintf("%s %s %s", gui.Tr.LcShowingGitDiff, "git diff "+gui.diffStr(), utils.ColoredString(gui.Tr.ResetInParentheses, color.Underline)),
					color.FgMagenta,
				)
			},
			reset: gui.exitDiffMode,
		},
		{
			isActive: gui.State.Modes.Filtering.Active,
			description: func() string {
				return utils.ColoredString(
					fmt.Sprintf("%s '%s' %s", gui.Tr.LcFilteringBy, gui.State.Modes.Filtering.Path, utils.ColoredString(gui.Tr.ResetInParentheses, color.Underline)),
					color.FgRed,
					color.Bold,
				)
			},
			reset: gui.exitFilterMode,
		},
		{
			isActive: gui.GitCommand.PatchManager.Active,
			description: func() string {
				return utils.ColoredString(
					fmt.Sprintf("%s %s", gui.Tr.LcBuildingPatch, utils.ColoredString(gui.Tr.ResetInParentheses, color.Underline)),
					color.FgYellow,
					color.Bold,
				)
			},
			reset: gui.handleResetPatch,
		},
		{
			isActive: gui.State.Modes.CherryPicking.Active,
			description: func() string {
				return utils.ColoredString(
					fmt.Sprintf("%d commits copied %s", len(gui.State.Modes.CherryPicking.CherryPickedCommits), utils.ColoredString(gui.Tr.ResetInParentheses, color.Underline)),
					color.FgCyan,
				)
			},
			reset: gui.exitCherryPickingMode,
		},
	}
}
