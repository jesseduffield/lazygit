package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() *context.MenuContext {
	return context.NewMenuContext(
		newGuiContextStateFetcher(gui, context.MENU_CONTEXT_KEY),
		gui.Views.Menu,
		gui.c,
		gui.getMenuOptions,
		func(content string) {
			gui.Views.Tooltip.SetContent(content)
		},
	)
}

func (gui *Gui) filesListContext() *context.WorkingTreeContext {
	return context.NewWorkingTreeContext(
		func() []*models.File { return gui.State.Model.Files },
		gui.Views.Files,
		func(startIdx int, length int) [][]string {
			lines := presentation.RenderFileTree(gui.State.Contexts.Files.FileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.State.Model.Submodules)
			return slices.Map(lines, func(line string) []string {
				return []string{line}
			})
		},
		OnFocusWrapper(gui.onFocusFile),
		OnFocusWrapper(gui.withDiffModeCheck(gui.filesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchesListContext() *context.BranchesContext {
	return context.NewBranchesContext(
		func() []*models.Branch { return gui.State.Model.Branches },
		newGuiContextStateFetcher(gui, context.LOCAL_BRANCHES_CONTEXT_KEY),
		gui.Views.Branches,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.branchesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) remotesListContext() *context.RemotesContext {
	return context.NewRemotesContext(
		func() []*models.Remote { return gui.State.Model.Remotes },
		newGuiContextStateFetcher(gui, context.REMOTES_CONTEXT_KEY),
		gui.Views.Branches,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.remotesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) remoteBranchesListContext() *context.RemoteBranchesContext {
	return context.NewRemoteBranchesContext(
		func() []*models.RemoteBranch { return gui.State.Model.RemoteBranches },
		newGuiContextStateFetcher(gui, context.REMOTE_BRANCHES_CONTEXT_KEY),
		gui.Views.RemoteBranches,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.remoteBranchesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) withDiffModeCheck(f func() error) func() error {
	return func() error {
		if gui.State.Modes.Diffing.Active() {
			return gui.renderDiff()
		}

		return f()
	}
}

func (gui *Gui) tagsListContext() *context.TagsContext {
	return context.NewTagsContext(
		func() []*models.Tag { return gui.State.Model.Tags },
		newGuiContextStateFetcher(gui, context.TAGS_CONTEXT_KEY),
		gui.Views.Branches,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.tagsRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchCommitsListContext() *context.LocalCommitsContext {
	return context.NewLocalCommitsContext(
		// TODO: standardise naming for branch commits vs local commits
		func() []*models.Commit { return gui.State.Model.Commits },
		newGuiContextStateFetcher(gui, context.LOCAL_COMMITS_CONTEXT_KEY),
		gui.Views.Commits,
		OnFocusWrapper(gui.onCommitFocus),
		OnFocusWrapper(gui.withDiffModeCheck(gui.branchCommitsRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) subCommitsListContext() *context.SubCommitsContext {
	return context.NewSubCommitsContext(
		func() []*models.Commit { return gui.State.Model.SubCommits },
		newGuiContextStateFetcher(gui, context.SUB_COMMITS_CONTEXT_KEY),
		gui.Views.SubCommits,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.subCommitsRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) reflogCommitsListContext() *context.ReflogCommitsContext {
	return context.NewReflogCommitsContext(
		func() []*models.Commit { return gui.State.Model.FilteredReflogCommits },
		newGuiContextStateFetcher(gui, context.REFLOG_COMMITS_CONTEXT_KEY),
		gui.Views.Commits,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.reflogCommitsRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) stashListContext() *context.StashContext {
	return context.NewStashContext(
		func() []*models.StashEntry { return gui.State.Model.StashEntries },
		newGuiContextStateFetcher(gui, context.STASH_CONTEXT_KEY),
		gui.Views.Stash,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.stashRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) commitFilesListContext() *context.CommitFilesContext {
	return context.NewCommitFilesContext(
		func() []*models.CommitFile { return gui.State.Model.CommitFiles },
		gui.Views.CommitFiles,
		func(startIdx int, length int) [][]string {
			if gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.Len() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := presentation.RenderCommitFileTree(gui.State.Contexts.CommitFiles.CommitFileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.git.Patch.PatchManager)
			return slices.Map(lines, func(line string) []string {
				return []string{line}
			})
		},
		OnFocusWrapper(gui.onCommitFileFocus),
		OnFocusWrapper(gui.withDiffModeCheck(gui.commitFilesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) submodulesListContext() *context.SubmodulesContext {
	return context.NewSubmodulesContext(
		func() []*models.SubmoduleConfig { return gui.State.Model.Submodules },
		newGuiContextStateFetcher(gui, context.SUBMODULES_CONTEXT_KEY),
		gui.Views.Files,
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.submodulesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) suggestionsListContext() *context.SuggestionsContext {
	return context.NewSuggestionsContext(
		func() []*types.Suggestion { return gui.State.Suggestions },
		newGuiContextStateFetcher(gui, context.SUGGESTIONS_CONTEXT_KEY),
		gui.Views.Suggestions,
		nil,
		nil,
		nil,
		gui.c,
	)
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
