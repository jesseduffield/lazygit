package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
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
				return gui.withResetButton(
					fmt.Sprintf(
						"%s %s",
						gui.Tr.LcShowingGitDiff,
						"git diff "+gui.diffStr(),
					),
					style.FgMagenta,
				)
			},
			reset: gui.exitDiffMode,
		},
		{
			isActive: gui.Git.Patch.PatchManager.Active,
			description: func() string {
				return gui.withResetButton(gui.Tr.LcBuildingPatch, style.FgYellow.SetBold())
			},
			reset: gui.handleResetPatch,
		},
		{
			isActive: gui.State.Modes.Filtering.Active,
			description: func() string {
				return gui.withResetButton(
					fmt.Sprintf(
						"%s '%s'",
						gui.Tr.LcFilteringBy,
						gui.State.Modes.Filtering.GetPath(),
					),
					style.FgRed,
				)
			},
			reset: gui.exitFilterMode,
		},
		{
			isActive: gui.State.Modes.CherryPicking.Active,
			description: func() string {
				return gui.withResetButton(
					fmt.Sprintf(
						"%d commits copied",
						len(gui.State.Modes.CherryPicking.CherryPickedCommits),
					),
					style.FgCyan,
				)
			},
			reset: gui.exitCherryPickingMode,
		},
		{
			isActive: func() bool {
				return gui.Git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE
			},
			description: func() string {
				workingTreeState := gui.Git.Status.WorkingTreeState()
				return gui.withResetButton(
					formatWorkingTreeState(workingTreeState), style.FgYellow,
				)
			},
			reset: gui.abortMergeOrRebaseWithConfirm,
		},
		{
			isActive: func() bool {
				return gui.State.BisectInfo.Started()
			},
			description: func() string {
				return gui.withResetButton("bisecting", style.FgGreen)
			},
			reset: gui.resetBisect,
		},
	}
}

func (gui *Gui) withResetButton(content string, textStyle style.TextStyle) string {
	return textStyle.Sprintf(
		"%s %s",
		content,
		style.AttrUnderline.Sprint(gui.Tr.ResetInParentheses),
	)
}
