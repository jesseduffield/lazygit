package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

const preCommitHook = `#!/bin/bash

exit 1
`

var CommitSkipHook = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with pre-commit hook and skip hook config option in various situations.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(testConfig *config.AppConfig) {
		testConfig.UserConfig.Git.SkipHookPrefix = "skip! "

	},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("user.email", "Bill@example.com")
		shell.SetConfig("user.name", "Bill Smith")

		shell.CreateFileAndAdd("initial file", "initial content")
		shell.Commit("initial commit")

		shell.SetConfig("user.email", "John@example.com")
		shell.SetConfig("user.name", "John Smith")

		shell.CreateFile(".git/hooks/pre-commit", preCommitHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.CreateFileAndAdd("testfile", "I'm just testing pre-commit hooks")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		// hook should trigger when creating a regular commit
		t.ExpectPopup().CommitMessagePanel().Type("my commit message").Confirm()
		t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("Git command failed")).Confirm()

		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

			// we should be able to skip hooks when creating a regular commit
		t.ExpectPopup().CommitMessagePanel().Clear().Type("skip! my commit message").Confirm()
		t.Views().Commits().Focus().Lines(
			Contains("skip! my commit message"),
			Contains("initial commit"),
		)

		// we should be able to skip hooks when rewording a commit
		t.Views().Commits().Focus().Press(keys.Commits.RenameCommit)
		t.ExpectPopup().CommitMessagePanel().Type(" (reworded)").Confirm()

		/* EXPECTED:
						t.Views().Commits().IsFocused().
							Lines(
		                        Contains("skip! my commit message (reworded)"),
		                        Contains("initial commit"),
		                    )
				            ACTUAL:
		*/
		t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("exit status 1")).Confirm()

		// we should be able to skip hooks when changing authors
		t.Views().Commits().IsFocused().SelectedLine(Contains("CI").IsSelected())
		t.Views().Commits().Focus().Press(keys.Commits.ResetCommitAuthor)
		/* EXPECTED:
		        t.Views().Commits().IsFocused().Lines(Contains("JS").IsSelected())
				   ACTUAL:
		*/
		t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("exit status 1")).Confirm()

	},
})
