package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ContinuePromptDismissedWhenResolvedExternally = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When the prompt to continue a merge is showing and the merge is then continued outside lazygit, dismiss the prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Common().PretendMergeOrRebaseStartedInLazygit()

		// Resolve the conflict and refresh so lazygit prompts us to continue.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved content")
			}).
			Press(keys.Universal.Refresh)

		t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue the merge?"))

		// While the prompt is up, the merge is continued outside lazygit (e.g. by
		// a coding agent).
		t.Shell().ContinueMerge()

		// Simulate lazygit noticing the change (as it would on its next refresh or
		// when the window regains focus); the stale prompt is dismissed.
		t.FocusIn()

		t.Views().Files().
			IsFocused().
			IsEmpty()

		t.Views().Information().Content(DoesNotContain("Merging"))
	},
})
