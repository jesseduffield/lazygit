package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var Filter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensures that when there are merge conflicts, the files panel only shows conflicted files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFiles(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  UU file1"),
				Equals("  UU file2"),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("No filter")).
					Confirm()
			}).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  UU file1"),
				Equals("  UU file2"),
				// now we see the non-merge conflict file
				Equals("  A  file3"),
			)
	},
})
