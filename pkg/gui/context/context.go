package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	GLOBAL_CONTEXT_KEY                   types.ContextKey = "global"
	STATUS_CONTEXT_KEY                   types.ContextKey = "status"
	SNAKE_CONTEXT_KEY                    types.ContextKey = "snake"
	FILES_CONTEXT_KEY                    types.ContextKey = "files"
	LOCAL_BRANCHES_CONTEXT_KEY           types.ContextKey = "localBranches"
	REMOTES_CONTEXT_KEY                  types.ContextKey = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY          types.ContextKey = "remoteBranches"
	TAGS_CONTEXT_KEY                     types.ContextKey = "tags"
	LOCAL_COMMITS_CONTEXT_KEY            types.ContextKey = "commits"
	REFLOG_COMMITS_CONTEXT_KEY           types.ContextKey = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY              types.ContextKey = "subCommits"
	COMMIT_FILES_CONTEXT_KEY             types.ContextKey = "commitFiles"
	STASH_CONTEXT_KEY                    types.ContextKey = "stash"
	NORMAL_MAIN_CONTEXT_KEY              types.ContextKey = "normal"
	NORMAL_SECONDARY_CONTEXT_KEY         types.ContextKey = "normalSecondary"
	STAGING_MAIN_CONTEXT_KEY             types.ContextKey = "staging"
	STAGING_SECONDARY_CONTEXT_KEY        types.ContextKey = "stagingSecondary"
	PATCH_BUILDING_MAIN_CONTEXT_KEY      types.ContextKey = "patchBuilding"
	PATCH_BUILDING_SECONDARY_CONTEXT_KEY types.ContextKey = "patchBuildingSecondary"
	MERGE_CONFLICTS_CONTEXT_KEY          types.ContextKey = "mergeConflicts"

	// these shouldn't really be needed for anything but I'm giving them unique keys nonetheless
	OPTIONS_CONTEXT_KEY       types.ContextKey = "options"
	APP_STATUS_CONTEXT_KEY    types.ContextKey = "appStatus"
	SEARCH_PREFIX_CONTEXT_KEY types.ContextKey = "searchPrefix"
	INFORMATION_CONTEXT_KEY   types.ContextKey = "information"
	LIMIT_CONTEXT_KEY         types.ContextKey = "limit"

	MENU_CONTEXT_KEY               types.ContextKey = "menu"
	CONFIRMATION_CONTEXT_KEY       types.ContextKey = "confirmation"
	SEARCH_CONTEXT_KEY             types.ContextKey = "search"
	COMMIT_MESSAGE_CONTEXT_KEY     types.ContextKey = "commitMessage"
	COMMIT_DESCRIPTION_CONTEXT_KEY types.ContextKey = "commitDescription"
	SUBMODULES_CONTEXT_KEY         types.ContextKey = "submodules"
	SUGGESTIONS_CONTEXT_KEY        types.ContextKey = "suggestions"
	COMMAND_LOG_CONTEXT_KEY        types.ContextKey = "cmdLog"
)

var AllContextKeys = []types.ContextKey{
	GLOBAL_CONTEXT_KEY,
	STATUS_CONTEXT_KEY,
	FILES_CONTEXT_KEY,
	LOCAL_BRANCHES_CONTEXT_KEY,
	REMOTES_CONTEXT_KEY,
	REMOTE_BRANCHES_CONTEXT_KEY,
	TAGS_CONTEXT_KEY,
	LOCAL_COMMITS_CONTEXT_KEY,
	REFLOG_COMMITS_CONTEXT_KEY,
	SUB_COMMITS_CONTEXT_KEY,
	COMMIT_FILES_CONTEXT_KEY,
	STASH_CONTEXT_KEY,
	NORMAL_MAIN_CONTEXT_KEY,
	NORMAL_SECONDARY_CONTEXT_KEY,
	STAGING_MAIN_CONTEXT_KEY,
	STAGING_SECONDARY_CONTEXT_KEY,
	PATCH_BUILDING_MAIN_CONTEXT_KEY,
	PATCH_BUILDING_SECONDARY_CONTEXT_KEY,
	MERGE_CONFLICTS_CONTEXT_KEY,

	MENU_CONTEXT_KEY,
	CONFIRMATION_CONTEXT_KEY,
	SEARCH_CONTEXT_KEY,
	COMMIT_MESSAGE_CONTEXT_KEY,
	SUBMODULES_CONTEXT_KEY,
	SUGGESTIONS_CONTEXT_KEY,
	COMMAND_LOG_CONTEXT_KEY,
}

type ContextTree struct {
	Global                      types.Context
	Status                      types.Context
	Snake                       types.Context
	Files                       *WorkingTreeContext
	Menu                        *MenuContext
	Branches                    *BranchesContext
	Tags                        *TagsContext
	LocalCommits                *LocalCommitsContext
	CommitFiles                 *CommitFilesContext
	Remotes                     *RemotesContext
	Submodules                  *SubmodulesContext
	RemoteBranches              *RemoteBranchesContext
	ReflogCommits               *ReflogCommitsContext
	SubCommits                  *SubCommitsContext
	Stash                       *StashContext
	Suggestions                 *SuggestionsContext
	Normal                      types.Context
	NormalSecondary             types.Context
	Staging                     *PatchExplorerContext
	StagingSecondary            *PatchExplorerContext
	CustomPatchBuilder          *PatchExplorerContext
	CustomPatchBuilderSecondary types.Context
	MergeConflicts              *MergeConflictsContext
	Confirmation                types.Context
	CommitMessage               *CommitMessageContext
	CommitDescription           types.Context
	CommandLog                  types.Context

	// display contexts
	AppStatus    types.Context
	Options      types.Context
	SearchPrefix types.Context
	Search       types.Context
	Information  types.Context
	Limit        types.Context
}

// the order of this decides which context is initially at the top of its window
func (self *ContextTree) Flatten() []types.Context {
	return []types.Context{
		self.Global,
		self.Status,
		self.Snake,
		self.Submodules,
		self.Files,
		self.SubCommits,
		self.Remotes,
		self.RemoteBranches,
		self.Tags,
		self.Branches,
		self.CommitFiles,
		self.ReflogCommits,
		self.LocalCommits,
		self.Stash,
		self.Menu,
		self.Confirmation,
		self.CommitMessage,
		self.CommitDescription,

		self.MergeConflicts,
		self.StagingSecondary,
		self.Staging,
		self.CustomPatchBuilderSecondary,
		self.CustomPatchBuilder,
		self.NormalSecondary,
		self.Normal,

		self.Suggestions,
		self.CommandLog,
		self.AppStatus,
		self.Options,
		self.SearchPrefix,
		self.Search,
		self.Information,
		self.Limit,
	}
}

type TabView struct {
	Tab      string
	ViewName string
}
