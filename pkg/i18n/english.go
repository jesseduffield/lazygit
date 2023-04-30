/*

Todo list when making a new translation
- Copy this file and rename it to the language you want to translate to like someLanguage.go
- Change the addEnglish() name to the language you want to translate to like addSomeLanguage()
- change the first function argument of i18nObject.AddMessages( to the language you want to translate to like language.SomeLanguage
- Remove this todo and the about section

*/

package i18n

type TranslationSet struct {
	NotEnoughSpace                      string
	DiffTitle                           string
	FilesTitle                          string
	BranchesTitle                       string
	CommitsTitle                        string
	StashTitle                          string
	SnakeTitle                          string
	EasterEgg                           string
	UnstagedChanges                     string
	StagedChanges                       string
	MainTitle                           string
	StagingTitle                        string
	MergingTitle                        string
	MergeConfirmTitle                   string
	NormalTitle                         string
	LogTitle                            string
	CommitSummary                       string
	CredentialsUsername                 string
	CredentialsPassword                 string
	CredentialsPassphrase               string
	CredentialsPIN                      string
	PassUnameWrong                      string
	CommitChanges                       string
	AmendLastCommit                     string
	AmendLastCommitTitle                string
	SureToAmend                         string
	NoCommitToAmend                     string
	CommitChangesWithEditor             string
	StatusTitle                         string
	GlobalTitle                         string
	LcNavigate                          string
	LcMenu                              string
	LcExecute                           string
	LcToggleStaged                      string
	LcToggleStagedAll                   string
	LcToggleTreeView                    string
	LcOpenMergeTool                     string
	LcRefresh                           string
	LcPush                              string
	LcPull                              string
	LcScroll                            string
	LcFileFilter                        string
	FilterStagedFiles                   string
	FilterUnstagedFiles                 string
	ResetCommitFilterState              string
	MergeConflictsTitle                 string
	LcCheckout                          string
	NoChangedFiles                      string
	PullWait                            string
	PushWait                            string
	FetchWait                           string
	LcSoftReset                         string
	AlreadyCheckedOutBranch             string
	SureForceCheckout                   string
	ForceCheckoutBranch                 string
	BranchName                          string
	NewBranchNameBranchOff              string
	CantDeleteCheckOutBranch            string
	DeleteBranch                        string
	DeleteBranchMessage                 string
	ForceDeleteBranchMessage            string
	LcRebaseBranch                      string
	CantRebaseOntoSelf                  string
	CantMergeBranchIntoItself           string
	LcForceCheckout                     string
	LcCheckoutByName                    string
	LcNewBranch                         string
	LcDeleteBranch                      string
	NoBranchesThisRepo                  string
	CommitMessageConfirm                string
	CommitWithoutMessageErr             string
	CloseConfirm                        string
	LcClose                             string
	LcQuit                              string
	LcSquashDown                        string
	LcFixupCommit                       string
	CannotSquashOrFixupFirstCommit      string
	Fixup                               string
	SureFixupThisCommit                 string
	SureSquashThisCommit                string
	Squash                              string
	LcPickCommit                        string
	LcRevertCommit                      string
	LcRewordCommit                      string
	LcDeleteCommit                      string
	LcMoveDownCommit                    string
	LcMoveUpCommit                      string
	LcEditCommit                        string
	LcAmendToCommit                     string
	LcResetCommitAuthor                 string
	SetAuthorPromptTitle                string
	SureResetCommitAuthor               string
	LcRenameCommitEditor                string
	NoCommitsThisBranch                 string
	UpdateRefHere                       string
	Error                               string
	LcSelectHunk                        string
	LcNavigateConflicts                 string
	LcPickHunk                          string
	LcPickAllHunks                      string
	LcUndo                              string
	LcUndoReflog                        string
	LcRedoReflog                        string
	UndoTooltip                         string
	RedoTooltip                         string
	DiscardAllTooltip                   string
	DiscardUnstagedTooltip              string
	LcPop                               string
	LcDrop                              string
	LcApply                             string
	NoStashEntries                      string
	StashDrop                           string
	SureDropStashEntry                  string
	StashPop                            string
	SurePopStashEntry                   string
	StashApply                          string
	SureApplyStashEntry                 string
	NoTrackedStagedFilesStash           string
	NoFilesToStash                      string
	StashChanges                        string
	LcRenameStash                       string
	RenameStashPrompt                   string
	OpenConfig                          string
	EditConfig                          string
	ForcePush                           string
	ForcePushPrompt                     string
	ForcePushDisabled                   string
	UpdatesRejectedAndForcePushDisabled string
	LcCheckForUpdate                    string
	CheckingForUpdates                  string
	UpdateAvailableTitle                string
	UpdateAvailable                     string
	UpdateInProgressWaitingStatus       string
	UpdateCompletedTitle                string
	UpdateCompleted                     string
	FailedToRetrieveLatestVersionErr    string
	OnLatestVersionErr                  string
	MajorVersionErr                     string
	CouldNotFindBinaryErr               string
	UpdateFailedErr                     string
	ConfirmQuitDuringUpdateTitle        string
	ConfirmQuitDuringUpdate             string
	MergeToolTitle                      string
	MergeToolPrompt                     string
	IntroPopupMessage                   string
	DeprecatedEditConfigWarning         string
	GitconfigParseErr                   string
	LcEditFile                          string
	LcOpenFile                          string
	LcIgnoreFile                        string
	LcExcludeFile                       string
	LcRefreshFiles                      string
	LcMergeIntoCurrentBranch            string
	ConfirmQuit                         string
	SwitchRepo                          string
	LcAllBranchesLogGraph               string
	UnsupportedGitService               string
	LcCreatePullRequest                 string
	LcCopyPullRequestURL                string
	NoBranchOnRemote                    string
	LcFetch                             string
	NoAutomaticGitFetchTitle            string
	NoAutomaticGitFetchBody             string
	FileEnter                           string
	FileStagingRequirements             string
	StageSelection                      string
	ResetSelection                      string
	ToggleDragSelect                    string
	ToggleSelectHunk                    string
	ToggleSelectionForPatch             string
	EditHunk                            string
	ToggleStagingPanel                  string
	ReturnToFilesPanel                  string
	FastForward                         string
	Fetching                            string
	FoundConflicts                      string
	FoundConflictsTitle                 string
	PickHunk                            string
	PickAllHunks                        string
	ViewMergeRebaseOptions              string
	NotMergingOrRebasing                string
	AlreadyRebasing                     string
	RecentRepos                         string
	MergeOptionsTitle                   string
	RebaseOptionsTitle                  string
	CommitMessageTitle                  string
	CommitDescriptionTitle              string
	CommitDescriptionSubTitle           string
	LocalBranchesTitle                  string
	SearchTitle                         string
	TagsTitle                           string
	MenuTitle                           string
	RemotesTitle                        string
	RemoteBranchesTitle                 string
	PatchBuildingTitle                  string
	InformationTitle                    string
	SecondaryTitle                      string
	ReflogCommitsTitle                  string
	ConflictsResolved                   string
	RebasingTitle                       string
	SimpleRebase                        string
	InteractiveRebase                   string
	InteractiveRebaseTooltip            string
	ConfirmMerge                        string
	FwdNoUpstream                       string
	FwdNoLocalUpstream                  string
	FwdCommitsToPush                    string
	ErrorOccurred                       string
	NoRoom                              string
	YouAreHere                          string
	YouDied                             string
	LcRewordNotSupported                string
	LcChangingThisActionIsNotAllowed    string
	LcCherryPickCopy                    string
	LcCherryPickCopyRange               string
	LcPasteCommits                      string
	SureCherryPick                      string
	CherryPick                          string
	Donate                              string
	AskQuestion                         string
	PrevLine                            string
	NextLine                            string
	PrevHunk                            string
	NextHunk                            string
	PrevConflict                        string
	NextConflict                        string
	SelectPrevHunk                      string
	SelectNextHunk                      string
	ScrollDown                          string
	ScrollUp                            string
	LcScrollUpMainPanel                 string
	LcScrollDownMainPanel               string
	AmendCommitTitle                    string
	AmendCommitPrompt                   string
	DeleteCommitTitle                   string
	DeleteCommitPrompt                  string
	SquashingStatus                     string
	FixingStatus                        string
	DeletingStatus                      string
	MovingStatus                        string
	RebasingStatus                      string
	AmendingStatus                      string
	CherryPickingStatus                 string
	UndoingStatus                       string
	RedoingStatus                       string
	CheckingOutStatus                   string
	CommittingStatus                    string
	CommitFiles                         string
	SubCommitsDynamicTitle              string
	CommitFilesDynamicTitle             string
	RemoteBranchesDynamicTitle          string
	LcViewItemFiles                     string
	CommitFilesTitle                    string
	LcCheckoutCommitFile                string
	LcDiscardOldFileChange              string
	DiscardFileChangesTitle             string
	DiscardFileChangesPrompt            string
	DisabledForGPG                      string
	CreateRepo                          string
	BareRepo                            string
	InitialBranch                       string
	NoRecentRepositories                string
	IncorrectNotARepository             string
	AutoStashTitle                      string
	AutoStashPrompt                     string
	StashPrefix                         string
	LcViewDiscardOptions                string
	LcCancel                            string
	LcDiscardAllChanges                 string
	LcDiscardUnstagedChanges            string
	LcDiscardAllChangesToAllFiles       string
	LcDiscardAnyUnstagedChanges         string
	LcDiscardUntrackedFiles             string
	LcDiscardStagedChanges              string
	LcHardReset                         string
	LcViewResetOptions                  string
	LcCreateFixupCommit                 string
	LcSquashAboveCommits                string
	SquashAboveCommits                  string
	SureSquashAboveCommits              string
	CreateFixupCommit                   string
	SureCreateFixupCommit               string
	LcExecuteCustomCommand              string
	CustomCommand                       string
	LcCommitChangesWithoutHook          string
	SkipHookPrefixNotConfigured         string
	LcResetTo                           string
	PressEnterToReturn                  string
	LcViewStashOptions                  string
	LcStashAllChanges                   string
	LcStashStagedChanges                string
	LcStashAllChangesKeepIndex          string
	LcStashUnstagedChanges              string
	LcStashIncludeUntrackedChanges      string
	LcStashOptions                      string
	NotARepository                      string
	LcJump                              string
	LcScrollLeftRight                   string
	LcScrollLeft                        string
	LcScrollRight                       string
	DiscardPatch                        string
	DiscardPatchConfirm                 string
	CantPatchWhileRebasingError         string
	LcToggleAddToPatch                  string
	LcToggleAllInPatch                  string
	LcUpdatingPatch                     string
	ViewPatchOptions                    string
	PatchOptionsTitle                   string
	NoPatchError                        string
	LcEnterFile                         string
	ExitCustomPatchBuilder              string
	EnterUpstream                       string
	InvalidUpstream                     string
	ReturnToRemotesList                 string
	LcAddNewRemote                      string
	LcNewRemoteName                     string
	LcNewRemoteUrl                      string
	LcEditRemoteName                    string
	LcEditRemoteUrl                     string
	LcRemoveRemote                      string
	LcRemoveRemotePrompt                string
	DeleteRemoteBranch                  string
	DeleteRemoteBranchMessage           string
	LcSetAsUpstream                     string
	LcSetUpstream                       string
	LcUnsetUpstream                     string
	SetUpstreamTitle                    string
	SetUpstreamMessage                  string
	LcEditRemote                        string
	LcTagCommit                         string
	TagMenuTitle                        string
	TagNameTitle                        string
	TagMessageTitle                     string
	LcLightweightTag                    string
	LcAnnotatedTag                      string
	LcDeleteTag                         string
	DeleteTagTitle                      string
	DeleteTagPrompt                     string
	PushTagTitle                        string
	LcPushTag                           string
	LcCreateTag                         string
	CreateTagTitle                      string
	LcFetchRemote                       string
	FetchingRemoteStatus                string
	LcCheckoutCommit                    string
	SureCheckoutThisCommit              string
	LcGitFlowOptions                    string
	NotAGitFlowBranch                   string
	NewBranchNamePrompt                 string
	IgnoreTracked                       string
	ExcludeTracked                      string
	IgnoreTrackedPrompt                 string
	ExcludeTrackedPrompt                string
	LcViewResetToUpstreamOptions        string
	LcNextScreenMode                    string
	LcPrevScreenMode                    string
	LcStartSearch                       string
	Panel                               string
	Keybindings                         string
	LcRenameBranch                      string
	LcSetUnsetUpstream                  string
	NewGitFlowBranchPrompt              string
	RenameBranchWarning                 string
	LcOpenMenu                          string
	LcResetCherryPick                   string
	LcNextTab                           string
	LcPrevTab                           string
	LcCantUndoWhileRebasing             string
	LcCantRedoWhileRebasing             string
	MustStashWarning                    string
	MustStashTitle                      string
	ConfirmationTitle                   string
	LcPrevPage                          string
	LcNextPage                          string
	LcGotoTop                           string
	LcGotoBottom                        string
	LcFilteringBy                       string
	ResetInParentheses                  string
	LcOpenFilteringMenu                 string
	LcFilterBy                          string
	LcExitFilterMode                    string
	LcFilterPathOption                  string
	EnterFileName                       string
	FilteringMenuTitle                  string
	MustExitFilterModeTitle             string
	MustExitFilterModePrompt            string
	LcDiff                              string
	LcEnterRefToDiff                    string
	LcEnteRefName                       string
	LcExitDiffMode                      string
	DiffingMenuTitle                    string
	LcSwapDiff                          string
	LcOpenDiffingMenu                   string
	LcOpenExtrasMenu                    string
	LcShowingGitDiff                    string
	LcCommitDiff                        string
	LcCopyCommitShaToClipboard          string
	LcCommitSha                         string
	LcCommitURL                         string
	LcCopyCommitMessageToClipboard      string
	LcCommitMessage                     string
	LcCommitAuthor                      string
	LcCopyCommitAttributeToClipboard    string
	LcCopyBranchNameToClipboard         string
	LcCopyFileNameToClipboard           string
	LcCopyCommitFileNameToClipboard     string
	LcCommitPrefixPatternError          string
	LcCopySelectedTexToClipboard        string
	NoFilesStagedTitle                  string
	NoFilesStagedPrompt                 string
	BranchNotFoundTitle                 string
	BranchNotFoundPrompt                string
	LcBranchUnknown                     string
	UnstageLinesTitle                   string
	UnstageLinesPrompt                  string
	LcCreateNewBranchFromCommit         string
	LcBuildingPatch                     string
	LcViewCommits                       string
	MinGitVersionError                  string
	LcRunningCustomCommandStatus        string
	LcSubmoduleStashAndReset            string
	LcAndResetSubmodules                string
	LcEnterSubmodule                    string
	LcCopySubmoduleNameToClipboard      string
	RemoveSubmodule                     string
	LcRemoveSubmodule                   string
	RemoveSubmodulePrompt               string
	LcResettingSubmoduleStatus          string
	LcNewSubmoduleName                  string
	LcNewSubmoduleUrl                   string
	LcNewSubmodulePath                  string
	LcAddSubmodule                      string
	LcAddingSubmoduleStatus             string
	LcUpdateSubmoduleUrl                string
	LcUpdatingSubmoduleUrlStatus        string
	LcEditSubmoduleUrl                  string
	LcInitializingSubmoduleStatus       string
	LcInitSubmodule                     string
	LcSubmoduleUpdate                   string
	LcUpdatingSubmoduleStatus           string
	LcBulkInitSubmodules                string
	LcBulkUpdateSubmodules              string
	LcBulkDeinitSubmodules              string
	LcViewBulkSubmoduleOptions          string
	LcBulkSubmoduleOptions              string
	LcRunningCommand                    string
	SubCommitsTitle                     string
	SubmodulesTitle                     string
	NavigationTitle                     string
	SuggestionsCheatsheetTitle          string
	// Unlike the cheatsheet title above, the real suggestions title has a little message saying press tab to focus
	SuggestionsTitle                    string
	ExtrasTitle                         string
	PushingTagStatus                    string
	PullRequestURLCopiedToClipboard     string
	CommitDiffCopiedToClipboard         string
	CommitSHACopiedToClipboard          string
	CommitURLCopiedToClipboard          string
	CommitMessageCopiedToClipboard      string
	CommitAuthorCopiedToClipboard       string
	PatchCopiedToClipboard              string
	LcCopiedToClipboard                 string
	ErrCannotEditDirectory              string
	ErrStageDirWithInlineMergeConflicts string
	ErrRepositoryMovedOrDeleted         string
	CommandLog                          string
	ToggleShowCommandLog                string
	FocusCommandLog                     string
	CommandLogHeader                    string
	RandomTip                           string
	SelectParentCommitForMerge          string
	ToggleWhitespaceInDiffView          string
	IgnoringWhitespaceInDiffView        string
	ShowingWhitespaceInDiffView         string
	IncreaseContextInDiffView           string
	DecreaseContextInDiffView           string
	CreatePullRequestOptions            string
	LcCreatePullRequestOptions          string
	LcDefaultBranch                     string
	LcSelectBranch                      string
	CreatePullRequest                   string
	SelectConfigFile                    string
	NoConfigFileFoundErr                string
	LcLoadingFileSuggestions            string
	LcLoadingCommits                    string
	MustSpecifyOriginError              string
	GitOutput                           string
	GitCommandFailed                    string
	AbortTitle                          string
	AbortPrompt                         string
	LcOpenLogMenu                       string
	LogMenuTitle                        string
	ToggleShowGitGraphAll               string
	ShowGitGraph                        string
	SortCommits                         string
	CantChangeContextSizeError          string
	LcOpenCommitInBrowser               string
	LcViewBisectOptions                 string
	ConfirmRevertCommit                 string
	RewordInEditorTitle                 string
	RewordInEditorPrompt                string
	CheckoutPrompt                      string
	HardResetAutostashPrompt            string
	UpstreamGone                        string
	NukeDescription                     string
	DiscardStagedChangesDescription     string
	EmptyOutput                         string
	Patch                               string
	CustomPatch                         string
	LcCommitsCopied                     string
	LcCommitCopied                      string
	Actions                             Actions
	Bisect                              Bisect
}

