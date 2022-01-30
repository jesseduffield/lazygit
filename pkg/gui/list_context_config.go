package gui

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:        "menu",
			Key:             "menu",
			Kind:            types.PERSISTENT_POPUP,
			OnGetOptionsMap: gui.getMenuOptions,
		}),
		GetItemsLength:  func() int { return gui.Views.Menu.LinesHeight() },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Menu },
		Gui:             gui,

		// no GetDisplayStrings field because we do a custom render on menu creation
	}
}

func (gui *Gui) filesListContext() *context.WorkingTreeContext {
	return context.NewWorkingTreeContext(
		func() []*models.File { return gui.State.Files },
		func() *gocui.View { return gui.Views.Files },
		func(startIdx int, length int) [][]string {
			lines := presentation.RenderFileTree(gui.State.Contexts.Files.FileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.State.Submodules)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		OnFocusWrapper(gui.onFocusFile),
		OnFocusWrapper(gui.withDiffModeCheck(gui.filesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchesListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        context.LOCAL_BRANCHES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.Branches) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Branches },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.branchesRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedBranch()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
}

func (gui *Gui) remotesListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        context.REMOTES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.Remotes) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Remotes },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.remotesRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Modes.Diffing.Ref)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedRemote()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
}

func (gui *Gui) remoteBranchesListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        context.REMOTE_BRANCHES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.RemoteBranches) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.RemoteBranches },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.remoteBranchesRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedRemoteBranch()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
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
		func() []*models.Tag { return gui.State.Tags },
		func() *gocui.View { return gui.Views.Branches },
		func(startIdx int, length int) [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Tags, gui.State.Modes.Diffing.Ref)
		},
		nil,
		OnFocusWrapper(gui.withDiffModeCheck(gui.tagsRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) branchCommitsListContext() types.IListContext {
	parseEmoji := gui.c.UserConfig.Git.ParseEmoji
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        context.BRANCH_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.Commits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Commits },
		OnFocus:         OnFocusWrapper(gui.onCommitFocus),
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.branchCommitsRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == context.BRANCH_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedLocalCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.Commits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.cherryPick.CherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				gui.State.BisectInfo,
			)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedLocalCommit()
			if item == nil {
				return ""
			}
			return item.ID()
		},
		RenderSelection: true,
	}
}

func (gui *Gui) subCommitsListContext() types.IListContext {
	parseEmoji := gui.c.UserConfig.Git.ParseEmoji
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        context.SUB_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.SubCommits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.SubCommits },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.subCommitsRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == context.SUB_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedSubCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.SubCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.cherryPick.CherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
				git_commands.NewNullBisectInfo(),
			)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedSubCommit()
			if item == nil {
				return ""
			}
			return item.ID()
		},
		RenderSelection: true,
	}
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

func (gui *Gui) reflogCommitsListContext() types.IListContext {
	parseEmoji := gui.c.UserConfig.Git.ParseEmoji
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        context.REFLOG_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.FilteredReflogCommits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.ReflogCommits },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.reflogCommitsRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.FilteredReflogCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.helpers.cherryPick.CherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
			)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedReflogCommit()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
}

func (gui *Gui) stashListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "stash",
			WindowName: "stash",
			Key:        context.STASH_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.StashEntries) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Stash },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.stashRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedStashEntry()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
}

func (gui *Gui) commitFilesListContext() *context.CommitFilesContext {
	return context.NewCommitFilesContext(
		func() []*models.CommitFile { return gui.State.CommitFiles },
		func() *gocui.View { return gui.Views.CommitFiles },
		func(startIdx int, length int) [][]string {
			if gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.GetItemsLength() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := presentation.RenderCommitFileTree(gui.State.Contexts.CommitFiles.CommitFileTreeViewModel, gui.State.Modes.Diffing.Ref, gui.git.Patch.PatchManager)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		OnFocusWrapper(gui.onCommitFileFocus),
		OnFocusWrapper(gui.withDiffModeCheck(gui.commitFilesRenderToMain)),
		nil,
		gui.c,
	)
}

func (gui *Gui) submodulesListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "files",
			WindowName: "files",
			Key:        context.SUBMODULES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		}),
		GetItemsLength:  func() int { return len(gui.State.Submodules) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Submodules },
		OnRenderToMain:  OnFocusWrapper(gui.withDiffModeCheck(gui.submodulesRenderToMain)),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Submodules)
		},
		OnGetSelectedItemId: func() string {
			item := gui.getSelectedSubmodule()
			if item == nil {
				return ""
			}
			return item.ID()
		},
	}
}

func (gui *Gui) suggestionsListContext() types.IListContext {
	return &ListContext{
		BaseContext: context.NewBaseContext(context.NewBaseContextOpts{
			ViewName:   "suggestions",
			WindowName: "suggestions",
			Key:        context.SUGGESTIONS_CONTEXT_KEY,
			Kind:       types.PERSISTENT_POPUP,
		}),
		GetItemsLength:  func() int { return len(gui.State.Suggestions) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Suggestions },
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSuggestionListDisplayStrings(gui.State.Suggestions)
		},
	}
}

func (gui *Gui) getListContexts() []types.IListContext {
	return []types.IListContext{
		gui.State.Contexts.Menu,
		gui.State.Contexts.Files,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Suggestions,
	}
}
