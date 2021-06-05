package types

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

var AllContextKeys = []ContextKey{
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
