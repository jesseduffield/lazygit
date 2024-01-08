package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var preCommitHook = `#!/bin/bash

if [[ -f bad ]]; then
  exit 1
fi
`

var RememberCommitMessageAfterFail = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the commit message is remembered after a failed attempt at committing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", preCommitHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.CreateFileAndAdd("one", "one")

		// the presence of this file will cause the pre-commit hook to fail
		shell.CreateFile("bad", "bad")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("bad"),
				Contains("one"),
			).
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().Type("my message").Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("Git command failed")).Confirm()
			}).
			Press(keys.Universal.Remove). // remove file that triggers pre-commit hook to fail
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Lines(
				Contains("one"),
			).
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					InitialText(Equals("my message")). // it remembered the commit message
					Confirm()

				t.Views().Commits().
					Lines(
						Contains("my message"),
					)
			})
	},
})
