package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CheckForConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Run a command and check for conflicts after",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "m",
				Context: "localBranches",
				Command: "git merge {{ .SelectedLocalBranch.Name | quote }}",
				After: config.CustomCommandAfterHook{
					CheckForConflicts: true,
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			TopLines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
			).
			NavigateToLine(Contains("second-change-branch")).
			Press("m")

		t.Common().AcknowledgeConflicts()
	},
})
