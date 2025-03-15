package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var fileModHook = `#!/bin/bash

# For this test all we need is a hook that always fails
exit 1
`

var CommitSkipHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with skip hook using CommitChangesWithoutHook",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", fileModHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.CreateFileAndAdd("file.txt", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("A  file.txt"),
			).
			// Sanity check to make sure the hook is working:
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("foo bar").
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Git command failed.")).
			Confirm()

		// Committing without hook works:
		t.Views().Files().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("foo bar").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("foo bar"))
	},
})
