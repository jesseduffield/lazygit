package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/services/custom_commands"
)

func (gui *Gui) resetControllers() {
	helperCommon := gui.c
	osCommand := gui.os
	model := gui.State.Model
	refsHelper := helpers.NewRefsHelper(
		helperCommon,
		gui.git,
		gui.State.Contexts,
		model,
	)

	rebaseHelper := helpers.NewMergeAndRebaseHelper(helperCommon, gui.State.Contexts, gui.git, gui.takeOverMergeConflictScrolling, refsHelper)
	gui.helpers = &helpers.Helpers{
		Refs:           refsHelper,
		PatchBuilding:  helpers.NewPatchBuildingHelper(helperCommon, gui.git),
		Bisect:         helpers.NewBisectHelper(helperCommon, gui.git),
		Suggestions:    helpers.NewSuggestionsHelper(helperCommon, model, gui.refreshSuggestions),
		Files:          helpers.NewFilesHelper(helperCommon, gui.git, osCommand),
		WorkingTree:    helpers.NewWorkingTreeHelper(model),
		Tags:           helpers.NewTagsHelper(helperCommon, gui.git),
		GPG:            helpers.NewGpgHelper(helperCommon, gui.os, gui.git),
		MergeAndRebase: rebaseHelper,
		CherryPick: helpers.NewCherryPickHelper(
			helperCommon,
			gui.git,
			gui.State.Contexts,
			func() *cherrypicking.CherryPicking { return gui.State.Modes.CherryPicking },
			rebaseHelper,
		),
	}

	gui.CustomCommandsClient = custom_commands.NewClient(
		helperCommon,
		gui.os,
		gui.git,
		gui.State.Contexts,
		gui.helpers,
		gui.getKey,
	)

	common := controllers.NewControllerCommon(
		helperCommon,
		osCommand,
		gui.git,
		gui.helpers,
		model,
		gui.State.Contexts,
		gui.State.Modes,
	)

	syncController := controllers.NewSyncController(
		common,
		gui.getSuggestedRemote,
	)

	submodulesController := controllers.NewSubmodulesController(
		common,
		gui.enterSubmodule,
	)

	bisectController := controllers.NewBisectController(common)

	reflogController := controllers.NewReflogController(common)
	subCommitsController := controllers.NewSubCommitsController(common)

	getSavedCommitMessage := func() string {
		return gui.State.savedCommitMessage
	}

	getCommitMessage := func() string {
		return strings.TrimSpace(gui.Views.CommitMessage.TextArea.GetContent())
	}

	setCommitMessage := gui.getSetTextareaTextFn(func() *gocui.View { return gui.Views.CommitMessage })

	onCommitAttempt := func(message string) {
		gui.State.savedCommitMessage = message
		gui.Views.CommitMessage.ClearTextArea()
	}

	onCommitSuccess := func() {
		gui.State.savedCommitMessage = ""
	}

	commitMessageController := controllers.NewCommitMessageController(
		common,
		getCommitMessage,
		onCommitAttempt,
		onCommitSuccess,
	)

	remoteBranchesController := controllers.NewRemoteBranchesController(common)

	gui.Controllers = Controllers{
		Submodules: submodulesController,
		Global:     controllers.NewGlobalController(common),
		Files: controllers.NewFilesController(
			common,
			gui.enterSubmodule,
			setCommitMessage,
			getSavedCommitMessage,
			gui.switchToMerge,
		),
		Tags:         controllers.NewTagsController(common),
		LocalCommits: controllers.NewLocalCommitsController(common, syncController.HandlePull),
		Remotes: controllers.NewRemotesController(
			common,
			func(branches []*models.RemoteBranch) { gui.State.Model.RemoteBranches = branches },
		),
		Menu: controllers.NewMenuController(common),
		Undo: controllers.NewUndoController(common),
		Sync: syncController,
	}

	branchesController := controllers.NewBranchesController(common)
	gitFlowController := controllers.NewGitFlowController(common)
	filesRemoveController := controllers.NewFilesRemoveController(common)
	stashController := controllers.NewStashController(common)
	commitFilesController := controllers.NewCommitFilesController(common)

	switchToSubCommitsControllerFactory := controllers.NewSubCommitsSwitchControllerFactory(
		common,
		func(commits []*models.Commit) { gui.State.Model.SubCommits = commits },
	)

	for _, context := range []controllers.ContextWithRefName{
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
	} {
		controllers.AttachControllers(context, switchToSubCommitsControllerFactory.Create(context))
	}

	commitishControllerFactory := controllers.NewCommitishControllerFactory(
		common,
		gui.SwitchToCommitFilesContext,
	)

	for _, context := range []controllers.Commitish{
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
	} {
		controllers.AttachControllers(context, commitishControllerFactory.Create(context))
	}

	controllers.AttachControllers(gui.State.Contexts.Branches, branchesController, gitFlowController)
	controllers.AttachControllers(gui.State.Contexts.Files, gui.Controllers.Files, filesRemoveController)
	controllers.AttachControllers(gui.State.Contexts.Tags, gui.Controllers.Tags)
	controllers.AttachControllers(gui.State.Contexts.Submodules, gui.Controllers.Submodules)
	controllers.AttachControllers(gui.State.Contexts.LocalCommits, gui.Controllers.LocalCommits, bisectController)
	controllers.AttachControllers(gui.State.Contexts.ReflogCommits, reflogController)
	controllers.AttachControllers(gui.State.Contexts.SubCommits, subCommitsController)
	controllers.AttachControllers(gui.State.Contexts.CommitFiles, commitFilesController)
	controllers.AttachControllers(gui.State.Contexts.Remotes, gui.Controllers.Remotes)
	controllers.AttachControllers(gui.State.Contexts.Stash, stashController)
	controllers.AttachControllers(gui.State.Contexts.Menu, gui.Controllers.Menu)
	controllers.AttachControllers(gui.State.Contexts.CommitMessage, commitMessageController)
	controllers.AttachControllers(gui.State.Contexts.RemoteBranches, remoteBranchesController)
	controllers.AttachControllers(gui.State.Contexts.Global, gui.Controllers.Sync, gui.Controllers.Undo, gui.Controllers.Global)

	// this must come last so that we've got our click handlers defined against the context
	listControllerFactory := controllers.NewListControllerFactory(gui.c)
	for _, context := range gui.getListContexts() {
		controllers.AttachControllers(context, listControllerFactory.Create(context))
	}
}
