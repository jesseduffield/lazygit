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
	LogTitle                            string
	FilesTitle                          string
	BranchesTitle                       string
	CommitsTitle                        string
	StashTitle                          string
	UnstagedChanges                     string
	StagedChanges                       string
	PatchBuildingMainTitle              string
	MergingMainTitle                    string
	MainTitle                           string
	StagingTitle                        string
	MergingTitle                        string
	NormalTitle                         string
	CommitMessage                       string
	CredentialsUsername                 string
	CredentialsPassword                 string
	CredentialsPassphrase               string
	PassUnameWrong                      string
	CommitChanges                       string
	AmendLastCommit                     string
	SureToAmend                         string
	NoCommitToAmend                     string
	CommitChangesWithEditor             string
	StatusTitle                         string
	GlobalTitle                         string
	LcNavigate                          string
	LcMenu                              string
	LcExecute                           string
	LcOpen                              string
	LcIgnore                            string
	LcDelete                            string
	LcToggleStaged                      string
	LcToggleStagedAll                   string
	LcToggleTreeView                    string
	LcOpenMergeTool                     string
	LcRefresh                           string
	LcPush                              string
	LcPull                              string
	LcEdit                              string
	LcScroll                            string
	LcAbortMerge                        string
	LcResolveMergeConflicts             string
	MergeConflictsTitle                 string
	LcCheckout                          string
	NoChangedFiles                      string
	FileHasNoUnstagedChanges            string
	CannotGitAdd                        string
	NoFilesDisplay                      string
	NotAFile                            string
	PullWait                            string
	PushWait                            string
	FetchWait                           string
	FileNoMergeCons                     string
	LcSoftReset                         string
	SureTo                              string
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
	LcMerge                             string
	LcCheckoutByName                    string
	LcNewBranch                         string
	LcDeleteBranch                      string
	LcForceDeleteBranch                 string
	NoBranchesThisRepo                  string
	NoTrackingThisBranch                string
	CommitMessageConfirm                string
	CommitWithoutMessageErr             string
	CloseConfirm                        string
	LcClose                             string
	LcQuit                              string
	SureResetThisCommit                 string
	ResetToCommit                       string
	LcSquashDown                        string
	LcRename                            string
	LcResetToThisCommit                 string
	LcFixupCommit                       string
	OnlySquashTopmostCommit             string
	YouNoCommitsToSquash                string
	CantFixupWhileUnstagedChanges       string
	Fixup                               string
	SureFixupThisCommit                 string
	SureSquashThisCommit                string
	Squash                              string
	LcPickCommit                        string
	LcRevertCommit                      string
	OnlyRenameTopCommit                 string
	LcRenameCommit                      string
	LcDeleteCommit                      string
	LcMoveDownCommit                    string
	LcMoveUpCommit                      string
	LcEditCommit                        string
	LcAmendToCommit                     string
	LcRenameCommitEditor                string
	PotentialErrInGetselectedCommit     string
	NoCommitsThisBranch                 string
	Error                               string
	RunningSubprocess                   string
	LcSelectHunk                        string
	LcNavigateConflicts                 string
	LcPickHunk                          string
	LcPickBothHunks                     string
	LcUndo                              string
	LcUndoReflog                        string
	LcRedoReflog                        string
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
	NoStashTo                           string
	NoTrackedStagedFilesStash           string
	StashChanges                        string
	IssntListOfViews                    string
	LcNewFocusedViewIs                  string
	MergeAborted                        string
	OpenConfig                          string
	EditConfig                          string
	ForcePush                           string
	ForcePushPrompt                     string
	ForcePushDisabled                   string
	UpdatesRejectedAndForcePushDisabled string
	LcCheckForUpdate                    string
	CheckingForUpdates                  string
	OnLatestVersionErr                  string
	MajorVersionErr                     string
	CouldNotFindBinaryErr               string
	AnonymousReportingTitle             string
	AnonymousReportingPrompt            string
	MergeToolTitle                      string
	MergeToolPrompt                     string
	IntroPopupMessage                   string
	GitconfigParseErr                   string
	LcEditFile                          string
	LcOpenFile                          string
	LcIgnoreFile                        string
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
	SelectHunk                          string
	StageSelection                      string
	ResetSelection                      string
	ToggleDragSelect                    string
	ToggleSelectHunk                    string
	ToggleSelectionForPatch             string
	TogglePanel                         string
	CantStageStaged                     string
	ReturnToFilesPanel                  string
	CantFindHunks                       string
	CantFindHunk                        string
	FastForward                         string
	Fetching                            string
	FoundConflicts                      string
	FoundConflictsTitle                 string
	Undo                                string
	PickHunk                            string
	PickBothHunks                       string
	ViewMergeRebaseOptions              string
	NotMergingOrRebasing                string
	RecentRepos                         string
	MergeOptionsTitle                   string
	RebaseOptionsTitle                  string
	CommitMessageTitle                  string
	LocalBranchesTitle                  string
	SearchTitle                         string
	TagsTitle                           string
	BranchCommitsTitle                  string
	MenuTitle                           string
	RemotesTitle                        string
	CredentialsTitle                    string
	RemoteBranchesTitle                 string
	PatchBuildingTitle                  string
	InformationTitle                    string
	SecondaryTitle                      string
	ReflogCommitsTitle                  string
	Title                               string
	ConflictsResolved                   string
	RebasingTitle                       string
	ConfirmRebase                       string
	ConfirmMerge                        string
	FwdNoUpstream                       string
	FwdCommitsToPush                    string
	ErrorOccurred                       string
	NoRoom                              string
	YouAreHere                          string
	LcRewordNotSupported                string
	LcCherryPickCopy                    string
	LcCherryPickCopyRange               string
	LcPasteCommits                      string
	SureCherryPick                      string
	CherryPick                          string
	CannotRebaseOntoFirstCommit         string
	CannotSquashOntoSecondCommit        string
	Donate                              string
	AskQuestion                         string
	PrevLine                            string
	NextLine                            string
	PrevHunk                            string
	NextHunk                            string
	PrevConflict                        string
	NextConflict                        string
	SelectTop                           string
	SelectBottom                        string
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
	LcViewCommitFiles                   string
	CommitFilesTitle                    string
	LcGoBack                            string
	NoCommiteFiles                      string
	LcCheckoutCommitFile                string
	LcDiscardOldFileChange              string
	DiscardFileChangesTitle             string
	DiscardFileChangesPrompt            string
	DisabledForGPG                      string
	CreateRepo                          string
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
	LcHardReset                         string
	LcHardResetUpstream                 string
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
	LcStashOptions                      string
	NotARepository                      string
	LcJump                              string
	DiscardPatch                        string
	DiscardPatchConfirm                 string
	CantPatchWhileRebasingError         string
	LcToggleAddToPatch                  string
	ViewPatchOptions                    string
	PatchOptionsTitle                   string
	NoPatchError                        string
	LcEnterFile                         string
	ExitLineByLineMode                  string
	EnterUpstream                       string
	EnterUpstreamWithSlash              string
	LcNotTrackingRemote                 string
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
	LcSetUpstream                       string
	SetUpstreamTitle                    string
	SetUpstreamMessage                  string
	LcEditRemote                        string
	LcTagCommit                         string
	TagNameTitle                        string
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
	IgnoreTrackedPrompt                 string
	LcViewResetToUpstreamOptions        string
	LcNextScreenMode                    string
	LcPrevScreenMode                    string
	LcStartSearch                       string
	Panel                               string
	Keybindings                         string
	LcRenameBranch                      string
	NewGitFlowBranchPrompt              string
	RenameBranchWarning                 string
	LcOpenMenu                          string
	LcCloseMenu                         string
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
	LcEnterFileName                     string
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
	LcCopyCommitShaToClipboard          string
	LcCopyCommitMessageToClipboard      string
	LcCopyBranchNameToClipboard         string
	LcCopyFileNameToClipboard           string
	LcCopyCommitFileNameToClipboard     string
	LcCommitPrefixPatternError          string
	NoFilesStagedTitle                  string
	NoFilesStagedPrompt                 string
	BranchNotFoundTitle                 string
	BranchNotFoundPrompt                string
	UnstageLinesTitle                   string
	UnstageLinesPrompt                  string
	LcCreateNewBranchFromCommit         string
	LcViewStashFiles                    string
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
	LcViewResetAndRemoveOptions         string
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
	SuggestionsTitle                    string
	PushingTagStatus                    string
	PullRequestURLCopiedToClipboard     string
	CommitMessageCopiedToClipboard      string
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
	Spans                               Spans
}

