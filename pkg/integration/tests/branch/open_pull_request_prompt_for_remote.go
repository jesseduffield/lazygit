package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestPromptForRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a pull request and choose among multiple remotes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "echo {{link}} > /tmp/openlink"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")

		shell.NewBranch("feature")
		shell.CloneIntoRemote("upstream")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("feature", "origin/feature")

		shell.RunCommand([]string{"git", "remote", "set-url", "origin", "https://github.com/my-personal-fork/lazygit"})
		shell.RunCommand([]string{"git", "remote", "set-url", "upstream", "https://github.com/jesseduffield/lazygit"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().
			Branches().
			Focus().
			Press(keys.Branches.CreatePullRequest)

		t.ExpectPopup().
			Menu().
			Title(Equals("Select remote")).
			TopLines(
				Contains("origin"),
				Contains("upstream")).
			Select(Contains("upstream")).
			Confirm()

		t.FileSystem().FileContent(
			"/tmp/openlink",
			Equals("https://github.com/jesseduffield/lazygit/compare/feature?expand=1\n"))
	},
})
