package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ResolveMultipleFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensures that a file whose conflicts have been resolved keeps being shown while other files still have conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFiles(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Common().PretendMergeOrRebaseStartedInLazygit()

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  UU file1"),
				Equals("  UU file2"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			SelectNextItem().
			PressPrimaryAction()

		// The resolved file is still shown, and stays selected so that its diff
		// can be reviewed
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  M  file1").IsSelected(),
				Equals("  UU file2"),
			).
			SelectNextItem().
			PressEnter()

		// coincidentally these files have the same conflict
		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("======="),
				Contains("Second Change"),
				Contains(">>>>>>>"),
			).
			PressPrimaryAction()

		// Now that all conflicts are resolved, the filter is turned off again
		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Equals("  M  file1"),
				Equals("  M  file2").IsSelected(),
				Equals("  A  file3"),
			)

		t.Common().ContinueOnConflictsResolved("merge")
	},
})
