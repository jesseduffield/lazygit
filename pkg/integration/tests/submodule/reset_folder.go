package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetFolder = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset submodule changes located in a nested folder.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CreateDir("dir")
		shell.CloneIntoSubmodule("submodule1", "dir/submodule1")
		shell.CloneIntoSubmodule("submodule2", "dir/submodule2")
		shell.GitAddAll()
		shell.Commit("add submodules")

		shell.CreateFile("dir/submodule1/file", "")
		shell.CreateFile("dir/submodule2/file", "")
		shell.CreateFile("dir/file", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  ?? file"),
				Equals("   M submodule1 (submodule)"),
				Equals("   M submodule2 (submodule)"),
			).
			// Verify we cannot reset the entire folder (has nested file and submodule changes).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Multiselection not supported for submodules"))
			}).
			// Verify we cannot reset submodule + file or submodule + submodule via range select.
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Lines(
				Equals("▼ dir"),
				Equals("  ?? file").IsSelected(),
				Equals("   M submodule1 (submodule)").IsSelected(),
				Equals("   M submodule2 (submodule)"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Multiselection not supported for submodules"))
			}).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Lines(
				Equals("▼ dir"),
				Equals("  ?? file"),
				Equals("   M submodule1 (submodule)").IsSelected(),
				Equals("   M submodule2 (submodule)").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Multiselection not supported for submodules"))
			}).
			// Reset the file change.
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			NavigateToLine(Contains("▼ dir")).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("   M submodule1 (submodule)"),
				Equals("   M submodule2 (submodule)"),
			).
			// Verify we still cannot reset the entire folder (has two submodule changes).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Multiselection not supported for submodules"))
			}).
			// Reset one of the submodule changes.
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("dir/submodule1")).
					Select(Contains("Stash uncommitted submodule changes and update")).
					Confirm()
			}).
			NavigateToLine(Contains("▼ dir")).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("   M submodule2 (submodule)"),
			).
			// Now we can reset the folder (equivalent to resetting just the nested submodule change).
			// Range selecting both the folder and submodule change is allowed.
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("   M submodule2 (submodule)").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("dir/submodule2")).
					Select(Contains("Stash uncommitted submodule changes and update")).
					Cancel()
			}).
			// Or just selecting the folder itself.
			NavigateToLine(Contains("▼ dir")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("dir/submodule2")).
					Select(Contains("Stash uncommitted submodule changes and update")).
					Confirm()
			}).
			IsEmpty()
	},
})
