package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExcludeWithoutInfoDir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Exclude a file when .git/info directory does not exist",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("toExclude", "")
		shell.DeleteFile(".git/info")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("toExclude").IsSelected(),
			).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .git/info/exclude")).Confirm()
			}).
			IsEmpty()

		t.FileSystem().FileContent(".git/info/exclude", Contains("/toExclude"))
	},
})
