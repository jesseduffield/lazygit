package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/services/custom_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
	suggestionsHelper := helpers.NewSuggestionsHelper(helperCommon, model, gui.State.Contexts)
	setCommitMessage := gui.getSetTextareaTextFn(func() *gocui.View { return gui.Views.CommitMessage })
	getSavedCommitMessage := func() string {
		return gui.State.savedCommitMessage
	}
	gpgHelper := helpers.NewGpgHelper(helperCommon, gui.os, gui.git)
	viewHelper := helpers.NewViewHelper(helperCommon, gui.State.Contexts)
	recordDirectoryHelper := helpers.NewRecordDirectoryHelper(helperCommon)
	patchBuildingHelper := helpers.NewPatchBuildingHelper(helperCommon, gui.git, gui.State.Contexts)
	stagingHelper := helpers.NewStagingHelper(helperCommon, gui.git, gui.State.Contexts)
	mergeConflictsHelper := helpers.NewMergeConflictsHelper(helperCommon, gui.State.Contexts, gui.git)
	refreshHelper := helpers.NewRefreshHelper(helperCommon, gui.State.Contexts, gui.git, refsHelper, rebaseHelper, patchBuildingHelper, stagingHelper, mergeConflictsHelper, gui.fileWatcher)
	gui.helpers = &helpers.Helpers{
		Refs:           refsHelper,
		Host:           helpers.NewHostHelper(helperCommon, gui.git),
		PatchBuilding:  patchBuildingHelper,
		Staging:        stagingHelper,
		Bisect:         helpers.NewBisectHelper(helperCommon),
		Suggestions:    suggestionsHelper,
		Files:          helpers.NewFilesHelper(helperCommon, gui.git, osCommand),
		WorkingTree:    helpers.NewWorkingTreeHelper(helperCommon, gui.git, gui.State.Contexts, refsHelper, model, setCommitMessage, getSavedCommitMessage),
		Tags:           helpers.NewTagsHelper(helperCommon, gui.git),
		GPG:            gpgHelper,
		MergeAndRebase: rebaseHelper,
		MergeConflicts: mergeConflictsHelper,
		CherryPick: helpers.NewCherryPickHelper(
			helperCommon,
			gui.State.Contexts,
			rebaseHelper,
		),
		Upstream:        helpers.NewUpstreamHelper(helperCommon, model, suggestionsHelper.GetRemoteBranchesSuggestionsFunc),
		AmendHelper:     helpers.NewAmendHelper(helperCommon, gui.git, gpgHelper),
		Snake:           helpers.NewSnakeHelper(helperCommon),
		Diff:            helpers.NewDiffHelper(helperCommon),
		Repos:           helpers.NewRecentReposHelper(helperCommon, recordDirectoryHelper, gui.onNewRepo),
		RecordDirectory: recordDirectoryHelper,
		Update:          helpers.NewUpdateHelper(helperCommon, gui.Updater),
		Window:          helpers.NewWindowHelper(helperCommon, viewHelper, gui.State.Contexts),
		View:            viewHelper,
		Refresh:         refreshHelper,
		Confirmation:    helpers.NewConfirmationHelper(helperCommon, gui.State.Contexts),
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

	submodulesController := controllers.NewSubmodulesController(common)

	bisectController := controllers.NewBisectController(common)

	getCommitMessage := func() string {
		return strings.TrimSpace(gui.Views.CommitMessage.TextArea.GetContent())
	}

	onCommitAttempt := func(message string) {
		gui.State.savedCommitMessage = message
		gui.Views.CommitMessage.ClearTextArea()
	}

	onCommitSuccess := func() {
		gui.State.savedCommitMessage = ""
		_ = gui.c.Refresh(types.RefreshOptions{
			Scope: []types.RefreshableView{types.STAGING},
		})
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
	verticalScrollControllerFactory := controllers.NewVerticalScrollControllerFactory(common, &gui.viewBufferManagerMap)

	branchesController := controllers.NewBranchesController(common)
	gitFlowController := controllers.NewGitFlowController(common)
	filesRemoveController := controllers.NewFilesRemoveController(common)
	stashController := controllers.NewStashController(common)
	commitFilesController := controllers.NewCommitFilesController(common)
	patchExplorerControllerFactory := controllers.NewPatchExplorerControllerFactory(common)
	stagingController := controllers.NewStagingController(common, gui.State.Contexts.Staging, gui.State.Contexts.StagingSecondary, false)
	stagingSecondaryController := controllers.NewStagingController(common, gui.State.Contexts.StagingSecondary, gui.State.Contexts.Staging, true)
	patchBuildingController := controllers.NewPatchBuildingController(common)
	snakeController := controllers.NewSnakeController(common)
	reflogCommitsController := controllers.NewReflogCommitsController(common, gui.State.Contexts.ReflogCommits)
	subCommitsController := controllers.NewSubCommitsController(common, gui.State.Contexts.SubCommits)
	statusController := controllers.NewStatusController(common)
	commandLogController := controllers.NewCommandLogController(common)
	confirmationController := controllers.NewConfirmationController(common)
	suggestionsController := controllers.NewSuggestionsController(common)

	setSubCommits := func(commits []*models.Commit) {
		gui.Mutexes.SubCommitsMutex.Lock()
		defer gui.Mutexes.SubCommitsMutex.Unlock()

		gui.State.Model.SubCommits = commits
	}

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
			common, context, gui.State.Contexts.CommitFiles,
		))
	}

	for _, context := range []controllers.ContainsCommits{
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
	} {
		controllers.AttachControllers(context, controllers.NewBasicCommitsController(common, context))
	}

	controllers.AttachControllers(gui.State.Contexts.ReflogCommits,
		reflogCommitsController,
	)

	controllers.AttachControllers(gui.State.Contexts.SubCommits,
		subCommitsController,
	)

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

	controllers.AttachControllers(gui.State.Contexts.Status,
		statusController,
	)

	controllers.AttachControllers(gui.State.Contexts.CommandLog,
		commandLogController,
	)

	controllers.AttachControllers(gui.State.Contexts.Confirmation,
		confirmationController,
	)

	controllers.AttachControllers(gui.State.Contexts.Suggestions,
		suggestionsController,
	)

	controllers.AttachControllers(gui.State.Contexts.Global,
		syncController,
		undoController,
		globalController,
		contextLinesController,
	)

	controllers.AttachControllers(gui.State.Contexts.Snake,
		snakeController,
	)

	// this must come last so that we've got our click handlers defined against the context
	listControllerFactory := controllers.NewListControllerFactory(common)
	for _, context := range gui.getListContexts() {
		controllers.AttachControllers(context, listControllerFactory.Create(context))
	}
}

func (gui *Gui) getSetTextareaTextFn(getView func() *gocui.View) func(string) {
	return func(text string) {
		// using a getView function so that we don't need to worry about when the view is created
		view := getView()
		view.ClearTextArea()
		view.TextArea.TypeString(text)
		_ = gui.helpers.Confirmation.ResizePopupPanel(view, view.TextArea.GetContent())
		view.RenderTextArea()
	}
}
