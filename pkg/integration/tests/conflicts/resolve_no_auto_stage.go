package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ResolveNoAutoStage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resolving conflicts without auto-staging",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoStageResolvedConflicts = false
	},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFiles(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1").IsSelected(),
				Contains("UU").Contains("file2"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			PressPrimaryAction()

		t.Views().Files().
			IsFocused().
			// Resolving the conflict didn't auto-stage it
			Lines(
				Contains("UU").Contains("file1").IsSelected(),
				Contains("UU").Contains("file2"),
			).
			// So do that manually
			PressPrimaryAction().
			Lines(
				Contains("UU").Contains("file2").IsSelected(),
			).
			// Trying to stage a file that still has conflicts is not allowed:
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("Cannot stage/unstage directory containing files with inline merge conflicts.")).
					Confirm()
			}).
			PressEnter()

		// coincidentally these files have the same conflict
		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			PressPrimaryAction()

		t.Views().Files().
			IsFocused().
			// Again, resolving the conflict didn't auto-stage it
			Lines(
				Contains("UU").Contains("file2").IsSelected(),
			).
			// Doing that manually now works:
			PressPrimaryAction().
			Lines(
				Contains("A").Contains("file3").IsSelected(),
			)
	},
})
