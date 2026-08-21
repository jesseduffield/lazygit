package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

// Falling back to git's own diff when the configured one can't be acted on.
//
// A diff renderer is free to lay a diff out however it likes, and once it has, the only
// way to know which line of which file a row shows is for the renderer to say so. One
// that doesn't produces a diff that can be read but not staged, edited or copied from —
// so when the user focuses the main view to act on it, we show git's own diff instead.
// Browsing keeps the renderer's version; only acting on it needs one we can follow.

// MainViewDiffMode says how a side panel should produce the diff it renders into the
// main view: as the user configured it, or as git's own — while the main view holds
// focus and what the renderer would produce couldn't be acted on.
//
// Every panel that renders a diff into the main view asks, so that a re-render while
// focused — after staging a hunk, say — stays with git's own diff rather than flipping
// back to the renderer's.
func (self *DiffLineHelper) MainViewDiffMode() git_commands.DiffMode {
	if self.mainViewIsFocused() && self.diffNeedsMetadata() && !self.diffRendererEmitsMetadata() {
		return git_commands.DiffModeRaw
	}
	return git_commands.DiffModeRendered
}

// RenderFocusedMainViewAgain has the panel beneath the focused main view render its
// diff again — which, the main view now holding focus, is git's own diff rather than
// the renderer's — and calls place once that is on screen.
//
// The whole diff is read before it is shown, rather than the first screenful: place
// looks at what is there to decide where to put the selection, and a change line
// further down would otherwise be missed.
//
// It is the same diff of the same files, so the view keeps the scroll position it has
// rather than starting from the top: git lays the changes out its own way and the line
// the user was on is somewhere else now, but the offset still puts them among the same
// part of the file — and the selection is then established from what that leaves on
// screen.
func (self *DiffLineHelper) RenderFocusedMainViewAgain(view *gocui.View, sidePanel types.Context, place func()) {
	manager := self.c.GetOrCreateViewBufferManagerForView(view)
	if manager == nil {
		return
	}

	manager.SetKeepScrollPositionForNextTask()
	manager.SetRestoreForNextTask(&tasks.RenderRestore{
		FirstPaintReady: func() bool { return false },
		Apply: func(swapIn func()) {
			swapIn()
			place()
		},
	})

	sidePanel.HandleRenderToMain()
}

func (self *DiffLineHelper) mainViewIsFocused() bool {
	current := self.c.Context().CurrentStatic().GetKey()
	return current == self.c.Contexts().Normal.GetKey() ||
		current == self.c.Contexts().NormalSecondary.GetKey()
}

// diffNeedsMetadata reports whether the diff we would show is one whose rows can only
// be placed in the file by the records the renderer states. Any custom renderer may
// restructure the diff; so may git itself, once the renderer's arguments ask for a word
// diff, whose markup is inline. Plain git output describes itself, and needs no records.
func (self *DiffLineHelper) diffNeedsMetadata() bool {
	manager := self.c.State().GetDiffRendererConfigManager()
	if manager.GetDiffRendererType() != config.DiffRendererType_RawGit {
		return true
	}
	return len(manager.GetRawGitArgs()) > 0
}

// diffRendererEmitsMetadata is the probed verdict about the current diff renderer, asked
// once and remembered until the renderer changes — the user cycling to another one, or a
// changed config being reloaded.
func (self *DiffLineHelper) diffRendererEmitsMetadata() bool {
	signature := self.diffRendererSignature()
	if self.rendererEmitsMetadata == nil || signature != self.rendererSignature {
		verdict := self.c.Git().Diff.ProbeDiffRendererEmitsMetadata()
		self.rendererEmitsMetadata = &verdict
		self.rendererSignature = signature
	}
	return *self.rendererEmitsMetadata
}

// diffRendererSignature is what identifies the current diff renderer, so that the
// remembered verdict is dropped when it stops describing the renderer we have. The width
// a command is asked for is no part of its identity, so a fixed one is used.
func (self *DiffLineHelper) diffRendererSignature() string {
	manager := self.c.State().GetDiffRendererConfigManager()
	index, _ := manager.CurrentDiffRendererIndex()
	return fmt.Sprintf("%d\x00%s\x00%s\x00%s",
		index,
		manager.GetExternalDiffCommand(3),
		manager.GetStdinFilterCommand(0),
		strings.Join(manager.GetRawGitArgs(), "\x00"))
}
