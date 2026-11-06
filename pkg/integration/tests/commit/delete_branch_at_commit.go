package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteBranchAtCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a branch pointing at the selected commit from the Commits panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.NewBranch("branch1")
		shell.NewBranch("branch2")
		shell.EmptyCommit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.Remove)

		// only the checked-out branch points here, and that one can't be deleted
		t.ExpectPopup().Menu().
			Title(Equals("Drop commit or delete branch")).
			Lines(
				Contains("Drop commit").IsSelected(),
				Contains("Delete branch"),
				Contains("Cancel"),
			).
			Select(Contains("Delete branch")).
			Tooltip(Contains("Disabled: You cannot delete the checked out branch!")).
			Cancel()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("one")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Drop commit or delete branch")).
			Lines(
				Contains("Drop commit").IsSelected(),
				Contains("Delete branch"),
				Contains("Cancel"),
			).
			Select(Contains("Delete branch")).
			Tooltip(Contains("Disabled: No branches found at selected commit.")).
			Cancel()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("two")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Drop commit or delete branch")).
			Lines(
				Contains("Drop commit").IsSelected(),
				Contains("Delete branch 'branch1'"),
				Contains("Delete branch 'master'"),
				Contains("Delete branches"),
				Contains("Cancel"),
			).
			Select(Contains("Delete branches")).
			Tooltip(Equals("branch1, master")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"1"})
			})

		t.ExpectPopup().Menu().
			Title(Equals("Delete branch 'branch1'?")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"c"})
			})

		t.Views().Branches().
			Lines(
				Contains("branch2").IsSelected(),
				Contains("master"),
			)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("three"),
				Contains("two").DoesNotContain("branch1").IsSelected(),
				Contains("one"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Drop commit or delete branch")).
			Lines(
				Contains("Drop commit").IsSelected(),
				Contains("Delete branch 'master'"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
