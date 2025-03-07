package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var postCheckoutHook = `#!/bin/bash

echo "post-checkout hook called" > hook-result
`

var CheckoutFileFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout an individual file from a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/post-checkout", postCheckoutHook)
		shell.MakeExecutable(".git/hooks/post-checkout")

		shell.CreateFileAndAdd("file", "one\n")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file", "one\ntwo\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file").IsSelected(),
			).
			Press(keys.CommitFiles.CheckoutCommitFile)

		t.Views().Files().
			Focus().
			Lines(
				Contains("M  file"),
				/* EXPECTED:
				Contains("?? hook-result"),
				*/
			)

		t.FileSystem().FileContent("file", Equals("one\n"))
		/* EXPECTED:
		t.FileSystem().FileContent("hook-result", Equals("post-checkout hook called\n"))
		ACTUAL: */
		t.FileSystem().PathNotPresent("hook-result")
	},
},
)
