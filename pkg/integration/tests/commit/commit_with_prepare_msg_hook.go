package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var branchPrefixHook = `#!/bin/bash
branch=$(git symbolic-ref --short HEAD)
if [[ "$branch" =~ feature/(.*) ]]; then
    ticket="${BASH_REMATCH[1]}"
    echo "[$ticket]: $(cat "$1")" > "$1"
fi`

var CommitWithPrepareMsgHook = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with a prepare-commit-msg defined and no user commit prefixes. The hook should prefix the commit msg.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/prepare-commit-msg", branchPrefixHook)
		shell.MakeExecutable(".git/hooks/prepare-commit-msg")

		shell.NewBranch("feature/TEST-001")
		shell.CreateFile("test-prepare-commit-msg", "This is foo bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("[TEST-001]: ")).
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("[TEST-001]: "))
	},
})
