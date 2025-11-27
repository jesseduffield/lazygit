package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteRemoteTagWhenBranchWithSameNameExists = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a remote tag when a remote branch with the same name exists",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CloneIntoRemote("origin")
		shell.CreateLightweightTag("xyz", "HEAD")
		shell.PushBranch("origin", "HEAD:refs/tags/xyz") // abusing PushBranch to push a tag
		shell.PushBranch("origin", "HEAD:refs/heads/xyz")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("xyz").IsSelected(),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().
			Menu().
			Title(Equals("Delete tag 'xyz'?")).
			Select(Contains("Delete remote tag")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Remote from which to remove tag 'xyz':")).
			InitialText(Equals("origin")).
			SuggestionLines(
				Contains("origin"),
			).
			Confirm()

		t.ExpectPopup().
			Confirmation().
			Title(Equals("Delete tag 'xyz'?")).
			Content(Equals("Are you sure you want to delete the remote tag 'xyz' from 'origin'?")).
			Confirm()

		t.ExpectToast(Equals("Remote tag deleted"))

		t.Shell().AssertRemoteTagNotFound("origin", "xyz")
	},
})
