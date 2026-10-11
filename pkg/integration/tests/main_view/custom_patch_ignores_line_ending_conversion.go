package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomPatchIgnoresLineEndingConversion = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The custom patch is shown as git states it on a machine that checks files out with CRLF",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// The patch is materialized into trees of its own outside the repo, so what the repo
	// says about line endings never reaches the commands over them. These variables are
	// how a setting reaches a git command wherever it runs.
	ExtraEnvVars: map[string]string{
		"GIT_CONFIG_COUNT":   "1",
		"GIT_CONFIG_KEY_0":   "core.autocrlf",
		"GIT_CONFIG_VALUE_0": "true",
	},
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction()

		// The patch is one line and its context, so that is all the two trees differ in.
		// A tree written in the form the platform checks files out in would differ from
		// the other in every line instead.
		t.Views().Secondary().
			ContainsLines(
				Contains(" one"),
				Contains("-two"),
				Contains(" three"),
			).
			// git says nothing about a round trip these files never take.
			Content(DoesNotContain("warning"))
	},
})
