package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitWipWithPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with skip hook and config commitPrefix is defined. Prefix is ignored when creating WIP commits.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.CommitPrefixes = map[string][]config.CommitPrefixConfig{"repo": {{Pattern: "^\\w+\\/(\\w+-\\w+).*", Replace: "[$1]: "}}}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", blockingHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.NewBranch("feature/TEST-002")
		shell.CreateFile("test-wip-commit-prefix", "This is foo bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		checkBlockingHook(t, keys)

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("WIP")).
			Type(" foo").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("WIP foo")).
			Type(" bar").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("WIP foo bar")).
			Type(". Added something else").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("WIP foo bar. Added something else"))
	},
})
