package context

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	GLOBAL_CONTEXT_KEY              types.ContextKey = "global"
	STATUS_CONTEXT_KEY              types.ContextKey = "status"
	FILES_CONTEXT_KEY               types.ContextKey = "files"
	LOCAL_BRANCHES_CONTEXT_KEY      types.ContextKey = "localBranches"
	REMOTES_CONTEXT_KEY             types.ContextKey = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY     types.ContextKey = "remoteBranches"
	TAGS_CONTEXT_KEY                types.ContextKey = "tags"
	LOCAL_COMMITS_CONTEXT_KEY       types.ContextKey = "commits"
	REFLOG_COMMITS_CONTEXT_KEY      types.ContextKey = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY         types.ContextKey = "subCommits"
	COMMIT_FILES_CONTEXT_KEY        types.ContextKey = "commitFiles"
	STASH_CONTEXT_KEY               types.ContextKey = "stash"
	MAIN_NORMAL_CONTEXT_KEY         types.ContextKey = "normal"
	MAIN_MERGING_CONTEXT_KEY        types.ContextKey = "merging"
	MAIN_PATCH_BUILDING_CONTEXT_KEY types.ContextKey = "patchBuilding"
	MAIN_STAGING_CONTEXT_KEY        types.ContextKey = "staging"
	MENU_CONTEXT_KEY                types.ContextKey = "menu"
	CONFIRMATION_CONTEXT_KEY        types.ContextKey = "confirmation"
	SEARCH_CONTEXT_KEY              types.ContextKey = "search"
	COMMIT_MESSAGE_CONTEXT_KEY      types.ContextKey = "commitMessage"
	SUBMODULES_CONTEXT_KEY          types.ContextKey = "submodules"
	SUGGESTIONS_CONTEXT_KEY         types.ContextKey = "suggestions"
	COMMAND_LOG_CONTEXT_KEY         types.ContextKey = "cmdLog"
)

var AllContextKeys = []types.ContextKey{
	GLOBAL_CONTEXT_KEY, // not focusable
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
	MAIN_NORMAL_CONTEXT_KEY, // not focusable
	MAIN_MERGING_CONTEXT_KEY,
	MAIN_PATCH_BUILDING_CONTEXT_KEY,
	MAIN_STAGING_CONTEXT_KEY, // not focusable for secondary view
	MENU_CONTEXT_KEY,
	CONFIRMATION_CONTEXT_KEY,
	SEARCH_CONTEXT_KEY,
	COMMIT_MESSAGE_CONTEXT_KEY,
	SUBMODULES_CONTEXT_KEY,
	SUGGESTIONS_CONTEXT_KEY,
	COMMAND_LOG_CONTEXT_KEY,
}

type ContextTree struct {
	Global         types.Context
	Status         types.Context
	Files          *WorkingTreeContext
	Menu           *MenuContext
	Branches       *BranchesContext
	Tags           *TagsContext
	LocalCommits   *LocalCommitsContext
	CommitFiles    *CommitFilesContext
	Remotes        *RemotesContext
	Submodules     *SubmodulesContext
	RemoteBranches *RemoteBranchesContext
	ReflogCommits  *ReflogCommitsContext
	SubCommits     *SubCommitsContext
	Stash          *StashContext
	Suggestions    *SuggestionsContext
	Normal         types.Context
	Staging        types.Context
	PatchBuilding  types.Context
	Merging        types.Context
	Confirmation   types.Context
	CommitMessage  types.Context
	Search         types.Context
	CommandLog     types.Context
}

func (self *ContextTree) Flatten() []types.Context {
	return []types.Context{
		self.Global,
		self.Status,
		self.Files,
		self.Submodules,
		self.Branches,
		self.Remotes,
		self.RemoteBranches,
		self.Tags,
		self.LocalCommits,
		self.CommitFiles,
		self.ReflogCommits,
		self.Stash,
		self.Menu,
		self.Confirmation,
		self.CommitMessage,
		self.Normal,
		self.Staging,
		self.Merging,
		self.PatchBuilding,
		self.SubCommits,
		self.Suggestions,
		self.CommandLog,
	}
}

type ViewContextMap struct {
	content map[string]types.Context
	sync.RWMutex
}

func NewViewContextMap() *ViewContextMap {
	return &ViewContextMap{content: map[string]types.Context{}}
}

func (self *ViewContextMap) Get(viewName string) types.Context {
	self.RLock()
	defer self.RUnlock()

	return self.content[viewName]
}

func (self *ViewContextMap) Set(viewName string, context types.Context) {
	self.Lock()
	defer self.Unlock()
	self.content[viewName] = context
}

func (self *ViewContextMap) Entries() map[string]types.Context {
	self.Lock()
	defer self.Unlock()
	return self.content
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
				Tab:      "Remotes",
				Contexts: []types.Context{tree.Remotes},
			},
			{
				Tab:      "Tags",
				Contexts: []types.Context{tree.Tags},
			},
		},
		"commits": {
			{
				Tab:      "Commits",
				Contexts: []types.Context{tree.LocalCommits},
			},
			{
				Tab:      "Reflog",
				Contexts: []types.Context{tree.ReflogCommits},
			},
		},
		"files": {
			{
				Tab:      "Files",
				Contexts: []types.Context{tree.Files},
			},
			{
				Tab:      "Submodules",
				Contexts: []types.Context{tree.Submodules},
			},
		},
	}
}
