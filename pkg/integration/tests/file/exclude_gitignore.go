package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExcludeGitignore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Failed attempt at excluding and ignoring the .gitignore file",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(`?? .gitignore`).IsSelected(),
			).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .git/info/exclude")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot exclude .gitignore")).Confirm()
			}).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("ignore or exclude file")).Select(Contains("add to .gitignore")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot ignore .gitignore")).Confirm()
			})

		t.FileSystem().FileContent(".gitignore", Equals(""))
		t.FileSystem().FileContent(".git/info/exclude", DoesNotContain(".gitignore"))
	},
})
