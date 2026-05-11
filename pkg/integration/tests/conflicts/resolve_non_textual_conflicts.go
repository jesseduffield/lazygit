package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResolveNonTextualConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resolve non-textual merge conflicts (e.g. one side modified, the other side deleted)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.RunShellCommand(`echo test1 > both-deleted1.txt`)
		shell.RunShellCommand(`echo test2 > both-deleted2.txt`)
		shell.RunShellCommand(`git checkout -b conflict && git add both-deleted1.txt both-deleted2.txt`)
		shell.RunShellCommand(`echo haha1 > deleted-them1.txt && git add deleted-them1.txt`)
		shell.RunShellCommand(`echo haha2 > deleted-them2.txt && git add deleted-them2.txt`)
		shell.RunShellCommand(`echo haha1 > deleted-us1.txt && git add deleted-us1.txt`)
		shell.RunShellCommand(`echo haha2 > deleted-us2.txt && git add deleted-us2.txt`)
		shell.RunShellCommand(`git commit -m one`)

		// stuff on other branch
		shell.RunShellCommand(`git branch conflict_second`)
		shell.RunShellCommand(`git mv both-deleted1.txt added-them-changed-us1.txt`)
		shell.RunShellCommand(`git mv both-deleted2.txt added-them-changed-us2.txt`)
		shell.RunShellCommand(`git rm deleted-them1.txt deleted-them2.txt`)
		shell.RunShellCommand(`echo modded1 > deleted-us1.txt && git add deleted-us1.txt`)
		shell.RunShellCommand(`echo modded2 > deleted-us2.txt && git add deleted-us2.txt`)
		shell.RunShellCommand(`git commit -m "two"`)

		// stuff on our branch
		shell.RunShellCommand(`git checkout conflict_second`)
		shell.RunShellCommand(`git mv both-deleted1.txt changed-them-added-us1.txt`)
		shell.RunShellCommand(`git mv both-deleted2.txt changed-them-added-us2.txt`)
		shell.RunShellCommand(`echo modded1 > deleted-them1.txt && git add deleted-them1.txt`)
		shell.RunShellCommand(`echo modded2 > deleted-them2.txt && git add deleted-them2.txt`)
		shell.RunShellCommand(`git rm deleted-us1.txt deleted-us2.txt`)
		shell.RunShellCommand(`git commit -m "three"`)
		shell.RunShellCommand(`git reset --hard conflict_second`)
		shell.RunCommandExpectError([]string{"git", "merge", "conflict"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		resolve := func(filename string, menuChoice string) {
			t.Views().Files().
				NavigateToLine(Contains(filename)).
				Tap(func() {
					t.Views().Main().Content(Contains("Conflict:"))
				}).
				Press(keys.Universal.GoInto).
				Tap(func() {
					t.ExpectPopup().Menu().Title(Equals("Merge conflicts")).
						Select(Contains(menuChoice)).
						Confirm()
				})
		}

		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  UA added-them-changed-us1.txt"),
				Equals("  UA added-them-changed-us2.txt"),
				Equals("  DD both-deleted1.txt"),
				Equals("  DD both-deleted2.txt"),
				Equals("  AU changed-them-added-us1.txt"),
				Equals("  AU changed-them-added-us2.txt"),
				Equals("  UD deleted-them1.txt"),
				Equals("  UD deleted-them2.txt"),
				Equals("  DU deleted-us1.txt"),
				Equals("  DU deleted-us2.txt"),
			).
			Tap(func() {
				resolve("added-them-changed-us1.txt", "Delete file")
				resolve("added-them-changed-us2.txt", "Keep file")
				resolve("both-deleted1.txt", "Delete file")
				resolve("both-deleted2.txt", "Delete file")
				resolve("changed-them-added-us1.txt", "Delete file")
				resolve("changed-them-added-us2.txt", "Keep file")
				resolve("deleted-them1.txt", "Delete file")
				resolve("deleted-them2.txt", "Keep file")
				resolve("deleted-us1.txt", "Delete file")
				resolve("deleted-us2.txt", "Keep file")
			}).
			Lines(
				Equals("▼ /"),
				Equals("  A  added-them-changed-us2.txt"),
				Equals("  D  changed-them-added-us1.txt"),
				Equals("  D  deleted-them1.txt"),
				Equals("  A  deleted-us2.txt"),
			)

		t.FileSystem().
			PathNotPresent("added-them-changed-us1.txt").
			FileContent("added-them-changed-us2.txt", Equals("test2\n")).
			PathNotPresent("both-deleted1.txt").
			PathNotPresent("both-deleted2.txt").
			PathNotPresent("changed-them-added-us1.txt").
			FileContent("changed-them-added-us2.txt", Equals("test2\n")).
			PathNotPresent("deleted-them1.txt").
			FileContent("deleted-them2.txt", Equals("modded2\n")).
			PathNotPresent("deleted-us1.txt").
			FileContent("deleted-us2.txt", Equals("modded2\n"))
	},
})
