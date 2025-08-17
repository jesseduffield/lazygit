package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var blockingHook = `#!/bin/bash

# For this test all we need is a hook that always fails
exit 1
`

var CommitSkipHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with skip hook using CommitChangesWithoutHook",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", blockingHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.CreateFile("file.txt", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		checkBlockingHook(t, keys)

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Lines(
				Equals("A  file.txt"),
			).
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("foo bar").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("foo bar"))
	},
})
