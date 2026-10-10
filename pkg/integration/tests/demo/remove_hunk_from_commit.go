package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var greetingFileContent = `package main

import "fmt"

func greet(name string) {
	fmt.Println("hello", name)
}

func farewell(name string) {
	fmt.Println("bye", name)
}
`

var greetingFileContentWithDebugLine = `package main

import "fmt"

func greet(name string) {
	fmt.Println("hello there", name)
}

func farewell(name string) {
	fmt.Println("DEBUG: saying bye to", name)
	fmt.Println("bye", name)
}
`

var RemoveHunkFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a hunk from an old commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(cfg *config.AppConfig) {
		setDefaultDemoConfig(cfg)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/user-authentication")
		shell.CreateFileAndAdd("src/greeting.go", greetingFileContent)
		shell.Commit("Greet users when they sign in")
		shell.UpdateFileAndAdd("src/greeting.go", greetingFileContentWithDebugLine)
		shell.Commit("Make the greeting friendlier")
		shell.EmptyCommit("Expire sessions after a day")
		shell.EmptyCommit("Enhance user authentication feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Remove a hunk from an old commit")
		t.Wait(1000)

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("Make the greeting friendlier")).
			Wait(1000).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Wait(1000).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains(`fmt.Println("DEBUG: saying bye to", name)`),
			).
			Wait(1000).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Discard lines from commit")).
					Content(AnyString()).
					Wait(1000).
					Confirm()
			}).
			Wait(1000).
			Content(Contains(`fmt.Println("hello there", name)`).DoesNotContain("DEBUG"))
	},
})
