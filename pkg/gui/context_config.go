package gui

type ContextKey string

const (
	STATUS_CONTEXT_KEY              ContextKey = "status"
	FILES_CONTEXT_KEY               ContextKey = "files"
	LOCAL_BRANCHES_CONTEXT_KEY      ContextKey = "localBranches"
	REMOTES_CONTEXT_KEY             ContextKey = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY     ContextKey = "remoteBranches"
	TAGS_CONTEXT_KEY                ContextKey = "tags"
	BRANCH_COMMITS_CONTEXT_KEY      ContextKey = "commits"
	REFLOG_COMMITS_CONTEXT_KEY      ContextKey = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY         ContextKey = "subCommits"
	COMMIT_FILES_CONTEXT_KEY        ContextKey = "commitFiles"
	STASH_CONTEXT_KEY               ContextKey = "stash"
	MAIN_NORMAL_CONTEXT_KEY         ContextKey = "normal"
	MAIN_MERGING_CONTEXT_KEY        ContextKey = "merging"
	MAIN_PATCH_BUILDING_CONTEXT_KEY ContextKey = "patchBuilding"
	MAIN_STAGING_CONTEXT_KEY        ContextKey = "staging"
	MENU_CONTEXT_KEY                ContextKey = "menu"
	CREDENTIALS_CONTEXT_KEY         ContextKey = "credentials"
	CONFIRMATION_CONTEXT_KEY        ContextKey = "confirmation"
	SEARCH_CONTEXT_KEY              ContextKey = "search"
	COMMIT_MESSAGE_CONTEXT_KEY      ContextKey = "commitMessage"
	SUBMODULES_CONTEXT_KEY          ContextKey = "submodules"
	SUGGESTIONS_CONTEXT_KEY         ContextKey = "suggestions"
	COMMAND_LOG_CONTEXT_KEY         ContextKey = "cmdLog"
)

var allContextKeys = []ContextKey{
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

type ContextTree struct {
	Status         Context
	Files          IListContext
	Submodules     IListContext
	Menu           IListContext
	Branches       IListContext
	Remotes        IListContext
	RemoteBranches IListContext
	Tags           IListContext
	BranchCommits  IListContext
	CommitFiles    IListContext
	ReflogCommits  IListContext
	SubCommits     IListContext
	Stash          IListContext
	Suggestions    IListContext
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
			OnRenderToMain: OnFocusWrapper(gui.statusRenderToMain),
			Kind:           SIDE_CONTEXT,
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
			OnFocus: func(opts ...OnFocusOpts) error {
				return nil // TODO: should we do something here? We should allow for scrolling the panel
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_NORMAL_CONTEXT_KEY,
		},
		Staging: &BasicContext{
			OnFocus: func(opts ...OnFocusOpts) error {
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
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_STAGING_CONTEXT_KEY,
		},
		PatchBuilding: &BasicContext{
			OnFocus: func(opts ...OnFocusOpts) error {
				selectedLineIdx := -1
				if len(opts) > 0 && (opts[0].ClickedViewName == "main" || opts[0].ClickedViewName == "secondary") {
					selectedLineIdx = opts[0].ClickedViewLineIdx
				}

				return gui.onPatchBuildingFocus(selectedLineIdx)
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_PATCH_BUILDING_CONTEXT_KEY,
		},
		Merging: &BasicContext{
			OnFocus:         OnFocusWrapper(func() error { return gui.renderConflictsWithLock(true) }),
			Kind:            MAIN_CONTEXT,
			ViewName:        "main",
			Key:             MAIN_MERGING_CONTEXT_KEY,
			OnGetOptionsMap: gui.getMergingOptions,
		},
		Credentials: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleAskFocused),
			Kind:     PERSISTENT_POPUP,
			ViewName: "credentials",
			Key:      CREDENTIALS_CONTEXT_KEY,
		},
		Confirmation: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleAskFocused),
			Kind:     TEMPORARY_POPUP,
			ViewName: "confirmation",
			Key:      CONFIRMATION_CONTEXT_KEY,
		},
		Suggestions: gui.suggestionsListContext(),
		CommitMessage: &BasicContext{
			OnFocus:  OnFocusWrapper(gui.handleCommitMessageFocused),
			Kind:     PERSISTENT_POPUP,
			ViewName: "commitMessage",
			Key:      COMMIT_MESSAGE_CONTEXT_KEY,
		},
		Search: &BasicContext{
			Kind:     PERSISTENT_POPUP,
			ViewName: "search",
			Key:      SEARCH_CONTEXT_KEY,
		},
		CommandLog: &BasicContext{
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

// using this wrapper for when an onFocus function doesn't care about any potential
// props that could be passed
func OnFocusWrapper(f func() error) func(opts ...OnFocusOpts) error {
	return func(opts ...OnFocusOpts) error {
		return f()
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
