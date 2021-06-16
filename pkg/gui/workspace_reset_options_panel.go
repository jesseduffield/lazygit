package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateResetMenu() error {
	red := color.New(color.FgRed)

	nukeStr := "reset --hard HEAD && git clean -fd"
	if len(gui.State.Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, gui.Tr.LcAndResetSubmodules)
	}

	menuItems := []*menuItem{
		{
			displayStrings: []string{
				gui.Tr.LcDiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.NukeWorkingTree).Worktree().ResetAndClean(); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardUnstagedFileChanges).Worktree().DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.RemoveUntrackedFiles).Worktree().RemoveUntrackedFiles(); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.SoftReset).Branches().ResetToRef("HEAD", commands.SOFT, commands.ResetToRefOpts{}); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.MixedReset).Branches().ResetToRef("HEAD", commands.MIXED, commands.ResetToRefOpts{}); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			onPress: func() error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.HardReset).Branches().ResetToRef("HEAD", commands.HARD, commands.ResetToRefOpts{}); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
