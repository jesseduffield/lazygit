package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardNestedRepositories = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a nested repository asks first if that would discard its history or changes that it has that aren't committed",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")
		shell.RunCommand([]string{"git", "rm", "--cached", "modules/sub"})
		shell.RunCommand([]string{"git", "rm", ".gitmodules"})
		shell.Commit("remove submodule")
		shell.CreateFile("modules/sub/new-file", "content")

		shell.RunCommand([]string{"git", "init", "nested"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? modules/sub (removed submodule)").IsSelected(),
				Equals("?? nested (nested repository)"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Discard changes")).
					Content(Equals("Discarding this also discards the following. Are you sure?\n\n" +
						"- 'modules/sub': uncommitted changes\n" +
						"- 'nested': the repository with all of its history")).
					Cancel()
			}).
			Lines(
				Equals("?? modules/sub (removed submodule)").IsSelected(),
				Equals("?? nested (nested repository)").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Discard changes")).
					Content(Contains("Are you sure?")).
					Confirm()
			}).
			IsEmpty()

		t.FileSystem().PathNotPresent("modules/sub")
		t.FileSystem().PathNotPresent("nested")
	},
})
