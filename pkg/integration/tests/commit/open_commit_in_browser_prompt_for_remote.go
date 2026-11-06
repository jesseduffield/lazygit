package commit

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenCommitInBrowserPromptForRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a commit URL and choose among multiple remotes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "echo {{link}} > /tmp/openlink"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("upstream")
		shell.CloneIntoRemote("origin")

		shell.RunCommand([]string{"git", "remote", "set-url", "origin", "https://github.com/my-personal-fork/lazygit"})
		shell.RunCommand([]string{"git", "remote", "set-url", "upstream", "https://github.com/jesseduffield/lazygit"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Commits.OpenInBrowser)

		t.ExpectPopup().
			Menu().
			Title(Equals("Select remote")).
			TopLines(
				Contains("origin"),
				Contains("upstream")).
			Select(Contains("upstream")).
			Confirm()

		hash := t.Git().GetCommitHash("HEAD")
		t.FileSystem().FileContent(
			"/tmp/openlink",
			Equals(fmt.Sprintf("https://github.com/jesseduffield/lazygit/commit/%s\n", hash)))
	},
})
