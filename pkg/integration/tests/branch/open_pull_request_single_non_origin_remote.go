package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestSingleNonOriginRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a pull request when the only remote is not named origin",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "echo {{link}} > /tmp/openlink"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")

		shell.NewBranch("feature")
		shell.CloneIntoRemote("goran")
		shell.SetBranchUpstream("feature", "goran/feature")
		shell.RunCommand([]string{"git", "remote", "set-url", "goran", "https://github.com/my-org/lazygit"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().
			Branches().
			Focus().
			Press(keys.Branches.CreatePullRequest)

		// With a single remote there should be no remote-selection prompt.
		t.FileSystem().FileContent(
			"/tmp/openlink",
			Equals("https://github.com/my-org/lazygit/compare/feature?expand=1\n"))
	},
})
