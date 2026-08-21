package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RawFallbackUnderAnExternalDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view under an external diff that says nothing about its rows brings git's own diff, keeping the scroll position",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// An external diff whose output has nothing to do with the diff it was given,
		// let alone anything to say about which line of which file each row shows. It
		// is long enough to be scrolled about in.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "opaque", Type: "extDiff", Command: `sh -c 'seq -f "EXT-%g" 40'`},
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
		// Browsing shows what the renderer produced, whatever that is.
		t.Views().Main().
			Content(Contains("EXT-1")).
			Tap(func() {
				t.Views().Files().
					IsFocused().
					Press(keys.Universal.ScrollDownMain).
					Press(keys.Universal.ScrollDownMain).
					Press(keys.Universal.ScrollDownMain)
			}).
			OriginY(6).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.FocusMainView)
			}).
			IsFocused().
			// Focusing the view to act on it brings git's own diff instead — the
			// renderer's rows can't be placed in the file — and leaves the view at the
			// offset it was at, rather than at the top.
			Content(Contains("+LINE05")).
			Content(DoesNotContain("EXT-1")).
			OriginY(6).
			SelectionIsActive()
	},
})
