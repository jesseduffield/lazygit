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
			ID:    "UnstagedChanges",
			Other: `Unstaged wijzigingen`,
		}, &i18n.Message{
			ID:    "StagedChanges",
			Other: `Staged Wijzigingen`,
		}, &i18n.Message{
			ID:    "PatchBuildingMainTitle",
			Other: `Voeg lijnen/hunks toe aan Patch`,
		}, &i18n.Message{
			ID:    "MergingMainTitle",
			Other: "Los merge conflicten op",
		}, &i18n.Message{
			ID:    "MainTitle",
			Other: "Hooft",
		}, &i18n.Message{
			ID:    "StagingTitle",
			Other: "Staging",
		}, &i18n.Message{
			ID:    "MergingTitle",
			Other: "Merging",
		}, &i18n.Message{
			ID:    "NormalTitle",
			Other: "Normaal",
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
			Other: "Globaal",
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
			ID:    "MergeConflictsTitle",
			Other: "Merge Conflicten",
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
			ID:    "softReset",
			Other: "zacht reset",
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
			ID:    "CantRebaseOntoSelf",
			Other: "Je kan niet een branch rebasen op zichzelf",
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
			ID:    "CommitMessageConfirm",
			Other: "{{.keyBindClose}}: Sluiten, {{.keyBindNewLine}}: Nieuwe lijn, {{.keyBindConfirm}}: Bevestig",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "Je kan geen commit maken zonder commit bericht",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: Sluiten, {{.keyBindConfirm}}: Bevestig",
		}, &i18n.Message{
			ID:    "close",
			Other: "sluiten",
		}, &i18n.Message{
			ID:    "quit",
			Other: "quit",
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
			ID:    "SureSquashThisCommit",
			Other: "Weet je zeker dat je deze commit wil samenvoegen met de commit hieronder?",
		}, &i18n.Message{
			ID:    "Squash",
			Other: "Squash",
		}, &i18n.Message{
			ID:    "pickCommit",
			Other: "kies commit (wanneer midden in rebase)",
		}, &i18n.Message{
			ID:    "revertCommit",
			Other: "commit ongedaan maken",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Je kan alleen de bovenste commit hernoemen",
		}, &i18n.Message{
			ID:    "renameCommit",
			Other: "hernoem commit",
		}, &i18n.Message{
			ID:    "deleteCommit",
			Other: "verwijder commit",
		}, &i18n.Message{
			ID:    "moveDownCommit",
			Other: "verplaats commit 1 naar beneden",
		}, &i18n.Message{
			ID:    "moveUpCommit",
			Other: "verplaats commit 1 naar boven",
		}, &i18n.Message{
			ID:    "editCommit",
			Other: "wijzig commit",
		}, &i18n.Message{
			ID:    "amendToCommit",
			Other: "wijzig commit met staged veranderingen",
		}, &i18n.Message{
			ID:    "renameCommitEditor",
			Other: "hernoem commit met editor",
		}, &i18n.Message{
			ID:    "PotentialErrInGetselectedCommit",
			Other: "Er is mogelijk een error in getSelected Commit (geen match tussen ui en state)",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "Geen commits in deze branch",
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
			ID:    "undoReflog",
			Other: "ongedaan maken (via reflog) (experimenteel)",
		}, &i18n.Message{
			ID:    "redoReflog",
			Other: "redo (via reflog) (experimenteel)",
		}, &i18n.Message{
			ID:    "pop",
			Other: "pop",
		}, &i18n.Message{
			ID:    "drop",
			Other: "laten vallen",
		}, &i18n.Message{
			ID:    "apply",
			Other: "toepassen",
		}, &i18n.Message{
			ID:    "NoStashEntries",
			Other: "Geen stash items",
		}, &i18n.Message{
			ID:    "StashDrop",
			Other: "Stash laten vallen",
		}, &i18n.Message{
			ID:    "SureDropStashEntry",
			Other: "Weet je het zeker dat je deze stash entry wil laten vallen?",
		}, &i18n.Message{
			ID:    "StashPop",
			Other: "Stash pop",
		}, &i18n.Message{
			ID:    "SurePopStashEntry",
			Other: "Weet je zeker dat je deze stash entry wil poppen?",
		}, &i18n.Message{
			ID:    "StashApply",
			Other: "Stash toepassen",
		}, &i18n.Message{
			ID:    "SureApplyStashEntry",
			Other: "Weet je zeker dat je deze stash entry wil toepassen?",
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
			ID:    "newFocusedViewIs",
			Other: "nieuw gefocussed weergave is {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Geen veranderde bestanden",
		}, &i18n.Message{
			ID:    "MergeAborted",
			Other: "Merge afgebroken",
		}, &i18n.Message{
			ID:    "OpenConfig",
			Other: "open config bestand",
		}, &i18n.Message{
			ID:    "EditConfig",
			Other: "verander config bestand",
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
			ID: "IntroPopupMessage",
			Other: `Bedankt voor het gebruik maken van lazygit! 2 dingen die je moet weten:

1) Als je meer van lazygit zijn features wilt leren bekijk dan deze video:
   https://youtu.be/CPLdltN7wgE

3) Als je git gebruikt, ben je een programmeur! Met jouw hulp kunnen we lazygit verbeteren, dus overweeg om een ​​donateur te worden en mee te doen aan het plezier op
   https://github.com/jesseduffield/lazygit`,
		}, &i18n.Message{
			ID:    "GitconfigParseErr",
			Other: `Gogit kon je gitconfig bestand niet goed parsen door de aanwezigheid van losstaande '\' tekens. Het weghalen van deze tekens zou het probleem moeten oplossen. `,
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
			Other: `Geen automatische git fetch`,
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
			ID:    "SelectHunk",
			Other: `selecteer hunk`,
		}, &i18n.Message{
			ID:    "StageSelection",
			Other: `toggle lijnen staged / unstaged`,
		}, &i18n.Message{
			ID:    "ResetSelection",
			Other: `verwijdert change (git reset)`,
		}, &i18n.Message{
			ID:    "ToggleDragSelect",
			Other: `toggle drag  selecteer`,
		}, &i18n.Message{
			ID:    "ToggleSelectHunk",
			Other: `toggle  selecteer hunk`,
		}, &i18n.Message{
			ID:    "ToggleSelectionForPatch",
			Other: `voeg toe/verwijder lijn(en) in patch`,
		}, &i18n.Message{
			ID:    "TogglePanel",
			Other: `ga naar een ander paneel`,
		}, &i18n.Message{
			ID:    "CantStageStaged",
			Other: `Je kan niet al gestaged verandering stagen!`,
		}, &i18n.Message{
			ID:    "ReturnToFilesPanel",
			Other: `ga terug naar het bestanden paneel`,
		}, &i18n.Message{
			ID:    "CantFindHunks",
			Other: `Kan geen hunks vinden in deze patch`,
		}, &i18n.Message{
			ID:    "CantFindHunk",
			Other: `Kan geen hunk vinden`,
		}, &i18n.Message{
			ID:    "FastForward",
			Other: `fast-forward deze branch vanaf zijn upstream`,
		}, &i18n.Message{
			ID:    "Fetching",
			Other: "fetching en fast-forwarding {{.from}} -> {{.to}} ...",
		}, &i18n.Message{
			ID:    "FoundConflicts",
			Other: "Conflicten!, Om af te breken druk 'esc', anders druk op 'enter'",
		}, &i18n.Message{
			ID:    "FoundConflictsTitle",
			Other: "Auto-merge mislukt",
		}, &i18n.Message{
			ID:    "Undo",
			Other: "ongedaan maken",
		}, &i18n.Message{
			ID:    "PickHunk",
			Other: "kies hunk",
		}, &i18n.Message{
			ID:    "PickBothHunks",
			Other: "kies bijde hunks",
		}, &i18n.Message{
			ID:    "ViewMergeRebaseOptions",
			Other: "bekijk merge/rebase opties",
		}, &i18n.Message{
			ID:    "NotMergingOrRebasing",
			Other: "Je bent momenteel niet aan het rebasen of mergen",
		}, &i18n.Message{
			ID:    "RecentRepos",
			Other: "recente repositories",
		}, &i18n.Message{
			ID:    "MergeOptionsTitle",
			Other: "Merge Opties",
		}, &i18n.Message{
			ID:    "RebaseOptionsTitle",
			Other: "Rebase Opties",
		}, &i18n.Message{
			ID:    "CommitMessageTitle",
			Other: "Commit Bericht",
		}, &i18n.Message{
			ID:    "Local-BranchesTitle",
			Other: "Branches Tab",
		}, &i18n.Message{
			ID:    "SearchTitle",
			Other: "Zoek",
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
			Other: "Patch Bouwen",
		}, &i18n.Message{
			ID:    "InformationTitle",
			Other: "Informatie",
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
			Other: "Globaale Sneltoetsen",
		}, &i18n.Message{
			ID:    "MerginTitle",
			Other: "Mergen",
		}, &i18n.Message{
			ID:    "ConflictsResolved",
			Other: "alle merge conflicten zijn opgelost. Wilt je verder gaan?",
		}, &i18n.Message{
			ID:    "RebasingTitle",
			Other: "Rebasen",
		}, &i18n.Message{
			ID:    "MergingTitle",
			Other: "Merggen",
		}, &i18n.Message{
			ID:    "ConfirmRebase",
			Other: "Weet je zeker dat je {{.checkedOutBranch}} op {{.selectedBranch}} wil rebasen?",
		}, &i18n.Message{
			ID:    "ConfirmMerge",
			Other: "Weet je zeker dat je {{.selectedBranch}} in {{.checkedOutBranch}} wil mergen?",
		}, &i18n.Message{
			ID:    "FwdNoUpstream",
			Other: "Kan niet de branch vooruitspoelen zonder upstream",
		}, &i18n.Message{
			ID:    "FwdCommitsToPush",
			Other: "Je kan niet vooruitspoelen als de branch geen nieuwe commits heeft",
		}, &i18n.Message{
			ID:    "ErrorOccurred",
			Other: "Er is iets fout gegaan! Zou je hier een issue aan willen maken: https://github.com/jesseduffield/lazygit/issues",
		}, &i18n.Message{
			ID:    "NoRoom",
			Other: "Niet genoeg ruimte",
		}, &i18n.Message{
			ID:    "YouAreHere",
			Other: "JE BENT HIER",
		}, &i18n.Message{
			ID:    "rewordNotSupported",
			Other: "herformatteren van commits in interactief rebasen is nog niet ondersteund",
		}, &i18n.Message{
			ID:    "cherryPickCopy",
			Other: "kopiëer commit (cherry-pick)",
		}, &i18n.Message{
			ID:    "cherryPickCopyRange",
			Other: "kopiëer commit reeks (cherry-pick)",
		}, &i18n.Message{
			ID:    "pasteCommits",
			Other: "plak commits (cherry-pick)",
		}, &i18n.Message{
			ID:    "SureCherryPick",
			Other: "Weet je zeker dat je de gekopieerde commits naar deze branch wil cherry-picken?",
		}, &i18n.Message{
			ID:    "CherryPick",
			Other: "Cherry-Pick",
		}, &i18n.Message{
			ID:    "CannotRebaseOntoFirstCommit",
			Other: "Je kan niet interactief rebasen naar de eerste commit",
		}, &i18n.Message{
			ID:    "CannotSquashOntoSecondCommit",
			Other: "Je kan niet een squash/fixup doen naar de 2de commit",
		}, &i18n.Message{
			ID:    "Donate",
			Other: "Doneer",
		}, &i18n.Message{
			ID:    "PrevLine",
			Other: "selecteer de vorige lijn",
		}, &i18n.Message{
			ID:    "NextLine",
			Other: "selecteer de volgende lijn",
		}, &i18n.Message{
			ID:    "PrevHunk",
			Other: "selecteer de vorige hunk",
		}, &i18n.Message{
			ID:    "NextHunk",
			Other: "selecteer de volgende hunk",
		}, &i18n.Message{
			ID:    "PrevConflict",
			Other: "selecteer voorgaand conflict",
		}, &i18n.Message{
			ID:    "NextConflict",
			Other: "selecteer volgende conflict",
		}, &i18n.Message{
			ID:    "SelectTop",
			Other: "selecteer bovenste hunk",
		}, &i18n.Message{
			ID:    "SelectBottom",
			Other: "selecteer onderste hunk",
		}, &i18n.Message{
			ID:    "ScrollDown",
			Other: "scroll omlaag",
		}, &i18n.Message{
			ID:    "ScrollUp",
			Other: "scroll omhoog",
		}, &i18n.Message{
			ID:    "scrollUpMainPanel",
			Other: "scroll naar beneden vanaf hooft paneel",
		}, &i18n.Message{
			ID:    "scrollDownMainPanel",
			Other: "scroll naar beneden vabaf hooft paneel",
		}, &i18n.Message{
			ID:    "AmendCommitTitle",
			Other: "Commit wijzigen",
		}, &i18n.Message{
			ID:    "AmendCommitPrompt",
			Other: "Weet je zeker dat je deze commit wil wijzigen met de vorige staged bestanden?",
		}, &i18n.Message{
			ID:    "DeleteCommitTitle",
			Other: "Verwijder Commit",
		}, &i18n.Message{
			ID:    "DeleteCommitPrompt",
			Other: "Weet je zeker dat je deze commit wil verwijderen?",
		}, &i18n.Message{
			ID:    "SquashingStatus",
			Other: "squashen",
		}, &i18n.Message{
			ID:    "FixingStatus",
			Other: "fixing up",
		}, &i18n.Message{
			ID:    "DeletingStatus",
			Other: "verwijderen",
		}, &i18n.Message{
			ID:    "MovingStatus",
			Other: "verplaatsen",
		}, &i18n.Message{
			ID:    "RebasingStatus",
			Other: "rebasen",
		}, &i18n.Message{
			ID:    "AmendingStatus",
			Other: "wijzigen",
		}, &i18n.Message{
			ID:    "CherryPickingStatus",
			Other: "cherry-picken",
		}, &i18n.Message{
			ID:    "UndoingStatus",
			Other: "ongedaan maken",
		}, &i18n.Message{
			ID:    "RedoingStatus",
			Other: "redoing",
		}, &i18n.Message{
			ID:    "CheckingOutStatus",
			Other: "uitchecken",
		}, &i18n.Message{
			ID:    "CommitFiles",
			Other: "Commit bestanden",
		}, &i18n.Message{
			ID:    "viewCommitFiles",
			Other: "bekijk gecommite bestanden",
		}, &i18n.Message{
			ID:    "CommitFilesTitle",
			Other: "Commit bestanden",
		}, &i18n.Message{
			ID:    "goBack",
			Other: "ga terug",
		}, &i18n.Message{
			ID:    "NoCommiteFiles",
			Other: "Geen bestanden voor deze commit",
		}, &i18n.Message{
			ID:    "checkoutCommitFile",
			Other: "bestand uitchecken",
		}, &i18n.Message{
			ID:    "discardOldFileChange",
			Other: "uitsluit deze commit zijn veranderingen aan dit bestand",
		}, &i18n.Message{
			ID:    "DiscardFileChangesTitle",
			Other: "uitsluit bestand zijn veranderingen",
		}, &i18n.Message{
			ID:    "DiscardFileChangesPrompt",
			Other: "Weet je zeker dat je de wijzigingen van deze commit in dit bestand wilt weggooien? Als dit bestand is gecreëerd in deze commit dan zal dit bestand worden verwijdert",
		}, &i18n.Message{
			ID:    "DisabledForGPG",
			Other: "Onderdelen niet beschikbaar voor gebruikers die GPG gebruiken",
		}, &i18n.Message{
			ID:    "CreateRepo",
			Other: "Niet in een git repository. Creëer een nieuwe git repository? (y/n): ",
		}, &i18n.Message{
			ID:    "AutoStashTitle",
			Other: "Autostash?",
		}, &i18n.Message{
			ID:    "AutoStashPrompt",
			Other: "Je moet je veranderingen stashen en poppen om ze over te bregen. Dit automatisch doen? (enter/esc)",
		}, &i18n.Message{
			ID:    "StashPrefix",
			Other: "Auto-stashing veranderingen voor ",
		}, &i18n.Message{
			ID:    "viewDiscardOptions",
			Other: "bekijk 'veranderingen ongedaan maken' opties",
		}, &i18n.Message{
			ID:    "cancel",
			Other: "anuleren",
		}, &i18n.Message{
			ID:    "discardAllChanges",
			Other: "negeer alle wijzigingen",
		}, &i18n.Message{
			ID:    "discardUnstagedChanges",
			Other: "negeer unstaged wijzigingen",
		}, &i18n.Message{
			ID:    "discardAllChangesToAllFiles",
			Other: "verwijder werkende tree",
		}, &i18n.Message{
			ID:    "discardAnyUnstagedChanges",
			Other: "discard unstaged wijzigingen",
		}, &i18n.Message{
			ID:    "discardUntrackedFiles",
			Other: "negeer niet-gevonden bestanden",
		}, &i18n.Message{
			ID:    "viewResetOptions",
			Other: `bekijk reset opties`,
		}, &i18n.Message{
			ID:    "hardReset",
			Other: "harde reset",
		}, &i18n.Message{
			ID:    "hardResetUpstream",
			Other: "harde naar upstream branch",
		}, &i18n.Message{
			ID:    "viewResetOptions",
			Other: `bekijk reset opties`,
		}, &i18n.Message{
			ID:    "createFixupCommit",
			Other: `creëer fixup commit voor deze commit`,
		}, &i18n.Message{
			ID:    "squashAboveCommits",
			Other: `squash bovenstaande commits`,
		}, &i18n.Message{
			ID:    "SquashAboveCommits",
			Other: `Squash bovenstaande commits`,
		}, &i18n.Message{
			ID:    "SureSquashAboveCommits",
			Other: `Weet je zeker dat je alles wil squash/fixup! voor de bovenstaand commits {{.commit}}?`,
		}, &i18n.Message{
			ID:    "CreateFixupCommit",
			Other: `Creëer fixup commit`,
		}, &i18n.Message{
			ID:    "SureCreateFixupCommit",
			Other: `Weet je zeker dat je een fixup wil maken! commit voor commit {{.commit}}?`,
		}, &i18n.Message{
			ID:    "executeCustomCommand",
			Other: "voor aangepast commando uit",
		}, &i18n.Message{
			ID:    "CustomCommand",
			Other: "Aangepast commando:",
		}, &i18n.Message{
			ID:    "commitChangesWithoutHook",
			Other: "commit veranderingen zonder pre-commit hook",
		}, &i18n.Message{
			ID:    "SkipHookPrefixNotConfigured",
			Other: "Je hebt nog niet een commit bericht voorvoegsel ingesteld voor het overslaan van hooks. Set `git.skipHookPrefix = 'WIP'` in je config",
		}, &i18n.Message{
			ID:    "resetTo",
			Other: `reset naar`,
		}, &i18n.Message{
			ID:    "pressEnterToReturn",
			Other: "Press om terug te gaan naar lazygit",
		}, &i18n.Message{
			ID:    "viewStashOptions",
			Other: "bekijk stash opties",
		}, &i18n.Message{
			ID:    "stashAllChanges",
			Other: "stash-bestanden",
		}, &i18n.Message{
			ID:    "stashStagedChanges",
			Other: "stash staged wijzigingen",
		}, &i18n.Message{
			ID:    "stashOptions",
			Other: "Stash opties",
		}, &i18n.Message{
			ID:    "notARepository",
			Other: "Fout: must be run inside a git repository",
		}, &i18n.Message{
			ID:    "jump",
			Other: "ga naar paneel",
		}, &i18n.Message{
			ID:    "DiscardPatch",
			Other: "Patch weg gooien",
		}, &i18n.Message{
			ID:    "DiscardPatchConfirm",
			Other: "Je kan alleen maar een patch bouwen van 1 commit. actueel patch weg gooien?",
		}, &i18n.Message{
			ID:    "CantPatchWhileRebasingError",
			Other: "Je kan geen patch bouwen of patch commando uitvoeren wanneer je in een merging of rebasing state zit",
		}, &i18n.Message{
			ID:    "toggleAddToPatch",
			Other: "toggle bestand inbegrepen in patch",
		}, &i18n.Message{
			ID:    "ViewPatchOptions",
			Other: "bekijk aangepaste patch opties",
		}, &i18n.Message{
			ID:    "PatchOptionsTitle",
			Other: "Patch Opties",
		}, &i18n.Message{
			ID:    "NoPatchError",
			Other: "Nog geen patch gecreëerd. Om een patch te bouwen gebruik 'space' op een commit bestand of 'enter' om een spesiefieke lijnen toe te voegen",
		}, &i18n.Message{
			ID:    "enterFile",
			Other: "enter bestand to add selecteered lines to the patch",
		}, &i18n.Message{
			ID:    "ExitLineByLineMode",
			Other: `sluit lijn-bij-lijn mode`,
		}, &i18n.Message{
			ID:    "EnterUpstream",
			Other: `Enter upstream als '<remote> <branchnaam>'`,
		}, &i18n.Message{
			ID:    "EnterUpstreamWithSlash",
			Other: `Enter upstream als '<remote>/<branchnaam>'`,
		}, &i18n.Message{
			ID:    "notTrackingRemote",
			Other: "(nog geen remote aan het volgen)",
		}, &i18n.Message{
			ID:    "ReturnToRemotesList",
			Other: `Ga terug naar remotes lijst`,
		}, &i18n.Message{
			ID:    "addNewRemote",
			Other: `voeg een nieuwe remote toe`,
		}, &i18n.Message{
			ID:    "newRemoteName",
			Other: `Nieuwe remote name:`,
		}, &i18n.Message{
			ID:    "newRemoteUrl",
			Other: `Nieuwe remote url:`,
		}, &i18n.Message{
			ID:    "editRemoteName",
			Other: `Enter updated remote naam voor {{ .remoteName }}:`,
		}, &i18n.Message{
			ID:    "editRemoteUrl",
			Other: `Enter updated remote url voor {{ .remoteName }}:`,
		}, &i18n.Message{
			ID:    "removeRemote",
			Other: `verwijder remote`,
		}, &i18n.Message{
			ID:    "removeRemotePrompt",
			Other: "Weet je zeker dat je deze remote wilt verwijderen",
		}, &i18n.Message{
			ID:    "DeleteRemoteBranch",
			Other: "Verwijder Remote Branch",
		}, &i18n.Message{
			ID:    "DeleteRemoteBranchMessage",
			Other: "Weet je zeker dat je deze remote branch wilt verwijderen",
		}, &i18n.Message{
			ID:    "setUpstream",
			Other: "stel in als upstream van uitgecheckte branch",
		}, &i18n.Message{
			ID:    "SetUpstreamTitle",
			Other: "Stel in als upstream branch",
		}, &i18n.Message{
			ID:    "SetUpstreamMessage",
			Other: "Weet je zeker dat je de upstream branch van '{{.checkedOut}}' naar '{{.selected}}' wilt zetten",
		}, &i18n.Message{
			ID:    "editRemote",
			Other: "wijzig remote",
		}, &i18n.Message{
			ID:    "tagCommit",
			Other: "tag commit",
		}, &i18n.Message{
			ID:    "TagNameTitle",
			Other: "Tag naam:",
		}, &i18n.Message{
			ID:    "deleteTag",
			Other: "verwijder tag",
		}, &i18n.Message{
			ID:    "DeleteTagTitle",
			Other: "Verwijder tag",
		}, &i18n.Message{
			ID:    "DeleteTagPrompt",
			Other: "Weet je zeker dat je '{{.tagName}}' wil verwijderen?",
		}, &i18n.Message{
			ID:    "PushTagTitle",
			Other: "remote om tag '{{.tagName}}' te pushen naar:",
		}, &i18n.Message{
			ID:    "pushTag",
			Other: "push tag",
		}, &i18n.Message{
			ID:    "createTag",
			Other: "creëer tag",
		}, &i18n.Message{
			ID:    "CreateTagTitle",
			Other: "Tag naam:",
		}, &i18n.Message{
			ID:    "fetchRemote",
			Other: "fetch remote",
		}, &i18n.Message{
			ID:    "FetchingRemoteStatus",
			Other: "remote fetchen",
		}, &i18n.Message{
			ID:    "checkoutCommit",
			Other: "checkout commit",
		}, &i18n.Message{
			ID:    "SureCheckoutThisCommit",
			Other: "Weet je zeker dat je deze commit wil uitchecken?",
		}, &i18n.Message{
			ID:    "gitFlowOptions",
			Other: "laat git-flow opties zien",
		}, &i18n.Message{
			ID:    "NotAGitFlowBranch",
			Other: "Dit lijkt geen git flow branch te zijn",
		}, &i18n.Message{
			ID:    "NewBranchNamePrompt",
			Other: "nieuwe {{.branchType}} naam:",
		}, &i18n.Message{
			ID:    "IgnoreTracked",
			Other: "Negeer tracked bestand",
		}, &i18n.Message{
			ID:    "IgnoreTrackedPrompt",
			Other: "weet je zeker dat je een getracked bestand wil negeeren?",
		}, &i18n.Message{
			ID:    "viewResetToUpstreamOptions",
			Other: "bekijk upstream reset opties",
		}, &i18n.Message{
			ID:    "nextScreenMode",
			Other: "volgende schermmode (normaal/half/groot )",
		}, &i18n.Message{
			ID:    "prevScreenMode",
			Other: "vorige schermmode",
		}, &i18n.Message{
			ID:    "startSearch",
			Other: "start met zoekken",
		}, &i18n.Message{
			ID:    "Panel",
			Other: "Paneel",
		}, &i18n.Message{
			ID:    "Keybindings",
			Other: "Sneltoetsen",
		}, &i18n.Message{
			ID:    "renameBranch",
			Other: "hernoem branch",
		}, &i18n.Message{
			ID:    "NewBranchNamePrompt",
			Other: "Noem een nieuwe branch naam",
		}, &i18n.Message{
			ID:    "RenameBranchWarning",
			Other: "Deze branch volgt een remote. Deze actie zal alleen de locale branch name wijzigen niet de naam van de remote branch. Verder gaan?",
		}, &i18n.Message{
			ID:    "openMenu",
			Other: "open menu",
		}, &i18n.Message{
			ID:    "closeMenu",
			Other: "sluit menu",
		}, &i18n.Message{
			ID:    "resetCherryPick",
			Other: "reset cherry-picked (gecopieerde) commits selectie",
		}, &i18n.Message{
			ID:    "nextTab",
			Other: "volgende tab",
		}, &i18n.Message{
			ID:    "prevTab",
			Other: "vorige tab",
		}, &i18n.Message{
			ID:    "cantUndoWhileRebasing",
			Other: "Kan niet ongedaan maken terwijl je aan het rebasen bent",
		}, &i18n.Message{
			ID:    "cantRedoWhileRebasing",
			Other: "Kan niet opnieuw doen (redo) terwijl je aan het rebasen bent",
		}, &i18n.Message{
			ID:    "MustStashWarning",
			Other: "Een patch in de index stoppen verijst stashen en onstashen van je wijzigingen. Als iets verkeert gaat kan je je bestanden terug vinden in de stash. Verder gaan?",
		}, &i18n.Message{
			ID:    "MustStashTitle",
			Other: "Moet stashen",
		}, &i18n.Message{
			ID:    "ConfirmationTitle",
			Other: "Bevestigings Paneel",
		}, &i18n.Message{
			ID:    "prevPage",
			Other: "vorige pagina",
		}, &i18n.Message{
			ID:    "nextPage",
			Other: "volgende pagina",
		}, &i18n.Message{
			ID:    "gotoTop",
			Other: "scroll naar boven",
		}, &i18n.Message{
			ID:    "gotoBottom",
			Other: "scroll naar beneden",
		}, &i18n.Message{
			ID:    "filteringBy",
			Other: "filteren bij",
		}, &i18n.Message{
			ID:    "(reset)",
			Other: "(reset)",
		}, &i18n.Message{
			ID:    "openScopingMenu",
			Other: "bekijk scoping opties",
		}, &i18n.Message{
			ID:    "filterBy",
			Other: "filter bij",
		}, &i18n.Message{
			ID:    "exitFilterMode",
			Other: "stop met filteren bij pad",
		}, &i18n.Message{
			ID:    "filterPathOption",
			Other: "vulin pad om op te filteren",
		}, &i18n.Message{
			ID:    "enterFileName",
			Other: "vulin path:",
		}, &i18n.Message{
			ID:    "FilteringMenuTitle",
			Other: "Filteren",
		}, &i18n.Message{
			ID:    "MustExitFilterModeTitle",
			Other: "Command niet beschikbaar",
		}, &i18n.Message{
			ID:    "MustExitFilterModePrompt",
			Other: "Command niet beschikbaar in filter mode. Sluit filter mode?",
		}, &i18n.Message{
			ID:    "diff",
			Other: "diff",
		}, &i18n.Message{
			ID:    "enterRefToDiff",
			Other: "vulin ref to diff",
		}, &i18n.Message{
			ID:    "enteRefName",
			Other: "vulin ref:",
		}, &i18n.Message{
			ID:    "exitDiffMode",
			Other: "sluit diff mode",
		}, &i18n.Message{
			ID:    "DiffingMenuTitle",
			Other: "Diffen",
		}, &i18n.Message{
			ID:    "swapDiff",
			Other: "keer diff richting om",
		}, &i18n.Message{
			ID:    "openDiffingMenu",
			Other: "open diff menu",
		}, &i18n.Message{
			ID:    "showingGitDiff",
			Other: "laat output zien voor:",
		}, &i18n.Message{
			ID:    "copyCommitShaToClipboard",
			Other: "copieer commit SHA naar clipboard",
		}, &i18n.Message{
			ID:    "copyBranchNameToClipboard",
			Other: "copieer branch name naar clipboard",
		}, &i18n.Message{
			ID:    "commitPrefixPatternError",
			Other: "Fout in commitPrefix patroon",
		}, &i18n.Message{
			ID:    "NoFilesStagedTitle",
			Other: "geen bestanden gestaged",
		}, &i18n.Message{
			ID:    "NoFilesStagedPrompt",
			Other: "Je hebt geen bestanden gestaged. Commit alle bestanden?",
		}, &i18n.Message{
			ID:    "BranchNotFoundTitle",
			Other: "Branch niet gevonden",
		}, &i18n.Message{
			ID:    "BranchNotFoundPrompt",
			Other: "Branch niet gevonden. Creëer een nieuwe branch genaamd",
		},
	)
}
