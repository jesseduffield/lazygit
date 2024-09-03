package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutAutostash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check out a branch that requires performing autostash",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "a\n\nb")
		shell.Commit("add file")
		shell.UpdateFileAndAdd("file", "a\n\nc")
		shell.Commit("edit last line")

		shell.Checkout("HEAD^")
		shell.UpdateFile("file", "b\n\nb")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Lines(
				Contains("file"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				MatchesRegexp(`\*.*HEAD`).IsSelected(),
				Contains("master"),
			).
			NavigateToLine(Contains("master")).
			PressPrimaryAction()

		t.ExpectPopup().Confirmation().
			Title(Contains("Autostash?")).
			Content(Contains("You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master").IsSelected(),
			)

		t.Git().CurrentBranchName("master")

		t.Views().Files().
			Lines(
				Contains("file"),
			)
	},
})
