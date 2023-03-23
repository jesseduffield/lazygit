package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) contextTree() *context.ContextTree {
	contextCommon := &context.ContextCommon{
		IGuiCommon: gui.c.IGuiCommon,
		Common:     gui.c.Common,
	}
	return context.NewContextTree(contextCommon)
}

// using this wrapper for when an onFocus function doesn't care about any potential
// props that could be passed
func OnFocusWrapper(f func() error) func(opts types.OnFocusOpts) error {
	return func(opts types.OnFocusOpts) error {
		return f()
	}
}

func (gui *Gui) popupViewNames() []string {
	popups := slices.Filter(gui.State.Contexts.Flatten(), func(c types.Context) bool {
		return c.GetKind() == types.PERSISTENT_POPUP || c.GetKind() == types.TEMPORARY_POPUP
	})

	return slices.Map(popups, func(c types.Context) string {
		return c.GetViewName()
	})
}

func (gui *Gui) defaultSideContext() types.Context {
	if gui.State.Modes.Filtering.Active() {
		return gui.State.Contexts.LocalCommits
	} else {
		return gui.State.Contexts.Files
	}
}

func (gui *Gui) TransientContexts() []types.Context {
	return slices.Filter(gui.State.Contexts.Flatten(), func(context types.Context) bool {
		return context.IsTransient()
	})
}

func (gui *Gui) getListContexts() []types.IListContext {
	return []types.IListContext{
		gui.State.Contexts.Menu,
		gui.State.Contexts.Files,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Suggestions,
	}
}

func (gui *Gui) getPatchExplorerContexts() []types.IPatchExplorerContext {
	return []types.IPatchExplorerContext{
		gui.State.Contexts.Staging,
		gui.State.Contexts.StagingSecondary,
		gui.State.Contexts.CustomPatchBuilder,
	}
}
