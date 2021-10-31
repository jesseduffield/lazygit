package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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
				return style.FgMagenta.Sprintf(
					"%s %s %s",
					gui.Tr.LcShowingGitDiff,
					"git diff "+gui.diffStr(),
					style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
				)
			},
			reset: gui.exitDiffMode,
		},
		{
			isActive: gui.GitCommand.PatchManager.Active,
			description: func() string {
				return style.FgYellow.SetBold().Sprintf(
					"%s %s",
					gui.Tr.LcBuildingPatch,
					style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
				)
			},
			reset: gui.handleResetPatch,
		},
		{
			isActive: gui.State.Modes.Filtering.Active,
			description: func() string {
				return style.FgRed.SetBold().Sprintf(
					"%s '%s' %s",
					gui.Tr.LcFilteringBy,
					gui.State.Modes.Filtering.GetPath(),
					style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
				)
			},
			reset: gui.exitFilterMode,
		},
		{
			isActive: gui.State.Modes.CherryPicking.Active,
			description: func() string {
				return style.FgCyan.Sprintf(
					"%d commits copied %s",
					len(gui.State.Modes.CherryPicking.CherryPickedCommits),
					style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
				)
			},
			reset: gui.exitCherryPickingMode,
		},
		{
			isActive: func() bool {
				return gui.GitCommand.WorkingTreeState() != commands.REBASE_MODE_NORMAL
			},
			description: func() string {
				workingTreeState := gui.GitCommand.WorkingTreeState()
				return style.FgYellow.Sprintf(
					"%s %s",
					workingTreeState,
					style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
				)
			},
			reset: gui.abortMergeOrRebaseWithConfirm,
		},
	}
}
