package conflicts

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
	"github.com/lobes/lazytask/pkg/integration/tests/shared"
)

var ResolveExternally = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensures that when merge conflicts are resolved outside of lazygit, lazygit prompts you to continue",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved content")
			}).
			Press(keys.Universal.Refresh)

		t.Common().ContinueOnConflictsResolved()

		t.Views().Files().
			IsEmpty()
	},
})
