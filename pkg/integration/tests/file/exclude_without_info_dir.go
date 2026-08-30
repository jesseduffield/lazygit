package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExcludeWithoutInfoDir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Exclude a file when .git/info directory does not exist",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// Remove .git/info directory to reproduce #5302
		shell.RunCommand([]string{"rm", "-rf", ".git/info"})
		shell.CreateFile("toExclude", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Focus().
			NavigateToLine(Contains("toExclude")).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .git/info/exclude")).Confirm()

				// Should succeed without error, creating .git/info/ directory automatically
				t.FileSystem().FileContent(".git/info/exclude", Contains("/toExclude"))
			})
	},
})
