package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitIgnore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that we can't ignore the .gitignore file, then ignore/exclude other files",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("toExclude", "")
		shell.CreateFile("toIgnore", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(`?? .gitignore`).IsSelected(),
				Contains(`?? toExclude`),
				Contains(`?? toIgnore`),
			).
			Press(keys.Files.IgnoreFile).
			// ensure we can't exclude the .gitignore file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .git/info/exclude")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot exclude .gitignore")).Confirm()
			}).
			Press(keys.Files.IgnoreFile).
			// ensure we can't ignore the .gitignore file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .gitignore")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot ignore .gitignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", DoesNotContain(".gitignore"))
			}).
			SelectNextItem().
			Press(keys.Files.IgnoreFile).
			// exclude a file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .git/info/exclude")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", Contains("toExclude"))
			}).
			SelectNextItem().
			Press(keys.Files.IgnoreFile).
			// ignore a file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .gitignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals("toIgnore\n"))
				t.FileSystem().FileContent(".git/info/exclude", Contains("toExclude"))
			})
	},
})
