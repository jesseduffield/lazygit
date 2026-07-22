package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a custom patch by toggling a hunk into it from the focused main view of a commit's files, then apply it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("first commit")

		// Two separate change blocks, far enough apart to stay distinct hunks.
		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.NextItem).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			// `a` extends the selection to the whole change block, then space toggles
			// just that block into the custom patch — the other block stays out.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction().
			// After toggling the block in, the selection advances to the next stageable
			// hunk (the NINE block) — like staging advances to the next hunk — across the
			// re-render (which, when the secondary view first appears, also re-wraps the
			// narrower diff).
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			)

		t.Views().Information().Content(Contains("Building patch"))

		// The secondary view updates live with the cumulative patch — just the toggled
		// block, not the other.
		t.Views().Secondary().
			ContainsLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			Content(DoesNotContain("NINE"))

		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		// Only the toggled block reached the working tree: THREE is applied, NINE isn't.
		t.Views().Files().
			Focus().
			Lines(
				Contains("file1").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("THREE")).
			Content(DoesNotContain("NINE"))
	},
})
