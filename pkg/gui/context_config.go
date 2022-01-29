package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

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
		Status: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.SIDE_CONTEXT,
				ViewName:   "status",
				Key:        context.STATUS_CONTEXT_KEY,
				WindowName: "status",
			}),
			NewSimpleContextOpts{
				OnRenderToMain: OnFocusWrapper(gui.statusRenderToMain),
			},
		),
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
		Suggestions:    gui.suggestionsListContext(),
		Normal: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				ViewName:   "main",
				WindowName: "main",
				Key:        context.MAIN_NORMAL_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
				OnFocus: func(opts ...types.OnFocusOpts) error {
					return nil // TODO: should we do something here? We should allow for scrolling the panel
				},
			},
		),
		Staging: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				ViewName:   "main",
				WindowName: "main",
				Key:        context.MAIN_STAGING_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
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
			},
		),
		PatchBuilding: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				ViewName:   "main",
				WindowName: "main",
				Key:        context.MAIN_PATCH_BUILDING_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
				OnFocus: func(opts ...types.OnFocusOpts) error {
					selectedLineIdx := -1
					if len(opts) > 0 && (opts[0].ClickedViewName == "main" || opts[0].ClickedViewName == "secondary") {
						selectedLineIdx = opts[0].ClickedViewLineIdx
					}

					return gui.onPatchBuildingFocus(selectedLineIdx)
				},
			},
		),
		Merging: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:            types.MAIN_CONTEXT,
				ViewName:        "main",
				WindowName:      "main",
				Key:             context.MAIN_MERGING_CONTEXT_KEY,
				OnGetOptionsMap: gui.getMergingOptions,
			}),
			NewSimpleContextOpts{
				OnFocus: OnFocusWrapper(func() error { return gui.renderConflictsWithLock(true) }),
			},
		),
		Credentials: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.PERSISTENT_POPUP,
				ViewName:   "credentials",
				WindowName: "credentials",
				Key:        context.CREDENTIALS_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
				OnFocus: OnFocusWrapper(gui.handleAskFocused),
			},
		),
		Confirmation: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.TEMPORARY_POPUP,
				ViewName:   "confirmation",
				WindowName: "confirmation",
				Key:        context.CONFIRMATION_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
				OnFocus: OnFocusWrapper(gui.handleAskFocused),
			},
		),
		CommitMessage: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.PERSISTENT_POPUP,
				ViewName:   "commitMessage",
				WindowName: "commitMessage",
				Key:        context.COMMIT_MESSAGE_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{
				OnFocus: OnFocusWrapper(gui.handleCommitMessageFocused),
			},
		),
		Search: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:       types.PERSISTENT_POPUP,
				ViewName:   "search",
				WindowName: "search",
				Key:        context.SEARCH_CONTEXT_KEY,
			}),
			NewSimpleContextOpts{},
		),
		CommandLog: NewSimpleContext(
			context.NewBaseContext(context.NewBaseContextOpts{
				Kind:            types.EXTRAS_CONTEXT,
				ViewName:        "extras",
				WindowName:      "extras",
				Key:             context.COMMAND_LOG_CONTEXT_KEY,
				OnGetOptionsMap: gui.getMergingOptions,
			}),
			NewSimpleContextOpts{
				OnFocusLost: func() error {
					gui.Views.Extras.Autoscroll = true
					return nil
				},
			},
		),
	}
}

// using this wrapper for when an onFocus function doesn't care about any potential
// props that could be passed
func OnFocusWrapper(f func() error) func(opts ...types.OnFocusOpts) error {
	return func(opts ...types.OnFocusOpts) error {
		return f()
	}
}
