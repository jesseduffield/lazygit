package tag

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var (
	tagName       = "v1.2.3"
	commitMessage = "Tagged commit"
)

var JumpToCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a branch and jump to the commit associated with it.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("blah").
			EmptyCommit("blah").
			EmptyCommit(commitMessage).
			Tag(tagName).
			EmptyCommit("blah").
			EmptyCommit("blah")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.
			Views().
			Branches().
			Focus().
			Press(keys.Universal.NextTab).
			Press(keys.Universal.NextTab).
			Press(keys.Tags.JumpToCommit)

		t.
			Views().
			Commits().
			SelectedLine(MatchesRegexp(fmt.Sprintf(`^.* CI %s %s\s*`, tagName, commitMessage)))
	},
})
