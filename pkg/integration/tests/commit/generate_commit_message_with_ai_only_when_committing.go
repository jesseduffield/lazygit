package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GenerateCommitMessageWithAiOnlyWhenCommitting = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The commit menu offers AI generation only when committing staged changes, not when rewording or tagging",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Commit.AI.APIKeyEnv = ""
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file", "staged content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			ContainsLines(Contains("Generate commit message with AI")).
			Cancel()

		t.ExpectPopup().CommitMessagePanel().
			Cancel()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("initial commit").IsSelected(),
			).
			Press(keys.Commits.RenameCommit)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Reword commit")).
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Lines(
				Contains("Open in editor"),
				Contains("Add co-author"),
				Contains("Copy commit message to clipboard"),
				Contains("Paste commit message from clipboard"),
				Contains("Clear commit message"),
				Contains("Cancel"),
			).
			Cancel()

		t.ExpectPopup().CommitMessagePanel().
			Cancel()

		t.Views().Tags().
			Focus().
			Press(keys.Universal.New)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Tag name")).
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Lines(
				Contains("Open in editor"),
				Contains("Add co-author"),
				Contains("Copy commit message to clipboard"),
				Contains("Paste commit message from clipboard"),
				Contains("Clear commit message"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
