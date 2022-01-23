package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	STATUS_CONTEXT_KEY              types.ContextKey = "status"
	FILES_CONTEXT_KEY               types.ContextKey = "files"
	LOCAL_BRANCHES_CONTEXT_KEY      types.ContextKey = "localBranches"
	REMOTES_CONTEXT_KEY             types.ContextKey = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY     types.ContextKey = "remoteBranches"
	TAGS_CONTEXT_KEY                types.ContextKey = "tags"
	BRANCH_COMMITS_CONTEXT_KEY      types.ContextKey = "commits"
	REFLOG_COMMITS_CONTEXT_KEY      types.ContextKey = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY         types.ContextKey = "subCommits"
	COMMIT_FILES_CONTEXT_KEY        types.ContextKey = "commitFiles"
	STASH_CONTEXT_KEY               types.ContextKey = "stash"
	MAIN_NORMAL_CONTEXT_KEY         types.ContextKey = "normal"
	MAIN_MERGING_CONTEXT_KEY        types.ContextKey = "merging"
	MAIN_PATCH_BUILDING_CONTEXT_KEY types.ContextKey = "patchBuilding"
	MAIN_STAGING_CONTEXT_KEY        types.ContextKey = "staging"
	MENU_CONTEXT_KEY                types.ContextKey = "menu"
	CREDENTIALS_CONTEXT_KEY         types.ContextKey = "credentials"
	CONFIRMATION_CONTEXT_KEY        types.ContextKey = "confirmation"
	SEARCH_CONTEXT_KEY              types.ContextKey = "search"
	COMMIT_MESSAGE_CONTEXT_KEY      types.ContextKey = "commitMessage"
	SUBMODULES_CONTEXT_KEY          types.ContextKey = "submodules"
	SUGGESTIONS_CONTEXT_KEY         types.ContextKey = "suggestions"
	COMMAND_LOG_CONTEXT_KEY         types.ContextKey = "cmdLog"
)

var AllContextKeys = []types.ContextKey{
	STATUS_CONTEXT_KEY,
	FILES_CONTEXT_KEY,
	LOCAL_BRANCHES_CONTEXT_KEY,
	REMOTES_CONTEXT_KEY,
	REMOTE_BRANCHES_CONTEXT_KEY,
	TAGS_CONTEXT_KEY,
	BRANCH_COMMITS_CONTEXT_KEY,
	REFLOG_COMMITS_CONTEXT_KEY,
	SUB_COMMITS_CONTEXT_KEY,
	COMMIT_FILES_CONTEXT_KEY,
	STASH_CONTEXT_KEY,
	MAIN_NORMAL_CONTEXT_KEY,
	MAIN_MERGING_CONTEXT_KEY,
	MAIN_PATCH_BUILDING_CONTEXT_KEY,
	MAIN_STAGING_CONTEXT_KEY,
	MENU_CONTEXT_KEY,
	CREDENTIALS_CONTEXT_KEY,
	CONFIRMATION_CONTEXT_KEY,
	SEARCH_CONTEXT_KEY,
	COMMIT_MESSAGE_CONTEXT_KEY,
	SUBMODULES_CONTEXT_KEY,
	SUGGESTIONS_CONTEXT_KEY,
	COMMAND_LOG_CONTEXT_KEY,
}

func (gui *Gui) allContexts() []types.Context {
	return []types.Context{
		gui.State.Contexts.Status,
		gui.State.Contexts.Files,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.Menu,
		gui.State.Contexts.Confirmation,
		gui.State.Contexts.Credentials,
		gui.State.Contexts.CommitMessage,
		gui.State.Contexts.Normal,
		gui.State.Contexts.Staging,
		gui.State.Contexts.Merging,
		gui.State.Contexts.PatchBuilding,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Suggestions,
		gui.State.Contexts.CommandLog,
	}
}

