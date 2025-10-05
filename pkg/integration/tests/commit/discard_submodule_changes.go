package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardSubmoduleChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding changes to a submodule from an old commit.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("Initial commit")
		shell.CloneIntoSubmodule("submodule", "submodule")
		shell.Commit("Add submodule")

		shell.AddFileInWorktreeOrSubmodule("submodule", "file", "content")
		shell.CommitInWorktreeOrSubmodule("submodule", "add file in submodule")
		shell.GitAdd("submodule")
		shell.Commit("Update submodule")

		shell.UpdateFileInWorktreeOrSubmodule("submodule", "file", "changed content")
		shell.CommitInWorktreeOrSubmodule("submodule", "change file in submodule")
		shell.GitAdd("submodule")
		shell.Commit("Update submodule again")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Update submodule again").IsSelected(),
				Contains("Update submodule"),
				Contains("Add submodule"),
				Contains("Initial commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M submodule").IsSelected(),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Contains("Are you sure you want to remove changes to the selected file(s) from this commit?")).
			Confirm()

		t.Shell().RunCommand([]string{"git", "submodule", "update"})
		t.FileSystem().FileContent("submodule/file", Equals("content"))
	},
})
