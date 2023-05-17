/*

Todo list when making a new translation
- Copy this file and rename it to the language you want to translate to like someLanguage.go
- Change the EnglishTranslationSet() name to the language you want to translate to like SomeLanguageTranslationSet()
- Add an entry of someLanguage in GetTranslationSets()
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
	Menu                                string
	Execute                             string
	ToggleStaged                        string
	ToggleStagedAll                     string
	ToggleTreeView                      string
	OpenMergeTool                       string
	Refresh                             string
	Push                                string
	Pull                                string
	Scroll                              string
	FileFilter                          string
	FilterStagedFiles                   string
	FilterUnstagedFiles                 string
	ResetFilter                         string
	MergeConflictsTitle                 string
	Checkout                            string
	NoChangedFiles                      string
	PullWait                            string
	PushWait                            string
	FetchWait                           string
	SoftReset                           string
	AlreadyCheckedOutBranch             string
	SureForceCheckout                   string
	ForceCheckoutBranch                 string
	BranchName                          string
	NewBranchNameBranchOff              string
	CantDeleteCheckOutBranch            string
	DeleteBranch                        string
	DeleteBranchMessage                 string
	ForceDeleteBranchMessage            string
	RebaseBranch                        string
	CantRebaseOntoSelf                  string
	CantMergeBranchIntoItself           string
	ForceCheckout                       string
	CheckoutByName                      string
	NewBranch                           string
	NoBranchesThisRepo                  string
	CommitWithoutMessageErr             string
	Close                               string
	CloseCancel                         string
	Confirm                             string
	Quit                                string
	SquashDown                          string
	FixupCommit                         string
	CannotSquashOrFixupFirstCommit      string
	Fixup                               string
	SureFixupThisCommit                 string
	SureSquashThisCommit                string
	Squash                              string
	PickCommit                          string
	RevertCommit                        string
	RewordCommit                        string
	DeleteCommit                        string
	MoveDownCommit                      string
	MoveUpCommit                        string
	EditCommit                          string
	AmendToCommit                       string
	ResetAuthor                         string
	SetAuthor                           string
	SetResetCommitAuthor                string
	SetAuthorPromptTitle                string
	SureResetCommitAuthor               string
	RenameCommitEditor                  string
	NoCommitsThisBranch                 string
	UpdateRefHere                       string
	Error                               string
	Undo                                string
	UndoReflog                          string
	RedoReflog                          string
	UndoTooltip                         string
	RedoTooltip                         string
	DiscardAllTooltip                   string
	DiscardUnstagedTooltip              string
	Pop                                 string
	Drop                                string
	Apply                               string
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
	RenameStash                         string
	RenameStashPrompt                   string
	OpenConfig                          string
	EditConfig                          string
	ForcePush                           string
	ForcePushPrompt                     string
	ForcePushDisabled                   string
	UpdatesRejectedAndForcePushDisabled string
	CheckForUpdate                      string
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
	EditFile                            string
	OpenFile                            string
	IgnoreFile                          string
	ExcludeFile                         string
	RefreshFiles                        string
	MergeIntoCurrentBranch              string
	ConfirmQuit                         string
	SwitchRepo                          string
	AllBranchesLogGraph                 string
	UnsupportedGitService               string
	CopyPullRequestURL                  string
	NoBranchOnRemote                    string
	Fetch                               string
	NoAutomaticGitFetchTitle            string
	NoAutomaticGitFetchBody             string
	FileEnter                           string
	FileStagingRequirements             string
	StageSelection                      string
	DiscardSelection                    string
	ToggleDragSelect                    string
	ToggleSelectHunk                    string
	ToggleSelectionForPatch             string
	EditHunk                            string
	ToggleStagingPanel                  string
	ReturnToFilesPanel                  string
	FastForward                         string
	Fetching                            string
	FoundConflictsTitle                 string
	ViewConflictsMenuItem               string
	AbortMenuItem                       string
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
	Continue                            string
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
	RewordNotSupported                  string
	ChangingThisActionIsNotAllowed      string
	CherryPickCopy                      string
	CherryPickCopyRange                 string
	PasteCommits                        string
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
	ScrollUpMainPanel                   string
	ScrollDownMainPanel                 string
	AmendCommitTitle                    string
	AmendCommitPrompt                   string
	DeleteCommitTitle                   string
	DeleteCommitPrompt                  string
	SquashingStatus                     string
	FixingStatus                        string
	DeletingStatus                      string
	MovingStatus                        string
	RebasingStatus                      string
	MergingStatus                       string
	LowercaseRebasingStatus             string
	LowercaseMergingStatus              string
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
	ViewItemFiles                       string
	CommitFilesTitle                    string
	CheckoutCommitFile                  string
	CanOnlyDiscardFromLocalCommits      string
	DiscardOldFileChange                string
	DiscardFileChangesTitle             string
	DiscardFileChangesPrompt            string
	DiscardAddedFileChangesPrompt       string
	DiscardDeletedFileChangesPrompt     string
	DiscardNotSupportedForDirectory     string
	DisabledForGPG                      string
	CreateRepo                          string
	BareRepo                            string
	InitialBranch                       string
	NoRecentRepositories                string
	IncorrectNotARepository             string
	AutoStashTitle                      string
	AutoStashPrompt                     string
	StashPrefix                         string
	ViewDiscardOptions                  string
	Cancel                              string
	DiscardAllChanges                   string
	DiscardUnstagedChanges              string
	DiscardAllChangesToAllFiles         string
	DiscardAnyUnstagedChanges           string
	DiscardUntrackedFiles               string
	DiscardStagedChanges                string
	HardReset                           string
	ViewResetOptions                    string
	CreateFixupCommit                   string
	CreateFixupCommitDescription        string
	SquashAboveCommits                  string
	SureSquashAboveCommits              string
	SureCreateFixupCommit               string
	ExecuteCustomCommand                string
	CustomCommand                       string
	CommitChangesWithoutHook            string
	SkipHookPrefixNotConfigured         string
	ResetTo                             string
	PressEnterToReturn                  string
	ViewStashOptions                    string
	StashAllChanges                     string
	StashStagedChanges                  string
	StashAllChangesKeepIndex            string
	StashUnstagedChanges                string
	StashIncludeUntrackedChanges        string
	StashOptions                        string
	NotARepository                      string
	Jump                                string
	ScrollLeftRight                     string
	ScrollLeft                          string
	ScrollRight                         string
	DiscardPatch                        string
	DiscardPatchConfirm                 string
	CantPatchWhileRebasingError         string
	ToggleAddToPatch                    string
	ToggleAllInPatch                    string
	UpdatingPatch                       string
	ViewPatchOptions                    string
	PatchOptionsTitle                   string
	NoPatchError                        string
	EnterFile                           string
	ExitCustomPatchBuilder              string
	EnterUpstream                       string
	InvalidUpstream                     string
	ReturnToRemotesList                 string
	AddNewRemote                        string
	NewRemoteName                       string
	NewRemoteUrl                        string
	EditRemoteName                      string
	EditRemoteUrl                       string
	RemoveRemote                        string
	RemoveRemotePrompt                  string
	DeleteRemoteBranch                  string
	DeleteRemoteBranchMessage           string
	SetAsUpstream                       string
	SetUpstream                         string
	UnsetUpstream                       string
	SetUpstreamTitle                    string
	SetUpstreamMessage                  string
	EditRemote                          string
	TagCommit                           string
	TagMenuTitle                        string
	TagNameTitle                        string
	TagMessageTitle                     string
	LightweightTag                      string
	AnnotatedTag                        string
	DeleteTag                           string
	DeleteTagTitle                      string
	DeleteTagPrompt                     string
	PushTagTitle                        string
	PushTag                             string
	CreateTag                           string
	CreateTagTitle                      string
	FetchRemote                         string
	FetchingRemoteStatus                string
	CheckoutCommit                      string
	SureCheckoutThisCommit              string
	GitFlowOptions                      string
	NotAGitFlowBranch                   string
	NewBranchNamePrompt                 string
	IgnoreTracked                       string
	ExcludeTracked                      string
	IgnoreTrackedPrompt                 string
	ExcludeTrackedPrompt                string
	ViewResetToUpstreamOptions          string
	NextScreenMode                      string
	PrevScreenMode                      string
	StartSearch                         string
	StartFilter                         string
	Panel                               string
	Keybindings                         string
	KeybindingsLegend                   string
	RenameBranch                        string
	SetUnsetUpstream                    string
	NewGitFlowBranchPrompt              string
	RenameBranchWarning                 string
	OpenMenu                            string
	ResetCherryPick                     string
	NextTab                             string
	PrevTab                             string
	CantUndoWhileRebasing               string
	CantRedoWhileRebasing               string
	MustStashWarning                    string
	MustStashTitle                      string
	ConfirmationTitle                   string
	PrevPage                            string
	NextPage                            string
	GotoTop                             string
	GotoBottom                          string
	FilteringBy                         string
	ResetInParentheses                  string
	OpenFilteringMenu                   string
	FilterBy                            string
	ExitFilterMode                      string
	FilterPathOption                    string
	EnterFileName                       string
	FilteringMenuTitle                  string
	MustExitFilterModeTitle             string
	MustExitFilterModePrompt            string
	Diff                                string
	EnterRefToDiff                      string
	EnterRefName                        string
	ExitDiffMode                        string
	DiffingMenuTitle                    string
	SwapDiff                            string
	OpenDiffingMenu                     string
	OpenExtrasMenu                      string
	ShowingGitDiff                      string
	CommitDiff                          string
	CopyCommitShaToClipboard            string
	CommitSha                           string
	CommitURL                           string
	CopyCommitMessageToClipboard        string
	CommitMessage                       string
	CommitAuthor                        string
	CopyCommitAttributeToClipboard      string
	CopyBranchNameToClipboard           string
	CopyFileNameToClipboard             string
	CopyCommitFileNameToClipboard       string
	CommitPrefixPatternError            string
	CopySelectedTexToClipboard          string
	NoFilesStagedTitle                  string
	NoFilesStagedPrompt                 string
	BranchNotFoundTitle                 string
	BranchNotFoundPrompt                string
	BranchUnknown                       string
	DiscardChangeTitle                  string
	DiscardChangePrompt                 string
	CreateNewBranchFromCommit           string
	BuildingPatch                       string
	ViewCommits                         string
	MinGitVersionError                  string
	RunningCustomCommandStatus          string
	SubmoduleStashAndReset              string
	AndResetSubmodules                  string
	EnterSubmodule                      string
	CopySubmoduleNameToClipboard        string
	RemoveSubmodule                     string
	RemoveSubmodulePrompt               string
	ResettingSubmoduleStatus            string
	NewSubmoduleName                    string
	NewSubmoduleUrl                     string
	NewSubmodulePath                    string
	AddSubmodule                        string
	AddingSubmoduleStatus               string
	UpdateSubmoduleUrl                  string
	UpdatingSubmoduleUrlStatus          string
	EditSubmoduleUrl                    string
	InitializingSubmoduleStatus         string
	InitSubmodule                       string
	SubmoduleUpdate                     string
	UpdatingSubmoduleStatus             string
	BulkInitSubmodules                  string
	BulkUpdateSubmodules                string
	BulkDeinitSubmodules                string
	ViewBulkSubmoduleOptions            string
	BulkSubmoduleOptions                string
	RunningCommand                      string
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
	CopiedToClipboard                   string
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
	IgnoreWhitespaceDiffViewSubTitle    string
	IgnoreWhitespaceNotSupportedHere    string
	IncreaseContextInDiffView           string
	DecreaseContextInDiffView           string
	CreatePullRequestOptions            string
	DefaultBranch                       string
	SelectBranch                        string
	CreatePullRequest                   string
	SelectConfigFile                    string
	NoConfigFileFoundErr                string
	LoadingFileSuggestions              string
	LoadingCommits                      string
	MustSpecifyOriginError              string
	GitOutput                           string
	GitCommandFailed                    string
	AbortTitle                          string
	AbortPrompt                         string
	OpenLogMenu                         string
	LogMenuTitle                        string
	ToggleShowGitGraphAll               string
	ShowGitGraph                        string
	SortCommits                         string
	CantChangeContextSizeError          string
	OpenCommitInBrowser                 string
	ViewBisectOptions                   string
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
	CommitsCopied                       string
	CommitCopied                        string
	ResetPatch                          string
	ApplyPatch                          string
	ApplyPatchInReverse                 string
	RemovePatchFromOriginalCommit       string
	MovePatchOutIntoIndex               string
	MovePatchIntoNewCommit              string
	MovePatchToSelectedCommit           string
	CopyPatchToClipboard                string
	NoMatchesFor                        string
	MatchesFor                          string
	SearchKeybindings                   string
	SearchPrefix                        string
	FilterPrefix                        string
	ExitSearchMode                      string
	ExitTextFilterMode                  string
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
	Bisecting                   string
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
	IgnoreExcludeFile                 string
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
		EasterEgg:                           "Easter egg",
		UnstagedChanges:                     "Unstaged changes",
		StagedChanges:                       "Staged changes",
		MainTitle:                           "Main",
		MergeConfirmTitle:                   "Merge",
		StagingTitle:                        "Main panel (staging)",
		MergingTitle:                        "Main panel (merging)",
		NormalTitle:                         "Main panel (normal)",
		LogTitle:                            "Log",
		CommitSummary:                       "Commit summary",
		CredentialsUsername:                 "Username",
		CredentialsPassword:                 "Password",
		CredentialsPassphrase:               "Enter passphrase for SSH key",
		CredentialsPIN:                      "Enter PIN for SSH key",
		PassUnameWrong:                      "Password, passphrase and/or username wrong",
		CommitChanges:                       "Commit changes",
		AmendLastCommit:                     "Amend last commit",
		AmendLastCommitTitle:                "Amend last commit",
		SureToAmend:                         "Are you sure you want to amend last commit? Afterwards, you can change the commit message from the commits panel.",
		NoCommitToAmend:                     "There's no commit to amend.",
		CommitChangesWithEditor:             "Commit changes using git editor",
		StatusTitle:                         "Status",
		Menu:                                "Menu",
		Execute:                             "Execute",
		ToggleStaged:                        "Toggle staged",
		ToggleStagedAll:                     "Stage/unstage all",
		ToggleTreeView:                      "Toggle file tree view",
		OpenMergeTool:                       "Open external merge tool (git mergetool)",
		Refresh:                             "Refresh",
		Push:                                "Push",
		Pull:                                "Pull",
		Scroll:                              "Scroll",
		MergeConflictsTitle:                 "Merge conflicts",
		Checkout:                            "Checkout",
		FileFilter:                          "Filter files by status",
		FilterStagedFiles:                   "Show only staged files",
		FilterUnstagedFiles:                 "Show only unstaged files",
		ResetFilter:                         "Reset filter",
		NoChangedFiles:                      "No changed files",
		PullWait:                            "Pulling...",
		PushWait:                            "Pushing...",
		FetchWait:                           "Fetching...",
		SoftReset:                           "Soft reset",
		AlreadyCheckedOutBranch:             "You have already checked out this branch",
		SureForceCheckout:                   "Are you sure you want force checkout? You will lose all local changes",
		ForceCheckoutBranch:                 "Force checkout branch",
		BranchName:                          "Branch name",
		NewBranchNameBranchOff:              "New branch name (branch is off of '{{.branchName}}')",
		CantDeleteCheckOutBranch:            "You cannot delete the checked out branch!",
		DeleteBranch:                        "Delete branch",
		DeleteBranchMessage:                 "Are you sure you want to delete the branch '{{.selectedBranchName}}'?",
		ForceDeleteBranchMessage:            "'{{.selectedBranchName}}' is not fully merged. Are you sure you want to delete it?",
		RebaseBranch:                        "Rebase checked-out branch onto this branch",
		CantRebaseOntoSelf:                  "You cannot rebase a branch onto itself",
		CantMergeBranchIntoItself:           "You cannot merge a branch into itself",
		ForceCheckout:                       "Force checkout",
		CheckoutByName:                      "Checkout by name",
		NewBranch:                           "New branch",
		NoBranchesThisRepo:                  "No branches for this repo",
		CommitWithoutMessageErr:             "You cannot commit without a commit message",
		Close:                               "Close",
		CloseCancel:                         "Close/Cancel",
		Confirm:                             "Confirm",
		Quit:                                "Quit",
		SquashDown:                          "Squash down",
		FixupCommit:                         "Fixup commit",
		NoCommitsThisBranch:                 "No commits for this branch",
		UpdateRefHere:                       "Update branch '{{.ref}}' here",
		CannotSquashOrFixupFirstCommit:      "There's no commit below to squash into",
		Fixup:                               "Fixup",
		SureFixupThisCommit:                 "Are you sure you want to 'fixup' this commit? It will be merged into the commit below",
		SureSquashThisCommit:                "Are you sure you want to squash this commit into the commit below?",
		Squash:                              "Squash",
		PickCommit:                          "Pick commit (when mid-rebase)",
		RevertCommit:                        "Revert commit",
		RewordCommit:                        "Reword commit",
		DeleteCommit:                        "Delete commit",
		MoveDownCommit:                      "Move commit down one",
		MoveUpCommit:                        "Move commit up one",
		EditCommit:                          "Edit commit",
		AmendToCommit:                       "Amend commit with staged changes",
		ResetAuthor:                         "Reset author",
		SetAuthor:                           "Set author",
		SetResetCommitAuthor:                "Set/Reset commit author",
		SetAuthorPromptTitle:                "Set author (must look like 'Name <Email>')",
		SureResetCommitAuthor:               "The author field of this commit will be updated to match the configured user. This also renews the author timestamp. Continue?",
		RenameCommitEditor:                  "Reword commit with editor",
		Error:                               "Error",
		PickHunk:                            "Pick hunk",
		PickAllHunks:                        "Pick all hunks",
		Undo:                                "Undo",
		UndoReflog:                          "Undo",
		RedoReflog:                          "Redo",
		UndoTooltip:                         "The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		RedoTooltip:                         "The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		DiscardAllTooltip:                   "Discard both staged and unstaged changes in '{{.path}}'.",
		DiscardUnstagedTooltip:              "Discard unstaged changes in '{{.path}}'.",
		Pop:                                 "Pop",
		Drop:                                "Drop",
		Apply:                               "Apply",
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
		RenameStash:                         "Rename stash",
		RenameStashPrompt:                   "Rename stash: {{.stashName}}",
		OpenConfig:                          "Open config file",
		EditConfig:                          "Edit config file",
		ForcePush:                           "Force push",
		ForcePushPrompt:                     "Your branch has diverged from the remote branch. Press 'esc' to cancel, or 'enter' to force push.",
		ForcePushDisabled:                   "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejectedAndForcePushDisabled: "Updates were rejected and you have disabled force pushing",
		CheckForUpdate:                      "Check for update",
		CheckingForUpdates:                  "Checking for updates...",
		UpdateAvailableTitle:                "Update available!",
		UpdateAvailable:                     "Download and install version {{.newVersion}}?",
		UpdateInProgressWaitingStatus:       "Updating",
		UpdateCompletedTitle:                "Update completed!",
		UpdateCompleted:                     "Update has been installed successfully. Restart lazygit for it to take effect.",
		FailedToRetrieveLatestVersionErr:    "Failed to retrieve version information",
		OnLatestVersionErr:                  "You already have the latest version",
		MajorVersionErr:                     "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		CouldNotFindBinaryErr:               "Could not find any binary at {{.url}}",
		UpdateFailedErr:                     "Update failed: {{.errMessage}}",
		ConfirmQuitDuringUpdateTitle:        "Currently updating",
		ConfirmQuitDuringUpdate:             "An update is in progress. Are you sure you want to quit?",
		MergeToolTitle:                      "Merge tool",
		MergeToolPrompt:                     "Are you sure you want to open `git mergetool`?",
		IntroPopupMessage:                   englishIntroPopupMessage,
		DeprecatedEditConfigWarning:         englishDeprecatedEditConfigWarning,
		GitconfigParseErr:                   `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		EditFile:                            `Edit file`,
		OpenFile:                            `Open file`,
		IgnoreFile:                          `Add to .gitignore`,
		ExcludeFile:                         `Add to .git/info/exclude`,
		RefreshFiles:                        `Refresh files`,
		MergeIntoCurrentBranch:              `Merge into currently checked out branch`,
		ConfirmQuit:                         `Are you sure you want to quit?`,
		SwitchRepo:                          `Switch to a recent repo`,
		AllBranchesLogGraph:                 `Show all branch logs`,
		UnsupportedGitService:               `Unsupported git service`,
		CreatePullRequest:                   `Create pull request`,
		CopyPullRequestURL:                  `Copy pull request URL to clipboard`,
		NoBranchOnRemote:                    `This branch doesn't exist on remote. You need to push it to remote first.`,
		Fetch:                               `Fetch`,
		NoAutomaticGitFetchTitle:            `No automatic git fetch`,
		NoAutomaticGitFetchBody:             `Lazygit can't use "git fetch" in a private repo; use 'f' in the files panel to run "git fetch" manually`,
		FileEnter:                           `Stage individual hunks/lines for file, or collapse/expand for directory`,
		FileStagingRequirements:             `Can only stage individual lines for tracked files`,
		StageSelection:                      `Toggle line staged / unstaged`,
		DiscardSelection:                    `Discard change (git reset)`,
		ToggleDragSelect:                    `Toggle drag select`,
		ToggleSelectHunk:                    `Toggle select hunk`,
		ToggleSelectionForPatch:             `Add/Remove line(s) to patch`,
		EditHunk:                            `Edit hunk`,
		ToggleStagingPanel:                  `Switch to other panel (staged/unstaged changes)`,
		ReturnToFilesPanel:                  `Return to files panel`,
		FastForward:                         `Fast-forward this branch from its upstream`,
		Fetching:                            "Fetching and fast-forwarding {{.from}} -> {{.to}} ...",
		FoundConflictsTitle:                 "Conflicts!",
		ViewConflictsMenuItem:               "View conflicts",
		AbortMenuItem:                       "Abort the %s",
		ViewMergeRebaseOptions:              "View merge/rebase options",
		NotMergingOrRebasing:                "You are currently neither rebasing nor merging",
		AlreadyRebasing:                     "Can't perform this action during a rebase",
		RecentRepos:                         "Recent repositories",
		MergeOptionsTitle:                   "Merge options",
		RebaseOptionsTitle:                  "Rebase options",
		CommitMessageTitle:                  "Commit summary",
		CommitDescriptionTitle:              "Commit description",
		CommitDescriptionSubTitle:           "Press {{.togglePanelKeyBinding}} to toggle focus",
		LocalBranchesTitle:                  "Local branches",
		SearchTitle:                         "Search",
		TagsTitle:                           "Tags",
		MenuTitle:                           "Menu",
		RemotesTitle:                        "Remotes",
		RemoteBranchesTitle:                 "Remote branches",
		PatchBuildingTitle:                  "Main panel (patch building)",
		InformationTitle:                    "Information",
		SecondaryTitle:                      "Secondary",
		ReflogCommitsTitle:                  "Reflog",
		GlobalTitle:                         "Global keybindings",
		ConflictsResolved:                   "All merge conflicts resolved. Continue?",
		Continue:                            "Continue",
		Keybindings:                         "Keybindings",
		RebasingTitle:                       "Rebase '{{.checkedOutBranch}}' onto '{{.ref}}'",
		SimpleRebase:                        "Simple rebase",
		InteractiveRebase:                   "Interactive rebase",
		InteractiveRebaseTooltip:            "Begin an interactive rebase with a break at the start, so you can update the TODO commits before continuing",
		ConfirmMerge:                        "Are you sure you want to merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}'?",
		FwdNoUpstream:                       "Cannot fast-forward a branch with no upstream",
		FwdNoLocalUpstream:                  "Cannot fast-forward a branch whose remote is not registered locally",
		FwdCommitsToPush:                    "Cannot fast-forward a branch with commits to push",
		ErrorOccurred:                       "An error occurred! Please create an issue at",
		NoRoom:                              "Not enough room",
		YouAreHere:                          "YOU ARE HERE",
		YouDied:                             "YOU DIED!",
		RewordNotSupported:                  "Rewording commits while interactively rebasing is not currently supported",
		ChangingThisActionIsNotAllowed:      "Changing this kind of rebase todo entry is not allowed",
		CherryPickCopy:                      "Copy commit (cherry-pick)",
		CherryPickCopyRange:                 "Copy commit range (cherry-pick)",
		PasteCommits:                        "Paste commits (cherry-pick)",
		SureCherryPick:                      "Are you sure you want to cherry-pick the copied commits onto this branch?",
		CherryPick:                          "Cherry-pick",
		Donate:                              "Donate",
		AskQuestion:                         "Ask Question",
		PrevLine:                            "Select previous line",
		NextLine:                            "Select next line",
		PrevHunk:                            "Select previous hunk",
		NextHunk:                            "Select next hunk",
		PrevConflict:                        "Select previous conflict",
		NextConflict:                        "Select next conflict",
		SelectPrevHunk:                      "Select previous hunk",
		SelectNextHunk:                      "Select next hunk",
		ScrollDown:                          "Scroll down",
		ScrollUp:                            "Scroll up",
		ScrollUpMainPanel:                   "Scroll up main panel",
		ScrollDownMainPanel:                 "Scroll down main panel",
		AmendCommitTitle:                    "Amend commit",
		AmendCommitPrompt:                   "Are you sure you want to amend this commit with your staged files?",
		DeleteCommitTitle:                   "Delete commit",
		DeleteCommitPrompt:                  "Are you sure you want to delete this commit?",
		SquashingStatus:                     "Squashing",
		FixingStatus:                        "Fixing up",
		DeletingStatus:                      "Deleting",
		MovingStatus:                        "Moving",
		RebasingStatus:                      "Rebasing",
		MergingStatus:                       "Merging",
		LowercaseRebasingStatus:             "rebasing", // lowercase because it shows up in parentheses
		LowercaseMergingStatus:              "merging",  // lowercase because it shows up in parentheses
		AmendingStatus:                      "Amending",
		CherryPickingStatus:                 "Cherry-picking",
		UndoingStatus:                       "Undoing",
		RedoingStatus:                       "Redoing",
		CheckingOutStatus:                   "Checking out",
		CommittingStatus:                    "Committing",
		CommitFiles:                         "Commit files",
		SubCommitsDynamicTitle:              "Commits (%s)",
		CommitFilesDynamicTitle:             "Diff files (%s)",
		RemoteBranchesDynamicTitle:          "Remote branches (%s)",
		ViewItemFiles:                       "View selected item's files",
		CommitFilesTitle:                    "Commit files",
		CheckoutCommitFile:                  "Checkout file",
		CanOnlyDiscardFromLocalCommits:      "Changes can only be discarded from local commits",
		DiscardOldFileChange:                "Discard this commit's changes to this file",
		DiscardFileChangesTitle:             "Discard file changes",
		DiscardFileChangesPrompt:            "Are you sure you want to discard this commit's changes to this file?",
		DiscardAddedFileChangesPrompt:       "Are you sure you want to discard this commit's changes to this file? The file was added in this commit, so it will be deleted again.",
		DiscardDeletedFileChangesPrompt:     "Are you sure you want to discard this commit's changes to this file? The file was deleted in this commit, so it will reappear.",
		DiscardNotSupportedForDirectory:     "Discarding changes is not supported for entire directories. Please use a custom patch for this.",
		DisabledForGPG:                      "Feature not available for users using GPG",
		CreateRepo:                          "Not in a git repository. Create a new git repository? (y/n): ",
		BareRepo:                            "You've attempted to open Lazygit in a bare repo but Lazygit does not yet support bare repos. Open most recent repo? (y/n) ",
		InitialBranch:                       "Branch name? (leave empty for git's default): ",
		NoRecentRepositories:                "Must open lazygit in a git repository. No valid recent repositories. Exiting.",
		IncorrectNotARepository:             "The value of 'notARepository' is incorrect. It should be one of 'prompt', 'create', 'skip', or 'quit'.",
		AutoStashTitle:                      "Autostash?",
		AutoStashPrompt:                     "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		StashPrefix:                         "Auto-stashing changes for ",
		ViewDiscardOptions:                  "View 'discard changes' options",
		Cancel:                              "Cancel",
		DiscardAllChanges:                   "Discard all changes",
		DiscardUnstagedChanges:              "Discard unstaged changes",
		DiscardAllChangesToAllFiles:         "Nuke working tree",
		DiscardAnyUnstagedChanges:           "Discard unstaged changes",
		DiscardUntrackedFiles:               "Discard untracked files",
		DiscardStagedChanges:                "Discard staged changes",
		HardReset:                           "Hard reset",
		ViewResetOptions:                    `View reset options`,
		CreateFixupCommitDescription:        `Create fixup commit for this commit`,
		SquashAboveCommits:                  `Squash all 'fixup!' commits above selected commit (autosquash)`,
		SureSquashAboveCommits:              `Are you sure you want to squash all fixup! commits above {{.commit}}?`,
		CreateFixupCommit:                   `Create fixup commit`,
		SureCreateFixupCommit:               `Are you sure you want to create a fixup! commit for commit {{.commit}}?`,
		ExecuteCustomCommand:                "Execute custom command",
		CustomCommand:                       "Custom command:",
		CommitChangesWithoutHook:            "Commit changes without pre-commit hook",
		SkipHookPrefixNotConfigured:         "You have not configured a commit message prefix for skipping hooks. Set `git.skipHookPrefix = 'WIP'` in your config",
		ResetTo:                             `Reset to`,
		PressEnterToReturn:                  "Press enter to return to lazygit",
		ViewStashOptions:                    "View stash options",
		StashAllChanges:                     "Stash all changes",
		StashStagedChanges:                  "Stash staged changes",
		StashAllChangesKeepIndex:            "Stash all changes and keep index",
		StashUnstagedChanges:                "Stash unstaged changes",
		StashIncludeUntrackedChanges:        "Stash all changes including untracked files",
		StashOptions:                        "Stash options",
		NotARepository:                      "Error: must be run inside a git repository",
		Jump:                                "Jump to panel",
		ScrollLeftRight:                     "Scroll left/right",
		ScrollLeft:                          "Scroll left",
		ScrollRight:                         "Scroll right",
		DiscardPatch:                        "Discard patch",
		DiscardPatchConfirm:                 "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		CantPatchWhileRebasingError:         "You cannot build a patch or run patch commands while in a merging or rebasing state",
		ToggleAddToPatch:                    "Toggle file included in patch",
		ToggleAllInPatch:                    "Toggle all files included in patch",
		UpdatingPatch:                       "Updating patch",
		ViewPatchOptions:                    "View custom patch options",
		PatchOptionsTitle:                   "Patch options",
		NoPatchError:                        "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		EnterFile:                           "Enter file to add selected lines to the patch (or toggle directory collapsed)",
		ExitCustomPatchBuilder:              `Exit custom patch builder`,
		EnterUpstream:                       `Enter upstream as '<remote> <branchname>'`,
		InvalidUpstream:                     "Invalid upstream. Must be in the format '<remote> <branchname>'",
		ReturnToRemotesList:                 `Return to remotes list`,
		AddNewRemote:                        `Add new remote`,
		NewRemoteName:                       `New remote name:`,
		NewRemoteUrl:                        `New remote url:`,
		EditRemoteName:                      `Enter updated remote name for {{.remoteName}}:`,
		EditRemoteUrl:                       `Enter updated remote url for {{.remoteName}}:`,
		RemoveRemote:                        `Remove remote`,
		RemoveRemotePrompt:                  "Are you sure you want to remove remote",
		DeleteRemoteBranch:                  "Delete Remote Branch",
		DeleteRemoteBranchMessage:           "Are you sure you want to delete remote branch",
		SetAsUpstream:                       "Set as upstream of checked-out branch",
		SetUpstream:                         "Set upstream of selected branch",
		UnsetUpstream:                       "Unset upstream of selected branch",
		SetUpstreamTitle:                    "Set upstream branch",
		SetUpstreamMessage:                  "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		EditRemote:                          "Edit remote",
		TagCommit:                           "Tag commit",
		TagMenuTitle:                        "Create tag",
		TagNameTitle:                        "Tag name:",
		TagMessageTitle:                     "Tag message:",
		AnnotatedTag:                        "Annotated tag",
		LightweightTag:                      "Lightweight tag",
		DeleteTag:                           "Delete tag",
		DeleteTagTitle:                      "Delete tag",
		DeleteTagPrompt:                     "Are you sure you want to delete tag '{{.tagName}}'?",
		PushTagTitle:                        "Remote to push tag '{{.tagName}}' to:",
		PushTag:                             "Push tag",
		CreateTag:                           "Create tag",
		CreateTagTitle:                      "Tag name:",
		FetchRemote:                         "Fetch remote",
		FetchingRemoteStatus:                "Fetching remote",
		CheckoutCommit:                      "Checkout commit",
		SureCheckoutThisCommit:              "Are you sure you want to checkout this commit?",
		GitFlowOptions:                      "Show git-flow options",
		NotAGitFlowBranch:                   "This does not seem to be a git flow branch",
		NewGitFlowBranchPrompt:              "New {{.branchType}} name:",

		IgnoreTracked:              "Ignore tracked file",
		IgnoreTrackedPrompt:        "Are you sure you want to ignore a tracked file?",
		ExcludeTracked:             "Exclude tracked file",
		ExcludeTrackedPrompt:       "Are you sure you want to exclude a tracked file?",
		ViewResetToUpstreamOptions: "View upstream reset options",
		NextScreenMode:             "Next screen mode (normal/half/fullscreen)",
		PrevScreenMode:             "Prev screen mode",
		StartSearch:                "Search the current view by text",
		StartFilter:                "Filter the current view by text",
		Panel:                      "Panel",
		KeybindingsLegend:          "Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b",
		RenameBranch:               "Rename branch",
		SetUnsetUpstream:           "Set/Unset upstream",
		NewBranchNamePrompt:        "Enter new branch name for branch",
		RenameBranchWarning:        "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		OpenMenu:                   "Open menu",
		ResetCherryPick:            "Reset cherry-picked (copied) commits selection",
		NextTab:                    "Next tab",
		PrevTab:                    "Previous tab",
		CantUndoWhileRebasing:      "Can't undo while rebasing",
		CantRedoWhileRebasing:      "Can't redo while rebasing",
		MustStashWarning:           "Pulling a patch out into the index requires stashing and unstashing your changes. If something goes wrong, you'll be able to access your files from the stash. Continue?",
		MustStashTitle:             "Must stash",
		ConfirmationTitle:          "Confirmation panel",
		PrevPage:                   "Previous page",
		NextPage:                   "Next page",
		GotoTop:                    "Scroll to top",
		GotoBottom:                 "Scroll to bottom",
		FilteringBy:                "Filtering by",
		ResetInParentheses:         "(Reset)",
		OpenFilteringMenu:          "View filter-by-path options",
		FilterBy:                   "Filter by",
		ExitFilterMode:             "Stop filtering by path",
		FilterPathOption:           "Enter path to filter by",
		EnterFileName:              "Enter path:",
		FilteringMenuTitle:         "Filtering",
		MustExitFilterModeTitle:    "Command not available",
		MustExitFilterModePrompt:   "Command not available in filter-by-path mode. Exit filter-by-path mode?",
		Diff:                       "Diff",
		EnterRefToDiff:             "Enter ref to diff",
		EnterRefName:               "Enter ref:",
		ExitDiffMode:               "Exit diff mode",
		DiffingMenuTitle:           "Diffing",
		SwapDiff:                   "Reverse diff direction",
		OpenDiffingMenu:            "Open diff menu",
		// the actual view is the extras view which I intend to give more tabs in future but for now we'll only mention the command log part
		OpenExtrasMenu:                      "Open command log menu",
		ShowingGitDiff:                      "Showing output for:",
		CommitDiff:                          "Commit diff",
		CopyCommitShaToClipboard:            "Copy commit SHA to clipboard",
		CommitSha:                           "Commit SHA",
		CommitURL:                           "Commit URL",
		CopyCommitMessageToClipboard:        "Copy commit message to clipboard",
		CommitMessage:                       "Commit message",
		CommitAuthor:                        "Commit author",
		CopyCommitAttributeToClipboard:      "Copy commit attribute",
		CopyBranchNameToClipboard:           "Copy branch name to clipboard",
		CopyFileNameToClipboard:             "Copy the file name to the clipboard",
		CopyCommitFileNameToClipboard:       "Copy the committed file name to the clipboard",
		CopySelectedTexToClipboard:          "Copy the selected text to the clipboard",
		CommitPrefixPatternError:            "Error in commitPrefix pattern",
		NoFilesStagedTitle:                  "No files staged",
		NoFilesStagedPrompt:                 "You have not staged any files. Commit all files?",
		BranchNotFoundTitle:                 "Branch not found",
		BranchNotFoundPrompt:                "Branch not found. Create a new branch named",
		BranchUnknown:                       "Branch unknown",
		DiscardChangeTitle:                  "Discard change",
		DiscardChangePrompt:                 "Are you sure you want to discard this change (git reset)? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipDiscardChangeWarning' to true",
		CreateNewBranchFromCommit:           "Create new branch off of commit",
		BuildingPatch:                       "Building patch",
		ViewCommits:                         "View commits",
		MinGitVersionError:                  "Git version must be at least 2.20 (i.e. from 2018 onwards). Please upgrade your git version. Alternatively raise an issue at https://github.com/jesseduffield/lazygit/issues for lazygit to be more backwards compatible.",
		RunningCustomCommandStatus:          "Running custom command",
		SubmoduleStashAndReset:              "Stash uncommitted submodule changes and update",
		AndResetSubmodules:                  "And reset submodules",
		EnterSubmodule:                      "Enter submodule",
		CopySubmoduleNameToClipboard:        "Copy submodule name to clipboard",
		RemoveSubmodule:                     "Remove submodule",
		RemoveSubmodulePrompt:               "Are you sure you want to remove submodule '%s' and its corresponding directory? This is irreversible.",
		ResettingSubmoduleStatus:            "Resetting submodule",
		NewSubmoduleName:                    "New submodule name:",
		NewSubmoduleUrl:                     "New submodule URL:",
		NewSubmodulePath:                    "New submodule path:",
		AddSubmodule:                        "Add new submodule",
		AddingSubmoduleStatus:               "Adding submodule",
		UpdateSubmoduleUrl:                  "Update URL for submodule '%s'",
		UpdatingSubmoduleUrlStatus:          "Updating URL",
		EditSubmoduleUrl:                    "Update submodule URL",
		InitializingSubmoduleStatus:         "Initializing submodule",
		InitSubmodule:                       "Initialize submodule",
		SubmoduleUpdate:                     "Update submodule",
		UpdatingSubmoduleStatus:             "Updating submodule",
		BulkInitSubmodules:                  "Bulk init submodules",
		BulkUpdateSubmodules:                "Bulk update submodules",
		BulkDeinitSubmodules:                "Bulk deinit submodules",
		ViewBulkSubmoduleOptions:            "View bulk submodule options",
		BulkSubmoduleOptions:                "Bulk submodule options",
		RunningCommand:                      "Running command",
		SubCommitsTitle:                     "Sub-commits",
		SubmodulesTitle:                     "Submodules",
		NavigationTitle:                     "List panel navigation",
		SuggestionsCheatsheetTitle:          "Suggestions",
		SuggestionsTitle:                    "Suggestions (press %s to focus)",
		ExtrasTitle:                         "Command log",
		PushingTagStatus:                    "Pushing tag",
		PullRequestURLCopiedToClipboard:     "Pull request URL copied to clipboard",
		CommitDiffCopiedToClipboard:         "Commit diff copied to clipboard",
		CommitSHACopiedToClipboard:          "Commit SHA copied to clipboard",
		CommitURLCopiedToClipboard:          "Commit URL copied to clipboard",
		CommitMessageCopiedToClipboard:      "Commit message copied to clipboard",
		CommitAuthorCopiedToClipboard:       "Commit author copied to clipboard",
		PatchCopiedToClipboard:              "Patch copied to clipboard",
		CopiedToClipboard:                   "Copied to clipboard",
		ErrCannotEditDirectory:              "Cannot edit directory: you can only edit individual files",
		ErrStageDirWithInlineMergeConflicts: "Cannot stage/unstage directory containing files with inline merge conflicts. Please fix up the merge conflicts first",
		ErrRepositoryMovedOrDeleted:         "Cannot find repo. It might have been moved or deleted ¯\\_(ツ)_/¯",
		CommandLog:                          "Command log",
		ToggleShowCommandLog:                "Toggle show/hide command log",
		FocusCommandLog:                     "Focus command log",
		CommandLogHeader:                    "You can hide/focus this panel by pressing '%s'\n",
		RandomTip:                           "Random tip",
		SelectParentCommitForMerge:          "Select parent commit for merge",
		ToggleWhitespaceInDiffView:          "Toggle whether or not whitespace changes are shown in the diff view",
		IgnoreWhitespaceDiffViewSubTitle:    "(ignoring whitespace)",
		IgnoreWhitespaceNotSupportedHere:    "Ignoring whitespace is not supported in this view",
		IncreaseContextInDiffView:           "Increase the size of the context shown around changes in the diff view",
		DecreaseContextInDiffView:           "Decrease the size of the context shown around changes in the diff view",
		CreatePullRequestOptions:            "Create pull request options",
		DefaultBranch:                       "Default branch",
		SelectBranch:                        "Select branch",
		SelectConfigFile:                    "Select config file",
		NoConfigFileFoundErr:                "No config file found",
		LoadingFileSuggestions:              "Loading file suggestions",
		LoadingCommits:                      "Loading commits",
		MustSpecifyOriginError:              "Must specify a remote if specifying a branch",
		GitOutput:                           "Git output:",
		GitCommandFailed:                    "Git command failed. Check command log for details (open with %s)",
		AbortTitle:                          "Abort %s",
		AbortPrompt:                         "Are you sure you want to abort the current %s?",
		OpenLogMenu:                         "Open log menu",
		LogMenuTitle:                        "Commit Log Options",
		ToggleShowGitGraphAll:               "Toggle show whole git graph (pass the `--all` flag to `git log`)",
		ShowGitGraph:                        "Show git graph",
		SortCommits:                         "Commit sort order",
		CantChangeContextSizeError:          "Cannot change context while in patch building mode because we were too lazy to support it when releasing the feature. If you really want it, please let us know!",
		OpenCommitInBrowser:                 "Open commit in browser",
		ViewBisectOptions:                   "View bisect options",
		ConfirmRevertCommit:                 "Are you sure you want to revert {{.selectedCommit}}?",
		RewordInEditorTitle:                 "Reword in editor",
		RewordInEditorPrompt:                "Are you sure you want to reword this commit in your editor?",
		HardResetAutostashPrompt:            "Are you sure you want to hard reset to '%s'? An auto-stash will be performed if necessary.",
		CheckoutPrompt:                      "Are you sure you want to checkout '%s'?",
		UpstreamGone:                        "(upstream gone)",
		NukeDescription:                     "If you want to make all the changes in the worktree go away, this is the way to do it. If there are dirty submodule changes this will stash those changes in the submodule(s).",
		DiscardStagedChangesDescription:     "This will create a new stash entry containing only staged files and then drop it, so that the working tree is left with only unstaged changes",
		EmptyOutput:                         "<Empty output>",
		Patch:                               "Patch",
		CustomPatch:                         "Custom patch",
		CommitsCopied:                       "commits copied", // lowercase because it's used in a sentence
		CommitCopied:                        "commit copied",  // lowercase because it's used in a sentence
		ResetPatch:                          "Reset patch",
		ApplyPatch:                          "Apply patch",
		ApplyPatchInReverse:                 "Apply patch in reverse",
		RemovePatchFromOriginalCommit:       "Remove patch from original commit (%s)",
		MovePatchOutIntoIndex:               "Move patch out into index",
		MovePatchIntoNewCommit:              "Move patch into new commit",
		MovePatchToSelectedCommit:           "Move patch to selected commit (%s)",
		CopyPatchToClipboard:                "Copy patch to clipboard",
		NoMatchesFor:                        "No matches for '%s' %s",
		ExitSearchMode:                      "%s: Exit search mode",
		ExitTextFilterMode:                  "%s: Exit filter mode",
		MatchesFor:                          "matches for '%s' (%d of %d) %s", // lowercase because it's after other text
		SearchKeybindings:                   "%s: Next match, %s: Previous match, %s: Exit search mode",
		SearchPrefix:                        "Search: ",
		FilterPrefix:                        "Filter: ",
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
			SetUnsetUpstream:                  "Set/Unset upstream",
			CreateBranch:                      "Create branch",
			CherryPick:                        "(Cherry-pick) paste commits",
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
			IgnoreExcludeFile:                 "Ignore or exclude file",
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
			GitFlowFinish:                     "git flow finish",
			GitFlowStart:                      "git flow start",
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
			Mark:                        "Mark %s as %s",
			MarkStart:                   "Mark %s as %s (start bisect)",
			Skip:                        "Skip %s",
			ResetTitle:                  "Reset 'git bisect'",
			ResetPrompt:                 "Are you sure you want to reset 'git bisect'?",
			ResetOption:                 "Reset bisect",
			BisectMenuTitle:             "Bisect",
			CompleteTitle:               "Bisect complete",
			CompletePrompt:              "Bisect complete! The following commit introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			CompletePromptIndeterminate: "Bisect complete! Some commits were skipped, so any of the following commits may have introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			Bisecting:                   "Bisecting",
		},
	}
}
