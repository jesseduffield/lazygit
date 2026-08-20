package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoSelectionOverAConflictHint = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The hint for a conflict that has to be resolved by picking a side shows no selection, diff and all",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.RunShellCommand(`echo 1 > foo && echo 1 > bar`)
		shell.RunShellCommand(`git checkout -b base && git add . && git commit -m base`)

		// theirs: delete foo, modify bar
		shell.RunShellCommand(`git checkout -b theirs`)
		shell.RunShellCommand(`git rm foo && echo 2 > bar && git add bar && git commit -m theirs`)

		// ours: modify foo, delete bar
		shell.RunShellCommand(`git checkout base && git checkout -b ours`)
		shell.RunShellCommand(`echo 2 > foo && git add foo && git rm bar && git commit -m ours`)

		shell.RunCommandExpectError([]string{"git", "merge", "theirs"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The hint for a file deleted on one side and modified on the other explains
		// itself with a diff of what the other side did. That diff is no more ours to
		// stage than the words above it, so the pane holds nothing to point at.
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("DU bar")).
			Tap(func() {
				t.Views().Main().
					Content(Contains("Conflict: this file was deleted in the current changes")).
					Content(Contains("Incoming changes:")).
					Content(Contains("+2"))
			}).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			PressPrimaryAction().
			Tap(func() {
				t.ExpectToast(Contains("There is nothing to select here"))
			}).
			// A pane with nothing to select still has lines to point at. A modified
			// click names its own line rather than acting on the selection, so it
			// opens the file there. The file is in the working tree for a conflict
			// like this one, holding the modified side; you may want to copy a piece
			// of it elsewhere before resolving the conflict by deleting the file.
			//
			// The cursor moves where a plain click points even here, so the assertion
			// below says which row the modified click then lands on.
			Click(0, 17).
			SelectedLine(Contains("+2")).
			AltClick(0, 17).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/bar:1\n"))
			})
	},
})