type Bisect struct {
	MarkStart                   string
	MarkSkipCurrent             string
	MarkSkipSelected            string
	ResetTitle                  string
	ResetPrompt                 string
	ResetOption                 string
	BisectMenuTitle             string
	Mark                        string
	Skip                        string
	CompleteTitle               string
	CompletePrompt              string
	CompletePromptIndeterminate string
}

type Actions struct {
	CheckoutCommit                    string
	CheckoutTag                       string
	CheckoutBranch                    string
	ForceCheckoutBranch               string
	DeleteBranch                      string
	Merge                             string
	RebaseBranch                      string
	RenameBranch                      string
	SetUnsetUpstream                  string
	CreateBranch                      string
	FastForwardBranch                 string
	CherryPick                        string
	CheckoutFile                      string
	DiscardOldFileChange              string
	SquashCommitDown                  string
	FixupCommit                       string
	RewordCommit                      string
	DropCommit                        string
	EditCommit                        string
	AmendCommit                       string
	ResetCommitAuthor                 string
	SetCommitAuthor                   string
	RevertCommit                      string
	CreateFixupCommit                 string
	SquashAllAboveFixupCommits        string
	MoveCommitUp                      string
	MoveCommitDown                    string
	CopyCommitMessageToClipboard      string
	CopyCommitDiffToClipboard         string
	CopyCommitSHAToClipboard          string
	CopyCommitURLToClipboard          string
	CopyCommitAuthorToClipboard       string
	CopyCommitAttributeToClipboard    string
	CopyPatchToClipboard              string
	CustomCommand                     string
	DiscardAllChangesInDirectory      string
	DiscardUnstagedChangesInDirectory string
	DiscardAllChangesInFile           string
	DiscardAllUnstagedChangesInFile   string
	StageFile                         string
	StageResolvedFiles                string
	UnstageFile                       string
	UnstageAllFiles                   string
	StageAllFiles                     string
	LcIgnoreExcludeFile               string
	IgnoreFileErr                     string
	ExcludeFile                       string
	ExcludeFileErr                    string
	ExcludeGitIgnoreErr               string
	Commit                            string
	EditFile                          string
	Push                              string
	Pull                              string
	OpenFile                          string
	StashAllChanges                   string
	StashAllChangesKeepIndex          string
	StashStagedChanges                string
	StashUnstagedChanges              string
	StashIncludeUntrackedChanges      string
	GitFlowFinish                     string
	GitFlowStart                      string
	CopyToClipboard                   string
	CopySelectedTextToClipboard       string
	RemovePatchFromCommit             string
	MovePatchToSelectedCommit         string
	MovePatchIntoIndex                string
	MovePatchIntoNewCommit            string
	DeleteRemoteBranch                string
	SetBranchUpstream                 string
	AddRemote                         string
	RemoveRemote                      string
	UpdateRemote                      string
	ApplyPatch                        string
	Stash                             string
	RenameStash                       string
	RemoveSubmodule                   string
	ResetSubmodule                    string
	AddSubmodule                      string
	UpdateSubmoduleUrl                string
	InitialiseSubmodule               string
	BulkInitialiseSubmodules          string
	BulkUpdateSubmodules              string
	BulkDeinitialiseSubmodules        string
	UpdateSubmodule                   string
	CreateLightweightTag              string
	CreateAnnotatedTag                string
	DeleteTag                         string
	PushTag                           string
	NukeWorkingTree                   string
	DiscardUnstagedFileChanges        string
	RemoveUntrackedFiles              string
	RemoveStagedFiles                 string
	SoftReset                         string
	MixedReset                        string
	HardReset                         string
	Undo                              string
	Redo                              string
	CopyPullRequestURL                string
	OpenMergeTool                     string
	OpenCommitInBrowser               string
	OpenPullRequest                   string
	StartBisect                       string
	ResetBisect                       string
	BisectSkip                        string
	BisectMark                        string
}

