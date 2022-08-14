package gui

import (
	"log"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() *context.MenuContext {
	return context.NewMenuContext(
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
		nil,
		gui.withDiffModeCheck(gui.filesRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchesListContext() *context.BranchesContext {
	return context.NewBranchesContext(
		func() []*models.Branch { return gui.State.Model.Branches },
		gui.Views.Branches,
		func(startIdx int, length int) [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Model.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref, gui.Tr)
		},
		nil,
		gui.withDiffModeCheck(gui.branchesRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) remotesListContext() *context.RemotesContext {
	return context.NewRemotesContext(
		func() []*models.Remote { return gui.State.Model.Remotes },
		gui.Views.Remotes,
		func(startIdx int, length int) [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Model.Remotes, gui.State.Modes.Diffing.Ref)
		},
		nil,
		gui.withDiffModeCheck(gui.remotesRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) remoteBranchesListContext() *context.RemoteBranchesContext {
	return context.NewRemoteBranchesContext(
		func() []*models.RemoteBranch { return gui.State.Model.RemoteBranches },
		gui.Views.RemoteBranches,
		func(startIdx int, length int) [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.Model.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		nil,
		gui.withDiffModeCheck(gui.remoteBranchesRenderToMain),
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
		gui.Views.Tags,
		func(startIdx int, length int) [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Model.Tags, gui.State.Modes.Diffing.Ref)
		},
		nil,
		gui.withDiffModeCheck(gui.tagsRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchCommitsListContext() *context.LocalCommitsContext {
	return context.NewLocalCommitsContext(
		func() []*models.Commit { return gui.State.Model.Commits },
		gui.Views.Commits,
		func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.State.Contexts.LocalCommits.GetSelected()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.Model.Commits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.CherryPick.CherryPickedCommitShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				gui.State.Model.BisectInfo,
			)
		},
		OnFocusWrapper(gui.onCommitFocus),
		gui.withDiffModeCheck(gui.branchCommitsRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) subCommitsListContext() *context.SubCommitsContext {
	return context.NewSubCommitsContext(
		func() []*models.Commit { return gui.State.Model.SubCommits },
		gui.Views.SubCommits,
		func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == context.SUB_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.State.Contexts.SubCommits.GetSelected()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.Model.SubCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.CherryPick.CherryPickedCommitShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				git_commands.NewNullBisectInfo(),
			)
		},
		nil,
		gui.withDiffModeCheck(gui.subCommitsRenderToMain),
		nil,
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
		return gui.State.ScreenMode != SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}

func (gui *Gui) reflogCommitsListContext() *context.ReflogCommitsContext {
	return context.NewReflogCommitsContext(
		func() []*models.Commit { return gui.State.Model.FilteredReflogCommits },
		gui.Views.ReflogCommits,
		func(startIdx int, length int) [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.Model.FilteredReflogCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.CherryPick.CherryPickedCommitShaSet(),
				gui.State.Modes.Diffing.Ref,
				gui.c.UserConfig.Gui.TimeFormat,
				gui.c.UserConfig.Git.ParseEmoji,
			)
		},
		nil,
		gui.withDiffModeCheck(gui.reflogCommitsRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) stashListContext() *context.StashContext {
	return context.NewStashContext(
		func() []*models.StashEntry { return gui.State.Model.StashEntries },
		gui.Views.Stash,
		func(startIdx int, length int) [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.Model.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		nil,
		gui.withDiffModeCheck(gui.stashRenderToMain),
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
		nil,
		gui.withDiffModeCheck(gui.commitFilesRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) submodulesListContext() *context.SubmodulesContext {
	return context.NewSubmodulesContext(
		func() []*models.SubmoduleConfig { return gui.State.Model.Submodules },
		gui.Views.Submodules,
		func(startIdx int, length int) [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Model.Submodules)
		},
		nil,
		gui.withDiffModeCheck(gui.submodulesRenderToMain),
		nil,
		gui.c,
	)
}

func (gui *Gui) suggestionsListContext() *context.SuggestionsContext {
	return context.NewSuggestionsContext(
		func() []*types.Suggestion { return gui.State.Suggestions },
		gui.Views.Suggestions,
		func(startIdx int, length int) [][]string {
			return presentation.GetSuggestionListDisplayStrings(gui.State.Suggestions)
		},
		nil,
		nil,
		func(types.OnFocusLostOpts) error {
			gui.deactivateConfirmationPrompt()
			return nil
		},
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
