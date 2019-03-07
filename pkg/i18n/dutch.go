package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// addDutch will add all dutch translations
func addDutch(i18nObject *i18n.Bundle) error {

	// add the translations
	return i18nObject.AddMessages(language.Dutch,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Niet genoeg ruimte om de panelen te renderen",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Diff",
		}, &i18n.Message{
			ID:    "LogTitle",
			Other: "Log",
		}, &i18n.Message{
			ID:    "FilesTitle",
			Other: "Bestanden",
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
			ID:    "StagingMainTitle",
			Other: `Stage Lines/Hunks`,
		}, &i18n.Message{
			ID:    "MergingMainTitle",
			Other: "Resolve merge conflicts",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Commit bericht",
		}, &i18n.Message{
			ID:    "CredentialsUsername",
			Other: "Gebruikersnaam",
		}, &i18n.Message{
			ID:    "CredentialsPassword",
			Other: "Wachtwoord",
		}, &i18n.Message{
			ID:    "PassUnameWrong",
			Other: "Wachtwoord en/of gebruikersnaam verkeert",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "Commit veranderingen",
		}, &i18n.Message{
			ID:    "AmendLastCommit",
			Other: "wijzig laatste commit",
		}, &i18n.Message{
			ID:    "SureToAmend",
			Other: "Weet je zeker dat je de laatste commit wilt wijzigen? U kunt het commit-bericht wijzigen vanuit het commits-paneel.",
		}, &i18n.Message{
			ID:    "NoCommitToAmend",
			Other: "Er is geen commits om te wijzigen.",
		}, &i18n.Message{
			ID:    "CommitChangesWithEditor",
			Other: "commit veranderingen met de git editor",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "GlobalTitle",
			Other: "Global",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "navigeer",
		}, &i18n.Message{
			ID:    "menu",
			Other: "menu",
		}, &i18n.Message{
			ID:    "execute",
			Other: "uitvoeren",
		}, &i18n.Message{
			ID:    "stashFiles",
			Other: "stash-bestanden",
		}, &i18n.Message{
			ID:    "open",
			Other: "open",
		}, &i18n.Message{
			ID:    "ignore",
			Other: "negeren",
		}, &i18n.Message{
			ID:    "delete",
			Other: "verwijderen",
		}, &i18n.Message{
			ID:    "toggleStaged",
			Other: "toggle staged",
		}, &i18n.Message{
			ID:    "toggleStagedAll",
			Other: "toggle staged alle",
		}, &i18n.Message{
			ID:    "refresh",
			Other: "verversen",
		}, &i18n.Message{
			ID:    "push",
			Other: "push",
		}, &i18n.Message{
			ID:    "pull",
			Other: "pull",
		}, &i18n.Message{
			ID:    "addPatch",
			Other: "bewerkingen toevoegen",
		}, &i18n.Message{
			ID:    "edit",
			Other: "bewerken",
		}, &i18n.Message{
			ID:    "scroll",
			Other: "scroll",
		}, &i18n.Message{
			ID:    "abortMerge",
			Other: "samenvoegen afbreken",
		}, &i18n.Message{
			ID:    "resolveMergeConflicts",
			Other: "los merge conflicten op",
		}, &i18n.Message{
			ID:    "checkout",
			Other: "uitchecken",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Geen bestanden veranderd",
		}, &i18n.Message{
			ID:    "FileHasNoUnstagedChanges",
			Other: "Het bestand heeft geen unstaged veranderingen om toe te voegen",
		}, &i18n.Message{
			ID:    "CannotGitAdd",
			Other: "Kan commando niet uitvoeren git add --path untracked files",
		}, &i18n.Message{
			ID:    "CantIgnoreTrackFiles",
			Other: "Kan gevolgde bestanden niet negeren",
		}, &i18n.Message{
			ID:    "NoStagedFilesToCommit",
			Other: "Er zijn geen staged bestanden om te commiten",
		}, &i18n.Message{
			ID:    "NoFilesDisplay",
			Other: "Geen bestanden om te laten zien",
		}, &i18n.Message{
			ID:    "NotAFile",
			Other: "Dit is geen bestand",
		}, &i18n.Message{
			ID:    "PullWait",
			Other: "Pullen...",
		}, &i18n.Message{
			ID:    "PushWait",
			Other: "Pushen...",
		}, &i18n.Message{
			ID:    "FetchWait",
			Other: "Fetchen...",
		}, &i18n.Message{
			ID:    "FileNoMergeCons",
			Other: "Dit bestand heeft geen merge conflicten",
		}, &i18n.Message{
			ID:    "SureResetHardHead",
			Other: "Weet je het zeker dat je `reset --hard HEAD` en `clean -fd` wil uitvoeren? Het kan dat je hierdoor bestanden verliest",
		}, &i18n.Message{
			ID:    "SureTo",
			Other: "Weet je het zeker dat je {{.fileName}} wilt {{.deleteVerb}} (je veranderingen zullen worden verwijderd)",
		}, &i18n.Message{
			ID:    "AlreadyCheckedOutBranch",
			Other: "Je hebt deze branch al uitgecheckt",
		}, &i18n.Message{
			ID:    "SureForceCheckout",
			Other: "Weet je zeker dat je het uitchecken wil forceren? Al je lokale verandering zullen worden verwijdert",
		}, &i18n.Message{
			ID:    "ForceCheckoutBranch",
			Other: "Forceer uitchecken op deze branch",
		}, &i18n.Message{
			ID:    "BranchName",
			Other: "Branch naam",
		}, &i18n.Message{
			ID:    "NewBranchNameBranchOff",
			Other: "Nieuw branch naam (Branch is afgeleid van {{.branchName}})",
		}, &i18n.Message{
			ID:    "CantDeleteCheckOutBranch",
			Other: "Je kan een uitgecheckte branch niet verwijderen!",
		}, &i18n.Message{
			ID:    "DeleteBranch",
			Other: "Verwijder branch",
		}, &i18n.Message{
			ID:    "DeleteBranchMessage",
			Other: "Weet je zeker dat je branch {{.selectedBranchName}} wilt verwijderen?",
		}, &i18n.Message{
			ID:    "ForceDeleteBranchMessage",
			Other: "Weet je zeker dat je branch {{.selectedBranchName}} geforceerd wil verwijderen?",
		}, &i18n.Message{
			ID:    "rebaseBranch",
			Other: "rebase branch",
		}, &i18n.Message{
			ID:    "CantMergeBranchIntoItself",
			Other: "Je kan niet een branch in zichzelf mergen",
		}, &i18n.Message{
			ID:    "forceCheckout",
			Other: "forceer checkout",
		}, &i18n.Message{
			ID:    "merge",
			Other: "samenvoegen",
		}, &i18n.Message{
			ID:    "checkoutByName",
			Other: "uitchecken bij naam",
		}, &i18n.Message{
			ID:    "newBranch",
			Other: "nieuwe branch",
		}, &i18n.Message{
			ID:    "deleteBranch",
			Other: "verwijder branch",
		}, &i18n.Message{
			ID:    "forceDeleteBranch",
			Other: "verwijder branch (forceer)",
		}, &i18n.Message{
			ID:    "NoBranchesThisRepo",
			Other: "Geen branches voor deze repo",
		}, &i18n.Message{
			ID:    "NoTrackingThisBranch",
			Other: "deze branch wordt niet gevolgd",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "Je kan geen commit maken zonder commit bericht",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: Sluiten, {{.keyBindConfirm}}: Bevestigen",
		}, &i18n.Message{
			ID:    "close",
			Other: "sluiten",
		}, &i18n.Message{
			ID:    "SureResetThisCommit",
			Other: "Weet je het zeker dat je wil resetten naar deze commit?",
		}, &i18n.Message{
			ID:    "ResetToCommit",
			Other: "Reset Naar Commit",
		}, &i18n.Message{
			ID:    "squashDown",
			Other: "squash beneden",
		}, &i18n.Message{
			ID:    "rename",
			Other: "hernoemen",
		}, &i18n.Message{
			ID:    "resetToThisCommit",
			Other: "reset naar deze commit",
		}, &i18n.Message{
			ID:    "fixupCommit",
			Other: "Fixup commit",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "Er zijn geen commits voor deze branch",
		}, &i18n.Message{
			ID:    "OnlySquashTopmostCommit",
			Other: "Kan alleen bovenste commit squashen",
		}, &i18n.Message{
			ID:    "YouNoCommitsToSquash",
			Other: "Je hebt geen commits om mee te squashen",
		}, &i18n.Message{
			ID:    "CantFixupWhileUnstagedChanges",
			Other: "Kan geen Fixup uitvoeren op unstaged veranderingen",
		}, &i18n.Message{
			ID:    "Fixup",
			Other: "Fixup",
		}, &i18n.Message{
			ID:    "SureFixupThisCommit",
			Other: "Weet je zeker dat je fixup wil uitvoeren op deze commit? De commit hieronder zol worden squashed in deze",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Je kan alleen de bovenste commit hernoemen",
		}, &i18n.Message{
			ID:    "renameCommit",
			Other: "hernoem commit",
		}, &i18n.Message{
			ID:    "renameCommitEditor",
			Other: "rename commit with editor",
		}, &i18n.Message{
			ID:    "PotentialErrInGetselectedCommit",
			Other: "Er is mogelijk een error in getSelected Commit (geen match tussen ui en state)",
		}, &i18n.Message{
			ID:    "Error",
			Other: "Foutmelding",
		}, &i18n.Message{
			ID:    "resizingPopupPanel",
			Other: "resizen popup paneel",
		}, &i18n.Message{
			ID:    "RunningSubprocess",
			Other: "subprocess lopend",
		}, &i18n.Message{
			ID:    "selectHunk",
			Other: "selecteer stuk",
		}, &i18n.Message{
			ID:    "navigateConflicts",
			Other: "navigeer conflicts",
		}, &i18n.Message{
			ID:    "pickHunk",
			Other: "kies stuk",
		}, &i18n.Message{
			ID:    "pickBothHunks",
			Other: "kies beide stukken",
		}, &i18n.Message{
			ID:    "undo",
			Other: "ongedaan maken",
		}, &i18n.Message{
			ID:    "pop",
			Other: "pop",
		}, &i18n.Message{
			ID:    "drop",
			Other: "drop",
		}, &i18n.Message{
			ID:    "apply",
			Other: "toepassen",
		}, &i18n.Message{
			ID:    "NoStashEntries",
			Other: "Geen stash items",
		}, &i18n.Message{
			ID:    "StashDrop",
			Other: "Stash drop",
		}, &i18n.Message{
			ID:    "SureDropStashEntry",
			Other: "Weet je het zeker dat je deze stash entry wil laten vallen?",
		}, &i18n.Message{
			ID:    "NoStashTo",
			Other: "Geen stash voor {{.method}}",
		}, &i18n.Message{
			ID:    "NoTrackedStagedFilesStash",
			Other: "Je hebt geen tracked/staged bestanden om te laten stashen",
		}, &i18n.Message{
			ID:    "StashChanges",
			Other: "Stash veranderingen",
		}, &i18n.Message{
			ID:    "IssntListOfViews",
			Other: "{{.name}} is niet in de lijst van weergaves",
		}, &i18n.Message{
			ID:    "NoViewMachingNewLineFocusedSwitchStatement",
			Other: "Er machen geen weergave met de newLineFocused switch declaratie",
		}, &i18n.Message{
			ID:    "newFocusedViewIs",
			Other: "nieuw gefocussed weergave is {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "CantCloseConfirmationPrompt",
			Other: "Kon de bevestiging prompt niet sluiten: {{.error}}",
		}, &i18n.Message{
			ID:    "ClearFilePanel",
			Other: "maak bestandsvenster leeg",
		}, &i18n.Message{
			ID:    "MergeAborted",
			Other: "Merge afgebroken",
		}, &i18n.Message{
			ID:    "OpenConfig",
			Other: "open config file",
		}, &i18n.Message{
			ID:    "EditConfig",
			Other: "verander config file",
		}, &i18n.Message{
			ID:    "ForcePush",
			Other: "Forceer push",
		}, &i18n.Message{
			ID:    "ForcePushPrompt",
			Other: "Jouw branch is afgeweken van de remote branch. Druk 'esc' om te annuleren, of 'enter' om geforceert te pushen.",
		}, &i18n.Message{
			ID:    "checkForUpdate",
			Other: "check voor updates",
		}, &i18n.Message{
			ID:    "CheckingForUpdates",
			Other: "zoeken naar updates...",
		}, &i18n.Message{
			ID:    "OnLatestVersionErr",
			Other: "Je hebt al de laatste versie",
		}, &i18n.Message{
			ID:    "MajorVersionErr",
			Other: "Nieuwe versie ({{.newVersion}}) is niet backwards compatibele vergeleken met de huidige versie ({{.currentVersion}})",
		}, &i18n.Message{
			ID:    "CouldNotFindBinaryErr",
			Other: "Kon geen binary vinden op {{.url}}",
		}, &i18n.Message{
			ID:    "AnonymousReportingTitle",
			Other: "Help lazygit te verbeteren",
		}, &i18n.Message{
			ID:    "AnonymousReportingPrompt",
			Other: "Zou je anonieme data rapportage willen aanzetten om lazygit beter te kunnen maken? (enter/esc)",
		}, &i18n.Message{
			ID:    "GitconfigParseErr",
			Other: `Gogit kon je gitconfig bestand niet goed parsen door de aanwezigheid van losstaande '\' tekens. Het weghalen van deze tekens zou het probleem moeten oplossen. `,
		}, &i18n.Message{
			ID:    "removeFile",
			Other: `Verwijder als untracked / uitchecken wordt gevolgd (ga weg)`,
		}, &i18n.Message{
			ID:    "editFile",
			Other: `verander bestand`,
		}, &i18n.Message{
			ID:    "openFile",
			Other: `open bestand`,
		}, &i18n.Message{
			ID:    "ignoreFile",
			Other: `voeg toe aan .gitignore`,
		}, &i18n.Message{
			ID:    "refreshFiles",
			Other: `refresh bestanden`,
		}, &i18n.Message{
			ID:    "resetHard",
			Other: `harde reset and verwijderen ongevolgde bestanden`,
		}, &i18n.Message{
			ID:    "mergeIntoCurrentBranch",
			Other: `merge in met huidige checked out branch`,
		}, &i18n.Message{
			ID:    "ConfirmQuit",
			Other: `Weet je zeker dat je dit programma wil sluiten?`,
		}, &i18n.Message{
			ID:    "SwitchRepo",
			Other: "wissel naar een recente repo",
		}, &i18n.Message{
			ID:    "UnsupportedGitService",
			Other: `Niet-ondersteunde git-service`,
		}, &i18n.Message{
			ID:    "createPullRequest",
			Other: `maak een pull-aanvraag`,
		}, &i18n.Message{
			ID:    "NoBranchOnRemote",
			Other: `Deze branch bestaat niet op de remote. U moet het eerst naar de remote pushen.`,
		}, &i18n.Message{
			ID:    "fetch",
			Other: `fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchTitle",
			Other: `Geen automatiese git fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchBody",
			Other: `Lazygit kan niet "git fetch" uitvoeren in een privé repository, gebruik f in het branches paneel om "git fetch" manueel uit te voeren`,
		}, &i18n.Message{
			ID:    "StageLines",
			Other: `stage individuele hunks/lijnen`,
		}, &i18n.Message{
			ID:    "FileStagingRequirements",
			Other: `Kan alleen individuele lijnen stagen van getrackte bestanden met onstaged veranderingen`,
		}, &i18n.Message{
			ID:    "StagingTitle",
			Other: `Stage Lines/Hunks`,
		}, &i18n.Message{
			ID:    "StageHunk",
			Other: `stage hunk`,
		}, &i18n.Message{
			ID:    "StageLine",
			Other: `stage lijn`,
		}, &i18n.Message{
			ID:    "EscapeStaging",
			Other: `ga terug naar het bestanden paneel`,
		}, &i18n.Message{
			ID:    "CantFindHunks",
			Other: `Kan geen hunks vinden in deze patch`,
		}, &i18n.Message{
			ID:    "CantFindHunk",
			Other: `Kan geen hunk vinden`,
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
		}, &i18n.Message{
			ID:    "FwdNoUpstream",
			Other: "Cannot fast-forward a branch with no upstream",
		}, &i18n.Message{
			ID:    "ErrorOccurred",
			Other: "An error occurred! Please create an issue at https://github.com/jesseduffield/lazygit/issues",
		}, &i18n.Message{
			ID:    "FwdCommitsToPush",
			Other: "Cannot fast-forward a branch with commits to push",
		}, &i18n.Message{
			ID:    "MainTitle",
			Other: "Main",
		}, &i18n.Message{
			ID:    "NormalTitle",
			Other: "Normal",
		}, &i18n.Message{
			ID:    "softReset",
			Other: "soft reset to last commit",
		}, &i18n.Message{
			ID:    "SoftReset",
			Other: "Soft reset",
		}, &i18n.Message{
			ID:    "ConfirmSoftReset",
			Other: "Are you sure you want to `reset --soft HEAD^`? The changes in your topmost commit will be placed in your working tree",
		}, &i18n.Message{
			ID:    "CantRebaseOntoSelf",
			Other: "You cannot rebase a branch onto itself",
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
			ID:    "FoundConflicts",
			Other: "Damn, conflicts! To abort press 'esc', otherwise press 'enter'",
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
			ID:    "ConflictsResolved",
			Other: "all merge conflicts resolved. Continue?",
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
		},
	)
}
