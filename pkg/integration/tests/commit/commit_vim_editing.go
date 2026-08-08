package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitVimEditing = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit using vim-style modal editing in the commit message panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.VimStyleEditing = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("fix the the bug").
			Content(Equals("fix the the bug"))

		t.Views().CommitMessage().PressEscape()

		t.ExpectPopup().CommitMessagePanel().
			Type("0wdaw").
			Content(Equals("fix the bug")).
			Type("u").
			Content(Equals("fix the the bug"))

		t.Views().CommitMessage().Press(config.Keybinding{"<c-r>"})

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("fix the bug")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("fix the bug"),
			)
	},
})
