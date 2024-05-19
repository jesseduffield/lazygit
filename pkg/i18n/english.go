/*

Todo list when making a new translation
- Copy this file and rename it to the language you want to translate to like someLanguage.go
- Change the EnglishTranslationSet() name to the language you want to translate to like SomeLanguageTranslationSet()
- Add an entry of someLanguage in GetTranslationSets()
- Remove this todo and the about section

*/

package i18n

type TranslationSet struct {
	NotEnoughSpace                        string
	DiffTitle                             string
	FilesTitle                            string
	BranchesTitle                         string
	CommitsTitle                          string
	StashTitle                            string
	SnakeTitle                            string
	EasterEgg                             string
	UnstagedChanges                       string
	StagedChanges                         string
	MainTitle                             string
	StagingTitle                          string
	MergingTitle                          string
	MergeConfirmTitle                     string
	NormalTitle                           string
	LogTitle                              string
	CommitSummary                         string
	CredentialsUsername                   string
	CredentialsPassword                   string
	CredentialsPassphrase                 string
	CredentialsPIN                        string
	PassUnameWrong                        string
	Commit                                string
	CommitTooltip                         string
	AmendLastCommit                       string
	AmendLastCommitTitle                  string
	SureToAmend                           string
	NoCommitToAmend                       string
	CommitChangesWithEditor               string
	FindBaseCommitForFixup                string
	FindBaseCommitForFixupTooltip         string
	NoDeletedLinesInDiff                  string
	NoBaseCommitsFound                    string
	MultipleBaseCommitsFoundStaged        string
	MultipleBaseCommitsFoundUnstaged      string
	BaseCommitIsAlreadyOnMainBranch       string
	BaseCommitIsNotInCurrentView          string
	HunksWithOnlyAddedLinesWarning        string
	StatusTitle                           string
	GlobalTitle                           string
	Menu                                  string
	Execute                               string
	Stage                                 string
	StageTooltip                          string
	ToggleStagedAll                       string
	ToggleStagedAllTooltip                string
	ToggleTreeView                        string
	ToggleTreeViewTooltip                 string
	OpenDiffTool                          string
	OpenMergeTool                         string
	OpenMergeToolTooltip                  string
	Refresh                               string
	RefreshTooltip                        string
	Push                                  string
	Pull                                  string
	PushTooltip                           string
	PullTooltip                           string
	Scroll                                string
	FileFilter                            string
	CopyToClipboardMenu                   string
	CopyFileName                          string
	CopyFilePath                          string
	CopyFileDiffTooltip                   string
	CopySelectedDiff                      string
	CopyAllFilesDiff                      string
	NoContentToCopyError                  string
	FileNameCopiedToast                   string
	FilePathCopiedToast                   string
	FileDiffCopiedToast                   string
	AllFilesDiffCopiedToast               string
	FilterStagedFiles                     string
	FilterUnstagedFiles                   string
	ResetFilter                           string
	MergeConflictsTitle                   string
	Checkout                              string
	CheckoutTooltip                       string
	CantCheckoutBranchWhilePulling        string
	TagCheckoutTooltip                    string
	RemoteBranchCheckoutTooltip           string
	CantPullOrPushSameBranchTwice         string
	NoChangedFiles                        string
	SoftReset                             string
	AlreadyCheckedOutBranch               string
	SureForceCheckout                     string
	ForceCheckoutBranch                   string
	BranchName                            string
	NewBranchNameBranchOff                string
	CantDeleteCheckOutBranch              string
	DeleteBranchTitle                     string
	DeleteLocalBranch                     string
	DeleteRemoteBranchOption              string
	DeleteRemoteBranchPrompt              string
	ForceDeleteBranchTitle                string
	ForceDeleteBranchMessage              string
	RebaseBranch                          string
	RebaseBranchTooltip                   string
	CantRebaseOntoSelf                    string
	CantMergeBranchIntoItself             string
	ForceCheckout                         string
	ForceCheckoutTooltip                  string
	CheckoutByName                        string
	CheckoutByNameTooltip                 string
	RemoteBranchCheckoutTitle             string
	CheckoutTypeNewBranch                 string
	CheckoutTypeNewBranchTooltip          string
	CheckoutTypeDetachedHead              string
	CheckoutTypeDetachedHeadTooltip       string
	NewBranch                             string
	NewBranchFromStashTooltip             string
	NoBranchesThisRepo                    string
	CommitWithoutMessageErr               string
	Close                                 string
	CloseCancel                           string
	Confirm                               string
	Quit                                  string
	SquashTooltip                         string
	CannotSquashOrFixupFirstCommit        string
	Fixup                                 string
	FixupTooltip                          string
	SureFixupThisCommit                   string
	SureSquashThisCommit                  string
	Squash                                string
	PickCommitTooltip                     string
	Pick                                  string
	CantPickDisabledReason                string
	Edit                                  string
	RevertCommit                          string
	Revert                                string
	RevertCommitTooltip                   string
	Reword                                string
	CommitRewordTooltip                   string
	DropCommit                            string
	DropCommitTooltip                     string
	MoveDownCommit                        string
	MoveUpCommit                          string
	CannotMoveAnyFurther                  string
	EditCommit                            string
	EditCommitTooltip                     string
	AmendCommitTooltip                    string
	Amend                                 string
	ResetAuthor                           string
	ResetAuthorTooltip                    string
	SetAuthor                             string
	SetAuthorTooltip                      string
	AddCoAuthor                           string
	AmendCommitAttribute                  string
	AmendCommitAttributeTooltip           string
	SetAuthorPromptTitle                  string
	AddCoAuthorPromptTitle                string
	AddCoAuthorTooltip                    string
	SureResetCommitAuthor                 string
	RewordCommitEditor                    string
	NoCommitsThisBranch                   string
	UpdateRefHere                         string
	Error                                 string
	Undo                                  string
	UndoReflog                            string
	RedoReflog                            string
	UndoTooltip                           string
	RedoTooltip                           string
	UndoMergeResolveTooltip               string
	DiscardAllTooltip                     string
	DiscardUnstagedTooltip                string
	Pop                                   string
	StashPopTooltip                       string
	Drop                                  string
	StashDropTooltip                      string
	Apply                                 string
	StashApplyTooltip                     string
	NoStashEntries                        string
	StashDrop                             string
	SureDropStashEntry                    string
	StashPop                              string
	SurePopStashEntry                     string
	StashApply                            string
	SureApplyStashEntry                   string
	NoTrackedStagedFilesStash             string
	NoFilesToStash                        string
	StashChanges                          string
	RenameStash                           string
	RenameStashPrompt                     string
	OpenConfig                            string
	EditConfig                            string
	ForcePush                             string
	ForcePushPrompt                       string
	ForcePushDisabled                     string
	UpdatesRejected                       string
	CheckForUpdate                        string
	CheckingForUpdates                    string
	UpdateAvailableTitle                  string
	UpdateAvailable                       string
	UpdateInProgressWaitingStatus         string
	UpdateCompletedTitle                  string
	UpdateCompleted                       string
	FailedToRetrieveLatestVersionErr      string
	OnLatestVersionErr                    string
	MajorVersionErr                       string
	CouldNotFindBinaryErr                 string
	UpdateFailedErr                       string
	ConfirmQuitDuringUpdateTitle          string
	ConfirmQuitDuringUpdate               string
	MergeToolTitle                        string
	MergeToolPrompt                       string
	IntroPopupMessage                     string
	DeprecatedEditConfigWarning           string
	GitconfigParseErr                     string
	EditFile                              string
	EditFileTooltip                       string
	OpenFile                              string
	OpenFileTooltip                       string
	OpenInEditor                          string
	IgnoreFile                            string
	ExcludeFile                           string
	RefreshFiles                          string
	Merge                                 string
	MergeBranchTooltip                    string
	ConfirmQuit                           string
	SwitchRepo                            string
	AllBranchesLogGraph                   string
	UnsupportedGitService                 string
	CopyPullRequestURL                    string
	NoBranchOnRemote                      string
	Fetch                                 string
	FetchTooltip                          string
	NoAutomaticGitFetchTitle              string
	NoAutomaticGitFetchBody               string
	FileEnter                             string
	FileEnterTooltip                      string
	FileStagingRequirements               string
	StageSelectionTooltip                 string
	DiscardSelection                      string
	DiscardSelectionTooltip               string
	ToggleSelectHunk                      string
	ToggleSelectHunkTooltip               string
	ToggleSelectionForPatch               string
	ToggleSelectionForPatchTooltip        string
	EditHunk                              string
	EditHunkTooltip                       string
	ToggleStagingView                     string
	ToggleStagingViewTooltip              string
	ReturnToFilesPanel                    string
	FastForward                           string
	FastForwardTooltip                    string
	FastForwarding                        string
	FoundConflictsTitle                   string
	ViewConflictsMenuItem                 string
	AbortMenuItem                         string
	PickHunk                              string
	PickAllHunks                          string
	ViewMergeRebaseOptions                string
	ViewMergeRebaseOptionsTooltip         string
	ViewMergeOptions                      string
	ViewRebaseOptions                     string
	NotMergingOrRebasing                  string
	AlreadyRebasing                       string
	RecentRepos                           string
	MergeOptionsTitle                     string
	RebaseOptionsTitle                    string
	CommitSummaryTitle                    string
	CommitDescriptionTitle                string
	CommitDescriptionSubTitle             string
	LocalBranchesTitle                    string
	SearchTitle                           string
	TagsTitle                             string
	MenuTitle                             string
	CommitMenuTitle                       string
	RemotesTitle                          string
	RemoteBranchesTitle                   string
	PatchBuildingTitle                    string
	InformationTitle                      string
	SecondaryTitle                        string
	ReflogCommitsTitle                    string
	ConflictsResolved                     string
	Continue                              string
	RebasingTitle                         string
	RebasingFromBaseCommitTitle           string
	SimpleRebase                          string
	InteractiveRebase                     string
	InteractiveRebaseTooltip              string
	MustSelectTodoCommits                 string
	ConfirmMerge                          string
	FwdNoUpstream                         string
	FwdNoLocalUpstream                    string
	FwdCommitsToPush                      string
	PullRequestNoUpstream                 string
	ErrorOccurred                         string
	NoRoom                                string
	YouAreHere                            string
	YouDied                               string
	RewordNotSupported                    string
	ChangingThisActionIsNotAllowed        string
	CherryPickCopy                        string
	CherryPickCopyTooltip                 string
	CherryPickCopyRange                   string
	CherryPickCopyRangeTooltip            string
	PasteCommits                          string
	SureCherryPick                        string
	CherryPick                            string
	CannotCherryPickNonCommit             string
	CannotCherryPickMergeCommit           string
	Donate                                string
	AskQuestion                           string
	PrevLine                              string
	NextLine                              string
	PrevHunk                              string
	NextHunk                              string
	PrevConflict                          string
	NextConflict                          string
	SelectPrevHunk                        string
	SelectNextHunk                        string
	ScrollDown                            string
	ScrollUp                              string
	ScrollUpMainWindow                    string
	ScrollDownMainWindow                  string
	AmendCommitTitle                      string
	AmendCommitPrompt                     string
	DropCommitTitle                       string
	DropCommitPrompt                      string
	DropUpdateRefPrompt                   string
	PullingStatus                         string
	PushingStatus                         string
	FetchingStatus                        string
	SquashingStatus                       string
	FixingStatus                          string
	DeletingStatus                        string
	DroppingStatus                        string
	MovingStatus                          string
	RebasingStatus                        string
	MergingStatus                         string
	LowercaseRebasingStatus               string
	LowercaseMergingStatus                string
	AmendingStatus                        string
	CherryPickingStatus                   string
	UndoingStatus                         string
	RedoingStatus                         string
	CheckingOutStatus                     string
	CommittingStatus                      string
	RevertingStatus                       string
	CreatingFixupCommitStatus             string
	CommitFiles                           string
	SubCommitsDynamicTitle                string
	CommitFilesDynamicTitle               string
	RemoteBranchesDynamicTitle            string
	ViewItemFiles                         string
	ViewItemFilesTooltip                  string
	CommitFilesTitle                      string
	CheckoutCommitFileTooltip             string
	CanOnlyDiscardFromLocalCommits        string
	Remove                                string
	DiscardOldFileChangeTooltip           string
	DiscardFileChangesTitle               string
	DiscardFileChangesPrompt              string
	DisabledForGPG                        string
	CreateRepo                            string
	BareRepo                              string
	InitialBranch                         string
	NoRecentRepositories                  string
	IncorrectNotARepository               string
	AutoStashTitle                        string
	AutoStashPrompt                       string
	StashPrefix                           string
	Discard                               string
	DiscardChangesTitle                   string
	DiscardFileChangesTooltip             string
	Cancel                                string
	DiscardAllChanges                     string
	DiscardUnstagedChanges                string
	DiscardAllChangesToAllFiles           string
	DiscardAnyUnstagedChanges             string
	DiscardUntrackedFiles                 string
	DiscardStagedChanges                  string
	HardReset                             string
	BranchDeleteTooltip                   string
	TagDeleteTooltip                      string
	Delete                                string
	Reset                                 string
	ResetTooltip                          string
	ViewResetOptions                      string
	FileResetOptionsTooltip               string
	CreateFixupCommit                     string
	CreateFixupCommitTooltip              string
	CreateAmendCommit                     string
	FixupMenu_Fixup                       string
	FixupMenu_FixupTooltip                string
	FixupMenu_AmendWithChanges            string
	FixupMenu_AmendWithChangesTooltip     string
	FixupMenu_AmendWithoutChanges         string
	FixupMenu_AmendWithoutChangesTooltip  string
	SquashAboveCommitsTooltip             string
	SquashCommitsAboveSelectedTooltip     string
	SquashCommitsInCurrentBranchTooltip   string
	SquashAboveCommits                    string
	SquashCommitsInCurrentBranch          string
	SquashCommitsAboveSelectedCommit      string
	CannotSquashCommitsInCurrentBranch    string
	ExecuteCustomCommand                  string
	ExecuteCustomCommandTooltip           string
	CustomCommand                         string
	CommitChangesWithoutHook              string
	SkipHookPrefixNotConfigured           string
	ResetTo                               string
	ResetSoftTooltip                      string
	ResetMixedTooltip                     string
	ResetHardTooltip                      string
	PressEnterToReturn                    string
	ViewStashOptions                      string
	ViewStashOptionsTooltip               string
	Stash                                 string
	StashTooltip                          string
	StashAllChanges                       string
	StashStagedChanges                    string
	StashAllChangesKeepIndex              string
	StashUnstagedChanges                  string
	StashIncludeUntrackedChanges          string
	StashOptions                          string
	NotARepository                        string
	WorkingDirectoryDoesNotExist          string
	Jump                                  string
	ScrollLeftRight                       string
	ScrollLeft                            string
	ScrollRight                           string
	DiscardPatch                          string
	DiscardPatchConfirm                   string
	DiscardPatchSameCommitConfirm         string
	CantPatchWhileRebasingError           string
	ToggleAddToPatch                      string
	ToggleAddToPatchTooltip               string
	ToggleAllInPatch                      string
	ToggleAllInPatchTooltip               string
	UpdatingPatch                         string
	ViewPatchOptions                      string
	PatchOptionsTitle                     string
	NoPatchError                          string
	EmptyPatchError                       string
	EnterCommitFile                       string
	EnterCommitFileTooltip                string
	ExitCustomPatchBuilder                string
	EnterUpstream                         string
	InvalidUpstream                       string
	ReturnToRemotesList                   string
	NewRemote                             string
	NewRemoteName                         string
	NewRemoteUrl                          string
	ViewBranches                          string
	EditRemoteName                        string
	EditRemoteUrl                         string
	RemoveRemote                          string
	RemoveRemoteTooltip                   string
	RemoveRemotePrompt                    string
	DeleteRemoteBranch                    string
	DeleteRemoteBranchMessage             string
	DeleteRemoteBranchTooltip             string
	SetAsUpstream                         string
	SetAsUpstreamTooltip                  string
	SetUpstream                           string
	UnsetUpstream                         string
	ViewDivergenceFromUpstream            string
	DivergenceSectionHeaderLocal          string
	DivergenceSectionHeaderRemote         string
	ViewUpstreamResetOptions              string
	ViewUpstreamResetOptionsTooltip       string
	ViewUpstreamRebaseOptions             string
	ViewUpstreamRebaseOptionsTooltip      string
	UpstreamGenericName                   string
	SetUpstreamTitle                      string
	SetUpstreamMessage                    string
	EditRemoteTooltip                     string
	TagCommit                             string
	TagCommitTooltip                      string
	TagMenuTitle                          string
	TagNameTitle                          string
	TagMessageTitle                       string
	LightweightTag                        string
	AnnotatedTag                          string
	DeleteTagTitle                        string
	DeleteLocalTag                        string
	DeleteRemoteTag                       string
	SelectRemoteTagUpstream               string
	DeleteRemoteTagPrompt                 string
	RemoteTagDeletedMessage               string
	PushTagTitle                          string
	PushTag                               string
	PushTagTooltip                        string
	NewTag                                string
	NewTagTooltip                         string
	CreatingTag                           string
	ForceTag                              string
	ForceTagPrompt                        string
	FetchRemoteTooltip                    string
	FetchingRemoteStatus                  string
	CheckoutCommit                        string
	CheckoutCommitTooltip                 string
	SureCheckoutThisCommit                string
	GitFlowOptions                        string
	NotAGitFlowBranch                     string
	NewBranchNamePrompt                   string
	IgnoreTracked                         string
	ExcludeTracked                        string
	IgnoreTrackedPrompt                   string
	ExcludeTrackedPrompt                  string
	ViewResetToUpstreamOptions            string
	NextScreenMode                        string
	PrevScreenMode                        string
	StartSearch                           string
	StartFilter                           string
	Panel                                 string
	Keybindings                           string
	KeybindingsLegend                     string
	KeybindingsMenuSectionLocal           string
	KeybindingsMenuSectionGlobal          string
	KeybindingsMenuSectionNavigation      string
	RenameBranch                          string
	Upstream                              string
	UpstreamTooltip                       string
	BranchUpstreamOptionsTitle            string
	ViewBranchUpstreamOptions             string
	ViewBranchUpstreamOptionsTooltip      string
	UpstreamNotSetError                   string
	NewGitFlowBranchPrompt                string
	RenameBranchWarning                   string
	OpenKeybindingsMenu                   string
	ResetCherryPick                       string
	NextTab                               string
	PrevTab                               string
	CantUndoWhileRebasing                 string
	CantRedoWhileRebasing                 string
	MustStashWarning                      string
	MustStashTitle                        string
	ConfirmationTitle                     string
	PrevPage                              string
	NextPage                              string
	GotoTop                               string
	GotoBottom                            string
	FilteringBy                           string
	ResetInParentheses                    string
	OpenFilteringMenu                     string
	OpenFilteringMenuTooltip              string
	FilterBy                              string
	ExitFilterMode                        string
	ExitFilterModeAuthor                  string
	FilterPathOption                      string
	FilterAuthorOption                    string
	EnterFileName                         string
	EnterAuthor                           string
	FilteringMenuTitle                    string
	WillCancelExistingFilterTooltip       string
	MustExitFilterModeTitle               string
	MustExitFilterModePrompt              string
	Diff                                  string
	EnterRefToDiff                        string
	EnterRefName                          string
	ExitDiffMode                          string
	DiffingMenuTitle                      string
	SwapDiff                              string
	ViewDiffingOptions                    string
	ViewDiffingOptionsTooltip             string
	OpenCommandLogMenu                    string
	OpenCommandLogMenuTooltip             string
	ShowingGitDiff                        string
	CommitDiff                            string
	CopyCommitHashToClipboard             string
	CommitHash                            string
	CommitURL                             string
	CopyCommitMessageToClipboard          string
	CommitMessage                         string
	CommitSubject                         string
	CommitAuthor                          string
	CopyCommitAttributeToClipboard        string
	CopyCommitAttributeToClipboardTooltip string
	CopyBranchNameToClipboard             string
	CopyPathToClipboard                   string
	CommitPrefixPatternError              string
	CopySelectedTextToClipboard           string
	NoFilesStagedTitle                    string
	NoFilesStagedPrompt                   string
	BranchNotFoundTitle                   string
	BranchNotFoundPrompt                  string
	BranchUnknown                         string
	DiscardChangeTitle                    string
	DiscardChangePrompt                   string
	CreateNewBranchFromCommit             string
	BuildingPatch                         string
	ViewCommits                           string
	MinGitVersionError                    string
	RunningCustomCommandStatus            string
	SubmoduleStashAndReset                string
	AndResetSubmodules                    string
	EnterSubmoduleTooltip                 string
	Enter                                 string
	CopySubmoduleNameToClipboard          string
	RemoveSubmodule                       string
	RemoveSubmoduleTooltip                string
	RemoveSubmodulePrompt                 string
	ResettingSubmoduleStatus              string
	NewSubmoduleName                      string
	NewSubmoduleUrl                       string
	NewSubmodulePath                      string
	NewSubmodule                          string
	AddingSubmoduleStatus                 string
	UpdateSubmoduleUrl                    string
	UpdatingSubmoduleUrlStatus            string
	EditSubmoduleUrl                      string
	InitializingSubmoduleStatus           string
	InitSubmoduleTooltip                  string
	Update                                string
	Initialize                            string
	SubmoduleUpdateTooltip                string
	UpdatingSubmoduleStatus               string
	BulkInitSubmodules                    string
	BulkUpdateSubmodules                  string
	BulkDeinitSubmodules                  string
	ViewBulkSubmoduleOptions              string
	BulkSubmoduleOptions                  string
	RunningCommand                        string
	SubCommitsTitle                       string
	SubmodulesTitle                       string
	NavigationTitle                       string
	SuggestionsCheatsheetTitle            string
	// Unlike the cheatsheet title above, the real suggestions title has a little message saying press tab to focus
	SuggestionsTitle                     string
	SuggestionsSubtitle                  string
	ExtrasTitle                          string
	PushingTagStatus                     string
	PullRequestURLCopiedToClipboard      string
	CommitDiffCopiedToClipboard          string
	CommitURLCopiedToClipboard           string
	CommitMessageCopiedToClipboard       string
	CommitSubjectCopiedToClipboard       string
	CommitAuthorCopiedToClipboard        string
	PatchCopiedToClipboard               string
	CopiedToClipboard                    string
	ErrCannotEditDirectory               string
	ErrStageDirWithInlineMergeConflicts  string
	ErrRepositoryMovedOrDeleted          string
	ErrWorktreeMovedOrRemoved            string
	CommandLog                           string
	ToggleShowCommandLog                 string
	FocusCommandLog                      string
	CommandLogHeader                     string
	RandomTip                            string
	SelectParentCommitForMerge           string
	ToggleWhitespaceInDiffView           string
	ToggleWhitespaceInDiffViewTooltip    string
	IgnoreWhitespaceDiffViewSubTitle     string
	IgnoreWhitespaceNotSupportedHere     string
	IncreaseContextInDiffView            string
	IncreaseContextInDiffViewTooltip     string
	DecreaseContextInDiffView            string
	DecreaseContextInDiffViewTooltip     string
	DiffContextSizeChanged               string
	CreatePullRequestOptions             string
	DefaultBranch                        string
	SelectBranch                         string
	CreatePullRequest                    string
	SelectConfigFile                     string
	NoConfigFileFoundErr                 string
	LoadingFileSuggestions               string
	LoadingCommits                       string
	MustSpecifyOriginError               string
	GitOutput                            string
	GitCommandFailed                     string
	AbortTitle                           string
	AbortPrompt                          string
	OpenLogMenu                          string
	OpenLogMenuTooltip                   string
	LogMenuTitle                         string
	ToggleShowGitGraphAll                string
	ShowGitGraph                         string
	SortOrder                            string
	SortAlphabetical                     string
	SortByDate                           string
	SortByRecency                        string
	SortBasedOnReflog                    string
	SortCommits                          string
	CantChangeContextSizeError           string
	OpenCommitInBrowser                  string
	ViewBisectOptions                    string
	ConfirmRevertCommit                  string
	RewordInEditorTitle                  string
	RewordInEditorPrompt                 string
	CheckoutPrompt                       string
	HardResetAutostashPrompt             string
	UpstreamGone                         string
	NukeDescription                      string
	DiscardStagedChangesDescription      string
	EmptyOutput                          string
	Patch                                string
	CustomPatch                          string
	CommitsCopied                        string
	CommitCopied                         string
	ResetPatch                           string
	ResetPatchTooltip                    string
	ApplyPatch                           string
	ApplyPatchTooltip                    string
	ApplyPatchInReverse                  string
	ApplyPatchInReverseTooltip           string
	RemovePatchFromOriginalCommit        string
	RemovePatchFromOriginalCommitTooltip string
	MovePatchOutIntoIndex                string
	MovePatchOutIntoIndexTooltip         string
	MovePatchIntoNewCommit               string
	MovePatchIntoNewCommitTooltip        string
	MovePatchToSelectedCommit            string
	MovePatchToSelectedCommitTooltip     string
	CopyPatchToClipboard                 string
	NoMatchesFor                         string
	MatchesFor                           string
	SearchKeybindings                    string
	SearchPrefix                         string
	FilterPrefix                         string
	ExitSearchMode                       string
	ExitTextFilterMode                   string
	Switch                               string
	SwitchToWorktree                     string
	SwitchToWorktreeTooltip              string
	AlreadyCheckedOutByWorktree          string
	BranchCheckedOutByWorktree           string
	DetachWorktreeTooltip                string
	Switching                            string
	RemoveWorktree                       string
	RemoveWorktreeTitle                  string
	DetachWorktree                       string
	DetachingWorktree                    string
	WorktreesTitle                       string
	WorktreeTitle                        string
	RemoveWorktreePrompt                 string
	ForceRemoveWorktreePrompt            string
	RemovingWorktree                     string
	AddingWorktree                       string
	CantDeleteCurrentWorktree            string
	AlreadyInWorktree                    string
	CantDeleteMainWorktree               string
	NoWorktreesThisRepo                  string
	MissingWorktree                      string
	MainWorktree                         string
	NewWorktree                          string
	NewWorktreePath                      string
	NewWorktreeBase                      string
	RemoveWorktreeTooltip                string
	BranchNameCannotBeBlank              string
	NewBranchName                        string
	NewBranchNameLeaveBlank              string
	ViewWorktreeOptions                  string
	CreateWorktreeFrom                   string
	CreateWorktreeFromDetached           string
	LcWorktree                           string
	ChangingDirectoryTo                  string
	Name                                 string
	Branch                               string
	Path                                 string
	MarkedBaseCommitStatus               string
	MarkAsBaseCommit                     string
	MarkAsBaseCommitTooltip              string
	MarkedCommitMarker                   string
	PleaseGoToURL                        string
	NoCopiedCommits                      string
	DisabledMenuItemPrefix               string
	QuickStartInteractiveRebase          string
	QuickStartInteractiveRebaseTooltip   string
	CannotQuickStartInteractiveRebase    string
	ToggleRangeSelect                    string
	RangeSelectUp                        string
	RangeSelectDown                      string
	RangeSelectNotSupported              string
	NoItemSelected                       string
	SelectedItemIsNotABranch             string
	SelectedItemDoesNotHaveFiles         string
	RangeSelectNotSupportedForSubmodules string
	OldCherryPickKeyWarning              string
	CommandDoesNotSupportOpeningInEditor string
	Actions                              Actions
	Bisect                               Bisect
	Log                                  Log
	BreakingChangesTitle                 string
	BreakingChangesMessage               string
	BreakingChangesByVersion             map[string]string
}

