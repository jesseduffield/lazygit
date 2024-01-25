package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardAllDirChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding all changes in a directory",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// typically we would use more bespoke shell methods here, but I struggled to find a way to do that,
		// and this is copied over from a legacy integration test which did everything in a big shell script
		// so I'm just copying it across.

		shell.CreateDir("dir")

		// common stuff
		shell.RunShellCommand(`echo test > dir/both-deleted.txt`)
		shell.RunShellCommand(`git checkout -b conflict && git add dir/both-deleted.txt`)
		shell.RunShellCommand(`echo bothmodded > dir/both-modded.txt && git add dir/both-modded.txt`)
		shell.RunShellCommand(`echo haha > dir/deleted-them.txt && git add dir/deleted-them.txt`)
		shell.RunShellCommand(`echo haha2 > dir/deleted-us.txt && git add dir/deleted-us.txt`)
		shell.RunShellCommand(`echo mod > dir/modded.txt && git add dir/modded.txt`)
		shell.RunShellCommand(`echo mod > dir/modded-staged.txt && git add dir/modded-staged.txt`)
		shell.RunShellCommand(`echo del > dir/deleted.txt && git add dir/deleted.txt`)
		shell.RunShellCommand(`echo del > dir/deleted-staged.txt && git add dir/deleted-staged.txt`)
		shell.RunShellCommand(`echo change-delete > dir/change-delete.txt && git add dir/change-delete.txt`)
		shell.RunShellCommand(`echo delete-change > dir/delete-change.txt && git add dir/delete-change.txt`)
		shell.RunShellCommand(`echo double-modded > dir/double-modded.txt && git add dir/double-modded.txt`)
		shell.RunShellCommand(`echo "renamed\nhaha" > dir/renamed.txt && git add dir/renamed.txt`)
		shell.RunShellCommand(`git commit -m one`)

		// stuff on other branch
		shell.RunShellCommand(`git branch conflict_second && git mv dir/both-deleted.txt dir/added-them-changed-us.txt`)
		shell.RunShellCommand(`git commit -m "dir/both-deleted.txt renamed in dir/added-them-changed-us.txt"`)
		shell.RunShellCommand(`echo blah > dir/both-added.txt && git add dir/both-added.txt`)
		shell.RunShellCommand(`echo mod1 > dir/both-modded.txt && git add dir/both-modded.txt`)
		shell.RunShellCommand(`rm dir/deleted-them.txt && git add dir/deleted-them.txt`)
		shell.RunShellCommand(`echo modded > dir/deleted-us.txt && git add dir/deleted-us.txt`)
		shell.RunShellCommand(`git commit -m "two"`)

		// stuff on our branch
		shell.RunShellCommand(`git checkout conflict_second`)
		shell.RunShellCommand(`git mv dir/both-deleted.txt dir/changed-them-added-us.txt`)
		shell.RunShellCommand(`git commit -m "both-deleted.txt renamed in dir/changed-them-added-us.txt"`)
		shell.RunShellCommand(`echo mod2 > dir/both-modded.txt && git add dir/both-modded.txt`)
		shell.RunShellCommand(`echo blah2 > dir/both-added.txt && git add dir/both-added.txt`)
		shell.RunShellCommand(`echo modded > dir/deleted-them.txt && git add dir/deleted-them.txt`)
		shell.RunShellCommand(`rm dir/deleted-us.txt && git add dir/deleted-us.txt`)
		shell.RunShellCommand(`git commit -m "three"`)
		shell.RunShellCommand(`git reset --hard conflict_second`)
		shell.RunCommandExpectError([]string{"git", "merge", "conflict"})

		shell.RunShellCommand(`echo "new" > dir/new.txt`)
		shell.RunShellCommand(`echo "new staged" > dir/new-staged.txt && git add dir/new-staged.txt`)
		shell.RunShellCommand(`echo mod2 > dir/modded.txt`)
		shell.RunShellCommand(`echo mod2 > dir/modded-staged.txt && git add dir/modded-staged.txt`)
		shell.RunShellCommand(`rm dir/deleted.txt`)
		shell.RunShellCommand(`rm dir/deleted-staged.txt && git add dir/deleted-staged.txt`)
		shell.RunShellCommand(`echo change-delete2 > dir/change-delete.txt && git add dir/change-delete.txt`)
		shell.RunShellCommand(`rm dir/change-delete.txt`)
		shell.RunShellCommand(`rm dir/delete-change.txt && git add dir/delete-change.txt`)
		shell.RunShellCommand(`echo "changed" > dir/delete-change.txt`)
		shell.RunShellCommand(`echo "change1" > dir/double-modded.txt && git add dir/double-modded.txt`)
		shell.RunShellCommand(`echo "change2" > dir/double-modded.txt`)
		shell.RunShellCommand(`echo before > dir/added-changed.txt && git add dir/added-changed.txt`)
		shell.RunShellCommand(`echo after > dir/added-changed.txt`)
		shell.RunShellCommand(`rm dir/renamed.txt && git add dir/renamed.txt`)
		shell.RunShellCommand(`echo "renamed\nhaha" > dir/renamed2.txt && git add dir/renamed2.txt`)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("dir").IsSelected(),
				Contains("UA").Contains("added-them-changed-us.txt"),
				Contains("AA").Contains("both-added.txt"),
				Contains("DD").Contains("both-deleted.txt"),
				Contains("UU").Contains("both-modded.txt"),
				Contains("AU").Contains("changed-them-added-us.txt"),
				Contains("UD").Contains("deleted-them.txt"),
				Contains("DU").Contains("deleted-us.txt"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			Tap(func() {
				t.Common().ContinueOnConflictsResolved()
			}).
			Lines(
				Contains("dir").IsSelected(),
				Contains(" M").Contains("added-changed.txt"),
				Contains(" D").Contains("change-delete.txt"),
				Contains("??").Contains("delete-change.txt"),
				Contains(" D").Contains("deleted.txt"),
				Contains(" M").Contains("double-modded.txt"),
				Contains(" M").Contains("modded.txt"),
				Contains("??").Contains("new.txt"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}).
			IsEmpty()
	},
})
