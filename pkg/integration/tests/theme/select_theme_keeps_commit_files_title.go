package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectThemeKeepsCommitFilesTitle = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme while the commit files of a commit are shown, but another panel is focused, keeps the commit in the title of the commit files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "content").
			Commit("add file")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("add file").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Title(MatchesRegexp(`^Diff files \([0-9a-f]{7} add file\)$`))

		t.Views().Files().
			Focus().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")
		t.Views().CommitFiles().
			IsVisible().
			Title(MatchesRegexp(`^Diff files \([0-9a-f]{7} add file\)$`))
	},
})
