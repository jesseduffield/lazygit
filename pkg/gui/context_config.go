package gui

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) contextTree() *context.ContextTree {
	return &context.ContextTree{
		Global: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:                  types.GLOBAL_CONTEXT,
				View:                  nil, // TODO: see if this breaks anything
				WindowName:            "",
				Key:                   context.GLOBAL_CONTEXT_KEY,
				Focusable:             false,
				HasUncontrolledBounds: true, // setting to true because the global context doesn't even have a view
			}),
			context.ContextCallbackOpts{},
		),
		Status: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.SIDE_CONTEXT,
				View:       gui.Views.Status,
				WindowName: "status",
				Key:        context.STATUS_CONTEXT_KEY,
				Focusable:  true,
			}),
			context.ContextCallbackOpts{},
		),
		Snake: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.SIDE_CONTEXT,
				View:       gui.Views.Snake,
				WindowName: "files",
				Key:        context.SNAKE_CONTEXT_KEY,
				Focusable:  true,
			}),
			context.ContextCallbackOpts{},
		),
		Files:          gui.filesListContext(),
		Submodules:     gui.submodulesListContext(),
		Menu:           gui.menuListContext(),
		Remotes:        gui.remotesListContext(),
		RemoteBranches: gui.remoteBranchesListContext(),
		LocalCommits:   gui.branchCommitsListContext(),
		CommitFiles:    gui.commitFilesListContext(),
		ReflogCommits:  gui.reflogCommitsListContext(),
		SubCommits:     gui.subCommitsListContext(),
		Branches:       gui.branchesListContext(),
		Tags:           gui.tagsListContext(),
		Stash:          gui.stashListContext(),
		Suggestions:    gui.suggestionsListContext(),
		Normal: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       gui.Views.Main,
				WindowName: "main",
				Key:        context.NORMAL_MAIN_CONTEXT_KEY,
				Focusable:  false,
			}),
			context.ContextCallbackOpts{},
		),
		NormalSecondary: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       gui.Views.Secondary,
				WindowName: "secondary",
				Key:        context.NORMAL_SECONDARY_CONTEXT_KEY,
				Focusable:  false,
			}),
			context.ContextCallbackOpts{},
		),
		Staging: context.NewPatchExplorerContext(
			gui.Views.Staging,
			"main",
			context.STAGING_MAIN_CONTEXT_KEY,
			func(opts types.OnFocusOpts) error {
				gui.Views.Staging.Wrap = false
				gui.Views.StagingSecondary.Wrap = false

				return gui.helpers.Staging.RefreshStagingPanel(opts)
			},
			func(opts types.OnFocusLostOpts) error {
				gui.State.Contexts.Staging.SetState(nil)

				if opts.NewContextKey != context.STAGING_SECONDARY_CONTEXT_KEY {
					gui.Views.Staging.Wrap = true
					gui.Views.StagingSecondary.Wrap = true
					_ = gui.State.Contexts.Staging.Render(false)
					_ = gui.State.Contexts.StagingSecondary.Render(false)
				}
				return nil
			},
			func() []int { return nil },
			gui.c,
		),
		StagingSecondary: context.NewPatchExplorerContext(
			gui.Views.StagingSecondary,
			"secondary",
			context.STAGING_SECONDARY_CONTEXT_KEY,
			func(opts types.OnFocusOpts) error {
				gui.Views.Staging.Wrap = false
				gui.Views.StagingSecondary.Wrap = false

				return gui.helpers.Staging.RefreshStagingPanel(opts)
			},
			func(opts types.OnFocusLostOpts) error {
				gui.State.Contexts.StagingSecondary.SetState(nil)

				if opts.NewContextKey != context.STAGING_MAIN_CONTEXT_KEY {
					gui.Views.Staging.Wrap = true
					gui.Views.StagingSecondary.Wrap = true
					_ = gui.State.Contexts.Staging.Render(false)
					_ = gui.State.Contexts.StagingSecondary.Render(false)
				}
				return nil
			},
			func() []int { return nil },
			gui.c,
		),
		CustomPatchBuilder: context.NewPatchExplorerContext(
			gui.Views.PatchBuilding,
			"main",
			context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
			func(opts types.OnFocusOpts) error {
				// no need to change wrap on the secondary view because it can't be interacted with
				gui.Views.PatchBuilding.Wrap = false

				return gui.helpers.PatchBuilding.RefreshPatchBuildingPanel(opts)
			},
			func(opts types.OnFocusLostOpts) error {
				gui.Views.PatchBuilding.Wrap = true

				if gui.git.Patch.PatchBuilder.IsEmpty() {
					gui.git.Patch.PatchBuilder.Reset()
				}

				return nil
			},
			func() []int {
				filename := gui.State.Contexts.CommitFiles.GetSelectedPath()
				includedLineIndices, err := gui.git.Patch.PatchBuilder.GetFileIncLineIndices(filename)
				if err != nil {
					gui.Log.Error(err)
					return nil
				}

				return includedLineIndices
			},
			gui.c,
		),
		CustomPatchBuilderSecondary: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       gui.Views.PatchBuildingSecondary,
				WindowName: "secondary",
				Key:        context.PATCH_BUILDING_SECONDARY_CONTEXT_KEY,
				Focusable:  false,
			}),
			context.ContextCallbackOpts{},
		),
		MergeConflicts: context.NewMergeConflictsContext(
			gui.Views.MergeConflicts,
			context.ContextCallbackOpts{},
			gui.c,
			func() map[string]string {
				// wrapping in a function because contexts are initialized before helpers
				return gui.helpers.MergeConflicts.GetMergingOptions()
			},
		),
		Confirmation: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:                  types.TEMPORARY_POPUP,
				View:                  gui.Views.Confirmation,
				WindowName:            "confirmation",
				Key:                   context.CONFIRMATION_CONTEXT_KEY,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}),
			context.ContextCallbackOpts{
				OnFocus: OnFocusWrapper(gui.handleAskFocused),
			},
		),
		CommitMessage: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:                  types.PERSISTENT_POPUP,
				View:                  gui.Views.CommitMessage,
				WindowName:            "commitMessage",
				Key:                   context.COMMIT_MESSAGE_CONTEXT_KEY,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}),
			context.ContextCallbackOpts{
				OnFocus: OnFocusWrapper(gui.handleCommitMessageFocused),
			},
		),
		Search: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.PERSISTENT_POPUP,
				View:       gui.Views.Search,
				WindowName: "search",
				Key:        context.SEARCH_CONTEXT_KEY,
				Focusable:  true,
			}),
			context.ContextCallbackOpts{},
		),
		CommandLog: context.NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.EXTRAS_CONTEXT,
				View:       gui.Views.Extras,
				WindowName: "extras",
				Key:        context.COMMAND_LOG_CONTEXT_KEY,
				Focusable:  true,
			}),
			context.ContextCallbackOpts{
				OnFocusLost: func(opts types.OnFocusLostOpts) error {
					gui.Views.Extras.Autoscroll = true
					return nil
				},
			},
		),
		Options:      context.NewDisplayContext(context.OPTIONS_CONTEXT_KEY, gui.Views.Options, "options"),
		AppStatus:    context.NewDisplayContext(context.APP_STATUS_CONTEXT_KEY, gui.Views.AppStatus, "appStatus"),
		SearchPrefix: context.NewDisplayContext(context.SEARCH_PREFIX_CONTEXT_KEY, gui.Views.SearchPrefix, "searchPrefix"),
		Information:  context.NewDisplayContext(context.INFORMATION_CONTEXT_KEY, gui.Views.Information, "information"),
		Limit:        context.NewDisplayContext(context.LIMIT_CONTEXT_KEY, gui.Views.Limit, "limit"),
	}
}

