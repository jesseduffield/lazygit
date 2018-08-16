package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func addPolish(i18nObject *i18n.Bundle) {

	i18nObject.AddMessages(language.Polish,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Za mało miejsca do wyświetlenia paneli",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Różnice",
		}, &i18n.Message{
			ID:    "FilesTitle",
			Other: "Pliki",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Gałęzie",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Zatwierdzenia",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Skrytka",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Wiadomość zatwierdzenia",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "zatwierdź zmiany",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "nawiguj",
		}, &i18n.Message{
			ID:    "stashFiles",
			Other: "przechowaj pliki",
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
			ID:    "toggleStaged", //TODO
			Other: "toggle staged",
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
			ID:    "FileHasNoUnstagedChanges", //TODO
			Other: "File has no unstaged changes to add",
		}, &i18n.Message{
			ID:    "CannotGitAdd",
			Other: "Nie można git add --patch nieśledzonych plików",
		}, &i18n.Message{
			ID:    "CantIgnoreTrackFiles",
			Other: "Nie można zignorować nieśledzonych plików",
		}, &i18n.Message{
			ID:    "NoStagedFilesToCommit", //TODO
			Other: "There are no staged files to commit",
		}, &i18n.Message{
			ID:    "NoFilesDisplay",
			Other: "Brak pliku do wyświetlenia",
		}, &i18n.Message{
			ID:    "PullWait",
			Other: "Wciaganie...",
		}, &i18n.Message{
			ID:    "PushWait",
			Other: "Wypychanie...",
		}, &i18n.Message{
			ID:    "FileNoMergeCons",
			Other: "Ten plik nie powoduje konfliktów scalania",
		}, &i18n.Message{
			ID:    "SureResetHardHead",
			Other: "Jesteś pewny, że chcesz wykonać `reset --hard HEAD`? Możesz stracić wprowadzone zmiany",
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
			ID:    "NoBranchesThisRepo",
			Other: "Brak gałęzi dla tego repozytorium",
		}, &i18n.Message{
			ID:    "NoTrackingThisBranch",
			Other: "Brak śledzenia dla tej gałęzi",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "Nie możesz zatwierdzić bez podania wiadomości",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: zamknij, {{.keyBindConfirm}}: potwierdź",
		}, &i18n.Message{
			ID:    "SureResetThisCommit",
			Other: "Jesteś pewny, że chcesz zresetować to zatwierdzenie?",
		}, &i18n.Message{
			ID:    "ResetToCommit",
			Other: "Zresetuj, aby zatwierdzić",
		}, &i18n.Message{
			ID:    "squashDown",
			Other: "ściśnij w dół",
		}, &i18n.Message{
			ID:    "rename",
			Other: "przemianuj",
		}, &i18n.Message{
			ID:    "resetToThisCommit",
			Other: "zresetuj do tego zatwierdzenia",
		}, &i18n.Message{
			ID:    "fixupCommit", //TODO
			Other: "fixup commit",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "Brak zatwierdzeń dla tej gałęzi",
		}, &i18n.Message{
			ID:    "OnlySquashTopmostCommit",
			Other: "Można tylko ścisnąć najwyższe zatwierdzenie",
		}, &i18n.Message{
			ID:    "YouNoCommitsToSquash",
			Other: "Nie masz zatwierdzeń do ściśnięcia",
		}, &i18n.Message{
			ID:    "CantFixupWhileUnstagedChanges", //TODO
			Other: "Can't fixup while there are unstaged changes",
		}, &i18n.Message{
			ID:    "Fixup", //TODO
			Other: "Fixup",
		}, &i18n.Message{
			ID:    "SureFixupThisCommit", //TODO
			Other: "Are you sure you want to fixup this commit? The commit beneath will be squashed up into this one",
		}, &i18n.Message{
			ID:    "OnlyRenameTopCommit",
			Other: "Można przmianować tylko najwyższe zatwierdzenie",
		}, &i18n.Message{
			ID:    "RenameCommit",
			Other: "Przemianuj zatwierdzenie",
		}, &i18n.Message{
			ID:    "PotentialErrInGetselectedCommit", //TODO
			Other: "potential error in getSelected Commit (mismatched ui and state)",
		}, &i18n.Message{
			ID:    "NoCommitsThisBranch",
			Other: "Brak zatwierdzeń dla tej gałęzi",
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
			Other: "Brak pozycji w skrytce",
		}, &i18n.Message{
			ID:    "StashDrop",
			Other: "Porzuć skrytkę",
		}, &i18n.Message{
			ID:    "SureDropStashEntry",
			Other: "Jesteś pewny, że chcesz porzucić tę pozycję w skrytce?",
		}, &i18n.Message{
			ID:    "NoStashTo", //TODO
			Other: "No stash to {{.method}}",
		}, &i18n.Message{
			ID:    "NoTrackedStagedFilesStash", //TODO
			Other: "You have no tracked/staged files to stash",
		}, &i18n.Message{
			ID:    "StashChanges",
			Other: "Przechowaj zmiany",
		}, &i18n.Message{
			ID:    "IssntListOfViews",
			Other: "{{.name}} nie jest na liście widoków",
		}, &i18n.Message{
			ID:    "NoViewMachingNewLineFocusedSwitchStatement", //TODO
			Other: "No view matching newLineFocused switch statement",
		}, &i18n.Message{
			ID:    "settingPreviewsViewTo",
			Other: "ustawianie poprzedniego widoku na: {{.oldViewName}}",
		}, &i18n.Message{
			ID:    "newFocusedViewIs",
			Other: "nowy skupiony widok to {{.newFocusedView}}",
		}, &i18n.Message{
			ID:    "CantCloseConfirmationPrompt",
			Other: "Nie można zamknąć monitu potwierdzenia: {{.error}}",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Brak zmienionych plików",
		}, &i18n.Message{
			ID:    "ClearFilePanel",
			Other: "Wyczyść panel plików",
		}, &i18n.Message{
			ID:    "MergeAborted",
			Other: "Scalanie anulowane",
		},
	)
}
