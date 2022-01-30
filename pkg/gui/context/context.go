package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

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

type ContextTree struct {
	Status         types.Context
	Files          *WorkingTreeContext
	Submodules     types.IListContext
	Menu           types.IListContext
	Branches       types.IListContext
	Remotes        types.IListContext
	RemoteBranches types.IListContext
	Tags           *TagsContext
	BranchCommits  types.IListContext
	CommitFiles    types.IListContext
	ReflogCommits  types.IListContext
	SubCommits     types.IListContext
	Stash          types.IListContext
	Suggestions    types.IListContext
	Normal         types.Context
	Staging        types.Context
	PatchBuilding  types.Context
	Merging        types.Context
	Credentials    types.Context
	Confirmation   types.Context
	CommitMessage  types.Context
	Search         types.Context
	CommandLog     types.Context
}

func (tree ContextTree) InitialViewContextMap() map[string]types.Context {
	return map[string]types.Context{
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

type TabContext struct {
	Tab      string
	Contexts []types.Context
}

func (tree ContextTree) InitialViewTabContextMap() map[string][]TabContext {
	return map[string][]TabContext{
		"branches": {
			{
				Tab:      "Local Branches",
				Contexts: []types.Context{tree.Branches},
			},
			{
				Tab: "Remotes",
				Contexts: []types.Context{
					tree.Remotes,
					tree.RemoteBranches,
				},
			},
			{
				Tab:      "Tags",
				Contexts: []types.Context{tree.Tags},
			},
		},
		"commits": {
			{
				Tab:      "Commits",
				Contexts: []types.Context{tree.BranchCommits},
			},
			{
				Tab: "Reflog",
				Contexts: []types.Context{
					tree.ReflogCommits,
				},
			},
		},
		"files": {
			{
				Tab:      "Files",
				Contexts: []types.Context{tree.Files},
			},
			{
				Tab: "Submodules",
				Contexts: []types.Context{
					tree.Submodules,
				},
			},
		},
	}
}
