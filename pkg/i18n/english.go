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
			ID:    "FilesTitle",
			Other: "Files",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Branches",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Commits",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Stash",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Commit message",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "commit changes",
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
			ID:    "stashFiles",
			Other: "stash files",
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
			ID:    "addPatch",
			Other: "add patch",
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
			ID:    "CantIgnoreTrackFiles",
			Other: "Cannot ignore tracked files",
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
			ID:    "FileNoMergeCons",
			Other: "This file has no merge conflicts",
		}, &i18n.Message{
			ID:    "SureResetHardHead",
			Other: "Are you sure you want `reset --hard HEAD`? You may lose changes",
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
			Other: "Are you sure you want to force delete the branch {{.selectedBranchName}}?",
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
			Other: "Are you sure you want to fixup this commit? The commit beneath will be squashed up into this one",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Can only rename topmost commit",
		}, &i18n.Message{
			ID:    "renameCommit",
			Other: "rename commit",
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
			ID:    "settingPreviewsViewTo",
			Other: "setting previous view to: {{.oldViewName}}",
		}, &i18n.Message{
			ID:    "newFocusedViewIs",
			Other: "new focused view is {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "CantCloseConfirmationPrompt",
			Other: "Could not close confirmation prompt: {{.error}}",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "No changed files",
		}, &i18n.Message{
			ID:    "ClearFilePanel",
			Other: "Clear file panel",
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
			ID:    "GitconfigParseErr",
			Other: `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		}, &i18n.Message{
			ID:    "removeFile",
			Other: `delete if untracked / checkout if tracked`,
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
			ID:    "resetHard",
			Other: `reset hard`,
		}, &i18n.Message{
			ID:    "mergeIntoCurrentBranch",
			Other: `merge into currently checked out branch`,
		}, &i18n.Message{
			ID:    "ConfirmQuit",
			Other: `Are you sure you want to quit?`,
		}, &i18n.Message{
			ID:    "SwitchRepo",
			Other: `switch to a recent repo`,
		},
	)
}
