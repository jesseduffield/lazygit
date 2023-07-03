package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NestedFilter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter in the several nested panels and verify the filters are preserved as you escape back to the surface",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// need to create some branches, each with their own commits
		shell.NewBranch("branch-gold")
		shell.CreateFileAndAdd("apple", "apple")
		shell.CreateFileAndAdd("orange", "orange")
		shell.CreateFileAndAdd("grape", "grape")
		shell.Commit("commit-knife")

		shell.NewBranch("branch-silver")
		shell.UpdateFileAndAdd("apple", "apple-2")
		shell.UpdateFileAndAdd("orange", "orange-2")
		shell.UpdateFileAndAdd("grape", "grape-2")
		shell.Commit("commit-spoon")

		shell.NewBranch("branch-bronze")
		shell.UpdateFileAndAdd("apple", "apple-3")
		shell.UpdateFileAndAdd("orange", "orange-3")
		shell.UpdateFileAndAdd("grape", "grape-3")
		shell.Commit("commit-fork")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains(`branch-bronze`).IsSelected(),
				Contains(`branch-silver`),
				Contains(`branch-gold`),
			).
			FilterOrSearch("sil").
			Lines(
				Contains(`branch-silver`).IsSelected(),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains(`commit-spoon`).IsSelected(),
				Contains(`commit-knife`),
			).
			FilterOrSearch("knife").
			Lines(
				// sub-commits view searches, it doesn't filter, so we haven't filtered down the list
				Contains(`commit-spoon`),
				Contains(`commit-knife`).IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains(`apple`).IsSelected(),
				Contains(`grape`),
				Contains(`orange`),
			).
			FilterOrSearch("grape").
			Lines(
				Contains(`apple`),
				Contains(`grape`).IsSelected(),
				Contains(`orange`),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			FilterOrSearch("newline").
			SelectedLine(Contains("No newline at end of file")).
			PressEscape(). // cancel search
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			// escape to commit-files view
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains(`apple`),
				Contains(`grape`).IsSelected(),
				Contains(`orange`),
			).
			Tap(func() {
				t.Views().Search().IsVisible().Content(Contains("matches for 'grape'"))
			}).
			// cancel search
			PressEscape().
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Lines(
				Contains(`apple`),
				Contains(`grape`).IsSelected(),
				Contains(`orange`),
			).
			// escape to sub-commits view
			PressEscape()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains(`commit-spoon`),
				Contains(`commit-knife`).IsSelected(),
			).
			Tap(func() {
				t.Views().Search().IsVisible().Content(Contains("matches for 'knife'"))
			}).
			// cancel search
			PressEscape().
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Lines(
				Contains(`commit-spoon`),
				// still selected
				Contains(`commit-knife`).IsSelected(),
			).
			// escape to branches view
			PressEscape()

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains(`branch-silver`).IsSelected(),
			).
			Tap(func() {
				t.Views().Search().IsVisible().Content(Contains("matches for 'sil'"))
			}).
			// cancel search
			PressEscape().
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Lines(
				Contains(`branch-bronze`),
				Contains(`branch-silver`).IsSelected(),
				Contains(`branch-gold`),
			)
	},
})
