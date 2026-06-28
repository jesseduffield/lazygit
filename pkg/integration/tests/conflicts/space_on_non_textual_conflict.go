package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SpaceOnNonTextualConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing space on a non-textual conflict opens the resolution menu; staging is disabled for a range that includes one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
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
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("DU bar"),
				Contains("UD foo"),
			).
			// Pressing space on a single non-textual conflict opens the
			// resolution menu rather than trying to stage it.
			NavigateToLine(Contains("bar")).
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Merge conflicts")).Cancel()
			}).
			// Staging is disabled for a range selection that includes a conflict.
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("foo")).
			PressPrimaryAction().
			Tap(func() {
				t.ExpectToast(Contains("Cannot stage a selection that includes files with merge conflicts"))
			}).
			// Entering a range selection is disabled too, with the usual toast.
			Press(keys.Universal.GoInto).
			Tap(func() {
				t.ExpectToast(Contains("does not support range selection"))
			})
	},
})
