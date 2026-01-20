package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AutoWrapMessageFootnotesAndTrailers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit, and test how auto-wrap preserves footnotes and trailers",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Use a ridiculously small width so that we don't have to use so much test data
		config.GetUserConfig().Git.Commit.AutoWrapWidth = 20
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("subject").
			SwitchToDescription().
			Type("[1]: https://github.com/jesseduffield/lazygit").
			AddNewline().
			Type("Co-authored-by: John Smith <jsmith@gmail.com>").
			Content(Equals("[1]: https://github.com/jesseduffield/lazygit\nCo-authored-by: John Smith <jsmith@gmail.com>")).
			SwitchToSummary().
			Confirm()
	},
})
