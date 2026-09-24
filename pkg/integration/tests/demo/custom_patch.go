package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var usersFileContent = `package main

import "fmt"

func main() {
	fmt.Println(shims.Greeting())
	serve()
}

func serve() {
	fmt.Println("listening on :8080")
}
`

var usersFileContentWithLogging = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
	serve()
}

func logRequest(path string) {
	fmt.Println("request:", path)
}

func serve() {
	fmt.Println("listening on :8080")
}
`

var CustomPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Split a change out of an old commit into a new one",
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
		shell.CreateFileAndAdd("src/users.go", usersFileContent)
		shell.Commit("Fix local session storage")
		shell.CreateFile("src/authentication.go", "package main")
		shell.CreateFile("src/session.go", "package main")
		shell.UpdateFileAndAdd("src/users.go", usersFileContentWithLogging)
		shell.Commit("Stop using shims")
		shell.UpdateFileAndAdd("src/authentication.go", "package authentication")
		shell.UpdateFileAndAdd("src/session.go", "package session")
		shell.Commit("Enhance user authentication feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Split a change out of an old commit")
		t.Wait(1000)

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("Stop using shims")).
			Wait(1000).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Wait(1000).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+func logRequest(path string) {"),
				Contains(`fmt.Println("request:", path)`),
				Contains("+}"),
				Equals("+"),
			).
			Wait(500).
			SetCaptionPrefix("Add the hunk to a custom patch").
			PressPrimaryAction().
			Wait(1000).
			SetCaptionPrefix("Move the patch into a new commit").
			Press(keys.Universal.CreatePatchOptionsMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Patch options")).
					Select(Contains("Move patch into new commit after the original commit")).
					Wait(500).
					Confirm()

				t.ExpectPopup().CommitMessagePanel().
					Type("Add request logging").
					Confirm()
			})

		t.Views().Commits().
			IsFocused().
			TopLines(
				Contains("Enhance user authentication feature"),
				Contains("Add request logging"),
				Contains("Stop using shims"),
			)
	},
})
