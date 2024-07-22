package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/services/custom_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) Helpers() *helpers.Helpers {
	return gui.helpers
}

// Note, the order of controllers determines the order in which keybindings appear
// in the keybinding menu: the earlier that the controller is attached to a context,
// the lower in the list the keybindings will appear.
func (gui *Gui) resetHelpersAndControllers() {
	for _, context := range gui.Contexts().Flatten() {
		context.ClearAllBindingsFn()
	}

	helperCommon := gui.c
	recordDirectoryHelper := helpers.NewRecordDirectoryHelper(helperCommon)
	reposHelper := helpers.NewRecentReposHelper(helperCommon, recordDirectoryHelper, gui.onNewRepo)
	refsHelper := helpers.NewRefsHelper(helperCommon)
	suggestionsHelper := helpers.NewSuggestionsHelper(helperCommon)
	worktreeHelper := helpers.NewWorktreeHelper(helperCommon, reposHelper, refsHelper, suggestionsHelper)

	rebaseHelper := helpers.NewMergeAndRebaseHelper(helperCommon, refsHelper)

	setCommitSummary := gui.getCommitMessageSetTextareaTextFn(func() *gocui.View { return gui.Views.CommitMessage })
	setCommitDescription := gui.getCommitMessageSetTextareaTextFn(func() *gocui.View { return gui.Views.CommitDescription })
	getCommitSummary := func() string {
		return strings.TrimSpace(gui.Views.CommitMessage.TextArea.GetContent())
	}

	getCommitDescription := func() string {
		return strings.TrimSpace(gui.Views.CommitDescription.TextArea.GetContent())
	}
	getUnwrappedCommitDescription := func() string {
		return strings.TrimSpace(gui.Views.CommitDescription.TextArea.GetUnwrappedContent())
	}
	commitsHelper := helpers.NewCommitsHelper(helperCommon,
		getCommitSummary,
		setCommitSummary,
		getCommitDescription,
		getUnwrappedCommitDescription,
		setCommitDescription,
	)

	gpgHelper := helpers.NewGpgHelper(helperCommon)
	viewHelper := helpers.NewViewHelper(helperCommon, gui.State.Contexts)
	patchBuildingHelper := helpers.NewPatchBuildingHelper(helperCommon)
	stagingHelper := helpers.NewStagingHelper(helperCommon)
	mergeConflictsHelper := helpers.NewMergeConflictsHelper(helperCommon)
	searchHelper := helpers.NewSearchHelper(helperCommon)

	refreshHelper := helpers.NewRefreshHelper(
		helperCommon,
		refsHelper,
		rebaseHelper,
		patchBuildingHelper,
		stagingHelper,
		mergeConflictsHelper,
		worktreeHelper,
		searchHelper,
	)
	diffHelper := helpers.NewDiffHelper(helperCommon)
	cherryPickHelper := helpers.NewCherryPickHelper(
		helperCommon,
		rebaseHelper,
	)
	bisectHelper := helpers.NewBisectHelper(helperCommon)
	windowHelper := helpers.NewWindowHelper(helperCommon, viewHelper)
	modeHelper := helpers.NewModeHelper(
		helperCommon,
		diffHelper,
		patchBuildingHelper,
		cherryPickHelper,
		rebaseHelper,
		bisectHelper,
	)
	appStatusHelper := helpers.NewAppStatusHelper(
		helperCommon,
		func() *status.StatusManager { return gui.statusManager },
		modeHelper,
	)

	setSubCommits := func(commits []*models.Commit) {
		gui.Mutexes.SubCommitsMutex.Lock()
		defer gui.Mutexes.SubCommitsMutex.Unlock()

		gui.State.Model.SubCommits = commits
	}
	gui.helpers = &helpers.Helpers{
		Refs:            refsHelper,
		Host:            helpers.NewHostHelper(helperCommon),
		PatchBuilding:   patchBuildingHelper,
		Staging:         stagingHelper,
		Bisect:          bisectHelper,
		Suggestions:     suggestionsHelper,
		Files:           helpers.NewFilesHelper(helperCommon),
		WorkingTree:     helpers.NewWorkingTreeHelper(helperCommon, refsHelper, commitsHelper, gpgHelper),
		Tags:            helpers.NewTagsHelper(helperCommon, commitsHelper),
		BranchesHelper:  helpers.NewBranchesHelper(helperCommon),
		GPG:             helpers.NewGpgHelper(helperCommon),
		MergeAndRebase:  rebaseHelper,
		MergeConflicts:  mergeConflictsHelper,
		CherryPick:      cherryPickHelper,
		Upstream:        helpers.NewUpstreamHelper(helperCommon, suggestionsHelper.GetRemoteBranchesSuggestionsFunc),
		AmendHelper:     helpers.NewAmendHelper(helperCommon, gpgHelper),
		FixupHelper:     helpers.NewFixupHelper(helperCommon),
		Commits:         commitsHelper,
		Snake:           helpers.NewSnakeHelper(helperCommon),
		Diff:            diffHelper,
		Repos:           reposHelper,
		RecordDirectory: recordDirectoryHelper,
		Update:          helpers.NewUpdateHelper(helperCommon, gui.Updater),
		Window:          windowHelper,
		View:            viewHelper,
		Refresh:         refreshHelper,
		Confirmation:    helpers.NewConfirmationHelper(helperCommon),
		Mode:            modeHelper,
		AppStatus:       appStatusHelper,
		InlineStatus:    helpers.NewInlineStatusHelper(helperCommon, windowHelper),
		WindowArrangement: helpers.NewWindowArrangementHelper(
			gui.c,
			windowHelper,
			modeHelper,
			appStatusHelper,
		),
		Search:     searchHelper,
		Worktree:   worktreeHelper,
		SubCommits: helpers.NewSubCommitsHelper(helperCommon, refreshHelper, setSubCommits),
	}

	gui.CustomCommandsClient = custom_commands.NewClient(
		helperCommon,
		gui.helpers,
	)

	common := controllers.NewControllerCommon(helperCommon, gui)

	syncController := controllers.NewSyncController(
		common,
	)

	submodulesController := controllers.NewSubmodulesController(common)

	bisectController := controllers.NewBisectController(common)

	commitMessageController := controllers.NewCommitMessageController(
		common,
	)

	commitDescriptionController := controllers.NewCommitDescriptionController(
		common,
	)

	remoteBranchesController := controllers.NewRemoteBranchesController(common)

	menuController := controllers.NewMenuController(common)
	localCommitsController := controllers.NewLocalCommitsController(common, syncController.HandlePull)
	tagsController := controllers.NewTagsController(common)
	filesController := controllers.NewFilesController(
		common,
	)
	mergeConflictsController := controllers.NewMergeConflictsController(common)
	remotesController := controllers.NewRemotesController(
		common,
		func(branches []*models.RemoteBranch) { gui.State.Model.RemoteBranches = branches },
	)
	worktreesController := controllers.NewWorktreesController(common)
	undoController := controllers.NewUndoController(common)
	globalController := controllers.NewGlobalController(common)
	contextLinesController := controllers.NewContextLinesController(common)
	renameSimilarityThresholdController := controllers.NewRenameSimilarityThresholdController(common)
	verticalScrollControllerFactory := controllers.NewVerticalScrollControllerFactory(common, &gui.viewBufferManagerMap)

	branchesController := controllers.NewBranchesController(common)
	gitFlowController := controllers.NewGitFlowController(common)
	stashController := controllers.NewStashController(common)
	commitFilesController := controllers.NewCommitFilesController(common)
	patchExplorerControllerFactory := controllers.NewPatchExplorerControllerFactory(common)
	stagingController := controllers.NewStagingController(common, gui.State.Contexts.Staging, gui.State.Contexts.StagingSecondary, false)
	stagingSecondaryController := controllers.NewStagingController(common, gui.State.Contexts.StagingSecondary, gui.State.Contexts.Staging, true)
	patchBuildingController := controllers.NewPatchBuildingController(common)
	snakeController := controllers.NewSnakeController(common)
	reflogCommitsController := controllers.NewReflogCommitsController(common)
	subCommitsController := controllers.NewSubCommitsController(common)
	statusController := controllers.NewStatusController(common)
	commandLogController := controllers.NewCommandLogController(common)
	confirmationController := controllers.NewConfirmationController(common)
	suggestionsController := controllers.NewSuggestionsController(common)
	jumpToSideWindowController := controllers.NewJumpToSideWindowController(common)

	sideWindowControllerFactory := controllers.NewSideWindowControllerFactory(common)

	filterControllerFactory := controllers.NewFilterControllerFactory(common)
	for _, context := range gui.c.Context().AllFilterable() {
		controllers.AttachControllers(context, filterControllerFactory.Create(context))
	}

	searchControllerFactory := controllers.NewSearchControllerFactory(common)
	for _, context := range gui.c.Context().AllSearchable() {
		controllers.AttachControllers(context, searchControllerFactory.Create(context))
	}

	for _, context := range []controllers.CanViewWorktreeOptions{
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
	} {
		controllers.AttachControllers(context, controllers.NewWorktreeOptionsController(common, context))
	}

	// allow for navigating between side window contexts
	for _, context := range []types.Context{
		gui.State.Contexts.Status,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.Worktrees,
		gui.State.Contexts.Tags,
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Files,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.LocalCommits,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
	} {
		controllers.AttachControllers(context, sideWindowControllerFactory.Create(context))
	}

	for _, context := range []controllers.CanSwitchToSubCommits{
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.ReflogCommits,
	} {
		controllers.AttachControllers(context, controllers.NewSwitchToSubCommitsController(
			common, context,
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
	)

	controllers.AttachControllers(gui.State.Contexts.Tags,
		tagsController,
	)

	controllers.AttachControllers(gui.State.Contexts.Submodules,
		submodulesController,
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

	controllers.AttachControllers(gui.State.Contexts.Worktrees,
		worktreesController,
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

	controllers.AttachControllers(gui.State.Contexts.CommitDescription,
		commitDescriptionController,
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

	controllers.AttachControllers(gui.State.Contexts.Search,
		controllers.NewSearchPromptController(common),
	)

	controllers.AttachControllers(gui.State.Contexts.Global,
		undoController,
		globalController,
		contextLinesController,
		renameSimilarityThresholdController,
		jumpToSideWindowController,
		syncController,
	)

	controllers.AttachControllers(gui.State.Contexts.Snake,
		snakeController,
	)

	// this must come last so that we've got our click handlers defined against the context
	listControllerFactory := controllers.NewListControllerFactory(common)
	for _, context := range gui.c.Context().AllList() {
		controllers.AttachControllers(context, listControllerFactory.Create(context))
	}
}

func (gui *Gui) getCommitMessageSetTextareaTextFn(getView func() *gocui.View) func(string) {
	return func(text string) {
		// using a getView function so that we don't need to worry about when the view is created
		view := getView()
		view.ClearTextArea()
		view.TextArea.TypeString(text)
		view.RenderTextArea()
	}
}
