/*

Todo list when making a new translation
- Copy this file and rename it to the language you want to translate to like someLanguage.go
- Change the addEnglish() name to the language you want to translate to like addSomeLanguage()
- change the first function argument of i18nObject.AddMessages( to the language you want to translate to like language.SomeLanguage
- Remove this todo and the about section

*/

package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func addEnglish(i18nObject *i18n.Bundle) error {

	return i18nObject.AddMessages(language.English,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Not enough space to render panels",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Diff",
		}, &i18n.Message{
			ID:    "LogTitle",
			Other: "Log",
		}, &i18n.Message{
			ID:    "FilesTitle",
			Other: "Files",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Branches",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Commits",
		}, &i18n.Message{
			ID:    "CommitsDiffTitle",
			Other: "Commits (specific diff mode)",
		}, &i18n.Message{
			ID:    "CommitsDiff",
			Other: "select commit to diff with another commit",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Stash",
		}, &i18n.Message{
			ID:    "UnstagedChanges",
			Other: `Unstaged Changes`,
		}, &i18n.Message{
			ID:    "StagedChanges",
			Other: `Staged Changes`,
		}, &i18n.Message{
			ID:    "PatchBuildingMainTitle",
			Other: `Add Lines/Hunks To Patch`,
		}, &i18n.Message{
			ID:    "MergingMainTitle",
			Other: "Resolve merge conflicts",
		}, &i18n.Message{
			ID:    "MainTitle",
			Other: "Main",
		}, &i18n.Message{
			ID:    "StagingTitle",
			Other: "Staging",
		}, &i18n.Message{
			ID:    "MergingTitle",
			Other: "Merging",
		}, &i18n.Message{
			ID:    "NormalTitle",
			Other: "Normal",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Commit message",
		}, &i18n.Message{
			ID:    "CredentialsUsername",
			Other: "Username",
		}, &i18n.Message{
			ID:    "CredentialsPassword",
			Other: "Password",
		}, &i18n.Message{
			ID:    "PassUnameWrong",
			Other: "Password and/or username wrong",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "commit changes",
		}, &i18n.Message{
			ID:    "AmendLastCommit",
			Other: "amend last commit",
		}, &i18n.Message{
			ID:    "SureToAmend",
			Other: "Are you sure you want to amend last commit? Afterwards, you can change commit message from the commits panel.",
		}, &i18n.Message{
			ID:    "NoCommitToAmend",
			Other: "There's no commit to amend.",
		}, &i18n.Message{
			ID:    "CommitChangesWithEditor",
			Other: "commit changes using git editor",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "GlobalTitle",
			Other: "Global",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "navigate",
		}, &i18n.Message{
			ID:    "menu",
			Other: "menu",
		}, &i18n.Message{
			ID:    "execute",
			Other: "execute",
		}, &i18n.Message{
			ID:    "open",
			Other: "open",
		}, &i18n.Message{
			ID:    "ignore",
			Other: "ignore",
		}, &i18n.Message{
			ID:    "delete",
			Other: "delete",
		}, &i18n.Message{
			ID:    "toggleStaged",
			Other: "toggle staged",
		}, &i18n.Message{
			ID:    "toggleStagedAll",
			Other: "stage/unstage all",
		}, &i18n.Message{
			ID:    "refresh",
			Other: "refresh",
		}, &i18n.Message{
			ID:    "push",
			Other: "push",
		}, &i18n.Message{
			ID:    "pull",
			Other: "pull",
		}, &i18n.Message{
			ID:    "edit",
			Other: "edit",
		}, &i18n.Message{
			ID:    "scroll",
			Other: "scroll",
		}, &i18n.Message{
			ID:    "abortMerge",
			Other: "abort merge",
		}, &i18n.Message{
			ID:    "resolveMergeConflicts",
			Other: "resolve merge conflicts",
		}, &i18n.Message{
			ID:    "MergeConflictsTitle",
			Other: "Merge Conflicts",
		}, &i18n.Message{
			ID:    "checkout",
			Other: "checkout",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "No changed files",
		}, &i18n.Message{
			ID:    "FileHasNoUnstagedChanges",
			Other: "File has no unstaged changes to add",
		}, &i18n.Message{
			ID:    "CannotGitAdd",
			Other: "Cannot git add --patch untracked files",
		}, &i18n.Message{
			ID:    "NoStagedFilesToCommit",
			Other: "There are no staged files to commit",
		}, &i18n.Message{
			ID:    "NoFilesDisplay",
			Other: "No file to display",
		}, &i18n.Message{
			ID:    "NotAFile",
			Other: "Not a file",
		}, &i18n.Message{
			ID:    "PullWait",
			Other: "Pulling...",
		}, &i18n.Message{
			ID:    "PushWait",
			Other: "Pushing...",
		}, &i18n.Message{
			ID:    "FetchWait",
			Other: "Fetching...",
		}, &i18n.Message{
			ID:    "FileNoMergeCons",
			Other: "This file has no inline merge conflicts",
		}, &i18n.Message{
			ID:    "softReset",
			Other: "soft reset",
		}, &i18n.Message{
			ID:    "SureTo",
			Other: "Are you sure you want to {{.deleteVerb}} {{.fileName}} (you will lose your changes)?",
		}, &i18n.Message{
			ID:    "AlreadyCheckedOutBranch",
			Other: "You have already checked out this branch",
		}, &i18n.Message{
			ID:    "SureForceCheckout",
			Other: "Are you sure you want force checkout? You will lose all local changes",
		}, &i18n.Message{
			ID:    "ForceCheckoutBranch",
			Other: "Force Checkout Branch",
		}, &i18n.Message{
			ID:    "BranchName",
			Other: "Branch name",
		}, &i18n.Message{
			ID:    "NewBranchNameBranchOff",
			Other: "New Branch Name (Branch is off of {{.branchName}})",
		}, &i18n.Message{
			ID:    "CantDeleteCheckOutBranch",
			Other: "You cannot delete the checked out branch!",
		}, &i18n.Message{
			ID:    "DeleteBranch",
			Other: "Delete Branch",
		}, &i18n.Message{
			ID:    "DeleteBranchMessage",
			Other: "Are you sure you want to delete the branch {{.selectedBranchName}}?",
		}, &i18n.Message{
			ID:    "ForceDeleteBranchMessage",
			Other: "{{.selectedBranchName}} is not fully merged. Are you sure you want to delete it?",
		}, &i18n.Message{
			ID:    "rebaseBranch",
			Other: "rebase checked-out branch onto this branch",
		}, &i18n.Message{
			ID:    "CantRebaseOntoSelf",
			Other: "You cannot rebase a branch onto itself",
		}, &i18n.Message{
			ID:    "CantMergeBranchIntoItself",
			Other: "You cannot merge a branch into itself",
		}, &i18n.Message{
			ID:    "forceCheckout",
			Other: "force checkout",
		}, &i18n.Message{
			ID:    "merge",
			Other: "merge",
		}, &i18n.Message{
			ID:    "checkoutByName",
			Other: "checkout by name",
		}, &i18n.Message{
			ID:    "newBranch",
			Other: "new branch",
		}, &i18n.Message{
			ID:    "deleteBranch",
			Other: "delete branch",
		}, &i18n.Message{
			ID:    "forceDeleteBranch",
			Other: "delete branch (force)",
		}, &i18n.Message{
			ID:    "NoBranchesThisRepo",
			Other: "No branches for this repo",
		}, &i18n.Message{
			ID:    "NoTrackingThisBranch",
			Other: "There is no tracking for this branch",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "You cannot commit without a commit message",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: close, {{.keyBindConfirm}}: confirm",
		}, &i18n.Message{
			ID:    "close",
			Other: "close",
		}, &i18n.Message{
			ID:    "SureResetThisCommit",
			Other: "Are you sure you want to reset to this commit?",
		}, &i18n.Message{
			ID:    "ResetToCommit",
			Other: "Reset To Commit",
		}, &i18n.Message{
			ID:    "squashDown",
			Other: "squash down",
		}, &i18n.Message{
			ID:    "rename",
			Other: "rename",
		}, &i18n.Message{
			ID:    "resetToThisCommit",
			Other: "reset to this commit",
		}, &i18n.Message{
			ID:    "fixupCommit",
			Other: "fixup commit",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "No commits for this branch",
		}, &i18n.Message{
			ID:    "OnlySquashTopmostCommit",
			Other: "Can only squash topmost commit",
		}, &i18n.Message{
			ID:    "YouNoCommitsToSquash",
			Other: "You have no commits to squash with",
		}, &i18n.Message{
			ID:    "CantFixupWhileUnstagedChanges",
			Other: "Can't fixup while there are unstaged changes",
		}, &i18n.Message{
			ID:    "Fixup",
			Other: "Fixup",
		}, &i18n.Message{
			ID:    "SureFixupThisCommit",
			Other: "Are you sure you want to 'fixup' this commit? It will be merged into the commit below",
		}, &i18n.Message{
			ID:    "SureSquashThisCommit",
			Other: "Are you sure you want to squash this commit into the commit below?",
		}, &i18n.Message{
			ID:    "Squash",
			Other: "Squash",
		}, &i18n.Message{
			ID:    "pickCommit",
			Other: "pick commit (when mid-rebase)",
		}, &i18n.Message{
			ID:    "revertCommit",
			Other: "revert commit",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Can only reword topmost commit from within lazygit. Use shift+R instead",
		}, &i18n.Message{
			ID:    "renameCommit",
			Other: "reword commit",
		}, &i18n.Message{
			ID:    "deleteCommit",
			Other: "delete commit",
		}, &i18n.Message{
			ID:    "moveDownCommit",
			Other: "move commit down one",
		}, &i18n.Message{
			ID:    "moveUpCommit",
			Other: "move commit up one",
		}, &i18n.Message{
			ID:    "editCommit",
			Other: "edit commit",
		}, &i18n.Message{
			ID:    "amendToCommit",
			Other: "amend commit with staged changes",
		}, &i18n.Message{
			ID:    "renameCommitEditor",
			Other: "rename commit with editor",
		}, &i18n.Message{
			ID:    "PotentialErrInGetselectedCommit",
			Other: "potential error in getSelected Commit (mismatched ui and state)",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "No commits for this branch",
		}, &i18n.Message{
			ID:    "Error",
			Other: "Error",
		}, &i18n.Message{
			ID:    "resizingPopupPanel",
			Other: "resizing popup panel",
		}, &i18n.Message{
			ID:    "RunningSubprocess",
			Other: "running subprocess",
		}, &i18n.Message{
			ID:    "selectHunk",
			Other: "select hunk",
		}, &i18n.Message{
			ID:    "navigateConflicts",
			Other: "navigate conflicts",
		}, &i18n.Message{
			ID:    "pickHunk",
			Other: "pick hunk",
		}, &i18n.Message{
			ID:    "pickBothHunks",
			Other: "pick both hunks",
		}, &i18n.Message{
			ID:    "undo",
			Other: "undo",
		}, &i18n.Message{
			ID:    "undoReflog",
			Other: "undo (via reflog) (experimental)",
		}, &i18n.Message{
			ID:    "redoReflog",
			Other: "redo (via reflog) (experimental)",
		}, &i18n.Message{
			ID:    "pop",
			Other: "pop",
		}, &i18n.Message{
			ID:    "drop",
			Other: "drop",
		}, &i18n.Message{
			ID:    "apply",
			Other: "apply",
		}, &i18n.Message{
			ID:    "NoStashEntries",
			Other: "No stash entries",
		}, &i18n.Message{
			ID:    "StashDrop",
			Other: "Stash drop",
		}, &i18n.Message{
			ID:    "SureDropStashEntry",
			Other: "Are you sure you want to drop this stash entry?",
		}, &i18n.Message{
			ID:    "NoStashTo",
			Other: "No stash to {{.method}}",
		}, &i18n.Message{
			ID:    "NoTrackedStagedFilesStash",
			Other: "You have no tracked/staged files to stash",
		}, &i18n.Message{
			ID:    "StashChanges",
			Other: "Stash changes",
		}, &i18n.Message{
			ID:    "IssntListOfViews",
			Other: "{{.name}} is not in the list of views",
		}, &i18n.Message{
			ID:    "NoViewMachingNewLineFocusedSwitchStatement",
			Other: "No view matching newLineFocused switch statement",
		}, &i18n.Message{
			ID:    "newFocusedViewIs",
			Other: "new focused view is {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "No changed files",
		}, &i18n.Message{
			ID:    "MergeAborted",
			Other: "Merge aborted",
		}, &i18n.Message{
			ID:    "OpenConfig",
			Other: "open config file",
		}, &i18n.Message{
			ID:    "EditConfig",
			Other: "edit config file",
		}, &i18n.Message{
			ID:    "ForcePush",
			Other: "Force push",
		}, &i18n.Message{
			ID:    "ForcePushPrompt",
			Other: "Your branch has diverged from the remote branch. Press 'esc' to cancel, or 'enter' to force push.",
		}, &i18n.Message{
			ID:    "checkForUpdate",
			Other: "check for update",
		}, &i18n.Message{
			ID:    "CheckingForUpdates",
			Other: "Checking for updates...",
		}, &i18n.Message{
			ID:    "OnLatestVersionErr",
			Other: "You already have the latest version",
		}, &i18n.Message{
			ID:    "MajorVersionErr",
			Other: "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		}, &i18n.Message{
			ID:    "CouldNotFindBinaryErr",
			Other: "Could not find any binary at {{.url}}",
		}, &i18n.Message{
			ID:    "AnonymousReportingTitle",
			Other: "Help make lazygit better",
		}, &i18n.Message{
			ID:    "AnonymousReportingPrompt",
			Other: "Would you like to enable anonymous reporting data to help improve lazygit? (enter/esc)",
		}, &i18n.Message{
			ID:    "ShamelessSelfPromotionTitle",
			Other: "Shameless Self Promotion",
		}, &i18n.Message{
			ID: "ShamelessSelfPromotionMessage",
			Other: `Thanks for using lazygit! Three things to share with you:

1) lazygit now has basic mouse support!

2) If you want to learn about lazygit's features, watch this vid:
   https://youtu.be/CPLdltN7wgE

3) Github are now matching any donations dollar-for-dollar for the next 12 months, so if you've been tossing up over whether to click the donate link in the bottom right corner, now is the time!`,
		}, &i18n.Message{
			ID:    "GitconfigParseErr",
			Other: `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		}, &i18n.Message{
			ID:    "editFile",
			Other: `edit file`,
		}, &i18n.Message{
			ID:    "openFile",
			Other: `open file`,
		}, &i18n.Message{
			ID:    "ignoreFile",
			Other: `add to .gitignore`,
		}, &i18n.Message{
			ID:    "refreshFiles",
			Other: `refresh files`,
		}, &i18n.Message{
			ID:    "mergeIntoCurrentBranch",
			Other: `merge into currently checked out branch`,
		}, &i18n.Message{
			ID:    "ConfirmQuit",
			Other: `Are you sure you want to quit?`,
		}, &i18n.Message{
			ID:    "SwitchRepo",
			Other: `switch to a recent repo`,
		}, &i18n.Message{
			ID:    "UnsupportedGitService",
			Other: `Unsupported git service`,
		}, &i18n.Message{
			ID:    "createPullRequest",
			Other: `create pull request`,
		}, &i18n.Message{
			ID:    "NoBranchOnRemote",
			Other: `This branch doesn't exist on remote. You need to push it to remote first.`,
		}, &i18n.Message{
			ID:    "fetch",
			Other: `fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchTitle",
			Other: `No automatic git fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchBody",
			Other: `Lazygit can't use "git fetch" in a private repo; use 'f' in the files panel to run "git fetch" manually`,
		}, &i18n.Message{
			ID:    "StageLines",
			Other: `stage individual hunks/lines`,
		}, &i18n.Message{
			ID:    "FileStagingRequirements",
			Other: `Can only stage individual lines for tracked files`,
		}, &i18n.Message{
			ID:    "SelectHunk",
			Other: `select hunk`,
		}, &i18n.Message{
			ID:    "StageSelection",
			Other: `toggle line staged / unstaged`,
		}, &i18n.Message{
			ID:    "ResetSelection",
			Other: `delete change (git reset)`,
		}, &i18n.Message{
			ID:    "ToggleDragSelect",
			Other: `toggle drag select`,
		}, &i18n.Message{
			ID:    "ToggleSelectHunk",
			Other: `toggle select hunk`,
		}, &i18n.Message{
			ID:    "ToggleSelectionForPatch",
			Other: `add/remove line(s) to patch`,
		},
		&i18n.Message{
			ID:    "TogglePanel",
			Other: `switch to other panel`,
		},
		&i18n.Message{
			ID:    "CantStageStaged",
			Other: `You can't stage an already staged change!`,
		}, &i18n.Message{
			ID:    "ReturnToFilesPanel",
			Other: `return to files panel`,
		}, &i18n.Message{
			ID:    "CantFindHunks",
			Other: `Could not find any hunks in this patch`,
		}, &i18n.Message{
			ID:    "CantFindHunk",
			Other: `Could not find hunk`,
		}, &i18n.Message{
			ID:    "FastForward",
			Other: `fast-forward this branch from its upstream`,
		}, &i18n.Message{
			ID:    "Fetching",
			Other: "fetching and fast-forwarding {{.from}} -> {{.to}} ...",
		}, &i18n.Message{
			ID:    "FoundConflicts",
			Other: "Conflicts! To abort press 'esc', otherwise press 'enter'",
		}, &i18n.Message{
			ID:    "FoundConflictsTitle",
			Other: "Auto-merge failed",
		}, &i18n.Message{
			ID:    "Undo",
			Other: "undo",
		}, &i18n.Message{
			ID:    "PickHunk",
			Other: "pick hunk",
		}, &i18n.Message{
			ID:    "PickBothHunks",
			Other: "pick both hunks",
		}, &i18n.Message{
			ID:    "ViewMergeRebaseOptions",
			Other: "view merge/rebase options",
		}, &i18n.Message{
			ID:    "NotMergingOrRebasing",
			Other: "You are currently neither rebasing nor merging",
		}, &i18n.Message{
			ID:    "RecentRepos",
			Other: "recent repositories",
		}, &i18n.Message{
			ID:    "MergeOptionsTitle",
			Other: "Merge Options",
		}, &i18n.Message{
			ID:    "RebaseOptionsTitle",
			Other: "Rebase Options",
		}, &i18n.Message{
			ID:    "CommitMessageTitle",
			Other: "Commit Message",
		}, &i18n.Message{
			ID:    "Local-BranchesTitle",
			Other: "Branches Tab",
		}, &i18n.Message{
			ID:    "SearchTitle",
			Other: "Search",
		}, &i18n.Message{
			ID:    "TagsTitle",
			Other: "Tags Tab",
		}, &i18n.Message{
			ID:    "Branch-CommitsTitle",
			Other: "Commits Tab",
		}, &i18n.Message{
			ID:    "MenuTitle",
			Other: "Menu",
		}, &i18n.Message{
			ID:    "RemotesTitle",
			Other: "Remotes Tab",
		}, &i18n.Message{
			ID:    "CredentialsTitle",
			Other: "Credentials",
		}, &i18n.Message{
			ID:    "Remote-BranchesTitle",
			Other: "Remote Branches (in Remotes tab)",
		}, &i18n.Message{
			ID:    "Patch-BuildingTitle",
			Other: "Patch Building",
		}, &i18n.Message{
			ID:    "InformationTitle",
			Other: "Information",
		}, &i18n.Message{
			ID:    "SecondaryTitle",
			Other: "Secondary",
		}, &i18n.Message{
			ID:    "Reflog-CommitsTitle",
			Other: "Reflog Tab",
		}, &i18n.Message{
			ID:    "Title",
			Other: "Title",
		}, &i18n.Message{
			ID:    "GlobalTitle",
			Other: "Global Keybindings",
		}, &i18n.Message{
			ID:    "MerginTitle",
			Other: "Mergin",
		}, &i18n.Message{
			ID:    "ConflictsResolved",
			Other: "all merge conflicts resolved. Continue?",
		}, &i18n.Message{
			ID:    "RebasingTitle",
			Other: "Rebasing",
		}, &i18n.Message{
			ID:    "MergingTitle",
			Other: "Merging",
		}, &i18n.Message{
			ID:    "ConfirmRebase",
			Other: "Are you sure you want to rebase {{.checkedOutBranch}} onto {{.selectedBranch}}?",
		}, &i18n.Message{
			ID:    "ConfirmMerge",
			Other: "Are you sure you want to merge {{.selectedBranch}} into {{.checkedOutBranch}}?",
		}, &i18n.Message{}, &i18n.Message{
			ID:    "FwdNoUpstream",
			Other: "Cannot fast-forward a branch with no upstream",
		}, &i18n.Message{
			ID:    "FwdCommitsToPush",
			Other: "Cannot fast-forward a branch with commits to push",
		}, &i18n.Message{
			ID:    "ErrorOccurred",
			Other: "An error occurred! Please create an issue at https://github.com/jesseduffield/lazygit/issues",
		}, &i18n.Message{
			ID:    "NoRoom",
			Other: "Not enough room",
		}, &i18n.Message{
			ID:    "YouAreHere",
			Other: "YOU ARE HERE",
		}, &i18n.Message{
			ID:    "rewordNotSupported",
			Other: "rewording commits while interactively rebasing is not currently supported",
		}, &i18n.Message{
			ID:    "cherryPickCopy",
			Other: "copy commit (cherry-pick)",
		}, &i18n.Message{
			ID:    "cherryPickCopyRange",
			Other: "copy commit range (cherry-pick)",
		}, &i18n.Message{
			ID:    "pasteCommits",
			Other: "paste commits (cherry-pick)",
		}, &i18n.Message{
			ID:    "SureCherryPick",
			Other: "Are you sure you want to cherry-pick the copied commits onto this branch?",
		}, &i18n.Message{
			ID:    "CherryPick",
			Other: "Cherry-Pick",
		}, &i18n.Message{
			ID:    "CannotRebaseOntoFirstCommit",
			Other: "You cannot interactive rebase onto the first commit",
		}, &i18n.Message{
			ID:    "CannotSquashOntoSecondCommit",
			Other: "You cannot squash/fixup onto the second commit",
		}, &i18n.Message{
			ID:    "Donate",
			Other: "Donate",
		}, &i18n.Message{
			ID:    "PrevLine",
			Other: "select previous line",
		}, &i18n.Message{
			ID:    "NextLine",
			Other: "select next line",
		}, &i18n.Message{
			ID:    "PrevHunk",
			Other: "select previous hunk",
		}, &i18n.Message{
			ID:    "NextHunk",
			Other: "select next hunk",
		}, &i18n.Message{
			ID:    "PrevConflict",
			Other: "select previous conflict",
		}, &i18n.Message{
			ID:    "NextConflict",
			Other: "select next conflict",
		}, &i18n.Message{
			ID:    "SelectTop",
			Other: "select top hunk",
		}, &i18n.Message{
			ID:    "SelectBottom",
			Other: "select bottom hunk",
		}, &i18n.Message{
			ID:    "ScrollDown",
			Other: "scroll down",
		}, &i18n.Message{
			ID:    "ScrollUp",
			Other: "scroll up",
		}, &i18n.Message{
			ID:    "scrollUpMainPanel",
			Other: "scroll up main panel",
		}, &i18n.Message{
			ID:    "scrollDownMainPanel",
			Other: "scroll down main panel",
		}, &i18n.Message{
			ID:    "AmendCommitTitle",
			Other: "Amend Commit",
		}, &i18n.Message{
			ID:    "AmendCommitPrompt",
			Other: "Are you sure you want to amend this commit with your staged files?",
		}, &i18n.Message{
			ID:    "DeleteCommitTitle",
			Other: "Delete Commit",
		}, &i18n.Message{
			ID:    "DeleteCommitPrompt",
			Other: "Are you sure you want to delete this commit?",
		}, &i18n.Message{
			ID:    "SquashingStatus",
			Other: "squashing",
		}, &i18n.Message{
			ID:    "FixingStatus",
			Other: "fixing up",
		}, &i18n.Message{
			ID:    "DeletingStatus",
			Other: "deleting",
		}, &i18n.Message{
			ID:    "MovingStatus",
			Other: "moving",
		}, &i18n.Message{
			ID:    "RebasingStatus",
			Other: "rebasing",
		}, &i18n.Message{
			ID:    "AmendingStatus",
			Other: "amending",
		}, &i18n.Message{
			ID:    "CherryPickingStatus",
			Other: "cherry-picking",
		}, &i18n.Message{
			ID:    "UndoingStatus",
			Other: "undoing",
		}, &i18n.Message{
			ID:    "RedoingStatus",
			Other: "redoing",
		}, &i18n.Message{
			ID:    "CheckingOutStatus",
			Other: "checking out",
		}, &i18n.Message{
			ID:    "CommitFiles",
			Other: "Commit files",
		}, &i18n.Message{
			ID:    "viewCommitFiles",
			Other: "view commit's files",
		}, &i18n.Message{
			ID:    "CommitFilesTitle",
			Other: "Commit Files",
		}, &i18n.Message{
			ID:    "goBack",
			Other: "go back",
		}, &i18n.Message{
			ID:    "NoCommiteFiles",
			Other: "No files for this commit",
		}, &i18n.Message{
			ID:    "checkoutCommitFile",
			Other: "checkout file",
		}, &i18n.Message{
			ID:    "discardOldFileChange",
			Other: "discard this commit's changes to this file",
		}, &i18n.Message{
			ID:    "DiscardFileChangesTitle",
			Other: "Discard file changes",
		}, &i18n.Message{
			ID:    "DiscardFileChangesPrompt",
			Other: "Are you sure you want to discard this commit's changes to this file? If this file was created in this commit, it will be deleted",
		}, &i18n.Message{
			ID:    "DisabledForGPG",
			Other: "Feature not available for users using GPG",
		}, &i18n.Message{
			ID:    "CreateRepo",
			Other: "Not in a git repository. Create a new git repository? (y/n): ",
		}, &i18n.Message{
			ID:    "AutoStashTitle",
			Other: "Autostash?",
		}, &i18n.Message{
			ID:    "AutoStashPrompt",
			Other: "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		}, &i18n.Message{
			ID:    "StashPrefix",
			Other: "Auto-stashing changes for ",
		}, &i18n.Message{
			ID:    "viewDiscardOptions",
			Other: "view 'discard changes' options",
		}, &i18n.Message{
			ID:    "cancel",
			Other: "cancel",
		}, &i18n.Message{
			ID:    "discardAllChanges",
			Other: "discard all changes",
		}, &i18n.Message{
			ID:    "discardUnstagedChanges",
			Other: "discard unstaged changes",
		}, &i18n.Message{
			ID:    "discardAllChangesToAllFiles",
			Other: "nuke working tree",
		}, &i18n.Message{
			ID:    "discardAnyUnstagedChanges",
			Other: "discard unstaged changes",
		}, &i18n.Message{
			ID:    "discardUntrackedFiles",
			Other: "discard untracked files",
		}, &i18n.Message{
			ID:    "hardReset",
			Other: "hard reset",
		}, &i18n.Message{
			ID:    "hardResetUpstream",
			Other: "hard reset to upstream branch",
		}, &i18n.Message{
			ID:    "viewResetOptions",
			Other: `view reset options`,
		}, &i18n.Message{
			ID:    "createFixupCommit",
			Other: `create fixup commit for this commit`,
		}, &i18n.Message{
			ID:    "squashAboveCommits",
			Other: `squash above commits`,
		}, &i18n.Message{
			ID:    "SquashAboveCommits",
			Other: `Squash above commits`,
		}, &i18n.Message{
			ID:    "SureSquashAboveCommits",
			Other: `Are you sure you want to squash all fixup! commits above {{.commit}}?`,
		}, &i18n.Message{
			ID:    "CreateFixupCommit",
			Other: `Create fixup commit`,
		}, &i18n.Message{
			ID:    "SureCreateFixupCommit",
			Other: `Are you sure you want to create a fixup! commit for commit {{.commit}}?`,
		}, &i18n.Message{
			ID:    "executeCustomCommand",
			Other: "execute custom command",
		}, &i18n.Message{
			ID:    "CustomCommand",
			Other: "Custom Command:",
		}, &i18n.Message{
			ID:    "commitChangesWithoutHook",
			Other: "commit changes without pre-commit hook",
		}, &i18n.Message{
			ID:    "SkipHookPrefixNotConfigured",
			Other: "You have not configured a commit message prefix for skipping hooks. Set `git.skipHookPrefix = 'WIP'` in your config",
		}, &i18n.Message{
			ID:    "resetTo",
			Other: `reset to`,
		}, &i18n.Message{
			ID:    "pressEnterToReturn",
			Other: "Press enter to return to lazygit",
		}, &i18n.Message{
			ID:    "viewStashOptions",
			Other: "view stash options",
		}, &i18n.Message{
			ID:    "stashAllChanges",
			Other: "stash changes",
		}, &i18n.Message{
			ID:    "stashStagedChanges",
			Other: "stash staged changes",
		}, &i18n.Message{
			ID:    "stashOptions",
			Other: "Stash options",
		}, &i18n.Message{
			ID:    "notARepository",
			Other: "Error: must be run inside a git repository",
		}, &i18n.Message{
			ID:    "jump",
			Other: "jump to panel",
		}, &i18n.Message{
			ID:    "DiscardPatch",
			Other: "Discard Patch",
		}, &i18n.Message{
			ID:    "DiscardPatchConfirm",
			Other: "You can only build a patch from one commit at a time. Discard current patch?",
		}, &i18n.Message{
			ID:    "CantPatchWhileRebasingError",
			Other: "You cannot build a patch or run patch commands while in a merging or rebasing state",
		}, &i18n.Message{
			ID:    "toggleAddToPatch",
			Other: "toggle file included in patch",
		}, &i18n.Message{
			ID:    "ViewPatchOptions",
			Other: "view custom patch options",
		}, &i18n.Message{
			ID:    "PatchOptionsTitle",
			Other: "Patch Options",
		}, &i18n.Message{
			ID:    "NoPatchError",
			Other: "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		}, &i18n.Message{
			ID:    "enterFile",
			Other: "enter file to add selected lines to the patch",
		}, &i18n.Message{
			ID:    "ExitLineByLineMode",
			Other: `exit line-by-line mode`,
		}, &i18n.Message{
			ID:    "EnterUpstream",
			Other: `Enter upstream as '<remote> <branchname>'`,
		}, &i18n.Message{
			ID:    "EnterUpstreamWithSlash",
			Other: `Enter upstream as '<remote>/<branchname>'`,
		}, &i18n.Message{
			ID:    "notTrackingRemote",
			Other: "(not tracking any remote)",
		}, &i18n.Message{
			ID:    "ReturnToRemotesList",
			Other: `return to remotes list`,
		}, &i18n.Message{
			ID:    "addNewRemote",
			Other: `add new remote`,
		}, &i18n.Message{
			ID:    "newRemoteName",
			Other: `New remote name:`,
		}, &i18n.Message{
			ID:    "newRemoteUrl",
			Other: `New remote url:`,
		}, &i18n.Message{
			ID:    "editRemoteName",
			Other: `Enter updated remote name for {{ .remoteName }}:`,
		}, &i18n.Message{
			ID:    "editRemoteUrl",
			Other: `Enter updated remote url for {{ .remoteName }}:`,
		}, &i18n.Message{
			ID:    "removeRemote",
			Other: `remove remote`,
		}, &i18n.Message{
			ID:    "removeRemotePrompt",
			Other: "Are you sure you want to remove remote",
		}, &i18n.Message{
			ID:    "DeleteRemoteBranch",
			Other: "Delete Remote Branch",
		}, &i18n.Message{
			ID:    "DeleteRemoteBranchMessage",
			Other: "Are you sure you want to delete remote branch",
		}, &i18n.Message{
			ID:    "setUpstream",
			Other: "set as upstream of checked-out branch",
		}, &i18n.Message{
			ID:    "SetUpstreamTitle",
			Other: "Set upstream branch",
		}, &i18n.Message{
			ID:    "SetUpstreamMessage",
			Other: "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		}, &i18n.Message{
			ID:    "editRemote",
			Other: "edit remote",
		}, &i18n.Message{
			ID:    "tagCommit",
			Other: "tag commit",
		}, &i18n.Message{
			ID:    "TagNameTitle",
			Other: "Tag name:",
		}, &i18n.Message{
			ID:    "deleteTag",
			Other: "delete tag",
		}, &i18n.Message{
			ID:    "DeleteTagTitle",
			Other: "Delete tag",
		}, &i18n.Message{
			ID:    "DeleteTagPrompt",
			Other: "Are you sure you want to delete tag '{{.tagName}}'?",
		}, &i18n.Message{
			ID:    "PushTagTitle",
			Other: "remote to push tag '{{.tagName}}' to:",
		}, &i18n.Message{
			ID:    "pushTag",
			Other: "push tag",
		}, &i18n.Message{
			ID:    "createTag",
			Other: "create tag",
		}, &i18n.Message{
			ID:    "CreateTagTitle",
			Other: "Tag name:",
		}, &i18n.Message{
			ID:    "fetchRemote",
			Other: "fetch remote",
		}, &i18n.Message{
			ID:    "FetchingRemoteStatus",
			Other: "fetching remote",
		}, &i18n.Message{
			ID:    "checkoutCommit",
			Other: "checkout commit",
		}, &i18n.Message{
			ID:    "SureCheckoutThisCommit",
			Other: "Are you sure you want to checkout this commit?",
		}, &i18n.Message{
			ID:    "gitFlowOptions",
			Other: "show git-flow options",
		}, &i18n.Message{
			ID:    "NotAGitFlowBranch",
			Other: "This does not seem to be a git flow branch",
		}, &i18n.Message{
			ID:    "NewBranchNamePrompt",
			Other: "new {{.branchType}} name:",
		}, &i18n.Message{
			ID:    "IgnoreTracked",
			Other: "Ignore tracked file",
		}, &i18n.Message{
			ID:    "IgnoreTrackedPrompt",
			Other: "Are you sure you want to ignore a tracked file?",
		}, &i18n.Message{
			ID:    "viewResetToUpstreamOptions",
			Other: "view upstream reset options",
		}, &i18n.Message{
			ID:    "nextScreenMode",
			Other: "next screen mode (normal/half/fullscreen)",
		}, &i18n.Message{
			ID:    "prevScreenMode",
			Other: "prev screen mode",
		}, &i18n.Message{
			ID:    "startSearch",
			Other: "start search",
		}, &i18n.Message{
			ID:    "Panel",
			Other: "Panel",
		}, &i18n.Message{
			ID:    "Keybindings",
			Other: "Keybindings",
		}, &i18n.Message{
			ID:    "renameBranch",
			Other: "rename branch",
		}, &i18n.Message{
			ID:    "NewBranchNamePrompt",
			Other: "Enter new branch name for branch",
		}, &i18n.Message{
			ID:    "RenameBranchWarning",
			Other: "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		}, &i18n.Message{
			ID:    "openMenu",
			Other: "open menu",
		}, &i18n.Message{
			ID:    "closeMenu",
			Other: "close menu",
		}, &i18n.Message{
			ID:    "resetCherryPick",
			Other: "reset cherry-picked (copied) commits selection",
		}, &i18n.Message{
			ID:    "nextTab",
			Other: "next tab",
		}, &i18n.Message{
			ID:    "prevTab",
			Other: "previous tab",
		}, &i18n.Message{
			ID:    "cantUndoWhileRebasing",
			Other: "Can't undo while rebasing",
		}, &i18n.Message{
			ID:    "cantRedoWhileRebasing",
			Other: "Can't redo while rebasing",
		},
	)
}
