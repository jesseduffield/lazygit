package conflicts

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
	"github.com/lobes/lazytask/pkg/integration/tests/shared"
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
				Contains("UU").Contains("file1").IsSelected(),
				Contains("UU").Contains("file2"),
			).
			Press(keys.Files.OpenStatusFilter).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Reset filter")).
					Confirm()
			}).
			Lines(
				Contains("UU").Contains("file1").IsSelected(),
				Contains("UU").Contains("file2"),
				// now we see the non-merge conflict file
				Contains("A ").Contains("file3"),
			)
	},
})
