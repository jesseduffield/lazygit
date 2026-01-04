package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Gitignore = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that we can't ignore the .gitignore file, then ignore/exclude other files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("toExclude", "")
		shell.CreateFile("toIgnore", "")
		shell.CreateFile("toIgnoreGlobally", "")

		// using a relative path here in order to not pollute the home dir of the user running the tests
		shell.SetConfig("core.excludesFile", "../my-global-git-ignore")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ?? .gitignore"),
				Equals("  ?? toExclude"),
				Equals("  ?? toIgnore"),
				Equals("  ?? toIgnoreGlobally"),
			).
			SelectNextItem().
			Press(keys.Files.IgnoreFile).
			// ensure we can't exclude the .gitignore file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .git/info/exclude")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot exclude .gitignore")).Confirm()
			}).
			Press(keys.Files.IgnoreFile).
			// ensure we can't ignore the .gitignore file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .gitignore")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot ignore .gitignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", DoesNotContain(".gitignore"))
			}).
			Press(keys.Files.IgnoreFile).
			// ensure we can't globally ignore the .gitignore file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to ../my-global-git-ignore")).Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Equals("Cannot ignore .gitignore")).Confirm()

				t.FileSystem().PathNotPresent("../my-global-git-ignore")
			}).
			SelectNextItem().
			Press(keys.Files.IgnoreFile).
			// exclude a file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .git/info/exclude")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals(""))
				t.FileSystem().FileContent(".git/info/exclude", Contains("toExclude"))
			}).
			Lines(
				Equals("▼ /"),
				Equals("  ?? .gitignore"),
				Equals("  ?? toIgnore").IsSelected(),
				Equals("  ?? toIgnoreGlobally"),
			).
			Press(keys.Files.IgnoreFile).
			// ignore a file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .gitignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals("toIgnore\n"))
				t.FileSystem().FileContent(".git/info/exclude", Contains("toExclude"))
			}).
			Lines(
				Equals("▼ /"),
				Equals("  ?? .gitignore"),
				Equals("  ?? toIgnoreGlobally").IsSelected(),
			).
			Press(keys.Files.IgnoreFile).
			// ignore a file
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to ../my-global-git-ignore")).Confirm()

				t.FileSystem().FileContent(".gitignore", Equals("toIgnore\n"))
				t.FileSystem().FileContent(".git/info/exclude", Contains("toExclude"))
				t.FileSystem().FileContent("../my-global-git-ignore", Contains("toIgnoreGlobally"))
			}).
			Lines(
				Equals("?? .gitignore"),
			)
	},
})
