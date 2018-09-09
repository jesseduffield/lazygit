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
			ID:    "FilesTitle",
			Other: "Pliki",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Gałęzie",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Commity",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Schowek",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Wiadomość commita",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "commituj zmiany",
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
			ID:    "ForceDeleteBranchMessage",
			Other: "Na pewno wymusić usunięcie gałęzi {{.selectedBranchName}}?",
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
			ID:    "NoCommitsThisBranch",
			Other: "Brak commitów dla tej gałęzi",
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
			ID:    "removeFile",
			Other: `usuń jeśli nie śledzony / przełącz jeśli śledzony`,
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
			ID:    "resetHard",
			Other: `zresetuj twardo`,
		}, &i18n.Message{
			ID:    "mergeIntoCurrentBranch",
			Other: `scal do obecnej gałęzi`,
		}, &i18n.Message{
			ID:    "ConfirmQuit",
			Other: `Na pewno chcesz wyjść z programu?`,
		},
	)
}
