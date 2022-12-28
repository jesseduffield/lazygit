package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding all possible permutations of changed files",
	ExtraCmdArgs: "",
	Skip:         true, // failing due to index.lock file being created
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// typically we would use more bespoke shell methods here, but I struggled to find a way to do that,
		// and this is copied over from a legacy integration test which did everything in a big shell script
		// so I'm just copying it across.

		// common stuff
		shell.RunShellCommand(`echo test > both-deleted.txt`)
		shell.RunShellCommand(`git checkout -b conflict && git add both-deleted.txt`)
		shell.RunShellCommand(`echo bothmodded > both-modded.txt && git add both-modded.txt`)
		shell.RunShellCommand(`echo haha > deleted-them.txt && git add deleted-them.txt`)
		shell.RunShellCommand(`echo haha2 > deleted-us.txt && git add deleted-us.txt`)
		shell.RunShellCommand(`echo mod > modded.txt & git add modded.txt`)
		shell.RunShellCommand(`echo mod > modded-staged.txt & git add modded-staged.txt`)
		shell.RunShellCommand(`echo del > deleted.txt && git add deleted.txt`)
		shell.RunShellCommand(`echo del > deleted-staged.txt && git add deleted-staged.txt`)
		shell.RunShellCommand(`echo change-delete > change-delete.txt && git add change-delete.txt`)
		shell.RunShellCommand(`echo delete-change > delete-change.txt && git add delete-change.txt`)
		shell.RunShellCommand(`echo double-modded > double-modded.txt && git add double-modded.txt`)
		shell.RunShellCommand(`echo "renamed\nhaha" > renamed.txt && git add renamed.txt`)
		shell.RunShellCommand(`git commit -m one`)

		// stuff on other branch
		shell.RunShellCommand(`git branch conflict_second && git mv both-deleted.txt added-them-changed-us.txt`)
		shell.RunShellCommand(`git commit -m "both-deleted.txt renamed in added-them-changed-us.txt"`)
		shell.RunShellCommand(`echo blah > both-added.txt && git add both-added.txt`)
		shell.RunShellCommand(`echo mod1 > both-modded.txt && git add both-modded.txt`)
		shell.RunShellCommand(`rm deleted-them.txt && git add deleted-them.txt`)
		shell.RunShellCommand(`echo modded > deleted-us.txt && git add deleted-us.txt`)
		shell.RunShellCommand(`git commit -m "two"`)

		// stuff on our branch
		shell.RunShellCommand(`git checkout conflict_second`)
		shell.RunShellCommand(`git mv both-deleted.txt changed-them-added-us.txt`)
		shell.RunShellCommand(`git commit -m "both-deleted.txt renamed in changed-them-added-us.txt"`)
		shell.RunShellCommand(`echo mod2 > both-modded.txt && git add both-modded.txt`)
		shell.RunShellCommand(`echo blah2 > both-added.txt && git add both-added.txt`)
		shell.RunShellCommand(`echo modded > deleted-them.txt && git add deleted-them.txt`)
		shell.RunShellCommand(`rm deleted-us.txt && git add deleted-us.txt`)
		shell.RunShellCommand(`git commit -m "three"`)
		shell.RunShellCommand(`git reset --hard conflict_second`)
		shell.RunShellCommandExpectError(`git merge conflict`)

		shell.RunShellCommand(`echo "new" > new.txt`)
		shell.RunShellCommand(`echo "new staged" > new-staged.txt && git add new-staged.txt`)
		shell.RunShellCommand(`echo mod2 > modded.txt`)
		shell.RunShellCommand(`echo mod2 > modded-staged.txt && git add modded-staged.txt`)
		shell.RunShellCommand(`rm deleted.txt`)
		shell.RunShellCommand(`rm deleted-staged.txt && git add deleted-staged.txt`)
		shell.RunShellCommand(`echo change-delete2 > change-delete.txt && git add change-delete.txt`)
		shell.RunShellCommand(`rm change-delete.txt`)
		shell.RunShellCommand(`rm delete-change.txt && git add delete-change.txt`)
		shell.RunShellCommand(`echo "changed" > delete-change.txt`)
		shell.RunShellCommand(`echo "change1" > double-modded.txt && git add double-modded.txt`)
		shell.RunShellCommand(`echo "change2" > double-modded.txt`)
		shell.RunShellCommand(`echo before > added-changed.txt && git add added-changed.txt`)
		shell.RunShellCommand(`echo after > added-changed.txt`)
		shell.RunShellCommand(`rm renamed.txt && git add renamed.txt`)
		shell.RunShellCommand(`echo "renamed\nhaha" > renamed2.txt && git add renamed2.txt`)
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		type statusFile struct {
			status    string
			label     string
			menuTitle string
		}

		discardOneByOne := func(files []statusFile) {
			for _, file := range files {
				t.Views().Files().
					IsFocused().
					SelectedLine(Contains(file.status + " " + file.label)).
					Press(keys.Universal.Remove)

				t.ExpectPopup().Menu().Title(Equals(file.menuTitle)).Select(Contains("discard all changes")).Confirm()
			}
		}

		discardOneByOne([]statusFile{
			{status: "UA", label: "added-them-changed-us.txt", menuTitle: "added-them-changed-us.txt"},
			{status: "AA", label: "both-added.txt", menuTitle: "both-added.txt"},
			{status: "DD", label: "both-deleted.txt", menuTitle: "both-deleted.txt"},
			{status: "UU", label: "both-modded.txt", menuTitle: "both-modded.txt"},
			{status: "AU", label: "changed-them-added-us.txt", menuTitle: "changed-them-added-us.txt"},
			{status: "UD", label: "deleted-them.txt", menuTitle: "deleted-them.txt"},
			{status: "DU", label: "deleted-us.txt", menuTitle: "deleted-us.txt"},
		})

		t.ExpectPopup().Confirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Cancel()

		discardOneByOne([]statusFile{
			{status: "MD", label: "change-delete.txt", menuTitle: "change-delete.txt"},
			{status: "D ", label: "delete-change.txt", menuTitle: "delete-change.txt"},
			{status: "D ", label: "deleted-staged.txt", menuTitle: "deleted-staged.txt"},
			{status: " D", label: "deleted.txt", menuTitle: "deleted.txt"},
			{status: "MM", label: "double-modded.txt", menuTitle: "double-modded.txt"},
			{status: "M ", label: "modded-staged.txt", menuTitle: "modded-staged.txt"},
			{status: " M", label: "modded.txt", menuTitle: "modded.txt"},
			// the menu title only includes the new file
			{status: "R ", label: "renamed.txt → renamed2.txt", menuTitle: "renamed2.txt"},
			{status: "AM", label: "added-changed.txt", menuTitle: "added-changed.txt"},
			{status: "A ", label: "new-staged.txt", menuTitle: "new-staged.txt"},
			{status: "??", label: "new.txt", menuTitle: "new.txt"},
		})

		t.Views().Files().IsEmpty()
	},
})
