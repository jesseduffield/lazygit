package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// The second keypress arrives before the refresh triggered by the first one
// has rebuilt the commits model. The handler reads the selected todo from the
// model at the already-advanced selection index, so with the stale, pre-move
// model it grabs the todo the first move swapped with and moves that one back
// down — turning the two presses into a net no-op instead of moving the
// selected todo down two slots. This is what happens when holding down the
// move-down key to move a todo several slots.
//
// We continue the rebase and assert the resulting commit order rather than
// asserting the todo list, because the two presses also spawn two racing
// refreshes whose updates can land in either order, so what the todo list
// shows in the broken state is not deterministic (it can even disagree with
// the todo file). The rebase replays what's in the file.
var MoveTodoDownWithRapidKeypresses = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a todo down two slots with two keypresses in rapid succession",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			NavigateToLine(Contains("commit-01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("--- Commits ---"),
				Contains("commit-01").IsSelected(),
			).
			NavigateToLine(Contains("commit-04")).
			PressRapidly(keys.Commits.MoveDownCommit, keys.Commits.MoveDownCommit).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			/* EXPECTED:
			Lines(
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-04"),
				Contains("commit-01"),
			)
			ACTUAL: */
			Lines(
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			)
	},
})
