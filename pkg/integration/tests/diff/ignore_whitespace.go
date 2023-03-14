package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View diff with and without ignoring whitespace",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("first commit")
		// First line has a real change, second line changes only indentation:
		shell.UpdateFileAndAdd("file1", "first line changed\n  second line\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Tap(func() {
				// By default, both changes are shown in the diff:
				t.Views().Main().Content(Contains("-first line\n-second line\n+first line changed\n+  second line\n"))
			}).
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			Tap(func() {
				// After enabling ignore whitespace, only the real change remains:
				t.Views().Main().Content(Contains("-first line\n+first line changed\n"))
			})
	},
})
