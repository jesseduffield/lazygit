package filter_by_author

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MailmapSelectAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Filter commits by a mailmapped author when the selected commit's author is shown",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.SetAuthor("Legacy Name", "old@example.com")
		shell.EmptyCommit("mapped commit")
		shell.CreateFile(".mailmap", "Canonical Name <new@example.com> Legacy Name <old@example.com>\n")
		shell.RunCommand([]string{"git", "config", "--local", "log.mailmap", "false"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(Contains("mapped commit")).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			/* EXPECTED:
			Select(Contains("Filter by 'Canonical Name <new@example.com>'")).
			ACTUAL: */
			Select(Contains("Filter by 'Canonical Name <old@example.com>'")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			/* EXPECTED:
			Lines(Contains("mapped commit"))
			ACTUAL: */
			IsEmpty()
	},
})
