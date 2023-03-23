package gui

import (
	"log"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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
	return context.NewRemotesContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Model.Remotes, gui.State.Modes.Diffing.Ref)
		},
		gui.c,
	)
}

func (gui *Gui) remoteBranchesListContext() *context.RemoteBranchesContext {
	return context.NewRemoteBranchesContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.Model.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		gui.c,
	)
}

func (gui *Gui) tagsListContext() *context.TagsContext {
	return context.NewTagsContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Model.Tags, gui.State.Modes.Diffing.Ref)
		},
		gui.c,
	)
}

func (gui *Gui) branchCommitsListContext() *context.LocalCommitsContext {
	return context.NewLocalCommitsContext(
		func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.c.CurrentContext().GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.State.Contexts.LocalCommits.GetSelected()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}

			showYouAreHereLabel := gui.State.Model.WorkingTreeStateAtLastCommitRefresh == enums.REBASE_MODE_REBASING

			return presentation.GetCommitListDisplayStrings(
				gui.Common,
				gui.State.Model.Commits,
				gui.State.ScreenMode != types.SCREEN_NORMAL,
				gui.c.Modes().CherryPicking.SelectedShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				gui.State.Model.BisectInfo,
				showYouAreHereLabel,
			)
		},
		gui.c,
	)
}

func (gui *Gui) subCommitsListContext() *context.SubCommitsContext {
	return context.NewSubCommitsContext(
		func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.c.CurrentContext().GetKey() == context.SUB_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.State.Contexts.SubCommits.GetSelected()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.Common,
				gui.State.Model.SubCommits,
				gui.State.ScreenMode != types.SCREEN_NORMAL,
				gui.c.Modes().CherryPicking.SelectedShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				git_commands.NewNullBisectInfo(),
				false,
			)
		},
		gui.c,
	)
}

func (gui *Gui) shouldShowGraph() bool {
	if gui.State.Modes.Filtering.Active() {
		return false
	}

	value := gui.c.UserConfig.Git.Log.ShowGraph
	switch value {
	case "always":
		return true
	case "never":
		return false
	case "when-maximised":
		return gui.State.ScreenMode != types.SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}

func (gui *Gui) reflogCommitsListContext() *context.ReflogCommitsContext {
	return context.NewReflogCommitsContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.Model.FilteredReflogCommits,
				gui.State.ScreenMode != types.SCREEN_NORMAL,
				gui.c.Modes().CherryPicking.SelectedShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
			)
		},
		gui.c,
	)
}

func (gui *Gui) stashListContext() *context.StashContext {
	return context.NewStashContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.Model.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		gui.c,
	)
}

func (gui *Gui) commitFilesListContext() *context.CommitFilesContext {
	return context.NewCommitFilesContext(
		func(startIdx int, length int) [][]string {
			if gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.Len() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := presentation.RenderCommitFileTree(gui.State.Contexts.CommitFiles.CommitFileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.git.Patch.PatchBuilder)
			return slices.Map(lines, func(line string) []string {
				return []string{line}
			})
		},
		gui.c,
	)
}

func (gui *Gui) submodulesListContext() *context.SubmodulesContext {
	return context.NewSubmodulesContext(
		func(startIdx int, length int) [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Model.Submodules)
		},
		gui.c,
	)
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
