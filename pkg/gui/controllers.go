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

	rebaseHelper := helpers.NewMergeAndRebaseHelper(helperCommon, gui.State.Contexts, gui.git, refsHelper)
	suggestionsHelper := helpers.NewSuggestionsHelper(helperCommon, model, gui.refreshSuggestions)
	gui.helpers = &helpers.Helpers{
		Refs:           refsHelper,
		Host:           helpers.NewHostHelper(helperCommon, gui.git),
		PatchBuilding:  helpers.NewPatchBuildingHelper(helperCommon, gui.git, gui.State.Contexts),
		Bisect:         helpers.NewBisectHelper(helperCommon, gui.git),
		Suggestions:    suggestionsHelper,
		Files:          helpers.NewFilesHelper(helperCommon, gui.git, osCommand),
		WorkingTree:    helpers.NewWorkingTreeHelper(helperCommon, gui.git, model),
		Tags:           helpers.NewTagsHelper(helperCommon, gui.git),
		GPG:            helpers.NewGpgHelper(helperCommon, gui.os, gui.git),
		MergeAndRebase: rebaseHelper,
		MergeConflicts: helpers.NewMergeConflictsHelper(helperCommon, gui.State.Contexts, gui.git),
		CherryPick: helpers.NewCherryPickHelper(
			helperCommon,
			gui.git,
			gui.State.Contexts,
			func() *cherrypicking.CherryPicking { return gui.State.Modes.CherryPicking },
			rebaseHelper,
		),
		Upstream: helpers.NewUpstreamHelper(helperCommon, model, suggestionsHelper.GetRemoteBranchesSuggestionsFunc),
	}

	gui.CustomCommandsClient = custom_commands.NewClient(
		helperCommon,
		gui.os,
		gui.git,
		gui.State.Contexts,
		gui.helpers,
	)

	common := controllers.NewControllerCommon(
		helperCommon,
		osCommand,
		gui.git,
		gui.helpers,
		model,
		gui.State.Contexts,
		gui.State.Modes,
		&gui.Mutexes,
	)

	syncController := controllers.NewSyncController(
		common,
	)

	submodulesController := controllers.NewSubmodulesController(
		common,
		gui.enterSubmodule,
	)

	bisectController := controllers.NewBisectController(common)

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

	menuController := controllers.NewMenuController(common)
	localCommitsController := controllers.NewLocalCommitsController(common, syncController.HandlePull)
	tagsController := controllers.NewTagsController(common)
	filesController := controllers.NewFilesController(
		common,
		gui.enterSubmodule,
		setCommitMessage,
		getSavedCommitMessage,
	)
	mergeConflictsController := controllers.NewMergeConflictsController(common)
	remotesController := controllers.NewRemotesController(
		common,
		func(branches []*models.RemoteBranch) { gui.State.Model.RemoteBranches = branches },
	)
	undoController := controllers.NewUndoController(common)
	globalController := controllers.NewGlobalController(common)
	contextLinesController := controllers.NewContextLinesController(common)
	verticalScrollControllerFactory := controllers.NewVerticalScrollControllerFactory(common)

	branchesController := controllers.NewBranchesController(common)
	gitFlowController := controllers.NewGitFlowController(common)
	filesRemoveController := controllers.NewFilesRemoveController(common)
	stashController := controllers.NewStashController(common)
	commitFilesController := controllers.NewCommitFilesController(common)
	patchExplorerControllerFactory := controllers.NewPatchExplorerControllerFactory(common)
	stagingController := controllers.NewStagingController(common, gui.State.Contexts.Staging, gui.State.Contexts.StagingSecondary, false)
	stagingSecondaryController := controllers.NewStagingController(common, gui.State.Contexts.StagingSecondary, gui.State.Contexts.Staging, true)
	patchBuildingController := controllers.NewPatchBuildingController(common)

	setSubCommits := func(commits []*models.Commit) { gui.State.Model.SubCommits = commits }

	for _, context := range []controllers.CanSwitchToSubCommits{
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.ReflogCommits,
	} {
		controllers.AttachControllers(context, controllers.NewSwitchToSubCommitsController(
			common, setSubCommits, context,
		))
	}

	for _, context := range []controllers.CanSwitchToDiffFiles{
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
	} {
		controllers.AttachControllers(context, controllers.NewSwitchToDiffFilesController(
			common, gui.SwitchToCommitFilesContext, context,
		))
	}

	for _, context := range []controllers.ContainsCommits{
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
	} {
		controllers.AttachControllers(context, controllers.NewBasicCommitsController(common, context))
	}

	// TODO: add scroll controllers for main panels (need to bring some more functionality across for that e.g. reading more from the currently displayed git command)
	controllers.AttachControllers(gui.State.Contexts.Staging,
		stagingController,
		patchExplorerControllerFactory.Create(gui.State.Contexts.Staging),
		verticalScrollControllerFactory.Create(gui.State.Contexts.Staging),
	)

	controllers.AttachControllers(gui.State.Contexts.StagingSecondary,
		stagingSecondaryController,
		patchExplorerControllerFactory.Create(gui.State.Contexts.StagingSecondary),
		verticalScrollControllerFactory.Create(gui.State.Contexts.StagingSecondary),
	)

	controllers.AttachControllers(gui.State.Contexts.CustomPatchBuilder,
		patchBuildingController,
		patchExplorerControllerFactory.Create(gui.State.Contexts.CustomPatchBuilder),
		verticalScrollControllerFactory.Create(gui.State.Contexts.CustomPatchBuilder),
	)

	controllers.AttachControllers(gui.State.Contexts.CustomPatchBuilderSecondary,
		verticalScrollControllerFactory.Create(gui.State.Contexts.CustomPatchBuilder),
	)

	controllers.AttachControllers(gui.State.Contexts.MergeConflicts,
		mergeConflictsController,
	)

	controllers.AttachControllers(gui.State.Contexts.Files,
		filesController,
		filesRemoveController,
	)

	controllers.AttachControllers(gui.State.Contexts.Tags,
		tagsController,
	)

	controllers.AttachControllers(gui.State.Contexts.Submodules,
		submodulesController,
	)

	controllers.AttachControllers(gui.State.Contexts.LocalCommits,
		localCommitsController,
		bisectController,
	)

	controllers.AttachControllers(gui.State.Contexts.Branches,
		branchesController,
		gitFlowController,
	)

	controllers.AttachControllers(gui.State.Contexts.LocalCommits,
		localCommitsController,
		bisectController,
	)

	controllers.AttachControllers(gui.State.Contexts.CommitFiles,
		commitFilesController,
	)

	controllers.AttachControllers(gui.State.Contexts.Remotes,
		remotesController,
	)

	controllers.AttachControllers(gui.State.Contexts.Stash,
		stashController,
	)

	controllers.AttachControllers(gui.State.Contexts.Menu,
		menuController,
	)

	controllers.AttachControllers(gui.State.Contexts.CommitMessage,
		commitMessageController,
	)

	controllers.AttachControllers(gui.State.Contexts.RemoteBranches,
		remoteBranchesController,
	)

	controllers.AttachControllers(gui.State.Contexts.Global,
		syncController,
		undoController,
		globalController,
		contextLinesController,
	)

	// this must come last so that we've got our click handlers defined against the context
	listControllerFactory := controllers.NewListControllerFactory(gui.c)
	for _, context := range gui.getListContexts() {
		controllers.AttachControllers(context, listControllerFactory.Create(context))
	}
}
