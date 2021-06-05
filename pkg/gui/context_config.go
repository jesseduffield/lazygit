package gui

import . "github.com/jesseduffield/lazygit/pkg/gui/types"

type ContextTree struct {
	Status         Context
	Files          *ListContext
	Submodules     *ListContext
	Menu           *ListContext
	Branches       *ListContext
	Remotes        *ListContext
	RemoteBranches *ListContext
	Tags           *ListContext
	BranchCommits  *ListContext
	CommitFiles    *ListContext
	ReflogCommits  *ListContext
	SubCommits     *ListContext
	Stash          *ListContext
	Suggestions    *ListContext
	Normal         Context
	Staging        Context
	PatchBuilding  Context
	Merging        Context
	Credentials    Context
	Confirmation   Context
	CommitMessage  Context
	Search         Context
	CommandLog     Context
}

func (gui *Gui) allContexts() []Context {
	return []Context{
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

func (gui *Gui) contextTree() ContextTree {
	return ContextTree{
		Status: &BasicContext{
			OnFocus:  gui.handleStatusSelect,
			Kind:     SIDE_CONTEXT,
			ViewName: "status",
			Key:      STATUS_CONTEXT_KEY,
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
			OnFocus: func() error {
				return nil // TODO: should we do something here? We should allow for scrolling the panel
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_NORMAL_CONTEXT_KEY,
		},
		Staging: &BasicContext{
			OnFocus: func() error {
				return nil
				// TODO: centralise the code here
				// return gui.refreshStagingPanel(false, -1)
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_STAGING_CONTEXT_KEY,
		},
		PatchBuilding: &BasicContext{
			OnFocus: func() error {
				return nil
				// TODO: centralise the code here
				// return gui.refreshPatchBuildingPanel(-1)
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_PATCH_BUILDING_CONTEXT_KEY,
		},
		Merging: &BasicContext{
			OnFocus:         gui.refreshMergePanelWithLock,
			Kind:            MAIN_CONTEXT,
			ViewName:        "main",
			Key:             MAIN_MERGING_CONTEXT_KEY,
			OnGetOptionsMap: gui.getMergingOptions,
		},
		Credentials: &BasicContext{
			OnFocus:  gui.handleCredentialsViewFocused,
			Kind:     PERSISTENT_POPUP,
			ViewName: "credentials",
			Key:      CREDENTIALS_CONTEXT_KEY,
		},
		Confirmation: &BasicContext{
			OnFocus:  func() error { return nil },
			Kind:     TEMPORARY_POPUP,
			ViewName: "confirmation",
			Key:      CONFIRMATION_CONTEXT_KEY,
		},
		Suggestions: gui.suggestionsListContext(),
		CommitMessage: &BasicContext{
			OnFocus:  gui.handleCommitMessageFocused,
			Kind:     PERSISTENT_POPUP,
			ViewName: "commitMessage",
			Key:      COMMIT_MESSAGE_CONTEXT_KEY,
		},
		Search: &BasicContext{
			OnFocus:  func() error { return nil },
			Kind:     PERSISTENT_POPUP,
			ViewName: "search",
			Key:      SEARCH_CONTEXT_KEY,
		},
		CommandLog: &BasicContext{
			OnFocus:         func() error { return nil },
			Kind:            EXTRAS_CONTEXT,
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

func (tree ContextTree) initialViewContextMap() map[string]Context {
	return map[string]Context{
		"status":        tree.Status,
		"files":         tree.Files,
		"branches":      tree.Branches,
		"commits":       tree.BranchCommits,
		"commitFiles":   tree.CommitFiles,
		"stash":         tree.Stash,
		"menu":          tree.Menu,
		"confirmation":  tree.Confirmation,
		"credentials":   tree.Credentials,
		"commitMessage": tree.CommitMessage,
		"main":          tree.Normal,
		"secondary":     tree.Normal,
		"extras":        tree.CommandLog,
	}
}

func (tree ContextTree) initialViewTabContextMap() map[string][]tabContext {
	return map[string][]tabContext{
		"branches": {
			{
				tab:      "Local Branches",
				contexts: []Context{tree.Branches},
			},
			{
				tab: "Remotes",
				contexts: []Context{
					tree.Remotes,
					tree.RemoteBranches,
				},
			},
			{
				tab:      "Tags",
				contexts: []Context{tree.Tags},
			},
		},
		"commits": {
			{
				tab:      "Commits",
				contexts: []Context{tree.BranchCommits},
			},
			{
				tab: "Reflog",
				contexts: []Context{
					tree.ReflogCommits,
				},
			},
		},
		"files": {
			{
				tab:      "Files",
				contexts: []Context{tree.Files},
			},
			{
				tab: "Submodules",
				contexts: []Context{
					tree.Submodules,
				},
			},
		},
	}
}
