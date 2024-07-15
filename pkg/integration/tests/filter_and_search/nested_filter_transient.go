package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This one requires some explanation: the sub-commits and diff-file contexts are
// 'transient' in that they are spawned inside a window when you need them, but
// can be relocated elsewhere if you need them somewhere else. So for example if
// I hit enter on a branch I'll see the sub-commits view, but if I then navigate
// to the reflog context and hit enter on a reflog, the sub-commits view is moved
// to the reflog window. This is because we reuse the same view (it's a limitation
// that would be nice to remove in the future).
// Nonetheless, we need to ensure that upon moving the view, the filter is cancelled.

var NestedFilterTransient = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter in a transient panel (sub-commits and diff-files) and ensure filter is cancelled when the panel is moved",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// need to create some branches, each with their own commits
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("file-one", "file-one")
		shell.CreateFileAndAdd("file-two", "file-two")
		shell.Commit("commit-one")
		shell.EmptyCommit("commit-two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains(`mybranch`).IsSelected(),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains(`commit-two`).IsSelected(),
				Contains(`commit-one`),
			).
			FilterOrSearch("one").
			Lines(
				Contains(`commit-two`),
				Contains(`commit-one`).IsSelected(),
			)

		t.Views().ReflogCommits().
			Focus().
			SelectedLine(Contains("commit: commit-two")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			// the search on the sub-commits context has been cancelled
			Lines(
				Contains(`commit-two`).IsSelected(),
				Contains(`commit-one`),
			).
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			NavigateToLine(Contains("commit-one")).
			PressEnter()

		// Now let's test the commit files context
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains(`file-one`).IsSelected(),
				Contains(`file-two`),
			).
			FilterOrSearch("two").
			Lines(
				Contains(`file-one`),
				Contains(`file-two`).IsSelected(),
			)

		t.Views().Branches().
			Focus().
			SelectedLine(Contains("mybranch")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains(`commit-two`).IsSelected(),
				Contains(`commit-one`),
			).
			NavigateToLine(Contains("commit-one")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			// the search on the commit-files context has been cancelled
			Lines(
				Contains(`file-one`).IsSelected(),
				Contains(`file-two`),
			).
			Tap(func() {
				t.Views().Search().IsInvisible()
			})
	},
})
