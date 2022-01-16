package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type ContextTree struct {
	Status         types.Context
	Files          types.IListContext
	Submodules     types.IListContext
	Menu           types.IListContext
	Branches       types.IListContext
	Remotes        types.IListContext
	RemoteBranches types.IListContext
	Tags           types.IListContext
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
