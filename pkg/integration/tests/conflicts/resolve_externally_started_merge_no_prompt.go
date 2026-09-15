package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ResolveExternallyStartedMergeNoPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When a merge started outside lazygit has its conflicts resolved, don't prompt to continue it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Start the merge by running git directly and never tell lazygit it was
		// the one to start it, so from lazygit's point of view it was started
		// externally (e.g. by a coding agent in another terminal).
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

		// No prompt to continue the merge appears; we stay in the files view
		// with the conflict resolved and the merge still in progress.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file"),
			)

		t.Views().Information().Content(Contains("Merging"))
	},
})
