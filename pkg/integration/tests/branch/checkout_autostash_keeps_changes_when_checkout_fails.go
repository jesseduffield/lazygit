package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutAutostashKeepsChangesWhenCheckoutFails = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check out a branch with autostash when an untracked file is in the way, keeping the local changes in the working tree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.55.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "a\n\nb")
		shell.Commit("add file")
		shell.UpdateFileAndAdd("file", "a\n\nc")
		shell.CreateFileAndAdd("untracked", "theirs")
		shell.Commit("edit last line and add untracked")

		shell.Checkout("HEAD^")
		shell.UpdateFile("file", "b\n\nb")
		shell.CreateFile("untracked", "mine")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Equals("   M file"),
				Equals("  ?? untracked"),
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

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("untracked working tree files would be overwritten by checkout")).
			Confirm()

		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Equals("   M file"),
				Equals("  ?? untracked"),
			)

		t.FileSystem().FileContent("file", Equals("b\n\nb"))

		t.Views().Stash().IsEmpty()
	},
})
