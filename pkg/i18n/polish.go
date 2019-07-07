package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func addPolish(i18nObject *i18n.Bundle) error {

	return i18nObject.AddMessages(language.Polish,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Za mało miejsca do wyświetlenia paneli",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Różnice",
		}, &i18n.Message{
			ID:    "LogTitle",
			Other: "Log",
		}, &i18n.Message{
			ID:    "FilesTitle",
			Other: "Pliki",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Gałęzie",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Commity",
		}, &i18n.Message{
			ID:    "CommitsDiffTitle",
			Other: "Commits (specific diff mode)",
		}, &i18n.Message{
			ID:    "CommitsDiff",
			Other: "select commit to diff with another commit",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Schowek",
		}, &i18n.Message{
			ID:    "StagingMainTitle",
			Other: `Stage Lines/Hunks`,
		}, &i18n.Message{
			ID:    "MergingMainTitle",
			Other: "Resolve merge conflicts",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Wiadomość commita",
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
			Other: "commituj zmiany",
		}, &i18n.Message{
			ID:    "AmendLastCommit",
			Other: "zmień ostatnie zatwierdzenie",
		}, &i18n.Message{
			ID:    "SureToAmend",
			Other: "Czy na pewno chcesz zmienić ostatnie zatwierdzenie? Możesz zmienić komunikat zatwierdzenia z panelu zatwierdzeń.",
		}, &i18n.Message{
			ID:    "NoCommitToAmend",
			Other: "Nie ma zobowiązania do zmiany.",
		}, &i18n.Message{
			ID:    "CommitChangesWithEditor",
			Other: "commituj zmiany używając edytora z gita",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "GlobalTitle",
			Other: "Globalne",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "nawiguj",
		}, &i18n.Message{
			ID:    "menu",
			Other: "menu",
		}, &i18n.Message{
			ID:    "execute",
			Other: "wykonaj",
		}, &i18n.Message{
			ID:    "open",
			Other: "otwórz",
		}, &i18n.Message{
			ID:    "ignore",
			Other: "ignoruj",
		}, &i18n.Message{
			ID:    "delete",
			Other: "usuń",
		}, &i18n.Message{
			ID:    "toggleStaged",
			Other: "przełącz zatwierdzenie",
		}, &i18n.Message{
			ID:    "toggleStagedAll",
			Other: "przełącz wszystkie zatwierdzenia",
		}, &i18n.Message{
			ID:    "refresh",
			Other: "odśwież",
		}, &i18n.Message{
			ID:    "addPatch",
			Other: "dodaj łatkę",
		}, &i18n.Message{
			ID:    "edit",
			Other: "edytuj",
		}, &i18n.Message{
			ID:    "scroll",
			Other: "przewiń",
		}, &i18n.Message{
			ID:    "abortMerge",
			Other: "o scalaniu",
		}, &i18n.Message{
			ID:    "resolveMergeConflicts",
			Other: "rozwiąż konflikty scalania",
		}, &i18n.Message{
			ID:    "checkout",
			Other: "przełącz",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Brak zmienionych plików",
		}, &i18n.Message{
			ID:    "FileHasNoUnstagedChanges",
			Other: "Plik nie zawiera żadnych nieopublikowanych zmian do dodania",
		}, &i18n.Message{
			ID:    "CannotGitAdd",
			Other: "Nie można git add --patch nieśledzonych plików",
		}, &i18n.Message{
			ID:    "CantIgnoreTrackFiles",
			Other: "Nie można zignorować nieśledzonych plików",
		}, &i18n.Message{
			ID:    "NoStagedFilesToCommit",
			Other: "Brak zatwierdzonych plików do commita",
		}, &i18n.Message{
			ID:    "NoFilesDisplay",
			Other: "Brak pliku do wyświetlenia",
		}, &i18n.Message{
			ID:    "PullWait",
			Other: "Wciąganie zmian...",
		}, &i18n.Message{
			ID:    "PushWait",
			Other: "Wypychanie zmian...",
		}, &i18n.Message{
			ID:    "FetchWait",
			Other: "Fetching...",
		}, &i18n.Message{
			ID:    "FileNoMergeCons",
			Other: "Ten plik nie powoduje konfliktów scalania",
		}, &i18n.Message{
			ID:    "SureTo",
			Other: "Jesteś pewny, że chcesz {{.deleteVerb}} {{.fileName}} (stracisz swoje wprowadzone zmiany)?",
		}, &i18n.Message{
			ID:    "AlreadyCheckedOutBranch",
			Other: "Już przęłączono na tą gałąź",
		}, &i18n.Message{
			ID:    "SureForceCheckout",
			Other: "Jesteś pewny, że chcesz wymusić przełączenie? Stracisz wszystkie lokalne zmiany",
		}, &i18n.Message{
			ID:    "ForceCheckoutBranch",
			Other: "Wymuś przełączenie gałęzi",
		}, &i18n.Message{
			ID:    "BranchName",
			Other: "Nazwa gałęzi",
		}, &i18n.Message{
			ID:    "NewBranchNameBranchOff",
			Other: "Nazwa nowej gałęzi (gałąź na bazie {{.branchName}})",
		}, &i18n.Message{
			ID:    "CantDeleteCheckOutBranch",
			Other: "Nie możesz usunąć obecnej przełączonej gałęzi!",
		}, &i18n.Message{
			ID:    "DeleteBranch",
			Other: "Usuń gałąź",
		}, &i18n.Message{
			ID:    "DeleteBranchMessage",
			Other: "Jesteś pewien, że chcesz usunąć gałąź {{.selectedBranchName}} ?",
		}, &i18n.Message{
			ID:    "ForceDeleteBranchMessage",
			Other: "Na pewno wymusić usunięcie gałęzi {{.selectedBranchName}}?",
		}, &i18n.Message{
			ID:    "rebaseBranch",
			Other: "rebase branch",
		}, &i18n.Message{
			ID:    "CantRebaseOntoSelf",
			Other: "You cannot rebase a branch onto itself",
		}, &i18n.Message{
			ID:    "CantMergeBranchIntoItself",
			Other: "Nie możesz scalić gałęzi do samej siebie",
		}, &i18n.Message{
			ID:    "forceCheckout",
			Other: "wymuś przełączenie",
		}, &i18n.Message{
			ID:    "merge",
			Other: "scal",
		}, &i18n.Message{
			ID:    "checkoutByName",
			Other: "przełącz używając nazwy",
		}, &i18n.Message{
			ID:    "newBranch",
			Other: "nowa gałąź",
		}, &i18n.Message{
			ID:    "deleteBranch",
			Other: "usuń gałąź",
		}, &i18n.Message{
			ID:    "forceDeleteBranch",
			Other: "usuń gałąź (wymuś)",
		}, &i18n.Message{
			ID:    "NoBranchesThisRepo",
			Other: "Brak gałęzi dla tego repozytorium",
		}, &i18n.Message{
			ID:    "NoTrackingThisBranch",
			Other: "Brak śledzenia dla tej gałęzi",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "Nie możesz commitować bez podania wiadomości",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: zamknij, {{.keyBindConfirm}}: potwierdź",
		}, &i18n.Message{
			ID:    "close",
			Other: "zamknij",
		}, &i18n.Message{
			ID:    "SureResetThisCommit",
			Other: "Jesteś pewny, że chcesz zresetować ten commit?",
		}, &i18n.Message{
			ID:    "ResetToCommit",
			Other: "Zresetuj, aby commitować",
		}, &i18n.Message{
			ID:    "squashDown",
			Other: "ściśnij w dół",
		}, &i18n.Message{
			ID:    "rename",
			Other: "przemianuj",
		}, &i18n.Message{
			ID:    "resetToThisCommit",
			Other: "zresetuj do tego commita",
		}, &i18n.Message{
			ID:    "fixupCommit",
			Other: "napraw commit",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "Brak commitów dla tej gałęzi",
		}, &i18n.Message{
			ID:    "OnlySquashTopmostCommit",
			Other: "Można tylko ścisnąć najwyższy commit",
		}, &i18n.Message{
			ID:    "YouNoCommitsToSquash",
			Other: "Nie masz commitów do ściśnięcia",
		}, &i18n.Message{
			ID:    "CantFixupWhileUnstagedChanges",
			Other: "Nie można wykonać naprawy, kiedy istnieją niezatwierdzone zmiany",
		}, &i18n.Message{
			ID:    "Fixup",
			Other: "Napraw",
		}, &i18n.Message{
			ID:    "SureFixupThisCommit",
			Other: "Jesteś pewny, ze chcesz naprawić ten commit? Commit poniżej zostanie ściśnięty w górę wraz z tym",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Można przmianować tylko najwyższy commit",
		}, &i18n.Message{
			ID:    "renameCommit",
			Other: "przemianuj commit",
		}, &i18n.Message{
			ID:    "renameCommitEditor",
			Other: "przemianuj commit w edytorze",
		}, &i18n.Message{
			ID:    "PotentialErrInGetselectedCommit",
			Other: "potencjalny błąd w getSelected Commit (niedopasowane ui i stan)",
		}, &i18n.Message{
			ID:    "Error",
			Other: "Błąd",
		}, &i18n.Message{
			ID:    "resizingPopupPanel",
			Other: "skalowanie wyskakującego panelu",
		}, &i18n.Message{
			ID:    "RunningSubprocess",
			Other: "uruchomiony podproces",
		}, &i18n.Message{
			ID:    "selectHunk",
			Other: "wybierz kawałek",
		}, &i18n.Message{
			ID:    "navigateConflicts",
			Other: "nawiguj konflikty",
		}, &i18n.Message{
			ID:    "pickHunk",
			Other: "wybierz kawałek",
		}, &i18n.Message{
			ID:    "pickBothHunks",
			Other: "wybierz oba kawałki",
		}, &i18n.Message{
			ID:    "undo",
			Other: "cofnij",
		}, &i18n.Message{
			ID:    "pop",
			Other: "wyciągnij",
		}, &i18n.Message{
			ID:    "drop",
			Other: "porzuć",
		}, &i18n.Message{
			ID:    "apply",
			Other: "zastosuj",
		}, &i18n.Message{
			ID:    "NoStashEntries",
			Other: "Brak pozycji w schowku",
		}, &i18n.Message{
			ID:    "StashDrop",
			Other: "Porzuć schowek",
		}, &i18n.Message{
			ID:    "SureDropStashEntry",
			Other: "Jesteś pewny, że chcesz porzucić tę pozycję w schowku?",
		}, &i18n.Message{
			ID:    "NoStashTo",
			Other: "Brak schowka dla {{.method}}",
		}, &i18n.Message{
			ID:    "NoTrackedStagedFilesStash",
			Other: "Nie masz śledzonych/zatwierdzonych plików do przechowania",
		}, &i18n.Message{
			ID:    "StashChanges",
			Other: "Przechowaj zmiany",
		}, &i18n.Message{
			ID:    "IssntListOfViews",
			Other: "{{.name}} nie jest na liście widoków",
		}, &i18n.Message{
			ID:    "NoViewMachingNewLineFocusedSwitchStatement",
			Other: "Brak widoku pasującego do instrukcji przełączania newLineFocused",
		}, &i18n.Message{
			ID:    "newFocusedViewIs",
			Other: "nowy skupiony widok to {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "CantCloseConfirmationPrompt",
			Other: "Nie można zamknąć monitu potwierdzenia: {{.error}}",
		}, &i18n.Message{
			ID:    "MergeAborted",
			Other: "Scalanie anulowane",
		}, &i18n.Message{
			ID:    "OpenConfig",
			Other: "otwórz plik konfiguracyjny",
		}, &i18n.Message{
			ID:    "EditConfig",
			Other: "edytuj plik konfiguracyjny",
		}, &i18n.Message{
			ID:    "ForcePush",
			Other: "Wymuś wypchnięcie",
		}, &i18n.Message{
			ID:    "ForcePushPrompt",
			Other: "Twoja gałąź rozeszła się z gałęzią zdalną. Wciśnij 'esc' aby anulować lub 'enter' aby wymusić wypchnięcie.",
		}, &i18n.Message{
			ID:    "checkForUpdate",
			Other: "sprawdź aktualizacje",
		}, &i18n.Message{
			ID:    "CheckingForUpdates",
			Other: "Sprawdzanie aktualizacji...",
		}, &i18n.Message{
			ID:    "OnLatestVersionErr",
			Other: "Już posiadasz najnowszą wersję",
		}, &i18n.Message{
			ID:    "MajorVersionErr",
			Other: "Nowa wersja ({{.newVersion}}) posiada niekompatybilne zmiany w porównaniu do obecnej wersji ({{.currentVersion}})",
		}, &i18n.Message{
			ID:    "CouldNotFindBinaryErr",
			Other: "Nie można znaleźć pliku binarnego w {{.url}}",
		}, &i18n.Message{
			ID:    "AnonymousReportingTitle",
			Other: "Help make lazygit better",
		}, &i18n.Message{
			ID:    "AnonymousReportingPrompt",
			Other: "Włączyć anonimowe raportowanie błędów w celu pomocy w usprawnianiu lazygita (enter/esc)?",
		}, &i18n.Message{
			ID:    "editFile",
			Other: `edytuj plik`,
		}, &i18n.Message{
			ID:    "openFile",
			Other: `otwórz plik`,
		}, &i18n.Message{
			ID:    "ignoreFile",
			Other: `dodaj do .gitignore`,
		}, &i18n.Message{
			ID:    "refreshFiles",
			Other: `odśwież pliki`,
		}, &i18n.Message{
			ID:    "mergeIntoCurrentBranch",
			Other: `scal do obecnej gałęzi`,
		}, &i18n.Message{
			ID:    "ConfirmQuit",
			Other: `Na pewno chcesz wyjść z programu?`,
		}, &i18n.Message{
			ID:    "UnsupportedGitService",
			Other: `Nieobsługiwana usługa git`,
		}, &i18n.Message{
			ID:    "createPullRequest",
			Other: `utwórz żądanie wyciągnięcia`,
		}, &i18n.Message{
			ID:    "NoBranchOnRemote",
			Other: `Ta gałąź nie istnieje na zdalnym. Najpierw musisz go odepchnąć na odległość.`,
		}, &i18n.Message{
			ID:    "fetch",
			Other: `fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchTitle",
			Other: `No automatic git fetch`,
		}, &i18n.Message{
			ID:    "NoAutomaticGitFetchBody",
			Other: `Lazygit can't use "git fetch" in a private repo use f in the branches panel to run "git fetch" manually`,
		}, &i18n.Message{
			ID:    "StageLines",
			Other: `zatwierdź pojedyncze linie`,
		}, &i18n.Message{
			ID:    "FileStagingRequirements",
			Other: `Można tylko zatwierdzić pojedyncze linie dla śledzonych plików z niezatwierdzonymi zmianami`,
		}, &i18n.Message{
			ID:    "StagingTitle",
			Other: `Zatwierdzanie`,
		}, &i18n.Message{
			ID:    "StageHunk",
			Other: `zatwierdź kawałek`,
		}, &i18n.Message{
			ID:    "StageLine",
			Other: `zatwierdź linię`,
		}, &i18n.Message{
			ID:    "EscapeStaging",
			Other: `wróć do panelu plików`,
		}, &i18n.Message{
			ID:    "CantFindHunks",
			Other: `Nie można znaleźć żadnych kawałków w tej łatce`,
		}, &i18n.Message{
			ID:    "CantFindHunk",
			Other: `Nie można znaleźć kawałka`,
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
			ID:    "MainTitle",
			Other: "Main",
		}, &i18n.Message{
			ID:    "NormalTitle",
			Other: "Normal",
		}, &i18n.Message{
			ID:    "softReset",
			Other: "soft reset",
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
			ID:    "CommitFiles",
			Other: "Commit files",
		}, &i18n.Message{
			ID:    "viewCommitFiles",
			Other: "view commit's files",
		}, &i18n.Message{
			ID:    "CommitFilesTitle",
			Other: "Commit files",
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
			Other: "przechowaj pliki",
		}, &i18n.Message{
			ID:    "stashStagedChanges",
			Other: "stash staged changes",
		}, &i18n.Message{
			ID:    "stashOptions",
			Other: "Stash options",
		}, &i18n.Message{
			ID:    "notARepository",
			Other: "Error: must be run inside a git repository",
		},
	)
}
