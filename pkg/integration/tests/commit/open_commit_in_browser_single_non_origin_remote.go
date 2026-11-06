package commit

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenCommitInBrowserSingleNonOriginRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a commit URL when the only remote is not named origin",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "echo {{link}} > /tmp/openlink"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("goran")
		shell.RunCommand([]string{"git", "remote", "set-url", "goran", "https://github.com/my-org/lazygit"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Commits.OpenInBrowser)

		// With a single remote there should be no remote-selection prompt.
		hash := t.Git().GetCommitHash("HEAD")
		t.FileSystem().FileContent(
			"/tmp/openlink",
			Equals(fmt.Sprintf("https://github.com/my-org/lazygit/commit/%s\n", hash)))
	},
})
