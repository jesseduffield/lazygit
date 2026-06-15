package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// sidePanelViewNames maps each gui.sidePanels name to the gocui view it controls.
// A panel's window name is the name of its first tab, so for a panel's first tab
// this also gives the default view of its window. The keys must match
// config.ValidSidePanelTabs (enforced by a test).
var sidePanelViewNames = map[string]string{
	"status":     "status",
	"files":      "files",
	"worktrees":  "worktrees",
	"submodules": "submodules",
	"branches":   "localBranches",
	"remotes":    "remotes",
	"tags":       "tags",
	"commits":    "commits",
	"reflog":     "reflogCommits",
	"stash":      "stash",
}

// sidePanelTabTitles maps each gui.sidePanels name to the title shown on its tab.
func (gui *Gui) sidePanelTabTitles() map[string]string {
	tr := gui.c.Tr
	return map[string]string{
		"status":     tr.StatusTitle,
		"files":      tr.FilesTitle,
		"worktrees":  tr.WorktreesTitle,
		"submodules": tr.SubmodulesTitle,
		"branches":   tr.LocalBranchesTitle,
		"remotes":    tr.RemotesTitle,
		"tags":       tr.TagsTitle,
		"commits":    tr.CommitsTitle,
		"reflog":     tr.ReflogCommitsTitle,
		"stash":      tr.StashTitle,
	}
}

// sidePanelContexts maps each gui.sidePanels name to the context it controls.
func sidePanelContexts(contextTree *context.ContextTree) map[string]types.Context {
	return map[string]types.Context{
		"status":     contextTree.Status,
		"files":      contextTree.Files,
		"worktrees":  contextTree.Worktrees,
		"submodules": contextTree.Submodules,
		"branches":   contextTree.Branches,
		"remotes":    contextTree.Remotes,
		"tags":       contextTree.Tags,
		"commits":    contextTree.LocalCommits,
		"reflog":     contextTree.ReflogCommits,
		"stash":      contextTree.Stash,
	}
}

// applySidePanelConfig (re)assigns each side context's window and resets each
// window's default view from the current gui.sidePanels config. It runs against
// the current repo's contexts, so gui.State must already be set. We call it on
// every repo entry (a repo's per-repo config can differ from the previous one's)
// and on a live config reload.
func (gui *Gui) applySidePanelConfig() {
	contextTree := gui.State.Contexts
	gui.assignSidePanelWindows(contextTree)
	gui.State.WindowViewNameMap = gui.initialWindowViewNameMap(contextTree)
}

// moveDefaultTabsToTop brings each panel's first configured tab to the top of
// its window, so the configured default tab is the one shown when a panel hasn't
// been focused yet (the view z-order is otherwise set from a fixed list that
// need not match the configured tab order).
func (gui *Gui) moveDefaultTabsToTop() {
	contexts := sidePanelContexts(gui.State.Contexts)
	for _, panel := range gui.c.UserConfig().Gui.SidePanels {
		gui.helpers.Window.MoveToTopOfWindow(contexts[panel[0]])
	}
}

// reloadSidePanels re-applies the side panel config to the current repo after a
// live config reload: it reassigns windows and default views, restores each
// panel's default tab, and keeps the focused panel in a consistent state.
func (gui *Gui) reloadSidePanels() {
	gui.applySidePanelConfig()
	gui.moveDefaultTabsToTop()

	// applySidePanelConfig reset every window to show its first configured tab,
	// which would leave the focused tab hidden behind its panel's default tab
	// (the panel would look unfocused even though its tab is selected). Re-focus
	// the current context so its tab stays shown and highlighted. If the new
	// config has hidden the focused panel entirely, move focus to the default
	// side panel instead.
	current := gui.c.Context().Current()
	if current.GetKind() != types.SIDE_CONTEXT {
		return
	}

	if lo.Contains(gui.helpers.Window.SideWindows(), current.GetWindowName()) {
		gui.c.Context().Activate(current, types.OnFocusOpts{})
	} else {
		gui.c.Context().Push(gui.defaultSideContext(), types.OnFocusOpts{})
	}
}

// assignSidePanelWindows sets each side context's window name from the config so
// that contexts grouped into one panel share a window (the window name being the
// panel's first tab). Side panels the user hasn't listed get their own window
// name; since the layout produces no dimensions for those windows, their views
// stay hidden rather than overlapping a visible panel.
func (gui *Gui) assignSidePanelWindows(contextTree *context.ContextTree) {
	contexts := sidePanelContexts(contextTree)
	assigned := make(map[string]bool, len(contexts))

	for _, panel := range gui.c.UserConfig().Gui.SidePanels {
		windowName := panel[0]
		for _, name := range panel {
			contexts[name].SetWindowName(windowName)
			assigned[name] = true
		}
	}

	for name, ctx := range contexts {
		if !assigned[name] {
			ctx.SetWindowName(name)
		}
	}
}
