package custom_commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Export the selected diff from staging and patch building to custom commands",
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     config.Keybinding{"X"},
				Context: "global",
				Command: "printf '%s' {{.SelectedDiff | quote}} > .git/selected-diff",
			},
			{
				Key:     config.Keybinding{"Y"},
				Context: "staging",
				Command: "printf '%s' {{.SelectedDiff | quote}} > .git/selected-diff",
				Prompts: []config.CustomCommandPrompt{
					{Type: "input", Title: "Comment", Key: "Comment"},
				},
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "old\n")
		shell.Commit("initial")
		shell.UpdateFile("file1", "new 'quoted' $(echo substituted)\n"+strings.Repeat("long line ", 30)+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().PressEnter()
		t.Views().Staging().
			IsFocused().
			SelectedLines(Contains("-old")).
			Press(config.Keybinding{"X"})
		t.FileSystem().FileContent(".git/selected-diff", Equals("-old\n"))

		t.Views().Staging().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Press(config.Keybinding{"Y"})
		t.ExpectPopup().Prompt().Title(Equals("Comment")).Type("Please simplify").Confirm()
		t.FileSystem().FileContent(".git/selected-diff", Equals("-old\n+new 'quoted' $(echo substituted)\n"))

		t.Views().Staging().
			Press(keys.Main.ToggleSelectHunk).
			Press(config.Keybinding{"X"})
		t.FileSystem().FileContent(".git/selected-diff", Equals("-old\n+new 'quoted' $(echo substituted)\n+"+strings.Repeat("long line ", 30)+"\n"))

		t.Views().Staging().PressPrimaryAction()
		t.Views().StagingSecondary().
			IsFocused().
			Press(config.Keybinding{"X"})
		t.FileSystem().FileContent(".git/selected-diff", Equals("-old\n"))

		t.Views().StagingSecondary().PressEscape()
		t.Views().Files().IsFocused().Press(config.Keybinding{"X"})
		t.FileSystem().FileContent(".git/selected-diff", Equals(""))

		t.Views().Commits().Focus().PressEnter()
		t.Views().CommitFiles().IsFocused().PressEnter()
		t.Views().PatchBuilding().IsFocused().Press(config.Keybinding{"X"})
		t.FileSystem().FileContent(".git/selected-diff", Equals("+old\n"))
	},
})
