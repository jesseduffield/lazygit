package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectionCommandTooltipsFollowTheDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The tooltips of the selection commands describe what they do over the diff they are offered on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nTWO\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Over the working tree's diff the two keys act on the index.
		t.Views().Files().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Select(Contains("Stage")).
			Tooltip(Equals("Toggle selection staged / unstaged.")).
			Select(Contains("Discard")).
			Tooltip(Contains("discard the change using `git reset`")).
			Cancel()

		t.Views().Main().PressEscape()

		// Over a commit's diff they build a custom patch and rewrite the commit, so
		// the index wording would be wrong; taking lines out of a commit is worth a
		// warning of its own.
		t.Views().Commits().
			Focus().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Select(Contains("Toggle lines in patch")).
			Tooltip(Equals("")).
			Select(Contains("Remove lines from commit")).
			Tooltip(Contains("runs an interactive rebase in the background")).
			Cancel()
	},
})
