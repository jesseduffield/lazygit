package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func NewContextTree(c *ContextCommon) *ContextTree {
	commitFilesContext := NewCommitFilesContext(c)

	return &ContextTree{
		Global: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:                  types.GLOBAL_CONTEXT,
				View:                  nil, // TODO: see if this breaks anything
				WindowName:            "",
				Key:                   GLOBAL_CONTEXT_KEY,
				Focusable:             false,
				HasUncontrolledBounds: true, // setting to true because the global context doesn't even have a view
			}),
		),
		Status: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.SIDE_CONTEXT,
				View:       c.Views().Status,
				WindowName: "status",
				Key:        STATUS_CONTEXT_KEY,
				Focusable:  true,
			}),
		),
		Files:          NewWorkingTreeContext(c),
		Submodules:     NewSubmodulesContext(c),
		Menu:           NewMenuContext(c),
		Remotes:        NewRemotesContext(c),
		RemoteBranches: NewRemoteBranchesContext(c),
		LocalCommits:   NewLocalCommitsContext(c),
		CommitFiles:    commitFilesContext,
		ReflogCommits:  NewReflogCommitsContext(c),
		SubCommits:     NewSubCommitsContext(c),
		Branches:       NewBranchesContext(c),
		Tags:           NewTagsContext(c),
		Stash:          NewStashContext(c),
		Suggestions:    NewSuggestionsContext(c),
		Normal: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       c.Views().Main,
				WindowName: "main",
				Key:        NORMAL_MAIN_CONTEXT_KEY,
				Focusable:  false,
			}),
		),
		NormalSecondary: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       c.Views().Secondary,
				WindowName: "secondary",
				Key:        NORMAL_SECONDARY_CONTEXT_KEY,
				Focusable:  false,
			}),
		),
		Staging: NewPatchExplorerContext(
			c.Views().Staging,
			"main",
			STAGING_MAIN_CONTEXT_KEY,
			func() []int { return nil },
			c,
		),
		StagingSecondary: NewPatchExplorerContext(
			c.Views().StagingSecondary,
			"secondary",
			STAGING_SECONDARY_CONTEXT_KEY,
			func() []int { return nil },
			c,
		),
		CustomPatchBuilder: NewPatchExplorerContext(
			c.Views().PatchBuilding,
			"main",
			PATCH_BUILDING_MAIN_CONTEXT_KEY,
			func() []int {
				filename := commitFilesContext.GetSelectedPath()
				includedLineIndices, err := c.Git().Patch.PatchBuilder.GetFileIncLineIndices(filename)
				if err != nil {
					c.Log.Error(err)
					return nil
				}

				return includedLineIndices
			},
			c,
		),
		CustomPatchBuilderSecondary: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       c.Views().PatchBuildingSecondary,
				WindowName: "secondary",
				Key:        PATCH_BUILDING_SECONDARY_CONTEXT_KEY,
				Focusable:  false,
			}),
		),
		MergeConflicts: NewMergeConflictsContext(
			c,
		),
		Confirmation:  NewConfirmationContext(c),
		CommitMessage: NewCommitMessageContext(c),
		CommitDescription: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:                  types.PERSISTENT_POPUP,
				View:                  c.Views().CommitDescription,
				WindowName:            "commitDescription",
				Key:                   COMMIT_DESCRIPTION_CONTEXT_KEY,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}),
		),
		Search: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.PERSISTENT_POPUP,
				View:       c.Views().Search,
				WindowName: "search",
				Key:        SEARCH_CONTEXT_KEY,
				Focusable:  true,
			}),
		),
		CommandLog: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.EXTRAS_CONTEXT,
				View:       c.Views().Extras,
				WindowName: "extras",
				Key:        COMMAND_LOG_CONTEXT_KEY,
				Focusable:  true,
			}),
		),
		Snake: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.SIDE_CONTEXT,
				View:       c.Views().Snake,
				WindowName: "files",
				Key:        SNAKE_CONTEXT_KEY,
				Focusable:  true,
			}),
		),
		Options:      NewDisplayContext(OPTIONS_CONTEXT_KEY, c.Views().Options, "options"),
		AppStatus:    NewDisplayContext(APP_STATUS_CONTEXT_KEY, c.Views().AppStatus, "appStatus"),
		SearchPrefix: NewDisplayContext(SEARCH_PREFIX_CONTEXT_KEY, c.Views().SearchPrefix, "searchPrefix"),
		Information:  NewDisplayContext(INFORMATION_CONTEXT_KEY, c.Views().Information, "information"),
		Limit:        NewDisplayContext(LIMIT_CONTEXT_KEY, c.Views().Limit, "limit"),
	}
}
