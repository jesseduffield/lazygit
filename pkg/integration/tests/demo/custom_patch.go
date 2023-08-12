package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var usersFileContent = `package main

import "fmt"

func main() {
	// TODO: verify that this actuall works
	fmt.Println("hello world")
}
`

var CustomPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a line from an old commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(cfg *config.AppConfig) {
		setDefaultDemoConfig(cfg)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/user-authentication")
		shell.EmptyCommit("Add user authentication feature")
		shell.CreateFileAndAdd("src/users.go", "package main\n")
		shell.Commit("Fix local session storage")
		shell.CreateFile("src/authentication.go", "package main")
		shell.CreateFile("src/session.go", "package main")
		shell.UpdateFileAndAdd("src/users.go", usersFileContent)
		shell.EmptyCommit("Stop using shims")
		shell.UpdateFileAndAdd("src/authentication.go", "package authentication")
		shell.UpdateFileAndAdd("src/session.go", "package session")
		shell.Commit("Enhance user authentication feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Remove a line from an old commit")
		t.Wait(1000)

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("Stop using shims")).
			Wait(1000).
			PressEnter().
			Tap(func() {
				t.Views().CommitFiles().
					IsFocused().
					NavigateToLine(Contains("users.go")).
					Wait(1000).
					PressEnter().
					Tap(func() {
						t.Views().PatchBuilding().
							IsFocused().
							NavigateToLine(Contains("TODO")).
							Wait(500).
							PressPrimaryAction().
							PressEscape()
					}).
					Press(keys.Universal.CreatePatchOptionsMenu).
					Tap(func() {
						t.ExpectPopup().Menu().
							Title(Equals("Patch options")).
							Select(Contains("Remove patch from original commit")).
							Wait(500).
							Confirm()
					}).
					PressEscape()
			})
	},
})
