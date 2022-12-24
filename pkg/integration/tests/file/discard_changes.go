package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding all possible permutations of changed files",
	ExtraCmdArgs: "",
	Skip:         false,
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

	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(3)

		type statusFile struct {
			status string
			path   string
		}

		discardOneByOne := func(files []statusFile) {
			for _, file := range files {
				assert.MatchSelectedLine(Contains(file.status + " " + file.path))
				input.PressKeys(keys.Universal.Remove)
				assert.InMenu()
				assert.MatchCurrentViewContent(Contains("discard all changes"))
				input.Confirm()
			}
		}

		discardOneByOne([]statusFile{
			{"UA", "added-them-changed-us.txt"},
			{"AA", "both-added.txt"},
			{"DD", "both-deleted.txt"},
			{"UU", "both-modded.txt"},
			{"AU", "changed-them-added-us.txt"},
			{"UD", "deleted-them.txt"},
			{"DU", "deleted-us.txt"},
		})

		assert.InConfirm()
		assert.MatchCurrentViewTitle(Contains("continue"))
		assert.MatchCurrentViewContent(Contains("all merge conflicts resolved. Continue?"))
		input.PressKeys(keys.Universal.Return)

		discardOneByOne([]statusFile{
			{"MD", "change-delete.txt"},
			{"D ", "delete-change.txt"},
			{"D ", "deleted-staged.txt"},
			{" D", "deleted.txt"},
			{"MM", "double-modded.txt"},
			{"M ", "modded-staged.txt"},
			{" M", "modded.txt"},
			{"R ", "renamed.txt → renamed2.txt"},
			{"AM", "added-changed.txt"},
			{"A ", "new-staged.txt"},
			{"??", "new.txt"},
		})

		assert.WorkingTreeFileCount(0)
	},
})
