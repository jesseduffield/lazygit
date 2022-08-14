package gui

// const diffForTest = `diff --git a/pkg/gui/diff_context_size.go b/pkg/gui/diff_context_size.go
// index 0da0a982..742b7dcf 100644
// --- a/pkg/gui/diff_context_size.go
// +++ b/pkg/gui/diff_context_size.go
// @@ -9,12 +9,12 @@ func getRefreshFunction(gui *Gui) func()error {
//                 }
//         } else if key == context.MAIN_STAGING_CONTEXT_KEY {
//                 return func() error {
// -                       selectedLine := gui.Views.Secondary.SelectedLineIdx()
// +                       selectedLine := gui.State.Panels.LineByLine.GetSelectedLineIdx()
//                         return gui.handleRefreshStagingPanel(false, selectedLine)
//                 }
//         } else if key == context.MAIN_PATCH_BUILDING_CONTEXT_KEY {
// `

// func setupGuiForTest(gui *Gui) {
// 	gui.g = &gocui.Gui{}
// 	gui.Views.Main, _ = gui.prepareView("main")
// 	gui.Views.Secondary, _ = gui.prepareView("secondary")
// 	gui.Views.Options, _ = gui.prepareView("options")
// 	gui.git.Patch.PatchManager = &patch.PatchManager{}
// 	_, _ = gui.refreshLineByLinePanel(diffForTest, "", false, 11)
// }

// func TestIncreasesContextInDiffViewByOneInContextWithDiff(t *testing.T) {
// 	contexts := []func(gui *Gui) types.Context{
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Files },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.BranchCommits },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.CommitFiles },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Stash },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Staging },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.PatchBuilding },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.SubCommits },
// 	}

// 	for _, c := range contexts {
// 		gui := NewDummyGui()
// 		context := c(gui)
// 		setupGuiForTest(gui)
// 		gui.c.UserConfig.Git.DiffContextSize = 1
// 		_ = gui.c.PushContext(context)

// 		_ = gui.IncreaseContextInDiffView()

// 		assert.Equal(t, 2, gui.c.UserConfig.Git.DiffContextSize, string(context.GetKey()))
// 	}
// }

// func TestDoesntIncreaseContextInDiffViewInContextWithoutDiff(t *testing.T) {
// 	contexts := []func(gui *Gui) types.Context{
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Status },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Submodules },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Remotes },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Normal },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.ReflogCommits },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.RemoteBranches },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Tags },
// 		// not testing this because it will kick straight back to the files context
// 		// upon pushing the context
// 		// func(gui *Gui) types.Context { return gui.State.Contexts.Merging },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.CommandLog },
// 	}

// 	for _, c := range contexts {
// 		gui := NewDummyGui()
// 		context := c(gui)
// 		setupGuiForTest(gui)
// 		gui.c.UserConfig.Git.DiffContextSize = 1
// 		_ = gui.c.PushContext(context)

// 		_ = gui.IncreaseContextInDiffView()

// 		assert.Equal(t, 1, gui.c.UserConfig.Git.DiffContextSize, string(context.GetKey()))
// 	}
// }

// func TestDecreasesContextInDiffViewByOneInContextWithDiff(t *testing.T) {
// 	contexts := []func(gui *Gui) types.Context{
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Files },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.BranchCommits },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.CommitFiles },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Stash },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Staging },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.PatchBuilding },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.SubCommits },
// 	}

// 	for _, c := range contexts {
// 		gui := NewDummyGui()
// 		context := c(gui)
// 		setupGuiForTest(gui)
// 		gui.c.UserConfig.Git.DiffContextSize = 2
// 		_ = gui.c.PushContext(context)

// 		_ = gui.DecreaseContextInDiffView()

// 		assert.Equal(t, 1, gui.c.UserConfig.Git.DiffContextSize, string(context.GetKey()))
// 	}
// }

// func TestDoesntDecreaseContextInDiffViewInContextWithoutDiff(t *testing.T) {
// 	contexts := []func(gui *Gui) types.Context{
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Status },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Submodules },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Remotes },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Normal },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.ReflogCommits },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.RemoteBranches },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.Tags },
// 		// not testing this because it will kick straight back to the files context
// 		// upon pushing the context
// 		// func(gui *Gui) types.Context { return gui.State.Contexts.Merging },
// 		func(gui *Gui) types.Context { return gui.State.Contexts.CommandLog },
// 	}

// 	for _, c := range contexts {
// 		gui := NewDummyGui()
// 		context := c(gui)
// 		setupGuiForTest(gui)
// 		gui.c.UserConfig.Git.DiffContextSize = 2
// 		_ = gui.c.PushContext(context)

// 		_ = gui.DecreaseContextInDiffView()

// 		assert.Equal(t, 2, gui.c.UserConfig.Git.DiffContextSize, string(context.GetKey()))
// 	}
// }

// func TestDoesntIncreaseContextInDiffViewInContextWhenInPatchBuildingMode(t *testing.T) {
// 	gui := NewDummyGui()
// 	setupGuiForTest(gui)
// 	gui.c.UserConfig.Git.DiffContextSize = 2
// 	_ = gui.c.PushContext(gui.State.Contexts.CommitFiles)
// 	gui.git.Patch.PatchManager.Start("from", "to", false, false)

// 	errorCount := 0
// 	gui.PopupHandler = &popup.TestPopupHandler{
// 		OnErrorMsg: func(message string) error {
// 			assert.Equal(t, gui.c.Tr.CantChangeContextSizeError, message)
// 			errorCount += 1
// 			return nil
// 		},
// 	}

// 	_ = gui.IncreaseContextInDiffView()

// 	assert.Equal(t, 1, errorCount)
// 	assert.Equal(t, 2, gui.c.UserConfig.Git.DiffContextSize)
// }

// func TestDoesntDecreaseContextInDiffViewInContextWhenInPatchBuildingMode(t *testing.T) {
// 	gui := NewDummyGui()
// 	setupGuiForTest(gui)
// 	gui.c.UserConfig.Git.DiffContextSize = 2
// 	_ = gui.c.PushContext(gui.State.Contexts.CommitFiles)
// 	gui.git.Patch.PatchManager.Start("from", "to", false, false)

// 	errorCount := 0
// 	gui.PopupHandler = &popup.TestPopupHandler{
// 		OnErrorMsg: func(message string) error {
// 			assert.Equal(t, gui.c.Tr.CantChangeContextSizeError, message)
// 			errorCount += 1
// 			return nil
// 		},
// 	}

// 	_ = gui.DecreaseContextInDiffView()

// 	assert.Equal(t, 2, gui.c.UserConfig.Git.DiffContextSize)
// }

// func TestDecreasesContextInDiffViewNoFurtherThanOne(t *testing.T) {
// 	gui := NewDummyGui()
// 	setupGuiForTest(gui)
// 	gui.c.UserConfig.Git.DiffContextSize = 1

// 	_ = gui.DecreaseContextInDiffView()

// 	assert.Equal(t, 1, gui.c.UserConfig.Git.DiffContextSize)
// }
