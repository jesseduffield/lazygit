package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// addDutch will add all dutch translations
func addDutch(i18nObject *i18n.Bundle) {

	// add the translations
	i18nObject.AddMessages(language.Dutch,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Niet genoeg ruimte om de panelen te renderen",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Diff",
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
			ID:    "CommitMessage",
			Other: "Commit Bericht",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "Commit Veranderingen",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "navigeer",
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
			ID:    "refresh",
			Other: "verversen",
		}, &i18n.Message{
			ID:    "addPatch",
			Other: "verandering toevoegen",
		}, &i18n.Message{
			ID:    "edit",
			Other: "veranderen",
		}, &i18n.Message{
			ID:    "scroll",
			Other: "scroll",
		}, &i18n.Message{
			ID:    "abortMerge",
			Other: "samenvoegen afbreken",
		}, &i18n.Message{
			ID:    "resolveMergeConflicts",
			Other: "verhelp samenvoegen fouten",
		}, &i18n.Message{
			ID:    "checkout",
			Other: "uitchecken",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Geen Bestanden verandert",
		}, &i18n.Message{
			ID:    "NoUsernameSetErr",
			Other: `Geen gebruikersnaam ingesteld. Doe: git config --global user.name "Jou Naam"`,
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
			ID:    "PullWait",
			Other: "Pulling...",
		}, &i18n.Message{
			ID:    "PushWait",
			Other: "Pushing...",
		}, &i18n.Message{
			ID:    "FileNoMergeCons",
			Other: "Dit bestand heeft geen merge conflicten",
		}, &i18n.Message{
			ID:    "SureResetHardHead",
			Other: "Weet je het zeker dat je `reset --hard HEAD` wil uitvoeren? het kan dat je hierdoor bestanden verliest",
		}, &i18n.Message{
			ID:    "SureTo",
			Other: "Weet je het zeker dat je {{.fileName}} wilt {{.deleteVerb}} (je veranderingen zullen worden verwijdert)",
		}, &i18n.Message{
			ID:    "AlreadyCheckedOutBranch",
			Other: "Je hebt uitgecheckt op deze branch",
		}, &i18n.Message{
			ID:    "SureForceCheckout",
			Other: "Weet je zeker dat je het uitchecken wil forceren? al je locale verandering zullen worden verwijdert",
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
			Other: "Weet je zeker dat je {{.selectedBranchName}} branch wil verwijderen?",
		}, &i18n.Message{
			ID:    "CantMergeBranchIntoItself",
			Other: "Je kan niet een branch in zichzelf mergen",
		}, &i18n.Message{
			ID:    "forceCheckout",
			Other: "forceren checkout",
		}, &i18n.Message{
			ID:    "merge",
			Other: "merge",
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
			ID:    "NoBranchesThisRepo",
			Other: "Geen branches voor deze repo",
		}, &i18n.Message{
			ID:    "NoTrackingThisRepo",
			Other: "deze branch wordt niet gevolgd",
		}, &i18n.Message{
			ID:    "CommitWithoutMessageErr",
			Other: "Je kan geen commit maken zonder commit bericht",
		}, &i18n.Message{
			ID:    "CloseConfirm",
			Other: "{{.keyBindClose}}: Sluiten, {{.keyBindConfirm}}: Bevestigen",
		},
	)
}