// using this wrapper for when an onFocus function doesn't care about any potential
// props that could be passed
func OnFocusWrapper(f func() error) func(opts types.OnFocusOpts) error {
	return func(opts types.OnFocusOpts) error {
		return f()
	}
}

func (gui *Gui) getPatchExplorerContexts() []types.IPatchExplorerContext {
	return []types.IPatchExplorerContext{
		gui.State.Contexts.Staging,
		gui.State.Contexts.StagingSecondary,
		gui.State.Contexts.CustomPatchBuilder,
	}
}

func (gui *Gui) popupViewNames() []string {
	popups := slices.Filter(gui.State.Contexts.Flatten(), func(c types.Context) bool {
		return c.GetKind() == types.PERSISTENT_POPUP || c.GetKind() == types.TEMPORARY_POPUP
	})

	return slices.Map(popups, func(c types.Context) string {
		return c.GetViewName()
	})
}

func (gui *Gui) defaultSideContext() types.Context {
	if gui.State.Modes.Filtering.Active() {
		return gui.State.Contexts.LocalCommits
	} else {
		return gui.State.Contexts.Files
	}
}

func (gui *Gui) TransientContexts() []types.Context {
	return slices.Filter(gui.State.Contexts.Flatten(), func(context types.Context) bool {
		return context.IsTransient()
	})
}
