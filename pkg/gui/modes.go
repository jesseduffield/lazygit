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
						gui.c.Tr.LcShowingGitDiff,
						"git diff "+gui.diffStr(),
					),
					style.FgMagenta,
				)
			},
			reset: gui.exitDiffMode,
		},
		{
			isActive: gui.git.Patch.PatchManager.Active,
			description: func() string {
				return gui.withResetButton(gui.c.Tr.LcBuildingPatch, style.FgYellow.SetBold())
			},
			reset: gui.helpers.PatchBuilding.Reset,
		},
		{
			isActive: gui.State.Modes.Filtering.Active,
			description: func() string {
				return gui.withResetButton(
					fmt.Sprintf(
						"%s '%s'",
						gui.c.Tr.LcFilteringBy,
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
				copiedCount := len(gui.State.Modes.CherryPicking.CherryPickedCommits)
				text := gui.c.Tr.LcCommitsCopied
				if copiedCount == 1 {
					text = gui.c.Tr.LcCommitCopied
				}

				return gui.withResetButton(
					fmt.Sprintf(
						"%d %s",
						copiedCount,
						text,
					),
					style.FgCyan,
				)
			},
			reset: gui.helpers.CherryPick.Reset,
		},
		{
			isActive: func() bool {
				return gui.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE
			},
			description: func() string {
				workingTreeState := gui.git.Status.WorkingTreeState()
				return gui.withResetButton(
					formatWorkingTreeState(workingTreeState), style.FgYellow,
				)
			},
			reset: gui.helpers.MergeAndRebase.AbortMergeOrRebaseWithConfirm,
		},
		{
			isActive: func() bool {
				return gui.State.Model.BisectInfo.Started()
			},
			description: func() string {
				return gui.withResetButton("bisecting", style.FgGreen)
			},
			reset: gui.helpers.Bisect.Reset,
		},
	}
}

func (gui *Gui) withResetButton(content string, textStyle style.TextStyle) string {
	return textStyle.Sprintf(
		"%s %s",
		content,
		style.AttrUnderline.Sprint(gui.c.Tr.ResetInParentheses),
	)
}
