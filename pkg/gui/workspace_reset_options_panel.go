package gui

import (
	"fmt"

	"github.com/fatih/color"
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.NukeWorkingTree).ResetAndClean(); err != nil {
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.DiscardUnstagedFileChanges).DiscardAnyUnstagedFileChanges(); err != nil {
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.RemoveUntrackedFiles).RemoveUntrackedFiles(); err != nil {
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.SoftReset).ResetSoft("HEAD"); err != nil {
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.MixedReset).ResetMixed("HEAD"); err != nil {
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
				if err := gui.Git.WithSpan(gui.Tr.Spans.HardReset).ResetHard("HEAD"); err != nil {
					return gui.SurfaceError(err)
				}

				return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{FILES}})
			},
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
