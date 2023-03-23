package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() *context.MenuContext {
	return context.NewMenuContext(gui.c)
}

func (gui *Gui) filesListContext() *context.WorkingTreeContext {
	return context.NewWorkingTreeContext(gui.c)
}

func (gui *Gui) branchesListContext() *context.BranchesContext {
	return context.NewBranchesContext(gui.c)
}

func (gui *Gui) remotesListContext() *context.RemotesContext {
	return context.NewRemotesContext(gui.c)
}

func (gui *Gui) remoteBranchesListContext() *context.RemoteBranchesContext {
	return context.NewRemoteBranchesContext(gui.c)
}

func (gui *Gui) tagsListContext() *context.TagsContext {
	return context.NewTagsContext(gui.c)
}

func (gui *Gui) branchCommitsListContext() *context.LocalCommitsContext {
	return context.NewLocalCommitsContext(gui.c)
}

func (gui *Gui) subCommitsListContext() *context.SubCommitsContext {
	return context.NewSubCommitsContext(gui.c)
}

func (gui *Gui) reflogCommitsListContext() *context.ReflogCommitsContext {
	return context.NewReflogCommitsContext(gui.c)
}

func (gui *Gui) stashListContext() *context.StashContext {
	return context.NewStashContext(gui.c)
}

func (gui *Gui) commitFilesListContext() *context.CommitFilesContext {
	return context.NewCommitFilesContext(gui.c)
}

func (gui *Gui) submodulesListContext() *context.SubmodulesContext {
	return context.NewSubmodulesContext(gui.c)
}

func (gui *Gui) suggestionsListContext() *context.SuggestionsContext {
	return context.NewSuggestionsContext(gui.c)
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