type Spans struct {
	CheckoutCommit                    string
	CheckoutReflogCommit              string
	CheckoutTag                       string
	CheckoutBranch                    string
	ForceCheckoutBranch               string
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
	RevertCommit                      string
	CreateFixupCommit                 string
	SquashAllAboveFixupCommits        string
	MoveCommitUp                      string
	MoveCommitDown                    string
	CopyCommitMessageToClipboard      string
	CustomCommand                     string
	DiscardAllChangesInDirectory      string
	DiscardUnstagedChangesInDirectory string
	DiscardAllChangesInFile           string
	DiscardAllUnstagedChangesInFile   string
	StageFile                         string
	UnstageFile                       string
	UnstageAllFiles                   string
	StageAllFiles                     string
	IgnoreFile                        string
	Commit                            string
	EditFile                          string
	Push                              string
	Pull                              string
	OpenFile                          string
	StashAllChanges                   string
	StashStagedChanges                string
	GitFlowFinish                     string
	GitFlowStart                      string
	CopyToClipboard                   string
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
	RemoveSubmodule                   string
	ResetSubmodule                    string
	AddSubmodule                      string
	UpdateSubmoduleUrl                string
	InitialiseSubmodule               string
	BulkInitialiseSubmodules          string
	BulkUpdateSubmodules              string
	BulkStashAndResetSubmodules       string
	BulkDeinitialiseSubmodules        string
	UpdateSubmodule                   string
	CreateLightweightTag              string
	DeleteTag                         string
	PushTag                           string
	NukeWorkingTree                   string
	DiscardUnstagedFileChanges        string
	RemoveUntrackedFiles              string
	SoftReset                         string
	MixedReset                        string
	HardReset                         string
	Undo                              string
	Redo                              string
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
    Or even just star the repo cos we're not far from 20k stars!
`

func EnglishTranslationSet() *TranslationSet {
	return &TranslationSet{
		NotEnoughSpace:                      "Not enough space to render panels",
		DiffTitle:                           "Diff",
		LogTitle:                            "Log",
		FilesTitle:                          "Files",
		BranchesTitle:                       "Branches",
		CommitsTitle:                        "Commits",
		StashTitle:                          "Stash",
		UnstagedChanges:                     `Unstaged Changes`,
		StagedChanges:                       `Staged Changes`,
		PatchBuildingMainTitle:              `Add Lines/Hunks To Patch`,
		MergingMainTitle:                    "Resolve merge conflicts",
		MainTitle:                           "Main",
		StagingTitle:                        "Staging",
		MergingTitle:                        "Merging",
		NormalTitle:                         "Normal",
		CommitMessage:                       "Commit message",
		CredentialsUsername:                 "Username",
		CredentialsPassword:                 "Password",
		CredentialsPassphrase:               "Enter passphrase for SSH key",
		PassUnameWrong:                      "Password, passphrase and/or username wrong",
		CommitChanges:                       "commit changes",
		AmendLastCommit:                     "amend last commit",
		SureToAmend:                         "Are you sure you want to amend last commit? Afterwards, you can change commit message from the commits panel.",
		NoCommitToAmend:                     "There's no commit to amend.",
		CommitChangesWithEditor:             "commit changes using git editor",
		StatusTitle:                         "Status",
		LcNavigate:                          "navigate",
		LcMenu:                              "menu",
		LcExecute:                           "execute",
		LcOpen:                              "open",
		LcIgnore:                            "ignore",
		LcDelete:                            "delete",
		LcToggleStaged:                      "toggle staged",
		LcToggleStagedAll:                   "stage/unstage all",
		LcToggleTreeView:                    "toggle file tree view",
		LcOpenMergeTool:                     "open external merge tool (git mergetool)",
		LcRefresh:                           "refresh",
		LcPush:                              "push",
		LcPull:                              "pull",
		LcEdit:                              "edit",
		LcScroll:                            "scroll",
		LcAbortMerge:                        "abort merge",
		LcResolveMergeConflicts:             "resolve merge conflicts",
		MergeConflictsTitle:                 "Merge Conflicts",
		LcCheckout:                          "checkout",
		NoChangedFiles:                      "No changed files",
		FileHasNoUnstagedChanges:            "File has no unstaged changes to add",
		CannotGitAdd:                        "Cannot git add --patch untracked files",
		NoFilesDisplay:                      "No file to display",
		NotAFile:                            "Not a file",
		PullWait:                            "Pulling...",
		PushWait:                            "Pushing...",
		FetchWait:                           "Fetching...",
		FileNoMergeCons:                     "This file has no inline merge conflicts",
		LcSoftReset:                         "soft reset",
		SureTo:                              "Are you sure you want to {{.deleteVerb}} {{.fileName}} (you will lose your changes)?",
		AlreadyCheckedOutBranch:             "You have already checked out this branch",
		SureForceCheckout:                   "Are you sure you want force checkout? You will lose all local changes",
		ForceCheckoutBranch:                 "Force Checkout Branch",
		BranchName:                          "Branch name",
		NewBranchNameBranchOff:              "New Branch Name (Branch is off of {{.branchName}})",
		CantDeleteCheckOutBranch:            "You cannot delete the checked out branch!",
		DeleteBranch:                        "Delete Branch",
		DeleteBranchMessage:                 "Are you sure you want to delete the branch {{.selectedBranchName}}?",
		ForceDeleteBranchMessage:            "{{.selectedBranchName}} is not fully merged. Are you sure you want to delete it?",
		LcRebaseBranch:                      "rebase checked-out branch onto this branch",
		CantRebaseOntoSelf:                  "You cannot rebase a branch onto itself",
		CantMergeBranchIntoItself:           "You cannot merge a branch into itself",
		LcForceCheckout:                     "force checkout",
		LcMerge:                             "merge",
		LcCheckoutByName:                    "checkout by name",
		LcNewBranch:                         "new branch",
		LcDeleteBranch:                      "delete branch",
		LcForceDeleteBranch:                 "delete branch (force)",
		NoBranchesThisRepo:                  "No branches for this repo",
		NoTrackingThisBranch:                "There is no tracking for this branch",
		CommitMessageConfirm:                "{{.keyBindClose}}: close, {{.keyBindNewLine}}: new line, {{.keyBindConfirm}}: confirm",
		CommitWithoutMessageErr:             "You cannot commit without a commit message",
		CloseConfirm:                        "{{.keyBindClose}}: close, {{.keyBindConfirm}}: confirm",
		LcClose:                             "close",
		LcQuit:                              "quit",
		SureResetThisCommit:                 "Are you sure you want to reset to this commit?",
		ResetToCommit:                       "Reset To Commit",
		LcSquashDown:                        "squash down",
		LcRename:                            "rename",
		LcResetToThisCommit:                 "reset to this commit",
		LcFixupCommit:                       "fixup commit",
		NoCommitsThisBranch:                 "No commits for this branch",
		OnlySquashTopmostCommit:             "Can only squash topmost commit",
		YouNoCommitsToSquash:                "You have no commits to squash with",
		CantFixupWhileUnstagedChanges:       "Can't fixup while there are unstaged changes",
		Fixup:                               "Fixup",
		SureFixupThisCommit:                 "Are you sure you want to 'fixup' this commit? It will be merged into the commit below",
		SureSquashThisCommit:                "Are you sure you want to squash this commit into the commit below?",
		Squash:                              "Squash",
		LcPickCommit:                        "pick commit (when mid-rebase)",
		LcRevertCommit:                      "revert commit",
		OnlyRenameTopCommit:                 "Can only reword topmost commit from within lazygit. Use shift+R instead",
		LcRenameCommit:                      "reword commit",
		LcDeleteCommit:                      "delete commit",
		LcMoveDownCommit:                    "move commit down one",
		LcMoveUpCommit:                      "move commit up one",
		LcEditCommit:                        "edit commit",
		LcAmendToCommit:                     "amend commit with staged changes",
		LcRenameCommitEditor:                "rename commit with editor",
		PotentialErrInGetselectedCommit:     "potential error in getSelected Commit (mismatched ui and state)",
		Error:                               "Error",
		RunningSubprocess:                   "running subprocess",
		LcSelectHunk:                        "select hunk",
		LcNavigateConflicts:                 "navigate conflicts",
		LcPickHunk:                          "pick hunk",
		LcPickBothHunks:                     "pick both hunks",
		LcUndo:                              "undo",
		LcUndoReflog:                        "undo (via reflog) (experimental)",
		LcRedoReflog:                        "redo (via reflog) (experimental)",
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
		NoStashTo:                           "No stash to {{.method}}",
		NoTrackedStagedFilesStash:           "You have no tracked/staged files to stash",
		StashChanges:                        "Stash changes",
		IssntListOfViews:                    "{{.name}} is not in the list of views",
		LcNewFocusedViewIs:                  "new focused view is {{.newFocusedView}}",
		MergeAborted:                        "Merge aborted",
		OpenConfig:                          "open config file",
		EditConfig:                          "edit config file",
		ForcePush:                           "Force push",
		ForcePushPrompt:                     "Your branch has diverged from the remote branch. Press 'esc' to cancel, or 'enter' to force push.",
		ForcePushDisabled:                   "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejectedAndForcePushDisabled: "Updates were rejected and you have disabled force pushing",
		LcCheckForUpdate:                    "check for update",
		CheckingForUpdates:                  "Checking for updates...",
		OnLatestVersionErr:                  "You already have the latest version",
		MajorVersionErr:                     "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		CouldNotFindBinaryErr:               "Could not find any binary at {{.url}}",
		AnonymousReportingTitle:             "Help make lazygit better",
		AnonymousReportingPrompt:            "Would you like to enable anonymous reporting data to help improve lazygit? (enter/esc)",
		MergeToolTitle:                      "Merge tool",
		MergeToolPrompt:                     "Are you sure you want to open `git mergetool`?",
		IntroPopupMessage:                   englishIntroPopupMessage,
		GitconfigParseErr:                   `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		LcEditFile:                          `edit file`,
		LcOpenFile:                          `open file`,
		LcIgnoreFile:                        `add to .gitignore`,
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
		SelectHunk:                          `select hunk`,
		StageSelection:                      `toggle line staged / unstaged`,
		ResetSelection:                      `delete change (git reset)`,
		ToggleDragSelect:                    `toggle drag select`,
		ToggleSelectHunk:                    `toggle select hunk`,
		ToggleSelectionForPatch:             `add/remove line(s) to patch`,
		TogglePanel:                         `switch to other panel`,
		CantStageStaged:                     `You can't stage an already staged change!`,
		ReturnToFilesPanel:                  `return to files panel`,
		CantFindHunks:                       `Could not find any hunks in this patch`,
		CantFindHunk:                        `Could not find hunk`,
		FastForward:                         `fast-forward this branch from its upstream`,
		Fetching:                            "fetching and fast-forwarding {{.from}} -> {{.to}} ...",
		FoundConflicts:                      "Conflicts! To abort press 'esc', otherwise press 'enter'",
		FoundConflictsTitle:                 "Auto-merge failed",
		Undo:                                "undo",
		PickHunk:                            "pick hunk",
		PickBothHunks:                       "pick both hunks",
		ViewMergeRebaseOptions:              "view merge/rebase options",
		NotMergingOrRebasing:                "You are currently neither rebasing nor merging",
		RecentRepos:                         "recent repositories",
		MergeOptionsTitle:                   "Merge Options",
		RebaseOptionsTitle:                  "Rebase Options",
		CommitMessageTitle:                  "Commit Message",
		LocalBranchesTitle:                  "Branches Tab",
		SearchTitle:                         "Search",
		TagsTitle:                           "Tags Tab",
		BranchCommitsTitle:                  "Commits Tab",
		MenuTitle:                           "Menu",
		RemotesTitle:                        "Remotes Tab",
		CredentialsTitle:                    "Credentials",
		RemoteBranchesTitle:                 "Remote Branches (in Remotes tab)",
		PatchBuildingTitle:                  "Patch Building",
		InformationTitle:                    "Information",
		SecondaryTitle:                      "Secondary",
		ReflogCommitsTitle:                  "Reflog Tab",
		Title:                               "Title",
		GlobalTitle:                         "Global Keybindings",
		ConflictsResolved:                   "all merge conflicts resolved. Continue?",
		RebasingTitle:                       "Rebasing",
		ConfirmRebase:                       "Are you sure you want to rebase {{.checkedOutBranch}} onto {{.selectedBranch}}?",
		ConfirmMerge:                        "Are you sure you want to merge {{.selectedBranch}} into {{.checkedOutBranch}}?",
		FwdNoUpstream:                       "Cannot fast-forward a branch with no upstream",
		FwdCommitsToPush:                    "Cannot fast-forward a branch with commits to push",
		ErrorOccurred:                       "An error occurred! Please create an issue at",
		NoRoom:                              "Not enough room",
		YouAreHere:                          "YOU ARE HERE",
		LcRewordNotSupported:                "rewording commits while interactively rebasing is not currently supported",
		LcCherryPickCopy:                    "copy commit (cherry-pick)",
		LcCherryPickCopyRange:               "copy commit range (cherry-pick)",
		LcPasteCommits:                      "paste commits (cherry-pick)",
		SureCherryPick:                      "Are you sure you want to cherry-pick the copied commits onto this branch?",
		CherryPick:                          "Cherry-Pick",
		CannotRebaseOntoFirstCommit:         "You cannot interactive rebase onto the first commit",
		CannotSquashOntoSecondCommit:        "You cannot squash/fixup onto the second commit",
		Donate:                              "Donate",
		AskQuestion:                         "Ask Question",
		PrevLine:                            "select previous line",
		NextLine:                            "select next line",
		PrevHunk:                            "select previous hunk",
		NextHunk:                            "select next hunk",
		PrevConflict:                        "select previous conflict",
		NextConflict:                        "select next conflict",
		SelectTop:                           "select top hunk",
		SelectBottom:                        "select bottom hunk",
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
		LcViewCommitFiles:                   "view commit's files",
		CommitFilesTitle:                    "Commit Files",
		LcGoBack:                            "go back",
		NoCommiteFiles:                      "No files for this commit",
		LcCheckoutCommitFile:                "checkout file",
		LcDiscardOldFileChange:              "discard this commit's changes to this file",
		DiscardFileChangesTitle:             "Discard file changes",
		DiscardFileChangesPrompt:            "Are you sure you want to discard this commit's changes to this file? If this file was created in this commit, it will be deleted",
		DisabledForGPG:                      "Feature not available for users using GPG",
		CreateRepo:                          "Not in a git repository. Create a new git repository? (y/n): ",
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
		LcHardReset:                         "hard reset",
		LcHardResetUpstream:                 "hard reset to upstream branch",
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
		LcStashAllChanges:                   "stash changes",
		LcStashStagedChanges:                "stash staged changes",
		LcStashOptions:                      "Stash options",
		NotARepository:                      "Error: must be run inside a git repository",
		LcJump:                              "jump to panel",
		DiscardPatch:                        "Discard Patch",
		DiscardPatchConfirm:                 "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		CantPatchWhileRebasingError:         "You cannot build a patch or run patch commands while in a merging or rebasing state",
		LcToggleAddToPatch:                  "toggle file included in patch",
		ViewPatchOptions:                    "view custom patch options",
		PatchOptionsTitle:                   "Patch Options",
		NoPatchError:                        "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		LcEnterFile:                         "enter file to add selected lines to the patch (or toggle directory collapsed)",
		ExitLineByLineMode:                  `exit line-by-line mode`,
		EnterUpstream:                       `Enter upstream as '<remote> <branchname>'`,
		EnterUpstreamWithSlash:              `Enter upstream as '<remote>/<branchname>'`,
		LcNotTrackingRemote:                 "(not tracking any remote)",
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
		LcSetUpstream:                       "set as upstream of checked-out branch",
		SetUpstreamTitle:                    "Set upstream branch",
		SetUpstreamMessage:                  "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		LcEditRemote:                        "edit remote",
		LcTagCommit:                         "tag commit",
		TagNameTitle:                        "Tag name:",
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
		LcViewResetToUpstreamOptions:        "view upstream reset options",
		LcNextScreenMode:                    "next screen mode (normal/half/fullscreen)",
		LcPrevScreenMode:                    "prev screen mode",
		LcStartSearch:                       "start search",
		Panel:                               "Panel",
		Keybindings:                         "Keybindings",
		LcRenameBranch:                      "rename branch",
		NewBranchNamePrompt:                 "Enter new branch name for branch",
		RenameBranchWarning:                 "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		LcOpenMenu:                          "open menu",
		LcCloseMenu:                         "close menu",
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
		LcEnterFileName:                     "enter path:",
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
		LcCopyCommitShaToClipboard:          "copy commit SHA to clipboard",
		LcCopyCommitMessageToClipboard:      "copy commit message to clipboard",
		LcCopyBranchNameToClipboard:         "copy branch name to clipboard",
		LcCopyFileNameToClipboard:           "copy the file name to the clipboard",
		LcCopyCommitFileNameToClipboard:     "copy the committed file name to the clipboard",
		LcCommitPrefixPatternError:          "Error in commitPrefix pattern",
		NoFilesStagedTitle:                  "No files staged",
		NoFilesStagedPrompt:                 "You have not staged any files. Commit all files?",
		BranchNotFoundTitle:                 "Branch not found",
		BranchNotFoundPrompt:                "Branch not found. Create a new branch named",
		UnstageLinesTitle:                   "Unstage lines",
		UnstageLinesPrompt:                  "Are you sure you want to delete the selected lines (git reset)? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipUnstageLineWarning' to true",
		LcCreateNewBranchFromCommit:         "create new branch off of commit",
		LcViewStashFiles:                    "view stash entry's files",
		LcBuildingPatch:                     "building patch",
		LcViewCommits:                       "view commits",
		MinGitVersionError:                  "Git version must be at least 2.0 (i.e. from 2014 onwards). Please upgrade your git version. Alternatively raise an issue at https://github.com/jesseduffield/lazygit/issues for lazygit to be more backwards compatible.",
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
		LcViewResetAndRemoveOptions:         "view reset and remove submodule options",
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
		SuggestionsTitle:                    "Suggestions",
		PushingTagStatus:                    "pushing tag",
		PullRequestURLCopiedToClipboard:     "Pull request URL copied to clipboard",
		CommitMessageCopiedToClipboard:      "Commit message copied to clipboard",
		LcCopiedToClipboard:                 "copied to clipboard",
		ErrCannotEditDirectory:              "Cannot edit directory: you can only edit individual files",
		ErrStageDirWithInlineMergeConflicts: "Cannot stage/unstage directory containing files with inline merge conflicts. Please fix up the merge conflicts first",
		ErrRepositoryMovedOrDeleted:         "Cannot find repo. It might have been moved or deleted ¯\\_(ツ)_/¯",
		CommandLog:                          "Command Log",
		ToggleShowCommandLog:                "Toggle show/hide command log",
		FocusCommandLog:                     "Focus command log",
		CommandLogHeader:                    "You can hide/focus this panel by pressing '%s' or hide it permanently in your config with `gui.showCommandLog: false`\n",
		RandomTip:                           "Random Tip",
		SelectParentCommitForMerge:          "Select parent commit for merge",
		Spans: Spans{
			// TODO: combine this with the original keybinding descriptions (those are all in lowercase atm)
			CheckoutCommit:                    "Checkout commit",
			CheckoutReflogCommit:              "Checkout reflog commit",
			CheckoutTag:                       "Checkout tag",
			CheckoutBranch:                    "Checkout branch",
			ForceCheckoutBranch:               "Force checkout branch",
			DeleteBranch:                      "Delete branch",
			Merge:                             "Merge",
			RebaseBranch:                      "Rebase branch",
			RenameBranch:                      "Rename branch",
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
			RevertCommit:                      "Revert commit",
			CreateFixupCommit:                 "Create fixup commit",
			SquashAllAboveFixupCommits:        "Squash all above fixup commits",
			CreateLightweightTag:              "Create lightweight tag",
			CopyCommitMessageToClipboard:      "Copy commit message to clipboard",
			MoveCommitUp:                      "Move commit up",
			MoveCommitDown:                    "Move commit down",
			CustomCommand:                     "Custom command",
			DiscardAllChangesInDirectory:      "Discard all changes in directory",
			DiscardUnstagedChangesInDirectory: "Discard unstaged changes in directory",
			DiscardAllChangesInFile:           "Discard all changes in file",
			DiscardAllUnstagedChangesInFile:   "Discard all unstaged changes in file",
			StageFile:                         "Stage file",
			UnstageFile:                       "Unstage file",
			UnstageAllFiles:                   "Unstage all files",
			StageAllFiles:                     "Stage all files",
			IgnoreFile:                        "Ignore file",
			Commit:                            "Commit",
			EditFile:                          "Edit file",
			Push:                              "Push",
			Pull:                              "Pull",
			OpenFile:                          "Open file",
			StashAllChanges:                   "Stash all changes",
			StashStagedChanges:                "Stash staged changes",
			GitFlowFinish:                     "Git flow finish",
			GitFlowStart:                      "Git Flow start",
			CopyToClipboard:                   "Copy to clipboard",
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
			RemoveSubmodule:                   "Remove submodule",
			ResetSubmodule:                    "Reset submodule",
			AddSubmodule:                      "Add submodule",
			UpdateSubmoduleUrl:                "Update submodule URL",
			InitialiseSubmodule:               "Initialise submodule",
			BulkInitialiseSubmodules:          "Bulk initialise submodules",
			BulkUpdateSubmodules:              "Bulk update submodules",
			BulkStashAndResetSubmodules:       "Bulk stash and reset submodules",
			BulkDeinitialiseSubmodules:        "Bulk deinitialise submodules",
			UpdateSubmodule:                   "Update submodule",
			DeleteTag:                         "Delete tag",
			PushTag:                           "Push tag",
			NukeWorkingTree:                   "Nuke working tree",
			DiscardUnstagedFileChanges:        "Discard unstaged file changes",
			RemoveUntrackedFiles:              "Remove untracked files",
			SoftReset:                         "Soft reset",
			MixedReset:                        "Mixed reset",
			HardReset:                         "Hard reset",
			FastForwardBranch:                 "Fast forward branch",
			Undo:                              "Undo",
			Redo:                              "Redo",
		},
	}
}