type Bisect struct {
	MarkStart                   string
	MarkSkipCurrent             string
	MarkSkipSelected            string
	ResetTitle                  string
	ResetPrompt                 string
	ResetOption                 string
	ChooseTerms                 string
	OldTermPrompt               string
	NewTermPrompt               string
	BisectMenuTitle             string
	Mark                        string
	SkipCurrent                 string
	SkipSelected                string
	CompleteTitle               string
	CompletePrompt              string
	CompletePromptIndeterminate string
	Bisecting                   string
}

type Log struct {
	EditRebase               string
	MoveCommitUp             string
	MoveCommitDown           string
	CherryPickCommits        string
	HandleUndo               string
	HandleMidRebaseCommand   string
	RemoveFile               string
	CopyToClipboard          string
	Remove                   string
	CreateFileWithContent    string
	AppendingLineToFile      string
	EditRebaseFromBaseCommit string
}

type Actions struct {
	CheckoutCommit                    string
	CheckoutTag                       string
	CheckoutBranch                    string
	ForceCheckoutBranch               string
	DeleteLocalBranch                 string
	DeleteBranch                      string
	Merge                             string
	RebaseBranch                      string
	RenameBranch                      string
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
	AddCommitCoAuthor                 string
	RevertCommit                      string
	CreateFixupCommit                 string
	SquashAllAboveFixupCommits        string
	MoveCommitUp                      string
	MoveCommitDown                    string
	CopyCommitMessageToClipboard      string
	CopyCommitSubjectToClipboard      string
	CopyCommitDiffToClipboard         string
	CopyCommitHashToClipboard         string
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
	IgnoreExcludeFile                 string
	IgnoreFileErr                     string
	ExcludeFile                       string
	ExcludeGitIgnoreErr               string
	Commit                            string
	EditFile                          string
	Push                              string
	Pull                              string
	OpenFile                          string
	OpenFileTooltip                   string
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
	DeleteLocalTag                    string
	DeleteRemoteTag                   string
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
	OpenDiffTool                      string
	OpenMergeTool                     string
	OpenCommitInBrowser               string
	OpenPullRequest                   string
	StartBisect                       string
	ResetBisect                       string
	BisectSkip                        string
	BisectMark                        string
	RemoveWorktree                    string
	AddWorktree                       string
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
		NotEnoughSpace:                       "Not enough space to render panels",
		DiffTitle:                            "Diff",
		FilesTitle:                           "Files",
		BranchesTitle:                        "Branches",
		CommitsTitle:                         "Commits",
		StashTitle:                           "Stash",
		SnakeTitle:                           "Snake",
		EasterEgg:                            "Easter egg",
		UnstagedChanges:                      "Unstaged changes",
		StagedChanges:                        "Staged changes",
		MainTitle:                            "Main",
		MergeConfirmTitle:                    "Merge",
		StagingTitle:                         "Main panel (staging)",
		MergingTitle:                         "Main panel (merging)",
		NormalTitle:                          "Main panel (normal)",
		LogTitle:                             "Log",
		CommitSummary:                        "Commit summary",
		CredentialsUsername:                  "Username",
		CredentialsPassword:                  "Password",
		CredentialsPassphrase:                "Enter passphrase for SSH key",
		CredentialsPIN:                       "Enter PIN for SSH key",
		PassUnameWrong:                       "Password, passphrase and/or username wrong",
		Commit:                               "Commit",
		CommitTooltip:                        "Commit staged changes.",
		AmendLastCommit:                      "Amend last commit",
		AmendLastCommitTitle:                 "Amend last commit",
		SureToAmend:                          "Are you sure you want to amend last commit? Afterwards, you can change the commit message from the commits panel.",
		NoCommitToAmend:                      "There's no commit to amend.",
		CommitChangesWithEditor:              "Commit changes using git editor",
		FindBaseCommitForFixup:               "Find base commit for fixup",
		FindBaseCommitForFixupTooltip:        "Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md>",
		NoDeletedLinesInDiff:                 "No deleted lines in diff",
		NoBaseCommitsFound:                   "No base commits found",
		MultipleBaseCommitsFoundStaged:       "Multiple base commits found. (Try staging fewer changes at once)",
		MultipleBaseCommitsFoundUnstaged:     "Multiple base commits found. (Try staging some of the changes)",
		BaseCommitIsAlreadyOnMainBranch:      "The base commit for this change is already on the main branch",
		BaseCommitIsNotInCurrentView:         "Base commit is not in current view",
		HunksWithOnlyAddedLinesWarning:       "There are ranges of only added lines in the diff; be careful to check that these belong in the found base commit.\n\nProceed?",
		StatusTitle:                          "Status",
		Menu:                                 "Menu",
		Execute:                              "Execute",
		Stage:                                "Stage",
		StageTooltip:                         "Toggle staged for selected file.",
		ToggleStagedAll:                      "Stage all",
		ToggleStagedAllTooltip:               "Toggle staged/unstaged for all files in working tree.",
		ToggleTreeView:                       "Toggle file tree view",
		ToggleTreeViewTooltip:                "Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.",
		OpenDiffTool:                         "Open external diff tool (git difftool)",
		OpenMergeTool:                        "Open external merge tool",
		OpenMergeToolTooltip:                 "Run `git mergetool`.",
		Refresh:                              "Refresh",
		RefreshTooltip:                       "Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`.",
		Push:                                 "Push",
		PushTooltip:                          "Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch.",
		Pull:                                 "Pull",
		PullTooltip:                          "Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch.",
		Scroll:                               "Scroll",
		MergeConflictsTitle:                  "Merge conflicts",
		Checkout:                             "Checkout",
		CheckoutTooltip:                      "Checkout selected item.",
		CantCheckoutBranchWhilePulling:       "You cannot checkout another branch while pulling the current branch",
		TagCheckoutTooltip:                   "Checkout the selected tag tag as a detached HEAD.",
		RemoteBranchCheckoutTooltip:          "Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head.",
		CantPullOrPushSameBranchTwice:        "You cannot push or pull a branch while it is already being pushed or pulled",
		FileFilter:                           "Filter files by status",
		CopyToClipboardMenu:                  "Copy to clipboard",
		CopyFileName:                         "File name",
		CopyFilePath:                         "Path",
		CopyFileDiffTooltip:                  "If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.",
		CopySelectedDiff:                     "Diff of selected file",
		CopyAllFilesDiff:                     "Diff of all files",
		NoContentToCopyError:                 "Nothing to copy",
		FileNameCopiedToast:                  "File name copied to clipboard",
		FilePathCopiedToast:                  "File path copied to clipboard",
		FileDiffCopiedToast:                  "File diff copied to clipboard",
		AllFilesDiffCopiedToast:              "All files diff copied to clipboard",
		FilterStagedFiles:                    "Show only staged files",
		FilterUnstagedFiles:                  "Show only unstaged files",
		ResetFilter:                          "Reset filter",
		NoChangedFiles:                       "No changed files",
		SoftReset:                            "Soft reset",
		AlreadyCheckedOutBranch:              "You have already checked out this branch",
		SureForceCheckout:                    "Are you sure you want force checkout? You will lose all local changes",
		ForceCheckoutBranch:                  "Force checkout branch",
		BranchName:                           "Branch name",
		NewBranchNameBranchOff:               "New branch name (branch is off of '{{.branchName}}')",
		CantDeleteCheckOutBranch:             "You cannot delete the checked out branch!",
		DeleteBranchTitle:                    "Delete branch '{{.selectedBranchName}}'?",
		DeleteLocalBranch:                    "Delete local branch",
		DeleteRemoteBranchOption:             "Delete remote branch",
		DeleteRemoteBranchPrompt:             "Are you sure you want to delete the remote branch '{{.selectedBranchName}}' from '{{.upstream}}'?",
		ForceDeleteBranchTitle:               "Force delete branch",
		ForceDeleteBranchMessage:             "'{{.selectedBranchName}}' is not fully merged. Are you sure you want to delete it?",
		RebaseBranch:                         "Rebase",
		RebaseBranchTooltip:                  "Rebase the checked-out branch onto the selected branch.",
		CantRebaseOntoSelf:                   "You cannot rebase a branch onto itself",
		CantMergeBranchIntoItself:            "You cannot merge a branch into itself",
		ForceCheckout:                        "Force checkout",
		ForceCheckoutTooltip:                 "Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch.",
		CheckoutByName:                       "Checkout by name",
		CheckoutByNameTooltip:                "Checkout by name. In the input box you can enter '-' to switch to the last branch.",
		RemoteBranchCheckoutTitle:            "Checkout {{.branchName}}",
		CheckoutTypeNewBranch:                "New local branch",
		CheckoutTypeNewBranchTooltip:         "Checkout the remote branch as a local branch, tracking the remote branch.",
		CheckoutTypeDetachedHead:             "Detached head",
		CheckoutTypeDetachedHeadTooltip:      "Checkout the remote branch as a detached head, which can be useful if you just want to test the branch but not work on it yourself. You can still create a local branch from it later.",
		NewBranch:                            "New branch",
		NewBranchFromStashTooltip:            "Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit.",
		NoBranchesThisRepo:                   "No branches for this repo",
		CommitWithoutMessageErr:              "You cannot commit without a commit message",
		Close:                                "Close",
		CloseCancel:                          "Close/Cancel",
		Confirm:                              "Confirm",
		Quit:                                 "Quit",
		SquashTooltip:                        "Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it.",
		NoCommitsThisBranch:                  "No commits for this branch",
		UpdateRefHere:                        "Update branch '{{.ref}}' here",
		CannotSquashOrFixupFirstCommit:       "There's no commit below to squash into",
		Fixup:                                "Fixup",
		SureFixupThisCommit:                  "Are you sure you want to 'fixup' the selected commit(s) into the commit below?",
		SureSquashThisCommit:                 "Are you sure you want to squash the selected commit(s) into the commit below?",
		Squash:                               "Squash",
		PickCommitTooltip:                    "Mark the selected commit to be picked (when mid-rebase). This means that the commit will be retained upon continuing the rebase.",
		Pick:                                 "Pick",
		CantPickDisabledReason:               "Cannot pick a commit when not mid-rebase",
		Edit:                                 "Edit",
		RevertCommit:                         "Revert commit",
		Revert:                               "Revert",
		RevertCommitTooltip:                  "Create a revert commit for the selected commit, which applies the selected commit's changes in reverse.",
		Reword:                               "Reword",
		CommitRewordTooltip:                  "Reword the selected commit's message.",
		DropCommit:                           "Drop",
		DropCommitTooltip:                    "Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts.",
		MoveDownCommit:                       "Move commit down one",
		MoveUpCommit:                         "Move commit up one",
		CannotMoveAnyFurther:                 "Cannot move any further",
		EditCommit:                           "Edit (start interactive rebase)",
		EditCommitTooltip:                    "Edit the selected commit. Use this to start an interactive rebase from the selected commit. When already mid-rebase, this will mark the selected commit for editing, which means that upon continuing the rebase, the rebase will pause at the selected commit to allow you to make changes.",
		AmendCommitTooltip:                   "Amend commit with staged changes. If the selected commit is the HEAD commit, this will perform `git commit --amend`. Otherwise the commit will be amended via a rebase.",
		Amend:                                "Amend",
		ResetAuthor:                          "Reset author",
		ResetAuthorTooltip:                   "Reset the commit's author to the currently configured user. This will also renew the author timestamp",
		SetAuthor:                            "Set author",
		SetAuthorTooltip:                     "Set the author based on a prompt",
		AddCoAuthor:                          "Add co-author",
		AmendCommitAttribute:                 "Amend commit attribute",
		AmendCommitAttributeTooltip:          "Set/Reset commit author or set co-author.",
		SetAuthorPromptTitle:                 "Set author (must look like 'Name <Email>')",
		AddCoAuthorPromptTitle:               "Add co-author (must look like 'Name <Email>')",
		AddCoAuthorTooltip:                   "Add co-author using the Github/Gitlab metadata Co-authored-by.",
		SureResetCommitAuthor:                "The author field of this commit will be updated to match the configured user. This also renews the author timestamp. Continue?",
		RewordCommitEditor:                   "Reword with editor",
		Error:                                "Error",
		PickHunk:                             "Pick hunk",
		PickAllHunks:                         "Pick all hunks",
		Undo:                                 "Undo",
		UndoReflog:                           "Undo",
		RedoReflog:                           "Redo",
		UndoTooltip:                          "The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		RedoTooltip:                          "The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		UndoMergeResolveTooltip:              "Undo last merge conflict resolution.",
		DiscardAllTooltip:                    "Discard both staged and unstaged changes in '{{.path}}'.",
		DiscardUnstagedTooltip:               "Discard unstaged changes in '{{.path}}'.",
		Pop:                                  "Pop",
		StashPopTooltip:                      "Apply the stash entry to your working directory and remove the stash entry.",
		Drop:                                 "Drop",
		StashDropTooltip:                     "Remove the stash entry from the stash list.",
		Apply:                                "Apply",
		StashApplyTooltip:                    "Apply the stash entry to your working directory.",
		NoStashEntries:                       "No stash entries",
		StashDrop:                            "Stash drop",
		SureDropStashEntry:                   "Are you sure you want to drop this stash entry?",
		StashPop:                             "Stash pop",
		SurePopStashEntry:                    "Are you sure you want to pop this stash entry?",
		StashApply:                           "Stash apply",
		SureApplyStashEntry:                  "Are you sure you want to apply this stash entry?",
		NoTrackedStagedFilesStash:            "You have no tracked/staged files to stash",
		NoFilesToStash:                       "You have no files to stash",
		StashChanges:                         "Stash changes",
		RenameStash:                          "Rename stash",
		RenameStashPrompt:                    "Rename stash: {{.stashName}}",
		OpenConfig:                           "Open config file",
		EditConfig:                           "Edit config file",
		ForcePush:                            "Force push",
		ForcePushPrompt:                      "Your branch has diverged from the remote branch. Press {{.cancelKey}} to cancel, or {{.confirmKey}} to force push.",
		ForcePushDisabled:                    "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejected:                      "Updates were rejected. Please fetch and examine the remote changes before pushing again.",
		CheckForUpdate:                       "Check for update",
		CheckingForUpdates:                   "Checking for updates...",
		UpdateAvailableTitle:                 "Update available!",
		UpdateAvailable:                      "Download and install version {{.newVersion}}?",
		UpdateInProgressWaitingStatus:        "Updating",
		UpdateCompletedTitle:                 "Update completed!",
		UpdateCompleted:                      "Update has been installed successfully. Restart lazygit for it to take effect.",
		FailedToRetrieveLatestVersionErr:     "Failed to retrieve version information",
		OnLatestVersionErr:                   "You already have the latest version",
		MajorVersionErr:                      "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		CouldNotFindBinaryErr:                "Could not find any binary at {{.url}}",
		UpdateFailedErr:                      "Update failed: {{.errMessage}}",
		ConfirmQuitDuringUpdateTitle:         "Currently updating",
		ConfirmQuitDuringUpdate:              "An update is in progress. Are you sure you want to quit?",
		MergeToolTitle:                       "Merge tool",
		MergeToolPrompt:                      "Are you sure you want to open `git mergetool`?",
		IntroPopupMessage:                    englishIntroPopupMessage,
		DeprecatedEditConfigWarning:          englishDeprecatedEditConfigWarning,
		GitconfigParseErr:                    `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		EditFile:                             `Edit file`,
		EditFileTooltip:                      "Open file in external editor.",
		OpenFile:                             `Open file`,
		OpenFileTooltip:                      "Open file in default application.",
		OpenInEditor:                         "Open in editor",
		IgnoreFile:                           `Add to .gitignore`,
		ExcludeFile:                          `Add to .git/info/exclude`,
		RefreshFiles:                         `Refresh files`,
		Merge:                                `Merge`,
		MergeBranchTooltip:                   "Merge selected branch into currently checked out branch.",
		ConfirmQuit:                          `Are you sure you want to quit?`,
		SwitchRepo:                           `Switch to a recent repo`,
		AllBranchesLogGraph:                  `Show all branch logs`,
		UnsupportedGitService:                `Unsupported git service`,
		CreatePullRequest:                    `Create pull request`,
		CopyPullRequestURL:                   `Copy pull request URL to clipboard`,
		NoBranchOnRemote:                     `This branch doesn't exist on remote. You need to push it to remote first.`,
		Fetch:                                `Fetch`,
		FetchTooltip:                         "Fetch changes from remote.",
		NoAutomaticGitFetchTitle:             `No automatic git fetch`,
		NoAutomaticGitFetchBody:              `Lazygit can't use "git fetch" in a private repo; use 'f' in the files panel to run "git fetch" manually`,
		FileEnter:                            `Stage lines / Collapse directory`,
		FileEnterTooltip:                     "If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it.",
		FileStagingRequirements:              `Can only stage individual lines for tracked files`,
		StageSelectionTooltip:                `Toggle selection staged / unstaged.`,
		DiscardSelection:                     `Discard`,
		DiscardSelectionTooltip:              "When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change.",
		ToggleRangeSelect:                    "Toggle range select",
		ToggleSelectHunk:                     "Select hunk",
		ToggleSelectHunkTooltip:              "Toggle hunk selection mode.",
		ToggleSelectionForPatch:              `Toggle lines in patch`,
		EditHunk:                             `Edit hunk`,
		EditHunkTooltip:                      "Edit selected hunk in external editor.",
		ToggleStagingView:                    "Switch view",
		ToggleStagingViewTooltip:             "Switch to other view (staged/unstaged changes).",
		ReturnToFilesPanel:                   `Return to files panel`,
		FastForward:                          `Fast-forward`,
		FastForwardTooltip:                   "Fast-forward selected branch from its upstream.",
		FastForwarding:                       "Fast-forwarding",
		FoundConflictsTitle:                  "Conflicts!",
		ViewConflictsMenuItem:                "View conflicts",
		AbortMenuItem:                        "Abort the %s",
		ViewMergeRebaseOptions:               "View merge/rebase options",
		ViewMergeRebaseOptionsTooltip:        "View options to abort/continue/skip the current merge/rebase.",
		ViewMergeOptions:                     "View merge options",
		ViewRebaseOptions:                    "View rebase options",
		NotMergingOrRebasing:                 "You are currently neither rebasing nor merging",
		AlreadyRebasing:                      "Can't perform this action during a rebase",
		RecentRepos:                          "Recent repositories",
		MergeOptionsTitle:                    "Merge options",
		RebaseOptionsTitle:                   "Rebase options",
		CommitSummaryTitle:                   "Commit summary",
		CommitDescriptionTitle:               "Commit description",
		CommitDescriptionSubTitle:            "Press {{.togglePanelKeyBinding}} to toggle focus, {{.commitMenuKeybinding}} to open menu",
		LocalBranchesTitle:                   "Local branches",
		SearchTitle:                          "Search",
		TagsTitle:                            "Tags",
		MenuTitle:                            "Menu",
		CommitMenuTitle:                      "Commit Menu",
		RemotesTitle:                         "Remotes",
		RemoteBranchesTitle:                  "Remote branches",
		PatchBuildingTitle:                   "Main panel (patch building)",
		InformationTitle:                     "Information",
		SecondaryTitle:                       "Secondary",
		ReflogCommitsTitle:                   "Reflog",
		GlobalTitle:                          "Global keybindings",
		ConflictsResolved:                    "All merge conflicts resolved. Continue?",
		Continue:                             "Continue",
		Keybindings:                          "Keybindings",
		KeybindingsMenuSectionLocal:          "Local",
		KeybindingsMenuSectionGlobal:         "Global",
		KeybindingsMenuSectionNavigation:     "Navigation",
		RebasingTitle:                        "Rebase '{{.checkedOutBranch}}' onto '{{.ref}}'",
		RebasingFromBaseCommitTitle:          "Rebase '{{.checkedOutBranch}}' from marked base onto '{{.ref}}'",
		SimpleRebase:                         "Simple rebase",
		InteractiveRebase:                    "Interactive rebase",
		InteractiveRebaseTooltip:             "Begin an interactive rebase with a break at the start, so you can update the TODO commits before continuing.",
		MustSelectTodoCommits:                "When rebasing, this action only works on a selection of TODO commits.",
		ConfirmMerge:                         "Are you sure you want to merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}'?",
		FwdNoUpstream:                        "Cannot fast-forward a branch with no upstream",
		FwdNoLocalUpstream:                   "Cannot fast-forward a branch whose remote is not registered locally",
		FwdCommitsToPush:                     "Cannot fast-forward a branch with commits to push",
		PullRequestNoUpstream:                "Cannot open a pull request for a branch with no upstream",
		ErrorOccurred:                        "An error occurred! Please create an issue at",
		NoRoom:                               "Not enough room",
		YouAreHere:                           "YOU ARE HERE",
		YouDied:                              "YOU DIED!",
		RewordNotSupported:                   "Rewording commits while interactively rebasing is not currently supported",
		ChangingThisActionIsNotAllowed:       "Changing this kind of rebase todo entry is not allowed",
		CherryPickCopy:                       "Copy (cherry-pick)",
		CherryPickCopyTooltip:                "Mark commit as copied. Then, within the local commits view, you can press `{{.paste}}` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `{{.escape}}` to cancel the selection.",
		CherryPickCopyRangeTooltip:           "Mark commits as copied from the last copied commit to the selected commit.",
		PasteCommits:                         "Paste (cherry-pick)",
		SureCherryPick:                       "Are you sure you want to cherry-pick the copied commits onto this branch?",
		CherryPick:                           "Cherry-pick",
		CannotCherryPickNonCommit:            "Cannot cherry-pick this kind of todo item",
		CannotCherryPickMergeCommit:          "Cherry-picking merge commits is not supported",
		Donate:                               "Donate",
		AskQuestion:                          "Ask Question",
		PrevLine:                             "Select previous line",
		NextLine:                             "Select next line",
		PrevHunk:                             "Go to previous hunk",
		NextHunk:                             "Go to next hunk",
		PrevConflict:                         "Previous conflict",
		NextConflict:                         "Next conflict",
		SelectPrevHunk:                       "Previous hunk",
		SelectNextHunk:                       "Next hunk",
		ScrollDown:                           "Scroll down",
		ScrollUp:                             "Scroll up",
		ScrollUpMainWindow:                   "Scroll up main window",
		ScrollDownMainWindow:                 "Scroll down main window",
		AmendCommitTitle:                     "Amend commit",
		AmendCommitPrompt:                    "Are you sure you want to amend this commit with your staged files?",
		DropCommitTitle:                      "Drop commit",
		DropCommitPrompt:                     "Are you sure you want to drop the selected commit(s)?",
		DropUpdateRefPrompt:                  "Are you sure you want to delete the selected update-ref todo(s)? This is irreversible except by aborting the rebase.",
		PullingStatus:                        "Pulling",
		PushingStatus:                        "Pushing",
		FetchingStatus:                       "Fetching",
		SquashingStatus:                      "Squashing",
		FixingStatus:                         "Fixing up",
		DeletingStatus:                       "Deleting",
		DroppingStatus:                       "Dropping",
		MovingStatus:                         "Moving",
		RebasingStatus:                       "Rebasing",
		MergingStatus:                        "Merging",
		LowercaseRebasingStatus:              "rebasing", // lowercase because it shows up in parentheses
		LowercaseMergingStatus:               "merging",  // lowercase because it shows up in parentheses
		AmendingStatus:                       "Amending",
		CherryPickingStatus:                  "Cherry-picking",
		UndoingStatus:                        "Undoing",
		RedoingStatus:                        "Redoing",
		CheckingOutStatus:                    "Checking out",
		CommittingStatus:                     "Committing",
		RevertingStatus:                      "Reverting",
		CreatingFixupCommitStatus:            "Creating fixup commit",
		CommitFiles:                          "Commit files",
		SubCommitsDynamicTitle:               "Commits (%s)",
		CommitFilesDynamicTitle:              "Diff files (%s)",
		RemoteBranchesDynamicTitle:           "Remote branches (%s)",
		ViewItemFiles:                        "View files",
		ViewItemFilesTooltip:                 "View the files modified by the selected item.",
		CommitFilesTitle:                     "Commit files",
		CheckoutCommitFileTooltip:            "Checkout file. This replaces the file in your working tree with the version from the selected commit.",
		CanOnlyDiscardFromLocalCommits:       "Changes can only be discarded from local commits",
		Remove:                               "Remove",
		DiscardOldFileChangeTooltip:          "Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file.",
		DiscardFileChangesTitle:              "Discard file changes",
		DiscardFileChangesPrompt:             "Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.",
		DisabledForGPG:                       "Feature not available for users using GPG",
		CreateRepo:                           "Not in a git repository. Create a new git repository? (y/n): ",
		BareRepo:                             "You've attempted to open Lazygit in a bare repo but Lazygit does not yet support bare repos. Open most recent repo? (y/n) ",
		InitialBranch:                        "Branch name? (leave empty for git's default): ",
		NoRecentRepositories:                 "Must open lazygit in a git repository. No valid recent repositories. Exiting.",
		IncorrectNotARepository:              "The value of 'notARepository' is incorrect. It should be one of 'prompt', 'create', 'skip', or 'quit'.",
		AutoStashTitle:                       "Autostash?",
		AutoStashPrompt:                      "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		StashPrefix:                          "Auto-stashing changes for ",
		Discard:                              "Discard",
		DiscardFileChangesTooltip:            "View options for discarding changes to the selected file.",
		DiscardChangesTitle:                  "Discard changes",
		Cancel:                               "Cancel",
		DiscardAllChanges:                    "Discard all changes",
		DiscardUnstagedChanges:               "Discard unstaged changes",
		DiscardAllChangesToAllFiles:          "Nuke working tree",
		DiscardAnyUnstagedChanges:            "Discard unstaged changes",
		DiscardUntrackedFiles:                "Discard untracked files",
		DiscardStagedChanges:                 "Discard staged changes",
		HardReset:                            "Hard reset",
		BranchDeleteTooltip:                  "View delete options for local/remote branch.",
		TagDeleteTooltip:                     "View delete options for local/remote tag.",
		Delete:                               "Delete",
		Reset:                                "Reset",
		ResetTooltip:                         "View reset options (soft/mixed/hard) for resetting onto selected item.",
		ResetSoftTooltip:                     "Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes.",
		ResetMixedTooltip:                    "Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes.",
		ResetHardTooltip:                     "Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree.",
		ViewResetOptions:                     `Reset`,
		FileResetOptionsTooltip:              "View reset options for working tree (e.g. nuking the working tree).",
		FixupTooltip:                         "Meld the selected commit into the commit below it. Similar to fixup, but the selected commit's message will be discarded.",
		CreateFixupCommit:                    "Create fixup commit",
		CreateFixupCommitTooltip:             "Create 'fixup!' commit for the selected commit. Later on, you can press `{{.squashAbove}}` on this same commit to apply all above fixup commits.",
		CreateAmendCommit:                    `Create "amend!" commit`,
		FixupMenu_Fixup:                      "fixup! commit",
		FixupMenu_FixupTooltip:               "Lets you fixup another commit and keep the original commit's message.",
		FixupMenu_AmendWithChanges:           "amend! commit with changes",
		FixupMenu_AmendWithChangesTooltip:    "Lets you fixup another commit and also change its commit message.",
		FixupMenu_AmendWithoutChanges:        "amend! commit without changes (pure reword)",
		FixupMenu_AmendWithoutChangesTooltip: "Lets you change the commit message of another commit without changing its content.",
		SquashAboveCommits:                   "Apply fixup commits",
		SquashAboveCommitsTooltip:            `Squash all 'fixup!' commits, either above the selected commit, or all in current branch (autosquash).`,
		SquashCommitsAboveSelectedTooltip:    `Squash all 'fixup!' commits above the selected commit (autosquash).`,
		SquashCommitsInCurrentBranchTooltip:  `Squash all 'fixup!' commits in the current branch (autosquash).`,
		SquashCommitsInCurrentBranch:         "In current branch",
		SquashCommitsAboveSelectedCommit:     "Above the selected commit",
		CannotSquashCommitsInCurrentBranch:   "Cannot squash commits in current branch: the HEAD commit is a merge commit or is present on the main branch.",
		ExecuteCustomCommand:                 "Execute custom command",
		ExecuteCustomCommandTooltip:          "Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands.",
		CustomCommand:                        "Custom command:",
		CommitChangesWithoutHook:             "Commit changes without pre-commit hook",
		SkipHookPrefixNotConfigured:          "You have not configured a commit message prefix for skipping hooks. Set `git.skipHookPrefix = 'WIP'` in your config",
		ResetTo:                              `Reset to`,
		PressEnterToReturn:                   "Press enter to return to lazygit",
		ViewStashOptions:                     "View stash options",
		ViewStashOptionsTooltip:              "View stash options (e.g. stash all, stash staged, stash unstaged).",
		Stash:                                "Stash",
		StashTooltip:                         "Stash all changes. For other variations of stashing, use the view stash options keybinding.",
		StashAllChanges:                      "Stash all changes",
		StashStagedChanges:                   "Stash staged changes",
		StashAllChangesKeepIndex:             "Stash all changes and keep index",
		StashUnstagedChanges:                 "Stash unstaged changes",
		StashIncludeUntrackedChanges:         "Stash all changes including untracked files",
		StashOptions:                         "Stash options",
		NotARepository:                       "Error: must be run inside a git repository",
		WorkingDirectoryDoesNotExist:         "Error: the current working directory does not exist",
		Jump:                                 "Jump to panel",
		ScrollLeftRight:                      "Scroll left/right",
		ScrollLeft:                           "Scroll left",
		ScrollRight:                          "Scroll right",
		DiscardPatch:                         "Discard patch",
		DiscardPatchConfirm:                  "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		DiscardPatchSameCommitConfirm:        "You currently have changes added to a patch for this commit. Discard current patch?",
		CantPatchWhileRebasingError:          "You cannot build a patch or run patch commands while in a merging or rebasing state",
		ToggleAddToPatch:                     "Toggle file included in patch",
		ToggleAddToPatchTooltip:              "Toggle whether the file is included in the custom patch. See {{.doc}}.",
		ToggleAllInPatch:                     "Toggle all files",
		ToggleAllInPatchTooltip:              "Add/remove all commit's files to custom patch. See {{.doc}}.",
		UpdatingPatch:                        "Updating patch",
		ViewPatchOptions:                     "View custom patch options",
		PatchOptionsTitle:                    "Patch options",
		NoPatchError:                         "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		EmptyPatchError:                      "Patch is still empty. Add some files or lines to your patch first.",
		EnterCommitFile:                      "Enter file / Toggle directory collapsed",
		EnterCommitFileTooltip:               "If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory.",
		ExitCustomPatchBuilder:               `Exit custom patch builder`,
		EnterUpstream:                        `Enter upstream as '<remote> <branchname>'`,
		InvalidUpstream:                      "Invalid upstream. Must be in the format '<remote> <branchname>'",
		ReturnToRemotesList:                  `Return to remotes list`,
		NewRemote:                            `New remote`,
		NewRemoteName:                        `New remote name:`,
		NewRemoteUrl:                         `New remote url:`,
		ViewBranches:                         "View branches",
		EditRemoteName:                       `Enter updated remote name for {{.remoteName}}:`,
		EditRemoteUrl:                        `Enter updated remote url for {{.remoteName}}:`,
		RemoveRemote:                         `Remove remote`,
		RemoveRemoteTooltip:                  `Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected.`,
		RemoveRemotePrompt:                   "Are you sure you want to remove remote",
		DeleteRemoteBranch:                   "Delete remote branch",
		DeleteRemoteBranchMessage:            "Are you sure you want to delete remote branch",
		DeleteRemoteBranchTooltip:            "Delete the remote branch from the remote.",
		SetAsUpstream:                        "Set as upstream",
		SetAsUpstreamTooltip:                 "Set the selected remote branch as the upstream of the checked-out branch.",
		SetUpstream:                          "Set upstream of selected branch",
		UnsetUpstream:                        "Unset upstream of selected branch",
		ViewDivergenceFromUpstream:           "View divergence from upstream",
		DivergenceSectionHeaderLocal:         "Local",
		DivergenceSectionHeaderRemote:        "Remote",
		ViewUpstreamResetOptions:             "Reset checked-out branch onto {{.upstream}}",
		ViewUpstreamResetOptionsTooltip:      "View options for resetting the checked-out branch onto {{upstream}}. Note: this will not reset the selected branch onto the upstream, it will reset the checked-out branch onto the upstream.",
		ViewUpstreamRebaseOptions:            "Rebase checked-out branch onto {{.upstream}}",
		ViewUpstreamRebaseOptionsTooltip:     "View options for rebasing the checked-out branch onto {{upstream}}. Note: this will not rebase the selected branch onto the upstream, it will rebased the checked-out branch onto the upstream.",
		UpstreamGenericName:                  "upstream of selected branch",
		SetUpstreamTitle:                     "Set upstream branch",
		SetUpstreamMessage:                   "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		EditRemoteTooltip:                    "Edit the selected remote's name or URL.",
		TagCommit:                            "Tag commit",
		TagCommitTooltip:                     "Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description.",
		TagMenuTitle:                         "Create tag",
		TagNameTitle:                         "Tag name",
		TagMessageTitle:                      "Tag description",
		AnnotatedTag:                         "Annotated tag",
		LightweightTag:                       "Lightweight tag",
		DeleteTagTitle:                       "Delete tag '{{.tagName}}'?",
		DeleteLocalTag:                       "Delete local tag",
		DeleteRemoteTag:                      "Delete remote tag",
		RemoteTagDeletedMessage:              "Remote tag deleted",
		SelectRemoteTagUpstream:              "Remote from which to remove tag '{{.tagName}}':",
		DeleteRemoteTagPrompt:                "Are you sure you want to delete the remote tag '{{.tagName}}' from '{{.upstream}}'?",
		PushTagTitle:                         "Remote to push tag '{{.tagName}}' to:",
		// Using 'push tag' rather than just 'push' to disambiguate from a global push
		PushTag:                "Push tag",
		PushTagTooltip:         "Push the selected tag to a remote. You'll be prompted to select a remote.",
		NewTag:                 "New tag",
		NewTagTooltip:          "Create new tag from current commit. You'll be prompted to enter a tag name and optional description.",
		CreatingTag:            "Creating tag",
		ForceTag:               "Force Tag",
		ForceTagPrompt:         "The tag '{{.tagName}}' exists already. Press {{.cancelKey}} to cancel, or {{.confirmKey}} to overwrite.",
		FetchRemoteTooltip:     "Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches.",
		FetchingRemoteStatus:   "Fetching remote",
		CheckoutCommit:         "Checkout commit",
		CheckoutCommitTooltip:  "Checkout the selected commit as a detached HEAD.",
		SureCheckoutThisCommit: "Are you sure you want to checkout this commit?",
		GitFlowOptions:         "Show git-flow options",
		NotAGitFlowBranch:      "This does not seem to be a git flow branch",
		NewGitFlowBranchPrompt: "New {{.branchType}} name:",

		IgnoreTracked:                    "Ignore tracked file",
		IgnoreTrackedPrompt:              "Are you sure you want to ignore a tracked file?",
		ExcludeTracked:                   "Exclude tracked file",
		ExcludeTrackedPrompt:             "Are you sure you want to exclude a tracked file?",
		ViewResetToUpstreamOptions:       "View upstream reset options",
		NextScreenMode:                   "Next screen mode (normal/half/fullscreen)",
		PrevScreenMode:                   "Prev screen mode",
		StartSearch:                      "Search the current view by text",
		StartFilter:                      "Filter the current view by text",
		Panel:                            "Panel",
		KeybindingsLegend:                "Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b",
		RenameBranch:                     "Rename branch",
		BranchUpstreamOptionsTitle:       "Upstream options",
		ViewBranchUpstreamOptions:        "View upstream options",
		ViewBranchUpstreamOptionsTooltip: "View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream.",
		UpstreamNotSetError:              "The selected branch has no upstream (or the upstream is not stored locally)",
		Upstream:                         "Upstream",
		UpstreamTooltip:                  "View upstream options for selected branch e.g. setting/unsetting the upstream and resetting to the upstream.",
		NewBranchNamePrompt:              "Enter new branch name for branch",
		RenameBranchWarning:              "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		OpenKeybindingsMenu:              "Open keybindings menu",
		ResetCherryPick:                  "Reset copied (cherry-picked) commits selection",
		NextTab:                          "Next tab",
		PrevTab:                          "Previous tab",
		CantUndoWhileRebasing:            "Can't undo while rebasing",
		CantRedoWhileRebasing:            "Can't redo while rebasing",
		MustStashWarning:                 "Pulling a patch out into the index requires stashing and unstashing your changes. If something goes wrong, you'll be able to access your files from the stash. Continue?",
		MustStashTitle:                   "Must stash",
		ConfirmationTitle:                "Confirmation panel",
		PrevPage:                         "Previous page",
		NextPage:                         "Next page",
		GotoTop:                          "Scroll to top",
		GotoBottom:                       "Scroll to bottom",
		FilteringBy:                      "Filtering by",
		ResetInParentheses:               "(Reset)",
		OpenFilteringMenu:                "View filter options",
		OpenFilteringMenuTooltip:         "View options for filtering the commit log, so that only commits matching the filter are shown.",
		FilterBy:                         "Filter by",
		ExitFilterMode:                   "Stop filtering",
		FilterPathOption:                 "Enter path to filter by",
		FilterAuthorOption:               "Enter author to filter by",
		EnterFileName:                    "Enter path:",
		EnterAuthor:                      "Enter author:",
		FilteringMenuTitle:               "Filtering",
		WillCancelExistingFilterTooltip:  "Note: this will cancel the existing filter",
		MustExitFilterModeTitle:          "Command not available",
		MustExitFilterModePrompt:         "Command not available in filter-by-path mode. Exit filter-by-path mode?",
		Diff:                             "Diff",
		EnterRefToDiff:                   "Enter ref to diff",
		EnterRefName:                     "Enter ref:",
		ExitDiffMode:                     "Exit diff mode",
		DiffingMenuTitle:                 "Diffing",
		SwapDiff:                         "Reverse diff direction",
		ViewDiffingOptions:               "View diffing options",
		ViewDiffingOptionsTooltip:        "View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction.",
		// the actual view is the extras view which I intend to give more tabs in future but for now we'll only mention the command log part
		OpenCommandLogMenu:                    "View command log options",
		OpenCommandLogMenuTooltip:             "View options for the command log e.g. show/hide the command log and focus the command log.",
		ShowingGitDiff:                        "Showing output for:",
		CommitDiff:                            "Commit diff",
		CopyCommitHashToClipboard:             "Copy commit hash to clipboard",
		CommitHash:                            "Commit hash",
		CommitURL:                             "Commit URL",
		CopyCommitMessageToClipboard:          "Copy commit message to clipboard",
		CommitMessage:                         "Commit message",
		CommitSubject:                         "Commit subject",
		CommitAuthor:                          "Commit author",
		CopyCommitAttributeToClipboard:        "Copy commit attribute to clipboard",
		CopyCommitAttributeToClipboardTooltip: "Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author).",
		CopyBranchNameToClipboard:             "Copy branch name to clipboard",
		CopyPathToClipboard:                   "Copy path to clipboard",
		CopySelectedTextToClipboard:           "Copy selected text to clipboard",
		CommitPrefixPatternError:              "Error in commitPrefix pattern",
		NoFilesStagedTitle:                    "No files staged",
		NoFilesStagedPrompt:                   "You have not staged any files. Commit all files?",
		BranchNotFoundTitle:                   "Branch not found",
		BranchNotFoundPrompt:                  "Branch not found. Create a new branch named",
		BranchUnknown:                         "Branch unknown",
		DiscardChangeTitle:                    "Discard change",
		DiscardChangePrompt:                   "Are you sure you want to discard this change (git reset)? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipDiscardChangeWarning' to true",
		CreateNewBranchFromCommit:             "Create new branch off of commit",
		BuildingPatch:                         "Building patch",
		ViewCommits:                           "View commits",
		MinGitVersionError:                    "Git version must be at least 2.20 (i.e. from 2018 onwards). Please upgrade your git version. Alternatively raise an issue at https://github.com/jesseduffield/lazygit/issues for lazygit to be more backwards compatible.",
		RunningCustomCommandStatus:            "Running custom command",
		SubmoduleStashAndReset:                "Stash uncommitted submodule changes and update",
		AndResetSubmodules:                    "And reset submodules",
		Enter:                                 "Enter",
		EnterSubmoduleTooltip:                 "Enter submodule. After entering the submodule, you can press `{{.escape}}` to escape back to the parent repo.",
		CopySubmoduleNameToClipboard:          "Copy submodule name to clipboard",
		RemoveSubmodule:                       "Remove submodule",
		RemoveSubmodulePrompt:                 "Are you sure you want to remove submodule '%s' and its corresponding directory? This is irreversible.",
		RemoveSubmoduleTooltip:                "Remove the selected submodule and its corresponding directory.",
		ResettingSubmoduleStatus:              "Resetting submodule",
		NewSubmoduleName:                      "New submodule name:",
		NewSubmoduleUrl:                       "New submodule URL:",
		NewSubmodulePath:                      "New submodule path:",
		NewSubmodule:                          "New submodule",
		AddingSubmoduleStatus:                 "Adding submodule",
		UpdateSubmoduleUrl:                    "Update URL for submodule '%s'",
		UpdatingSubmoduleUrlStatus:            "Updating URL",
		EditSubmoduleUrl:                      "Update submodule URL",
		InitializingSubmoduleStatus:           "Initializing submodule",
		InitSubmoduleTooltip:                  "Initialize the selected submodule to prepare for fetching. You probably want to follow this up by invoking the 'update' action to fetch the submodule.",
		Update:                                "Update",
		Initialize:                            "Initialize",
		SubmoduleUpdateTooltip:                "Update selected submodule.",
		UpdatingSubmoduleStatus:               "Updating submodule",
		BulkInitSubmodules:                    "Bulk init submodules",
		BulkUpdateSubmodules:                  "Bulk update submodules",
		BulkDeinitSubmodules:                  "Bulk deinit submodules",
		ViewBulkSubmoduleOptions:              "View bulk submodule options",
		BulkSubmoduleOptions:                  "Bulk submodule options",
		RunningCommand:                        "Running command",
		SubCommitsTitle:                       "Sub-commits",
		SubmodulesTitle:                       "Submodules",
		NavigationTitle:                       "List panel navigation",
		SuggestionsCheatsheetTitle:            "Suggestions",
		SuggestionsTitle:                      "Suggestions (press %s to focus)",
		SuggestionsSubtitle:                   "(press %s to delete, %s to edit)",
		ExtrasTitle:                           "Command log",
		PushingTagStatus:                      "Pushing tag",
		PullRequestURLCopiedToClipboard:       "Pull request URL copied to clipboard",
		CommitDiffCopiedToClipboard:           "Commit diff copied to clipboard",
		CommitURLCopiedToClipboard:            "Commit URL copied to clipboard",
		CommitMessageCopiedToClipboard:        "Commit message copied to clipboard",
		CommitSubjectCopiedToClipboard:        "Commit subject copied to clipboard",
		CommitAuthorCopiedToClipboard:         "Commit author copied to clipboard",
		PatchCopiedToClipboard:                "Patch copied to clipboard",
		CopiedToClipboard:                     "copied to clipboard",
		ErrCannotEditDirectory:                "Cannot edit directories: you can only edit individual files",
		ErrStageDirWithInlineMergeConflicts:   "Cannot stage/unstage directory containing files with inline merge conflicts. Please fix up the merge conflicts first",
		ErrRepositoryMovedOrDeleted:           "Cannot find repo. It might have been moved or deleted ¯\\_(ツ)_/¯",
		CommandLog:                            "Command log",
		ErrWorktreeMovedOrRemoved:             "Cannot find worktree. It might have been moved or removed ¯\\_(ツ)_/¯",
		ToggleShowCommandLog:                  "Toggle show/hide command log",
		FocusCommandLog:                       "Focus command log",
		CommandLogHeader:                      "You can hide/focus this panel by pressing '%s'\n",
		RandomTip:                             "Random tip",
		SelectParentCommitForMerge:            "Select parent commit for merge",
		ToggleWhitespaceInDiffView:            "Toggle whitespace",
		ToggleWhitespaceInDiffViewTooltip:     "Toggle whether or not whitespace changes are shown in the diff view.",
		IgnoreWhitespaceDiffViewSubTitle:      "(ignoring whitespace)",
		IgnoreWhitespaceNotSupportedHere:      "Ignoring whitespace is not supported in this view",
		IncreaseContextInDiffView:             "Increase diff context size",
		IncreaseContextInDiffViewTooltip:      "Increase the amount of the context shown around changes in the diff view.",
		DecreaseContextInDiffView:             "Decrease diff context size",
		DecreaseContextInDiffViewTooltip:      "Decrease the amount of the context shown around changes in the diff view.",
		DiffContextSizeChanged:                "Changed diff context size to %d",
		CreatePullRequestOptions:              "View create pull request options",
		DefaultBranch:                         "Default branch",
		SelectBranch:                          "Select branch",
		SelectConfigFile:                      "Select config file",
		NoConfigFileFoundErr:                  "No config file found",
		LoadingFileSuggestions:                "Loading file suggestions",
		LoadingCommits:                        "Loading commits",
		MustSpecifyOriginError:                "Must specify a remote if specifying a branch",
		GitOutput:                             "Git output:",
		GitCommandFailed:                      "Git command failed. Check command log for details (open with %s)",
		AbortTitle:                            "Abort %s",
		AbortPrompt:                           "Are you sure you want to abort the current %s?",
		OpenLogMenu:                           "View log options",
		OpenLogMenuTooltip:                    "View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph.",
		LogMenuTitle:                          "Commit Log Options",
		ToggleShowGitGraphAll:                 "Toggle show whole git graph (pass the `--all` flag to `git log`)",
		ShowGitGraph:                          "Show git graph",
		SortOrder:                             "Sort order",
		SortAlphabetical:                      "Alphabetical",
		SortByDate:                            "Date",
		SortByRecency:                         "Recency",
		SortBasedOnReflog:                     "(based on reflog)",
		SortCommits:                           "Commit sort order",
		CantChangeContextSizeError:            "Cannot change context while in patch building mode because we were too lazy to support it when releasing the feature. If you really want it, please let us know!",
		OpenCommitInBrowser:                   "Open commit in browser",
		ViewBisectOptions:                     "View bisect options",
		ConfirmRevertCommit:                   "Are you sure you want to revert {{.selectedCommit}}?",
		RewordInEditorTitle:                   "Reword in editor",
		RewordInEditorPrompt:                  "Are you sure you want to reword this commit in your editor?",
		HardResetAutostashPrompt:              "Are you sure you want to hard reset to '%s'? An auto-stash will be performed if necessary.",
		CheckoutPrompt:                        "Are you sure you want to checkout '%s'?",
		UpstreamGone:                          "(upstream gone)",
		NukeDescription:                       "If you want to make all the changes in the worktree go away, this is the way to do it. If there are dirty submodule changes this will stash those changes in the submodule(s).",
		DiscardStagedChangesDescription:       "This will create a new stash entry containing only staged files and then drop it, so that the working tree is left with only unstaged changes",
		EmptyOutput:                           "<Empty output>",
		Patch:                                 "Patch",
		CustomPatch:                           "Custom patch",
		CommitsCopied:                         "commits copied", // lowercase because it's used in a sentence
		CommitCopied:                          "commit copied",  // lowercase because it's used in a sentence
		ResetPatch:                            "Reset patch",
		ResetPatchTooltip:                     "Clear the current patch.",
		ApplyPatch:                            "Apply patch",
		ApplyPatchTooltip:                     "Apply the current patch to the working tree.",
		ApplyPatchInReverse:                   "Apply patch in reverse",
		ApplyPatchInReverseTooltip:            "Apply the current patch in reverse to the working tree.",
		RemovePatchFromOriginalCommit:         "Remove patch from original commit (%s)",
		RemovePatchFromOriginalCommitTooltip:  "Remove the current patch from its commit. This is achieved by starting an interactive rebase at the commit, applying the patch in reverse, and then continuing the rebase. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchOutIntoIndex:                 "Move patch out into index",
		MovePatchOutIntoIndexTooltip:          "Move the patch out of its commit and into the index. This is achieved by starting an interactive rebase at the commit, applying the patch in reverse, continuing the rebase to completion, and then applying the patch to the index. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchIntoNewCommit:                "Move patch into new commit",
		MovePatchIntoNewCommitTooltip:         "Move the patch out of its commit and into a new commit sitting on top of the original commit. This is achieved by starting an interactive rebase at the original commit, applying the patch in reverse, then applying the patch to the index and committing it as a new commit, before continuing the rebase to completion. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchToSelectedCommit:             "Move patch to selected commit (%s)",
		MovePatchToSelectedCommitTooltip:      "Move the patch out of its original commit and into the selected commit. This is achieved by starting an interactive rebase at the original commit, applying the patch in reverse, then continuing the rebase up to the selected commit, before applying the patch forward and amending the seleced commit. The rebase is then continued to completion. If commits between the source and destination commit depend on the patch, you may need to resolve conflicts.",
		CopyPatchToClipboard:                  "Copy patch to clipboard",
		NoMatchesFor:                          "No matches for '%s' %s",
		ExitSearchMode:                        "%s: Exit search mode",
		ExitTextFilterMode:                    "%s: Exit filter mode",
		MatchesFor:                            "matches for '%s' (%d of %d) %s", // lowercase because it's after other text
		SearchKeybindings:                     "%s: Next match, %s: Previous match, %s: Exit search mode",
		SearchPrefix:                          "Search: ",
		FilterPrefix:                          "Filter: ",
		WorktreesTitle:                        "Worktrees",
		WorktreeTitle:                         "Worktree",
		Switch:                                "Switch",
		SwitchToWorktree:                      "Switch to worktree",
		SwitchToWorktreeTooltip:               "Switch to the selected worktree.",
		AlreadyCheckedOutByWorktree:           "This branch is checked out by worktree {{.worktreeName}}. Do you want to switch to that worktree?",
		BranchCheckedOutByWorktree:            "Branch {{.branchName}} is checked out by worktree {{.worktreeName}}",
		DetachWorktreeTooltip:                 "This will run `git checkout --detach` on the worktree so that it stops hogging the branch, but the worktree's working tree will be left alone.",
		Switching:                             "Switching",
		RemoveWorktree:                        "Remove worktree",
		RemoveWorktreeTitle:                   "Remove worktree",
		RemoveWorktreePrompt:                  "Are you sure you want to remove worktree '{{.worktreeName}}'?",
		ForceRemoveWorktreePrompt:             "'{{.worktreeName}}' contains modified or untracked files (to be honest, it could contain both). Are you sure you want to remove it?",
		RemovingWorktree:                      "Deleting worktree",
		DetachWorktree:                        "Detach worktree",
		DetachingWorktree:                     "Detaching worktree",
		AddingWorktree:                        "Adding worktree",
		CantDeleteCurrentWorktree:             "You cannot remove the current worktree!",
		AlreadyInWorktree:                     "You are already in the selected worktree",
		CantDeleteMainWorktree:                "You cannot remove the main worktree!",
		NoWorktreesThisRepo:                   "No worktrees",
		MissingWorktree:                       "(missing)",
		MainWorktree:                          "(main)",
		NewWorktree:                           "New worktree",
		NewWorktreePath:                       "New worktree path",
		NewWorktreeBase:                       "New worktree base ref",
		RemoveWorktreeTooltip:                 "Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory.",
		BranchNameCannotBeBlank:               "Branch name cannot be blank",
		NewBranchName:                         "New branch name",
		NewBranchNameLeaveBlank:               "New branch name (leave blank to checkout {{.default}})",
		ViewWorktreeOptions:                   "View worktree options",
		CreateWorktreeFrom:                    "Create worktree from {{.ref}}",
		CreateWorktreeFromDetached:            "Create worktree from {{.ref}} (detached)",
		LcWorktree:                            "worktree",
		ChangingDirectoryTo:                   "Changing directory to {{.path}}",
		Name:                                  "Name",
		Branch:                                "Branch",
		Path:                                  "Path",
		MarkedBaseCommitStatus:                "Marked a base commit for rebase",
		MarkAsBaseCommit:                      "Mark as base commit for rebase",
		MarkAsBaseCommitTooltip:               "Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command.",
		MarkedCommitMarker:                    "↑↑↑ Will rebase from here ↑↑↑",
		PleaseGoToURL:                         "Please go to {{.url}}",
		DisabledMenuItemPrefix:                "Disabled: ",
		NoCopiedCommits:                       "No copied commits",
		QuickStartInteractiveRebase:           "Start interactive rebase",
		QuickStartInteractiveRebaseTooltip:    "Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.\nIf you would instead like to start an interactive rebase from the selected commit, press `{{.editKey}}`.",
		CannotQuickStartInteractiveRebase:     "Cannot start interactive rebase: the HEAD commit is a merge commit or is present on the main branch, so there is no appropriate base commit to start the rebase from. You can start an interactive rebase from a specific commit by selecting the commit and pressing `{{.editKey}}`.",
		RangeSelectUp:                         "Range select up",
		RangeSelectDown:                       "Range select down",
		RangeSelectNotSupported:               "Action does not support range selection, please select a single item",
		NoItemSelected:                        "No item selected",
		SelectedItemIsNotABranch:              "Selected item is not a branch",
		SelectedItemDoesNotHaveFiles:          "Selected item does not have files to view",
		RangeSelectNotSupportedForSubmodules:  "Range select not supported for submodules",
		OldCherryPickKeyWarning:               "The 'c' key is no longer the default key for copying commits to cherry pick. Please use `{{.copy}}` instead (and `{{.paste}}` to paste). The reason for this change is that the 'v' key for selecting a range of lines when staging is now also used for selecting a range of lines in any list view, meaning that we needed to find a new key for pasting commits, and if we're going to now use `{{.paste}}` for pasting commits, we may as well use `{{.copy}}` for copying them. If you want to configure the keybindings to get the old behaviour, set the following in your config:\n\nkeybinding:\n  universal:\n    toggleRangeSelect: <something other than v>\n  commits:\n    cherryPickCopy: 'c'\n    pasteCommits: 'v'",
		CommandDoesNotSupportOpeningInEditor:  "This command doesn't support switching to the editor",

		Actions: Actions{
			// TODO: combine this with the original keybinding descriptions (those are all in lowercase atm)
			CheckoutCommit:                 "Checkout commit",
			CheckoutTag:                    "Checkout tag",
			CheckoutBranch:                 "Checkout branch",
			ForceCheckoutBranch:            "Force checkout branch",
			DeleteLocalBranch:              "Delete local branch",
			DeleteBranch:                   "Delete branch",
			Merge:                          "Merge",
			RebaseBranch:                   "Rebase branch",
			RenameBranch:                   "Rename branch",
			CreateBranch:                   "Create branch",
			CherryPick:                     "(Cherry-pick) paste commits",
			CheckoutFile:                   "Checkout file",
			DiscardOldFileChange:           "Discard old file change",
			SquashCommitDown:               "Squash commit down",
			FixupCommit:                    "Fixup commit",
			RewordCommit:                   "Reword commit",
			DropCommit:                     "Drop commit",
			EditCommit:                     "Edit commit",
			AmendCommit:                    "Amend commit",
			ResetCommitAuthor:              "Reset commit author",
			SetCommitAuthor:                "Set commit author",
			RevertCommit:                   "Revert commit",
			CreateFixupCommit:              "Create fixup commit",
			SquashAllAboveFixupCommits:     "Squash all above fixup commits",
			CreateLightweightTag:           "Create lightweight tag",
			CreateAnnotatedTag:             "Create annotated tag",
			CopyCommitMessageToClipboard:   "Copy commit message to clipboard",
			CopyCommitSubjectToClipboard:   "Copy commit subject to clipboard",
			CopyCommitDiffToClipboard:      "Copy commit diff to clipboard",
			CopyCommitHashToClipboard:      "Copy full commit hash to clipboard",
			CopyCommitURLToClipboard:       "Copy commit URL to clipboard",
			CopyCommitAuthorToClipboard:    "Copy commit author to clipboard",
			CopyCommitAttributeToClipboard: "Copy to clipboard",
			CopyPatchToClipboard:           "Copy patch to clipboard",
			MoveCommitUp:                   "Move commit up",
			MoveCommitDown:                 "Move commit down",
			CustomCommand:                  "Custom command",

			// TODO: remove
			DiscardAllChangesInDirectory:      "Discard all changes in directory",
			DiscardUnstagedChangesInDirectory: "Discard unstaged changes in directory",

			DiscardAllChangesInFile:         "Discard all changes in selected file(s)",
			DiscardAllUnstagedChangesInFile: "Discard all unstaged changes selected file(s)",
			StageFile:                       "Stage file",
			StageResolvedFiles:              "Stage files whose merge conflicts were resolved",
			UnstageFile:                     "Unstage file",
			UnstageAllFiles:                 "Unstage all files",
			StageAllFiles:                   "Stage all files",
			IgnoreExcludeFile:               "Ignore or exclude file",
			IgnoreFileErr:                   "Cannot ignore .gitignore",
			ExcludeFile:                     "Exclude file",
			ExcludeGitIgnoreErr:             "Cannot exclude .gitignore",
			Commit:                          "Commit",
			EditFile:                        "Edit file",
			Push:                            "Push",
			Pull:                            "Pull",
			OpenFile:                        "Open file",
			StashAllChanges:                 "Stash all changes",
			StashAllChangesKeepIndex:        "Stash all changes and keep index",
			StashStagedChanges:              "Stash staged changes",
			StashUnstagedChanges:            "Stash unstaged changes",
			StashIncludeUntrackedChanges:    "Stash all changes including untracked files",
			GitFlowFinish:                   "git flow finish",
			GitFlowStart:                    "git flow start",
			CopyToClipboard:                 "Copy to clipboard",
			CopySelectedTextToClipboard:     "Copy selected text to clipboard",
			RemovePatchFromCommit:           "Remove patch from commit",
			MovePatchToSelectedCommit:       "Move patch to selected commit",
			MovePatchIntoIndex:              "Move patch into index",
			MovePatchIntoNewCommit:          "Move patch into new commit",
			DeleteRemoteBranch:              "Delete remote branch",
			SetBranchUpstream:               "Set branch upstream",
			AddRemote:                       "Add remote",
			RemoveRemote:                    "Remove remote",
			UpdateRemote:                    "Update remote",
			ApplyPatch:                      "Apply patch",
			Stash:                           "Stash",
			RenameStash:                     "Rename stash",
			RemoveSubmodule:                 "Remove submodule",
			ResetSubmodule:                  "Reset submodule",
			AddSubmodule:                    "Add submodule",
			UpdateSubmoduleUrl:              "Update submodule URL",
			InitialiseSubmodule:             "Initialise submodule",
			BulkInitialiseSubmodules:        "Bulk initialise submodules",
			BulkUpdateSubmodules:            "Bulk update submodules",
			BulkDeinitialiseSubmodules:      "Bulk deinitialise submodules",
			UpdateSubmodule:                 "Update submodule",
			DeleteLocalTag:                  "Delete local tag",
			DeleteRemoteTag:                 "Delete remote tag",
			PushTag:                         "Push tag",
			NukeWorkingTree:                 "Nuke working tree",
			DiscardUnstagedFileChanges:      "Discard unstaged file changes",
			RemoveUntrackedFiles:            "Remove untracked files",
			RemoveStagedFiles:               "Remove staged files",
			SoftReset:                       "Soft reset",
			MixedReset:                      "Mixed reset",
			HardReset:                       "Hard reset",
			FastForwardBranch:               "Fast forward branch",
			Undo:                            "Undo",
			Redo:                            "Redo",
			CopyPullRequestURL:              "Copy pull request URL",
			OpenDiffTool:                    "Open diff tool",
			OpenMergeTool:                   "Open merge tool",
			OpenCommitInBrowser:             "Open commit in browser",
			OpenPullRequest:                 "Open pull request in browser",
			StartBisect:                     "Start bisect",
			ResetBisect:                     "Reset bisect",
			BisectSkip:                      "Bisect skip",
			BisectMark:                      "Bisect mark",
			RemoveWorktree:                  "Remove worktree",
			AddWorktree:                     "Add worktree",
		},
		Bisect: Bisect{
			Mark:                        "Mark current commit (%s) as %s",
			MarkStart:                   "Mark %s as %s (start bisect)",
			SkipCurrent:                 "Skip current commit (%s)",
			SkipSelected:                "Skip selected commit (%s)",
			ResetTitle:                  "Reset 'git bisect'",
			ResetPrompt:                 "Are you sure you want to reset 'git bisect'?",
			ResetOption:                 "Reset bisect",
			ChooseTerms:                 "Choose bisect terms",
			OldTermPrompt:               "Term for old/good commit:",
			NewTermPrompt:               "Term for new/bad commit:",
			BisectMenuTitle:             "Bisect",
			CompleteTitle:               "Bisect complete",
			CompletePrompt:              "Bisect complete! The following commit introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			CompletePromptIndeterminate: "Bisect complete! Some commits were skipped, so any of the following commits may have introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			Bisecting:                   "Bisecting",
		},
		Log: Log{
			EditRebase:               "Beginning interactive rebase at '{{.ref}}'",
			MoveCommitUp:             "Moving TODO down: '{{.shortHash}}'",
			MoveCommitDown:           "Moving TODO down: '{{.shortHash}}'",
			CherryPickCommits:        "Cherry-picking commits:\n'{{.commitLines}}'",
			HandleUndo:               "Undoing last conflict resolution",
			HandleMidRebaseCommand:   "Updating rebase action of commit {{.shortHash}} to '{{.action}}'",
			RemoveFile:               "Deleting path '{{.path}}'",
			CopyToClipboard:          "Copying '{{.str}}' to clipboard",
			Remove:                   "Removing '{{.filename}}'",
			CreateFileWithContent:    "Creating file '{{.path}}'",
			AppendingLineToFile:      "Appending '{{.line}}' to file '{{.filename}}'",
			EditRebaseFromBaseCommit: "Beginning interactive rebase from '{{.baseCommit}}' onto '{{.targetBranchName}}",
		},
		BreakingChangesTitle: "Breaking Changes",
		BreakingChangesMessage: `You are updating to a new version of lazygit which contains breaking changes. Please review the notes below and update your configuration if necessary.
For more information, see the full release notes at <https://github.com/jesseduffield/lazygit/releases>.`,
		BreakingChangesByVersion: map[string]string{
			"0.41.0": `- When you press 'g' to bring up the git reset menu, the 'mixed' option is now the first and default, rather than 'soft'. This is because 'mixed' is the most commonly used option.
- The commit message panel now automatically hard-wraps by default (i.e. it adds newline characters when you reach the margin). You can adjust the config like so:

git:
  commit:
    autoWrapCommitMessage: true
    autoWrapWidth: 72

- The 'v' key was already being used in the staging view to start a range select, but now you can use it to start a range select in any view. Unfortunately this clashes with the 'v' keybinding for pasting commits (cherry-pick), so now pasting commits is done via 'shift+V' and for the sake of consistency, copying commits is now done via 'shift+C' instead of just 'c'. Note that the 'v' keybinding is only one way to start a range-select: you can use shift+up/down arrow instead. So, if you want to configure the cherry-pick keybindings to get the old behaviour, set the following in your config:

keybinding:
  universal:
      toggleRangeSelect: <something other than v>
    commits:
      cherryPickCopy: 'c'
      pasteCommits: 'v'

- Squashing fixups using 'shift-S' now brings up a menu, with the default option being to squash all fixup commits in the branch. The original behaviour of only squashing fixup commits above the selected commit is still available as the second option in that menu.
- Push/pull/fetch loading statuses are now shown against the branch rather than in a popup. This allows you to e.g. fetch multiple branches in parallel and see the status for each branch.
- The git log graph in the commits view is now always shown by default (previously it was only shown when the view was maximised). If you find this too noisy, you can change it back via ctrl+L -> 'Show git graph' -> 'when maximised'
- Pressing space on a remote branch used to show a prompt for entering a name for a new local branch to check out from the remote branch. Now it just checks out the remote branch directly, letting you choose between a new local branch with the same name, or a detached head. The old behavior is still available via the 'n' keybinding.
- Filtering (e.g. when pressing '/') is less fuzzy by default; it only matches substrings now. Multiple substrings can be matched by separating them with spaces. If you want to revert to the old behavior, set the following in your config:

gui:
  filterMode: 'fuzzy'
	  `,
		},
	}
}