func (gui *Gui) contextTree() context.ContextTree {
	return context.ContextTree{
		Status: &BasicContext{
			OnRenderToMain: OnFocusWrapper(gui.statusRenderToMain),
			Kind:           types.SIDE_CONTEXT,
			ViewName:       "status",
			Key:            STATUS_CONTEXT_KEY,
		},
		Files:          gui.filesListContext(),
		Submodules:     gui.submodulesListContext(),
		Menu:           gui.menuListContext(),
		Remotes:        gui.remotesListContext(),
		RemoteBranches: gui.remoteBranchesListContext(),
		BranchCommits:  gui.branchCommitsListContext(),
		CommitFiles:    gui.commitFilesListContext(),
		ReflogCommits:  gui.reflogCommitsListContext(),
		SubCommits:     gui.subCommitsListContext(),
		Branches:       gui.branchesListContext(),
		Tags:           gui.tagsListContext(),
		Stash:          gui.stashListContext(),
		Normal: &BasicContext{
			OnFocus: func(opts ...types.OnFocusOpts) error {
				return nil // TODO: should we do something here? We should allow for scrolling the panel
			},
			Kind:     types.MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_NORMAL_CONTEXT_KEY,
		},
		Staging: &BasicContext{
			OnFocus: func(opts ...types.OnFocusOpts) error {
				forceSecondaryFocused := false
				selectedLineIdx := -1
				if len(opts) > 0 && opts[0].ClickedViewName != "" {
					if opts[0].ClickedViewName == "main" || opts[0].ClickedViewName == "secondary" {
						selectedLineIdx = opts[0].ClickedViewLineIdx
					}
					if opts[0].ClickedViewName == "secondary" {
						forceSecondaryFocused = true
					}
				}
				return gui.onStagingFocus(forceSecondaryFocused, selectedLineIdx)
			},
			Kind:     types.MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_STAGING_CONTEXT_KEY,
		},
		PatchBuilding: &BasicContext{
			OnFocus: func(opts ...types.OnFocusOpts) error {
				selectedLineIdx := -1
				if len(opts) > 0 && (opts[0].ClickedViewName == "main" || opts[0].ClickedViewName == "secondary") {
					selectedLineIdx = opts[0].ClickedViewLineIdx
				}

				return gui.onPatchBuildingFocus(selectedLineIdx)
			},
			Kind:     types.MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_PATCH_BUILDING_CONTEXT_KEY,
		},
		Merging: &BasicContext{
			OnFocus:         OnFocusWrapper(func() error { return gui.renderConflictsWithLock(true) }),
			Kind:            types.MAIN_CONTEXT,
			ViewName:        "main",
			Key:             MAIN_MERGING_CONTEXT_KEY,
			OnGetOptionsMap: gui.getMergingOptions,
		},
		Credentials: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleAskFocused),
			Kind:     types.PERSISTENT_POPUP,
			ViewName: "credentials",
			Key:      CREDENTIALS_CONTEXT_KEY,
		},
		Confirmation: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleAskFocused),
			Kind:     types.TEMPORARY_POPUP,
			ViewName: "confirmation",
			Key:      CONFIRMATION_CONTEXT_KEY,
		},
		Suggestions: gui.suggestionsListContext(),
		CommitMessage: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleCommitMessageFocused),
			Kind:     types.PERSISTENT_POPUP,
			ViewName: "commitMessage",
			Key:      COMMIT_MESSAGE_CONTEXT_KEY,
		},
		Search: &BasicContext{
			Kind:     types.PERSISTENT_POPUP,
			ViewName: "search",
			Key:      SEARCH_CONTEXT_KEY,
		},
		CommandLog: &BasicContext{
			Kind:            types.EXTRAS_CONTEXT,
			ViewName:        "extras",
			Key:             COMMAND_LOG_CONTEXT_KEY,
			OnGetOptionsMap: gui.getMergingOptions,
			OnFocusLost: func() error {
				gui.Views.Extras.Autoscroll = true
				return nil
			},
		},
	}
}

// using this wrapper for when an onFocus function doesn't care about any potential
// props that could be passed
func OnFocusWrapper(f func() error) func(opts ...types.OnFocusOpts) error {
	return func(opts ...types.OnFocusOpts) error {
		return f()
	}
}
