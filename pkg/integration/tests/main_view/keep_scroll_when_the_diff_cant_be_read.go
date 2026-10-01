package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepScrollWhenTheDiffCantBeRead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Changing the context size under a diff renderer whose rows can't be placed keeps the scroll position rather than jumping to the top",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// A renderer that says nothing about which line of which file each row shows,
		// and mangles the diff enough that it can't be read back as one either: no line
		// of it can be looked for in the re-render.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "opaque", Command: `sed -e 's/^/| /'`},
		}
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		for _, i := range []int{5, 15, 25, 35} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain)

		t.Views().Main().
			Content(Contains("| +LINE05")).
			OriginY(6).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.IncreaseContextInDiffView)
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			// The re-render is a different command, and nothing in its output can be
			// matched up with what was on screen, so the offset is all there is to keep —
			// and it is a good deal closer than the top.
			OriginY(6)
	},
})