const englishIntroPopupMessage = `
Thanks for using lazygit! Seriously you rock. Three things to share with you:

 1) If you want to learn about lazygit's features, watch this vid:
      https://youtu.be/CPLdltN7wgE

 2) Be sure to read the latest release notes at:
      https://github.com/jesseduffield/lazygit/releases

 3) If you're using git, that makes you a programmer! With your help we can make
    lazygit better, so consider becoming a contributor and joining the fun at
      https://github.com/jesseduffield/lazygit
    You can also sponsor me and tell me what to work on by clicking the donate
    button at the bottom right.
    Or even just star the repo to share the love!
`

const englishDeprecatedEditConfigWarning = `
### Deprecated config warning ###

The following config settings are deprecated and will be removed in a future
version:
{{configs}}

Please refer to

  https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#configuring-file-editing

for up-to-date information how to configure your editor.

`

// exporting this so we can use it in tests
func EnglishTranslationSet() TranslationSet {
	return TranslationSet{
		NotEnoughSpace:                      "Not enough space to render panels",
		DiffTitle:                           "Diff",
		FilesTitle:                          "Files",
		BranchesTitle:                       "Branches",
		CommitsTitle:                        "Commits",
		StashTitle:                          "Stash",
		SnakeTitle:                          "Snake",
		EasterEgg:                           "easter egg",
		UnstagedChanges:                     `Unstaged Changes`,
		StagedChanges:                       `Staged Changes`,
		MainTitle:                           "Main",
		MergeConfirmTitle:                   "Merge",
		StagingTitle:                        "Main Panel (Staging)",
		MergingTitle:                        "Main Panel (Merging)",
		NormalTitle:                         "Main Panel (Normal)",
		LogTitle:                            "Log",
		CommitSummary:                       "Commit summary",
		CredentialsUsername:                 "Username",
		CredentialsPassword:                 "Password",
		CredentialsPassphrase:               "Enter passphrase for SSH key",
		CredentialsPIN:                      "Enter PIN for SSH key",
		PassUnameWrong:                      "Password, passphrase and/or username wrong",
		CommitChanges:                       "commit changes",
		AmendLastCommit:                     "amend last commit",
		AmendLastCommitTitle:                "Amend Last Commit",
		SureToAmend:                         "Are you sure you want to amend last commit? Afterwards, you can change commit message from the commits panel.",
		NoCommitToAmend:                     "There's no commit to amend.",
		CommitChangesWithEditor:             "commit changes using git editor",
		StatusTitle:                         "Status",
		LcNavigate:                          "navigate",
		LcMenu:                              "menu",
		LcExecute:                           "execute",
		LcToggleStaged:                      "toggle staged",
		LcToggleStagedAll:                   "stage/unstage all",
		LcToggleTreeView:                    "toggle file tree view",
		LcOpenMergeTool:                     "open external merge tool (git mergetool)",
		LcRefresh:                           "refresh",
		LcPush:                              "push",
		LcPull:                              "pull",
		LcScroll:                            "scroll",
		MergeConflictsTitle:                 "Merge Conflicts",
		LcCheckout:                          "checkout",
		LcFileFilter:                        "Filter files (staged/unstaged)",
		FilterStagedFiles:                   "Show only staged files",
		FilterUnstagedFiles:                 "Show only unstaged files",
		ResetCommitFilterState:              "Reset filter",
		NoChangedFiles:                      "No changed files",
		PullWait:                            "Pulling...",
		PushWait:                            "Pushing...",
		FetchWait:                           "Fetching...",
		LcSoftReset:                         "soft reset",
		AlreadyCheckedOutBranch:             "You have already checked out this branch",
		SureForceCheckout:                   "Are you sure you want force checkout? You will lose all local changes",
		ForceCheckoutBranch:                 "Force Checkout Branch",
		BranchName:                          "Branch name",
		NewBranchNameBranchOff:              "New Branch Name (Branch is off of '{{.branchName}}')",
		CantDeleteCheckOutBranch:            "You cannot delete the checked out branch!",
		DeleteBranch:                        "Delete Branch",
		DeleteBranchMessage:                 "Are you sure you want to delete the branch '{{.selectedBranchName}}'?",
		ForceDeleteBranchMessage:            "'{{.selectedBranchName}}' is not fully merged. Are you sure you want to delete it?",
		LcRebaseBranch:                      "rebase checked-out branch onto this branch",
		CantRebaseOntoSelf:                  "You cannot rebase a branch onto itself",
		CantMergeBranchIntoItself:           "You cannot merge a branch into itself",
		LcForceCheckout:                     "force checkout",
		LcCheckoutByName:                    "checkout by name",
		LcNewBranch:                         "new branch",
		LcDeleteBranch:                      "delete branch",
		NoBranchesThisRepo:                  "No branches for this repo",
		CommitMessageConfirm:                "{{.keyBindClose}}: close, {{.keyBindConfirm}}: confirm",
		CommitWithoutMessageErr:             "You cannot commit without a commit message",
		CloseConfirm:                        "{{.keyBindClose}}: close/cancel, {{.keyBindConfirm}}: confirm",
		LcClose:                             "close",
		LcQuit:                              "quit",
		LcSquashDown:                        "squash down",
		LcFixupCommit:                       "fixup commit",
		NoCommitsThisBranch:                 "No commits for this branch",
		UpdateRefHere:                       "Update branch '{{.ref}}' here",
		CannotSquashOrFixupFirstCommit:      "There's no commit below to squash into",
		Fixup:                               "Fixup",
		SureFixupThisCommit:                 "Are you sure you want to 'fixup' this commit? It will be merged into the commit below",
		SureSquashThisCommit:                "Are you sure you want to squash this commit into the commit below?",
		Squash:                              "Squash",
		LcPickCommit:                        "pick commit (when mid-rebase)",
		LcRevertCommit:                      "revert commit",
		LcRewordCommit:                      "reword commit",
		LcDeleteCommit:                      "delete commit",
		LcMoveDownCommit:                    "move commit down one",
		LcMoveUpCommit:                      "move commit up one",
		LcEditCommit:                        "edit commit",
		LcAmendToCommit:                     "amend commit with staged changes",
		LcResetCommitAuthor:                 "reset commit author",
		SetAuthorPromptTitle:                "Set author (must look like 'Name <Email>')",
		SureResetCommitAuthor:               "The author field of this commit will be updated to match the configured user. This also renews the author timestamp. Continue?",
		LcRenameCommitEditor:                "reword commit with editor",
		Error:                               "Error",
		LcSelectHunk:                        "select hunk",
		LcNavigateConflicts:                 "navigate conflicts",
		LcPickHunk:                          "pick hunk",
		LcPickAllHunks:                      "pick all hunks",
		LcUndo:                              "undo",
		LcUndoReflog:                        "undo (via reflog) (experimental)",
		LcRedoReflog:                        "redo (via reflog) (experimental)",
		UndoTooltip:                         "The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		RedoTooltip:                         "The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		DiscardAllTooltip:                   "Discard both staged and unstaged changes in '{{.path}}'.",
		DiscardUnstagedTooltip:              "Discard unstaged changes in '{{.path}}'.",
		LcPop:                               "pop",
		LcDrop:                              "drop",
		LcApply:                             "apply",
		NoStashEntries:                      "No stash entries",
		StashDrop:                           "Stash drop",
		SureDropStashEntry:                  "Are you sure you want to drop this stash entry?",
		StashPop:                            "Stash pop",
		SurePopStashEntry:                   "Are you sure you want to pop this stash entry?",
		StashApply:                          "Stash apply",
		SureApplyStashEntry:                 "Are you sure you want to apply this stash entry?",
		NoTrackedStagedFilesStash:           "You have no tracked/staged files to stash",
		NoFilesToStash:                      "You have no files to stash",
		StashChanges:                        "Stash changes",
		LcRenameStash:                       "rename stash",
		RenameStashPrompt:                   "Rename stash: {{.stashName}}",
		OpenConfig:                          "open config file",
		EditConfig:                          "edit config file",
		ForcePush:                           "Force push",
		ForcePushPrompt:                     "Your branch has diverged from the remote branch. Press 'esc' to cancel, or 'enter' to force push.",
		ForcePushDisabled:                   "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejectedAndForcePushDisabled: "Updates were rejected and you have disabled force pushing",
		LcCheckForUpdate:                    "check for update",
		CheckingForUpdates:                  "Checking for updates...",
		UpdateAvailableTitle:                "Update available!",
		UpdateAvailable:                     "Download and install version {{.newVersion}}?",
		UpdateInProgressWaitingStatus:       "updating",
		UpdateCompletedTitle:                "Update completed!",
		UpdateCompleted:                     "Update has been installed successfully. Restart lazygit for it to take effect.",
		FailedToRetrieveLatestVersionErr:    "Failed to retrieve version information",
		OnLatestVersionErr:                  "You already have the latest version",
		MajorVersionErr:                     "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		CouldNotFindBinaryErr:               "Could not find any binary at {{.url}}",
		UpdateFailedErr:                     "Update failed: {{.errMessage}}",
		ConfirmQuitDuringUpdateTitle:        "Currently Updating",
		ConfirmQuitDuringUpdate:             "An update is in progress. Are you sure you want to quit?",
		MergeToolTitle:                      "Merge tool",
		MergeToolPrompt:                     "Are you sure you want to open `git mergetool`?",
		IntroPopupMessage:                   englishIntroPopupMessage,
		DeprecatedEditConfigWarning:         englishDeprecatedEditConfigWarning,
		GitconfigParseErr:                   `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		LcEditFile:                          `edit file`,
		LcOpenFile:                          `open file`,
		LcIgnoreFile:                        `add to .gitignore`,
		LcExcludeFile:                       `add to .git/info/exclude`,
		LcRefreshFiles:                      `refresh files`,
		LcMergeIntoCurrentBranch:            `merge into currently checked out branch`,
		ConfirmQuit:                         `Are you sure you want to quit?`,
		SwitchRepo:                          `switch to a recent repo`,
		LcAllBranchesLogGraph:               `show all branch logs`,
		UnsupportedGitService:               `Unsupported git service`,
		LcCreatePullRequest:                 `create pull request`,
		LcCopyPullRequestURL:                `copy pull request URL to clipboard`,
		NoBranchOnRemote:                    `This branch doesn't exist on remote. You need to push it to remote first.`,
		LcFetch:                             `fetch`,
		NoAutomaticGitFetchTitle:            `No automatic git fetch`,
		NoAutomaticGitFetchBody:             `Lazygit can't use "git fetch" in a private repo; use 'f' in the files panel to run "git fetch" manually`,
		FileEnter:                           `stage individual hunks/lines for file, or collapse/expand for directory`,
		FileStagingRequirements:             `Can only stage individual lines for tracked files`,
		StageSelection:                      `toggle line staged / unstaged`,
		ResetSelection:                      `delete change (git reset)`,
		ToggleDragSelect:                    `toggle drag select`,
		ToggleSelectHunk:                    `toggle select hunk`,
		ToggleSelectionForPatch:             `add/remove line(s) to patch`,
		EditHunk:                            `edit hunk`,
		ToggleStagingPanel:                  `switch to other panel (staged/unstaged changes)`,
		ReturnToFilesPanel:                  `return to files panel`,
		FastForward:                         `fast-forward this branch from its upstream`,
		Fetching:                            "fetching and fast-forwarding {{.from}} -> {{.to}} ...",
		FoundConflicts:                      "Conflicts! To abort press 'esc', otherwise press 'enter'",
		FoundConflictsTitle:                 "Auto-merge failed",
		PickHunk:                            "pick hunk",
		PickAllHunks:                        "pick all hunks",
		ViewMergeRebaseOptions:              "view merge/rebase options",
		NotMergingOrRebasing:                "You are currently neither rebasing nor merging",
		AlreadyRebasing:                     "Can't perform this action during a rebase",
		RecentRepos:                         "recent repositories",
		MergeOptionsTitle:                   "Merge Options",
		RebaseOptionsTitle:                  "Rebase Options",
		CommitMessageTitle:                  "Commit Summary",
		CommitDescriptionTitle:              "Commit description",
		CommitDescriptionSubTitle:           "Press tab to toggle focus",
		LocalBranchesTitle:                  "Local Branches",
		SearchTitle:                         "Search",
		TagsTitle:                           "Tags",
		MenuTitle:                           "Menu",
		RemotesTitle:                        "Remotes",
		RemoteBranchesTitle:                 "Remote Branches",
		PatchBuildingTitle:                  "Main Panel (Patch Building)",
		InformationTitle:                    "Information",
		SecondaryTitle:                      "Secondary",
		ReflogCommitsTitle:                  "Reflog",
		GlobalTitle:                         "Global Keybindings",
		ConflictsResolved:                   "all merge conflicts resolved. Continue?",
		RebasingTitle:                       "Rebase '{{.checkedOutBranch}}' onto '{{.ref}}'",
		SimpleRebase:                        "simple rebase",
		InteractiveRebase:                   "interactive rebase",
		InteractiveRebaseTooltip:            "Begin an interactive rebase with a break at the start, so you can update the TODO commits before continuing",
		ConfirmMerge:                        "Are you sure you want to merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}'?",
		FwdNoUpstream:                       "Cannot fast-forward a branch with no upstream",
		FwdNoLocalUpstream:                  "Cannot fast-forward a branch whose remote is not registered locally",
		FwdCommitsToPush:                    "Cannot fast-forward a branch with commits to push",
		ErrorOccurred:                       "An error occurred! Please create an issue at",
		NoRoom:                              "Not enough room",
		YouAreHere:                          "YOU ARE HERE",
		YouDied:                             "YOU DIED!",
		LcRewordNotSupported:                "rewording commits while interactively rebasing is not currently supported",
		LcChangingThisActionIsNotAllowed:    "changing this kind of rebase todo entry is not allowed",
		LcCherryPickCopy:                    "copy commit (cherry-pick)",
		LcCherryPickCopyRange:               "copy commit range (cherry-pick)",
		LcPasteCommits:                      "paste commits (cherry-pick)",
		SureCherryPick:                      "Are you sure you want to cherry-pick the copied commits onto this branch?",
		CherryPick:                          "Cherry-Pick",
		Donate:                              "Donate",
		AskQuestion:                         "Ask Question",
		PrevLine:                            "select previous line",
		NextLine:                            "select next line",
		PrevHunk:                            "select previous hunk",
		NextHunk:                            "select next hunk",
		PrevConflict:                        "select previous conflict",
		NextConflict:                        "select next conflict",
		SelectPrevHunk:                      "select previous hunk",
		SelectNextHunk:                      "select next hunk",
		ScrollDown:                          "scroll down",
		ScrollUp:                            "scroll up",
		LcScrollUpMainPanel:                 "scroll up main panel",
		LcScrollDownMainPanel:               "scroll down main panel",
		AmendCommitTitle:                    "Amend Commit",
		AmendCommitPrompt:                   "Are you sure you want to amend this commit with your staged files?",
		DeleteCommitTitle:                   "Delete Commit",
		DeleteCommitPrompt:                  "Are you sure you want to delete this commit?",
		SquashingStatus:                     "squashing",
		FixingStatus:                        "fixing up",
		DeletingStatus:                      "deleting",
		MovingStatus:                        "moving",
		RebasingStatus:                      "rebasing",
		AmendingStatus:                      "amending",
		CherryPickingStatus:                 "cherry-picking",
		UndoingStatus:                       "undoing",
		RedoingStatus:                       "redoing",
		CheckingOutStatus:                   "checking out",
		CommittingStatus:                    "committing",
		CommitFiles:                         "Commit files",
		SubCommitsDynamicTitle:              "Commits (%s)",
		CommitFilesDynamicTitle:             "Diff files (%s)",
		RemoteBranchesDynamicTitle:          "Remote branches (%s)",
		LcViewItemFiles:                     "view selected item's files",
		CommitFilesTitle:                    "Commit Files",
		LcCheckoutCommitFile:                "checkout file",
		LcDiscardOldFileChange:              "discard this commit's changes to this file",
		DiscardFileChangesTitle:             "Discard file changes",
		DiscardFileChangesPrompt:            "Are you sure you want to discard this commit's changes to this file? If this file was created in this commit, it will be deleted",
		DisabledForGPG:                      "Feature not available for users using GPG",
		CreateRepo:                          "Not in a git repository. Create a new git repository? (y/n): ",
		BareRepo:                            "You've attempted to open Lazygit in a bare repo but Lazygit does not yet support bare repos. Open most recent repo? (y/n) ",
		InitialBranch:                       "Branch name? (leave empty for git's default): ",
		NoRecentRepositories:                "Must open lazygit in a git repository. No valid recent repositories. Exiting.",
		IncorrectNotARepository:             "The value of 'notARepository' is incorrect. It should be one of 'prompt', 'create', 'skip', or 'quit'.",
		AutoStashTitle:                      "Autostash?",
		AutoStashPrompt:                     "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		StashPrefix:                         "Auto-stashing changes for ",
		LcViewDiscardOptions:                "view 'discard changes' options",
		LcCancel:                            "cancel",
		LcDiscardAllChanges:                 "discard all changes",
		LcDiscardUnstagedChanges:            "discard unstaged changes",
		LcDiscardAllChangesToAllFiles:       "nuke working tree",
		LcDiscardAnyUnstagedChanges:         "discard unstaged changes",
		LcDiscardUntrackedFiles:             "discard untracked files",
		LcDiscardStagedChanges:              "discard staged changes",
		LcHardReset:                         "hard reset",
		LcViewResetOptions:                  `view reset options`,
		LcCreateFixupCommit:                 `create fixup commit for this commit`,
		LcSquashAboveCommits:                `squash all 'fixup!' commits above selected commit (autosquash)`,
		SquashAboveCommits:                  `Squash all 'fixup!' commits above selected commit (autosquash)`,
		SureSquashAboveCommits:              `Are you sure you want to squash all fixup! commits above {{.commit}}?`,
		CreateFixupCommit:                   `Create fixup commit`,
		SureCreateFixupCommit:               `Are you sure you want to create a fixup! commit for commit {{.commit}}?`,
		LcExecuteCustomCommand:              "execute custom command",
		CustomCommand:                       "Custom Command:",
		LcCommitChangesWithoutHook:          "commit changes without pre-commit hook",
		SkipHookPrefixNotConfigured:         "You have not configured a commit message prefix for skipping hooks. Set `git.skipHookPrefix = 'WIP'` in your config",
		LcResetTo:                           `reset to`,
		PressEnterToReturn:                  "Press enter to return to lazygit",
		LcViewStashOptions:                  "view stash options",
		LcStashAllChanges:                   "stash all changes",
		LcStashStagedChanges:                "stash staged changes",
		LcStashAllChangesKeepIndex:          "stash all changes and keep index",
		LcStashUnstagedChanges:              "stash unstaged changes",
		LcStashIncludeUntrackedChanges:      "stash all changes including untracked files",
		LcStashOptions:                      "Stash options",
		NotARepository:                      "Error: must be run inside a git repository",
		LcJump:                              "jump to panel",
		LcScrollLeftRight:                   "scroll left/right",
		LcScrollLeft:                        "scroll left",
		LcScrollRight:                       "scroll right",
		DiscardPatch:                        "Discard Patch",
		DiscardPatchConfirm:                 "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		CantPatchWhileRebasingError:         "You cannot build a patch or run patch commands while in a merging or rebasing state",
		LcToggleAddToPatch:                  "toggle file included in patch",
		LcToggleAllInPatch:                  "toggle all files included in patch",
		LcUpdatingPatch:                     "updating patch",
		ViewPatchOptions:                    "view custom patch options",
		PatchOptionsTitle:                   "Patch Options",
		NoPatchError:                        "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		LcEnterFile:                         "enter file to add selected lines to the patch (or toggle directory collapsed)",
		ExitCustomPatchBuilder:              `exit custom patch builder`,
		EnterUpstream:                       `Enter upstream as '<remote> <branchname>'`,
		InvalidUpstream:                     "Invalid upstream. Must be in the format '<remote> <branchname>'",
		ReturnToRemotesList:                 `Return to remotes list`,
		LcAddNewRemote:                      `add new remote`,
		LcNewRemoteName:                     `New remote name:`,
		LcNewRemoteUrl:                      `New remote url:`,
		LcEditRemoteName:                    `Enter updated remote name for {{.remoteName}}:`,
		LcEditRemoteUrl:                     `Enter updated remote url for {{.remoteName}}:`,
		LcRemoveRemote:                      `remove remote`,
		LcRemoveRemotePrompt:                "Are you sure you want to remove remote",
		DeleteRemoteBranch:                  "Delete Remote Branch",
		DeleteRemoteBranchMessage:           "Are you sure you want to delete remote branch",
		LcSetAsUpstream:                     "set as upstream of checked-out branch",
		LcSetUpstream:                       "set upstream of selected branch",
		LcUnsetUpstream:                     "unset upstream of selected branch",
		SetUpstreamTitle:                    "Set upstream branch",
		SetUpstreamMessage:                  "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		LcEditRemote:                        "edit remote",
		LcTagCommit:                         "tag commit",
		TagMenuTitle:                        "Create tag",
		TagNameTitle:                        "Tag name:",
		TagMessageTitle:                     "Tag message:",
		LcAnnotatedTag:                      "annotated tag",
		LcLightweightTag:                    "lightweight tag",
		LcDeleteTag:                         "delete tag",
		DeleteTagTitle:                      "Delete tag",
		DeleteTagPrompt:                     "Are you sure you want to delete tag '{{.tagName}}'?",
		PushTagTitle:                        "remote to push tag '{{.tagName}}' to:",
		LcPushTag:                           "push tag",
		LcCreateTag:                         "create tag",
		CreateTagTitle:                      "Tag name:",
		LcFetchRemote:                       "fetch remote",
		FetchingRemoteStatus:                "fetching remote",
		LcCheckoutCommit:                    "checkout commit",
		SureCheckoutThisCommit:              "Are you sure you want to checkout this commit?",
		LcGitFlowOptions:                    "show git-flow options",
		NotAGitFlowBranch:                   "This does not seem to be a git flow branch",
		NewGitFlowBranchPrompt:              "new {{.branchType}} name:",
		IgnoreTracked:                       "Ignore tracked file",
		IgnoreTrackedPrompt:                 "Are you sure you want to ignore a tracked file?",
		ExcludeTracked:                      "Exclude tracked file",
		ExcludeTrackedPrompt:                "Are you sure you want to exclude a tracked file?",
		LcViewResetToUpstreamOptions:        "view upstream reset options",
		LcNextScreenMode:                    "next screen mode (normal/half/fullscreen)",
		LcPrevScreenMode:                    "prev screen mode",
		LcStartSearch:                       "start search",
		Panel:                               "Panel",
		Keybindings:                         "Keybindings",
		LcRenameBranch:                      "rename branch",
		LcSetUnsetUpstream:                  "set/unset upstream",
		NewBranchNamePrompt:                 "Enter new branch name for branch",
		RenameBranchWarning:                 "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		LcOpenMenu:                          "open menu",
		LcResetCherryPick:                   "reset cherry-picked (copied) commits selection",
		LcNextTab:                           "next tab",
		LcPrevTab:                           "previous tab",
		LcCantUndoWhileRebasing:             "Can't undo while rebasing",
		LcCantRedoWhileRebasing:             "Can't redo while rebasing",
		MustStashWarning:                    "Pulling a patch out into the index requires stashing and unstashing your changes. If something goes wrong, you'll be able to access your files from the stash. Continue?",
		MustStashTitle:                      "Must stash",
		ConfirmationTitle:                   "Confirmation Panel",
		LcPrevPage:                          "previous page",
		LcNextPage:                          "next page",
		LcGotoTop:                           "scroll to top",
		LcGotoBottom:                        "scroll to bottom",
		LcFilteringBy:                       "filtering by",
		ResetInParentheses:                  "(reset)",
		LcOpenFilteringMenu:                 "view filter-by-path options",
		LcFilterBy:                          "filter by",
		LcExitFilterMode:                    "stop filtering by path",
		LcFilterPathOption:                  "enter path to filter by",
		EnterFileName:                       "Enter path:",
		FilteringMenuTitle:                  "Filtering",
		MustExitFilterModeTitle:             "Command not available",
		MustExitFilterModePrompt:            "Command not available in filtered mode. Exit filtered mode?",
		LcDiff:                              "diff",
		LcEnterRefToDiff:                    "enter ref to diff",
		LcEnteRefName:                       "enter ref:",
		LcExitDiffMode:                      "exit diff mode",
		DiffingMenuTitle:                    "Diffing",
		LcSwapDiff:                          "reverse diff direction",
		LcOpenDiffingMenu:                   "open diff menu",
		// the actual view is the extras view which I intend to give more tabs in future but for now we'll only mention the command log part
		LcOpenExtrasMenu:                    "open command log menu",
		LcShowingGitDiff:                    "showing output for:",
		LcCommitDiff:                        "commit diff",
		LcCopyCommitShaToClipboard:          "copy commit SHA to clipboard",
		LcCommitSha:                         "commit SHA",
		LcCommitURL:                         "commit URL",
		LcCopyCommitMessageToClipboard:      "copy commit message to clipboard",
		LcCommitMessage:                     "commit message",
		LcCommitAuthor:                      "commit author",
		LcCopyCommitAttributeToClipboard:    "copy commit attribute",
		LcCopyBranchNameToClipboard:         "copy branch name to clipboard",
		LcCopyFileNameToClipboard:           "copy the file name to the clipboard",
		LcCopyCommitFileNameToClipboard:     "copy the committed file name to the clipboard",
		LcCopySelectedTexToClipboard:        "copy the selected text to the clipboard",
		LcCommitPrefixPatternError:          "Error in commitPrefix pattern",
		NoFilesStagedTitle:                  "No files staged",
		NoFilesStagedPrompt:                 "You have not staged any files. Commit all files?",
		BranchNotFoundTitle:                 "Branch not found",
		BranchNotFoundPrompt:                "Branch not found. Create a new branch named",
		LcBranchUnknown:                     "branch unknown",
		UnstageLinesTitle:                   "Unstage lines",
		UnstageLinesPrompt:                  "Are you sure you want to delete the selected lines (git reset)? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipUnstageLineWarning' to true",
		LcCreateNewBranchFromCommit:         "create new branch off of commit",
		LcBuildingPatch:                     "building patch",
		LcViewCommits:                       "view commits",
		MinGitVersionError:                  "Git version must be at least 2.20 (i.e. from 2018 onwards). Please upgrade your git version. Alternatively raise an issue at https://github.com/jesseduffield/lazygit/issues for lazygit to be more backwards compatible.",
		LcRunningCustomCommandStatus:        "running custom command",
		LcSubmoduleStashAndReset:            "stash uncommitted submodule changes and update",
		LcAndResetSubmodules:                "and reset submodules",
		LcEnterSubmodule:                    "enter submodule",
		LcCopySubmoduleNameToClipboard:      "copy submodule name to clipboard",
		RemoveSubmodule:                     "Remove submodule",
		LcRemoveSubmodule:                   "remove submodule",
		RemoveSubmodulePrompt:               "Are you sure you want to remove submodule '%s' and its corresponding directory? This is irreversible.",
		LcResettingSubmoduleStatus:          "resetting submodule",
		LcNewSubmoduleName:                  "new submodule name:",
		LcNewSubmoduleUrl:                   "new submodule URL:",
		LcNewSubmodulePath:                  "new submodule path:",
		LcAddSubmodule:                      "add new submodule",
		LcAddingSubmoduleStatus:             "adding submodule",
		LcUpdateSubmoduleUrl:                "update URL for submodule '%s'",
		LcUpdatingSubmoduleUrlStatus:        "updating URL",
		LcEditSubmoduleUrl:                  "update submodule URL",
		LcInitializingSubmoduleStatus:       "initializing submodule",
		LcInitSubmodule:                     "initialize submodule",
		LcSubmoduleUpdate:                   "update submodule",
		LcUpdatingSubmoduleStatus:           "updating submodule",
		LcBulkInitSubmodules:                "bulk init submodules",
		LcBulkUpdateSubmodules:              "bulk update submodules",
		LcBulkDeinitSubmodules:              "bulk deinit submodules",
		LcViewBulkSubmoduleOptions:          "view bulk submodule options",
		LcBulkSubmoduleOptions:              "bulk submodule options",
		LcRunningCommand:                    "running command",
		SubCommitsTitle:                     "Sub-commits",
		SubmodulesTitle:                     "Submodules",
		NavigationTitle:                     "List Panel Navigation",
		SuggestionsCheatsheetTitle:          "Suggestions",
		SuggestionsTitle:                    "Suggestions (press %s to focus)",
		ExtrasTitle:                         "Command Log",
		PushingTagStatus:                    "pushing tag",
		PullRequestURLCopiedToClipboard:     "Pull request URL copied to clipboard",
		CommitDiffCopiedToClipboard:         "Commit diff copied to clipboard",
		CommitSHACopiedToClipboard:          "Commit SHA copied to clipboard",
		CommitURLCopiedToClipboard:          "Commit URL copied to clipboard",
		CommitMessageCopiedToClipboard:      "Commit message copied to clipboard",
		CommitAuthorCopiedToClipboard:       "Commit author copied to clipboard",
		PatchCopiedToClipboard:              "Patch copied to clipboard",
		LcCopiedToClipboard:                 "copied to clipboard",
		ErrCannotEditDirectory:              "Cannot edit directory: you can only edit individual files",
		ErrStageDirWithInlineMergeConflicts: "Cannot stage/unstage directory containing files with inline merge conflicts. Please fix up the merge conflicts first",
		ErrRepositoryMovedOrDeleted:         "Cannot find repo. It might have been moved or deleted ¯\\_(ツ)_/¯",
		CommandLog:                          "Command Log",
		ToggleShowCommandLog:                "Toggle show/hide command log",
		FocusCommandLog:                     "Focus command log",
		CommandLogHeader:                    "You can hide/focus this panel by pressing '%s'\n",
		RandomTip:                           "Random Tip",
		SelectParentCommitForMerge:          "Select parent commit for merge",
		ToggleWhitespaceInDiffView:          "Toggle whether or not whitespace changes are shown in the diff view",
		IgnoringWhitespaceInDiffView:        "Whitespace will be ignored in the diff view",
		ShowingWhitespaceInDiffView:         "Whitespace will be shown in the diff view",
		IncreaseContextInDiffView:           "Increase the size of the context shown around changes in the diff view",
		DecreaseContextInDiffView:           "Decrease the size of the context shown around changes in the diff view",
		CreatePullRequest:                   "Create pull request",
		CreatePullRequestOptions:            "Create pull request options",
		LcCreatePullRequestOptions:          "create pull request options",
		LcDefaultBranch:                     "default branch",
		LcSelectBranch:                      "select branch",
		SelectConfigFile:                    "Select config file",
		NoConfigFileFoundErr:                "No config file found",
		LcLoadingFileSuggestions:            "loading file suggestions",
		LcLoadingCommits:                    "loading commits",
		MustSpecifyOriginError:              "Must specify a remote if specifying a branch",
		GitOutput:                           "Git output:",
		GitCommandFailed:                    "Git command failed. Check command log for details (open with %s)",
		AbortTitle:                          "Abort %s",
		AbortPrompt:                         "Are you sure you want to abort the current %s?",
		LcOpenLogMenu:                       "open log menu",
		LogMenuTitle:                        "Commit Log Options",
		ToggleShowGitGraphAll:               "toggle show whole git graph (pass the `--all` flag to `git log`)",
		ShowGitGraph:                        "show git graph",
		SortCommits:                         "commit sort order",
		CantChangeContextSizeError:          "Cannot change context while in patch building mode because we were too lazy to support it when releasing the feature. If you really want it, please let us know!",
		LcOpenCommitInBrowser:               "open commit in browser",
		LcViewBisectOptions:                 "view bisect options",
		ConfirmRevertCommit:                 "Are you sure you want to revert {{.selectedCommit}}?",
		RewordInEditorTitle:                 "Reword in editor",
		RewordInEditorPrompt:                "Are you sure you want to reword this commit in your editor?",
		HardResetAutostashPrompt:            "Are you sure you want to hard reset to '%s'? An auto-stash will be performed if necessary.",
		CheckoutPrompt:                      "Are you sure you want to checkout '%s'?",
		UpstreamGone:                        "(upstream gone)",
		NukeDescription:                     "If you want to make all the changes in the worktree go away, this is the way to do it. If there are dirty submodule changes this will stash those changes in the submodule(s).",
		DiscardStagedChangesDescription:     "This will create a new stash entry containing only staged files and then drop it, so that the working tree is left with only unstaged changes",
		EmptyOutput:                         "<empty output>",
		Patch:                               "Patch",
		CustomPatch:                         "Custom patch",
		LcCommitsCopied:                     "commits copied",
		LcCommitCopied:                      "commit copied",
		Actions: Actions{
			// TODO: combine this with the original keybinding descriptions (those are all in lowercase atm)
			CheckoutCommit:                    "Checkout commit",
			CheckoutTag:                       "Checkout tag",
			CheckoutBranch:                    "Checkout branch",
			ForceCheckoutBranch:               "Force checkout branch",
			DeleteBranch:                      "Delete branch",
			Merge:                             "Merge",
			RebaseBranch:                      "Rebase branch",
			RenameBranch:                      "Rename branch",
			SetUnsetUpstream:                  "Set/unset upstream",
			CreateBranch:                      "Create branch",
			CherryPick:                        "(Cherry-pick) Paste commits",
			CheckoutFile:                      "Checkout file",
			DiscardOldFileChange:              "Discard old file change",
			SquashCommitDown:                  "Squash commit down",
			FixupCommit:                       "Fixup commit",
			RewordCommit:                      "Reword commit",
			DropCommit:                        "Drop commit",
			EditCommit:                        "Edit commit",
			AmendCommit:                       "Amend commit",
			ResetCommitAuthor:                 "Reset commit author",
			SetCommitAuthor:                   "Set commit author",
			RevertCommit:                      "Revert commit",
			CreateFixupCommit:                 "Create fixup commit",
			SquashAllAboveFixupCommits:        "Squash all above fixup commits",
			CreateLightweightTag:              "Create lightweight tag",
			CreateAnnotatedTag:                "Create annotated tag",
			CopyCommitMessageToClipboard:      "Copy commit message to clipboard",
			CopyCommitDiffToClipboard:         "Copy commit diff to clipboard",
			CopyCommitSHAToClipboard:          "Copy commit SHA to clipboard",
			CopyCommitURLToClipboard:          "Copy commit URL to clipboard",
			CopyCommitAuthorToClipboard:       "Copy commit author to clipboard",
			CopyCommitAttributeToClipboard:    "Copy to clipboard",
			CopyPatchToClipboard:              "Copy patch to clipboard",
			MoveCommitUp:                      "Move commit up",
			MoveCommitDown:                    "Move commit down",
			CustomCommand:                     "Custom command",
			DiscardAllChangesInDirectory:      "Discard all changes in directory",
			DiscardUnstagedChangesInDirectory: "Discard unstaged changes in directory",
			DiscardAllChangesInFile:           "Discard all changes in file",
			DiscardAllUnstagedChangesInFile:   "Discard all unstaged changes in file",
			StageFile:                         "Stage file",
			StageResolvedFiles:                "Stage files whose merge conflicts were resolved",
			UnstageFile:                       "Unstage file",
			UnstageAllFiles:                   "Unstage all files",
			StageAllFiles:                     "Stage all files",
			LcIgnoreExcludeFile:               "ignore or exclude file",
			IgnoreFileErr:                     "Cannot ignore .gitignore",
			ExcludeFile:                       "Exclude file",
			ExcludeFileErr:                    "Cannot exclude .git/info/exclude",
			ExcludeGitIgnoreErr:               "Cannot exclude .gitignore",
			Commit:                            "Commit",
			EditFile:                          "Edit file",
			Push:                              "Push",
			Pull:                              "Pull",
			OpenFile:                          "Open file",
			StashAllChanges:                   "Stash all changes",
			StashAllChangesKeepIndex:          "Stash all changes and keep index",
			StashStagedChanges:                "Stash staged changes",
			StashUnstagedChanges:              "Stash unstaged changes",
			StashIncludeUntrackedChanges:      "Stash all changes including untracked files",
			GitFlowFinish:                     "Git flow finish",
			GitFlowStart:                      "Git Flow start",
			CopyToClipboard:                   "Copy to clipboard",
			CopySelectedTextToClipboard:       "Copy selected text to clipboard",
			RemovePatchFromCommit:             "Remove patch from commit",
			MovePatchToSelectedCommit:         "Move patch to selected commit",
			MovePatchIntoIndex:                "Move patch into index",
			MovePatchIntoNewCommit:            "Move patch into new commit",
			DeleteRemoteBranch:                "Delete remote branch",
			SetBranchUpstream:                 "Set branch upstream",
			AddRemote:                         "Add remote",
			RemoveRemote:                      "Remove remote",
			UpdateRemote:                      "Update remote",
			ApplyPatch:                        "Apply patch",
			Stash:                             "Stash",
			RenameStash:                       "Rename stash",
			RemoveSubmodule:                   "Remove submodule",
			ResetSubmodule:                    "Reset submodule",
			AddSubmodule:                      "Add submodule",
			UpdateSubmoduleUrl:                "Update submodule URL",
			InitialiseSubmodule:               "Initialise submodule",
			BulkInitialiseSubmodules:          "Bulk initialise submodules",
			BulkUpdateSubmodules:              "Bulk update submodules",
			BulkDeinitialiseSubmodules:        "Bulk deinitialise submodules",
			UpdateSubmodule:                   "Update submodule",
			DeleteTag:                         "Delete tag",
			PushTag:                           "Push tag",
			NukeWorkingTree:                   "Nuke working tree",
			DiscardUnstagedFileChanges:        "Discard unstaged file changes",
			RemoveUntrackedFiles:              "Remove untracked files",
			RemoveStagedFiles:                 "Remove staged files",
			SoftReset:                         "Soft reset",
			MixedReset:                        "Mixed reset",
			HardReset:                         "Hard reset",
			FastForwardBranch:                 "Fast forward branch",
			Undo:                              "Undo",
			Redo:                              "Redo",
			CopyPullRequestURL:                "Copy pull request URL",
			OpenMergeTool:                     "Open merge tool",
			OpenCommitInBrowser:               "Open commit in browser",
			OpenPullRequest:                   "Open pull request in browser",
			StartBisect:                       "Start bisect",
			ResetBisect:                       "Reset bisect",
			BisectSkip:                        "Bisect skip",
			BisectMark:                        "Bisect mark",
		},
		Bisect: Bisect{
			Mark:                        "mark %s as %s",
			MarkStart:                   "mark %s as %s (start bisect)",
			Skip:                        "skip %s",
			ResetTitle:                  "Reset 'git bisect'",
			ResetPrompt:                 "Are you sure you want to reset 'git bisect'?",
			ResetOption:                 "reset bisect",
			BisectMenuTitle:             "Bisect",
			CompleteTitle:               "Bisect complete",
			CompletePrompt:              "Bisect complete! The following commit introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			CompletePromptIndeterminate: "Bisect complete! Some commits were skipped, so any of the following commits may have introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
		},
	}
}
